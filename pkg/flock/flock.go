package flock

import (
	"fmt"
	"github.com/gofrs/flock"
	"log"
	"os"
	"time"
)

func Case() {
	f := flock.New(os.TempDir() + "/go-lock.lock")
	lock, err := f.TryLock()
	go func() {
		log.Println("2:---")
		f2 := flock.New(os.TempDir() + "/go-lock.lock")
		time.Sleep(1 * time.Second)
		l := f2.Locked()
		log.Printf("2:--- : %v", l)
		// todo 试试 未获取到锁的情况下进行写入 case
	}()
	if err != nil {
		fmt.Printf("lock error: %v\n", err)
		return
	}
	file, err := os.OpenFile(os.TempDir()+"/go-lock.lock", os.O_RDWR, 0600)
	if err != nil {
		fmt.Printf("open file error: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString("Hello, world!")
	if err != nil {
		// 写入错误处理
	}

	log.Printf("%v\n", lock)
	time.Sleep(20 * time.Second)
	log.Printf("locked: %v\n", f.Locked())
	f.Unlock()
	log.Printf("locked: %v\n", f.Locked())
}
