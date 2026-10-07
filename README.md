# Distributed-Job-Scheduler-with-Rate-Limiting
A concurrent job scheduler and a rate limiter (combined into one system)


## Testing the load balancer schduling using Round robin
python -m http.server 8080 && python -m http.server 8081 && python -m http.server 8082

go run ./cmd/loadbalancer --port 9090 --backends http://localhost:8080,http://localhost:8081,http://localhost:8082

curl http://localhost:9090/   -  run repeatedly, watch it alternate

# Testing the rate limiter
go run ./cmd/ratelimiter --strategy fixed-window --port 8083

for i in $(seq 1 10); do curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:8083/limited & done; wait