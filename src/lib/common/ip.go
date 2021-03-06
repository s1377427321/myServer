package common

import (
	"github.com/ying32/qqwry"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data`
}

type IP struct {
	Country string `json:"country"`
	Area    string `json:"area"`
	Region  string `json:"region"`
	City    string `json:"city"`
	Isp     string `json:"isp"`
}

var g_qqwry *qqwry.QQWry

func TabaoAPI(ip string) *IPInfo {
	resp, err := http.Get(fmt.Sprintf("http://ip.taobao.com/service/getIpInfo.php?ip=%s", ip))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}

	return &result
}

func LoadIPData(ipdatapath string) bool {
	/*
		err := ipquery.Load(ipdatapath)
		if err != nil {
			return false
		}
		return true
	*/
	g_qqwry = qqwry.NewQQWry(ipdatapath)

	return true
}


func FindAddrFromIP(ip string) string {
	/*
		result, err := ipquery.Find(ip)
		if err != nil {
			fmt.Println("FindAddrFromIP 000", err)
			return ""
		}
		return strings.Split(string(result), "\t")[0]
	*/
	return g_qqwry.GetIPLocationOfString(ip)
}