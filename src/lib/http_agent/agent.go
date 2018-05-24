package http_agent

import (
	"net/http"
	"crypto/tls"
	"time"
	"lib/public_config"
	"io/ioutil"
)

//发送http请求
func AgentRequest(addr, path string) (ret []byte, err error) {

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	agent_client := &http.Client{
		Timeout:   time.Millisecond * time.Duration(public_config.GetCfgCommon().HttpTimeOut),
		Transport: tr,
	}
	resp, err := agent_client.Get(addr + path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}