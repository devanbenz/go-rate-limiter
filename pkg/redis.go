package pkg

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type CacheRules struct {
	kind   string
	window int
	number int
}

func NewRedisClient() (*redis.Client, context.Context) {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	}), context.Background()
}

func BuildCacheRules() *CacheRules {
	return &CacheRules{
		kind:   "host",
		window: 60,
		number: 5,
	}
}

func (c *CacheRules) IncrementHostTokens(client *redis.Client, ctx context.Context, key string) error {
	incrErr := client.Incr(ctx, key).Err()
	if incrErr != nil {
		return incrErr
	}
	expErr := client.Expire(ctx, key, time.Duration(c.window)*time.Second).Err()
	if expErr != nil {
		return expErr
	}
	return nil
}

func (c *CacheRules) GetHostTokenCount(client *redis.Client, ctx context.Context, key string) (string, error) {
	tokens, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return tokens, nil
}
