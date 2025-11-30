package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mageg-x/novel/src/log"
	"github.com/mageg-x/novel/src/model"
	"github.com/mageg-x/novel/src/util"
)

var (
	DB       *gorm.DB
	DataDir  = "data" // 默认数据目录，可通过InitDB函数参数覆盖
	BooksDir = filepath.Join(DataDir, "books")
	DBPath   = filepath.Join(DataDir, "books.db")
	logger   = log.GetLogger("novel")
)

// InitDB 初始化数据库连接，接受数据目录参数
func InitDB(dataDir string, bookDir string) error {
	// 如果提供了自定义数据目录，则更新全局变量
	if dataDir != "" {
		DataDir, _ = filepath.Abs(dataDir)
		DBPath = filepath.Join(DataDir, "books.db")
		BooksDir = filepath.Join(DataDir, "books")
	}
	if bookDir != "" {
		BooksDir, _ = filepath.Abs(bookDir)
	}

	// 创建数据目录
	os.MkdirAll(BooksDir, 0755)

	// 配置数据库连接
	config := &gorm.Config{}
	db, err := gorm.Open(sqlite.Open(DBPath), config)
	if err != nil {
		logger.Errorf("数据库连接失败: %v", err)
		return err
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorf("获取数据库连接池失败: %v", err)
		return err
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接的最大生命周期
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	logger.Infof("数据库连接成功")

	// SQLite 时间格式支持
	DB.Exec("PRAGMA foreign_keys = ON;")
	DB.Exec("PRAGMA journal_mode = WAL;")
	// 启用查询缓存
	DB.Exec("PRAGMA cache_size = -8000;") // 8MB缓存
	DB.Exec("PRAGMA synchronous = NORMAL;")

	// 自动迁移数据库表结构
	// if err := DB.AutoMigrate(&model.Book{}, &model.Chapter{}, &model.Rcmd{}, &model.Rank{}, &model.User{}, &model.Shelf{}, &model.History{}, &model.Comment{}); err != nil {
	// 	logger.Errorf("数据库迁移失败: %v", err)
	// 	return err
	// }

	return nil
}

// 书籍服务
type BookService struct{}

// 获取所有书籍
func (s *BookService) GetAllBooks(offset, limit int) ([]model.Book, int64, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取书籍总数失败: 数据库连接未初始化")
		return nil, 0, fmt.Errorf("数据库连接未初始化")
	}
	var books []model.Book
	var total int64

	// 优化Count查询，使用更高效的方式获取总数
	// 对于SQLite，直接使用COUNT(*) 是最高效的，但我们可以添加缓存来优化重复查询
	if err := DB.Model(&model.Book{}).Count(&total).Error; err != nil {
		logger.Errorf("获取书籍总数失败: %v", err)
		return nil, 0, err
	}

	// 优化分页查询，利用update_time索引提升排序性能
	// 移除或减少预加载以提升性能，Chapters可以在需要时单独查询
	err := DB.Offset(offset).Limit(limit).Order("update_time desc").Find(&books).Error
	if err != nil {
		logger.Errorf("获取书籍分页数据失败[偏移: %d, 限制: %d]: %v", offset, limit, err)
	}

	// 替换books 的cover url
	// https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/001/b_13806001.webp
	for i := range books {
		books[i].Cover = fmt.Sprintf("https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/%03d/b_%d.webp", books[i].ID%1000, books[i].ID)
	}
	return books, total, err
}

// 根据ID获取书籍
func (s *BookService) GetBookByID(id uint) (*model.Book, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取书籍失败[ID: %d]: 数据库连接未初始化", id)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var book model.Book
	err := DB.First(&book, id).Error
	if err != nil {
		logger.Errorf("获取书籍失败[ID: %d]: %v", id, err)
		return nil, err
	}
	book.Cover = fmt.Sprintf("https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/%03d/b_%d.webp", book.ID%1000, book.ID)

	return &book, nil
}

// 根据分类获取书籍
func (s *BookService) GetBooksByCategory(category string, offset, limit int) ([]model.Book, int64, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取分类[%s]书籍总数失败: 数据库连接未初始化", category)
		return nil, 0, fmt.Errorf("数据库连接未初始化")
	}
	var books []model.Book
	var total int64

	// 优化Count查询，利用category索引
	if err := DB.Model(&model.Book{}).Where("category = ?", category).Count(&total).Error; err != nil {
		logger.Errorf("获取分类[%s]书籍总数失败: %v", category, err)
		return nil, 0, err
	}

	// 优化分页查询，利用category和update_time索引提升性能
	// 移除预加载以提升性能
	err := DB.Where("category = ?", category).Offset(offset).Limit(limit).Order("update_time desc").Find(&books).Error
	if err != nil {
		logger.Errorf("获取分类[%s]书籍数据失败: %v", category, err)
	}

	// 替换books 的cover url
	// https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/001/b_13806001.webp
	for i := range books {
		books[i].Cover = fmt.Sprintf("https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/%03d/b_%d.webp", books[i].ID%1000, books[i].ID)
	}
	return books, total, err
}

// 添加书籍
func (s *BookService) AddBook(book *model.Book) error {
	// 检查指针是否为nil
	if book == nil {
		logger.Errorf("添加书籍失败: 书籍指针为nil")
		return fmt.Errorf("书籍指针为nil")
	}
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("添加书籍失败[标题: %s]: 数据库连接未初始化", book.Title)
		return fmt.Errorf("数据库连接未初始化")
	}
	book.CreateTime = time.Now()
	book.UpdateTime = time.Now()
	if err := DB.Create(book).Error; err != nil {
		logger.Errorf("添加书籍失败[标题: %s]: %v", book.Title, err)
		return err
	}
	logger.Infof("添加书籍成功[ID: %d, 标题: %s]", book.ID, book.Title)
	return nil
}

// 更新书籍
func (s *BookService) UpdateBook(book *model.Book) error {
	// 检查指针是否为nil
	if book == nil {
		logger.Errorf("更新书籍失败: 书籍指针为nil")
		return fmt.Errorf("书籍指针为nil")
	}
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("更新书籍失败[ID: %d, 标题: %s]: 数据库连接未初始化", book.ID, book.Title)
		return fmt.Errorf("数据库连接未初始化")
	}
	book.UpdateTime = time.Now()
	if err := DB.Save(book).Error; err != nil {
		logger.Errorf("更新书籍失败[ID: %d, 标题: %s]: %v", book.ID, book.Title, err)
		return err
	}
	logger.Infof("更新书籍成功[ID: %d, 标题: %s]", book.ID, book.Title)
	return nil
}

// 删除书籍
func (s *BookService) DeleteBook(id uint) error {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("删除书籍失败[ID: %d]: 数据库连接未初始化", id)
		return fmt.Errorf("数据库连接未初始化")
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		// 获取书籍信息用于日志
		var book model.Book
		if err := tx.First(&book, id).Error; err != nil {
			logger.Errorf("获取书籍信息失败[ID: %d]: %v", id, err)
			return err
		}

		// 删除书籍的章节
		if err := tx.Where("book_id = ?", id).Delete(&model.Chapter{}).Error; err != nil {
			logger.Errorf("删除书籍章节失败[书籍ID: %d]: %v", id, err)
			return err
		}

		// 删除书籍
		if err := tx.Delete(&model.Book{}, id).Error; err != nil {
			logger.Errorf("删除书籍失败[ID: %d, 标题: %s]: %v", id, book.Title, err)
			return err
		}

		logger.Infof("删除书籍成功[ID: %d, 标题: %s]", id, book.Title)
		return nil
	})
}

// 章节服务
// 获取书籍的所有章节
// GetChaptersByBookID 获取某本书的所有章节（按章节序号排序）
func (s *BookService) GetChaptersByBookID(bookID uint) (*model.Chapters, error) {
	if DB == nil {
		logger.Errorf("获取书籍章节失败[书籍ID: %d]: 数据库连接未初始化", bookID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	// 从books 表获取书籍信息
	book, err := s.GetBookByID(bookID)
	if err != nil || book == nil {
		logger.Errorf("获取书籍信息失败[ID: %d]: %v", bookID, err)
		return nil, err
	}
	// 读取.chapters.json 文件
	chapPath := filepath.Join(BooksDir, util.GetFirstChar(book.Title), book.Title, ".chapters.json")
	chapCont, err := os.ReadFile(chapPath)
	if err != nil || chapCont == nil {
		logger.Errorf("读取.chapters.json 文件失败[书籍ID: %d]: %v", bookID, err)
		return nil, err
	}
	var chapters model.Chapters
	err = json.Unmarshal(chapCont, &chapters)
	if err != nil {
		logger.Errorf("解析.chapters.json 文件失败[书籍ID: %d]: %v", bookID, err)
		return nil, err
	}
	return &chapters, nil
}

// GetChapterByNo 根据书籍ID和章节序号获取章节
func (s *BookService) GetChapterByNo(bookID, chapterNo uint) (*model.Chapter, error) {
	chapters, err := s.GetChaptersByBookID(bookID)
	if err != nil || chapters == nil {
		logger.Errorf("获取章节失败[书籍ID: %d,  章节No: %d]: %v", bookID, chapterNo, err)
		return nil, err
	}
	var chapter *model.Chapter
	for i := range chapters.Chapters {
		if chapters.Chapters[i].ChapterNo == chapterNo {
			chapter = &chapters.Chapters[i]
			break
		}
	}
	if chapter == nil {
		logger.Errorf("获取章节失败[书籍ID: %d,  章节No: %d]: 未找到章节", bookID, chapterNo)
		return nil, fmt.Errorf("未找到章节")
	}
	// 读取章节内容
	content, err := s.GetChapterContent(bookID, chapterNo)
	if err != nil {
		logger.Errorf("获取章节失败[书籍ID: %d,  章节No: %d]: %v", bookID, chapterNo, err)
		return nil, err
	}
	chapter.Content = content
	return chapter, nil
}

// AddChapter 添加章节
func (s *BookService) AddChapter(bookID uint, title, content string, isVip bool) error {
	book, err := s.GetBookByID(bookID)
	if err != nil || book == nil {
		logger.Errorf("获取书籍信息失败[ID: %d]: %v", bookID, err)
		return err
	}
	chapters, err := s.GetChaptersByBookID(bookID)
	if err != nil || chapters == nil {
		logger.Errorf("获取章节列表失败[书籍ID: %d]: %v", bookID, err)
		return err
	}
	chapNo := uint(chapters.MaxChapterNo) + 1
	filename := fmt.Sprintf("第%s章 %s.txt", util.NumToChinese(chapNo), title)
	// 写入chapter 文件内容
	contPath := filepath.Join(BooksDir, util.GetFirstChar(book.Title), book.Title, filename)
	if err := os.WriteFile(contPath, []byte(content), 0644); err != nil {
		logger.Errorf("保存章节内容失败[路径: %s]: %v", contPath, err)
		return err
	}

	// 更新 .chapters.json 文件
	ts := time.Now().Unix()
	chapter := model.Chapter{
		Title:      title,
		ChapterNo:  uint(chapters.MaxChapterNo) + 1,
		CreateTime: ts,
		UpdateTime: ts,
		IsVip:      isVip,
	}

	chapters.Chapters = append(chapters.Chapters, chapter)
	chapters.TotalChapters = uint(len(chapters.Chapters))
	chapters.MaxChapterNo = chapter.ChapterNo
	chapPath := filepath.Join(BooksDir, util.GetFirstChar(book.Title), book.Title, ".chapters.json")
	chapCont, err := json.Marshal(chapters)
	if err := os.WriteFile(chapPath, chapCont, 0644); err != nil {
		logger.Errorf("保存章节列表失败[路径: %s]: %v", chapPath, err)
		return err
	}

	// 更新 books 表
	book.UpdateTime = time.Now()
	book.WordCount += len(chapter.Content)
	if err := DB.Save(book).Error; err != nil {
		logger.Errorf("更新书籍信息失败[ID: %d]: %v", book.ID, err)
	}

	logger.Infof("添加章节成功[book: %s,  标题: %s]", chapters.BookName, chapter.Title)
	return nil
}

// GetChapterContent 获取章节内容
func (s *BookService) GetChapterContent(bookID, chapterNo uint) (string, error) {
	chapters, err := s.GetChaptersByBookID(bookID)
	if err != nil || chapters == nil {
		logger.Errorf("获取章节失败[书籍ID: %d,  章节No: %d]: %v", bookID, chapterNo, err)
		return "", err
	}
	var chapter *model.Chapter
	for i := range chapters.Chapters {
		if chapters.Chapters[i].ChapterNo == chapterNo {
			chapter = &chapters.Chapters[i]
			break
		}
	}
	if chapter == nil {
		logger.Errorf("获取章节失败[书籍ID: %d,  章节No: %d]: 未找到章节", bookID, chapterNo)
		return "", fmt.Errorf("未找到章节")
	}

	// 读取章节内容
	filename := fmt.Sprintf("第%s章 %s.txt", util.NumToChinese(chapter.ChapterNo), chapter.Title)
	contPath := filepath.Join(BooksDir, util.GetFirstChar(chapters.BookName), chapters.BookName, filename)
	cont, err := os.ReadFile(contPath)
	if err != nil {
		logger.Errorf("读取章节内容失败[路径: %s]: %v", contPath, err)
		return "", err
	}
	return string(cont), nil
}

// UpdateChapter 更新章节
func (s *BookService) UpdateChapter(bookID, chapterNo uint, title, content string, isVip bool) error {
	chapters, err := s.GetChaptersByBookID(bookID)
	if err != nil || chapters == nil {
		logger.Errorf("更新章节失败[书籍ID: %d,  章节No: %d]: %v", bookID, chapterNo, err)
		return err
	}
	var chapter *model.Chapter
	for i := range chapters.Chapters {
		if chapters.Chapters[i].ChapterNo == chapterNo {
			chapter = &chapters.Chapters[i]
			break
		}
	}
	if chapter == nil {
		logger.Errorf("更新章节失败[书籍ID: %d,  章节No: %d]: 未找到章节", bookID, chapterNo)
		return fmt.Errorf("未找到章节")
	}
	// 更新章节标题
	chapter.Title = title
	chapter.UpdateTime = time.Now().Unix()
	chapter.IsVip = isVip

	// 更新 .chapters.json 文件
	chapPath := filepath.Join(BooksDir, util.GetFirstChar(chapters.BookName), chapters.BookName, ".chapters.json")
	chapCont, err := json.Marshal(chapters)
	if err := os.WriteFile(chapPath, chapCont, 0644); err != nil {
		logger.Errorf("保存章节列表失败[路径: %s]: %v", chapPath, err)
		return err
	}

	// 更新 章节内容
	filename := fmt.Sprintf("第%s章 %s.txt", util.NumToChinese(chapter.ChapterNo), chapter.Title)
	contPath := filepath.Join(BooksDir, util.GetFirstChar(chapters.BookName), chapters.BookName, filename)
	if err := os.WriteFile(contPath, []byte(content), 0644); err != nil {
		logger.Errorf("保存章节内容失败[路径: %s]: %v", contPath, err)
		return err
	}

	// 更新 books 表
	book, err := s.GetBookByID(bookID)
	if err != nil || book == nil {
		logger.Errorf("更新书籍信息失败[ID: %d]: %v", bookID, err)
	}
	book.UpdateTime = time.Now()
	if err := DB.Save(book).Error; err != nil {
		logger.Errorf("更新书籍信息失败[ID: %d]: %v", book.ID, err)
	}
	logger.Infof("更新章节成功[book: %s,  标题: %s]", chapters.BookName, chapter.Title)
	return nil
}

// DeleteChapter 删除章节
func (s *BookService) DeleteChapter(bookID, chapterNo uint) error {
	chapters, err := s.GetChaptersByBookID(bookID)
	if err != nil || chapters == nil {
		logger.Errorf("删除章节失败[书籍ID: %d,  章节No: %d]: %v", bookID, chapterNo, err)
		return err
	}

	var chapter *model.Chapter
	for i, c := range chapters.Chapters {
		if c.ChapterNo == chapterNo {
			// 执行删除
			chapters.Chapters = append(chapters.Chapters[:i], chapters.Chapters[i+1:]...)
			chapter = &chapters.Chapters[i]
			break
		}
	}

	if chapter == nil {
		logger.Errorf("删除章节失败[书籍ID: %d, 章节No: %d]: 未找到章节", bookID, chapterNo)
		return fmt.Errorf("未找到章节")
	}

	// 删除章节内容
	filename := fmt.Sprintf("第%s章 %s.txt", util.NumToChinese(chapter.ChapterNo), chapter.Title)
	contPath := filepath.Join(BooksDir, util.GetFirstChar(chapters.BookName), chapters.BookName, filename)
	if err := os.Remove(contPath); err != nil {
		logger.Errorf("删除章节内容失败[路径: %s]: %v", contPath, err)
		return err
	}

	// 更新 .chapters.json 文件
	chapPath := filepath.Join(BooksDir, util.GetFirstChar(chapters.BookName), chapters.BookName, ".chapters.json")
	chapCont, err := json.Marshal(chapters)
	if err := os.WriteFile(chapPath, chapCont, 0644); err != nil {
		logger.Errorf("保存章节列表失败[路径: %s]: %v", chapPath, err)
		return err
	}

	// 更新 books 表
	book, err := s.GetBookByID(bookID)
	if err != nil || book == nil {
		logger.Errorf("更新书籍信息失败[ID: %d]: %v", bookID, err)
	}
	book.UpdateTime = time.Now()
	if err := DB.Save(book).Error; err != nil {
		logger.Errorf("更新书籍信息失败[ID: %d]: %v", book.ID, err)
	}

	logger.Infof("删除章节成功[book: %s,  标题: %s]", chapters.BookName, chapter.Title)
	return nil
}

// 排行榜服务
// 获取指定类型的排行榜
func (s *BookService) GetRankByType(rankType string, rankTypeName string) ([]model.Rank, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取%v失败: 数据库连接未初始化", rankTypeName)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var ranks []model.Rank
	// 查询时排除时间字段，避免类型转换错误
	err := DB.Preload("Book").Select("id, rank_type, book_id, \"order\"").Where("rank_type = ?", rankType).Order("\"order\" asc").Find(&ranks).Error
	if err != nil {
		logger.Errorf("获取%v失败: %v", rankTypeName, err)
		return nil, err
	}

	// 替换 ranks 中book的 cover_url
	for i := range ranks {
		book := &ranks[i].Book
		ranks[i].Book.Cover = fmt.Sprintf("https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/%03d/b_%d.webp", book.ID%1000, book.ID)
	}
	return ranks, nil
}

// 推荐服务
// 获取指定类型的推荐
func (s *BookService) GetRcmdByType(rcmdType string, rcmdTypeName string) ([]model.Rcmd, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取%v失败: 数据库连接未初始化", rcmdTypeName)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var rcmds []model.Rcmd
	err := DB.Preload("Book").Where("rcmd_type = ?", rcmdType).Order("\"order\" asc").Find(&rcmds).Error
	if err != nil {
		logger.Errorf("获取%v失败: %v", rcmdTypeName, err)
		return nil, err
	}

	// 替换 rcmds 中book的 cover_url
	for i := range rcmds {
		book := &rcmds[i].Book
		rcmds[i].Book.Cover = fmt.Sprintf("https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/%03d/b_%d.webp", book.ID%1000, book.ID)
	}
	return rcmds, nil
}

// 获取相关书籍（同类别，点击率高的书籍）
func (s *BookService) GetRelatedBooks(bookID uint, limit int) ([]model.Book, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取相关书籍失败: 数据库连接未初始化")
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	// 先获取当前书籍的类别
	var currentBook model.Book
	err := DB.Select("category").First(&currentBook, bookID).Error
	if err != nil {
		logger.Errorf("获取当前书籍信息失败[ID: %d]: %v", bookID, err)
		return nil, err
	}

	// 查询同类别中点击率高的书籍，排除当前书籍本身
	var relatedBooks []model.Book
	err = DB.Where("category = ? AND id != ?", currentBook.Category, bookID).
		Order("click_count desc").
		Limit(limit).
		Find(&relatedBooks).Error
	if err != nil {
		logger.Errorf("获取相关书籍失败[类别: %s, 排除ID: %d]: %v", currentBook.Category, bookID, err)
		return nil, err
	}

	// 替换books 的cover url
	// https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/001/b_13806001.webp
	for i := range relatedBooks {
		relatedBooks[i].Cover = fmt.Sprintf("https://cdn.jsdelivr.net/gh/mageg-x/novel-library/covers/%03d/b_%d.webp", relatedBooks[i].ID%1000, relatedBooks[i].ID)
	}
	return relatedBooks, nil
}

// 评论服务
// 获取书籍的评论列表
func (s *BookService) GetCommentsByBookID(bookID uint) ([]model.Comment, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取书籍评论失败[书籍ID: %d]: 数据库连接未初始化", bookID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var comments []model.Comment
	err := DB.Preload("User").Where("book_id = ?", bookID).Order("create_time desc").Find(&comments).Error
	if err != nil {
		logger.Errorf("获取书籍评论失败[书籍ID: %d]: %v", bookID, err)
		return nil, err
	}
	return comments, nil
}
