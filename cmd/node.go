package cmd

import (
	"fmt"

	"hayfrp-cli/api"

	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "节点相关操作",
	Long:  `节点查询相关操作，包括节点列表、节点信息等`,
}

var nodeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取节点探针信息",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewNodeAPIClient()
		resp, err := client.GetNodeInfo()
		if err != nil {
			fmt.Printf("获取节点信息失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("========== 节点探针信息 ==========\n")
			fmt.Printf("在线节点数: %d\n\n", resp.Number)

			for _, node := range resp.Servers {
				fmt.Printf("ID: %s\n", node.ID)
				fmt.Printf("名称: %s\n", node.Name)
				fmt.Printf("版本: %s\n", node.Version)
				fmt.Printf("绑定端口: %s\n", node.BindPort)
				fmt.Printf("HTTP端口: %s\n", node.VhostHTTPPort)
				fmt.Printf("HTTPS端口: %s\n", node.VhostHTTPSPort)
				fmt.Printf("连接数: %s\n", node.CurConns)
				fmt.Printf("客户端数: %s\n", node.ClientCounts)
				fmt.Printf("CPU使用率: %s\n", node.CPUUsage)
				fmt.Printf("内存使用率: %s\n", node.RAMUsage)
				fmt.Printf("磁盘使用率: %s\n", node.DiskUsage)
				fmt.Printf("今日入网: %s Bytes\n", node.TotalTrafficIn)
				fmt.Printf("今日出网: %s Bytes\n", node.TotalTrafficOut)
				fmt.Printf("状态: %s\n", node.Status)
				fmt.Println()
			}
		} else {
			fmt.Printf("✗ 获取节点信息失败: %s\n", resp.Message)
		}
	},
}

var nodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取节点列表",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewNodeAPIClient()
		resp, err := client.GetNodeList()
		if err != nil {
			fmt.Printf("获取节点列表失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("========== 节点列表 ==========\n")
			fmt.Printf("在线节点数: %d\n\n", resp.Number)

			for _, node := range resp.Servers {
				fmt.Printf("[%s] %s\n", node.ID, node.Name)
				fmt.Printf("    %s\n\n", node.Description)
			}
		} else {
			fmt.Printf("✗ 获取节点列表失败: %s\n", resp.Message)
		}
	},
}

var noticeCmd = &cobra.Command{
	Use:   "notice",
	Short: "获取公告",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewNodeAPIClient()
		notice, err := client.GetNotice()
		if err != nil {
			fmt.Printf("获取公告失败: %v\n", err)
			return
		}

		fmt.Println(notice)
	},
}

var hayfrpInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取HayFrp服务统计",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewNodeAPIClient()
		resp, err := client.GetHayFrpInfo()
		if err != nil {
			fmt.Printf("获取服务统计失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("========== HayFrp服务统计 ==========\n")
			fmt.Printf("总流量: %s MB\n", resp.Aflow)
			fmt.Printf("总入网流量: %s MB\n", resp.Aflowin)
			fmt.Printf("总出网流量: %s MB\n", resp.Aflowout)
			fmt.Printf("今日流量: %s MB\n", resp.Eflow)
			fmt.Printf("今日入网流量: %s MB\n", resp.Eflowin)
			fmt.Printf("今日出网流量: %s MB\n", resp.Eflowout)
			fmt.Printf("当前在线客户端: %d\n", resp.Oclient)
			fmt.Printf("总启动次数: %s\n", resp.Totalrun)
			fmt.Printf("今日启动次数: %s\n", resp.Todayrun)
		} else {
			fmt.Printf("✗ 获取服务统计失败\n")
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "获取下载列表",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewNodeAPIClient()
		resp, err := client.GetDownloadList()
		if err != nil {
			fmt.Printf("获取下载列表失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("========== 下载源 ==========\n")
			for _, source := range resp.Sources {
				fmt.Printf("%s: %s\n", source.Name, source.URL)
			}
			fmt.Printf("\n========== 文件列表 ==========\n")
			fmt.Println("frpc:")
			for _, item := range resp.Lists.Frpc {
				fmt.Printf("%s (%s) - %s\n", item.Name, item.Arch, item.Version)
			}
			fmt.Println("\nfrps:")
			for _, item := range resp.Lists.Frps {
				fmt.Printf("%s (%s) - %s\n", item.Name, item.Arch, item.Version)
			}
		} else {
			fmt.Printf("✗ 获取下载列表失败: %s\n", resp.Message)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "获取版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewNodeAPIClient()
		resp, err := client.GetVersion()
		if err != nil {
			fmt.Printf("获取版本信息失败: %v\n", err)
			return
		}

		fmt.Printf("========== 版本信息 ==========\n")
		fmt.Printf("HayFrps版本: %s\n", resp.VerHayfrps)
		fmt.Printf("Frpc版本: %s\n", resp.VerFrpc)
		fmt.Printf("启动器版本: %s\n", resp.VerLauncher)
		fmt.Printf("控制台版本: %s\n", resp.VerConsole)
		fmt.Printf("Dashboard版本: %s\n", resp.VerDashboard)
		fmt.Printf("启动器下载地址: %s\n", resp.UrlLauncher)
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.AddCommand(nodeInfoCmd)
	nodeCmd.AddCommand(nodeListCmd)
	nodeCmd.AddCommand(noticeCmd)
	nodeCmd.AddCommand(hayfrpInfoCmd)
	nodeCmd.AddCommand(downloadCmd)
	nodeCmd.AddCommand(versionCmd)
}
