package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mageg-x/novel/src/model"
	"github.com/mageg-x/novel/src/service"
)

// 评论管理
// 获取评论列表 - 支持搜索和分页
func GetComments(c *gin.Context) {
	bookIdStr := c.Query("bookId")
	search := c.Query("search") // 搜索关键词
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var comments []model.Comment
	var total int64

	// 构建查询
	query := service.DB.Model(&model.Comment{})

	// 根据书籍ID查询评论
	if bookIdStr != "" {
		bookId, _ := strconv.ParseUint(bookIdStr, 10, 32)
		query = query.Where("book_id = ?", bookId)
	}

	// 添加搜索条件
	if search != "" {
		searchPattern := "%" + search + "%"
		// 搜索评论内容或用户名
		query = query.Joins("LEFT JOIN users ON users.id = comments.user_id").
			Where("comments.content LIKE ? OR users.username LIKE ? OR users.nickname LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		logger.Errorf("获取评论总数失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论总数失败"})
		return
	}

	// 获取分页数据
	if err := query.Preload("User").Preload("Book").Offset(offset).Limit(limit).Order("create_time desc").Find(&comments).Error; err != nil {
		logger.Errorf("获取评论列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"comments": comments,
			"total":    total,
		},
	})
}

// 删除评论
func DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	if err := service.DB.Delete(&model.Comment{}, id).Error; err != nil {
		logger.Errorf("删除评论失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除评论成功",
	})
}