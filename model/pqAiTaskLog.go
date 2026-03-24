package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PqAiTaskLog  。
type PqAiTaskLog struct {
	Id            int64          `gorm:"column:id;primaryKey;" json:"Id"`             //type:int64            comment:主键               version:2026-02-24 15:24
	CreatedAt     *time.Time     `gorm:"column:created_at;" json:"CreatedAt"`         //type:*time.Time       comment:                   version:2026-02-24 15:24
	UpdatedAt     *time.Time     `gorm:"column:updated_at;" json:"UpdatedAt"`         //type:*time.Time       comment:                   version:2026-02-24 15:24
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;" json:"DeletedAt"`         //type:gorm.DeletedAt   comment:                   version:2026-02-24 15:24
	RequestParams datatypes.JSON `gorm:"column:request_params;" json:"RequestParams"` //type:datatypes.JSON   comment:请求参数           version:2026-02-24 15:24
	ResponseData  datatypes.JSON `gorm:"column:response_data;" json:"ResponseData"`   //type:datatypes.JSON   comment:返回结果/响应体    version:2026-02-24 15:24
	AiTaskId      int64          `gorm:"column:ai_task_id;" json:"AiTaskId"`          //type:int64            comment:任务id             version:2026-02-24 15:24
}

// TableName 表名:pq_ai_task_log，。
func (*PqAiTaskLog) TableName() string {
	return "pq_ai_task_log"
}
