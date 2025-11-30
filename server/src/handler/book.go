package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mageg-x/novel/src/log"
	"github.com/mageg-x/novel/src/model"
	"github.com/mageg-x/novel/src/service"
)

var (
	bs     = &service.BookService{}
	us     = &service.UserService{}
	ss     = service.SService
	logger = log.GetLogger("novel")
)

// 获取所有书籍
type BookQueryParams struct {
	Offset int `form:"offset,default=0" binding:"min=0"`
	Limit  int `form:"limit,default=10" binding:"min=1,max=100"`
}

func GetAllBooks(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}

	var params BookQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	books, total, err := bs.GetAllBooks(params.Offset, params.Limit)
	if err != nil {
		logger.Errorf("Failed to get book list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get book list",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"books": books,
			"total": total,
		},
	})
}

// 根据ID获取书籍
func GetBookByID(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	idStr := c.Param("book_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	book, err := bs.GetBookByID(uint(id))
	if err != nil {
		logger.Errorf("Failed to get book details: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Book not found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    book,
	})
}

// 根据分类获取书籍
func GetBooksByCategory(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	category := c.Param("category")

	var params BookQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	books, total, err := bs.GetBooksByCategory(category, params.Offset, params.Limit)
	if err != nil {
		logger.Errorf("Failed to get books by category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get book list",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"books": books,
			"total": total,
		},
	})
}

// 搜索书籍
type BookSearchParams struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
}

func SearchBooks(c *gin.Context) {
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

	var params BookSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Errorf("Invalid parameters: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	offset := (params.Page - 1) * params.PageSize
	books, total, err := ss.SearchBooks(params.Keyword, params.PageSize, offset)
	if err != nil {
		logger.Errorf("Failed to search books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to search books",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"books": books,
			"total": total,
		},
	})
}

// 添加书籍
type BookCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Cover       string `json:"cover"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=serializing completed"`
	WordCount   int    `json:"wordCount"`
}

func AddBook(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	var req BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	book := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Cover:       req.Cover,
		Category:    req.Category,
		Description: req.Description,
		Status:      req.Status,
		WordCount:   req.WordCount,
	}

	if err := bs.AddBook(book); err != nil {
		logger.Errorf("Failed to add book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add book",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully added book",
		"data":    book,
	})
}

// 更新书籍
type BookUpdateRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Cover       string `json:"cover"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=serializing completed"`
	WordCount   int    `json:"wordCount"`
}

func UpdateBook(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	idStr := c.Param("book_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	// 获取现有书籍
	book, err := bs.GetBookByID(uint(id))
	if err != nil {
		logger.Errorf("Failed to get book details: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Book not found",
			"data":    nil,
		})
		return
	}

	var req BookUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	// 更新字段
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Cover != "" {
		book.Cover = req.Cover
	}
	if req.Category != "" {
		book.Category = req.Category
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Status != "" {
		book.Status = req.Status
	}
	if req.WordCount != 0 {
		book.WordCount = req.WordCount
	}

	if err := bs.UpdateBook(book); err != nil {
		logger.Errorf("Failed to update book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update book",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully updated book",
		"data":    book,
	})
}

// 删除书籍
func DeleteBook(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	idStr := c.Param("book_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	if err := bs.DeleteBook(uint(id)); err != nil {
		logger.Errorf("Failed to delete book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete book",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully deleted book",
		"data":    nil,
	})
}

// 获取书籍的所有章节
func GetBookChapters(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	chapters, err := bs.GetChaptersByBookID(uint(bookID))
	if err != nil {
		logger.Errorf("Failed to get chapter list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get chapter list",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully retrieved chapter list",
		"data":    chapters,
	})
}

// 获取相关书籍（同类别，点击率高的书籍）
func GetRelatedBooks(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	// 获取相关书籍，默认返回4本
	relatedBooks, err := bs.GetRelatedBooks(uint(bookID), 4)
	if err != nil {
		logger.Errorf("Failed to get related books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get related books",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully retrieved related books",
		"data":    relatedBooks,
	})
}

// 获取书籍评论列表
func GetBookComments(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	// 获取书籍ID
	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		logger.Errorf("Failed to parse book ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	// 调用服务层获取评论列表
	comments, err := bs.GetCommentsByBookID(uint(bookID))
	if err != nil {
		logger.Errorf("Failed to get book comments [BookID: %d]: %v", bookID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get book comments",
			"data":    nil,
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully retrieved comments",
		"data":    comments,
	})
}

// 根据ID获取章节
func GetChapterByNo(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	bookIDStr := c.Param("book_id")
	chapterNoStr := c.Param("chapter_no")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	chapterNo, err := strconv.ParseUint(chapterNoStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid chapter ID",
			"data":    nil,
		})
		return
	}

	chapter, err := bs.GetChapterByNo(uint(bookID), uint(chapterNo))
	if err != nil {
		logger.Errorf("Failed to get chapter details: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Chapter not found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully retrieved chapter",
		"data":    chapter,
	})
}

// 添加章节
type ChapterCreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	IsVip   bool   `json:"is_vip"`
}

type ChapterContentResponse struct {
	ID      uint   `json:"id"`
	BookID  uint   `json:"bookId"`
	Title   string `json:"title"`
	Content string `json:"content"`
	IsVip   bool   `json:"isVip"`
}

func AddChapter(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	bookIDStr := c.Param("book_id")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	var req ChapterCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	if err := bs.AddChapter(uint(bookID), req.Title, req.Content, req.IsVip); err != nil {
		logger.Errorf("Failed to add chapter: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add chapter",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"book_id": bookID,
			"title":   req.Title,
			"content": req.Content,
			"is_vip":  req.IsVip,
		},
	})
}

// 更新章节
type ChapterUpdateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	IsVip   bool   `json:"is_vip"`
}

func UpdateChapter(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	bookIDStr := c.Param("book_id")
	chapterNoStr := c.Param("chapter_no")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	chapterNo, err := strconv.ParseUint(chapterNoStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid chapter ID",
			"data":    nil,
		})
		return
	}

	var req ChapterUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	if err := bs.UpdateChapter(uint(bookID), uint(chapterNo), req.Title, req.Content, req.IsVip); err != nil {
		logger.Errorf("Failed to update chapter: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update chapter",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data": gin.H{
			"book_id":    bookID,
			"chapter_no": chapterNo,
			"title":      req.Title,
			"content":    req.Content,
			"is_vip":     req.IsVip,
		},
	})
}

// 删除章节
func DeleteChapter(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}
	bookIDStr := c.Param("book_id")
	chapterNoStr := c.Param("chapter_no")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid book ID",
			"data":    nil,
		})
		return
	}

	chapterNo, err := strconv.ParseUint(chapterNoStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid chapter ID",
			"data":    nil,
		})
		return
	}

	if err := bs.DeleteChapter(uint(bookID), uint(chapterNo)); err != nil {
		logger.Errorf("Failed to delete chapter: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete chapter",
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

// 获取推荐书籍通用处理函数
func GetRcmds(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Book service not initialized",
			"data":    nil,
		})
		return
	}

	// 从路径参数获取推荐类型
	rcmdType := c.Param("type")
	var rcmds []model.Rcmd
	var err error

	// 根据推荐类型调用不同的服务方法
	switch rcmdType {
	case "hot":
		rcmds, err = bs.GetRcmdByType("hot", "Hot Recommendation")
	case "top":
		rcmds, err = bs.GetRcmdByType("top", "Top Recommendation")
	case "curated":
		rcmds, err = bs.GetRcmdByType("curated", "Curated Recommendation")
	case "featured":
		rcmds, err = bs.GetRcmdByType("featured", "Featured Recommendation")
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid recommendation type",
			"data":    nil,
		})
		return
	}

	if err != nil {
		logger.Errorf("Failed to get %v recommendation: %v", rcmdType, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get recommendation",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    rcmds,
	})
}

// 添加推荐
func AddRcmd(c *gin.Context) {
	// 从路径参数获取推荐类型
	rcmdType := c.Param("type")

	var req struct {
		BookID uint `json:"bookId" binding:"required"`
		Order  int  `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	// 检查书籍是否存在
	if _, err := bs.GetBookByID(req.BookID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Book not found",
			"data":    nil,
		})
		return
	}

	// 检查是否已存在
	var existingRcmd model.Rcmd
	if err := service.DB.Where("rcmd_type = ? AND book_id = ?", rcmdType, req.BookID).First(&existingRcmd).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Book already in recommendation list",
			"data":    nil,
		})
		return
	}

	// 创建新的推荐条目
	rcmd := model.Rcmd{
		RcmdType: rcmdType,
		BookID:   req.BookID,
		Order:    req.Order,
	}

	if err := service.DB.Create(&rcmd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add recommendation",
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

// 删除推荐
func DeleteRcmd(c *gin.Context) {
	// 从路径参数获取推荐类型和ID
	rcmdType := c.Param("type")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid recommendation ID",
			"data":    nil,
		})
		return
	}

	// 删除指定类型和ID的推荐条目
	if err := service.DB.Where("rcmd_type = ? AND id = ?", rcmdType, id).Delete(&model.Rcmd{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete recommendation",
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

// 批量更新推荐顺序
func UpdateRcmds(c *gin.Context) {
	// 从路径参数获取推荐类型
	rcmdType := c.Param("type")

	var rcmds []model.Rcmd
	if err := c.ShouldBindJSON(&rcmds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid parameters",
			"data":    nil,
		})
		return
	}

	// 批量更新顺序
	tx := service.DB.Begin()
	for _, rcmd := range rcmds {
		if err := tx.Model(&model.Rcmd{}).Where("id = ? AND rcmd_type = ?", rcmd.ID, rcmdType).Update("order", rcmd.Order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Failed to update recommendation order",
				"data":    nil,
			})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    rcmds,
	})
}
