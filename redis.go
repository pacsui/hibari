package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	dbNum, err := strconv.ParseInt(os.Getenv("REDIS_DB_NUM"), 10, 64)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_TOKEN"),
		DB:       int(dbNum),
	})
	if rdb == nil {
		return nil, fmt.Errorf("unable to create redis db client")
	}
	return rdb, nil
}
