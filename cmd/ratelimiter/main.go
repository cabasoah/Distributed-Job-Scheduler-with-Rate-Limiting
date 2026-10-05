package main

import (
	"flag"
	"log"
	"net/http"
)

func main()  {
	port := flag.String("port", "8083", "port for the rate limiter API to listen on")
	flag.Parse()

	http.HandleFunc("/unlimited", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Unlimited! Let's go!"))
	})

	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Limited! Don't over use me!"))
	})

	log.Printf("Rate limiter API listening on %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}