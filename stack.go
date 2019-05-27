package main

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func stackInit() {
	log.Println("Start redis")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Error on starting redis", err)
	}
}

func stackPush(key string, cap int64, value []byte) {
	if err := redisClient.LPush(key, value).Err(); err != nil {
		log.Println("Error push to stack", err)
		return
	}

	stackLen, err := redisClient.LLen(key).Result()
	if err != nil {
		log.Println("Error on getting stack length", err)
		return
	}

	if stackLen > cap {
		redisClient.RPop(key)
	}
}

func stackPutKeyValues(key string, value interface{}) {
	err := redisClient.Set(key, value, 0).Err()
	if err != nil {
		log.Println("Error on putting key values to redis", err)
		return
	}
}
