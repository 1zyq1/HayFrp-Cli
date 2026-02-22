# HayFrp CLI

HayFrp CLI 是一个命令行工具，用于管理 HayFrp 内网穿透服务。

## 功能特性

- **用户管理**：登录、注册、签到、查看信息、重置密码等
- **隧道管理**：创建、编辑、删除、列表、配置文件获取、状态切换等
- **节点查询**：节点列表、节点信息、服务统计等
- **自动登录**：保存登录状态，下次自动登录
- **自动下载**：自动下载对应平台的 frpc 并启动
- **API 容灾**：多端点自动故障转移

## 安装

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/yourusername/hayfrp-cli.git
cd hayfrp-cli

# 构建可执行文件
GOPROXY=https://goproxy.cn go build -o hayfrp main.go

# 安装到系统
sudo mv hayfrp /usr/local/bin/
```

### 直接下载

从 Releases 页面下载对应平台的可执行文件。

## 快速开始

### 启动隧道（推荐）

交互式启动流程，自动完成登录、选择隧道、启动：

```bash
hayfrp start
```

流程：
1. 输入用户名密码登录（首次）
2. 自动保存登录状态
3. 显示用户信息
4. 列出可用隧道
5. 选择要启动的隧道
6. 自动下载 frpc（如未安装）
7. 启动隧道

### 退出登录

```bash
hayfrp logout
```

## 命令详解

### 用户操作

```bash
# 用户登录
hayfrp user login username password

# 验证 Token
hayfrp user verify your-token

# 获取用户信息
hayfrp user info your-token

# 每日签到
hayfrp user sign your-token

# 更新 Token
hayfrp user retoken your-token

# 发送注册验证码
hayfrp user send-reg username email@example.com

# 用户注册
hayfrp user register username email@example.com password code

# 发送重置密码验证码
hayfrp user send-findpass username

# 重置密码
hayfrp user reset-pass token new-password
```

### 隧道操作

```bash
# 添加隧道
hayfrp proxy add your-token \
  --name "my-tunnel" \
  --type tcp \
  --local-ip 127.0.0.1 \
  --local-port 8080 \
  --remote-port 8081 \
  --node 1 \
  --encryption \
  --compression

# 编辑隧道
hayfrp proxy edit your-token proxy-id \
  --name "updated-tunnel" \
  --type tcp \
  --local-ip 127.0.0.1 \
  --local-port 8080 \
  --remote-port 8081 \
  --node 1

# 删除隧道
hayfrp proxy delete your-token proxy-id

# 列出所有隧道
hayfrp proxy list your-token

# 查看指定隧道
hayfrp proxy list your-token proxy-id

# 获取配置文件
hayfrp proxy config your-token --node 1 --format ini --output frpc.ini

# 切换隧道状态
hayfrp proxy toggle your-token proxy-id true

# 检查隧道状态
hayfrp proxy check your-token proxy-id

# 强制下线隧道
hayfrp proxy force-down your-token proxy-id
```

### 节点操作

```bash
# 获取节点探针信息
hayfrp node info

# 获取节点列表
hayfrp node list

# 获取公告
hayfrp node notice

# 获取 HayFrp 服务统计
hayfrp node info

# 获取下载列表
hayfrp node download

# 获取版本信息
hayfrp node version
```

## 配置文件

### 会话文件

登录状态保存在 `~/.hayfrp/session.json`：

```json
{
  "csrf": "your-csrf-token",
  "username": "myuser",
  "login_time": "2024-01-01T12:00:00Z"
}
```

### 配置文件

隧道配置保存在 `~/.hayfrp/frpc.toml`。

## API 容灾

CLI 内置多端点容灾机制，按优先级使用以下端点：

1. `https://api.hayfrp.1zyq1.com`
2. `https://v2.api.hayfrp.1zyq1.com`
3. `https://api.hayfrp.com`

当一个端点失败时，自动切换到下一个端点。

## 隧道类型

| 类型 | 说明 | 需要远程端口 | 需要域名 |
|------|------|-------------|---------|
| tcp | TCP 隧道 | 是 | 否 |
| udp | UDP 隧道 | 是 | 否 |
| http | HTTP 隧道 | 否 | 是 |
| https | HTTPS 隧道 | 否 | 是 |
| stcp | 密钥 TCP | 否 | 否 |
| xtcp | P2P TCP | 否 | 否 |

## 文件结构

```
~/.hayfrp/
├── session.json    # 登录会话
├── frpc.toml       # 隧道配置
└── frpc            # frpc 可执行文件
```

## 常见问题

### Token 过期

Token 有效期为 7 天，过期后需要重新登录：

```bash
hayfrp logout
hayfrp start
```

### frpc 未找到

使用 `hayfrp start` 会自动下载对应平台的 frpc。

### 网络问题

CLI 会自动尝试多个 API 端点，确保服务可用。

## 开发

### 项目结构

```
HayFrp-cli/
├── main.go           # 入口文件
├── go.mod            # 依赖管理
├── api/              # API 封装层
│   ├── client.go     # 公共客户端和容灾逻辑
│   ├── user.go       # 用户相关 API
│   ├── proxy.go      # 隧道相关 API
│   └── node.go       # 节点相关 API
└── cmd/              # 命令行定义
    ├── root.go       # 根命令
    ├── start.go      # 启动命令
    ├── user.go       # 用户命令
    ├── proxy.go      # 隧道命令
    └── node.go       # 节点命令
```

### 构建

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o hayfrp-linux-amd64 main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o hayfrp-windows-amd64.exe main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o hayfrp-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o hayfrp-darwin-arm64 main.go
```

## 许可证

MIT License

## 相关链接

- [HayFrp 官网](https://hayfrp.com)
- [API 文档](./api-doc.md)
