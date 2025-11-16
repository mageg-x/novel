package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户登录
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user, err := bs.Login(req.Username, req.Password)
	if err != nil {
		logger.Errorf("用户登录失败: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    user,
	})
}

// 用户注册
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Type     string `json:"type"`
	Desc     string `json:"desc"`
	Level    int    `json:"level"`
	Sex      int    `json:"sex"`
	IsVip    bool   `json:"isVip"`
	Location string `json:"location"`
	Status   int    `json:"status"`
	Email    string `json:"email"`
}

func Register(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user, err := bs.Register(req.Username, req.Password, req.Nickname, req.Avatar, req.Type, req.Desc, req.Level, req.Sex, req.IsVip, req.Location, req.Status, req.Email)
	if err != nil {
		logger.Errorf("用户注册失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data":    user,
	})
}

// 获取用户书架
func GetUserShelf(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	shelves, err := bs.GetUserShelf(uint(userID))
	if err != nil {
		logger.Errorf("获取用户书架失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户书架失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    shelves,
	})
}

// 添加书籍到书架
type AddShelfRequest struct {
	BookID uint `json:"bookId" binding:"required"`
}

func AddToShelf(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req AddShelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := bs.AddToShelf(uint(userID), req.BookID); err != nil {
		logger.Errorf("添加书籍到书架失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加书籍到书架失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加书籍到书架成功",
	})
}

// 从书架移除书籍
func RemoveFromShelf(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	if err := bs.RemoveFromShelf(uint(userID), uint(bookID)); err != nil {
		logger.Errorf("从书架移除书籍失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "从书架移除书籍失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "从书架移除书籍成功",
	})
}

// 获取用户阅读历史
func GetUserHistory(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	histories, err := bs.GetUserHistory(uint(userID))
	if err != nil {
		logger.Errorf("获取用户阅读历史失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户阅读历史失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    histories,
	})
}

// 更新阅读进度
type UpdateReadingProgressRequest struct {
	ChapterID       uint `json:"chapterId" binding:"required"`
	ReadingProgress int  `json:"readingProgress" binding:"min=0,max=100"`
}

func UpdateReadingProgress(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	var req UpdateReadingProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := bs.UpdateReadingProgress(uint(userID), uint(bookID), req.ChapterID, req.ReadingProgress); err != nil {
		logger.Errorf("更新阅读进度失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新阅读进度失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新阅读进度成功",
	})
}

// 获取用户信息
func GetUserByID(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := bs.GetUserByID(uint(userID))
	if err != nil {
		logger.Errorf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}

// 根据用户名获取用户信息
func GetUserByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名不能为空"})
		return
	}

	user, err := bs.GetUserByName(name)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取用户信息成功", "data": user})
}

// 更新用户信息
func UpdateUser(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取现有用户
	user, err := bs.GetUserByID(uint(userID))
	if err != nil {
		logger.Errorf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 绑定请求参数
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Desc     string `json:"desc"`
		Sex      int    `json:"sex"`
		Location string `json:"location"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Desc != "" {
		user.Desc = req.Desc
	}
	user.Sex = req.Sex
	if req.Location != "" {
		user.Location = req.Location
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := bs.UpdateUser(user); err != nil {
		logger.Errorf("更新用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新用户信息成功",
		"data":    user,
	})
}
