# Workshop with REST vd gRPC

## Build REST
```
$docker compose build rest
```

## Build gRPC
```
$docker compose build grpc
```


## Build Benchmark
```
$docker compose build bench
```

## Start REST server and gRPC server
```
$docker compose up -d rest grpc
$docker compose ps
```

Check REST API
```
$curl "http://localhost:8080/hello?message=hello%20world"

{"message":"hello world"}
```

## Benchmark
```
$docker compose run --rm bench -n 20000 -c 200
```

Result
```
Running with n=20000 c=200 message="hello world"

[REST (Echo)]
Requests: 20000 (errors: 0)
Wall time: 3.503073794s
Throughput: 5709.27 req/s
Latency avg: 34.83287ms
Latency p50: 8.855334ms
Latency p95: 140.195875ms
Latency p99: 324.344126ms

[gRPC]
Requests: 20000 (errors: 0)
Wall time: 452.190333ms
Throughput: 44229.16 req/s
Latency avg: 4.514317ms
Latency p50: 4.100875ms
Latency p95: 8.006625ms
Latency p99: 17.633125ms
```