package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func main() {
	// 服务器的Unix套接字地址
	socketPath := filepath.Join(os.TempDir(), "example.sock")

	// 连接到Unix域套接字
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("连接错误:", err)
		return
	}
	defer conn.Close()

	// 发送消息
	_, err = conn.Write([]byte("Hello from client"))
	if err != nil {
		fmt.Println("发送消息错误:", err)
		return
	}

	// 读取响应
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("读取响应错误:", err)
		return
	}
	fmt.Println("消息已接收:", string(buffer))
}
