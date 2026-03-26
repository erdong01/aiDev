package model

import (
	"time"

	"gorm.io/gorm"
)

// PqApiKey  api密钥。
// 说明:
// 表名:pq_api_key
// group: PqApiKey
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.PqApiKey
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.PqApiKey
// version:2026-02-26 14:22
type PqApiKey struct {
	Id          int64          `gorm:"column:id;primaryKey;" json:"Id"`         //type:int64            comment:                      version:2026-02-26 14:22
	UpdatedAt   *time.Time     `gorm:"column:updated_at;" json:"UpdatedAt"`     //type:*time.Time       comment:                      version:2026-02-26 14:22
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;" json:"DeletedAt"`     //type:gorm.DeletedAt   comment:                      version:2026-02-26 14:22
	CreatedAt   *time.Time     `gorm:"column:created_at;" json:"CreatedAt"`     //type:*time.Time       comment:                      version:2026-02-26 14:22
	UserId      int64          `gorm:"column:user_id;" json:"UserId"`           //type:int64            comment:用户id                version:2026-02-26 14:22
	AiModelId   int64          `gorm:"column:ai_model_id;" json:"AiModelId"`    //type:int64            comment:ai模型                version:2026-02-26 14:22
	Key         string         `gorm:"column:key;" json:"Key"`                  //type:string           comment:密钥                  version:2026-02-26 14:22
	TotalTokens int64          `gorm:"column:total_tokens;" json:"TotalTokens"` //type:int64            comment:拥有tokens数          version:2026-02-26 14:22
	UseTokens   string         `gorm:"column:use_tokens;" json:"UseTokens"`     //type:string           comment:已消耗tokens          version:2026-02-26 14:22
	UserKey     string         `gorm:"column:user_key;" json:"UserKey"`         //type:string           comment:外部用户key           version:2026-02-26 14:22
	UserName    string         `gorm:"column:user_name;" json:"UserName"`       //type:string           comment:用户名                version:2026-02-26 14:22
	Status      int            `gorm:"column:status;" json:"Status"`            //type:int              comment:状态 1 启用 2 禁用    version:2026-02-26 14:22
	Rate        int            `gorm:"column:rate;" json:"Rate"`                //type:int              comment:速率                  version:2026-02-26 14:22
}

// TableName 表名:pq_api_key，api密钥。
// 说明:
func (*PqApiKey) TableName() string {
	return "pq_api_key"
}
