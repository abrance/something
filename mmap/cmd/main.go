package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
)

func main() {
	// 打开文件
	file, err := os.OpenFile("example.dat", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 设置文件大小
	size := 1 << 12 // 例如，设置文件大小为 4096 字节
	err = file.Truncate(int64(size))
	if err != nil {
		log.Fatalf("Failed to truncate file: %v", err)
	}

	// 将文件映射到内存
	mmap, err := syscall.Mmap(int(file.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("Failed to mmap: %v", err)
	}
	defer syscall.Munmap(mmap)

	// 修改内存映射文件的内容
	copy(mmap, "hello")

	// 等待一段时间
	time.Sleep(1 * time.Second)

	fmt.Print(string(mmap))
}
