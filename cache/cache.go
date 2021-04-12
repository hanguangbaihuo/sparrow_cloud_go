package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var cache *redis.Client

// InitCache is to init redis cache
func InitCache(redisAddr, redisPasswd string, redisDb int) *redis.Client {
	cache = redis.NewClient(&redis.Options{
		Addr:     redisAddr,   // redis address, e.g. localhost:6379
		Password: redisPasswd, // redis password
		DB:       redisDb,     // redis database, select 0
	})
	_, err := cache.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return cache
}

func InitCustomCache(opt redis.Options) *redis.Client {
	cache = redis.NewClient(&opt)
	_, err := cache.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return cache
}

// Get for get cache client, panic when cache client is nil
func Get() *redis.Client {
	if cache != nil {
		return cache
	}
	log.Fatalln("cache is nil, please InitCache firstly")
	return nil
}

// GetOrNil return cache, nil return when not init cache
func GetOrNil() *redis.Client {
	return cache
}
