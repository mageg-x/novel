package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mageg-x/novel/src/util"
)

// 用户登录
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 搜索用户
func SearchUsers(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}

	keyword := c.Query("keyword")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, total, err := ss.SearchUsers(keyword, pageSize, offset)
	if err != nil {
		logger.Errorf("Failed to search users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to search users",
			"data":    nil,
		})
		return
	}

	// 移除每个user对象的password字段
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"users": users,
			"total": total,
		},
	})
}

// 搜索评论
func SearchComments(c *gin.Context) {
	// 检查SearchService指针是否为nil
	if ss == nil {
		logger.Errorf("SearchService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Search service not initialized",
			"data":    nil,
		})
		return
	}

	keyword := c.Query("keyword")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	comments, total, err := ss.SearchComments(keyword, pageSize, offset)
	if err != nil {
		logger.Errorf("Failed to search comments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to search comments",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"comments": comments,
			"total":    total,
		},
	})
}

func Login(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("User login parameter error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	user, err := us.Login(req.Username, req.Password)
	if err != nil {
		logger.Errorf("User login failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid username or password",
			"data":    nil,
		})
		return
	}

	// 生成JWT令牌
	token, err := util.GenerateToken(user)
	if err != nil {
		logger.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to generate token",
			"data":    nil,
		})
		return
	}

	// 移除password字段
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

// 用户注册
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
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
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Parameter error %+v: %v", req, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	user, err := us.Register(req.Username, req.Password, req.Nickname, req.Avatar, req.Type, req.Desc, req.Level, req.Sex, req.IsVip, req.Location, req.Status, req.Email)
	if err != nil {
		logger.Errorf("User registration failed: %v", err)
		// 根据错误类型返回不同的HTTP状态码
		if err.Error() == "用户名已存在" {
			c.JSON(http.StatusConflict, gin.H{
				"code":    http.StatusConflict,
				"message": "Username already exists",
				"data":    nil,
			})
		} else if err.Error() == "参数错误" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid parameters",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
				"data":    nil,
			})
		}
		return
	}

	// 移除password字段
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    user,
	})
}

// 获取用户书架
func GetUserShelf(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Get user shelf - invalid user ID: %s, error: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	shelves, err := us.GetUserShelf(uint(userID))
	if err != nil {
		logger.Errorf("Failed to get user shelf: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get user shelf",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    shelves,
	})
}

// 添加书籍到书架
type AddShelfRequest struct {
	BookID uint `json:"bookId" binding:"required"`
}

func AddToShelf(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Add book to shelf - invalid user ID: %s, error: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	var req AddShelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Add book to shelf - parameter error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	if err := us.AddToShelf(uint(userID), req.BookID); err != nil {
		logger.Errorf("Failed to add book to shelf: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add book to shelf",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    nil,
	})
}

// 从书架移除书籍
func RemoveFromShelf(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Remove book from shelf - invalid user ID: %s, error: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Remove book from shelf - invalid book ID: %s, error: %v", bookIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	if err := us.RemoveFromShelf(uint(userID), uint(bookID)); err != nil {
		logger.Errorf("Failed to remove book from shelf: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to remove book from shelf",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    nil,
	})
}

// 获取用户阅读历史
func GetUserHistory(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Get user reading history - invalid user ID: %s, error: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	histories, err := us.GetUserHistory(uint(userID))
	if err != nil {
		logger.Errorf("Failed to get user reading history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get user reading history",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    histories,
	})
}

// 更新阅读进度
type UpdateReadingProgressRequest struct {
	ChapterNo       uint `json:"chapterNo" binding:"required"`
	ReadingProgress int  `json:"readingProgress" binding:"min=0,max=100"`
}

func UpdateReadingProgress(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		logger.Errorf("BookService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Update reading progress - invalid user ID: %s, error: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Update reading progress - invalid book ID: %s, error: %v", bookIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	var req UpdateReadingProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Update reading progress - parameter error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	if err := us.UpdateReadingProgress(uint(userID), uint(bookID), req.ChapterNo, req.ReadingProgress); err != nil {
		logger.Errorf("Failed to update reading progress: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update reading progress",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    nil,
	})
}

// 获取用户信息
func GetUserByID(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Get user info - invalid user ID: %s, error: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	user, err := us.GetUserByID(uint(userID))
	if err != nil {
		logger.Errorf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get user info",
			"data":    nil,
		})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    user,
	})
}

// 根据用户名获取用户信息
func GetUserByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Username cannot be empty",
			"data":    nil,
		})
		return
	}

	user, err := us.GetUserByName(name)
	if err != nil {
		logger.Errorf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get user info",
			"data":    nil,
		})
		return
	}
	// 移除password字段
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    user,
	})
}

// 更新用户信息
func UpdateUser(c *gin.Context) {
	// 检查UserService指针是否为nil
	if us == nil {
		logger.Errorf("UserService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User service not initialized",
			"data":    nil,
		})
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Update user info - invalid user ID: %s, error: %v", userIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	// 获取现有用户
	user, err := us.GetUserByID(uint(userID))
	if err != nil {
		logger.Errorf("Update user info - failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get user info",
			"data":    nil,
		})
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
		logger.Errorf("Update user info - parameter error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
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

	if err := us.UpdateUser(user); err != nil {
		logger.Errorf("Failed to update user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update user info",
			"data":    nil,
		})
		return
	}

	// 移除password字段
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    user,
	})
}
