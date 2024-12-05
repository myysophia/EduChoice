package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/big-dust/DreamBridge/internal/pkg/common"
)

var (
	PROXY_HOST     = "as.n213.kdlfps.com:18866"
	PROXY_USERNAME = "f2693221377"
	PROXY_PASSWORD = "jx23httz"
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

type ProxyItem struct {
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
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
	if !common.CONFIG.Bool("proxy.switchon") {
		return
	}

	// 优先使用快代理
	if common.CONFIG.Bool("proxy.kuaidaili") {
		client, err := NewHttpClientWithProxy()
		if err != nil {
			common.LOG.Error("创建代理客户端失败: " + err.Error())
			return
		}

		// 使用一个可靠的测试URL，比如百度
		resp, err := client.Get("http://www.baidu.com")
		if err != nil {
			common.LOG.Error("代理连接测试失败: " + err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			// common.LOG.Info("快代理连接正常")
		} else {
			common.LOG.Error(fmt.Sprintf("快代理返回异常状态码: %d", resp.StatusCode))
		}
		return
	}

	// 使用公开代理
	proxyList, err := GetPublicProxies()
	if err != nil {
		common.LOG.Error("获取公开代理失败: " + err.Error())
		return
	}

	if len(proxyList) == 0 {
		common.LOG.Error("没有可用的公开代理")
		return
	}

	// 随机选择一个代理进行测试
	proxyURL := proxyList[time.Now().UnixNano()%int64(len(proxyList))]
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(proxyURL)
			},
		},
	}

	resp, err := client.Get("http://www.baidu.com")
	if err != nil {
		common.LOG.Error("公开代理连接测试失败: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		common.LOG.Info("公开代理连接正常")
	} else {
		common.LOG.Error(fmt.Sprintf("公开代理返回异常状态码: %d", resp.StatusCode))
	}
}

func GetPublicProxies() ([]string, error) {
	// 直接从本地文件获取代理列表
	proxyList, err := GetLocalProxies()
	if err != nil {
		return nil, fmt.Errorf("failed to get local proxies: %v", err)
	}

	if len(proxyList) == 0 {
		return nil, fmt.Errorf("no available proxies in local file")
	}

	return proxyList, nil
}

func NewHttpClientWithProxy() (*http.Client, error) {
	if !common.CONFIG.Bool("proxy.switchon") {
		return &http.Client{}, nil
	}

	// 优先使用快代理
	if common.CONFIG.Bool("proxy.kuaidaili") {
		proxyURL := fmt.Sprintf("http://%s:%s@%s", PROXY_USERNAME, PROXY_PASSWORD, PROXY_HOST)
		urli := url.URL{}
		urlproxy, err := urli.Parse(proxyURL)
		if err != nil {
			common.LOG.Error("URL Parse:" + err.Error())
			return &http.Client{}, nil
		}

		// 记录使用的代理IP
		// common.LOG.Info(fmt.Sprintf("Using KuaiDaiLi Proxy: %s", PROXY_HOST))

		return &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				Proxy:              http.ProxyURL(urlproxy),
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
			},
		}, nil
	}

	// 使用公开代理
	proxyList, err := GetLocalProxies()
	if err != nil {
		common.LOG.Error("Failed to get local proxies: " + err.Error())
		return &http.Client{}, nil
	}

	if len(proxyList) == 0 {
		common.LOG.Error("No available public proxies")
		return &http.Client{}, nil
	}

	// 随机选择一个代理
	proxyURL := proxyList[time.Now().UnixNano()%int64(len(proxyList))]

	// 记录使用的代理IP
	common.LOG.Info(fmt.Sprintf("Using Public Proxy: %s", proxyURL))

	urli := url.URL{}
	urlproxy, err := urli.Parse(proxyURL)
	if err != nil {
		common.LOG.Error("URL Parse:" + err.Error())
		return &http.Client{}, nil
	}

	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy:              http.ProxyURL(urlproxy),
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
		},
	}, nil
}

// 从本地文件读取代理列表
func GetLocalProxies() ([]string, error) {
	// 读取配置文件
	data, err := os.ReadFile("config/proxy.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read proxy.json: %v", err)
	}

	var proxyItems []ProxyItem
	if err := json.Unmarshal(data, &proxyItems); err != nil {
		return nil, fmt.Errorf("failed to parse proxy.json: %v", err)
	}

	var proxyList []string
	for _, proxy := range proxyItems {
		// 默认使用 http 协议
		proxyList = append(proxyList, fmt.Sprintf("http://%s:%d", proxy.IPAddress, proxy.Port))
	}

	return proxyList, nil
}
