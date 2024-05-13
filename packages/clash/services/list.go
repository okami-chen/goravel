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

	if in != "" {
		codes := strings.Split(in, ".")
		//query = query.Where(clause.IN{
		//	Column: "code",
		//	Values: StrToInterface(codes),
		//})
		var wheres []clause.Expression
		for _, c := range codes {
			//带m标识指定家宽或者原生
			if strings.Contains(c, "m") {
				val := strings.Replace(c, "m", "", 1)
				var conds []clause.Expression
				conds = append(conds, clause.Eq{Column: "code", Value: val})
				conds = append(conds, clause.Or(
					clause.Like{Column: "name", Value: "%家宽%"},
					clause.Like{Column: "name", Value: "%原生%"},
				))
				wheres = append(wheres, clause.Or(clause.And(conds...)))
			} else {
				wheres = append(wheres, clause.Or(clause.Eq{Column: "code", Value: c}))
			}
		}
		query = query.Where(clause.And(wheres...))
	}

	if names != nil && len(names) > 0 {
		var conds []clause.Expression
		for _, v := range names {
			lk := "%" + v.(string) + "%"
			conds = append(conds, clause.Like{Column: "name", Value: lk})
		}
		query = query.Where(clause.Or(conds...))
	}

	if out != "" {
		query = query.Where(clause.Not(clause.IN{
			Column: "code",
			Values: StrToInterface(strings.Split(out, ".")),
		}))
	}

	if ctx.Request().Input("tag") != "" {
		var more []clause.Expression
		more = append(more, clause.Like{Column: "name", Value: "%家宽%"})
		more = append(more, clause.Like{Column: "name", Value: "%原生%"})
		query = query.Where(clause.Or(more...))
	}

	query.Order("name").Find(&list)
	return list
}
