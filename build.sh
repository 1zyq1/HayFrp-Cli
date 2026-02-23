#!/bin/bash

# HayFrp-Cli 一键构建脚本

# 创建发布目录
mkdir -p releases

echo "开始构建 HayFrp-Cli..."

# Windows AMD64
echo "构建 Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -o releases/HayFrp-Cli-windows-amd64.exe main.go

# Windows 386
echo "构建 Windows 386..."
GOOS=windows GOARCH=386 go build -o releases/HayFrp-Cli-windows-386.exe main.go

# Linux AMD64
echo "构建 Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o releases/HayFrp-Cli-linux-amd64 main.go

# Linux 386
echo "构建 Linux 386..."
GOOS=linux GOARCH=386 go build -o releases/HayFrp-Cli-linux-386 main.go

# Linux ARM64
echo "构建 Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -o releases/HayFrp-Cli-linux-arm64 main.go

# Darwin (macOS) AMD64
echo "构建 macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -o releases/HayFrp-Cli-darwin-amd64 main.go

# Darwin (macOS) ARM64 (Apple Silicon)
echo "构建 macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -o releases/HayFrp-Cli-darwin-arm64 main.go

echo "构建完成！"
echo "所有构建文件已保存到 releases 目录"

# 列出构建的文件
ls -lh releases/
