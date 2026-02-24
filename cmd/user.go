package cmd

import (
	"fmt"
	"strconv"
	"syscall"

	"hayfrp-cli/api"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "用户相关操作",
	Long:  `用户登录、注册、信息查询等操作`,
}

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "用户登录",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user := args[0]

		fmt.Print("请输入密码: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\n读取密码失败: %v\n", err)
			return
		}
		fmt.Println() // 换行
		passwd := string(bytePassword)

		client := api.NewUserAPIClient()
		resp, err := client.Login(user, passwd)
		if err != nil {
			fmt.Printf("登录失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ 登录成功！\n")
			fmt.Printf("  Token: %s\n", resp.Token)
		} else {
			fmt.Printf("✗ 登录失败: %s\n", resp.Message)
		}
	},
}

var verifyCsrfCmd = &cobra.Command{
	Use:   "verify [csrf]",
	Short: "验证Token是否有效",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]

		client := api.NewUserAPIClient()
		resp, err := client.VerifyCsrf(csrf)
		if err != nil {
			fmt.Printf("验证失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ Token有效\n")
			fmt.Printf("  Token: %s\n", resp.Token)
		} else {
			fmt.Printf("✗ Token无效: %s\n", resp.Message)
		}
	},
}

var userInfoCmd = &cobra.Command{
	Use:   "info [csrf]",
	Short: "获取用户信息",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]

		client := api.NewUserAPIClient()
		resp, err := client.GetInfo(csrf)
		if err != nil {
			fmt.Printf("获取用户信息失败: %v\n", err)
			return
		}

		if resp.Status {
			fmt.Printf("========== 用户信息 ==========\n")
			fmt.Printf("用户ID: %v\n", resp.ID)
			fmt.Printf("用户名: %s\n", resp.Username)
			fmt.Printf("邮箱: %s\n", resp.Email)
			// 转换流量为GB
			var trafficGB float64
			switch v := resp.Traffic.(type) {
			case string:
				if val, err := strconv.ParseFloat(v, 64); err == nil {
					trafficGB = val / 1024
				}
			case float64:
				trafficGB = v / 1024
			}
			fmt.Printf("剩余流量: %.2f GB\n", trafficGB)
			fmt.Printf("今日使用流量: %v Bytes\n", resp.Todaytraffic)
			fmt.Printf("拥有隧道数: %v\n", resp.Proxies)
			fmt.Printf("已使用隧道: %v\n", resp.Useproxies)

			// 处理可能为 string 或 bool 的字段
			fmt.Printf("是否实名: %v\n", resp.Realname)
			fmt.Printf("是否服务商: %v\n", resp.Sprovider)

			fmt.Printf("UUID: %s\n", resp.UUID)
			fmt.Printf("Token: %s\n", resp.Token)
			if resp.Signdate != "" && resp.Signdate != "null" {
				fmt.Printf("上次签到时间: %s\n", resp.Signdate)
				fmt.Printf("总签到天数: %v\n", resp.Totalsign)
				fmt.Printf("总签到流量: %v GB\n", resp.Totaltraffic)
			}
		} else {
			fmt.Printf("获取用户信息失败: %s\n", resp.Message)
		}
	},
}

var signCmd = &cobra.Command{
	Use:   "sign [csrf]",
	Short: "每日签到",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]

		client := api.NewUserAPIClient()
		resp, err := client.Sign(csrf)
		if err != nil {
			fmt.Printf("签到失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
			fmt.Printf("  签到获得流量: %.2f GB\n", resp.Signflow)
			fmt.Printf("  剩余流量: %.2f GB\n", resp.Flow)
		} else {
			fmt.Printf("✗ 签到失败: %s\n", resp.Message)
		}
	},
}

var retokenCmd = &cobra.Command{
	Use:   "retoken [csrf]",
	Short: "更新用户Token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrf := args[0]

		client := api.NewUserAPIClient()
		resp, err := client.ReToken(csrf)
		if err != nil {
			fmt.Printf("更新Token失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
			fmt.Printf("  新Token: %s\n", resp.Token)
		} else {
			fmt.Printf("✗ 更新Token失败: %s\n", resp.Message)
		}
	},
}

var sendRegCodeCmd = &cobra.Command{
	Use:   "send-reg [username] [email]",
	Short: "发送注册验证码",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		user := args[0]
		email := args[1]

		client := api.NewUserAPIClient()
		resp, err := client.SendRegCode(user, user+"@"+email, email)
		if err != nil {
			fmt.Printf("发送验证码失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 发送验证码失败: %s\n", resp.Message)
		}
	},
}

var registerCmd = &cobra.Command{
	Use:   "register [username] [email] [password] [code]",
	Short: "用户注册",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		user := args[0]
		email := args[1]
		passwd := args[2]
		code := args[3]

		client := api.NewUserAPIClient()
		resp, err := client.Register(user, user+"@"+email, email, passwd, code)
		if err != nil {
			fmt.Printf("注册失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 注册失败: %s\n", resp.Message)
		}
	},
}

var sendFindPassCodeCmd = &cobra.Command{
	Use:   "send-findpass [username]",
	Short: "发送重置密码验证码",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user := args[0]

		client := api.NewUserAPIClient()
		resp, err := client.SendFindPassCode(user)
		if err != nil {
			fmt.Printf("发送验证码失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 发送验证码失败: %s\n", resp.Message)
		}
	},
}

var resetPassCmd = &cobra.Command{
	Use:   "reset-pass [token] [new-password]",
	Short: "重置密码",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		newpass := args[1]

		client := api.NewUserAPIClient()
		resp, err := client.ResetPass(token, newpass)
		if err != nil {
			fmt.Printf("重置密码失败: %v\n", err)
			return
		}

		if resp.Status == 200 {
			fmt.Printf("✓ %s\n", resp.Message)
		} else {
			fmt.Printf("✗ 重置密码失败: %s\n", resp.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(loginCmd)
	userCmd.AddCommand(verifyCsrfCmd)
	userCmd.AddCommand(userInfoCmd)
	userCmd.AddCommand(signCmd)
	userCmd.AddCommand(retokenCmd)
	userCmd.AddCommand(sendRegCodeCmd)
	userCmd.AddCommand(registerCmd)
	userCmd.AddCommand(sendFindPassCodeCmd)
	userCmd.AddCommand(resetPassCmd)
}
