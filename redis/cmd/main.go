package main

import (
	"log"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 添加分数

	client.ZAdd("leaderboard", redis.Z{
		Score:  100,
		Member: "user123",
	})

	// 获取排行榜
	result, err := client.ZRevRangeWithScores("leaderboard", 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range result {
		log.Printf("User: %v, Score: %v\n", item.Member, item.Score)
	}
}
