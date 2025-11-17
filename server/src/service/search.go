package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/yanyiwu/gojieba"

	"github.com/mageg-x/novel/src/model"
)

// SearchService 搜索服务结构体
type SearchService struct {
	bookIndex    bleve.Index
	userIndex    bleve.Index
	commentIndex bleve.Index
	indexDir     string
}

// GoJiebaTokenizer 自定义分词器
type GoJiebaTokenizer struct {
	jieba *gojieba.Jieba
}

// NewGoJiebaTokenizer 创建新的分词器
func NewGoJiebaTokenizer() *GoJiebaTokenizer {
	return &GoJiebaTokenizer{
		jieba: gojieba.NewJieba(),
	}
}

// Tokenize 实现分词接口
func (t *GoJiebaTokenizer) Tokenize(input []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)

	// 使用 gojieba 进行分词
	words := t.jieba.Cut(string(input), true)

	position := 0
	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		token := analysis.Token{
			Term:     []byte(word),
			Start:    position,
			End:      position + len(word),
			Position: position + 1,
			Type:     analysis.Ideographic,
		}
		result = append(result, &token)
		position += len(word)
	}

	return result
}

// Free 释放资源
func (t *GoJiebaTokenizer) Free() {
	if t.jieba != nil {
		t.jieba.Free()
	}
}

// GoJiebaTokenizerName 分词器名称
const GoJiebaTokenizerName = "gojieba"

// GoJiebaTokenizerConstructor 分词器构造函数
func GoJiebaTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	return NewGoJiebaTokenizer(), nil
}

// 简单的停用词过滤器
type ChineseStopFilter struct {
	stopWords map[string]struct{}
}

// NewChineseStopFilter 创建中文停用词过滤器
func NewChineseStopFilter(stopWords []string) *ChineseStopFilter {
	filter := &ChineseStopFilter{
		stopWords: make(map[string]struct{}),
	}
	for _, word := range stopWords {
		filter.stopWords[word] = struct{}{}
	}
	return filter
}

// Filter 实现停用词过滤
func (f *ChineseStopFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	output := make(analysis.TokenStream, 0, len(input))

	for _, token := range input {
		if _, isStop := f.stopWords[string(token.Term)]; !isStop {
			output = append(output, token)
		}
	}

	return output
}

const ChineseStopFilterName = "chinese_stop"

func ChineseStopFilterConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.TokenFilter, error) {
	// 中文停用词列表
	chineseStopWords := []string{"的", "了", "在", "是", "我", "有", "和", "就", "不", "人", "都", "一", "一个", "上", "也", "很", "到", "说", "要", "去", "你", "会", "着", "没有", "看", "好", "自己", "这", "那", "他", "她", "它"}
	return NewChineseStopFilter(chineseStopWords), nil
}

// 全局搜索服务实例
var (
	SService = newSearchService()
	mu       = sync.Mutex{}
)

func init() {
	// 注册自定义分词器
	registry.RegisterTokenizer(GoJiebaTokenizerName, GoJiebaTokenizerConstructor)
	// 注册中文停用词过滤器
	registry.RegisterTokenFilter(ChineseStopFilterName, ChineseStopFilterConstructor)
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
		if err := s.initBookIndex(); err != nil {
			logger.Errorf("初始化书籍索引失败: %w", err)
			return
		}
		if err := s.initUserIndex(); err != nil {
			logger.Errorf("初始化用户索引失败: %w", err)
			return
		}
		if err := s.initCommentIndex(); err != nil {
			logger.Errorf("初始化评论索引失败: %w", err)
			return
		}

		if err := s.indexExistingData(); err != nil {
			logger.Errorf("索引现有数据失败: %w", err)
			return
		}

		logger.Infof("搜索索引初始化成功")

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

	// 添加自定义中文分析器
	err := indexMapping.AddCustomAnalyzer("chinese_analyzer", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": GoJiebaTokenizerName,
		"token_filters": []string{
			"lowercase",
			ChineseStopFilterName,
		},
	})
	if err != nil {
		logger.Errorf("创建中文分析器失败: %v", err)
		return s.createBookMappingFallback()
	}

	bookMapping := bleve.NewDocumentMapping()

	// ID 字段
	idField := bleve.NewNumericFieldMapping()
	idField.Store = true
	idField.Index = true
	bookMapping.AddFieldMappingsAt("ID", idField)

	// 标题字段 - 使用中文分词
	titleField := bleve.NewTextFieldMapping()
	titleField.Store = true
	titleField.Index = true
	titleField.Analyzer = "chinese_analyzer"
	bookMapping.AddFieldMappingsAt("Title", titleField)

	// 作者字段 - 使用中文分词
	authorField := bleve.NewTextFieldMapping()
	authorField.Store = true
	authorField.Index = true
	authorField.Analyzer = "chinese_analyzer"
	bookMapping.AddFieldMappingsAt("Author", authorField)

	// 描述字段 - 使用中文分词
	descField := bleve.NewTextFieldMapping()
	descField.Store = true
	descField.Index = true
	descField.Analyzer = "chinese_analyzer"
	bookMapping.AddFieldMappingsAt("Description", descField)

	indexMapping.DefaultMapping = bookMapping
	indexMapping.DefaultAnalyzer = "chinese_analyzer"
	return indexMapping
}

func (s *SearchService) createUserMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()

	// 添加自定义中文分析器
	err := indexMapping.AddCustomAnalyzer("chinese_analyzer", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": GoJiebaTokenizerName,
		"token_filters": []string{
			"lowercase",
		},
	})
	if err != nil {
		logger.Errorf("创建中文分析器失败: %v", err)
		return s.createUserMappingFallback()
	}

	userMapping := bleve.NewDocumentMapping()

	// ID 字段
	idField := bleve.NewNumericFieldMapping()
	idField.Store = true
	idField.Index = true
	userMapping.AddFieldMappingsAt("ID", idField)

	// 用户名字段 - 使用中文分词
	usernameField := bleve.NewTextFieldMapping()
	usernameField.Store = true
	usernameField.Index = true
	usernameField.Analyzer = "chinese_analyzer"
	userMapping.AddFieldMappingsAt("Username", usernameField)

	// 昵称字段 - 使用中文分词
	nicknameField := bleve.NewTextFieldMapping()
	nicknameField.Store = true
	nicknameField.Index = true
	nicknameField.Analyzer = "chinese_analyzer"
	userMapping.AddFieldMappingsAt("Nickname", nicknameField)

	// 电话字段 - 不使用分词，使用关键字分析器
	phoneField := bleve.NewTextFieldMapping()
	phoneField.Store = true
	phoneField.Index = true
	phoneField.Analyzer = standard.Name
	userMapping.AddFieldMappingsAt("Phone", phoneField)

	indexMapping.DefaultMapping = userMapping
	indexMapping.DefaultAnalyzer = "chinese_analyzer"
	return indexMapping
}

func (s *SearchService) createCommentMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()

	// 添加自定义中文分析器
	err := indexMapping.AddCustomAnalyzer("chinese_analyzer", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": GoJiebaTokenizerName,
		"token_filters": []string{
			"lowercase",
			ChineseStopFilterName,
		},
	})
	if err != nil {
		logger.Errorf("创建中文分析器失败: %v", err)
		return s.createCommentMappingFallback()
	}

	commentMapping := bleve.NewDocumentMapping()

	// ID 字段
	idField := bleve.NewNumericFieldMapping()
	idField.Store = true
	idField.Index = true
	commentMapping.AddFieldMappingsAt("ID", idField)

	// 书籍标题字段 - 使用中文分词
	bookTitleField := bleve.NewTextFieldMapping()
	bookTitleField.Store = true
	bookTitleField.Index = true
	bookTitleField.Analyzer = "chinese_analyzer"
	commentMapping.AddFieldMappingsAt("BookTitle", bookTitleField)

	// 内容字段 - 使用中文分词
	contentField := bleve.NewTextFieldMapping()
	contentField.Store = true
	contentField.Index = true
	contentField.Analyzer = "chinese_analyzer"
	commentMapping.AddFieldMappingsAt("Content", contentField)

	// 用户昵称字段 - 使用中文分词
	userNicknameField := bleve.NewTextFieldMapping()
	userNicknameField.Store = true
	userNicknameField.Index = true
	userNicknameField.Analyzer = "chinese_analyzer"
	commentMapping.AddFieldMappingsAt("UserNickname", userNicknameField)

	indexMapping.DefaultMapping = commentMapping
	indexMapping.DefaultAnalyzer = "chinese_analyzer"
	return indexMapping
}

// 回退映射（如果中文分析器创建失败）
func (s *SearchService) createBookMappingFallback() *mapping.IndexMappingImpl {
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
	titleField.Analyzer = standard.Name
	bookMapping.AddFieldMappingsAt("Title", titleField)

	// 作者字段
	authorField := bleve.NewTextFieldMapping()
	authorField.Store = true
	authorField.Index = true
	authorField.Analyzer = standard.Name
	bookMapping.AddFieldMappingsAt("Author", authorField)

	// 描述字段
	descField := bleve.NewTextFieldMapping()
	descField.Store = true
	descField.Index = true
	descField.Analyzer = standard.Name
	bookMapping.AddFieldMappingsAt("Description", descField)

	indexMapping.DefaultMapping = bookMapping
	return indexMapping
}

func (s *SearchService) createUserMappingFallback() *mapping.IndexMappingImpl {
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
	usernameField.Analyzer = standard.Name
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

func (s *SearchService) createCommentMappingFallback() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()
	commentMapping := bleve.NewDocumentMapping()

	// 书籍标题字段
	bookTitleField := bleve.NewTextFieldMapping()
	bookTitleField.Store = true
	bookTitleField.Index = true
	bookTitleField.Analyzer = standard.Name
	commentMapping.AddFieldMappingsAt("BookTitle", bookTitleField)

	// 内容字段
	contentField := bleve.NewTextFieldMapping()
	contentField.Store = true
	contentField.Index = true
	contentField.Analyzer = standard.Name
	commentMapping.AddFieldMappingsAt("Content", contentField)

	indexMapping.DefaultMapping = commentMapping
	return indexMapping
}

// ==================== 索引操作 ====================

func (s *SearchService) indexExistingData() error {
	var books []model.Book
	if err := DB.Find(&books).Error; nil == err {
		for _, book := range books {
			s.IndexBook(book)
		}
	}

	var users []model.User
	if err := DB.Find(&users).Error; nil == err {
		for _, user := range users {
			s.IndexUser(user)
		}
	}

	var comments []model.Comment
	if err := DB.Preload("User").Preload("Book").Find(&comments).Error; nil == err {
		for _, comment := range comments {
			s.IndexComment(comment)
		}
	}
	return nil
}

func (s *SearchService) IndexBook(book model.Book) error {
	return s.bookIndex.Index(fmt.Sprintf("%d", book.ID), book)
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

	if !mu.TryLock() {
		logger.Error("服务太忙")
		return nil, 0, fmt.Errorf("服务太忙")
	}
	defer mu.Unlock()

	// 使用 MatchQuery 而不是 QueryStringQuery，更好地支持中文分词
	searchRequest := bleve.NewSearchRequest(bleve.NewMatchQuery(query))
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

	if !mu.TryLock() {
		logger.Error("服务太忙")
		return nil, 0, fmt.Errorf("服务太忙")
	}
	defer mu.Unlock()

	// 如果是纯数字，使用通配符查询电话号码
	if isNumeric(query) {
		wildcardQuery := bleve.NewWildcardQuery("*" + query + "*")
		searchRequest := bleve.NewSearchRequest(wildcardQuery)
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

	// 对于中文文本，使用 MatchQuery
	searchRequest := bleve.NewSearchRequest(bleve.NewMatchQuery(query))
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
	if !mu.TryLock() {
		logger.Error("服务太忙")
		return nil, 0, fmt.Errorf("服务太忙")
	}
	defer mu.Unlock()

	// 使用 MatchPhraseQuery 更好地处理中文短语
	var searchRequest *bleve.SearchRequest
	if isChinesePhrase(query) {
		searchRequest = bleve.NewSearchRequest(bleve.NewMatchPhraseQuery(query))
	} else {
		searchRequest = bleve.NewSearchRequest(bleve.NewMatchQuery(query))
	}
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

// isChinesePhrase 判断是否为中文短语
func isChinesePhrase(s string) bool {
	// 简单的判断：包含中文字符且长度适中
	hasChinese := false
	for _, r := range s {
		if r >= '\u4e00' && r <= '\u9fff' {
			hasChinese = true
			break
		}
	}
	return hasChinese && len([]rune(s)) <= 10
}
