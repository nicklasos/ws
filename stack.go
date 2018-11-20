package main

import (
	"github.com/go-redis/redis"
	"log"
	"os"
)

var client *redis.Client

const key = "ws.stack"

func StackInit() {
	log.Println("Start redis")

	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("Error on starting redis", err)
	}
}

func StackPush(value []byte) {
	if err := client.LPush(key, value).Err(); err != nil {
		log.Println("Error push to stack", err)
		return
	}

	stackLen, err := client.LLen(key).Result()
	if err != nil {
		log.Println("Error on getting stack length", err)
		return
	}

	if stackLen > 10 {
		client.RPop(key)
	}
}

func StackPutKeyValues(key string, value interface{}) {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		log.Println("Error on putting key values to redis", err)
		return
	}
}
