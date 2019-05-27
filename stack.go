package main

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var r *redis.Client

func stackInit() {
	log.Println("Start redis")

	r = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := r.Ping().Result()
	if err != nil {
		log.Fatal("Error on starting redis", err)
	}
}

func stackPush(key string, cap int64, value []byte) {
	if err := r.LPush(key, value).Err(); err != nil {
		log.Println("Error push to stack", err)
		return
	}

	stackLen, err := r.LLen(key).Result()
	if err != nil {
		log.Println("Error on getting stack length", err)
		return
	}

	if stackLen > cap {
		r.RPop(key)
	}
}

func stackPutKeyValues(key string, value interface{}) {
	err := r.Set(key, value, 0).Err()
	if err != nil {
		log.Println("Error on putting key values to redis", err)
		return
	}
}
