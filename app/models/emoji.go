package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type Emoji struct {
	orm.Model
	Emoji     string
	Code      string
	Country   string
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at;comment:'创建时间'"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at;comment:'更新时间'"`
}

func (p *Emoji) TableName() string {
	return "ue_emoji"
}
