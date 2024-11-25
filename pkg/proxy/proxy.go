package proxy

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/big-dust/DreamBridge/internal/pkg/common"
//	"io"
//	"net/http"
//	"net/url"
//	"sync"
//	"time"
//)
//
//var (
//	IP   string
//	PORT int
//	mu   sync.RWMutex
//)
//
//type GenProxyIPResponse struct {
//	Code    int    `json:"code"`
//	Data    []Data `json:"data"`
//	Msg     string `json:"msg"`
//	Success bool   `json:"success"`
//}
//type Data struct {
//	IP   string `json:"ip"`
//	Port int    `json:"port"`
//}
//
//func genProxyIP() (*GenProxyIPResponse, error) {
//	defer func() {
//		if r := recover(); r != nil {
//			common.LOG.Error(fmt.Sprintf("%v", r))
//		}
//	}()
//	client := &http.Client{
//		Timeout: 5 * time.Second,
//	}
//	resp, err := client.Get(common.CONFIG.String("proxy.pinyiLink"))
//	if err != nil {
//		common.LOG.Error("genProxyIP: " + err.Error())
//		return nil, err
//	}
//	bodyText, err := io.ReadAll(resp.Body)
//	if err != nil {
//		common.LOG.Error("genProxyIP: " + err.Error())
//		return nil, err
//	}
//	IPResp := &GenProxyIPResponse{}
//	if err = json.Unmarshal(bodyText, IPResp); err != nil {
//		common.LOG.Info("json.Unmarshal(bodyText, IPResp): " + string(bodyText))
//		return nil, err
//	}
//	return IPResp, nil
//}
//
//func ChangeHttpProxyIP() {
//	defer func() {
//		if r := recover(); r != nil {
//			common.LOG.Error(fmt.Sprintf("%v", r))
//		}
//	}()
//	resp, _ := genProxyIP()
//	for i := 0; i < 3; i++ {
//		if resp != nil {
//			if resp.Code == 0 {
//				break
//			}
//			common.LOG.Info(fmt.Sprintf("change ip: msg: %v", resp.Msg))
//		} else {
//			common.LOG.Info(fmt.Sprintf("change ip: resp: %v", resp))
//		}
//		time.Sleep(2 * time.Second)
//		resp, _ = genProxyIP()
//	}
//	// 尝试获取三次失败，则保持不变
//	if resp == nil || len(resp.Data) == 0 {
//		return
//	}
//
//	mu.Lock()
//	IP = resp.Data[0].IP
//	PORT = resp.Data[0].Port
//	mu.Unlock()
//}
//
//func NewHttpClientWithProxy() (*http.Client, error) {
//	if !common.CONFIG.Bool("proxy.switchon") {
//		return &http.Client{}, nil
//	}
//	mu.RLock()
//	ip := IP
//	port := PORT
//	mu.RUnlock()
//	if ip == "" {
//		return &http.Client{}, nil
//	}
//	common.LOG.Info(fmt.Sprintf("proxy url: https://%s:%d", ip, port))
//	urli := url.URL{}
//	urlproxy, _ := urli.Parse(fmt.Sprintf("http://%s:%d", ip, port))
//
//	client := &http.Client{
//		Timeout: 5 * time.Second,
//		Transport: &http.Transport{
//			Proxy: http.ProxyURL(urlproxy),
//			//TLSClientConfig:     &tls.Config{},
//			//TLSHandshakeTimeout: 0,
//		},
//	}
//	return client, nil
//}
