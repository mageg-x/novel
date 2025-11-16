package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mageg-x/novel/src/handler"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 设置Gin为发布模式，禁止调试日志输出
	gin.SetMode(gin.ReleaseMode)
	// 创建默认的Gin引擎
	r := gin.Default()

	// 允许跨域请求
	r.Use(corsMiddleware())

	// API路由组
	api := r.Group("/api")
	{
		api.GET("/categories/:category/books", handler.GetBooksByCategory)

		// 推荐相关路由
		api.GET("/rcmd/:type", handler.GetRcmdBooks)

		// 书籍相关路由
		api.GET("/books", handler.GetAllBooks)
		api.GET("/books/:book_id", handler.GetBookByID)
		api.POST("/books", handler.AddBook)
		api.PUT("/books/:book_id", handler.UpdateBook)
		api.DELETE("/books/:book_id", handler.DeleteBook)

		// 关联书籍
		api.GET("/books/:book_id/related", handler.GetRelatedBooks)

		// 书籍评论
		api.GET("/books/:book_id/comments", handler.GetBookComments)

		// 章节相关路由
		api.GET("/books/:book_id/chapters", handler.GetBookChapters)
		api.GET("/books/:book_id/chapters/:chapter_id", handler.GetChapterByID)
		api.POST("/books/:book_id/chapters/:chapter_id", handler.AddChapter)
		api.PUT("/books/:book_id/chapters/:chapter_id", handler.UpdateChapter)
		api.DELETE("/books/:book_id/chapters/:chapter_id", handler.DeleteChapter)

		// 排行榜相关路由
		api.GET("/ranks/:type", handler.GetRanks)

		// 用户相关路由
		api.POST("/users/login", handler.Login)
		api.POST("/users/register", handler.Register)
		api.GET("/users/:user_id", handler.AuthMiddleware(), handler.GetUserByID)
		api.PUT("/users/:user_id", handler.AuthMiddleware(), handler.UpdateUser)
		api.GET("/users/info/:name", handler.GetUserByName)

		// 书架相关路由（需要登录）
		api.GET("/users/:user_id/shelf", handler.AuthMiddleware(), handler.GetUserShelf)
		api.POST("/users/:user_id/shelf", handler.AuthMiddleware(), handler.AddToShelf)
		api.DELETE("/users/:user_id/shelf/:book_id", handler.AuthMiddleware(), handler.RemoveFromShelf)

		// 阅读历史相关路由（需要登录）
		api.GET("/users/:user_id/history", handler.AuthMiddleware(), handler.GetUserHistory)
		api.POST("/users/:user_id/books/:book_id/progress", handler.AuthMiddleware(), handler.UpdateReadingProgress)

		// 作者相关路由
		api.GET("/author/info/:id", handler.GetAuthorInfo)
		api.GET("/author/books/:id", handler.GetAuthorBooks)
		api.GET("/author/stats/:id", handler.GetAuthorStats)

		// 作者操作相关路由（需要登录且为作者身份）
		api.POST("/author/books", handler.AuthMiddleware(), handler.AuthorMiddleware(), handler.AddBook)
		api.PUT("/author/books/:book_id", handler.AuthMiddleware(), handler.AuthorMiddleware(), handler.UpdateBook)
		api.DELETE("/author/books/:book_id", handler.AuthMiddleware(), handler.AuthorMiddleware(), handler.DeleteBook)
		api.POST("/author/books/:book_id/chapters", handler.AuthMiddleware(), handler.AuthorMiddleware(), handler.AddChapter)
		api.PUT("/author/books/:book_id/chapters/:chapter_id", handler.AuthMiddleware(), handler.AuthorMiddleware(), handler.UpdateChapter)
		api.DELETE("/author/books/:book_id/chapters/:chapter_id", handler.AuthMiddleware(), handler.AuthorMiddleware(), handler.DeleteChapter)
	}

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return r
}

// corsMiddleware 跨域中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
