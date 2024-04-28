package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"goravel/app/models"
	"goravel/data"
	"goravel/pkg/proxy"
	"strings"
)

func getQuantumultX(l []models.Proxy, resp string) string {
	if len(l) < 1 {
		return resp
	}
	for _, row := range l {
		var j map[string]interface{}
		json.Unmarshal([]byte(row.Body), &j)
		p, _ := proxy.ParseProxyFromClashProxy(j)
		resp = resp + p.ToQuantumultX() + "\n"
	}
	return resp
}

func getSubscriptInfo(values string) string {
	s := facades.Orm()
	var infos []models.Info
	search := s.Query()
	if values != "" {
		values = strings.Replace(values, ".", "|", -1)
		values = strings.Replace(values, "+", "|", -1)
		codes := strings.Split(values, "|")
		c := models.Condition{}
		conditions := c.ConditionsEqOr("code", codes)
		search = search.Where(conditions)
	}
	search.Find(&infos)

	total := 0
	upload := 0
	download := 0
	now := carbon.Now().AddYears(5)
	for _, info := range infos {
		total = total + info.Total
		download = download + info.Download
		upload = upload + info.Upload
		if info.ExpireAt.Lt(now) {
			now = info.ExpireAt.Carbon
		}
	}

	fstr := "upload=%d; download=%d; total=%d ; expire=%d"

	return fmt.Sprintf(fstr, upload, download, total, now.Timestamp())
}

func getClash(clashYaml data.ClashYaml, proxies []models.Proxy) data.ClashYaml {
	for _, proxy := range proxies {
		clashYaml = processProxy(clashYaml, proxy)
	}
	return clashYaml
}

func processProxy(clashYaml data.ClashYaml, proxy models.Proxy) data.ClashYaml {
	var ret map[string]interface{}
	json.Unmarshal([]byte(proxy.Body), &ret)
	clashYaml.Proxies = append(clashYaml.Proxies, ret)
	clashYaml.ProxyGroups[0].Proxies = append(clashYaml.ProxyGroups[0].Proxies, proxy.Name)
	for i, pg1 := range clashYaml.ProxyGroups {
		if strings.Contains(pg1.Name, "狮城") && strings.Contains(proxy.Name, "新加坡") {
			clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
		}
		if strings.Contains(pg1.Name, "守候") && strings.Contains(proxy.Title, "守候") {
			clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
		}
		for _, country := range countries {
			if strings.Contains(pg1.Name, country) && strings.Contains(proxy.Name, country) {
				clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
			}
		}
	}
	return clashYaml
}

func buildCondition(field string, arr map[string]string, val string, query orm.Query) orm.Query {
	if val != "" {
		for k, v := range arr {
			val = strings.Replace(val, k, v, -1)
		}
		var conditions []string
		t := strings.Split(val, "|")
		for _, k := range t {
			conditions = append(conditions, field+" like '%"+k+"%'")
		}
		query = query.Where(strings.Join(conditions, " OR "))
	}
	return query
}

func buildNameQuery(values string, query orm.Query) orm.Query {
	values = strings.Replace(values, "+", "|", -1)
	values = strings.Replace(values, ".", "|", -1)
	var ems []models.Emoji
	c := models.Condition{}
	conditions := c.ConditionsEqOr("code", strings.Split(values, "|"))
	s := facades.Orm()

	s.Query().Where(conditions).Find(&ems)
	if len(ems) > 0 {
		var ns []string
		for _, i2 := range ems {
			ns = append(ns, i2.Country)
		}
		query = query.Where(c.ConditionsLikeOr("name", ns))
	}
	return query
}
