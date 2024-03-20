package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func main() {
	// 定义Unix套接字的地址
	socketPath := filepath.Join(os.TempDir(), "example.sock")

	// 确保之前的套接字文件被移除
	os.Remove(socketPath)

	// 监听Unix域套接字
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("监听错误:", err)
		return
	}
	defer l.Close()

	fmt.Println("监听Unix域套接字:", socketPath)

	// 接收连接
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("接受连接错误:", err)
			continue
		}

		// 启动一个goroutine来处理连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取数据（事件）
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("读取错误:", err)
		return
	}
	response := []byte("已收到你的消息")

	// 写回响应
	wn, err := conn.Write(response)
	if err != nil {
		fmt.Println("写入错误:", err)
		return
	}
	if wn != len(response) {
		fmt.Println("未能完全写入所有响应数据")
		return
	}
	// 处理事件（这里只是简单打印）
	fmt.Println("接收到事件:", string(buffer[:n]))
}
