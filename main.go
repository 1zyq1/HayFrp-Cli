package main

import (
	"os"
	"hayfrp-cli/cmd"
)

func main() {
	// 如果没有命令行参数，自动执行start命令
	if len(os.Args) <= 1 {
		cmd.ExecuteStart()
	} else {
		cmd.Execute()
	}
}
