package main

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	Payload = ""
	wg      = sync.WaitGroup{}
	// 总量 10 GB
	totalSize = 1024 * 1024 * 1024 * 10
	//每条 信息大小
	messageSize = 1024 * 1024
	// routine 数量
	senderCount = 8
	// loop count
	loopCount int = totalSize / messageSize / senderCount
)

type Message struct {
	Field1 string
}

func Insert10GB() {
	t := time.Now()
	for i := 0; i < senderCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ClientInsert1G()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(t))
}

func ClientInsert1G() {
	t0 := time.Now()
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017" +
		"/?&replicaSet=rs0&connect=direct")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < loopCount; i++ {
		Insert1Data(client)
	}

	fmt.Printf("time: %v\n", time.Since(t0))
}

func Insert1Data(client *mongo.Client) {
	var err error
	var insertDataLineNum int = 10

	var line int = messageSize / insertDataLineNum

	collection := client.Database("test").Collection("coll")
	err = faker.SetRandomStringLength(line)
	if err != nil {
		fmt.Println(err)
	}
	err = faker.FakeData(&Payload)
	if err != nil {
		fmt.Println(err)
	}
	var data []interface{}
	for i := 0; i < insertDataLineNum; i++ {
		data = append(data, Message{Field1: Payload})
	}

	_, err = collection.InsertMany(context.TODO(), data)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	Insert10GB()
}
