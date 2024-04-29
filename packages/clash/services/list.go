package services

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"gorm.io/gorm/clause"
	"strings"
)

func List(names []interface{}, in, out string, ctx http.Context) []models.Proxy {

	query := facades.Orm().WithContext(ctx).Query()
	var list []models.Proxy
	if names != nil && len(names) > 0 {
		var conds []clause.Expression
		for _, v := range names {
			lk := "%" + v.(string) + "%"
			conds = append(conds, clause.Like{Column: "name", Value: lk})
		}
		query = query.Where(clause.Or(conds...))
	}

	if in != "" {
		query = query.Where(clause.IN{
			Column: "code",
			Values: StrToInterface(strings.Split(in, ".")),
		})
	}

	if out != "" {
		query = query.Where(clause.Not(clause.IN{
			Column: "code",
			Values: StrToInterface(strings.Split(out, ".")),
		}))
	}
	query.Order("name").Find(&list)
	return list
}
