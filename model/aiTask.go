package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AiTask  AI 任务表。

type AiTask struct {
	Id               *int64          `gorm:"column:id;primaryKey" json:"Id"`                   //type:*int64            comment:                                                   version:2026-02-15 21:40
	CreatedAt        *time.Time      `gorm:"column:created_at" json:"CreatedAt"`               //type:*time.Time        comment:创建时间                                           version:2026-02-15 21:40
	UpdatedAt        *time.Time      `gorm:"column:updated_at" json:"UpdatedAt"`               //type:*time.Time        comment:更新时间                                           version:2026-02-15 21:40
	DeletedAt        *gorm.DeletedAt `gorm:"column:deleted_at" json:"DeletedAt"`               //type:*gorm.DeletedAt   comment:                                                   version:2026-02-15 21:40
	GenerateTaskId   string          `gorm:"column:generate_task_id" json:"GenerateTaskId"`    //type:string            comment:供应商任务id                                       version:2026-02-15 21:40
	UserId           *int64          `gorm:"column:user_id" json:"UserId"`                     //type:*int64            comment:用户id                                             version:2026-02-15 21:40
	ModelId          *int64          `gorm:"column:model_id" json:"ModelId"`                   //type:*int64            comment:ai模型id                                           version:2026-02-15 21:40
	Status           *int            `gorm:"column:status" json:"Status"`                      //type:*int              comment:状态                                               version:2026-02-15 21:40
	Params           datatypes.JSON  `gorm:"column:params" json:"Params"`                      //type:datatypes.JSON    comment:存储原始生成参数 (如 seed, motion_bucket_id 等)    version:2026-02-15 21:40
	ErrorMessage     string          `gorm:"column:error_message" json:"ErrorMessage"`         //type:string            comment:失败原因                                           version:2026-02-15 21:40
	CompletionTokens *int            `gorm:"column:completion_tokens" json:"CompletionTokens"` //type:*int              comment:输出视频花费的 token 数                            version:2026-02-15 21:40
	TotalTokens      *int            `gorm:"column:total_tokens" json:"TotalTokens"`           //type:*int              comment:本次请求消耗的总 token 数                          version:2026-02-15 21:40
	Key              string          `gorm:"column:key" json:"Key"`                            //type:string            comment:                                                   version:2026-02-16 20:08
	PqAiTaskLog      PqAiTaskLog     `gorm:"foreignKey:AiTaskId;references:Id" json:"PqAiTaskLog"`
}

// TableName 表名:ai_task，AI 任务表。
// 说明:
func (*AiTask) TableName() string {
	return "pq_ai_task"
}

type VideoGenerateRequest struct {
	Model         string    `json:"model"`          // 模型版本，例如: ep-20260326131643-jgvtf
	Ratio         string    `json:"ratio"`          // 视频画幅比例，例如: 9:16
	Content       []Content `json:"content"`        // 多态内容数组(文本分镜、参考图等)
	Duration      int       `json:"duration"`       // 生成时长(秒)
	Resolution    string    `json:"resolution"`     // 分辨率，例如: 480p
	GenerateAudio bool      `json:"generate_audio"` // 是否生成音频
}

// Content 多态内容项
// 💡 工程师提示：由于数组内元素结构不同，非通用字段需加上 omitempty，结构体尽量用指针
type Content struct {
	Type     string    `json:"type"`                // 类型: text 或 image_url (通用必填)
	Text     string    `json:"text,omitempty"`      // 文本内容 (仅当 type="text" 时存在)
	Role     string    `json:"role,omitempty"`      // 角色设定 (如 reference_image)
	ImageURL *ImageURL `json:"image_url,omitempty"` // 图片链接结构 (指针类型，仅当 type="image_url" 时存在)
}

// ImageURL 图片链接结构
type ImageURL struct {
	URL string `json:"url"` // 资产路径或公网图片URL
}
