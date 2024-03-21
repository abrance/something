package gen_host

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type kv struct {
	Key   string
	Value string
}

func GenerateHostsFile(entries map[string]string, hostFile string) error {
	// 打开文件
	file, err := os.OpenFile(hostFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

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
