package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mageg-x/novel/src/log"
	"github.com/mageg-x/novel/src/model"
)

var (
	DB       *gorm.DB
	DataDir  = "data"
	BooksDir = filepath.Join(DataDir, "books")
	DBPath   = filepath.Join(DataDir, "books.db")
	logger   = log.GetLogger("novel")
)

// InitDB 初始化数据库连接
func InitDB() error {
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

	// 自动迁移数据库表结构
	if err := DB.AutoMigrate(&model.Book{}, &model.Chapter{}, &model.Rcmd{}, &model.Rank{}, &model.User{}, &model.Shelf{}, &model.History{}, &model.Comment{}); err != nil {
		logger.Errorf("数据库迁移失败: %v", err)
		return err
	}

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

	// 计算总数
	if err := DB.Model(&model.Book{}).Count(&total).Error; err != nil {
		logger.Errorf("获取书籍总数失败: %v", err)
		return nil, 0, err
	}

	// 获取分页数据
	err := DB.Offset(offset).Limit(limit).Order("update_time desc").Find(&books).Error
	if err != nil {
		logger.Errorf("获取书籍分页数据失败[偏移: %d, 限制: %d]: %v", offset, limit, err)
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

	// 计算总数
	if err := DB.Model(&model.Book{}).Where("category = ?", category).Count(&total).Error; err != nil {
		logger.Errorf("获取分类[%s]书籍总数失败: %v", category, err)
		return nil, 0, err
	}

	// 获取分页数据
	err := DB.Where("category = ?", category).Offset(offset).Limit(limit).Order("update_time desc").Find(&books).Error
	if err != nil {
		logger.Errorf("获取分类[%s]书籍数据失败: %v", category, err)
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
func (s *BookService) GetChaptersByBookID(bookID uint) ([]model.Chapter, error) {
	if DB == nil {
		logger.Errorf("获取书籍章节失败[书籍ID: %d]: 数据库连接未初始化", bookID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var chapters []model.Chapter
	err := DB.Where("book_id = ?", bookID).Order("chapter_id ASC").Find(&chapters).Error
	if err != nil {
		logger.Errorf("获取书籍章节失败[书籍ID: %d]: %v", bookID, err)
	}
	return chapters, err
}

// GetChapterByID 根据书籍ID和章节序号获取章节
func (s *BookService) GetChapterByID(book_id, chapter_id uint) (*model.Chapter, error) {
	if DB == nil {
		logger.Errorf("获取章节失败[书籍ID: %d, 章节ID: %d]: 数据库连接未初始化", book_id, chapter_id)
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	var chapter model.Chapter
	err := DB.Where("book_id = ? AND chapter_id = ?", book_id, chapter_id).First(&chapter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("章节未找到[书籍ID: %d, 章节ID: %d]", book_id, chapter_id)
			return nil, fmt.Errorf("章节不存在")
		}
		logger.Errorf("获取章节失败[书籍ID: %d, 章节ID: %d]: %v", book_id, chapter_id, err)
		return nil, err
	}

	contentPath := fmt.Sprintf("%s/%d/chap_%05d.txt", BooksDir, chapter.BookID, chapter.ChapterID)
	content, err := os.ReadFile(contentPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("章节内容文件不存在[路径: %s]", contentPath)
			return &chapter, nil
		}
		logger.Errorf("读取章节内容失败[路径: %s]: %v", contentPath, err)
		return nil, err
	}

	chapter.Content = string(content)
	logger.Infof("成功读取章节内容[书籍ID: %d, 章节ID: %d]", book_id, chapter_id) // ✅ 只打业务ID
	return &chapter, nil
}

// AddChapter 添加章节
func (s *BookService) AddChapter(chapter *model.Chapter, content string) error {
	if chapter == nil {
		logger.Errorf("添加章节失败: 章节指针为nil")
		return fmt.Errorf("章节指针为nil")
	}
	if DB == nil {
		logger.Errorf("添加章节失败[书籍ID: %d, 章节ID: %d, 标题: %s]: 数据库连接未初始化",
			chapter.BookID, chapter.ChapterID, chapter.Title)
		return fmt.Errorf("数据库连接未初始化")
	}
	chapter.CreateTime = time.Now()
	chapter.UpdateTime = time.Now()

	if err := DB.Create(chapter).Error; err != nil {
		logger.Errorf("添加章节失败[书籍ID: %d, 章节ID: %d, 标题: %s]: %v",
			chapter.BookID, chapter.ChapterID, chapter.Title, err)
		return err
	}

	contentPath := fmt.Sprintf("%s/%d/chap_%05d.txt", BooksDir, chapter.BookID, chapter.ChapterID)
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		logger.Errorf("保存章节内容失败[路径: %s]: %v", contentPath, err)
		// 回滚：按业务字段删除（避免依赖刚插入的 ID）
		if err2 := DB.Where("book_id = ? AND chapter_id = ?", chapter.BookID, chapter.ChapterID).Delete(&model.Chapter{}).Error; err2 != nil {
			logger.Errorf("回滚章节添加失败[书籍ID: %d, 章节ID: %d]: %v", chapter.BookID, chapter.ChapterID, err2)
		}
		return err
	}

	logger.Infof("添加章节成功[书籍ID: %d, 章节ID: %d, 标题: %s]",
		chapter.BookID, chapter.ChapterID, chapter.Title) // ✅ 不打 ID
	return nil
}

// GetChapterContent 获取章节内容
func (s *BookService) GetChapterContent(bookID, chapterID uint) (string, error) {
	var chapter model.Chapter
	if err := DB.Where("book_id = ? AND chapter_id = ?", bookID, chapterID).First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("章节不存在")
		}
		logger.Errorf("获取章节元数据失败[书籍ID: %d, 章节ID: %d]: %v", bookID, chapterID, err)
		return "", err
	}

	contentPath := fmt.Sprintf("%s/%d/chap_%05d.txt", BooksDir, bookID, chapterID)
	contentBytes, err := os.ReadFile(contentPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("章节内容文件不存在[路径: %s]", contentPath)
			return "", nil
		}
		logger.Errorf("读取章节内容失败[路径: %s]: %v", contentPath, err)
		return "", err
	}
	return string(contentBytes), nil
}

// UpdateChapter 更新章节
func (s *BookService) UpdateChapter(chapter *model.Chapter, content string) error {
	if chapter == nil {
		logger.Errorf("更新章节失败: 章节指针为nil")
		return fmt.Errorf("章节指针为nil")
	}
	if DB == nil {
		logger.Errorf("更新章节失败[书籍ID: %d, 章节ID: %d, 标题: %s]: 数据库连接未初始化",
			chapter.BookID, chapter.ChapterID, chapter.Title)
		return fmt.Errorf("数据库连接未初始化")
	}

	// 先查出数据库中的完整记录（含 ID）
	var dbChapter model.Chapter
	if err := DB.Where("book_id = ? AND chapter_id = ?", chapter.BookID, chapter.ChapterID).First(&dbChapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("章节不存在")
		}
		logger.Errorf("查询章节失败[书籍ID: %d, 章节ID: %d]: %v", chapter.BookID, chapter.ChapterID, err)
		return err
	}

	// 用查到的 ID 来更新（内部使用，不暴露）
	chapter.ID = dbChapter.ID
	chapter.UpdateTime = time.Now()

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(chapter).Error; err != nil {
			logger.Errorf("更新章节失败[书籍ID: %d, 章节ID: %d, 标题: %s]: %v",
				chapter.BookID, chapter.ChapterID, chapter.Title, err)
			return err
		}

		if content != "" {
			contentPath := fmt.Sprintf("%s/%d/chap_%05d.txt", BooksDir, chapter.BookID, chapter.ChapterID)
			if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
				logger.Errorf("更新章节内容失败[路径: %s]: %v", contentPath, err)
				return err
			}
		}

		logger.Infof("更新章节成功[书籍ID: %d, 章节ID: %d, 标题: %s]",
			chapter.BookID, chapter.ChapterID, chapter.Title) // ✅ 不提 ID
		return nil
	})
}

// DeleteChapter 删除章节
func (s *BookService) DeleteChapter(book_id, chapter_id uint) error {
	if DB == nil {
		logger.Errorf("删除章节失败[书籍ID: %d, 章节ID: %d]: 数据库连接未初始化", book_id, chapter_id)
		return fmt.Errorf("数据库连接未初始化")
	}

	// 先查一下是否存在（可选，也可直接删）
	var chapter model.Chapter
	if err := DB.Where("book_id = ? AND chapter_id = ?", book_id, chapter_id).First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("章节不存在")
		}
		logger.Errorf("查询章节失败[书籍ID: %d, 章节ID: %d]: %v", book_id, chapter_id, err)
		return err
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("book_id = ? AND chapter_id = ?", book_id, chapter_id).Delete(&model.Chapter{}).Error; err != nil {
			logger.Errorf("删除章节失败[书籍ID: %d, 章节ID: %d]: %v", book_id, chapter_id, err)
			return err
		}

		contentPath := fmt.Sprintf("%s/%d/chap_%05d.txt", BooksDir, book_id, chapter_id)
		if err := os.Remove(contentPath); err != nil && !os.IsNotExist(err) {
			logger.Warnf("删除章节内容文件失败[路径: %s]: %v", contentPath, err)
		}

		logger.Infof("删除章节成功[书籍ID: %d, 章节ID: %d, 标题: %s]",
			book_id, chapter_id, chapter.Title) // ✅ 只用业务ID
		return nil
	})
}

// 用户服务
// 用户登录
func (s *BookService) Login(username, password string) (*model.User, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("用户登录失败[用户名: %s]: 数据库连接未初始化", username)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var user model.User
	err := DB.Where("username = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		logger.Errorf("用户登录失败[用户名: %s]: %v", username, err)
		return nil, err
	}
	logger.Infof("用户登录成功[用户名: %s, 用户ID: %d]", username, user.ID)
	return &user, nil
}

// 用户注册
func (s *BookService) Register(username, password, nickname, avatar, userType, desc string, level, sex int, isVip bool, location string, status int, email string) (*model.User, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("用户注册失败[用户名: %s]: 数据库连接未初始化", username)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	user := &model.User{
		Username: username,
		Password: password,
		Nickname: nickname,
		Avatar:   avatar,
		Type:     userType,
		Desc:     desc,
		Level:    level,
		Sex:      sex,
		IsVip:    isVip,
		Location: location,
		Status:   status,
		Email:    email,
	}
	if err := DB.Create(user).Error; err != nil {
		logger.Errorf("用户注册失败[用户名: %s]: %v", username, err)
		return nil, err
	}
	logger.Infof("用户注册成功[用户名: %s, 用户ID: %d]", username, user.ID)
	return user, nil
}

// 根据ID获取用户信息
func (s *BookService) GetUserByID(userID uint) (*model.User, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取用户信息失败[用户ID: %d]: 数据库连接未初始化", userID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var user model.User
	err := DB.First(&user, userID).Error
	if err != nil {
		logger.Errorf("获取用户信息失败[用户ID: %d]: %v", userID, err)
		return nil, err
	}
	return &user, nil
}

// 根据用户名获取用户信息
func (s *BookService) GetUserByName(name string) (*model.User, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取用户信息失败[用户名: %s]: 数据库连接未初始化", name)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var user model.User
	err := DB.First(&user, "username = ? or nickname = ?", name, name).Error
	if err != nil {
		logger.Errorf("获取用户信息失败[用户名: %s]: %v", name, err)
		return nil, err
	}
	return &user, nil
}

// 更新用户信息
func (s *BookService) UpdateUser(user *model.User) error {
	// 检查指针是否为nil
	if user == nil {
		logger.Errorf("更新用户信息失败: 用户指针为nil")
		return fmt.Errorf("用户指针为nil")
	}
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("更新用户信息失败[ID: %d]: 数据库连接未初始化", user.ID)
		return fmt.Errorf("数据库连接未初始化")
	}
	user.UpdateTime = time.Now()
	if err := DB.Save(user).Error; err != nil {
		logger.Errorf("更新用户信息失败[ID: %d]: %v", user.ID, err)
		return err
	}
	logger.Infof("更新用户信息成功[ID: %d]", user.ID)
	return nil
}

// 根据作者ID获取作品列表
func (s *BookService) GetBooksByAuthorID(userID uint) ([]model.Book, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取作者作品失败: 数据库连接未初始化")
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	// 获取用户信息
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	// 查询该作者的所有作品
	var books []model.Book
	err = DB.Where("author = ?", user.Nickname).Order("create_time desc").Find(&books).Error
	if err != nil {
		logger.Errorf("获取作者作品失败[作者ID: %d]: %v", userID, err)
		return nil, err
	}
	return books, nil
}

// 获取作者统计数据
func (s *BookService) GetAuthorStats(userID uint) (map[string]interface{}, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取作者统计数据失败: 数据库连接未初始化")
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	// 获取用户信息
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	// 获取作品数量
	var bookCount int64
	DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Count(&bookCount)
	// 获取总点击量
	var totalClicks int64
	DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Select("COALESCE(SUM(click_count), 0)").Scan(&totalClicks)
	// 获取总字数
	var totalWords int64
	DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Select("COALESCE(SUM(word_count), 0)").Scan(&totalWords)
	// 获取最新更新时间
	var latestUpdate time.Time
	DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Order("update_time desc").Select("update_time").First(&latestUpdate)
	// 返回统计数据
	stats := map[string]interface{}{
		"bookCount":    bookCount,
		"totalClicks":  totalClicks,
		"totalWords":   totalWords,
		"latestUpdate": latestUpdate,
	}
	return stats, nil
}

// 书架服务
// 获取用户书架
func (s *BookService) GetUserShelf(userID uint) ([]model.Shelf, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取用户书架失败[用户ID: %d]: 数据库连接未初始化", userID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var shelves []model.Shelf
	err := DB.Preload("Book").Where("user_id = ?", userID).Find(&shelves).Error
	if err != nil {
		logger.Errorf("获取用户书架失败[用户ID: %d]: %v", userID, err)
		return nil, err
	}
	return shelves, nil
}

// 添加到书架
func (s *BookService) AddToShelf(userID, bookID uint) error {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("添加到书架失败[用户ID: %d, 书籍ID: %d]: 数据库连接未初始化", userID, bookID)
		return fmt.Errorf("数据库连接未初始化")
	}
	// 检查是否已存在
	var existing model.Shelf
	err := DB.Where("user_id = ? AND book_id = ?", userID, bookID).First(&existing).Error
	if err == nil {
		// 已存在
		return nil
	}

	// 添加到书架
	shelf := model.Shelf{
		UserID:   userID,
		BookID:   bookID,
		AddTime:  time.Now(),
		LastRead: time.Now(),
	}
	err = DB.Create(&shelf).Error
	if err != nil {
		logger.Errorf("添加到书架失败[用户ID: %d, 书籍ID: %d]: %v", userID, bookID, err)
	}
	return err
}

// 从书架移除
func (s *BookService) RemoveFromShelf(userID, bookID uint) error {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("从书架移除失败[用户ID: %d, 书籍ID: %d]: 数据库连接未初始化", userID, bookID)
		return fmt.Errorf("数据库连接未初始化")
	}
	if err := DB.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&model.Shelf{}).Error; err != nil {
		logger.Errorf("从书架移除失败[用户ID: %d, 书籍ID: %d]: %v", userID, bookID, err)
		return err
	}
	logger.Infof("从书架移除成功[用户ID: %d, 书籍ID: %d]", userID, bookID)
	return nil
}

// 阅读历史服务
// 获取用户阅读历史
func (s *BookService) GetUserHistory(userID uint) ([]model.History, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("获取用户阅读历史失败[用户ID: %d]: 数据库连接未初始化", userID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var histories []model.History
	err := DB.Preload("Book").Preload("Chapter").Where("user_id = ?", userID).Order("update_time desc").Find(&histories).Error
	if err != nil {
		logger.Errorf("获取用户阅读历史失败[用户ID: %d]: %v", userID, err)
		return nil, err
	}
	return histories, nil
}

// 更新阅读进度
func (s *BookService) UpdateReadingProgress(userID, bookID, chapterID uint, readingProgress int) error {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("更新阅读进度失败[用户ID: %d, 书籍ID: %d]: 数据库连接未初始化", userID, bookID)
		return fmt.Errorf("数据库连接未初始化")
	}
	// 查找是否已存在阅读记录
	var history model.History
	err := DB.Where("user_id = ? AND book_id = ?", userID, bookID).First(&history).Error
	if err != nil {
		// 如果不存在，创建新记录
		history = model.History{
			UserID:          userID,
			BookID:          bookID,
			ChapterID:       chapterID,
			ReadingProgress: readingProgress,
		}
		if err := DB.Create(&history).Error; err != nil {
			logger.Errorf("创建阅读记录失败[用户ID: %d, 书籍ID: %d]: %v", userID, bookID, err)
			return err
		}
		logger.Infof("创建阅读记录成功[用户ID: %d, 书籍ID: %d, 章节ID: %d]", userID, bookID, chapterID)
	} else {
		// 如果存在，更新记录
		history.ChapterID = chapterID
		history.ReadingProgress = readingProgress
		if err := DB.Save(&history).Error; err != nil {
			logger.Errorf("更新阅读进度失败[用户ID: %d, 书籍ID: %d]: %v", userID, bookID, err)
			return err
		}
		logger.Infof("更新阅读进度成功[用户ID: %d, 书籍ID: %d, 章节ID: %d]", userID, bookID, chapterID)
	}
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
