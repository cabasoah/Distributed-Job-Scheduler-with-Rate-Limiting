package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/cabasoah/Distributed-Job-Scheduler-with-Rate-Limiting/loadbalancer"
)

//Load balancer routing to multple backend servers
func main()  {
	port := flag.String("port", "9090", "port for the load balancer to listen on")
	backends := flag.String("backends", "http://localhost:8080,http://localhost:8081", "comma-seperated backend URLs")
	flag.Parse()

	//pool
	pool, err := loadbalancer.NewPool(strings.Split(*backends, ","))
	if err != nil {
		log.Printf("Invalid backend list: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, pool)
	})

	log.Printf("Load balancer listening on :%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}

func handleRequest(w http.ResponseWriter, r *http.Request, pool *loadbalancer.Pool)  {
	backend := pool.NextBackend()
	if backend == nil {
		http.Error(w, "no backends available", http.StatusServiceUnavailable)
		return
	}

	log.Printf("Routing %s %s to %s", r.Method, r.URL.Path, backend.URL)

	resp, err := http.Get(backend.URL.String() + r.URL.Path)
	if err != nil {
		http.Error(w, "backend unreachable", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(name, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}
