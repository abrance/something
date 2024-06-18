package main

import (
	"bytes"
	"errors"
	"github.com/abrance/something/pkg/ginroutev"
	"io"
	"log"
	"sync"
)

func NodeIsSEM(b uint32) bool {
	return b>>7&1 == 1
}

// MessageBuffer 实现了 io.ReadWriteCloser 接口
type MessageBuffer struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
	lock        sync.Mutex
	Wait        chan struct{}
	CloseCh     chan struct{}
	readCh      chan struct{}
	closed      bool
}

// NewMessageBuffer 创建一个新的 MessageBuffer 实例
func NewMessageBuffer() *MessageBuffer {
	return &MessageBuffer{
		readBuffer:  bytes.NewBuffer(nil),
		writeBuffer: bytes.NewBuffer(nil),
		Wait:        make(chan struct{}),
		CloseCh:     make(chan struct{}),
		readCh:      make(chan struct{}, 1), // 缓冲大小为1的通道，用于通知读取操作
	}
}

func (mb *MessageBuffer) WriteToReadBuffer(data []byte) (int, error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	n, err := mb.readBuffer.Write(data)
	if err == nil {
		select {
		case mb.readCh <- struct{}{}:
		default:
		}
	}
	return n, err
}

// Read 从读取缓冲区中读取数据
func (mb *MessageBuffer) Read(p []byte) (n int, err error) {
	for {
		mb.lock.Lock()
		if mb.readBuffer.Len() > 0 {
			n, err = mb.readBuffer.Read(p)
			mb.lock.Unlock()
			return n, err
		}
		if mb.closed {
			mb.lock.Unlock()
			return 0, errors.New("buffer closed")
		}
		mb.lock.Unlock()

		select {
		case <-mb.readCh:
		case <-mb.CloseCh:
			return 0, errors.New("buffer closed")
		}
	}
}

// ReadNotBlock 非阻塞读取
func (mb *MessageBuffer) ReadNotBlock(b []byte) (n int, err error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	return mb.readBuffer.Read(b)
}

// Write 向写入缓冲区中写入数据
func (mb *MessageBuffer) Write(p []byte) (n int, err error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	return mb.writeBuffer.Write(p)
}

// Close 关闭缓冲区并解除阻塞
func (mb *MessageBuffer) Close() error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	if mb.closed {
		return errors.New("buffer already closed")
	}
	mb.closed = true
	close(mb.CloseCh)
	select {
	case mb.readCh <- struct{}{}:
	default:
	}
	log.Println("close 999999999999999")
	return nil
}

func (mb *MessageBuffer) ReadWBAll() ([]byte, error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	return io.ReadAll(mb.writeBuffer)
}

// Flush 主动发送数据
func (mb *MessageBuffer) Flush() {
	mb.Wait <- struct{}{}
}

func main() {
	ginroutev.Server()
	//// 示例用法
	//mb := NewMessageBuffer()
	//
	//// 启动一个写入协程
	//go func() {
	//	data := []byte("Hello, World!")
	//	mb.WriteToReadBuffer(data)
	//}()
	//
	//// 启动一个读取协程
	//go func() {
	//	buf := make([]byte, 1024)
	//	n, err := mb.Read(buf)
	//	if err != nil {
	//		log.Println("Read error:", err)
	//		return
	//	}
	//	log.Println("Read data:", string(buf[:n]))
	//}()
	//
	//// 模拟一些延迟，然后关闭缓冲区
	//// 等待 2 秒，以确保读取协程有机会运行
	//time.Sleep(2 * time.Second)
	//mb.Close()
	//log.Println(mb.closed)
}

//func main() {
//	// client_go.Connect()
//	//caseGenHost()
//	//fn, err := gz.ReadFileInGzip("./test.tar.gz", "test/test.txt")
//	//if err != nil {
//	//	fmt.Println(err)
//	//}
//	//fmt.Println(string(fn))
//	//conn.EpollCase()
//	//flock.Case()
//	fmt.Println(NodeIsSEM(8))
//}
//
//func caseGenHost() {
//	nh1 := gen_host.NodeHost{
//		IpLs: []string{"192.168.1.2", "192.168.1.1"},
//		TyId: gen_host.HostnameTypeM,
//	}
//	nh2 := gen_host.NodeHost{
//		IpLs: []string{"172.16.0.1"},
//		TyId: gen_host.HostnameSeByNodeId,
//	}
//	nh3 := gen_host.NodeHost{
//		IpLs: []string{"172.16.0.10"},
//		TyId: gen_host.HostnameSeByNodeId,
//	}
//	scHost := gen_host.NodeHost{
//		IpLs: []string{"172.24.0.1"},
//		TyId: gen_host.HostnameSC,
//	}
//	nh2.SetSeNodeId("$nodeid2")
//	nh3.SetSeNodeId("$nodeid3")
//
//	entries := gen_host.GenerateEntries([]gen_host.NodeHost{nh1, nh2, nh3, scHost})
//	fmt.Println(entries)
//	err := gen_host.GenerateHostsFile(entries, "./hosts")
//	if err != nil {
//		log.Fatal(err)
//	}
//}
