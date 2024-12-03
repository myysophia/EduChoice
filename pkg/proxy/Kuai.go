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

type ProxyScrapeResponse struct {
	ShownRecords int     `json:"shown_records"`
	TotalRecords int     `json:"total_records"`
	Limit        int     `json:"limit"`
	Skip         int     `json:"skip"`
	NextPage     bool    `json:"nextpage"`
	Proxies      []Proxy `json:"proxies"`
}

type Proxy struct {
	Alive          bool    `json:"alive"`
	AliveSince     float64 `json:"alive_since"`
	Anonymity      string  `json:"anonymity"`
	AverageTimeout float64 `json:"average_timeout"`
	FirstSeen      float64 `json:"first_seen"`
	IPData         IPData  `json:"ip_data"`
	Port           int     `json:"port"`
	Protocol       string  `json:"protocol"`
	ProxyURL       string  `json:"proxy"`
	SSL            bool    `json:"ssl"`
	IP             string  `json:"ip"`
}

type IPData struct {
	AS         string  `json:"as"`
	ASName     string  `json:"asname"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	ISP        string  `json:"isp"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Mobile     bool    `json:"mobile"`
	Org        string  `json:"org"`
	Proxy      bool    `json:"proxy"`
	RegionName string  `json:"regionName"`
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

func GetPublicProxies() ([]string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&proxy_format=protocolipport&format=json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch proxies: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var proxyResp ProxyScrapeResponse
	if err := json.Unmarshal(body, &proxyResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// 只返回活跃的HTTP/HTTPS代理
	var proxyList []string
	for _, proxy := range proxyResp.Proxies {
		if proxy.Alive && (proxy.Protocol == "http" || proxy.Protocol == "https") {
			proxyList = append(proxyList, fmt.Sprintf("%s://%s:%d", proxy.Protocol, proxy.IP, proxy.Port))
		}
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
	proxyList, err := GetPublicProxies()
	if err != nil {
		common.LOG.Error("Failed to get public proxies: " + err.Error())
		return &http.Client{}, nil
	}

	if len(proxyList) == 0 {
		common.LOG.Error("No available public proxies")
		return &http.Client{}, nil
	}

	// 随机选择一个代理
	proxyURL := proxyList[time.Now().UnixNano()%int64(len(proxyList))]
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
