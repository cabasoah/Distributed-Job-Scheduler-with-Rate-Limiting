package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/cabasoah/Distributed-Job-Scheduler-with-Rate-Limiting/ratelimiter"
)

func main() {
	port := flag.String("port", "8083", "port for the rate limiter API to listen on")
	flag.Parse()

	// limiter: capacity 10,1 token/sec
	limiter := ratelimiter.NewTokenBucketLimiter(10, 1)

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
