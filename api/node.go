package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NodeAPIClient 节点相关API客户端
type NodeAPIClient struct {
	client *http.Client
}

// NewNodeAPIClient 创建节点API客户端
func NewNodeAPIClient() *NodeAPIClient {
	return &NodeAPIClient{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// NodeInfo 节点探针信息
type NodeInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	BindPort        string `json:"bind_port"`
	BindUDPPort     string `json:"bind_udp_port"`
	VhostHTTPPort   string `json:"vhost_http_port"`
	VhostHTTPSPort  string `json:"vhost_https_port"`
	KCPBindPort     string `json:"kcp_bind_port"`
	SubdomainHost   string `json:"subdomain_host"`
	MaxPoolCount    string `json:"max_pool_count"`
	MaxPortsPerClient string `json:"max_ports_per_client"`
	HeartBeatTimeout string `json:"heart_beat_timeout"`
	TotalTrafficIn  string `json:"total_traffic_in"`
	TotalTrafficOut string `json:"total_traffic_out"`
	CurConns        string `json:"cur_conns"`
	ClientCounts    string `json:"client_counts"`
	CPUUsage        string `json:"cpu_usage"`
	RAMUsage        string `json:"ram_usage"`
	DiskUsage       string `json:"disk_usage"`
	Status          string `json:"status"`
}

// GetNodeInfoResponse 获取节点探针信息响应
type GetNodeInfoResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Number  int        `json:"number"`
	Servers []NodeInfo `json:"servers"`
}

// GetNodeInfo 获取节点探针信息
func (c *NodeAPIClient) GetNodeInfo() (*GetNodeInfoResponse, error) {
	httpReq, err := http.NewRequest("GET", BaseURL+"/node", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result GetNodeInfoResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// NodeListItem 节点列表项
type NodeListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetNodeListResponse 获取节点列表响应
type GetNodeListResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Number  int            `json:"number"`
	Servers []NodeListItem `json:"servers"`
}

// GetNodeList 获取节点列表
func (c *NodeAPIClient) GetNodeList() (*GetNodeListResponse, error) {
	httpReq, err := http.NewRequest("GET", BaseURL+"/nodes", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result GetNodeListResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetNotice 获取公告
func (c *NodeAPIClient) GetNotice() (string, error) {
	httpReq, err := http.NewRequest("GET", BaseURL+"/notice", nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查是否返回JSON错误
	if len(respBody) > 0 && respBody[0] == '{' {
		var result struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return "", fmt.Errorf("解析响应失败: %w", err)
		}
		if result.Status != 200 {
			return "", fmt.Errorf("%s", result.Message)
		}
	}

	return string(respBody), nil
}

// HayFrpInfo HayFrp服务统计
type HayFrpInfo struct {
	Status    int     `json:"status"`
	Aflow     string  `json:"aflow"`
	Aflowin   string  `json:"aflowin"`
	Aflowout  string  `json:"aflowout"`
	Eflow     string  `json:"eflow"`
	Eflowin   string  `json:"eflowin"`
	Eflowout  string  `json:"eflowout"`
	Oclient   int     `json:"oclient"`
	Totalrun  string  `json:"totalrun"`
	Todayrun  string  `json:"todayrun"`
}

// GetHayFrpInfo 获取HayFrp服务统计
func (c *NodeAPIClient) GetHayFrpInfo() (*HayFrpInfo, error) {
	httpReq, err := http.NewRequest("GET", BaseURL+"/info", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result HayFrpInfo
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// DownloadListItem 下载列表项
type DownloadListItem struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Platform string `json:"platform"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
}

// DownloadSource 下载源
type DownloadSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// DownloadLists 下载列表集合
type DownloadLists struct {
	Frps   []DownloadListItem `json:"frps"`
	Frpc   []DownloadListItem `json:"frpc"`
	Others []DownloadListItem `json:"others"`
}

// DownloadListResponse 下载列表响应
type DownloadListResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Sources []DownloadSource `json:"sources"`
	Lists   DownloadLists  `json:"lists"`
}

// GetDownloadList 获取下载列表
func (c *NodeAPIClient) GetDownloadList() (*DownloadListResponse, error) {
	httpReq, err := http.NewRequest("GET", BaseURL+"/downlist", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result DownloadListResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// VersionInfo 版本信息
type VersionInfo struct {
	VerHayfrps   string `json:"ver_hayfrps"`
	VerFrpc      string `json:"ver_frpc"`
	VerLauncher  string `json:"ver_launcher"`
	VerConsole   string `json:"ver_console"`
	VerDashboard string `json:"ver_dashboard"`
	UrlLauncher  string `json:"url_launcher"`
}

// GetVersion 获取版本信息
func (c *NodeAPIClient) GetVersion() (*VersionInfo, error) {
	httpReq, err := http.NewRequest("GET", BaseURL+"/version", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result VersionInfo
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}
