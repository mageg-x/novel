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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}

	var params BookQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	books, total, err := bs.GetAllBooks(params.Offset, params.Limit)
	if err != nil {
		logger.Errorf("获取书籍列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取书籍列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	idStr := c.Param("book_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	book, err := bs.GetBookByID(uint(id))
	if err != nil {
		logger.Errorf("获取书籍详情失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "书籍不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    book,
	})
}

// 根据分类获取书籍
func GetBooksByCategory(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	category := c.Param("category")

	var params BookQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	books, total, err := bs.GetBooksByCategory(category, params.Offset, params.Limit)
	if err != nil {
		logger.Errorf("根据分类获取书籍失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取书籍列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	var req BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
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
		logger.Errorf("添加书籍失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加书籍失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加书籍成功",
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	idStr := c.Param("book_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	// 获取现有书籍
	book, err := bs.GetBookByID(uint(id))
	if err != nil {
		logger.Errorf("获取书籍详情失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "书籍不存在"})
		return
	}

	var req BookUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
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
		logger.Errorf("更新书籍失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新书籍失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新书籍成功",
		"data":    book,
	})
}

// 删除书籍
func DeleteBook(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	idStr := c.Param("book_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	if err := bs.DeleteBook(uint(id)); err != nil {
		logger.Errorf("删除书籍失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除书籍失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除书籍成功",
	})
}

// 获取书籍的所有章节
func GetBookChapters(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	chapters, err := bs.GetChaptersByBookID(uint(bookID))
	if err != nil {
		logger.Errorf("获取章节列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取章节列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    chapters,
	})
}

// 获取相关书籍（同类别，点击率高的书籍）
func GetRelatedBooks(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	// 获取相关书籍，默认返回4本
	relatedBooks, err := bs.GetRelatedBooks(uint(bookID), 4)
	if err != nil {
		logger.Errorf("获取相关书籍失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取相关书籍失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    relatedBooks,
	})
}

// 获取书籍评论列表
func GetBookComments(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	// 获取书籍ID
	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		logger.Errorf("书籍ID解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的书籍ID",
		})
		return
	}

	// 调用服务层获取评论列表
	comments, err := bs.GetCommentsByBookID(uint(bookID))
	if err != nil {
		logger.Errorf("获取书籍评论失败[书籍ID: %d]: %v", bookID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取书籍评论失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    comments,
	})
}

// 根据ID获取章节
func GetChapterByID(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	bookIDStr := c.Param("book_id")
	chapterIDStr := c.Param("chapter_id")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节ID"})
		return
	}

	chapter, err := bs.GetChapterByID(uint(bookID), uint(chapterID))
	if err != nil {
		logger.Errorf("获取章节详情失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
		return
	}

	// 验证章节是否属于指定书籍
	if chapter.BookID != uint(bookID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "章节不属于指定书籍"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	bookIDStr := c.Param("book_id")
	chapterIDStr := c.Param("chapter_id")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节ID"})
		return
	}

	var req ChapterCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	chapter := &model.Chapter{
		ChapterID: uint(chapterID),
		BookID:    uint(bookID),
		Title:     req.Title,
		IsVip:     req.IsVip,
	}

	if err := bs.AddChapter(chapter, req.Content); err != nil {
		logger.Errorf("添加章节失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加章节失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加章节成功",
		"data":    chapter,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	bookIDStr := c.Param("book_id")
	chapterIDStr := c.Param("chapter_id")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节ID"})
		return
	}

	// 获取现有章节
	chapter, err := bs.GetChapterByID(uint(bookID), uint(chapterID))
	if err != nil {
		logger.Errorf("获取章节详情失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
		return
	}

	// 验证章节是否属于指定书籍
	if chapter.BookID != uint(bookID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "章节不属于指定书籍"})
		return
	}

	var req ChapterUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 更新字段
	if req.Title != "" {
		chapter.Title = req.Title
	}

	chapter.IsVip = req.IsVip

	if err := bs.UpdateChapter(chapter, req.Content); err != nil {
		logger.Errorf("更新章节失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新章节失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新章节成功",
		"data":    chapter,
	})
}

// 删除章节
func DeleteChapter(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}
	bookIDStr := c.Param("book_id")
	chapterIDStr := c.Param("chapter_id")

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的书籍ID"})
		return
	}

	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节ID"})
		return
	}

	// 获取章节信息，验证是否属于指定书籍
	chapter, err := bs.GetChapterByID(uint(bookID), uint(chapterID))
	if err != nil {
		logger.Errorf("获取章节详情失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
		return
	}

	// 验证章节是否属于指定书籍
	if chapter.BookID != uint(bookID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "章节不属于指定书籍"})
		return
	}

	if err := bs.DeleteChapter(uint(bookID), uint(chapterID)); err != nil {
		logger.Errorf("删除章节失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除章节失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除章节成功",
	})
}

// 获取推荐书籍通用处理函数
func GetRcmdBooks(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}

	// 从路径参数获取推荐类型
	rcmdType := c.Param("type")
	var rcmds []model.Rcmd
	var err error

	// 根据推荐类型调用不同的服务方法
	switch rcmdType {
	case "hot":
		rcmds, err = bs.GetRcmdByType("hot", "热门推荐")
	case "top":
		rcmds, err = bs.GetRcmdByType("top", "置顶推荐")
	case "curated":
		rcmds, err = bs.GetRcmdByType("curated", "精选推荐")
	case "featured":
		rcmds, err = bs.GetRcmdByType("featured", "特色推荐")
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的推荐类型"})
		return
	}

	if err != nil {
		logger.Errorf("获取%v推荐失败: %v", rcmdType, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取推荐失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rcmds,
	})
}
