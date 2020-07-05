package database

import "github.com/go-redis/redis"

var option = &redis.Options{
	Addr: "localhost:6379",
	DB:   0,
}

func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(option)

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}
