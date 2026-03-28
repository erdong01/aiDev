package main

import (
	"aiDev/lib/db"
	"aiDev/model"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 从环境变量构建 DSN
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")

	if dbUser == "" {
		dbUser = "root"
	}
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "ai_dev"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)
	db.InitMySQL(dsn)
	// 1. 定义单对象的结构体 (去掉了外层的数组)
	type ApisixLog struct {
		Request struct {
			Headers    map[string]string `json:"headers"`
			Body       json.RawMessage   `json:"body"`
			BodyBase64 bool              `json:"body_base64"`
		} `json:"request"`
		Response struct {
			Status     int    `json:"status"`
			Body       string `json:"body"`
			BodyBase64 bool   `json:"body_base64"`
		} `json:"response"`
		Consumer struct {
			Username string `json:"username"`
		} `json:"consumer"`
	}

	router.POST("/v1/volcengine/callback", func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(400)
			return
		}

		// 2. 解析为单对象 (不再是数组)
		var logItem ApisixLog
		if err := json.Unmarshal(bodyBytes, &logItem); err != nil {
			log.Printf("解析 APISIX Log 失败: %v", err)
			c.Status(400) // 这里可以返回 200 避免 APISIX 一直重试，但 400 方便我们看错
			return
		}

		// 如果响应不是 200 或者没有 Body，直接忽略
		if logItem.Response.Status != 200 || logItem.Response.Body == "" {
			c.Status(200)
			return
		}

		decodedRespBody, err := decodeUpstreamResponseBody(logItem.Response.Body, logItem.Response.BodyBase64)
		if err != nil {
			log.Printf("解压火山引擎 Body 失败: %v", err)
			c.Status(200)
			return
		}

		decodedReqBody, err := decodeRequestPayload(logItem.Request.Body, logItem.Request.BodyBase64)
		if err != nil {
			log.Printf("解析请求 Body 失败: %v", err)
			c.Status(200)
			return
		}

		// 3. 解析火山引擎的响应体 (这里的 Body 是内嵌的 JSON 字符串)
		var respMap map[string]interface{}
		if err := json.Unmarshal(decodedRespBody, &respMap); err != nil {
			prefix := decodedRespBody
			if len(prefix) > 16 {
				prefix = prefix[:16]
			}
			log.Printf("解析火山引擎 Body 失败: %v, body_prefix_hex=%s", err, hex.EncodeToString(prefix))
			c.Status(200)
			return
		}

		// 4. 精准提取任务 ID (根据你的日志，现在的 key 是 "id")
		var taskID string
		if id, ok := respMap["id"].(string); ok {
			taskID = id
		} else if data, ok := respMap["data"].(map[string]interface{}); ok {
			// 保留一个向下兼容的逻辑，万一以后火山又变回 data 结构了
			if tid, ok := data["task_id"].(string); ok {
				taskID = tid
			}
		}

		// 5. 完美入库
		if taskID != "" {

			consumerName := logItem.Consumer.Username
			pureKey := extractTaskKey(logItem.Request.Headers, c)

			if pureKey == "" {
				log.Printf(
					"任务入库未获取到 key: taskID=%s, consumer=%s, has_authorization=%t, has_apikey=%t, has_x_api_key=%t",
					taskID,
					consumerName,
					logItem.Request.Headers["authorization"] != "" || c.GetHeader("Authorization") != "",
					logItem.Request.Headers["apikey"] != "" || c.GetHeader("Apikey") != "",
					logItem.Request.Headers["x-api-key"] != "" || c.GetHeader("X-API-Key") != "",
				)
			}
			var apiKeyData model.PqApiKey
			db.DB.Unscoped().Where("user_key = ?", pureKey).First(&apiKeyData)

			var req model.VideoGenerateRequest
			unmarshalErr := json.Unmarshal(decodedReqBody, &req)

			var aiModel model.PqAiModel
			if unmarshalErr == nil {
				db.DB.Unscoped().Where("name = ?", req.Model).First(&aiModel)
			}

			aiTask := model.AiTask{
				GenerateTaskId: taskID,
				Key:            apiKeyData.Key,
				ModelId:        &aiModel.Id,
			}
			err := db.DB.Create(&aiTask).Error

			if err != nil {
				log.Printf("入库失败: %v", err)
			} else {
				if aiTask.Id != nil {
					logEntry := model.PqAiTaskLog{
						RequestParams: datatypes.JSON(normalizeJSONPayload(decodedReqBody)),
						ResponseData:  datatypes.JSON(normalizeJSONPayload(decodedRespBody)),
						AiTaskId:      *aiTask.Id,
					}

					if err := db.DB.Create(&logEntry).Error; err != nil {
						log.Printf("任务日志入库失败: taskID=%s, err=%v", taskID, err)
					}
				} else {
					log.Printf("任务日志未写入: taskID=%s, 原因=AiTask 主键为空", taskID)
				}

				log.Printf("✅ 任务完美入库: ID=%s, 用户=%s", taskID, consumerName)
			}
		} else {
			log.Printf("未找到任务 ID，原始 Body: %s", string(decodedRespBody))
		}

		// 无论业务逻辑如何，只要格式对了，就给 APISIX 返回 200，告诉它“我收到了，不用再重试了”
		c.Status(200)
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}

func extractTaskKey(headers map[string]string, c *gin.Context) string {
	candidates := []string{
		headers["authorization"],
		headers["Authorization"],
		headers["apikey"],
		headers["Apikey"],
		headers["x-api-key"],
		headers["X-API-Key"],
		c.GetHeader("Authorization"),
		c.GetHeader("Apikey"),
		c.GetHeader("X-API-Key"),
	}

	for _, candidate := range candidates {
		if key := normalizeTaskKey(candidate); key != "" {
			return key
		}
	}

	return ""
}

func normalizeTaskKey(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	parts := strings.Fields(trimmed)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}

	return trimmed
}
