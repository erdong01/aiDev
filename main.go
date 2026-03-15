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
		// 1. 读取原始 Body (最通用)
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}

		// 2. 解析成 Map 或者结构体
		var respMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &respMap); err != nil {
			log.Printf("解析 JSON 失败: %v", err)
			return
		}

		// 3. 提取 task_id (根据火山的实际返回路径提取)
		// 假设结构是 {"data": {"task_id": "xxx"}}
		if data, ok := respMap["data"].(map[string]interface{}); ok {
			if taskID, ok := data["task_id"].(string); ok {
				// 4. 使用 GORM 存入 MySQL
				db.DB.Create(&model.AiTask{
					GenerateTaskId: taskID,
					// RawResponse:    string(bodyBytes), // 甚至可以存下原始响应备查
				})
				log.Printf("成功保存任务 ID: %s", taskID)
			}
		}

		c.Status(200)
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
