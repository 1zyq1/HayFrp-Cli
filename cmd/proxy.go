package cmd

import (
	"fmt"
	"os"

	"hayfrp-cli/api"

	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "隧道相关操作",
	Long:  `隧道管理相关操作，包括创建、删除、编辑、查看隧道等`,
}

var addProxyCmd = &cobra.Command{
	Use:   "add [csrf]",
	Short: "添加隧道",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]

		name, _ := cmd.Flags().GetString("name")
		proxyType, _ := cmd.Flags().GetString("type")
		localIP, _ := cmd.Flags().GetString("local-ip")
		localPort, _ := cmd.Flags().GetInt("local-port")
		remotePort, _ := cmd.Flags().GetInt("remote-port")
		node, _ := cmd.Flags().GetString("node")
		domain, _ := cmd.Flags().GetString("domain")
		encryption, _ := cmd.Flags().GetBool("encryption")
		compression, _ := cmd.Flags().GetBool("compression")
		sk, _ := cmd.Flags().GetString("sk")

		if name == "" {
			fmt.Println("✗ 隧道名称不能为空")
			return
		}
		if proxyType == "" {
			fmt.Println("✗ 隧道类型不能为空 (tcp/udp/http/https/xtcp/stcp)")
			return
		}
		if localIP == "" {
			localIP = "127.0.0.1"
		}
		if localPort == 0 {
			fmt.Println("✗ 本地端口不能为空")
			return
		}
		if node == "" {
			fmt.Println("✗ 节点ID不能为空")
			return
		}

		// 布尔值转换为字符串
		encryptionStr := "false"
		if encryption {
			encryptionStr = "true"
		}
		compressionStr := "false"
		if compression {
			compressionStr = "true"
		}

		client := api.NewProxyAPIClient()
		req := &api.AddTunnelRequest{
			Type:            "add",
			Csrf:            csrf,
			ProxyName:       name,
			ProxyType:       proxyType,
			LocalIP:         localIP,
			LocalPort:       localPort,
			RemotePort:      remotePort,
			UseEncryption:   encryptionStr,
			UseCompression:  compressionStr,
			SK:              sk,
			Node:            node,
			Domain:          domain,
		}

		resp, err := client.AddTunnel(req)
		if err != nil {
			fmt.Printf("添加隧道失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
			fmt.Printf("  隧道ID: %s\n", resp.ID)
		} else {
			fmt.Printf("✗ 添加隧道失败: %s\n", resp.Message)
		}
	},
}

var editProxyCmd = &cobra.Command{
	Use:   "edit [csrf] [proxy-id]",
	Short: "编辑隧道",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]
		proxyID := args[1]

		name, _ := cmd.Flags().GetString("name")
		proxyType, _ := cmd.Flags().GetString("type")
		localIP, _ := cmd.Flags().GetString("local-ip")
		localPort, _ := cmd.Flags().GetInt("local-port")
		remotePort, _ := cmd.Flags().GetInt("remote-port")
		node, _ := cmd.Flags().GetString("node")
		domain, _ := cmd.Flags().GetString("domain")
		encryption, _ := cmd.Flags().GetBool("encryption")
		compression, _ := cmd.Flags().GetBool("compression")
		sk, _ := cmd.Flags().GetString("sk")

		if name == "" {
			fmt.Println("✗ 隧道名称不能为空")
			return
		}
		if proxyType == "" {
			fmt.Println("✗ 隧道类型不能为空 (tcp/udp/http/https/xtcp/stcp)")
			return
		}
		if localIP == "" {
			localIP = "127.0.0.1"
		}
		if localPort == 0 {
			fmt.Println("✗ 本地端口不能为空")
			return
		}
		if node == "" {
			fmt.Println("✗ 节点ID不能为空")
			return
		}

		encryptionStr := "false"
		if encryption {
			encryptionStr = "true"
		}
		compressionStr := "false"
		if compression {
			compressionStr = "true"
		}

		client := api.NewProxyAPIClient()
		req := &api.EditTunnelRequest{
			Type:            "edit",
			Csrf:            csrf,
			ID:              proxyID,
			ProxyName:       name,
			ProxyType:       proxyType,
			LocalIP:         localIP,
			LocalPort:       localPort,
			RemotePort:      remotePort,
			UseEncryption:   encryptionStr,
			UseCompression:  compressionStr,
			SK:              sk,
			Node:            node,
			Domain:          domain,
		}

		resp, err := client.EditTunnel(req)
		if err != nil {
			fmt.Printf("编辑隧道失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 编辑隧道失败: %s\n", resp.Message)
		}
	},
}

var deleteProxyCmd = &cobra.Command{
	Use:   "delete [csrf] [proxy-id]",
	Short: "删除隧道",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]
		proxyID := args[1]

		client := api.NewProxyAPIClient()
		resp, err := client.DeleteTunnel(csrf, proxyID)
		if err != nil {
			fmt.Printf("删除隧道失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 删除隧道失败: %s\n", resp.Message)
		}
	},
}

var listProxyCmd = &cobra.Command{
	Use:   "list [csrf] [proxy-id]",
	Short: "列出隧道",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]
		proxyID := ""
		if len(args) > 1 {
			proxyID = args[1]
		}

		client := api.NewProxyAPIClient()
		resp, err := client.ListTunnel(csrf, proxyID)
		if err != nil {
			fmt.Printf("列出隧道失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			if len(resp.Proxies) == 0 {
				fmt.Println("暂无隧道")
				return
			}
			for _, p := range resp.Proxies {
				fmt.Printf("========== 隧道 ==========\n")
				fmt.Printf("ID: %s\n", p.ID)
				fmt.Printf("名称: %s\n", p.ProxyName)
				fmt.Printf("类型: %s\n", p.ProxyType)
				fmt.Printf("本地地址: %s:%s\n", p.LocalIP, p.LocalPort)
				fmt.Printf("远程端口: %s\n", p.RemotePort)
				fmt.Printf("节点: %s (%s)\n", p.NodeName, p.Node)
				fmt.Printf("节点域名: %s\n", p.NodeDomain)
				if p.Domain != "" {
					fmt.Printf("域名: %s\n", p.Domain)
				}
				if p.SK != "" {
					fmt.Printf("SK密钥: %s\n", p.SK)
				}
				fmt.Printf("加密: %s\n", p.UseEncryption)
				fmt.Printf("压缩: %s\n", p.UseCompression)
				fmt.Printf("状态: %s\n", mapStatus(p.Status))
				fmt.Printf("最后更新: %s\n", p.LastUpdate)
				fmt.Println()
			}
		} else {
			fmt.Printf("✗ 列出隧道失败: %s\n", resp.Message)
		}
	},
}

var configProxyCmd = &cobra.Command{
	Use:   "config [csrf]",
	Short: "获取隧道配置文件",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]

		format, _ := cmd.Flags().GetString("format")
		node, _ := cmd.Flags().GetString("node")
		proxy, _ := cmd.Flags().GetString("proxy")
		output, _ := cmd.Flags().GetString("output")

		if format == "" {
			format = "ini"
		}
		if node == "" && proxy == "" {
			fmt.Println("✗ 请指定节点ID或隧道ID")
			return
		}

		client := api.NewProxyAPIClient()
		config, err := client.GetTunnelConfig(format, csrf, node, proxy)
		if err != nil {
			fmt.Printf("获取配置失败: %v\n", err)
			return
		}

		if output != "" {
			err := os.WriteFile(output, []byte(config), 0644)
			if err != nil {
				fmt.Printf("写入文件失败: %v\n", err)
				return
			}
			fmt.Printf("✓ 配置已保存到: %s\n", output)
		} else {
			fmt.Println(config)
		}
	},
}

var toggleProxyCmd = &cobra.Command{
	Use:   "toggle [csrf] [proxy-id] [true/false]",
	Short: "切换隧道状态",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]
		proxyID := args[1]
		toggle := args[2]

		if toggle != "true" && toggle != "false" {
			fmt.Println("✗ 状态必须是 true 或 false")
			return
		}

		client := api.NewProxyAPIClient()
		resp, err := client.ToggleTunnel(csrf, proxyID, toggle)
		if err != nil {
			fmt.Printf("切换隧道状态失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 切换隧道状态失败: %s\n", resp.Message)
		}
	},
}

var checkProxyCmd = &cobra.Command{
	Use:   "check [csrf] [proxy-id]",
	Short: "检查隧道状态",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]
		proxyID := args[1]

		client := api.NewProxyAPIClient()
		resp, err := client.CheckTunnel(csrf, proxyID)
		if err != nil {
			fmt.Printf("检查隧道状态失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s (状态: %s)\n", resp.Message, resp.OStatus)
		} else {
			fmt.Printf("✗ %s\n", resp.Message)
		}
	},
}

var forceDownProxyCmd = &cobra.Command{
	Use:   "force-down [csrf] [proxy-id]",
	Short: "强制下线隧道",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]
		proxyID := args[1]

		client := api.NewProxyAPIClient()
		resp, err := client.ForceDown(csrf, proxyID)
		if err != nil {
			fmt.Printf("强制下线隧道失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 强制下线隧道失败: %s\n", resp.Message)
		}
	},
}

func mapStatus(status string) string {
	if status == "true" {
		return "启用"
	}
	return "禁用"
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.AddCommand(addProxyCmd)
	proxyCmd.AddCommand(editProxyCmd)
	proxyCmd.AddCommand(deleteProxyCmd)
	proxyCmd.AddCommand(listProxyCmd)
	proxyCmd.AddCommand(configProxyCmd)
	proxyCmd.AddCommand(toggleProxyCmd)
	proxyCmd.AddCommand(checkProxyCmd)
	proxyCmd.AddCommand(forceDownProxyCmd)

	// add proxy flags
	addProxyCmd.Flags().String("name", "", "隧道名称")
	addProxyCmd.Flags().String("type", "", "隧道类型 (tcp/udp/http/https/xtcp/stcp)")
	addProxyCmd.Flags().String("local-ip", "", "本地IP (默认: 127.0.0.1)")
	addProxyCmd.Flags().Int("local-port", 0, "本地端口")
	addProxyCmd.Flags().Int("remote-port", 0, "远程端口")
	addProxyCmd.Flags().String("node", "", "节点ID")
	addProxyCmd.Flags().String("domain", "", "域名 (HTTP/HTTPS隧道)")
	addProxyCmd.Flags().Bool("encryption", false, "启用加密")
	addProxyCmd.Flags().Bool("compression", false, "启用压缩")
	addProxyCmd.Flags().String("sk", "", "SK密钥 (XTCP/STCP隧道)")

	// edit proxy flags
	editProxyCmd.Flags().String("name", "", "隧道名称")
	editProxyCmd.Flags().String("type", "", "隧道类型 (tcp/udp/http/https/xtcp/stcp)")
	editProxyCmd.Flags().String("local-ip", "", "本地IP (默认: 127.0.0.1)")
	editProxyCmd.Flags().Int("local-port", 0, "本地端口")
	editProxyCmd.Flags().Int("remote-port", 0, "远程端口")
	editProxyCmd.Flags().String("node", "", "节点ID")
	editProxyCmd.Flags().String("domain", "", "域名 (HTTP/HTTPS隧道)")
	editProxyCmd.Flags().Bool("encryption", false, "启用加密")
	editProxyCmd.Flags().Bool("compression", false, "启用压缩")
	editProxyCmd.Flags().String("sk", "", "SK密钥 (XTCP/STCP隧道)")

	// config proxy flags
	configProxyCmd.Flags().String("format", "ini", "配置文件格式 (ini/toml)")
	configProxyCmd.Flags().String("node", "", "节点ID")
	configProxyCmd.Flags().String("proxy", "", "隧道ID")
	configProxyCmd.Flags().String("output", "", "输出文件路径")
}
