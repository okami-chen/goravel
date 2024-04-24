package tool

func GetNameReplaces() map[string]string {
	locs := map[string]string{
		"+": "|",
		".": "|",
		//节点名称
		"mobile": "家宽",
		"HK":     "香港",
		"SG":     "新加坡",
		"US":     "美国",
		"JP":     "日本",
		"TW":     "台湾",
		"KR":     "韩国",
		"UK":     "英国",
		"GB":     "英国",
		"DE":     "德国",
		"FR":     "法国",
		"IT":     "意大利",
		"ES":     "西班牙",
		"CA":     "加拿大",
		"AU":     "澳大利亚",
		"RU":     "俄罗斯",
		"BR":     "巴西",
		"IN":     "印度",
		"AR":     "阿根廷",
		"MX":     "墨西哥",
		"NZ":     "新西兰",
		"ZA":     "南非",
		"CH":     "瑞士",
		"NL":     "荷兰",
		"SE":     "瑞典",
		"NO":     "挪威",
		"DK":     "丹麦",
		"PL":     "波兰",
		"BE":     "比利时",
		"FI":     "芬兰",
		"GR":     "希腊",
		"PT":     "葡萄牙",
		"TR":     "土耳其",
		"CN":     "中国",
		"VN":     "越南",
		"MY":     "马来西亚",
	}
	return locs
}

func GetGroupReplaces() map[string]string {
	locs := map[string]string{
		"+": "|",
		".": "|",
		//分组
		"s1": "Flower",
		"s2": "SSlinks",
		"s3": "Mesl",
		"s4": "守候",
		"s5": "Imm",
		"s6": "Tag",
		"s7": "CreamData",
		"s8": "奶优",
		"s9": "Cattle",
	}
	return locs
}
