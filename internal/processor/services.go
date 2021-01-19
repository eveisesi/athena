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
	scopeMap[athena.ReadLocationV1] = []athena.ScopeResolver{
		athena.ScopeResolver{
			Name: "MemberLocation",
			Func: s.location.EmptyMemberLocation,
		},
	}

	scopeMap[athena.ReadOnlineV1] = []athena.ScopeResolver{
		athena.ScopeResolver{
			Name: "MemberOnline",
			Func: s.location.EmptyMemberOnline,
		},
	}

	scopeMap[athena.ReadShipV1] = []athena.ScopeResolver{
		athena.ScopeResolver{
			Name: "MemberShip",
			Func: s.location.EmptyMemberShip,
		},
	}

	s.scopes = scopeMap

}

func (s *service) Run() {

	s.logger.Info("Processor is running")

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
	}

}

func (s *service) processMember(ctx context.Context, memberID string) {

	member, err := s.member.Member(ctx, memberID)
	if err != nil {
		s.logger.WithError(err).WithField("memberID", memberID).Errorln("failed to fetch member by ID")
		return
	}

	member, err = s.member.ValidateToken(ctx, member)
	if err != nil {
		s.logger.WithError(err).WithField("memberID", member.ID.Hex()).Errorln("failed to verify token is valid")
		return
	}

	// Member Retrieve successfully. Loop over the scopes array calling the functions in the scope map
	for _, scope := range member.Scopes {

		// If the scope expiry is valid, that means it has previously been called,
		// and if the expiry is after the current time, that means that the cache timer
		// has been expired yet, so updating the data now will not yield any fresh results
		if scope.Expiry.Valid && scope.Expiry.Time.After(time.Now()) {
			fmt.Println(scope.Scope, "is not valid")
			continue
		}

		if _, ok := s.scopes[scope.Scope]; !ok {
			s.logger.WithField("scope", scope.Scope.String()).Error("scope not supported")
			continue
		}

		for _, resolver := range s.scopes[scope.Scope] {
			s.logger.WithFields(logrus.Fields{
				"member": member.ID.Hex(),
				"scope":  scope.Scope,
				"name":   resolver.Name,
			}).Infoln()

			err := resolver.Func(ctx, member)
			if err != nil {
				s.logger.WithError(err).WithField("field", scope.Scope).Errorln()
			}
			time.Sleep(time.Second)
		}

	}

}
