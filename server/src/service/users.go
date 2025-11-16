package service

import (
	"fmt"
	"time"

	"github.com/mageg-x/novel/src/model"
)

// 用户服务
type UserService struct{}

// 用户登录
func (s *UserService) Login(username, password string) (*model.User, error) {
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
func (s *UserService) Register(username, password, nickname, avatar, userType, desc string, level, sex int, isVip bool, location string, status int, email string) (*model.User, error) {
	// 检查数据库连接是否已初始化
	if DB == nil {
		logger.Errorf("用户注册失败[用户名: %s]: 数据库连接未初始化", username)
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	// 检查用户名是否已存在
	var existingUser model.User
	if err := DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		logger.Errorf("用户注册失败[用户名: %s]: 用户名已存在", username)
		return nil, fmt.Errorf("用户名已存在")
	}
	// 设置默认值
	if nickname == "" {
		nickname = username
	}
	if avatar == "" {
		avatar = "https://api.dicebear.com/9.x/avataaars/svg?seed=default"
	}
	if desc == "" {
		desc = "这家伙很懒，什么都没有留下"
	}
	if location == "" {
		location = "中国"
	}
	if email == "" {
		email = fmt.Sprintf("%s@mageg.cn", username)
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
		// 检查是否是唯一约束错误
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return nil, fmt.Errorf("用户名已存在")
		}
		return nil, fmt.Errorf("注册失败，请稍后重试")
	}
	logger.Infof("用户注册成功[用户名: %s, 用户ID: %d]", username, user.ID)
	return user, nil
}

// 根据ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
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
func (s *UserService) GetUserByName(name string) (*model.User, error) {
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
func (s *UserService) UpdateUser(user *model.User) error {
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
func (s *UserService) GetBooksByAuthorID(userID uint) ([]model.Book, error) {
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
func (s *UserService) GetAuthorStats(userID uint) (map[string]interface{}, error) {
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
	if err := DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Count(&bookCount).Error; err != nil {
		logger.Errorf("获取作者作品数量失败[作者ID: %d]: %v", userID, err)
		bookCount = 0
	}
	// 获取总点击量
	var totalClicks int64
	if err := DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Select("COALESCE(SUM(click_count), 0)").Scan(&totalClicks).Error; err != nil {
		logger.Errorf("获取作者总点击量失败[作者ID: %d]: %v", userID, err)
		totalClicks = 0
	}
	// 获取总字数
	var totalWords int64
	if err := DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Select("COALESCE(SUM(word_count), 0)").Scan(&totalWords).Error; err != nil {
		logger.Errorf("获取作者总字数失败[作者ID: %d]: %v", userID, err)
		totalWords = 0
	}
	// 获取最新更新时间
	var latestUpdate time.Time
	if err := DB.Model(&model.Book{}).Where("author = ?", user.Nickname).Order("update_time desc").Select("update_time").First(&latestUpdate).Error; err != nil {
		// 如果没有作品，latestUpdate将是零值
		logger.Debugf("作者[%d]没有作品或获取最新更新时间失败: %v", userID, err)
	}

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
func (s *UserService) GetUserShelf(userID uint) ([]model.Shelf, error) {
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
func (s *UserService) AddToShelf(userID, bookID uint) error {
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
func (s *UserService) RemoveFromShelf(userID, bookID uint) error {
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
func (s *UserService) GetUserHistory(userID uint) ([]model.History, error) {
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
func (s *UserService) UpdateReadingProgress(userID, bookID, chapterID uint, readingProgress int) error {
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
