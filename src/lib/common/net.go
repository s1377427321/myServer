package common

import (
	"net/http"
	"time"
	"io/ioutil"
	"strings"
	"net"
)

// 获取外网IP地址
func GetNetAddr() string {
	agent_client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := agent_client.Get("http://icanhazip.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(ret))
}

// 获取本地IP地址
func GetLocalAddr() (IpAddr string) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		IpAddr = "127.0.0.1"
		return
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				IpAddr = ipnet.IP.String()
				return
			}
		}
	}
	IpAddr = "127.0.0.1"
	return
}

