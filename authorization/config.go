package authorization

import (
	"github.com/go-redis/redis/v8"
)

// type Config struct {
// 	// Redis address, e.g. localhost:6379
// 	RedisAddr string
// 	// redis password to auth
// 	RedisPassword string
// 	// redis database
// 	RedisDatabase string
// }

var tokenCache *redis.Client

func InitTokenCache(redisAddr, redisPasswd string, redisDb int) {
	tokenCache = redis.NewClient(&redis.Options{
		Addr:     redisAddr,   // redis address, e.g. localhost:6379
		Password: redisPasswd, // redis password
		DB:       redisDb,     // redis database, select 0
	})
}
