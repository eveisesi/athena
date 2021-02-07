package location

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/universe"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service interface {
	EmptyMemberLocation(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberLocation(ctx context.Context, member *athena.Member) (*athena.MemberLocation, *athena.Etag, error)

	EmptyMemberOnline(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberOnline(ctx context.Context, member *athena.Member) (*athena.MemberOnline, *athena.Etag, error)

	EmptyMemberShip(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberShip(ctx context.Context, member *athena.Member) (*athena.MemberShip, *athena.Etag, error)
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	universe universe.Service

	location athena.MemberLocationRepository
}

const (
	errPrefix = "[Location Service]"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe universe.Service, location athena.MemberLocationRepository) Service {
	return &service{
		logger: logger,

		cache:    cache,
		esi:      esi,
		universe: universe,

		location: location,
	}
}

func (s *service) EmptyMemberLocation(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterLocation, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("%s Failed to fetch Etag Object: %w", errPrefix, err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberLocation(ctx, member)

	return etag, err

}

func (s *service) MemberLocation(ctx context.Context, member *athena.Member) (*athena.MemberLocation, *athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterLocation, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("%s Failed to fetch Etag Object: %w", errPrefix, err)
	}

	exists := true
	cached := true

	location, err := s.cache.MemberLocation(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if location == nil {
		cached = false
		location, err = s.location.MemberLocation(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, err
		}

		if location == nil || errors.Is(err, sql.ErrNoRows) {
			exists = false
			location = &athena.MemberLocation{
				MemberID: member.ID,
			}
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetMemberLocation(ctx, member.ID, location)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return location, etag, nil

	}

	location, etag, _, err = s.esi.GetCharacterLocation(ctx, member, location)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, nil, err
	}

	s.resolveLocationAttributes(ctx, member, location)

	switch exists {
	case true:
		location, err = s.location.UpdateMemberLocation(ctx, member.ID, location)
		if err != nil {
			return nil, etag, err
		}

	case false:
		location, err = s.location.CreateMemberLocation(ctx, member.ID, location)
		if err != nil {
			return nil, etag, err
		}
	}

	err = s.cache.SetMemberLocation(ctx, location.MemberID, location)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return location, etag, nil

}

func (s *service) resolveLocationAttributes(ctx context.Context, member *athena.Member, location *athena.MemberLocation) {

	_, err := s.universe.SolarSystem(ctx, location.SolarSystemID)
	if err != nil {
		return
	}

	if location.StationID.Valid {
		_, err = s.universe.Station(ctx, location.StationID.Uint)
		if err != nil {
			s.logger.WithError(err).WithField("station_id", location.StationID.Uint).Error("failed to resolve station")
			return
		}
	}

	if location.StructureID.Valid {
		_, err := s.universe.Structure(ctx, member, location.StructureID.Uint64)
		if err != nil {
			s.logger.WithError(err).WithField("structure_id", location.StructureID.Uint64).Error("failed to resolve structure")
			return
		}
	}

}

func (s *service) EmptyMemberShip(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterShip, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("%s Failed to fetch Etag Object: %w", errPrefix, err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberShip(ctx, member)

	return etag, err

}

func (s *service) MemberShip(ctx context.Context, member *athena.Member) (*athena.MemberShip, *athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterShip, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("%s Failed to fetch Etag Object: %w", errPrefix, err)
	}

	exists := true
	cached := true

	ship, err := s.cache.MemberShip(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if ship == nil {
		cached = false
		ship, err = s.location.MemberShip(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, err
		}

		if ship == nil || errors.Is(err, sql.ErrNoRows) {
			exists = false
			ship = &athena.MemberShip{
				MemberID: member.ID,
			}
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetMemberShip(ctx, member.ID, ship)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return ship, etag, nil

	}

	ship, etag, _, err = s.esi.GetCharacterShip(ctx, member, ship)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, etag, err
	}

	s.resolveShipAttributes(ctx, member, ship)

	switch exists {
	case true:
		ship, err = s.location.UpdateMemberShip(ctx, member.ID, ship)
		if err != nil {
			return nil, etag, err
		}
	case false:
		ship, err = s.location.CreateMemberShip(ctx, member.ID, ship)
		if err != nil {
			return nil, etag, err
		}
	}

	err = s.cache.SetMemberShip(ctx, ship.MemberID, ship)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return ship, etag, nil

}

func (s *service) resolveShipAttributes(ctx context.Context, member *athena.Member, ship *athena.MemberShip) {

	_, err := s.universe.Type(ctx, ship.ShipTypeID)
	if err != nil {
		s.logger.WithError(err).WithField("ship_type_id", ship.ShipTypeID).Error("failed to resolve ship type id")
	}

}

func (s *service) EmptyMemberOnline(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterOnline, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("%s Failed to fetch Etag Object: %w", errPrefix, err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberOnline(ctx, member)

	return etag, err

}

func (s *service) MemberOnline(ctx context.Context, member *athena.Member) (*athena.MemberOnline, *athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterOnline, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("%s Failed to fetch Etag Object: %w", errPrefix, err)
	}

	exists := true
	cached := true

	online, err := s.cache.MemberOnline(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if online == nil {
		cached = false
		online, err = s.location.MemberOnline(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, err
		}

		if online == nil || errors.Is(err, sql.ErrNoRows) {
			exists = false
			online = &athena.MemberOnline{
				MemberID: member.ID,
			}
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && exists {

		if !cached {
			err = s.cache.SetMemberOnline(ctx, member.ID, online)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return online, etag, nil

	}

	online, etag, _, err = s.esi.GetCharacterOnline(ctx, member, online)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, nil, err
	}

	switch exists {
	case false:
		online, err = s.location.CreateMemberOnline(ctx, member.ID, online)
		if err != nil {
			return nil, etag, err
		}
	case true:
		online, err = s.location.UpdateMemberOnline(ctx, member.ID, online)
		if err != nil {
			return nil, etag, err
		}
	}

	err = s.cache.SetMemberOnline(ctx, online.MemberID, online)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return online, etag, nil

}
