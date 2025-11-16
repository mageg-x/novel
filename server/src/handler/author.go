package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAuthorInfo 获取作者信息
func GetAuthorInfo(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务未初始化"})
		return
	}
	// 获取作者ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的作者ID"})
		return
	}
	// 获取作者信息
	user, err := bs.GetUserByID(uint(id))
	if err != nil {
		logger.Errorf("获取作者信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取作者信息失败"})
		return
	}
	// 返回作者信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}

// GetAuthorBooks 获取作者作品列表
func GetAuthorBooks(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务未初始化"})
		return
	}
	// 获取作者ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的作者ID"})
		return
	}
	// 获取作者作品列表
	books, err := bs.GetBooksByAuthorID(uint(id))
	if err != nil {
		logger.Errorf("获取作者作品失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取作者作品失败"})
		return
	}
	// 返回作者作品列表
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    books,
	})
}

// GetAuthorStats 获取作者统计数据
func GetAuthorStats(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务未初始化"})
		return
	}
	// 获取作者ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的作者ID"})
		return
	}
	// 获取作者统计数据
	stats, err := bs.GetAuthorStats(uint(id))
	if err != nil {
		logger.Errorf("获取作者统计数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取作者统计数据失败"})
		return
	}
	// 返回作者统计数据
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}