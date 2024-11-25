package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	IP_PORT string
	mu      sync.RWMutex
)

type GenProxyIPResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data Data   `json:"data"`
}
type Data struct {
	Count     int      `json:"count"`
	ProxyList []string `json:"proxy_list"`
}

func genProxyIP() (*GenProxyIPResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			common.LOG.Error(fmt.Sprintf("%v", r))
		}
	}()
	resp, err := http.Get(common.CONFIG.String("proxy.hugeLink"))
	if err != nil {
		return nil, err
	}
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	IPResp := &GenProxyIPResponse{}
	if err = json.Unmarshal(bodyText, IPResp); err != nil {
		common.LOG.Error("GenIP BODY:" + string(bodyText))
		return nil, err
	}
	return IPResp, nil
}

func ChangeHttpProxyIP() {
	defer func() {
		if r := recover(); r != nil {
			common.LOG.Error(fmt.Sprintf("%v", r))
		}
	}()
	resp, err := genProxyIP()
	for i := 0; i < 3; i++ {
		if resp != nil && err == nil {
			if resp.Code == 200 {
				break
			}
			common.LOG.Info(fmt.Sprintf("change ip: msg: %v", resp.Msg))
		} else {
			common.LOG.Info(fmt.Sprintf("change ip: resp: %v err: %s", resp, err.Error()))
		}
		time.Sleep(3 * time.Second)
		resp, err = genProxyIP()
	}
	// 尝试获取三次失败，则保持不变
	if resp == nil || len(resp.Data.ProxyList) == 0 {
		return
	}

	mu.Lock()
	IP_PORT = resp.Data.ProxyList[0]
	mu.Unlock()
}

func NewHttpClientWithProxy() (*http.Client, error) {
	if !common.CONFIG.Bool("proxy.switchon") {
		return &http.Client{}, nil
	}
	mu.RLock()
	ip_port := IP_PORT
	mu.RUnlock()
	if ip_port == "" {
		return &http.Client{}, nil
	}
	//common.LOG.Info(fmt.Sprintf("proxy url: http://%s", ip_port))
	urli := url.URL{}
	urlproxy, err := urli.Parse(fmt.Sprintf("http://%s", ip_port))
	if err != nil {
		common.LOG.Error("URL Parse:" + err.Error())
		return &http.Client{}, nil
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
			//TLSClientConfig:     &tls.Config{},
			//TLSHandshakeTimeout: 0,
		},
	}
	return client, nil
}
