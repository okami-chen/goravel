package services

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"gorm.io/gorm/clause"
	"strings"
)

func List(values []interface{}, c string, ctx http.Context) []models.Proxy {
	query := facades.Orm().WithContext(ctx).Query()
	var list []models.Proxy
	if len(values) > 0 {
		var conds []clause.Expression
		for _, v := range values {
			lk := "%" + v.(string) + "%"
			conds = append(conds, clause.Like{Column: "name", Value: lk})
		}
		query = query.Where(clause.Or(conds...))
	}
	if c != "" {
		query = query.Where(clause.IN{
			Column: "code",
			Values: StrToInterface(strings.Split(c, ".")),
		})
	}
	query.Order("name").Find(&list)
	return list
}
