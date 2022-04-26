package gamestrg

import (
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	redisClient *redis.Client
}
