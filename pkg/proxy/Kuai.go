package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

var (
	PROXY_HOST     = "d162.kdltpspro.com:15818"
	PROXY_USERNAME = "t13284784715160"
	PROXY_PASSWORD = "zye1svvx"
	mu             sync.RWMutex
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
	resp, err := http.Get(common.CONFIG.String("proxy.link"))
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
	client, err := NewHttpClientWithProxy()
	if err != nil {
		common.LOG.Error("创建代理客户端失败: " + err.Error())
		return
	}

	resp, err := client.Get("http://d162.kdltpspro.com")
	if err != nil {
		common.LOG.Error("代理连接测试失败: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		common.LOG.Info("代理连接正常")
	} else {
		common.LOG.Error(fmt.Sprintf("代理返回异常状态码: %d", resp.StatusCode))
	}
}

func NewHttpClientWithProxy() (*http.Client, error) {
	if !common.CONFIG.Bool("proxy.switchon") {
		return &http.Client{}, nil
	}

	proxyURL := fmt.Sprintf("http://%s:%s@%s", PROXY_USERNAME, PROXY_PASSWORD, PROXY_HOST)
	urli := url.URL{}
	urlproxy, err := urli.Parse(proxyURL)
	if err != nil {
		common.LOG.Error("URL Parse:" + err.Error())
		return &http.Client{}, nil
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy:              http.ProxyURL(urlproxy),
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
		},
	}

	return client, nil
}
