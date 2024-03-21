package gen_host

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type kv struct {
	Key   string
	Value string
}

// GenerateHostsFile 生成 hosts 文件
// 生成后的 文件如下，已做了 hostname 排序
// 127.0.0.1	localhost
// ::1	ip6-localhost ip6-loopback
// fe00::0	ip6-localnet
// ff00::0	ip6-mcastprefix
// ff02::1	ip6-allnodes
// ff02::2	ip6-allrouters
//
// 172.16.0.1	$nodeid2.srhino.svc.local
// 172.16.0.10	$nodeid3.srhino.svc.local
// 172.24.0.1	sc.srhino.svc.local
// 192.168.1.2	sem.srhino.svc.local
// 192.168.1.1	sem.srhino.svc.local
func GenerateHostsFile(entries map[string]string, hostFile string) error {
	// 打开文件
	file, err := os.OpenFile(hostFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	var ih = [6][2]string{
		{"127.0.0.1", "localhost"},
		{"::1", "ip6-localhost ip6-loopback"},
		{"fe00::0", "ip6-localnet"},
		{"ff00::0", "ip6-mcastprefix"},
		{"ff02::1", "ip6-allnodes"},
		{"ff02::2", "ip6-allrouters"},
	}

	for i := 0; i < len(ih); i++ {
		_, err := writer.WriteString(fmt.Sprintf("%s\t%s\n", ih[i][0], ih[i][1]))
		if err != nil {
			return err
		}
	}
	_, err = writer.WriteString("\n\n")
	if err != nil {
		return err
	}
	// 排序并写入
	kvs := make([]kv, 0, len(entries))
	for k, v := range entries {
		kvs = append(kvs, kv{k, v})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value < kvs[j].Value
	})
	for _, kv := range kvs {
		fmt.Printf("%s: %s\n", kv.Key, kv.Value)
		_, err := writer.WriteString(fmt.Sprintf("%s\t%s\n", kv.Key, kv.Value))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

const (
	HostnameTypeM      = "sem.srhino.svc.local"
	HostnameSeByNodeId = "srhino.svc.local"
	HostnameSC         = "sc.srhino.svc.local"
)

type NodeHost struct {
	TyId     string
	IpLs     []string
	seNodeId string // 引擎才有
}

func (s *NodeHost) SetIps(ipLs []string) {
	s.IpLs = ipLs
}

func (s *NodeHost) SetSeNodeId(seNodeId string) {
	s.seNodeId = seNodeId
}

func (s *NodeHost) SetType(tyId string) error {
	switch tyId {
	case HostnameTypeM, HostnameSeByNodeId, HostnameSC:
		s.TyId = tyId
		return nil
	default:
		return fmt.Errorf("unsupported type Id: %s", tyId)
	}
}

func (s *NodeHost) InsertIntoEntries(entries map[string]string) {
	switch s.TyId {
	case HostnameSC, HostnameTypeM:
		for _, ip := range s.IpLs {
			entries[ip] = s.TyId
		}
	case HostnameSeByNodeId:
		for _, ip := range s.IpLs {
			entries[ip] = fmt.Sprintf("%s.%s", s.seNodeId, s.TyId)
		}
	default:
		log.Fatal("unsupported type Id: %s", s.TyId)
		return
	}
}

// GenerateEntries
// hostnameMIpLs : 管理面 ip 列表
// scHostIp : 总控 ip 列表 ，考虑后面做总控修改 ip
func GenerateEntries(nodeHostLs []NodeHost) map[string]string {
	entries := make(map[string]string)
	for _, node := range nodeHostLs {
		node.InsertIntoEntries(entries)
	}
	return entries
}
