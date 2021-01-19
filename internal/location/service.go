package location

import (
	"context"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	MemberLocation(ctx context.Context, member *athena.Member) (*athena.MemberLocation, error)
	EmptyMemberLocation(ctx context.Context, member *athena.Member) error

	MemberOnline(ctx context.Context, member *athena.Member) (*athena.MemberOnline, error)
	EmptyMemberOnline(ctx context.Context, member *athena.Member) error

	MemberShip(ctx context.Context, member *athena.Member) (*athena.MemberShip, error)
	EmptyMemberShip(ctx context.Context, member *athena.Member) error
}

type service struct {
	logger *logrus.Logger

	cache cache.Service
	esi   esi.Service

	location athena.MemberLocationRepository
}

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, location athena.MemberLocationRepository) Service {
	return &service{
		logger: logger,

		cache: cache,
		esi:   esi,

		location: location,
	}
}

func (s *service) EmptyMemberLocation(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberLocation(ctx, member)

	return err

}

func (s *service) MemberLocation(ctx context.Context, member *athena.Member) (*athena.MemberLocation, error) {

	var upsert string = "update"

	location, err := s.cache.MemberLocation(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if location == nil {
		location, err = s.location.MemberLocation(ctx, member.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, err
		}

		if err == mongo.ErrNoDocuments {
			upsert = "create"
			location = &athena.MemberLocation{
				MemberID: member.ID,
			}
		}
	}

	if location.CachedUntil.After(time.Now()) {
		return location, nil
	}

	location, _, err = s.esi.GetCharactersCharacterIDLocation(ctx, member, location)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, nil
	}

	switch upsert {
	case "create":
		location, err = s.location.CreateMemberLocation(ctx, location)
		if err != nil {
			return nil, err
		}

	case "update":
		location, err = s.location.UpdateMemberLocation(ctx, location.ID.Hex(), location)
		if err != nil {
			return nil, err
		}

	}

	err = s.cache.SetMemberLocation(ctx, location.MemberID.Hex(), location)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return location, nil
}

func (s *service) EmptyMemberShip(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberShip(ctx, member)

	return err

}

func (s *service) MemberShip(ctx context.Context, member *athena.Member) (*athena.MemberShip, error) {

	var upsert string = "update"

	ship, err := s.cache.MemberShip(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if ship == nil {
		ship, err = s.location.MemberShip(ctx, member.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, err
		}

		if err == mongo.ErrNoDocuments {
			upsert = "create"
			ship = &athena.MemberShip{
				MemberID: member.ID,
			}
		}
	}

	if ship.CachedUntil.After(time.Now()) {
		return ship, nil
	}

	ship, _, err = s.esi.GetCharactersCharacterIDShip(ctx, member, ship)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, err
	}

	switch upsert {
	case "create":
		ship, err = s.location.CreateMemberShip(ctx, ship)
		if err != nil {
			return nil, err
		}
	case "update":
		ship, err = s.location.UpdateMemberShip(ctx, ship.ID.Hex(), ship)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetMemberShip(ctx, ship.MemberID.Hex(), ship)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return ship, nil
}

func (s *service) EmptyMemberOnline(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberOnline(ctx, member)

	return err

}

func (s *service) MemberOnline(ctx context.Context, member *athena.Member) (*athena.MemberOnline, error) {

	var upsert string = "update"

	online, err := s.cache.MemberOnline(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if online == nil {
		online, err = s.location.MemberOnline(ctx, member.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, err
		}

		if err == mongo.ErrNoDocuments {
			upsert = "create"
			online = &athena.MemberOnline{
				MemberID: member.ID,
			}
		}
	}

	if online.CachedUntil.After(time.Now()) {
		return nil, err
	}

	online, _, err = s.esi.GetCharactersCharacterIDOnline(ctx, member, online)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, err
	}

	switch upsert {
	case "create":
		online, err = s.location.CreateMemberOnline(ctx, online)
		if err != nil {
			return nil, err
		}
	case "update":
		online, err = s.location.UpdateMemberOnline(ctx, online.ID.Hex(), online)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetMemberOnline(ctx, online.MemberID.Hex(), online)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return online, nil

}
