package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mageg-x/novel/src/model"
	"github.com/mageg-x/novel/src/service"
)

// 批量更新排行榜顺序
func UpdateRanks(c *gin.Context) {
	// 从路径参数获取排行榜类型
	rankType := c.Param("type")
	
	var ranks []model.Rank
	if err := c.ShouldBindJSON(&ranks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	
	// 批量更新顺序
	tx := service.DB.Begin()
	for _, rank := range ranks {
		if err := tx.Model(&model.Rank{}).Where("id = ? AND rank_type = ?", rank.ID, rankType).Update("order", rank.Order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新排行榜顺序失败"})
			return
		}
	}
	
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新排行榜顺序成功",
	})
}

// 添加排行榜
func AddRank(c *gin.Context) {
	// 从路径参数获取排行榜类型
	rankType := c.Param("type")
	
	var req struct {
		BookID uint `json:"bookId" binding:"required"`
		Order  int  `json:"order"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	
	// 检查书籍是否存在
	if _, err := bs.GetBookByID(req.BookID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "书籍不存在"})
		return
	}
	
	// 检查是否已存在
	var existingRank model.Rank
	if err := service.DB.Where("rank_type = ? AND book_id = ?", rankType, req.BookID).First(&existingRank).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该书籍已在排行榜中"})
		return
	}
	
	// 创建新的排行榜条目
	rank := model.Rank{
		RankType: rankType,
		BookID:   req.BookID,
		Order:    req.Order,
	}
	
	if err := service.DB.Create(&rank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加排行榜失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加排行榜成功",
		"data":    rank,
	})
}

// 删除排行榜
func DeleteRank(c *gin.Context) {
	// 从路径参数获取排行榜类型和ID
	rankType := c.Param("type")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的排行榜ID"})
		return
	}
	
	// 删除指定类型和ID的排行榜条目
	if err := service.DB.Where("rank_type = ? AND id = ?", rankType, id).Delete(&model.Rank{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除排行榜失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除排行榜成功",
	})
}

// 获取排行榜通用处理函数
func GetRanks(c *gin.Context) {
	// 检查BookService指针是否为nil
	if bs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "书籍服务未初始化"})
		return
	}

	// 从路径参数获取排行榜类型
	rankType := c.Param("type")
	var ranks []model.Rank
	var err error

	// 根据排行榜类型调用不同的服务方法
	switch rankType {
	case "click":
		ranks, err = bs.GetRankByType("click", "点击榜")
	case "new":
		ranks, err = bs.GetRankByType("new", "新书榜")
	case "update":
		ranks, err = bs.GetRankByType("update", "更新榜")
	case "comment":
		ranks, err = bs.GetRankByType("comment", "评论榜")
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的排行榜类型"})
		return
	}

	if err != nil {
		logger.Errorf("获取%v排行榜失败: %v", rankType, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排行榜失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    ranks,
	})
}
