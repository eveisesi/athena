package processor

import (
	"context"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/member"
	"github.com/korovkin/limiter"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Run()
	SetScopeMap(athena.ScopeMap)
}

type service struct {
	logger *logrus.Logger

	cache  cache.Service
	member member.Service

	scopes athena.ScopeMap
}

func NewService(logger *logrus.Logger, cache cache.Service, member member.Service) Service {

	s := &service{
		logger: logger,

		cache:  cache,
		member: member,
	}

	return s
}

func (s *service) SetScopeMap(scopes athena.ScopeMap) {
	s.scopes = scopes
}

func (s *service) Run() {

	if len(s.scopes) == 0 {
		panic("scopes are not set. Please run SetScopeMap and provide a list of scopes and pointers to their resolvers")
	}

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

func (s *service) processMember(ctx context.Context, memberID uint) {

	member, err := s.member.Member(ctx, memberID)
	if err != nil {
		s.logger.WithError(err).WithField("memberID", memberID).Errorln("failed to fetch member by ID")
		return
	}

	member, err = s.member.ValidateToken(ctx, member)
	if err != nil {
		s.logger.WithError(err).WithField("memberID", member.ID).Errorln("failed to verify token is valid")
		return
	}

	entry := s.logger.WithFields(logrus.Fields{
		"member": member.ID,
	})

	// Member Retrieve successfully. Loop over the scopes array calling the functions in the scope map
	for i, scope := range member.Scopes {

		entry := s.logger.WithField("scope", scope.Scope)

		// If the scope expiry is valid, that means it has previously been called,
		// and if the expiry is after the current time, that means that the cache timer
		// has not expired yet, so attempting to update the data now will not yield any fresh results
		if scope.Expiry.Valid && scope.Expiry.Time.After(time.Now()) {
			entry.Info("skipping valid and active scope")
			time.Sleep(time.Second)
			continue
		}

		if _, ok := s.scopes[scope.Scope]; !ok {
			// entry.Error("scope not supported")
			time.Sleep(time.Second)
			continue
		}

		for _, resolver := range s.scopes[scope.Scope] {
			entry := entry.WithField("name", resolver.Name)
			entry.Info()

			etag, err := resolver.Func(ctx, member)
			if err != nil {
				entry.WithError(err).Errorln()
				continue
			}

			entry.Info("scope resolved successfully")

			if etag != nil {
				scope.Expiry.SetValid(etag.CachedUntil)
			}

			member.Scopes[i] = scope

			time.Sleep(time.Second)
		}

	}

	_, err = s.member.UpdateMember(ctx, member)
	if err != nil {
		entry.WithError(err).Error("failed to update member")
	}

	entry.Info("member processed successfully")

}
