package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"syscall"
)

const DefaultEvents = 10
const NeverTimeout = -1

var (
	EpollFD int
)

func main() {
	var err error

	EpollFD, err = syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer syscall.Close(EpollFD)

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
	// 将文件描述符添加到 epoll 监听列表
	// 示例：syscall.EpollCtl(EpollFD, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Events: syscall.EPOLLIN, Fd: int32(fd)})
	// 事件循环
	events := make([]syscall.EpollEvent, DefaultEvents) // 调整大小以适应预期的负载
	for {
		n, err := syscall.EpollWait(EpollFD, events, NeverTimeout)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < n; i++ {
			fd := events[i].Fd
			// 处理文件描述符上的事件，例如读取数据、接受新连接等
			go handleEvent(fd)
		}
	}

	conn.Close()
}

func handleEvent(fd int32) {
	// 处理与该文件描述符相关的事件
	// 例如，如果它是一个网络套接字，你可能需要读取数据或接受连接
}
