package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/member"
	"github.com/korovkin/limiter"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Run()
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	member   member.Service
	location location.Service
	scopes   athena.ScopeMap
}

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, member member.Service, location location.Service) Service {

	s := &service{
		logger:   logger,
		cache:    cache,
		esi:      esi,
		member:   member,
		location: location,
	}

	s.buildScopeMap()

	return s
}

func (s *service) buildScopeMap() {

	scopeMap := make(athena.ScopeMap)
	scopeMap[athena.READ_LOCATION_V1] = s.location.MemberLocation
	scopeMap[athena.READ_ONLINE_V1] = s.location.MemberOnline
	scopeMap[athena.READ_SHIP_V1] = s.location.MemberShip

	s.scopes = scopeMap

}

func (s *service) Run() {

	limit := limiter.NewConcurrencyLimiter(10)
	for {
		ctx := context.Background()

		count, err := s.cache.ProcessorQueueCount(ctx)
		if err != nil {
			s.logger.WithError(err).Errorln("[processor.Run]")
			time.Sleep(time.Second)
			continue
		}

		if count == 0 {
			s.logger.Debug("record count is 0, sleep 5 seconds")
			time.Sleep(time.Second * 5)
			continue
		}

		results, err := s.cache.PopFromProcessorQueue(ctx, 10)
		if err != nil {
			s.logger.WithError(err).Errorln("[processor.Run]")
			time.Sleep(time.Second)
			continue
		}

		for _, result := range results {
			limit.Execute(func() {
				s.processMember(ctx, result)
			})
		}
		limit.Wait()
	}

}

func (s *service) processMember(ctx context.Context, memberID string) {

	fmt.Printf("Received %s\n", memberID)
	time.Sleep(time.Second * 10)

}
