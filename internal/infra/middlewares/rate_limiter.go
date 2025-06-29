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

	clientIsBlocked := rl.client.Exists(rl.context, fmt.Sprintf("%v:blocked", key))

	if clientIsBlocked != nil {
		return false, fmt.Errorf("client is blocked")
	}

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
	log.Printf("Current request %v", int(currentRequestCount))
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

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())

		}
		if !allowedRequest {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			rateLimiter.client.Set(rateLimiter.context, fmt.Sprintf("%v:blocked", clientIp), "blocked", rateLimiter.blockDuration)
			log.Println("Too many requests")
			return
		}

		next.ServeHTTP(w, r)
	})
}
