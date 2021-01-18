package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type processorService interface {
	PushIDToProcessorQueue(ctx context.Context, memberID primitive.ObjectID)
	PopFromProcessorQueue(ctx context.Context, count int) ([]string, error)
	ProcessorQueueCount(ctx context.Context) (int64, error)
}

const (
	keyProcessorMemberIDQueue = "athena::processor::members"
)

func (s *service) PushIDToProcessorQueue(ctx context.Context, memberID primitive.ObjectID) {

	mx.Lock()
	defer mx.Unlock()
	ts := time.Now().UnixNano()
	z := &redis.Z{Score: float64(ts), Member: memberID.Hex()}

	s.client.ZAdd(ctx, keyProcessorMemberIDQueue, z)

}

func (s *service) PopFromProcessorQueue(ctx context.Context, count int) ([]string, error) {

	mx.Lock()
	defer mx.Unlock()

	results, err := s.client.ZPopMin(ctx, keyProcessorMemberIDQueue, int64(count)).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[PopFromProcessorQueue] Failed to retrieve records from processor queue: %w", err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	slc := make([]string, len(results))
	for i, v := range results {
		msg := v.Member.(string)
		slc[i] = msg
	}

	return slc, nil

}

func (s *service) ProcessorQueueCount(ctx context.Context) (int64, error) {

	mx.Lock()
	defer mx.Unlock()

	results, err := s.client.ZCount(ctx, keyProcessorMemberIDQueue, "-inf", "+inf").Result()
	if err != nil {
		return 0, fmt.Errorf("[ProcessorQueueCount] Failed to retrieve the count of records from processor queue: %w", err)
	}

	return results, nil

}
