package main

import (
	"fmt"
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

	if _, err := redisClient.Ping().Result(); err != nil {
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
	if err := redisClient.Set(key, value, 0).Err(); err != nil {
		log.Println("Error on putting key values to redis", err)
	}
}

func incrementKey(key string) {
	if err := redisClient.Incr(key).Err(); err != nil {
		log.Println("Error on incrementing key value", err)
	}
}

func getKeys(key string) []string {
	res, err := redisClient.Keys(key).Result()
	if err == redis.Nil {
		return []string{}
	} else if err != nil {
		log.Println("Error on getting keys from set", err)
		return []string{}
	}

	return res
}

func keyDel(key string) {
	if err := redisClient.Del(key).Err(); err != nil {
		log.Println("Error on deleting key", err)
	}
}

func getKey(key string) string {
	res, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		return ""
	}

	return res
}

func setAdd(key string, value string) {
	if err := redisClient.SAdd(key, value).Err(); err != nil {
		log.Println("Error on adding key to set", err)
	}
}

func setCount(key string) int64 {
	fmt.Println(key)
	res, err := redisClient.SCard(key).Result()
	if err == redis.Nil {
		return 0
	} else if err != nil {
		log.Println("Error on count set", err)
		return 0
	}

	return res
}
