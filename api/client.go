package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// API端点列表，按优先级排序
var APIEndpoints = []string{
	"https://api.hayfrp.1zyq1.com",
	"https://v2.api.hayfrp.1zyq1.com",
	"https://api.hayfrp.com",
}

// BaseURL 当前使用的API端点
var BaseURL = APIEndpoints[0]
var currentEndpointIndex = 0
var endpointMutex sync.Mutex

// HTTPClient 公共HTTP客户端
var HTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

// SwitchToNextEndpoint 切换到下一个API端点
func SwitchToNextEndpoint() bool {
	endpointMutex.Lock()
	defer endpointMutex.Unlock()

	if currentEndpointIndex < len(APIEndpoints)-1 {
		currentEndpointIndex++
		BaseURL = APIEndpoints[currentEndpointIndex]
		fmt.Printf("[API] 切换到备用端点: %s\n", BaseURL)
		return true
	}
	return false
}

// ResetToPrimaryEndpoint 重置到第一个API端点
func ResetToPrimaryEndpoint() {
	endpointMutex.Lock()
	defer endpointMutex.Unlock()

	currentEndpointIndex = 0
	BaseURL = APIEndpoints[0]
}

// GetCurrentEndpoint 获取当前端点
func GetCurrentEndpoint() string {
	endpointMutex.Lock()
	defer endpointMutex.Unlock()
	return BaseURL
}

// DoRequestWithFallback 带故障转移的请求
func DoRequestWithFallback(httpReq *http.Request) (*http.Response, error) {
	var lastErr error

	for i := 0; i < len(APIEndpoints); i++ {
		// 获取当前尝试的端点
		endpointMutex.Lock()
		tryIndex := (currentEndpointIndex + i) % len(APIEndpoints)
		tryURL := APIEndpoints[tryIndex]
		endpointMutex.Unlock()

		// 更新请求URL
		httpReq.URL.Scheme = "https"
		httpReq.URL.Host = tryURL[8:] // 去掉 "https://"

		resp, err := HTTPClient.Do(httpReq)
		if err != nil {
			lastErr = err
			fmt.Printf("[API] 端点 %s 请求失败: %v\n", tryURL, err)
			continue
		}

		// 检查是否为服务器错误
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("服务器错误: %d", resp.StatusCode)
			fmt.Printf("[API] 端点 %s 返回错误: %d\n", tryURL, resp.StatusCode)
			continue
		}

		// 请求成功，更新当前端点
		if i > 0 {
			endpointMutex.Lock()
			currentEndpointIndex = tryIndex
			BaseURL = tryURL
			endpointMutex.Unlock()

			// 一段时间后尝试切回主端点
			go func() {
				time.Sleep(5 * time.Minute)
				ResetToPrimaryEndpoint()
			}()
		}

		return resp, nil
	}

	return nil, fmt.Errorf("所有API端点均不可用: %v", lastErr)
}
