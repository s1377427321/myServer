package common

import (
	"lib/utils"
	"encoding/json"
	"strings"
)

type DrawRateData struct {
	ID   int    `json:"ID"`
	Data string `json:"Data"`
}

func init() {
}

// 默认加载配置文件中的
func LoadDrawCardRate() {
	context, _ := utils.GetFileContext("config/draw_card_rate.json")

	var configs []*DrawRateData = make([]*DrawRateData, 0)

	err := json.Unmarshal([]byte(context), &configs)

	if err != nil {
		panic(err)
	} else {
		for _, v := range configs {
			SetRateLevel(v.ID, 1, strings.Replace(v.Data, "'", "\"", -1))
		}
	}
}