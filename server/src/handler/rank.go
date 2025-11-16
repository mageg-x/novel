package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mageg-x/novel/src/model"
)

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
