package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cabasoah/Distributed-Job-Scheduler-with-Rate-Limiting/ratelimiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	port := flag.String("port", "8083", "port for the rate limiter API to listen on")
	strategy := flag.String("strategy", "token-bucket", "token-bucket | fixed-window")
	redisAddr := flag.String("redis-addr", "localhost:6379", "redis address")
	flag.Parse()

	// limiter: capacity 10,1 token/sec
	// switch statement to switch between bucket and window limiter
	var limiter ratelimiter.Limiter

	switch *strategy {
	case "token-bucket":
		limiter = ratelimiter.NewTokenBucketLimiter(10,1)
	case "fixed-window":
		limiter = ratelimiter.NewFixedWindowLimiter(60*time.Second, 60) //60 req/min
	case "sliding-window-counter":
		limiter = ratelimiter.NewSlidingWindowCounterLimiter(60*time.Second, 60)
	case "redis-sliding-window":
		rdb := redis.NewClient(&redis.Options{Addr: *redisAddr})
		limiter = ratelimiter.NewRedisSlidingWindowLimiter(rdb, 60*time.Second, 60)
	default:
		log.Fatalf("Unkown Strategy: %s", *strategy)

	}

	http.HandleFunc("/unlimited", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Unlimited! Let's go!"))
	})

	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if !limiter.Allow(ip) {
			http.Error(w, "Too many request", http.StatusTooManyRequests)
			return
		}
		w.Write([]byte("Limited! Don't over use me!"))
	})

	log.Printf("Rate limiter API listening on %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}
