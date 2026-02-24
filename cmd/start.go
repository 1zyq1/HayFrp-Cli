package cmd

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"hayfrp-cli/api"

	"github.com/spf13/cobra"
)

// SavedSession 保存的会话信息
type SavedSession struct {
	CSRF      string    `json:"csrf"`
	Username  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动隧道（交互式）",
	Long:  `交互式启动流程：登录 -> 选择隧道 -> 启动隧道`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		homeDir, _ := os.UserHomeDir()
		configDir := filepath.Join(homeDir, ".hayfrp")
		sessionFile := filepath.Join(configDir, "session.json")

		// 主循环：支持退出账户后重新登录
		for {

		// 步骤1: 尝试自动登录
		fmt.Println("========== HayFrp 隧道启动器 ==========")

		csrf := ""
		userClient := api.NewUserAPIClient()

		// 主循环：支持退出账户后重新登录
		for {
			// 如果未登录，尝试登录
			if csrf == "" {
				// 尝试读取保存的会话
				if session := loadSession(sessionFile); session != nil {
					fmt.Printf("检测到保存的登录信息 (用户: %s)\n", session.Username)
					fmt.Print("正在验证 Token 有效性... ")

					// 验证 token 是否有效
					verifyResp, err := userClient.VerifyCsrf(session.CSRF)
					if err == nil && verifyResp.Status == 200 {
						fmt.Println("有效!")
						csrf = session.CSRF
						fmt.Printf("✓ 自动登录成功！\n\n")
					} else {
						fmt.Println("已过期")
						fmt.Println("请重新登录")
					}
				}

				// 如果自动登录失败，手动登录
				for csrf == "" {
					fmt.Print("用户名/邮箱: ")
					username, _ := reader.ReadString('\n')
					username = strings.TrimSpace(username)

					password, err := readPasswordWithMask("密码")
					if err != nil {
						fmt.Printf("读取密码失败: %v\n", err)
						continue
					}

					loginResp, err := userClient.Login(username, password)
					if err != nil {
						fmt.Printf("✗ 登录失败: %v\n", err)
						continue
					}

					if loginResp.Status != 200 {
						fmt.Printf("✗ 登录失败: %s\n", loginResp.Message)
						continue
					}

					csrf = loginResp.Token

					// 保存会话
					os.MkdirAll(configDir, 0755)
					session := &SavedSession{
						CSRF:      csrf,
						Username:  username,
						LoginTime: time.Now(),
					}
					if err := saveSession(sessionFile, session); err == nil {
						fmt.Printf("✓ 登录成功！(已保存登录状态)\n\n")
					} else {
						fmt.Printf("✓ 登录成功！\n\n")
					}
				}
			}

			// 步骤2: 获取用户信息
			infoResp, err := userClient.GetInfo(csrf)
			if err != nil {
				fmt.Printf("✗ 获取用户信息失败: %v\n", err)
				// 返回登录流程
				csrf = ""
				continue
			}
			if !infoResp.Status {
				fmt.Printf("✗ %s\n", infoResp.Message)
				// 返回登录流程
				csrf = ""
				continue
			}
			fmt.Printf("========== 用户信息 ==========\n")
			fmt.Printf("用户: %s\n", infoResp.Username)
			// 转换流量为GB
			var trafficGB float64
			switch v := infoResp.Traffic.(type) {
			case string:
				if val, err := strconv.ParseFloat(v, 64); err == nil {
					trafficGB = val / 1024
				}
			case float64:
				trafficGB = v / 1024
			}
			fmt.Printf("剩余流量: %.2f GB\n", trafficGB)
			fmt.Printf("拥有隧道: %v / 已使用: %v\n", infoResp.Proxies, infoResp.Useproxies)
			fmt.Printf("================================\n")
			fmt.Printf("[0] 退出账户\n\n")

			// 隧道选择循环
			proxyClient := api.NewProxyAPIClient()
			for {
				// 检查是否已登录
				if csrf == "" {
					break
				}
				// 步骤3: 获取隧道列表
				listResp, err := proxyClient.ListTunnel(csrf, "")
				if err != nil {
					fmt.Printf("✗ 获取隧道列表失败: %v\n", err)
					fmt.Print("\n按任意键重试...")
					reader.ReadString('\n')
					continue
				}

				if listResp.Status != 200 || len(listResp.Proxies) == 0 {
					fmt.Println("✗ 暂无可用隧道，请先在控制台创建隧道")
					fmt.Print("\n按任意键重试...")
					reader.ReadString('\n')
					continue
				}

				fmt.Println("========== 可用隧道列表 ==========")
				for i, p := range listResp.Proxies {
					status := "禁用"
					if p.Status == "true" {
						status = "启用"
					}
					fmt.Printf("%d. [%s] %s (%s)\n", i+1, p.ProxyType, p.ProxyName, status)
					fmt.Printf("   节点: %s\n", p.NodeName)
					fmt.Printf("   本地: %s:%s -> 远程: %s\n", p.LocalIP, p.LocalPort, p.RemotePort)
					if p.Domain != "" {
						fmt.Printf("   域名: %s\n", p.Domain)
					}
				}
				fmt.Println("================================")

				// 步骤4: 选择隧道
				fmt.Print("\n请选择要启动的隧道编号 [0退出]: ")
				choice, _ := reader.ReadString('\n')
				choice = strings.TrimSpace(choice)

				// 检查是否选择退出
				if choice == "0" {
					// 退出账户
					fmt.Print("\n确认退出账户? (y/n): ")
					confirm, _ := reader.ReadString('\n')
					confirm = strings.TrimSpace(confirm)
					if strings.ToLower(confirm) == "y" {
						// 删除会话文件
						os.Remove(sessionFile)
						fmt.Println("✓ 已退出账户")
						// 返回到登录流程
						csrf = ""
						break
					} else {
						// 用户取消退出，继续显示隧道列表
						continue
					}
				}

				var choiceIndex int
				_, err = fmt.Sscanf(choice, "%d", &choiceIndex)
				if err != nil || choiceIndex < 1 || choiceIndex > len(listResp.Proxies) {
					fmt.Println("✗ 无效的选择")
					fmt.Print("\n按任意键重试...")
					reader.ReadString('\n')
					continue
				}

				selectedProxy := listResp.Proxies[choiceIndex-1]

				// 检查隧道状态
				if selectedProxy.Status != "true" {
					fmt.Printf("隧道 %s 当前状态为禁用，正在启用...\n", selectedProxy.ProxyName)
					toggleResp, err := proxyClient.ToggleTunnel(csrf, selectedProxy.ID, "true")
					if err != nil {
						fmt.Printf("✗ 启用隧道失败: %v\n", err)
						fmt.Print("\n按任意键重试...")
						reader.ReadString('\n')
						continue
					}
					if toggleResp.Status != 200 {
						fmt.Printf("✗ 启用隧道失败: %s\n", toggleResp.Message)
						fmt.Print("\n按任意键重试...")
						reader.ReadString('\n')
						continue
					}
					fmt.Printf("✓ 隧道已启用\n")
				}

				// 步骤5: 生成配置文件
				fmt.Printf("\n正在为隧道 %s 生成配置文件...\n", selectedProxy.ProxyName)
				config, err := proxyClient.GetTunnelConfig("toml", csrf, "", selectedProxy.ID)
				if err != nil {
					fmt.Printf("✗ 生成配置文件失败: %v\n", err)
					fmt.Print("\n按任意键重试...")
					reader.ReadString('\n')
					continue
				}

				// 保存配置文件到用户目录
				if homeDir == "" {
					homeDir = "."
				}

				configDir = filepath.Join(homeDir, ".hayfrp")
				if err := os.MkdirAll(configDir, 0755); err != nil {
					configDir = "."
				}

				configFile := filepath.Join(configDir, "frpc.toml")
				if err := os.WriteFile(configFile, []byte(config), 0644); err != nil {
					fmt.Printf("✗ 保存配置文件失败: %v\n", err)
					fmt.Print("\n按任意键重试...")
					reader.ReadString('\n')
					continue
				}

				fmt.Printf("✓ 配置文件已保存: %s\n", configFile)

				// 步骤6: 启动frpc
				fmt.Printf("\n========== 启动frpc ==========\n")

				// 检查frpc可执行文件
				frpcPath := ""
				frpcName := "frpc"
				if runtime.GOOS == "windows" {
					frpcName = "frpc.exe"
				}

				possiblePaths := []string{
					filepath.Join(".", frpcName),
					filepath.Join(homeDir, ".hayfrp", frpcName),
				}

				// Unix 系统额外路径
				if runtime.GOOS != "windows" {
					possiblePaths = append(possiblePaths,
						"/usr/local/bin/frpc",
						"/usr/bin/frpc",
					)
				}

				for _, path := range possiblePaths {
					if _, err := os.Stat(path); err == nil {
						frpcPath = path
						break
					}
				}

				if frpcPath == "" {
					fmt.Println("未找到 frpc 可执行文件，正在尝试自动下载...")

					// 自动下载 frpc
					downloadResp, err := downloadFrpc(homeDir, configDir)
					if err != nil {
						fmt.Printf("✗ 自动下载 frpc 失败: %v\n", err)
						fmt.Println("\n请手动下载 frpc:")

						// 获取下载列表
						nodeClient := api.NewNodeAPIClient()
						downloadList, err := nodeClient.GetDownloadList()
						if err == nil && downloadList.Status == 200 {
							fmt.Println("\n下载源:")
							for _, source := range downloadList.Sources {
								fmt.Printf("  - %s: %s\n", source.Name, source.URL)
							}

							fmt.Println("\n推荐下载:")
							osName := runtime.GOOS
							arch := runtime.GOARCH

							fmt.Printf("  系统: %s, 架构: %s\n", osName, arch)
							for _, item := range downloadList.Lists.Frpc {
								if strings.ToLower(item.Platform) == strings.ToLower(osName) &&
									strings.Contains(strings.ToLower(item.Arch), strings.ToLower(arch)) {
									fmt.Printf("  - %s (版本: %s)\n", item.Name, item.Version)
									for _, source := range downloadList.Sources {
										fmt.Printf("    下载: %s%s\n", source.URL, item.URL)
									}
								}
							}
						}

						fmt.Printf("\n下载后请将 frpc 放到以下任一路径:\n")
						for _, path := range possiblePaths {
							fmt.Printf("  - %s\n", path)
						}
						fmt.Print("\n按任意键重试...")
						reader.ReadString('\n')
						continue
					}

					frpcPath = downloadResp
					fmt.Printf("✓ frpc 下载成功: %s\n", frpcPath)
				}

				fmt.Printf("使用 frpc: %s\n", frpcPath)
				fmt.Printf("配置文件: %s\n", configFile)
				fmt.Println("\n按 Ctrl+C 可停止隧道")
				fmt.Println("================================\n")

				// 启动frpc
				frpcExec := exec.Command(frpcPath, "-c", configFile)
				frpcExec.Stdout = os.Stdout
				frpcExec.Stderr = os.Stderr

				if err := frpcExec.Run(); err != nil {
					fmt.Printf("\n✗ frpc 启动失败: %v\n", err)
					fmt.Print("\n按任意键返回隧道列表...")
					reader.ReadString('\n')
				}
			}
		}
		}
	},
}

// downloadFrpc 自动下载对应平台的 frpc
func downloadFrpc(homeDir, configDir string) (string, error) {
	nodeClient := api.NewNodeAPIClient()
	downloadList, err := nodeClient.GetDownloadList()
	if err != nil {
		return "", fmt.Errorf("获取下载列表失败: %w", err)
	}

	if downloadList.Status != 200 || len(downloadList.Lists.Frpc) == 0 {
		return "", fmt.Errorf("未找到可用下载列表")
	}

	// 确定系统和架构
	osName := runtime.GOOS
	arch := runtime.GOARCH

	fmt.Printf("检测到系统: %s, 架构: %s\n", osName, arch)

	// 查找匹配的 frpc
	var matchedItem *api.DownloadListItem
	for i := range downloadList.Lists.Frpc {
		item := &downloadList.Lists.Frpc[i]
		// 匹配系统和架构
		osMatch := strings.ToLower(item.Platform) == strings.ToLower(osName)
		archMatch := strings.Contains(strings.ToLower(item.Arch), strings.ToLower(arch))
		if osMatch && archMatch {
			matchedItem = item
			break
		}
	}

	if matchedItem == nil {
		return "", fmt.Errorf("未找到匹配当前架构 %s 的 frpc 版本", arch)
	}

	// 使用第一个下载源
	if len(downloadList.Sources) == 0 {
		return "", fmt.Errorf("未找到下载源")
	}
	source := downloadList.Sources[0]

	// 清理 URL 中的换行符
	downloadURL := strings.TrimSpace(source.URL) + strings.TrimSpace(matchedItem.URL)
	fmt.Printf("版本: %s\n", matchedItem.Version)
	fmt.Printf("下载地址: %s\n", downloadURL)

	// 下载文件
	fmt.Println("正在下载 frpc...")
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 保存文件
	tempFile := filepath.Join(configDir, "frpc_download"+getFileExt(matchedItem.URL))
	file, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 显示进度
	size := resp.ContentLength
	buf := make([]byte, 32*1024)
	var downloaded int64
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			file.Write(buf[:n])
			downloaded += int64(n)
			if size > 0 {
				percent := float64(downloaded) / float64(size) * 100
				fmt.Printf("\r下载进度: %.2f%%", percent)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("下载失败: %w", err)
		}
	}
	fmt.Println("\r下载进度: 100%")

	file.Close()

	// 处理压缩包
	urlClean := strings.TrimSpace(matchedItem.URL)
	if strings.HasSuffix(urlClean, ".tar.gz") {
		fmt.Println("正在解压 tar.gz 文件...")
		if err := extractTarGz(tempFile, configDir); err != nil {
			os.Remove(tempFile)
			return "", fmt.Errorf("解压失败: %w", err)
		}
		os.Remove(tempFile)
	} else if strings.HasSuffix(urlClean, ".zip") {
		fmt.Println("正在解压 zip 文件...")
		if err := extractZip(tempFile, configDir); err != nil {
			os.Remove(tempFile)
			return "", fmt.Errorf("解压失败: %w", err)
		}
		os.Remove(tempFile)
	} else if strings.HasSuffix(urlClean, ".exe") {
		// Windows 可执行文件
		target := filepath.Join(configDir, "frpc.exe")
		os.Rename(tempFile, target)
	} else {
		// 直接是可执行文件
		target := filepath.Join(configDir, "frpc")
		os.Rename(tempFile, target)
	}

	// 检查 frpc 是否存在
	frpcPath := filepath.Join(configDir, "frpc")
	if _, err := os.Stat(frpcPath); os.IsNotExist(err) {
		// 尝试 frpc.exe (Windows)
		frpcPath = filepath.Join(configDir, "frpc.exe")
		if _, err := os.Stat(frpcPath); os.IsNotExist(err) {
			return "", fmt.Errorf("解压后未找到 frpc 可执行文件")
		}
	}

	// 设置执行权限 (Unix 系统需要)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(frpcPath, 0755); err != nil {
			return "", fmt.Errorf("设置执行权限失败: %w", err)
		}
	}

	return frpcPath, nil
}

// extractTarGz 解压 tar.gz 文件
func extractTarGz(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 只提取 frpc 文件
		name := filepath.Base(header.Name)
		if name == "frpc" || name == "frpc.exe" {
			target := filepath.Join(dest, name)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
			fmt.Printf("✓ 已解压: %s\n", target)
			return nil
		}
	}

	return fmt.Errorf("压缩包中未找到 frpc 文件")
}

// extractZip 解压 zip 文件
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		name := filepath.Base(f.Name)
		if name == "frpc" || name == "frpc.exe" {
			target := filepath.Join(dest, name)

			rc, err := f.Open()
			if err != nil {
				return err
			}

			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, f.Mode())
			if err != nil {
				rc.Close()
				return err
			}

			if _, err := io.Copy(file, rc); err != nil {
				file.Close()
				rc.Close()
				return err
			}

			file.Close()
			rc.Close()
			fmt.Printf("✓ 已解压: %s\n", target)
			return nil
		}
	}

	return fmt.Errorf("压缩包中未找到 frpc 文件")
}

// getFileExt 获取文件扩展名
func getFileExt(url string) string {
	if idx := strings.LastIndex(url, "."); idx != -1 {
		return url[idx:]
	}
	return ""
}

func init() {
	rootCmd.AddCommand(startCmd)
}

// logoutCmd 退出登录命令
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录",
	Long:  `清除保存的登录状态`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		sessionFile := filepath.Join(homeDir, ".hayfrp", "session.json")

		if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
			fmt.Println("当前没有保存的登录状态")
			return
		}

		if err := os.Remove(sessionFile); err != nil {
			fmt.Printf("✗ 退出登录失败: %v\n", err)
			return
		}

		fmt.Println("✓ 已退出登录")
	},
}

// loadSession 加载保存的会话
func loadSession(path string) *SavedSession {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var session SavedSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil
	}

	return &session
}

// saveSession 保存会话
func saveSession(path string, session *SavedSession) error {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
