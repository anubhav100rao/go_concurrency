package play

import (
	"context"
	"time"
)

const rateLimit = time.Second / 10 // 10 calls per second

// Client is an interface that calls something with a payload.
type Client interface {
	Call(*Payload)
}

// Payload is some payload a Client would send in a call.
type Payload struct{}

// BurstRateLimitCall allows burst rate limiting client calls with the
// payloads.
func BurstRateLimitCall(ctx context.Context, client Client, payloads []*Payload, burstLimit int) {
	throttle := make(chan time.Time, burstLimit)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		ticker := time.NewTicker(rateLimit)
		defer ticker.Stop()
		for t := range ticker.C {
			select {
			case throttle <- t:
			case <-ctx.Done():
				return // exit goroutine when surrounding function returns
			}
		}
	}()

	for _, payload := range payloads {
		<-throttle // rate limit our client calls
		go client.Call(payload)
	}
}
