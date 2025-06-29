package middlewares

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	limit   int
	window  time.Duration
	context context.Context
	client  redis.Client
}

func (rl *RateLimiter) Allow(key string) (bool, error) {

	pipe := rl.client.TxPipeline()

	incr := pipe.Incr(rl.context, key)

	_, err := pipe.Exec(rl.context)

	currentRequestCount := incr.Val()

	if currentRequestCount == 1 {
		pipe.Expire(rl.context, key, rl.window)
	}
	if err != nil {
		return false, err
	}

	return int(currentRequestCount) <= rl.limit, nil
}

func NewRateLimiter(limit int, window time.Duration, context context.Context, client redis.Client) *RateLimiter {

	return &RateLimiter{
		limit:   limit,
		window:  window,
		context: context,
		client:  client,
	}
}

func RateLimiterMiddleware(next http.Handler, rateLimiter RateLimiter) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clientIp, _, _ := net.SplitHostPort(r.RemoteAddr)
		allowedRequest, err := rateLimiter.Allow(clientIp)

		if err != nil {
			log.Println("Rate limiter error")
		}
		if !allowedRequest {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			log.Println("Too many requests")
			return
		}

		next.ServeHTTP(w, r)
	})
}
