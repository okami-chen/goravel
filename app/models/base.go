package models

import (
	"github.com/goravel/framework/support/json"
	"reflect"
	"strings"
)

type Base struct {
}

func (p *Base) ToArray() map[string]interface{} {
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
					continue
				}
			}
		}
		result[field.Name] = value
	}
	return result
}
