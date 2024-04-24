package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type Info struct {
	orm.Model
	Upload    int
	Download  int
	Total     int
	ExpireAt  carbon.DateTime `gorm:"autoCreateTime;column:expire_at;comment:'过期时间'"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at;comment:'创建时间'"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at;comment:'更新时间'"`
}

func (p *Info) TableName() string {
	return "ue_info"
}
