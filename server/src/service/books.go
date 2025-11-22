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
	// 启用查询缓存
	DB.Exec("PRAGMA cache_size = -8000;") // 8MB缓存
	DB.Exec("PRAGMA synchronous = NORMAL;")

	// 自动迁移数据库表结构
	if err := DB.AutoMigrate(&model.Book{}, &model.Volume{}, &model.Chapter{}, &model.Rcmd{}, &model.Rank{}, &model.User{}, &model.Shelf{}, &model.History{}, &model.Comment{}); err != nil {
		logger.Errorf("数据库迁移失败: %v", err)
		return err
	}

	// 索引已在模型定义中通过gorm标签指定(index:idx_books_update_time, index:idx_books_category)
	// AutoMigrate会自动创建这些索引，无需手动执行SQL语句

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
	// 移除或减少预加载以提升性能，Volumes和Chapters可以在需要时单独查询
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
	err := DB.Preload("Volumes").Preload("Chapters").First(&book, id).Error
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

		// 删除书籍的卷
		if err := tx.Where("book_id = ?", id).Delete(&model.Volume{}).Error; err != nil {
			logger.Errorf("删除书籍卷失败[书籍ID: %d]: %v", id, err)
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

// 卷服务
// 添加卷
func (s *BookService) AddVolume(volume *model.Volume) error {
	if volume == nil {
		logger.Errorf("添加卷失败: 卷指针为nil")
		return fmt.Errorf("卷指针为nil")
	}
	if DB == nil {
		logger.Errorf("添加卷失败[书籍ID: %d, 卷序号: %d, 标题: %s]: 数据库连接未初始化",
			volume.BookID, volume.VolumeNo, volume.Title)
		return fmt.Errorf("数据库连接未初始化")
	}

	err := DB.Create(volume).Error
	if err != nil {
		logger.Errorf("添加卷失败[书籍ID: %d, 卷序号: %d, 标题: %s]: %v",
			volume.BookID, volume.VolumeNo, volume.Title, err)
		return err
	}

	logger.Infof("添加卷成功[书籍ID: %d, 卷序号: %d, 标题: %s]",
		volume.BookID, volume.VolumeNo, volume.Title)
	return nil
}

// 获取书籍的所有卷
func (s *BookService) GetVolumesByBookID(bookID uint) ([]model.Volume, error) {
	if DB == nil {
		logger.Errorf("获取书籍卷失败[书籍ID: %d]: 数据库连接未初始化", bookID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	var volumes []model.Volume
	err := DB.Where("book_id = ?", bookID).Order("volume_no ASC").Find(&volumes).Error
	if err != nil {
		logger.Errorf("获取书籍卷失败[书籍ID: %d]: %v", bookID, err)
	}
	return volumes, err
}

// 根据ID获取卷
func (s *BookService) GetVolumeByNo(bookID, volumeNo uint) (*model.Volume, error) {
	if DB == nil {
		logger.Errorf("获取卷失败[书籍ID: %d, 卷序号: %d]: 数据库连接未初始化", bookID, volumeNo)
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	var volume model.Volume
	err := DB.Where("book_id = ? AND volume_no = ?", bookID, volumeNo).First(&volume).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("卷未找到[书籍ID: %d, 卷序号: %d]", bookID, volumeNo)
			return nil, fmt.Errorf("卷不存在")
		}
		logger.Errorf("获取卷失败[书籍ID: %d, 卷序号: %d]: %v", bookID, volumeNo, err)
		return nil, err
	}

	return &volume, nil
}

// 更新卷信息
func (s *BookService) UpdateVolume(volume *model.Volume) error {
	if volume == nil {
		logger.Errorf("更新卷失败: 卷指针为nil")
		return fmt.Errorf("卷指针为nil")
	}
	if DB == nil {
		logger.Errorf("更新卷失败[书籍ID: %d, 卷序号: %d, 标题: %s]: 数据库连接未初始化",
			volume.BookID, volume.VolumeNo, volume.Title)
		return fmt.Errorf("数据库连接未初始化")
	}

	// 先查出数据库中的完整记录（含 ID）
	var dbVolume model.Volume
	if err := DB.Where("book_id = ? AND volume_no = ?", volume.BookID, volume.VolumeNo).First(&dbVolume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("卷不存在")
		}
		logger.Errorf("查询卷失败[书籍ID: %d, 卷序号: %d]: %v", volume.BookID, volume.VolumeNo, err)
		return err
	}

	// 用查到的 ID 来更新（内部使用，不暴露）
	volume.ID = dbVolume.ID

	if err := DB.Save(volume).Error; err != nil {
		logger.Errorf("更新卷失败[书籍ID: %d, 卷序号: %d, 标题: %s]: %v",
			volume.BookID, volume.VolumeNo, volume.Title, err)
		return err
	}

	logger.Infof("更新卷成功[书籍ID: %d, 卷序号: %d, 标题: %s]",
		volume.BookID, volume.VolumeNo, volume.Title)
	return nil
}

// 删除卷
func (s *BookService) DeleteVolume(bookID, volumeNo uint) error {
	if DB == nil {
		logger.Errorf("删除卷失败[书籍ID: %d, 卷序号: %d]: 数据库连接未初始化", bookID, volumeNo)
		return fmt.Errorf("数据库连接未初始化")
	}

	// 先查一下是否存在（可选，也可直接删）
	var volume model.Volume
	if err := DB.Where("book_id = ? AND volume_no = ?", bookID, volumeNo).First(&volume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("卷不存在")
		}
		logger.Errorf("查询卷失败[书籍ID: %d, 卷序号: %d]: %v", bookID, volumeNo, err)
		return err
	}

	// 开启事务删除卷并更新相关章节的卷ID为NULL
	return DB.Transaction(func(tx *gorm.DB) error {
		// 将关联到此卷的章节的VolumeNo设为NULL
		if err := tx.Model(&model.Chapter{}).Where("book_id = ? AND volume_no = ?", bookID, volume.ID).Update("volume_no", nil).Error; err != nil {
			logger.Errorf("更新章节卷关联失败[书籍ID: %d, 卷ID: %d]: %v", bookID, volume.ID, err)
			return err
		}

		// 删除卷
		if err := tx.Where("book_id = ? AND volume_no = ?", bookID, volumeNo).Delete(&model.Volume{}).Error; err != nil {
			logger.Errorf("删除卷失败[书籍ID: %d, 卷序号: %d]: %v", bookID, volumeNo, err)
			return err
		}

		logger.Infof("删除卷成功[书籍ID: %d, 卷序号: %d, 标题: %s]",
			bookID, volumeNo, volume.Title)
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

// GetChaptersByVolumeNo 获取某一卷的所有章节
func (s *BookService) GetChaptersByVolumeNo(bookID, volumeNo uint) ([]model.Chapter, error) {
	if DB == nil {
		logger.Errorf("获取卷章节失败[书籍ID: %d, 卷序号: %d]: 数据库连接未初始化", bookID, volumeNo)
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	// 先获取卷的数据库ID
	var volume model.Volume
	if err := DB.Where("book_id = ? AND volume_no = ?", bookID, volumeNo).First(&volume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("卷未找到[书籍ID: %d, 卷序号: %d]", bookID, volumeNo)
			return nil, fmt.Errorf("卷不存在")
		}
		logger.Errorf("获取卷信息失败[书籍ID: %d, 卷序号: %d]: %v", bookID, volumeNo, err)
		return nil, err
	}

	var chapters []model.Chapter
	err := DB.Where("book_id = ? AND volume_no = ?", bookID, volume.ID).Order("chapter_id ASC").Find(&chapters).Error
	if err != nil {
		logger.Errorf("获取卷章节失败[书籍ID: %d, 卷ID: %d]: %v", bookID, volume.ID, err)
	}
	return chapters, err
}

// GetChapterByID 根据书籍ID和章节序号获取章节
func (s *BookService) GetChapterByID(bookID, volumeNo, chapterID uint) (*model.Chapter, error) {
	if DB == nil {
		logger.Errorf("获取章节失败[书籍ID: %d, 卷序号: %d, 章节No: %d]: 数据库连接未初始化", bookID, volumeNo, chapterID)
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	var chapter model.Chapter
	err := DB.Where("book_id = ? AND volume_no = ? AND chapter_id = ?", bookID, volumeNo, chapterID).First(&chapter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("章节未找到[书籍ID: %d, 卷序号: %d, 章节ID: %d]", bookID, volumeNo, chapterID)
			return nil, fmt.Errorf("章节不存在")
		}
		logger.Errorf("获取章节失败[书籍ID: %d, 卷序号: %d, 章节ID: %d]: %v", bookID, volumeNo, chapterID, err)
		return nil, err
	}

	contentPath := fmt.Sprintf("%s/%d/chap_%03d_%05d.txt", BooksDir, chapter.BookID, chapter.VolumeNo, chapter.ChapterID)
	content, err := os.ReadFile(contentPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("章节内容文件不存在[路径: %s, 书籍ID: %d, 卷序号: %d, 章节ID: %d]", contentPath, bookID, volumeNo, chapterID)
			return &chapter, nil
		}
		logger.Errorf("读取章节内容失败[路径: %s, 书籍ID: %d, 卷序号: %d, 章节ID: %d]: %v", contentPath, bookID, volumeNo, chapterID, err)
		return nil, err
	}

	chapter.Content = string(content)
	logger.Infof("成功读取章节内容[书籍ID: %d, 卷序号: %d, 章节ID: %d]", bookID, volumeNo, chapterID) // ✅ 只打业务ID
	return &chapter, nil
}

// AddChapter 添加章节
func (s *BookService) AddChapter(chapter *model.Chapter, content string) error {
	if chapter == nil {
		logger.Errorf("添加章节失败: 章节指针为nil")
		return fmt.Errorf("章节指针为nil")
	}
	if DB == nil {
		logger.Errorf("添加章节失败[书籍ID: %d, 卷序号: %d, 章节ID: %d, 标题: %s]: 数据库连接未初始化",
			chapter.BookID, chapter.VolumeNo, chapter.ChapterID, chapter.Title)
		return fmt.Errorf("数据库连接未初始化")
	}

	chapter.CreateTime = time.Now()
	chapter.UpdateTime = time.Now()

	if err := DB.Create(chapter).Error; err != nil {
		logger.Errorf("添加章节失败[书籍ID: %d, 卷ID: %d, 章节ID: %d, 标题: %s]: %v",
			chapter.BookID, chapter.VolumeNo, chapter.ChapterID, chapter.Title, err)
		return err
	}

	contentPath := fmt.Sprintf("%s/%d/chap_%03d_%05d.txt", BooksDir, chapter.BookID, chapter.VolumeNo, chapter.ChapterID)
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		logger.Errorf("保存章节内容失败[路径: %s, 书籍ID: %d, 卷序号: %d, 章节ID: %d]: %v", contentPath, chapter.BookID, chapter.VolumeNo, chapter.ChapterID, err)
		// 回滚：按业务字段删除（避免依赖刚插入的 ID）
		if err2 := DB.Where("book_id = ? AND chapter_id = ?", chapter.BookID, chapter.ChapterID).Delete(&model.Chapter{}).Error; err2 != nil {
			logger.Errorf("回滚章节添加失败[书籍ID: %d, 章节ID: %d]: %v", chapter.BookID, chapter.ChapterID, err2)
		}
		return err
	}

	logger.Infof("添加章节成功[书籍ID: %d, 卷序号: %d, 章节ID: %d, 标题: %s]", chapter.BookID, chapter.VolumeNo, chapter.ChapterID, chapter.Title) // ✅ 不打 ID
	return nil
}

// GetChapterContent 获取章节内容
func (s *BookService) GetChapterContent(bookID, volumeNo, chapterID uint) (string, error) {
	var chapter model.Chapter
	if err := DB.Where("book_id = ? AND volume_no = ? AND chapter_id = ?", bookID, volumeNo, chapterID).First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("章节不存在")
		}
		logger.Errorf("获取章节元数据失败[书籍ID: %d, 卷序号: %d, 章节ID: %d]: %v", bookID, volumeNo, chapterID, err)
		return "", err
	}

	contentPath := fmt.Sprintf("%s/%d/chap_%03d_%05d.txt", BooksDir, bookID, volumeNo, chapterID)
	contentBytes, err := os.ReadFile(contentPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("章节内容文件不存在[路径: %s, 书籍ID: %d, 卷序号: %d, 章节ID: %d]", contentPath, bookID, volumeNo, chapterID)
			return "", nil
		}
		logger.Errorf("读取章节内容失败[路径: %s, 书籍ID: %d, 卷序号: %d, 章节ID: %d]: %v", contentPath, bookID, volumeNo, chapterID, err)
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
	if err := DB.Where("book_id = ? AND volume_no = ? AND chapter_id = ?", chapter.BookID, chapter.VolumeNo, chapter.ChapterID).First(&dbChapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("章节不存在")
		}
		logger.Errorf("查询章节失败[书籍ID: %d, 卷号： %d, 章节ID: %d]: %v", chapter.BookID, chapter.VolumeNo, chapter.ChapterID, err)
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
			contentPath := fmt.Sprintf("%s/%d/chap_%03d_%05d.txt", BooksDir, chapter.BookID, chapter.VolumeNo, chapter.ChapterID)
			if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
				logger.Errorf("更新章节内容失败[路径: %s, 书籍ID: %d, 章节ID: %d]: %v", contentPath, chapter.BookID, chapter.ChapterID, err)
				return err
			}
		}

		logger.Infof("更新章节成功[书籍ID: %d, 章节ID: %d, 标题: %s]",
			chapter.BookID, chapter.ChapterID, chapter.Title) // ✅ 不提 ID
		return nil
	})
}

// DeleteChapter 删除章节
func (s *BookService) DeleteChapter(bookID, volumeNo, chapterID uint) error {
	if DB == nil {
		logger.Errorf("删除章节失败[书籍ID: %d, 卷序号: %d, 章节ID: %d]: 数据库连接未初始化", bookID, volumeNo, chapterID)
		return fmt.Errorf("数据库连接未初始化")
	}

	// 先查一下是否存在（可选，也可直接删）
	var chapter model.Chapter
	if err := DB.Where("book_id = ? AND volume_no = ? AND chapter_id = ?", bookID, volumeNo, chapterID).First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("章节不存在")
		}
		logger.Errorf("查询章节失败[书籍ID: %d, 卷ID: %d, 章节ID: %d]: %v", bookID, volumeNo, chapterID, err)
		return err
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("book_id = ? AND volume_no = ? AND chapter_id = ?", bookID, volumeNo, chapterID).Delete(&model.Chapter{}).Error; err != nil {
			logger.Errorf("删除章节失败[书籍ID: %d, 卷ID: %d, 章节ID: %d]: %v", bookID, volumeNo, chapterID, err)
			return err
		}

		contentPath := fmt.Sprintf("%s/%d/chap_%03d_%05d.txt", BooksDir, bookID, volumeNo, chapterID)
		if err := os.Remove(contentPath); err != nil && !os.IsNotExist(err) {
			logger.Warnf("删除章节内容文件失败[路径: %s, 书籍ID: %d, 卷ID: %d, 章节ID: %d]: %v", contentPath, bookID, volumeNo, chapterID, err)
		}

		logger.Infof("删除章节成功[书籍ID: %d, 卷ID: %d, 章节ID: %d, 标题: %s]",
			bookID, volumeNo, chapterID, chapter.Title) // ✅ 只用业务ID
		return nil
	})
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
