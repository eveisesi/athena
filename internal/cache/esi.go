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
	ATHENA_ESI_ERROR_RESET = "athena::esi::error::reset"
	ATHENA_ESI_ERROR_COUNT = "athena::esi::error::count"
	ATHENA_ESI_OK          = "athena::esi::ok"         // 200
	ATHENA_ESI_UNCHANGED   = "athena::esi::unchanged"  // 304
	ATHENA_ESI_RESTRICTED  = "athena::esi::restricted" // 420
	ATHENA_ESI_4XX         = "athena::esi::4xx"        // Does not include 420s. Those are in the calm down set
	ATHENA_ESI_5XX         = "athena::esi::5xx"
)

func (s *service) SetEsiErrorReset(ctx context.Context, reset int64) {
	mx.Lock()
	defer mx.Unlock()
	s.client.Set(ctx, ATHENA_ESI_ERROR_RESET, reset, 0)

}

func (s *service) SetESIErrCount(ctx context.Context, count int64) {
	mx.Lock()
	defer mx.Unlock()
	s.client.Set(ctx, ATHENA_ESI_ERROR_COUNT, count, 0)
}

func (s *service) SetESITracking(ctx context.Context, code int, ts int64) {
	mx.Lock()
	defer mx.Unlock()
	z := &redis.Z{Score: float64(ts), Member: strconv.FormatInt(ts, 10)}

	switch n := code; {
	case n == http.StatusOK:
		s.client.ZAdd(ctx, ATHENA_ESI_OK, z)
	case n == http.StatusNotModified:
		s.client.ZAdd(ctx, ATHENA_ESI_UNCHANGED, z)
	case n == 420:
		s.client.ZAdd(ctx, ATHENA_ESI_RESTRICTED, z)
	case n >= 400 && n < 500:
		s.client.ZAdd(ctx, ATHENA_ESI_4XX, z)
	case n >= 500:
		s.client.ZAdd(ctx, ATHENA_ESI_5XX, z)
	}

}
