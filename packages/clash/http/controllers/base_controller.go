package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"goravel/app/models"
	"goravel/data"
	"goravel/packages/clash/services"
	"goravel/pkg/proxy"
	"gorm.io/gorm/clause"
	"strings"
)

type BaseController struct {
	//Dependent services
	Countries []models.Emoji
	Clash     data.ClashYaml
	Ctx       http.Context
}

func (r BaseController) getSubInfo(values string) string {
	s := facades.Orm()
	var infos []models.Info
	search := s.Query()
	if values != "" {
		values = strings.Replace(values, "m", "", -1)
		search = search.Where(clause.IN{
			Column: "code",
			Values: services.StrToInterface(strings.Split(values, ".")),
		})
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
	//tb := 1024 * 1024 * 1024 * 1024 * 106

	fstr := "upload=%d; download=%d; total=%d ; expire=%d"

	return fmt.Sprintf(fstr, 0, total-download-upload, total, now.Timestamp())
}

func (r BaseController) getQuantumultX(l []models.Proxy, resp string) string {
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

func (r BaseController) getClash(clashYaml data.ClashYaml, proxies []models.Proxy) data.ClashYaml {
	facades.Orm().WithContext(r.Ctx).Query().Find(&r.Countries)
	for _, row := range proxies {
		clashYaml = r.processProxy(clashYaml, row)
	}
	return clashYaml
}

func (r BaseController) processProxy(clashYaml data.ClashYaml, proxy models.Proxy) data.ClashYaml {
	var ret map[string]interface{}
	json.Unmarshal([]byte(proxy.Body), &ret)
	clashYaml.Proxies = append(clashYaml.Proxies, ret)
	clashYaml.ProxyGroups[0].Proxies = append(clashYaml.ProxyGroups[0].Proxies, proxy.Name)
	_, eu, _ := sort()
	euMap := strings.Split(eu, ".")
	for i, group := range clashYaml.ProxyGroups {
		if strings.Contains(group.Name, "家宽") && strings.Contains(proxy.Name, "家宽") {
			if clashYaml.ProxyGroups[i].Proxies[0] == "DIRECT" {
				clashYaml.ProxyGroups[i].Proxies[0] = proxy.Name
			} else {
				clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
			}
		}
		if strings.Contains(group.Name, "原生") && strings.Contains(proxy.Name, "原生") {
			if clashYaml.ProxyGroups[i].Proxies[0] == "DIRECT" {
				clashYaml.ProxyGroups[i].Proxies[0] = proxy.Name
			} else {
				clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
			}
		}
		if strings.Contains(group.Name, "狮城") && strings.Contains(proxy.Name, "🇸🇬") {
			clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
		}
		if strings.Contains(group.Name, "守候") && proxy.Code == "h" {
			clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
		}
		if strings.Contains(group.Name, "欧盟") {
			for _, emoji := range euMap {
				if strings.Contains(proxy.Name, emoji) {
					if clashYaml.ProxyGroups[i].Proxies[0] == "DIRECT" {
						clashYaml.ProxyGroups[i].Proxies[0] = proxy.Name
					} else {
						clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
					}
					break
				}
			}
		}

		if strings.Contains(group.Name, "影院") {
			if strings.Contains(proxy.Name, "S0") || strings.Contains(proxy.Name, "S1") {
				if clashYaml.ProxyGroups[i].Proxies[0] == "DIRECT" {
					clashYaml.ProxyGroups[i].Proxies[0] = proxy.Name
				} else {
					clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
				}
			}
		}

		for _, country := range r.Countries {
			if country.Code == "SG" {
				continue
			}
			if strings.Contains(group.Name, country.Country) && strings.Contains(proxy.Name, country.Emoji) {
				if clashYaml.ProxyGroups[i].Proxies[0] == "DIRECT" {
					clashYaml.ProxyGroups[i].Proxies[0] = proxy.Name
				} else {
					clashYaml.ProxyGroups[i].Proxies = append(clashYaml.ProxyGroups[i].Proxies, proxy.Name)
				}
			}
		}
	}
	return clashYaml
}
