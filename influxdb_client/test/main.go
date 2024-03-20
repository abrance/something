package main

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"math/rand"
	"sync"
	"time"
)

var (
	Payload = ""
	wg      = sync.WaitGroup{}
	// 总量 10 GB
	totalSize = 1024 * 1024 * 1024 * 10
	//每条 信息大小
	messageSize = 1 * 10
	// routine 数量
	senderCount = 8
	// loop count
	loopCount int = totalSize / messageSize / senderCount
)

// example 1, 10亿条数据, 每条数据 100KB, 一共 10GB 1000000000 * 10B = 10GB

// 设计 tag 数据 , field 数据各占一般, 同样是 每条数据大小为 100KB, 插入 10000 条, 记录消耗时间.
// tag 50KB  1K 个 50B 数据
// field 50KB   1 个 50KB  数据
// 时间  64bit 忽略

func Insert10GB() {
	t := time.Now()
	for i := 0; i < senderCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Client1GBData()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(t))
}

func Client1GBData() {
	t0 := time.Now()
	// 连接到 InfluxDB
	client := influxdb2.NewClient("http://127.0.0.1:8086", "UwrvLbNFgDS-9cSf-ghepFuQB91iLoABQlbqqkQELHc2CGcG7jvfhxrCX_HI8K4Ni5olMVtuYRSLgsCcDxPFVw==")
	writeAPI := client.WriteAPI("ouryun", "default")
	for i := 0; i < loopCount; i++ {
		Insert1MBData(writeAPI)
	}

	client.Close()
	fmt.Printf("time: %v", time.Since(t0))
}

// TestApi  测试插入一条数据
// user password xiaoy 110isfaker.
// organization ouryun
// bucket default
func TestApi() {
	client := influxdb2.NewClient("http://127.0.0.1:8086", "UwrvLbNFgDS-9cSf-ghepFuQB91iLoABQlbqqkQELHc2CGcG7jvfhxrCX_HI8K4Ni5olMVtuYRSLgsCcDxPFVw==")
	writeAPI := client.WriteAPI("ouryun", "default")
	Insert1MBData(writeAPI)
}

func Insert1MBData(api api.WriteAPI) {
	tags := make(map[string]string)
	//faker.SetRandomStringLength(25)
	//faker.SetRandomMapAndSliceMaxSize(1000)
	//faker.SetRandomMapAndSliceMinSize(990)
	//faker.FakeData(&tags)

	fields := make(map[string]interface{})
	//for i := 0; i < 5; i++ {
	//	fieldKey := fmt.Sprintf("field%d", i)
	//	fieldValue := generateLargeData()
	//	fields[fieldKey] = fieldValue
	//}
	err := faker.SetRandomStringLength(messageSize)
	err = faker.SetRandomNumberBoundaries(0, 9)

	if err != nil {
		fmt.Printf("err: %v", err)
	}
	err = faker.FakeData(&Payload)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fields["z"] = Payload
	rand.Seed(time.Now().UnixNano())
	fields["data"] = rand.Intn(10)
	point := influxdb2.NewPoint("measurement", tags, fields, time.Now())
	// 写入数据点
	api.WritePoint(point)
	api.Flush()
}

func main() {
	//fmt.Printf(fmt.Sprintf("%s%05d", "n: ", 9))
	//Client10GBData()
	//faker.SetRandomNumberBoundaries()
	//Insert10GB()
	TestApi()
}
