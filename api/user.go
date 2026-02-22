package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UserAPIClient 用户相关API客户端
type UserAPIClient struct {
	client *http.Client
}

// NewUserAPIClient 创建用户API客户端
func NewUserAPIClient() *UserAPIClient {
	return &UserAPIClient{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Type   string `json:"type"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// Login 登录获取Token
func (c *UserAPIClient) Login(user, passwd string) (*LoginResponse, error) {
	req := LoginRequest{
		Type:   "login",
		User:   user,
		Passwd: passwd,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result LoginResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// CsrfRequest 验证Token请求
type CsrfRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
}

// CsrfResponse 验证Token响应
type CsrfResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// VerifyCsrf 验证Token是否有效
func (c *UserAPIClient) VerifyCsrf(csrf string) (*CsrfResponse, error) {
	req := CsrfRequest{
		Type: "csrf",
		Csrf: csrf,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result CsrfResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// SendRegCodeRequest 发送注册验证码请求
type SendRegCodeRequest struct {
	Type   string `json:"type"`
	User   string `json:"user"`
	Device string `json:"device"`
	Email  string `json:"email"`
}

// SendRegCodeResponse 发送注册验证码响应
type SendRegCodeResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// SendRegCode 发送注册验证码
func (c *UserAPIClient) SendRegCode(user, device, email string) (*SendRegCodeResponse, error) {
	req := SendRegCodeRequest{
		Type:   "sendregcode",
		User:   user,
		Device: device,
		Email:  email,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result SendRegCodeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Device  string `json:"device"`
	Email   string `json:"email"`
	Passwd  string `json:"passwd"`
	Regcode string `json:"regcode"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Register 注册新用户
func (c *UserAPIClient) Register(user, device, email, passwd, regcode string) (*RegisterResponse, error) {
	req := RegisterRequest{
		Type:    "register",
		User:    user,
		Device:  device,
		Email:   email,
		Passwd:  passwd,
		Regcode: regcode,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result RegisterResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetInfoRequest 获取用户信息请求
type GetInfoRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
}

// UserInfo 用户信息
type UserInfo struct {
	Status       bool        `json:"status"`
	Message      string      `json:"message"`
	ID           interface{} `json:"id"`
	Username     string      `json:"username"`
	Token        string      `json:"token"`
	Email        string      `json:"email"`
	Traffic      interface{} `json:"traffic"`
	Realname     interface{} `json:"realname"`     // 可能是 string 或 bool
	Proxies      interface{} `json:"proxies"`
	Useproxies   interface{} `json:"useproxies"`
	Regtime      interface{} `json:"regtime"`
	Signdate     string      `json:"signdate"`
	Totalsign    interface{} `json:"totalsign"`
	Totaltraffic interface{} `json:"totaltraffic"`
	Todaytraffic interface{} `json:"todaytraffic"`
	Qid          interface{} `json:"qid"`
	Sprovider    interface{} `json:"sprovider"`    // 可能是 string 或 bool
	UUID         string      `json:"uuid"`
}

// GetInfo 获取用户信息
func (c *UserAPIClient) GetInfo(csrf string) (*UserInfo, error) {
	req := GetInfoRequest{
		Type: "info",
		Csrf: csrf,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result UserInfo
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// SignRequest 签到请求
type SignRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
}

// SignResponse 签到响应
type SignResponse struct {
	Status   int     `json:"status"`
	Message  string  `json:"message"`
	Signflow float64 `json:"signflow"`
	Flow     float64 `json:"flow"`
}

// Sign 签到
func (c *UserAPIClient) Sign(csrf string) (*SignResponse, error) {
	req := SignRequest{
		Type: "sign",
		Csrf: csrf,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result SignResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// ReTokenRequest 更新Token请求
type ReTokenRequest struct {
	Type string `json:"type"`
	Csrf string `json:"csrf"`
}

// ReTokenResponse 更新Token响应
type ReTokenResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// ReToken 更新用户Token
func (c *UserAPIClient) ReToken(csrf string) (*ReTokenResponse, error) {
	req := ReTokenRequest{
		Type: "retoken",
		Csrf: csrf,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result ReTokenResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// FindPassEmRequest 重置密码发送验证码请求
type FindPassEmRequest struct {
	Type string `json:"type"`
	User string `json:"user"`
}

// FindPassEmResponse 重置密码发送验证码响应
type FindPassEmResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// SendFindPassCode 发送重置密码验证码
func (c *UserAPIClient) SendFindPassCode(user string) (*FindPassEmResponse, error) {
	req := FindPassEmRequest{
		Type: "findpassem",
		User: user,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result FindPassEmResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// FindPassCtRequest 重置密码请求
type FindPassCtRequest struct {
	Type    string `json:"type"`
	Token   string `json:"token"`
	Newpass string `json:"newpass"`
}

// FindPassCtResponse 重置密码响应
type FindPassCtResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ResetPass 重置密码
func (c *UserAPIClient) ResetPass(token, newpass string) (*FindPassCtResponse, error) {
	req := FindPassCtRequest{
		Type:    "findpassct",
		Token:   token,
		Newpass: newpass,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/user", bytes.NewBuffer(body))
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

	var result FindPassCtResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}
