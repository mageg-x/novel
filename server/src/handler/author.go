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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Service not initialized",
			"data":    nil,
		})
		return
	}
	// 获取作者ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid author ID",
			"data":    nil,
		})
		return
	}
	// 获取作者信息
	user, err := us.GetUserByID(uint(id))
	if err != nil {
		logger.Errorf("Failed to get author info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get author info",
			"data":    nil,
		})
		return
	}
	// 返回作者信息
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    user,
	})
}

// GetAuthorBooks 获取作者作品列表
func GetAuthorBooks(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Service not initialized",
			"data":    nil,
		})
		return
	}
	// 获取作者ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid author ID",
			"data":    nil,
		})
		return
	}
	// 获取作者作品列表
	books, err := us.GetBooksByAuthorID(uint(id))
	if err != nil {
		logger.Errorf("Failed to get author books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get author books",
			"data":    nil,
		})
		return
	}
	// 返回作者作品列表
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    books,
	})
}

// GetAuthorStats 获取作者统计数据
func GetAuthorStats(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Service not initialized",
			"data":    nil,
		})
		return
	}
	// 获取作者ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid author ID",
			"data":    nil,
		})
		return
	}
	// 获取作者统计数据
	stats, err := us.GetAuthorStats(uint(id))
	if err != nil {
		logger.Errorf("Failed to get author statistics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get author statistics",
			"data":    nil,
		})
		return
	}
	// 返回作者统计数据
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "succeed",
		"data":    stats,
	})
}
