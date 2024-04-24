package models

import (
	"fmt"
	event "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/json"
	"reflect"
	"strings"
)

type Condition struct {
}

func (p *Condition) ConditionsLikeOr(key string, conditions []string) string {
	str := []string{}
	for _, k := range conditions {
		str = append(str, key+" like '%"+k+"%'")
	}
	return strings.Join(str, " OR ")
}

func (p *Condition) ConditionsEqOr(key string, conditions []string) string {
	str := []string{}
	for _, k := range conditions {
		str = append(str, key+" = '"+k+"'")
	}
	return strings.Join(str, " OR ")
}

func (p *Condition) ConditionsLikeAnd(key string, conditions []string) string {
	str := []string{}
	for _, k := range conditions {
		str = append(str, key+" like '%"+k+"%'")
	}
	return strings.Join(str, " AND ")
}

type Proxy struct {
	Id        int             `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL;comment:'编号'"`
	Status    int             `gorm:"column:status;default:1;NOT NULL;comment:'状态'"`
	Title     string          `gorm:"column:title;NOT NULL;comment:'分组'"`
	Tag       *string         `gorm:"column:tag;default:NULL;comment:'标签'"`
	Ip        *string         `gorm:"column:ip;default:NULL;comment:'ip'"`
	Isp       *string         `gorm:"column:isp;default:NULL;comment:'isp'"`
	Name      string          `gorm:"column:name;NOT NULL;comment:'名称'"`
	Type      string          `gorm:"column:type;NOT NULL;comment:'类型'"`
	Server    string          `gorm:"column:server;NOT NULL;comment:'主机'"`
	Port      int             `gorm:"column:port;NOT NULL;comment:'端口'"`
	Body      string          `gorm:"column:body;NOT NULL;comment:'正文'"`
	CheckAt   carbon.DateTime `gorm:"column:check_at;NOT NULL;comment:'检测时间'"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at;comment:'创建时间'"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at;comment:'更新时间'"`
	orm.SoftDeletes
}

func (p *Proxy) TableName() string {
	return "ue_proxy"
}

func (p *Proxy) ToArray() map[string]interface{} {
	result := make(map[string]interface{})
	t := reflect.TypeOf(*p)
	v := reflect.ValueOf(*p)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if field.Name == "SoftDeletes" {
			continue
		}
		if reflect.TypeOf(value).Name() == "string" {
			strValue := value.(string)
			if strings.Contains(strValue, "{") {
				s := make(map[string]interface{})
				err := json.UnmarshalString(strValue, &s)
				if err == nil {
					result[field.Name] = s
					fmt.Println(field.Name + " : success")
					continue
				} else {
					fmt.Println(field.Name + " : " + err.Error())
				}
			}
		}
		result[field.Name] = value
	}
	return result
}

func (u *Proxy) DispatchesEvents() map[event.EventType]func(event.Event) error {
	return map[event.EventType]func(event.Event) error{
		event.EventCreating: func(event event.Event) error {
			return nil
		},
		event.EventCreated: func(event event.Event) error {
			return nil
		},
		event.EventSaving: func(event event.Event) error {
			return nil
		},
		event.EventSaved: func(event event.Event) error {
			return nil
		},
		event.EventUpdating: func(event event.Event) error {
			return nil
		},
		event.EventUpdated: func(event event.Event) error {
			return nil
		},
		event.EventDeleting: func(event event.Event) error {
			return nil
		},
		event.EventDeleted: func(event event.Event) error {
			return nil
		},
		event.EventForceDeleting: func(event event.Event) error {
			return nil
		},
		event.EventForceDeleted: func(event event.Event) error {
			return nil
		},
		event.EventRetrieved: func(event event.Event) error {
			return nil
		},
	}
}
