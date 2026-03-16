package main

import (
	"aiDev/lib/db"
	"aiDev/model"
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	db.InitMySQL("root:Tm@587973@tcp(47.100.85.27:3306)/ai_dev?charset=utf8mb4&parseTime=True&loc=Local")
	router.POST("/v1/volcengine/callback", func(c *gin.Context) {
		// 1. 获取令牌 (从 APISIX 转发的 Authorization header)
		authHeader := c.GetHeader("Authorization")
		consumerName := c.GetHeader("X-Consumer-Name")

		// 2. 读取原始 Body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("读取 Body 失败: %v", err)
			c.Status(400)
			return
		}

		// 3. 解析 JSON
		var respMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &respMap); err != nil {
			log.Printf("解析 JSON 失败: %v", err)
			c.Status(400)
			return
		}

		// 4. 提取 task_id，写入 key
		// 假设结构是 {"data": {"task_id": "xxx"}}
		if data, ok := respMap["data"].(map[string]interface{}); ok {
			if taskID, ok := data["task_id"].(string); ok {
				db.DB.Create(&model.AiTask{
					GenerateTaskId: taskID,
					Key:            authHeader, // 将令牌写入 key 字段
				})
				log.Printf("成功保存任务 ID: %s, Consumer: %s", taskID, consumerName)
			}
		}

		c.Status(200)
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
