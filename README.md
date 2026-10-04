# Distributed-Job-Scheduler-with-Rate-Limiting
A concurrent job scheduler and a rate limiter (combined into one system)


## Testing the load balancer schduling using Round robin
python -m http.server 8080
python -m http.server 8081  
go run ./cmd/lb --port 9090 --backends http://localhost:8080,http://localhost:8081
curl http://localhost:9090/   -  run repeatedly, watch it alternate