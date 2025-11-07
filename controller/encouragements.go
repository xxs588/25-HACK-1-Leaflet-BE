package controller

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/config"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/model"
	"github.com/gin-gonic/gin"
)

// 获取鼓励话语 根据时间段返回随机内容
func GetEncouragements(c *gin.Context) {
	timeNow := time.Now()
	hour := timeNow.Hour()

	// 早上时段 (0-12点)
	if hour < 12 {
		var encouragements []model.EncouragementMorning
		if err := config.DB.Find(&encouragements).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// 如果没有数据，返回空
		if len(encouragements) == 0 {
			c.JSON(http.StatusOK, gin.H{"encouragement": nil, "message": "暂无鼓励话语"})
			return
		}
		
		// 随机选择一条
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(encouragements))))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "随机数生成失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"encouragement": encouragements[randomIndex.Int64()]})
		return
	}

	// 下午时段 (12-18点)
	if hour < 18 {
		var encouragementsAfternoon []model.EncouragementAfternoon
		if err := config.DB.Find(&encouragementsAfternoon).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// 如果没有数据，返回空
		if len(encouragementsAfternoon) == 0 {
			c.JSON(http.StatusOK, gin.H{"encouragement": nil, "message": "暂无鼓励话语"})
			return
		}
		
		// 随机选择一条
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(encouragementsAfternoon))))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "随机数生成失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"encouragement": encouragementsAfternoon[randomIndex.Int64()]})
		return
	}

	// 晚上时段 (18-24点)
	var encouragementsEvening []model.EncouragementEvening
	if err := config.DB.Find(&encouragementsEvening).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 如果没有数据，返回空
	if len(encouragementsEvening) == 0 {
		c.JSON(http.StatusOK, gin.H{"encouragement": nil, "message": "暂无鼓励话语"})
		return
	}
	
	// 随机选择一条
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(encouragementsEvening))))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "随机数生成失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"encouragement": encouragementsEvening[randomIndex.Int64()]})
}