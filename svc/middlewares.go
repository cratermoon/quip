package svc

import (
	"context"
	"log"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
)

// MakeMetricsMiddleware instruments a call with a Histogram and a Counter
func MakeMetricsMiddleware(h metrics.Histogram, c metrics.Counter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			begin := time.Now()
			defer h.Observe(time.Since(begin).Seconds())
			return next(ctx, request)
		}
	}

}

func MakeLoggingMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			log.Println("msg", "start")
			defer log.Println("msg", "finish")
			return next(ctx, request)
		}
	}
}
