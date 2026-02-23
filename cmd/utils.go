package cmd

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// readPasswordWithMask 读取密码，输入时不显示
func readPasswordWithMask(prompt string) (string, error) {
	fmt.Print(prompt)
	fmt.Print(" (输入不会显示): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // 换行
	return string(bytePassword), nil
}
