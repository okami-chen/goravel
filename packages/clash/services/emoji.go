package services

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
)

func StrToInterface(str []string) []interface{} {
	values := []interface{}{}
	for _, v := range str {
		values = append(values, v)
	}
	return values
}

// FindEmojiByCode 通过代码获取国旗
func FindEmojiByCode(code string, ctx http.Context) []interface{} {

	if code == "" {
		return nil
	}

	//转换成数组
	code = strings.ToUpper(code)
	codes := strings.Split(code, ".")
	values := make([]interface{}, len(codes))
	for i, v := range codes {
		values[i] = v
	}

	// 如果是emoji直接返回
	match, _ := regexp.MatchString("^[a-zA-Z]+$", values[0].(string))
	if !match {
		return values
	}
	var emoji []models.Emoji
	query := facades.Orm().WithContext(ctx).Query()
	query = query.Where(clause.IN{
		Column: "code",
		Values: values,
	})
	query.Find(&emoji)
	// 按照入参排序
	sort := make(map[string]string, 0)
	em := []interface{}{}
	for _, v := range emoji {
		sort[v.Code] = v.Emoji
	}

	// 排序
	for _, val := range values {
		if exist, ok := sort[val.(string)]; ok {
			em = append(em, exist)
		}
	}

	return em
}

func SortByEmoji(emojis []interface{}, proxies []models.Proxy) []models.Proxy {
	if emojis == nil || len(emojis) == 0 {
		return proxies
	}
	result := []models.Proxy{}
	mGroup := make(map[string][]models.Proxy)
	mOther := make([]models.Proxy, 0)
	//放进分组
	for _, v2 := range proxies {
		var found bool
		for _, v1 := range emojis {
			val := v1.(string)
			if strings.Contains(v2.Name, v1.(string)) {
				mGroup[val] = append(mGroup[val], v2)
				found = true
				break
			}
		}
		if !found {
			mOther = append(mOther, v2)
		}
	}

	// 分组
	for _, v3 := range emojis {
		val := v3.(string)
		//不存在下标
		if _, ok := mGroup[val]; !ok {
			continue
		}
		for _, exist := range mGroup[val] {
			result = append(result, exist)
		}
	}
	// 其他
	if len(mOther) > 0 {
		for _, vm := range mOther {
			result = append(result, vm)
		}
	}
	return result
}
