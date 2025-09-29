package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	hello "demo/hellobench/gen/hello"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type stats struct {
	count     int
	errors    int64
	durations []time.Duration
	totalWall time.Duration
}

func runHTTP(n, c int, url string) stats {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: c}}
	var idx int64
	var errs int64
	durs := make([]time.Duration, n)
	wg := sync.WaitGroup{}
	wg.Add(c)
	start := time.Now()
	for i := 0; i < c; i++ {
		go func() {
			defer wg.Done()
			for {
				j := int(atomic.AddInt64(&idx, 1)) - 1
				if j >= n {
					return
				}
				t0 := time.Now()
				resp, err := client.Get(url)
				if err != nil {
					atomic.AddInt64(&errs, 1)
					continue
				}
				_ = resp.Body.Close()
				durs[j] = time.Since(t0)
			}
		}()
	}
	wg.Wait()
	return stats{count: n, errors: errs, durations: durs, totalWall: time.Since(start)}
}

func runGRPC(n, c int, addr, message string) stats {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return stats{count: n, errors: int64(n), durations: nil, totalWall: 0}
	}
	defer conn.Close()
	cli := hello.NewHelloServiceClient(conn)

	var idx int64
	var errs int64
	durs := make([]time.Duration, n)
	wg := sync.WaitGroup{}
	wg.Add(c)
	start := time.Now()
	for i := 0; i < c; i++ {
		go func() {
			defer wg.Done()
			for {
				j := int(atomic.AddInt64(&idx, 1)) - 1
				if j >= n {
					return
				}
				t0 := time.Now()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := cli.SayHello(ctx, &hello.HelloRequest{Message: message})
				cancel()
				if err != nil {
					atomic.AddInt64(&errs, 1)
					continue
				}
				durs[j] = time.Since(t0)
			}
		}()
	}
	wg.Wait()
	return stats{count: n, errors: errs, durations: durs, totalWall: time.Since(start)}
}

func percentile(durs []time.Duration, p float64) time.Duration {
	if len(durs) == 0 {
		return 0
	}
	s := append([]time.Duration(nil), durs...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	k := int(p*float64(len(s)-1) + 0.5)
	if k < 0 {
		k = 0
	}
	if k >= len(s) {
		k = len(s) - 1
	}
	return s[k]
}

func summarize(name string, s stats) string {
	ok := s.count - int(s.errors)
	rps := float64(ok) / s.totalWall.Seconds()
	p50 := percentile(s.durations, 0.50)
	p95 := percentile(s.durations, 0.95)
	p99 := percentile(s.durations, 0.99)
	avg := time.Duration(0)
	for _, d := range s.durations {
		avg += d
	}
	if ok > 0 {
		avg /= time.Duration(ok)
	}
	b := &strings.Builder{}
	fmt.Fprintf(b, "\n[%s]\n", name)
	fmt.Fprintf(b, "Requests: %d (errors: %d)\n", s.count, s.errors)
	fmt.Fprintf(b, "Wall time: %v\n", s.totalWall)
	fmt.Fprintf(b, "Throughput: %.2f req/s\n", rps)
	fmt.Fprintf(b, "Latency avg: %v\n", avg)
	fmt.Fprintf(b, "Latency p50: %v\n", p50)
	fmt.Fprintf(b, "Latency p95: %v\n", p95)
	fmt.Fprintf(b, "Latency p99: %v\n", p99)
	return b.String()
}

func main() {
	var (
		n        = flag.Int("n", 10000, "total requests per target")
		c        = flag.Int("c", 100, "concurrency")
		restURL  = flag.String("rest-url", "http://rest:8080/hello?message=hello%20world", "REST URL")
		grpcAddr = flag.String("grpc-addr", "grpc:50051", "gRPC host:port")
		message  = flag.String("message", "hello world", "message payload")
		target   = flag.String("target", "both", "one of: rest|grpc|both")
		warmup   = flag.Int("warmup", 200, "warmup requests per target (not counted)")
	)
	flag.Parse()

	fmt.Printf("Running with n=%d c=%d message=\"%s\"\n", *n, *c, *message)

	if *target == "rest" || *target == "both" {
		_ = runHTTP(*warmup, *c, *restURL) // warmup
		s := runHTTP(*n, *c, *restURL)
		fmt.Print(summarize("REST (Echo)", s))
	}
	if *target == "grpc" || *target == "both" {
		_ = runGRPC(*warmup, *c, *grpcAddr, *message) // warmup
		s := runGRPC(*n, *c, *grpcAddr, *message)
		fmt.Print(summarize("gRPC", s))
	}
}
