package main

import (
	"net/http"
	"os"
	"rate-limiter/internal"
	"rate-limiter/pkg"
)

func main() {
	addr := os.Getenv("ADDR")
	cacheRules := pkg.BuildCacheRules()
	redisClient, ctx := pkg.NewRedisClient()

	mux := http.NewServeMux()
	mux.HandleFunc("/time", internal.CurrentTimeHandler)

	wrappedMux := pkg.NewLogger(pkg.NewWrappedRateLimiter(mux, *cacheRules, ctx, *redisClient))

	err := http.ListenAndServe(addr, wrappedMux)
	if err != nil {
		return
	}
}
