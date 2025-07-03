package middlewares

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	limit         int
	window        time.Duration
	blockDuration time.Duration
	context       context.Context
	client        redis.Client
}

func (rl *RateLimiter) Allow(key string) (bool, error) {

	exists, err := rl.client.Exists(rl.context, fmt.Sprintf("%v:blocked", key)).Result()
	if err != nil {
		return false, fmt.Errorf("redis error: %w", err)
	}

	if exists > 0 {
		return false, fmt.Errorf("client is blocked")
	}

	pipe := rl.client.TxPipeline()

	incr := pipe.Incr(rl.context, key)

	pipe.Expire(rl.context, key, rl.window)

	_, err = pipe.Exec(rl.context)
	currentRequestCount := incr.Val()

	if err != nil {
		return false, err
	}

	return int(currentRequestCount) <= rl.limit, nil
}

func NewRateLimiter(limit int, window time.Duration, blockDuration time.Duration, context context.Context, client redis.Client) *RateLimiter {

	return &RateLimiter{
		limit:         limit,
		window:        window,
		context:       context,
		client:        client,
		blockDuration: blockDuration,
	}
}

func RateLimiterMiddleware(next http.Handler, rateLimiter RateLimiter) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clientIp, _, _ := net.SplitHostPort(r.RemoteAddr)
		allowedRequest, err := rateLimiter.Allow(clientIp)

		log.Printf("Client IP: %v", clientIp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())

		}
		if !allowedRequest {
			rateLimiter.client.Set(rateLimiter.context, fmt.Sprintf("%v:blocked", clientIp), "blocked", rateLimiter.blockDuration)
			log.Println("Too many requests")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
