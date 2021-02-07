package location

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/universe"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
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

	cache    cache.Service
	esi      esi.Service
	universe universe.Service

	location athena.MemberLocationRepository
}

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe universe.Service, location athena.MemberLocationRepository) Service {
	return &service{
		logger: logger,

		cache:    cache,
		esi:      esi,
		universe: universe,

		location: location,
	}
}

func (s *service) EmptyMemberLocation(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberLocation(ctx, member)

	return err

}

func (s *service) MemberLocation(ctx context.Context, member *athena.Member) (*athena.MemberLocation, error) {

	var upsert string = "update"

	location, err := s.cache.MemberLocation(ctx, member.ID)
	if err != nil {
		return nil, err
	}

	if location == nil {
		location, err = s.location.MemberLocation(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			upsert = "create"
			location = &athena.MemberLocation{
				MemberID: member.ID,
			}
		}
	}

	if location.CachedUntil.After(time.Now()) {
		return location, nil
	}

	location, _, err = s.esi.GetCharacterLocation(ctx, member, location)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, err
	}

	s.resolveLocationAttributes(ctx, member, location)

	switch upsert {
	case "create":
		location, err = s.location.CreateMemberLocation(ctx, location)
		if err != nil {
			return nil, err
		}

	case "update":
		location, err = s.location.UpdateMemberLocation(ctx, member.ID, location)
		if err != nil {
			return nil, err
		}

	}

	err = s.cache.SetMemberLocation(ctx, location.MemberID, location)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return location, nil
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

func (s *service) EmptyMemberShip(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberShip(ctx, member)

	return err

}

func (s *service) MemberShip(ctx context.Context, member *athena.Member) (*athena.MemberShip, error) {

	var upsert string = "update"

	ship, err := s.cache.MemberShip(ctx, member.ID)
	if err != nil {
		return nil, err
	}

	if ship == nil {
		ship, err = s.location.MemberShip(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			upsert = "create"
			ship = &athena.MemberShip{
				MemberID: member.ID,
			}
		}
	}

	if ship.CachedUntil.After(time.Now()) {
		return ship, nil
	}

	ship, _, err = s.esi.GetCharacterShip(ctx, member, ship)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch location for member")
		return nil, err
	}

	s.resolveShipAttributes(ctx, member, ship)

	switch upsert {
	case "create":
		ship, err = s.location.CreateMemberShip(ctx, ship)
		if err != nil {
			return nil, err
		}
	case "update":
		ship, err = s.location.UpdateMemberShip(ctx, member.ID, ship)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetMemberShip(ctx, ship.MemberID, ship)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return ship, nil
}

func (s *service) resolveShipAttributes(ctx context.Context, member *athena.Member, ship *athena.MemberShip) {

	_, err := s.universe.Type(ctx, ship.ShipTypeID)
	if err != nil {
		s.logger.WithError(err).WithField("ship_type_id", ship.ShipTypeID).Error("failed to resolve ship type id")
	}

}

func (s *service) EmptyMemberOnline(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberOnline(ctx, member)

	return err

}

func (s *service) MemberOnline(ctx context.Context, member *athena.Member) (*athena.MemberOnline, error) {

	var upsert string = "update"

	online, err := s.cache.MemberOnline(ctx, member.ID)
	if err != nil {
		return nil, err
	}

	if online == nil {
		online, err = s.location.MemberOnline(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			upsert = "create"
			online = &athena.MemberOnline{
				MemberID: member.ID,
			}
		}
	}

	if online.CachedUntil.After(time.Now()) {
		return nil, err
	}

	online, _, err = s.esi.GetCharacterOnline(ctx, member, online)
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
		online, err = s.location.UpdateMemberOnline(ctx, member.ID, online)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetMemberOnline(ctx, online.MemberID, online)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return online, nil

}
