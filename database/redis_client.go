package database

import (
	"context"
	"fmt"
	"log"
	"github.com/go-redis/redis/v8"
)

//connecting to RedisClient
func RedisClient()*redis.Client{
    redisHost := "hv-par6-004.clvrcld.net"
    redisPort := "14138"
    redisPassword := "tiKftxJztVt9G6kYeJ3"

    // Construct the Redis client options
    rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
        Password: redisPassword,
        DB:       0, // use default DB
    })

    // Test the connection
    ctx := context.Background()
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Printf("Failed to connect to Redis: %v", err)
    }
	return rdb
}