package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ProxyAPIClient 隧道相关API客户端
type ProxyAPIClient struct {
	client *http.Client
}

// NewProxyAPIClient 创建隧道API客户端
func NewProxyAPIClient() *ProxyAPIClient {
	return &ProxyAPIClient{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// AddTunnelRequest 添加隧道请求
type AddTunnelRequest struct {
	Type                 string `json:"type"`
	Csrf                 string `json:"csrf"`
	ProxyName            string `json:"proxy_name"`
	ProxyType            string `json:"proxy_type"`
	LocalIP              string `json:"local_ip"`
	LocalPort            int    `json:"local_port"`
	RemotePort           int    `json:"remote_port"`
	UseEncryption        string `json:"use_encryption"`
	UseCompression       string `json:"use_compression"`
	SK                   string `json:"sk"`
	Node                 string `json:"node"`
	Domain               string `json:"domain"`
	Locations            string `json:"locations"`
	HeaderXFromWhere     string `json:"header_X_From_Where"`
	HostHeaderRewrite    string `json:"host_header_rewrite"`
}

// AddTunnelResponse 添加隧道响应
type AddTunnelResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

// AddTunnel 添加隧道
func (c *ProxyAPIClient) AddTunnel(req *AddTunnelRequest) (*AddTunnelResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("waf", "off")

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result AddTunnelResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// EditTunnelRequest 编辑隧道请求
type EditTunnelRequest struct {
	Type                 string `json:"type"`
	Csrf                 string `json:"csrf"`
	ID                   string `json:"id"`
	ProxyName            string `json:"proxy_name"`
	ProxyType            string `json:"proxy_type"`
	LocalIP              string `json:"local_ip"`
	LocalPort            int    `json:"local_port"`
	RemotePort           int    `json:"remote_port"`
	UseEncryption        string `json:"use_encryption"`
	UseCompression       string `json:"use_compression"`
	SK                   string `json:"sk"`
	Node                 string `json:"node"`
	Domain               string `json:"domain"`
	Locations            string `json:"locations"`
	HeaderXFromWhere     string `json:"header_X_From_Where"`
	HostHeaderRewrite    string `json:"host_header_rewrite"`
}

// EditTunnelResponse 编辑隧道响应
type EditTunnelResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// EditTunnel 编辑隧道
func (c *ProxyAPIClient) EditTunnel(req *EditTunnelRequest) (*EditTunnelResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("waf", "off")

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result EditTunnelResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// DeleteTunnelRequest 删除隧道请求
type DeleteTunnelRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
	ID   string `json:"id"`
}

// DeleteTunnelResponse 删除隧道响应
type DeleteTunnelResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// DeleteTunnel 删除隧道
func (c *ProxyAPIClient) DeleteTunnel(csrf, id string) (*DeleteTunnelResponse, error) {
	req := DeleteTunnelRequest{
		Type: "remove",
		Csrf: csrf,
		ID:   id,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("waf", "off")

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result DeleteTunnelResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// ListTunnelRequest 列出隧道请求
type ListTunnelRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
	ID   string `json:"id"`
}

// TunnelInfo 隧道信息
type TunnelInfo struct {
	ID                 string `json:"id"`
	UUID               string `json:"uuid"`
	Username           string `json:"username"`
	ProxyName          string `json:"proxy_name"`
	ProxyType          string `json:"proxy_type"`
	LocalIP            string `json:"local_ip"`
	LocalPort          string `json:"local_port"`
	UseEncryption      string `json:"use_encryption"`
	UseCompression     string `json:"use_compression"`
	Domain             string `json:"domain"`
	Locations          string `json:"locations"`
	HostHeaderRewrite  string `json:"host_header_rewrite"`
	RemotePort         string `json:"remote_port"`
	SK                 string `json:"sk"`
	HeaderXFromWhere   string `json:"header_X-From-Where"`
	Status             string `json:"status"`
	LastUpdate         string `json:"lastupdate"`
	Node               string `json:"node"`
	NodeName           string `json:"node_name"`
	NodeDomain         string `json:"node_domain"`
}

// ListTunnelResponse 列出隧道响应
type ListTunnelResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Proxies []TunnelInfo `json:"proxies"`
}

// ListTunnel 列出隧道
func (c *ProxyAPIClient) ListTunnel(csrf, id string) (*ListTunnelResponse, error) {
	req := ListTunnelRequest{
		Type: "list",
		Csrf: csrf,
		ID:   id,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("waf", "off")

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ListTunnelResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// TunnelConfigRequest 获取隧道配置请求
type TunnelConfigRequest struct {
	Type   string `json:"type"`
	Format string `json:"format"`
	Csrf   string `json:"csrf"`
	Node   string `json:"node"`
	Proxy  string `json:"proxy"`
}

// GetTunnelConfig 获取隧道配置文件
func (c *ProxyAPIClient) GetTunnelConfig(format, csrf, node, proxy string) (string, error) {
	req := TunnelConfigRequest{
		Type:   "config",
		Format: format,
		Csrf:   csrf,
		Node:   node,
		Proxy:  proxy,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("waf", "off")

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

// ToggleTunnelRequest 切换隧道状态请求
type ToggleTunnelRequest struct {
	Type   string `json:"type"`
	Csrf   string `json:"csrf"`
	ID     string `json:"id"`
	Toggle string `json:"toggle"`
}

// ToggleTunnelResponse 切换隧道状态响应
type ToggleTunnelResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ToggleTunnel 切换隧道状态
func (c *ProxyAPIClient) ToggleTunnel(csrf, id, toggle string) (*ToggleTunnelResponse, error) {
	req := ToggleTunnelRequest{
		Type:   "toggle",
		Csrf:   csrf,
		ID:     id,
		Toggle: toggle,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("waf", "off")

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ToggleTunnelResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// CheckTunnelRequest 检查隧道状态请求
type CheckTunnelRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
	ID   string `json:"id"`
}

// CheckTunnelResponse 检查隧道状态响应
type CheckTunnelResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	OStatus string `json:"ostatus"`
}

// CheckTunnel 检查隧道状态
func (c *ProxyAPIClient) CheckTunnel(csrf, id string) (*CheckTunnelResponse, error) {
	req := CheckTunnelRequest{
		Type: "check",
		Csrf: csrf,
		ID:   id,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("waf", "off")

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result CheckTunnelResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// ForceDownRequest 强制下线隧道请求
type ForceDownRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
	ID   string `json:"id"`
}

// ForceDownResponse 强制下线隧道响应
type ForceDownResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ForceDown 强制下线隧道
func (c *ProxyAPIClient) ForceDown(csrf, id string) (*ForceDownResponse, error) {
	// 使用 form-urlencoded 格式，与参考代码一致
	formData := fmt.Sprintf("type=forcedown&csrf=%s&id=%s", csrf, id)

	httpReq, err := http.NewRequest("POST", BaseURL+"/proxy", bytes.NewBufferString(formData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("waf", "off")

	resp, err := DoRequestWithFallback(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ForceDownResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}
