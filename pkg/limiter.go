package pkg

import (
	"context"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
)

type RateLimiter struct {
	handler     http.Handler
	cacheRules  CacheRules
	ctx         context.Context
	redisClient redis.Client
}

func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var key string
	if rl.cacheRules.kind == "host" {
		key = r.Host
	} else {
		http.Error(w, "Unsupported limit kind", http.StatusBadRequest)
	}

	count, err := rl.cacheRules.GetHostTokenCount(&rl.redisClient, rl.ctx, key)
	if err != nil {
		if err != redis.Nil {
			http.Error(w, "There was an issue communicating with redis.", http.StatusFailedDependency)
			return
		}
	}
	countInt, _ := strconv.Atoi(count)

	w.Header().Add("X-Rate-Count", count)

	if countInt > rl.cacheRules.number {
		http.Error(w, "Rate limit exceeded.", http.StatusTooManyRequests)
		return
	}

	incrErr := rl.cacheRules.IncrementHostTokens(&rl.redisClient, rl.ctx, key)
	if incrErr != nil {
		http.Error(w, "There was an issue with an external dependency.", http.StatusFailedDependency)
		return
	}

	rl.handler.ServeHTTP(w, r)
}

func NewWrappedRateLimiter(handler http.Handler, cacheRules CacheRules, ctx context.Context, redisClient redis.Client) *RateLimiter {
	return &RateLimiter{
		handler,
		cacheRules,
		ctx,
		redisClient,
	}
}
