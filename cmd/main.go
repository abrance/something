package main

import (
	"github.com/abrance/something/pkg/gen_host"
	"log"
)

func main() {
	caseGenHost()
}

func caseGenHost() {
	entries := map[string]string{
		"127.0.0.1":   "localhost",
		"192.168.1.1": "example.com",
		"12.3.1.23":   "HOST1",
		"123.12.3.12": "HOST2",
		// 添加更多的 IP-Hostname 映射
	}
	err := gen_host.GenerateHostsFile(entries, "./hosts")
	if err != nil {
		log.Fatal(err)
	}
}
