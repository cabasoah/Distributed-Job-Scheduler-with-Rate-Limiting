# Distributed-Job-Scheduler-with-Rate-Limiting
A concurrent job scheduler and a rate limiter (combined into one system)


## Testing the load balancer schduling using Round robin
python -m http.server 8080 && python -m http.server 8081 && python -m http.server 8082

go run ./cmd/loadbalancer --port 9090 --backends http://localhost:8080,http://localhost:8081,http://localhost:8082

curl http://localhost:9090/   -  run repeatedly, watch it alternate

# Testing the rate limiter
go run ./cmd/ratelimiter --port 8083