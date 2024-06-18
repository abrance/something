package conn

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"time"
	"unsafe"
)

func triggerEvent(efd int) {
	// 2秒后向 eventfd 写入数据
	time.Sleep(2 * time.Second)
	value := uint64(1)
	_, err := unix.Write(efd, (*(*[8]byte)(unsafe.Pointer(&value)))[:])
	if err != nil {
		fmt.Println("Error writing to eventfd:", err)
	}
	fmt.Println("Event triggered")
}

func EpollCase() {
	var err error
	var epfd int
	var efd int

	if efd, err = unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC); err != nil {
		_ = unix.Close(efd)
		_ = os.NewSyscallError("eventfd", err)
	}
	defer unix.Close(efd)

	if epfd, err = unix.EpollCreate(unix.EPOLL_CLOEXEC); err != nil {
		_ = unix.Close(epfd)
		_ = os.NewSyscallError("epoll_create", err)
	}
	defer unix.Close(epfd)

	const ReadEvents = unix.EPOLLIN | unix.EPOLLPRI
	var ev uint32 = ReadEvents
	err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, efd, &unix.EpollEvent{
		Events: ev,
		Fd:     int32(efd),
	})
	if err != nil {
		_ = os.NewSyscallError("epoll_ctl", err)
	}

	go triggerEvent(efd)

	// 设置超时并监听事件
	events := make([]unix.EpollEvent, 1)
	fmt.Println("Waiting for events...")
	for {
		nevents, err := unix.EpollWait(epfd, events, 5000)
		if err != nil {
			fmt.Println(os.NewSyscallError("epoll_wait", err))
			return
		}
		if nevents == 0 {
			fmt.Println("No events received before timeout")
			return
		}
		for _, e := range events[:nevents] {
			if e.Events&unix.EPOLLIN != 0 {
				fmt.Println("Event received:", e)
				// 读取 eventfd 的值，确认事件
				var buf [8]byte
				_, err := unix.Read(efd, buf[:])
				if err != nil {
					fmt.Println("Failed to read from eventfd:", err)
				} else {
					val := *(*uint64)(unsafe.Pointer(&buf[0]))
					fmt.Println("Value from eventfd:", val)
				}
				return
			}
		}
	}
}
