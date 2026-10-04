package main

import (
	"flag"
	"log"
	"net/http"
)

// Basic backend that runs on port 9000
func main()  {
	port := flag.String("port", "9000", "port for backend to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recieved request from %s\n%s %s %s", r.RemoteAddr, r.Method, r.URL.Path, r.Proto)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Distributed job scheduler and rate limiting running"))
	})

	log.Printf("Backend server listening on :%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}