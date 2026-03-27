package model

import (
	"time"

	"gorm.io/gorm"
)

// PqAiModel  ai模型。
// 说明:
// 表名:pq_ai_model
// group: PqAiModel
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.PqAiModel
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.PqAiModel
// version:2026-02-27 21:51
type PqAiModel struct {
	Id            int64          `gorm:"column:id;primaryKey;" json:"Id"`             //type:int64            comment:            version:2026-02-27 21:51
	CreatedAt     *time.Time     `gorm:"column:created_at;" json:"CreatedAt"`         //type:*time.Time       comment:创建时间    version:2026-02-27 21:51
	UpdatedAt     *time.Time     `gorm:"column:updated_at;" json:"UpdatedAt"`         //type:*time.Time       comment:更新时间    version:2026-02-27 21:51
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;" json:"DeletedAt"`         //type:gorm.DeletedAt   comment:            version:2026-02-27 21:51
	Name          string         `gorm:"column:name;" json:"Name"`                    //type:string           comment:名称        version:2026-02-27 21:51
	Provider      string         `gorm:"column:provider;" json:"Provider"`            //type:string           comment:供应商      version:2026-02-27 21:51
	Version       string         `gorm:"column:version;" json:"Version"`              //type:string           comment:版本        version:2026-02-27 21:51
	WarningAmount float64        `gorm:"column:warning_amount;" json:"WarningAmount"` //type:float64          comment:预警金额    version:2026-02-27 21:51
	StopAmount    float64        `gorm:"column:stop_amount;" json:"StopAmount"`       //type:float64          comment:停止金额    version:2026-02-27 21:51
}

// TableName 表名:pq_ai_model，ai模型。
// 说明:
func (*PqAiModel) TableName() string {
	return "pq_ai_model"
}
