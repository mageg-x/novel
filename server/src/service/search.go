package service

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/yanyiwu/gojieba"
	"gorm.io/gorm"

	"github.com/mageg-x/novel/src/assets"
	"github.com/mageg-x/novel/src/model"
)

// 全局搜索服务实例
var (
	jieba    *gojieba.Jieba
	SService = newSearchService()
	// 使用读写锁代替互斥锁，允许多个读操作同时进行
	mu = sync.RWMutex{}
)

// SearchService 搜索服务结构体
type SearchService struct {
	bookIndex    bleve.Index
	userIndex    bleve.Index
	commentIndex bleve.Index
	indexDir     string
	// 索引状态字段
	indexReady bool
	indexMutex sync.RWMutex
}

type JiebaAnalyzer struct{}

func (a *JiebaAnalyzer) Analyze(input []byte) analysis.TokenStream {
	words := jieba.Cut(string(input), true)

	tokens := make(analysis.TokenStream, 0, len(words))
	position := 0
	for _, word := range words {
		token := &analysis.Token{
			Term:     []byte(word),
			Start:    position,
			End:      position + len(word),
			Position: position + 1,
			Type:     analysis.Ideographic,
		}
		tokens = append(tokens, token)
		position += len(word)
	}
	return tokens
}
func init() {
	dictDir := assets.DictTempDir
	logger.Infof("正在初始化分词器 dictDir: %s", dictDir)
	jieba = gojieba.NewJieba(
		filepath.Join(dictDir, "jieba.dict.utf8"),
		filepath.Join(dictDir, "hmm_model.utf8"),
		filepath.Join(dictDir, "user.dict.utf8"),
		filepath.Join(dictDir, "idf.utf8"),
		filepath.Join(dictDir, "stop_words.utf8"),
	)

	// 注册到bleve
	registry.RegisterAnalyzer("gojieba", func(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
		return &JiebaAnalyzer{}, nil
	})
}

func newSearchService() *SearchService {
	indexDir := filepath.Join(DataDir, "search")
	return &SearchService{indexDir: indexDir}
}

// ==================== 初始化相关 ====================

// InitSearch 初始化搜索索引
func (s *SearchService) InitSearch() error {
	os.MkdirAll(s.indexDir, 0755)

	go func() {
		mu.Lock()
		defer mu.Unlock()
		logger.Infof("正在初始化书籍索引索引目录: %s", s.indexDir)
		if err := s.initBookIndex(); err != nil {
			logger.Errorf("初始化书籍索引失败: %w", err)
			return
		}

		logger.Infof("正在初始化用户索引索引目录: %s", s.indexDir)
		if err := s.initUserIndex(); err != nil {
			logger.Errorf("初始化用户索引失败: %w", err)
			return
		}

		logger.Infof("正在初始化评论索引索引目录: %s", s.indexDir)
		if err := s.initCommentIndex(); err != nil {
			logger.Errorf("初始化评论索引失败: %w", err)
			return
		}

		logger.Infof("正在初始化索引...")

		// 检查索引是否为空，如果不为空则跳过完整重建
		if s.isIndexEmpty() {
			logger.Infof("检测到空索引，开始建立完整索引")
			if err := s.indexExistingData(); err != nil {
				logger.Errorf("索引现有数据失败: %w", err)
				return
			}
		} else {
			logger.Infof("索引已存在数据，跳过完整重建")
			// 执行增量更新
			if err := s.updateIndexIncrementally(); err != nil {
				logger.Errorf("增量更新索引失败: %w", err)
			}
		}

		// 设置索引状态为就绪
		s.setIndexReady(true)
		logger.Infof("搜索索引初始化成功，索引已就绪")

	}()

	go s.ListenEvents()
	return nil
}

// ListenEvents 监听模型事件并更新索引
func (s *SearchService) ListenEvents() {
	for event := range model.EventChannel {
		switch event.Type {
		case "created", "updated":
			s.handleModelUpdate(event)
		case "deleted":
			s.handleModelDelete(event)
		}
	}
}

func (s *SearchService) handleModelUpdate(event model.ModelEvent) {
	// 使用写锁进行模型更新
	if !mu.TryLock() {
		return
	}
	defer mu.Unlock()

	switch event.ModelType {
	case "user":
		var user model.User
		if err := DB.First(&user, event.ModelID).Error; nil == err {
			s.IndexUser(user)
		}
	case "book":
		var book model.Book
		if err := DB.First(&book, event.ModelID).Error; nil == err {
			s.IndexBook(book)
		}
	case "comment":
		var comment model.Comment
		if err := DB.Preload("User").Preload("Book").First(&comment, event.ModelID).Error; nil == err {
			s.IndexComment(comment)
		}
	}
}

func (s *SearchService) handleModelDelete(event model.ModelEvent) {
	// 使用写锁进行模型删除
	if !mu.TryLock() {
		return
	}
	defer mu.Unlock()
	switch event.ModelType {
	case "user":
		s.DeleteUserIndex(event.ModelID)
	case "book":
		s.DeleteBookIndex(event.ModelID)
	case "comment":
		s.DeleteCommentIndex(event.ModelID)
	}
}

// ==================== 索引初始化 ====================

func (s *SearchService) initBookIndex() error {
	bookIndexPath := filepath.Join(s.indexDir, "books.bleve")
	var err error
	s.bookIndex, err = s.openOrCreateIndex(bookIndexPath, s.createBookMapping)
	return err
}

func (s *SearchService) initUserIndex() error {
	userIndexPath := filepath.Join(s.indexDir, "users.bleve")
	var err error
	s.userIndex, err = s.openOrCreateIndex(userIndexPath, s.createUserMapping)
	return err
}

func (s *SearchService) initCommentIndex() error {
	commentIndexPath := filepath.Join(s.indexDir, "comments.bleve")
	var err error
	s.commentIndex, err = s.openOrCreateIndex(commentIndexPath, s.createCommentMapping)
	return err
}

func (s *SearchService) openOrCreateIndex(indexPath string, mappingFunc func() *mapping.IndexMappingImpl) (bleve.Index, error) {
	if _, err := os.Stat(indexPath); err == nil {
		return bleve.Open(indexPath)
	}
	return bleve.New(indexPath, mappingFunc())
}

// ==================== 索引映射配置 ====================

func (s *SearchService) createBookMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()
	bookMapping := bleve.NewDocumentMapping()

	// ID 字段
	idField := bleve.NewNumericFieldMapping()
	idField.Store = true
	idField.Index = true
	bookMapping.AddFieldMappingsAt("ID", idField)

	// 标题字段
	titleField := bleve.NewTextFieldMapping()
	titleField.Store = true
	titleField.Index = true
	titleField.Analyzer = "gojieba"
	bookMapping.AddFieldMappingsAt("Title", titleField)

	// 作者字段
	authorField := bleve.NewTextFieldMapping()
	authorField.Store = true
	authorField.Index = true
	authorField.Analyzer = "gojieba"
	bookMapping.AddFieldMappingsAt("Author", authorField)

	// 描述字段
	descField := bleve.NewTextFieldMapping()
	descField.Store = true
	descField.Index = true
	descField.Analyzer = "gojieba"
	bookMapping.AddFieldMappingsAt("Description", descField)

	indexMapping.DefaultMapping = bookMapping
	return indexMapping
}

func (s *SearchService) createUserMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()
	userMapping := bleve.NewDocumentMapping()

	// ID 字段
	idField := bleve.NewNumericFieldMapping()
	idField.Store = true
	idField.Index = true
	userMapping.AddFieldMappingsAt("ID", idField)

	// 用户名字段
	usernameField := bleve.NewTextFieldMapping()
	usernameField.Store = true
	usernameField.Index = true
	usernameField.Analyzer = "gojieba"
	userMapping.AddFieldMappingsAt("Username", usernameField)

	// 电话字段 - 支持部分匹配
	phoneField := bleve.NewTextFieldMapping()
	phoneField.Store = true
	phoneField.Index = true
	phoneField.Analyzer = standard.Name
	userMapping.AddFieldMappingsAt("Phone", phoneField)

	indexMapping.DefaultMapping = userMapping
	return indexMapping
}

func (s *SearchService) createCommentMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()
	commentMapping := bleve.NewDocumentMapping()

	// 书籍标题字段
	bookTitleField := bleve.NewTextFieldMapping()
	bookTitleField.Store = true
	bookTitleField.Index = true
	bookTitleField.Analyzer = "gojieba"
	commentMapping.AddFieldMappingsAt("BookTitle", bookTitleField)

	// 内容字段
	contentField := bleve.NewTextFieldMapping()
	contentField.Store = true
	contentField.Index = true
	contentField.Analyzer = "gojieba"
	commentMapping.AddFieldMappingsAt("Content", contentField)

	indexMapping.DefaultMapping = commentMapping
	return indexMapping
}

// ==================== 索引操作 ====================

func (s *SearchService) indexExistingData() error {
	// 分段加载书籍数据，每次加载1000条
	if err := s.batchIndexBooks(); err != nil {
		logger.Errorf("批量索引书籍数据失败: %w", err)
	}

	// 分段加载用户数据，每次加载1000条
	if err := s.batchIndexUsers(); err != nil {
		logger.Errorf("批量索引用户数据失败: %w", err)
	}

	// 分段加载评论数据，每次加载500条
	if err := s.batchIndexComments(); err != nil {
		logger.Errorf("批量索引评论数据失败: %w", err)
	}
	return nil
}

// 批量索引书籍数据
func (s *SearchService) batchIndexBooks() error {
	var total int64
	// 记录计数查询耗时
	countStart := time.Now()
	if err := DB.Model(&model.Book{}).Count(&total).Error; err != nil {
		return err
	}
	countDuration := time.Since(countStart)
	logger.Infof("书籍总数查询耗时: %v, 总数: %d", countDuration, total)

	// 批量索引优化：使用bleve的Batch API
	var lastID uint = 0
	indexedCount := 0

	// 动态批次大小调整
	// 根据总数据量设置初始批次大小
	batchSize := 1000
	if total > 50000 {
		batchSize = 2000 // 大数据集使用更大批次
	} else if total < 5000 {
		batchSize = 500 // 小数据集使用较小批次
	}

	// 记录总索引耗时
	totalStart := time.Now()

	// 自适应批次调整参数
	optimalDuration := 5 * time.Second   // 目标批次处理时间
	timeoutThreshold := 30 * time.Second // 超时阈值

	// 性能监控：记录每个阶段的耗时统计
	var queryTimes []time.Duration
	var prepareTimes []time.Duration
	var batchTimes []time.Duration
	var commitTimes []time.Duration

	// 使用基于ID的游标分页替代OFFSET分页，只查询必要字段
	for indexedCount < int(total) {
		var books []model.Book
		// 记录SQL查询耗时
		sqlStart := time.Now()

		// 使用WHERE ID > lastID的方式代替OFFSET，添加ORDER BY以保证结果一致性
		query := DB.Model(&model.Book{}).Where("id > ?", lastID).Order("id ASC").Limit(batchSize)

		// 只选择必要字段，减少数据传输和处理开销
		query = query.Select("id", "title", "author", "description")

		if err := query.Find(&books).Error; err != nil {
			return err
		}

		sqlDuration := time.Since(sqlStart)
		queryTimes = append(queryTimes, sqlDuration)
		logger.Infof("书籍数据查询耗时: %v, 查询条数: %d", sqlDuration, len(books))

		if len(books) == 0 {
			break // 没有更多数据
		}

		// 更新lastID用于下一页查询
		lastID = books[len(books)-1].ID

		// 记录索引操作耗时 - 修正计时逻辑，放在实际索引操作前
		indexStart := time.Now()

		// 详细性能监控：分阶段记录耗时
		prepareStart := time.Now()
		// 使用bleve的Batch API来批量索引，显著提升性能
		batch := s.bookIndex.NewBatch()

		// 监控内存使用情况
		for i, book := range books {
			bookID := fmt.Sprintf("%d", book.ID)
			if err := batch.Index(bookID, book); err != nil {
				logger.Errorf("批量索引添加书籍 %d 失败: %w", book.ID, err)
			}

			// 每100条记录一次，避免过多日志
			if (i+1)%100 == 0 || i+1 == len(books) {
				if i+1%500 == 0 {
					// 避免频繁获取内存统计
					var memStats runtime.MemStats
					runtime.ReadMemStats(&memStats)
					logger.Debugf("内存使用情况 - 已处理 %d/%d 条: 分配内存: %.2f MB, GC次数: %d",
						i+1, len(books), float64(memStats.Alloc)/1024/1024, memStats.NumGC)
				}
			}
		}

		prepareDuration := time.Since(prepareStart)
		prepareTimes = append(prepareTimes, prepareDuration)
		logger.Debugf("批处理准备耗时: %v, 准备 %d 条数据", prepareDuration, len(books))

		// 记录批处理提交耗时
		commitStart := time.Now()
		// 一次性提交整个批次，减少IO操作和锁竞争
		if err := s.bookIndex.Batch(batch); err != nil {
			logger.Errorf("批量索引提交失败: %w", err)
			return err
		}
		commitDuration := time.Since(commitStart)
		commitTimes = append(commitTimes, commitDuration)
		logger.Debugf("批处理提交耗时: %v, 提交 %d 条数据", commitDuration, len(books))

		// 清理批次对象，帮助GC
		batch.Reset()

		indexDuration := time.Since(indexStart)
		batchTimes = append(batchTimes, indexDuration)

		indexedCount += len(books)

		// 动态调整批次大小
		// 如果处理时间过长，减少批次大小
		if indexDuration > timeoutThreshold {
			newBatchSize := int(float64(batchSize) * 0.7) // 减少30%
			if newBatchSize < 100 {
				newBatchSize = 100 // 设置最小批次大小
			}
			logger.Infof("索引处理超时，调整批次大小从 %d 到 %d", batchSize, newBatchSize)
			batchSize = newBatchSize
		} else if indexDuration < optimalDuration && batchSize < 5000 {
			// 如果处理时间很短，增加批次大小以提高吞吐量
			newBatchSize := int(float64(batchSize) * 1.2) // 增加20%
			if newBatchSize > 5000 {
				newBatchSize = 5000 // 设置最大批次大小
			}
			logger.Infof("索引处理速度快，调整批次大小从 %d 到 %d", batchSize, newBatchSize)
			batchSize = newBatchSize
		}

		// 计算批次处理的各阶段耗时比例，帮助定位瓶颈
		queryPercent := float64(sqlDuration) / float64(indexDuration) * 100
		preparePercent := float64(prepareDuration) / float64(indexDuration) * 100
		commitPercent := float64(commitDuration) / float64(indexDuration) * 100

		logger.Infof("已索引书籍数据进度: %d/%d, 索引耗时: %v, 平均每秒索引: %.2f 条 [查询: %.1f%%, 准备: %.1f%%, 提交: %.1f%%]",
			indexedCount, total, indexDuration, float64(len(books))/indexDuration.Seconds(),
			queryPercent, preparePercent, commitPercent)

		// 每处理1000条数据，强制GC以避免内存占用过高
		if indexedCount%1000 == 0 {
			runtime.GC()
			logger.Debugf("处理完成 %d 条数据，执行GC清理内存", indexedCount)
		}
	}

	// 汇总统计信息
	totalDuration := time.Since(totalStart)
	logger.Infof("书籍索引总耗时: %v, 总处理条数: %d, 平均每秒处理: %.2f 条",
		totalDuration, indexedCount, float64(indexedCount)/totalDuration.Seconds())

	// 计算各阶段的平均耗时
	if len(queryTimes) > 0 {
		var avgQuery time.Duration
		for _, t := range queryTimes {
			avgQuery += t
		}
		avgQuery /= time.Duration(len(queryTimes))
		logger.Infof("平均SQL查询耗时: %v", avgQuery)
	}

	if len(prepareTimes) > 0 {
		var avgPrepare time.Duration
		for _, t := range prepareTimes {
			avgPrepare += t
		}
		avgPrepare /= time.Duration(len(prepareTimes))
		logger.Infof("平均批处理准备耗时: %v", avgPrepare)
	}

	if len(commitTimes) > 0 {
		var avgCommit time.Duration
		for _, t := range commitTimes {
			avgCommit += t
		}
		avgCommit /= time.Duration(len(commitTimes))
		logger.Infof("平均批处理提交耗时: %v", avgCommit)
	}

	// 清理数据，帮助GC
	runtime.GC()
	logger.Infof("书籍索引完成，执行最终GC")
	// 记录总耗时
	totalDuration = time.Since(totalStart)
	logger.Infof("书籍批量索引总耗时: %v, 平均每秒索引: %.2f 条",
		totalDuration, float64(total)/totalDuration.Seconds())
	return nil
}

// 批量索引用户数据
func (s *SearchService) batchIndexUsers() error {
	var total int64
	// 记录计数查询耗时
	countStart := time.Now()
	if err := DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return err
	}
	countDuration := time.Since(countStart)
	logger.Infof("用户总数查询耗时: %v, 总数: %d", countDuration, total)

	// 批量索引优化：使用bleve的Batch API
	var lastID uint = 0
	indexedCount := 0

	// 动态批次大小调整
	// 根据总数据量设置初始批次大小
	batchSize := 1000
	if total > 50000 {
		batchSize = 2000 // 大数据集使用更大批次
	} else if total < 5000 {
		batchSize = 500 // 小数据集使用较小批次
	}

	// 记录总索引耗时
	totalStart := time.Now()

	// 自适应批次调整参数
	optimalDuration := 5 * time.Second   // 目标批次处理时间
	timeoutThreshold := 30 * time.Second // 超时阈值

	// 使用基于ID的游标分页替代OFFSET分页，只查询必要字段
	for indexedCount < int(total) {
		var users []model.User
		// 记录SQL查询耗时
		sqlStart := time.Now()

		// 使用WHERE ID > lastID的方式代替OFFSET，添加ORDER BY以保证结果一致性
		query := DB.Model(&model.User{}).Where("id > ?", lastID).Order("id ASC").Limit(batchSize)

		// 只选择必要字段，减少数据传输和处理开销
		query = query.Select("id", "username", "phone")

		if err := query.Find(&users).Error; err != nil {
			return err
		}

		sqlDuration := time.Since(sqlStart)
		logger.Infof("用户数据查询耗时: %v, 查询条数: %d", sqlDuration, len(users))

		if len(users) == 0 {
			break // 没有更多数据
		}

		// 更新lastID用于下一页查询
		lastID = users[len(users)-1].ID

		// 记录索引操作耗时 - 修正计时逻辑，放在实际索引操作前
		indexStart := time.Now()

		// 使用bleve的Batch API来批量索引，显著提升性能
		batch := s.userIndex.NewBatch()
		for _, user := range users {
			userID := fmt.Sprintf("%d", user.ID)
			if err := batch.Index(userID, user); err != nil {
				logger.Errorf("批量索引添加用户 %d 失败: %w", user.ID, err)
			}
		}

		// 一次性提交整个批次，减少IO操作和锁竞争
		if err := s.userIndex.Batch(batch); err != nil {
			logger.Errorf("批量索引提交失败: %w", err)
			return err
		}

		indexDuration := time.Since(indexStart)

		indexedCount += len(users)

		// 动态调整批次大小
		// 如果处理时间过长，减少批次大小
		if indexDuration > timeoutThreshold {
			newBatchSize := int(float64(batchSize) * 0.7) // 减少30%
			if newBatchSize < 100 {
				newBatchSize = 100 // 设置最小批次大小
			}
			logger.Infof("索引处理超时，调整批次大小从 %d 到 %d", batchSize, newBatchSize)
			batchSize = newBatchSize
		} else if indexDuration < optimalDuration && batchSize < 5000 {
			// 如果处理时间很短，增加批次大小以提高吞吐量
			newBatchSize := int(float64(batchSize) * 1.2) // 增加20%
			if newBatchSize > 5000 {
				newBatchSize = 5000 // 设置最大批次大小
			}
			logger.Infof("索引处理速度快，调整批次大小从 %d 到 %d", batchSize, newBatchSize)
			batchSize = newBatchSize
		}

		logger.Infof("已索引用户数据进度: %d/%d, 索引耗时: %v, 平均每秒索引: %.2f 条",
			indexedCount, total, indexDuration, float64(len(users))/indexDuration.Seconds())
	}
	// 记录总耗时
	totalDuration := time.Since(totalStart)
	logger.Infof("用户批量索引总耗时: %v, 平均每秒索引: %.2f 条",
		totalDuration, float64(total)/totalDuration.Seconds())
	return nil
}

// 批量索引评论数据
func (s *SearchService) batchIndexComments() error {
	var total int64
	// 记录计数查询耗时
	countStart := time.Now()
	if err := DB.Model(&model.Comment{}).Count(&total).Error; err != nil {
		return err
	}
	countDuration := time.Since(countStart)
	logger.Infof("评论总数查询耗时: %v, 总数: %d", countDuration, total)

	// 批量索引优化：使用bleve的Batch API
	var lastID uint = 0
	indexedCount := 0

	// 动态批次大小调整
	// 根据总数据量设置初始批次大小，评论数据初始值设为500（与原代码一致）
	batchSize := 500
	if total > 50000 {
		batchSize = 1000 // 大数据集使用更大批次
	} else if total < 5000 {
		batchSize = 300 // 小数据集使用较小批次
	}

	// 记录总索引耗时
	totalStart := time.Now()

	// 自适应批次调整参数
	optimalDuration := 5 * time.Second   // 目标批次处理时间
	timeoutThreshold := 30 * time.Second // 超时阈值

	// 使用基于ID的游标分页替代OFFSET分页
	for indexedCount < int(total) {
		var comments []model.Comment
		// 记录SQL查询耗时
		sqlStart := time.Now()

		// 使用WHERE ID > lastID的方式代替OFFSET，添加ORDER BY以保证结果一致性
		query := DB.Model(&model.Comment{}).Where("id > ?", lastID).Order("id ASC").Limit(batchSize)

		// 预加载必要的关联数据，但优化查询
		query = query.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "nickname") // 只选择用户的必要字段
		}).Preload("Book", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title") // 只选择书籍的必要字段
		})

		if err := query.Find(&comments).Error; err != nil {
			return err
		}

		sqlDuration := time.Since(sqlStart)
		logger.Infof("评论数据查询耗时: %v, 查询条数: %d", sqlDuration, len(comments))

		if len(comments) == 0 {
			break // 没有更多数据
		}

		// 更新lastID用于下一页查询
		lastID = comments[len(comments)-1].ID

		// 记录索引操作耗时 - 修正计时逻辑
		indexStart := time.Now()

		// 使用bleve的Batch API来批量索引
		batch := s.commentIndex.NewBatch()
		for _, comment := range comments {
			// 转换为评论索引模型
			commentIndex := model.CommentIndex{
				ID:           comment.ID,
				BookTitle:    comment.Book.Title,
				UserNickname: comment.User.Nickname,
				Content:      comment.Content,
				CreateTime:   comment.CreateTime,
			}
			commentID := fmt.Sprintf("%d", comment.ID)
			if err := batch.Index(commentID, commentIndex); err != nil {
				logger.Errorf("批量索引添加评论 %d 失败: %w", comment.ID, err)
			}
		}

		// 一次性提交整个批次
		if err := s.commentIndex.Batch(batch); err != nil {
			logger.Errorf("批量索引提交失败: %w", err)
			return err
		}

		indexDuration := time.Since(indexStart)

		indexedCount += len(comments)

		// 动态调整批次大小
		// 如果处理时间过长，减少批次大小
		if indexDuration > timeoutThreshold {
			newBatchSize := int(float64(batchSize) * 0.7) // 减少30%
			if newBatchSize < 100 {
				newBatchSize = 100 // 设置最小批次大小
			}
			logger.Infof("索引处理超时，调整批次大小从 %d 到 %d", batchSize, newBatchSize)
			batchSize = newBatchSize
		} else if indexDuration < optimalDuration && batchSize < 2000 {
			// 评论数据相对较小，可以设置较小的最大批次上限
			newBatchSize := int(float64(batchSize) * 1.2) // 增加20%
			if newBatchSize > 2000 {
				newBatchSize = 2000 // 设置最大批次大小
			}
			logger.Infof("索引处理速度快，调整批次大小从 %d 到 %d", batchSize, newBatchSize)
			batchSize = newBatchSize
		}

		logger.Infof("已索引评论数据进度: %d/%d, 索引耗时: %v, 平均每秒索引: %.2f 条",
			indexedCount, total, indexDuration, float64(len(comments))/indexDuration.Seconds())
	}
	// 记录总耗时
	totalDuration := time.Since(totalStart)
	logger.Infof("评论批量索引总耗时: %v, 平均每秒索引: %.2f 条",
		totalDuration, float64(total)/totalDuration.Seconds())
	return nil
}

// 检查索引是否为空
func (s *SearchService) isIndexEmpty() bool {
	// 通过搜索一个简单查询来检查索引是否有数据
	bookSearchRequest := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
	bookSearchRequest.Size = 1
	bookResult, err := s.bookIndex.Search(bookSearchRequest)
	if err != nil || bookResult.Total == 0 {
		return true
	}
	return false
}

// 增量更新索引，只索引新增或修改的数据
func (s *SearchService) updateIndexIncrementally() error {
	startTime := time.Now()
	logger.Infof("开始增量更新索引")

	// 获取最后更新的时间戳
	var lastUpdateTime time.Time

	// 查询数据库中最近更新的数据
	var lastBook model.Book
	// 添加索引和限制以加速查询
	if err := DB.Order("update_time desc").Limit(1).First(&lastBook).Error; err == nil {
		lastUpdateTime = lastBook.UpdateTime
	} else {
		// 如果没有数据，使用当前时间的前一天
		lastUpdateTime = time.Now().AddDate(0, 0, -1)
	}

	logger.Infof("增量更新索引，更新时间点: %v 之后的数据", lastUpdateTime)

	// 增量更新书籍数据
	var updatedBooks []model.Book
	if err := DB.Where("update_time > ?", lastUpdateTime).Find(&updatedBooks).Error; err == nil {
		logger.Infof("找到 %d 条需要增量更新的书籍", len(updatedBooks))
		// 使用批量索引处理书籍数据
		if len(updatedBooks) > 0 {
			bookBatch := s.bookIndex.NewBatch()
			for _, book := range updatedBooks {
				bookID := fmt.Sprintf("%d", book.ID)
				if err := bookBatch.Index(bookID, book); err != nil {
					logger.Errorf("批量索引书籍失败: %v, 书籍ID: %d", err, book.ID)
				}
			}
			if err := s.bookIndex.Batch(bookBatch); err != nil {
				logger.Errorf("提交书籍批量索引失败: %v", err)
			} else {
				logger.Infof("成功批量更新 %d 条书籍索引", len(updatedBooks))
			}
		}
	}

	// 增量更新用户数据 - 使用批量索引优化
	var updatedUsers []model.User
	// 使用分页查询，避免一次加载过多数据
	var userTotal int64
	if err := DB.Model(&model.User{}).Where("update_time > ?", lastUpdateTime).Count(&userTotal).Error; err == nil {
		logger.Infof("找到 %d 条需要增量更新的用户", userTotal)

		// 分页处理用户数据，每页500条
		pageSize := 500
		for offset := 0; offset < int(userTotal); offset += pageSize {
			updatedUsers = []model.User{}
			if err := DB.Where("update_time > ?", lastUpdateTime).Offset(offset).Limit(pageSize).Find(&updatedUsers).Error; err == nil {
				logger.Infof("处理用户批次 %d-%d, 共 %d 条记录", offset, offset+len(updatedUsers), len(updatedUsers))

				// 使用批量索引处理当前页用户数据
				userBatch := s.userIndex.NewBatch()
				for _, user := range updatedUsers {
					userID := fmt.Sprintf("%d", user.ID)
					if err := userBatch.Index(userID, user); err != nil {
						logger.Errorf("批量索引用户失败: %v, 用户ID: %d", err, user.ID)
					}
				}
				if err := s.userIndex.Batch(userBatch); err != nil {
					logger.Errorf("提交用户批量索引失败: %v", err)
				} else {
					logger.Infof("成功批量更新 %d 条用户索引, 进度: %.1f%%", len(updatedUsers), float64(offset+len(updatedUsers))/float64(userTotal)*100)
				}
			}
			// 短暂休眠，避免过度占用系统资源
			if offset+pageSize < int(userTotal) {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	// 增量更新评论数据 - 修正字段名一致性问题
	var updatedComments []model.Comment
	// 使用与其他表一致的字段名update_time
	if err := DB.Where("update_time > ?", lastUpdateTime).Preload("User").Preload("Book").Find(&updatedComments).Error; err == nil {
		logger.Infof("找到 %d 条需要增量更新的评论", len(updatedComments))
		// 使用批量索引处理评论数据
		if len(updatedComments) > 0 {
			commentBatch := s.commentIndex.NewBatch()
			for _, comment := range updatedComments {
				commentID := fmt.Sprintf("%d", comment.ID)
				commentIndex := model.CommentIndex{
					ID:           comment.ID,
					BookTitle:    comment.Book.Title,
					UserNickname: comment.User.Nickname,
					Content:      comment.Content,
					CreateTime:   comment.CreateTime,
				}
				if err := commentBatch.Index(commentID, commentIndex); err != nil {
					logger.Errorf("批量索引评论失败: %v, 评论ID: %d", err, comment.ID)
				}
			}
			if err := s.commentIndex.Batch(commentBatch); err != nil {
				logger.Errorf("提交评论批量索引失败: %v", err)
			} else {
				logger.Infof("成功批量更新 %d 条评论索引", len(updatedComments))
			}
		}
	}

	duration := time.Since(startTime)
	logger.Infof("增量更新索引完成，总耗时: %v", duration)
	return nil
}

// 设置索引状态
func (s *SearchService) setIndexReady(ready bool) {
	s.indexMutex.Lock()
	defer s.indexMutex.Unlock()
	s.indexReady = ready
}

// 获取索引状态
func (s *SearchService) IsIndexReady() bool {
	s.indexMutex.RLock()
	defer s.indexMutex.RUnlock()
	return s.indexReady
}

// 优化的IndexBook方法，减少锁竞争
func (s *SearchService) IndexBook(book model.Book) error {
	// 预先生成文档ID，减少锁内操作
	bookID := fmt.Sprintf("%d", book.ID)
	return s.bookIndex.Index(bookID, book)
}

func (s *SearchService) IndexUser(user model.User) error {
	return s.userIndex.Index(fmt.Sprintf("%d", user.ID), user)
}

func (s *SearchService) IndexComment(comment model.Comment) error {
	commentIndex := model.CommentIndex{
		ID:           comment.ID,
		BookTitle:    comment.Book.Title,
		UserNickname: comment.User.Nickname,
		Content:      comment.Content,
		CreateTime:   comment.CreateTime,
	}
	return s.commentIndex.Index(fmt.Sprintf("%d", comment.ID), commentIndex)
}

func (s *SearchService) DeleteUserIndex(userID uint) error {
	return s.userIndex.Delete(fmt.Sprintf("%d", userID))
}

func (s *SearchService) DeleteBookIndex(bookID uint) error {
	return s.bookIndex.Delete(fmt.Sprintf("%d", bookID))
}

func (s *SearchService) DeleteCommentIndex(commentID uint) error {
	return s.commentIndex.Delete(fmt.Sprintf("%d", commentID))
}

// ==================== 搜索功能 ====================

func (s *SearchService) SearchBooks(query string, limit, offset int) ([]model.Book, int64, error) {
	if query == "" {
		return s.getAllBooks(limit, offset)
	}

	// 检查索引是否就绪
	if !s.IsIndexReady() {
		logger.Warn("索引尚未就绪，返回空结果")
		return []model.Book{}, 0, nil
	}

	// 使用读锁进行搜索操作，允许与其他读操作并发
	if !mu.TryRLock() {
		logger.Error("服务太忙")
		return nil, 0, fmt.Errorf("服务太忙")
	}
	defer mu.RUnlock()

	matchPhraseQuery := bleve.NewMatchPhraseQuery(query)

	searchRequest := bleve.NewSearchRequest(matchPhraseQuery)
	searchRequest.From = offset
	searchRequest.Size = limit

	searchResult, err := s.bookIndex.Search(searchRequest)
	if err != nil {
		return nil, 0, err
	}

	var books []model.Book
	for _, hit := range searchResult.Hits {
		var book model.Book
		if err := DB.First(&book, hit.ID).Error; nil == err {
			books = append(books, book)
		}
	}
	return books, int64(searchResult.Total), nil
}

func (s *SearchService) SearchUsers(query string, limit, offset int) ([]model.User, int64, error) {
	if query == "" {
		return s.getAllUsers(limit, offset)
	}

	// 检查索引是否就绪
	if !s.IsIndexReady() {
		logger.Warn("索引尚未就绪，返回空结果")
		return []model.User{}, 0, nil
	}

	// 使用读锁进行搜索操作
	if !mu.TryRLock() {
		logger.Error("服务太忙")
		return nil, 0, fmt.Errorf("服务太忙")
	}
	defer mu.RUnlock()

	// 电话号码支持部分匹配
	if isNumeric(query) {
		query = "*" + query + "*"
	}

	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchRequest.From = offset
	searchRequest.Size = limit

	searchResult, err := s.userIndex.Search(searchRequest)
	if err != nil {
		return nil, 0, err
	}

	var users []model.User
	for _, hit := range searchResult.Hits {
		var user model.User
		if err := DB.First(&user, hit.ID).Error; nil == err {
			users = append(users, user)
		}
	}
	return users, int64(searchResult.Total), nil
}

func (s *SearchService) SearchComments(query string, limit, offset int) ([]model.CommentIndex, int64, error) {
	if query == "" {
		return s.getAllComments(limit, offset)
	}

	// 检查索引是否就绪
	if !s.IsIndexReady() {
		logger.Warn("索引尚未就绪，返回空结果")
		return []model.CommentIndex{}, 0, nil
	}

	// 使用读锁进行搜索操作
	if !mu.TryRLock() {
		logger.Error("服务太忙")
		return nil, 0, fmt.Errorf("服务太忙")
	}
	defer mu.RUnlock()

	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchRequest.From = offset
	searchRequest.Size = limit

	searchResult, err := s.commentIndex.Search(searchRequest)
	if err != nil {
		return nil, 0, err
	}

	var commentIndexes []model.CommentIndex
	for _, hit := range searchResult.Hits {
		commentID, _ := strconv.ParseUint(hit.ID, 10, 32)
		var comment model.Comment
		if err := DB.Preload("User").Preload("Book").First(&comment, commentID).Error; nil == err {
			commentIndex := model.CommentIndex{
				ID:           comment.ID,
				BookTitle:    comment.Book.Title,
				UserNickname: comment.User.Nickname,
				Content:      comment.Content,
				CreateTime:   comment.CreateTime,
			}
			commentIndexes = append(commentIndexes, commentIndex)
		}
	}
	return commentIndexes, int64(searchResult.Total), nil
}

// ==================== 辅助函数 ====================

func (s *SearchService) getAllBooks(limit, offset int) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64
	DB.Model(&model.Book{}).Count(&total)
	DB.Offset(offset).Limit(limit).Find(&books)
	return books, total, nil
}

func (s *SearchService) getAllUsers(limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	DB.Model(&model.User{}).Count(&total)
	DB.Offset(offset).Limit(limit).Find(&users)
	return users, total, nil
}

func (s *SearchService) getAllComments(limit, offset int) ([]model.CommentIndex, int64, error) {
	var comments []model.Comment
	var total int64
	DB.Model(&model.Comment{}).Count(&total)
	DB.Preload("User").Preload("Book").Offset(offset).Limit(limit).Order("create_time desc").Find(&comments)

	var commentIndexes []model.CommentIndex
	for _, comment := range comments {
		commentIndexes = append(commentIndexes, model.CommentIndex{
			ID:           comment.ID,
			BookTitle:    comment.Book.Title,
			UserNickname: comment.User.Nickname,
			Content:      comment.Content,
			CreateTime:   comment.CreateTime,
		})
	}
	return commentIndexes, total, nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
