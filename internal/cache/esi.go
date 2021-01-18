package cache

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

type esiService interface {
	SetEsiErrorReset(ctx context.Context, reset int64)
	SetESIErrCount(ctx context.Context, count int64)
	SetESITracking(ctx context.Context, code int, ts int64)
}

var (
	mx sync.Mutex
)

const (
	keyESIErrorReset = "athena::esi::error::reset"
	keyESIErrorCount = "athena::esi::error::count"
	keyESIOk         = "athena::esi::ok"         // 200
	keyESIUnchanged  = "athena::esi::unchanged"  // 304
	keyESIRestricted = "athena::esi::restricted" // 420
	keyESI4xx        = "athena::esi::4xx"        // Does not include 420s. Those are in the calm down set
	keyESI5xx        = "athena::esi::5xx"
)

func (s *service) SetEsiErrorReset(ctx context.Context, reset int64) {
	mx.Lock()
	defer mx.Unlock()
	s.client.Set(ctx, keyESIErrorReset, reset, 0)

}

func (s *service) SetESIErrCount(ctx context.Context, count int64) {
	mx.Lock()
	defer mx.Unlock()
	s.client.Set(ctx, keyESIErrorCount, count, 0)
}

func (s *service) SetESITracking(ctx context.Context, code int, ts int64) {
	mx.Lock()
	defer mx.Unlock()
	z := &redis.Z{Score: float64(ts), Member: strconv.FormatInt(ts, 10)}

	switch n := code; {
	case n == http.StatusOK:
		s.client.ZAdd(ctx, keyESIOk, z)
	case n == http.StatusNotModified:
		s.client.ZAdd(ctx, keyESIUnchanged, z)
	case n == 420:
		s.client.ZAdd(ctx, keyESIRestricted, z)
	case n >= 400 && n < 500:
		s.client.ZAdd(ctx, keyESI4xx, z)
	case n >= 500:
		s.client.ZAdd(ctx, keyESI5xx, z)
	}

}
