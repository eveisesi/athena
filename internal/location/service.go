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

	"github.com/sirupsen/logrus"
)

type Service interface {
	EmptyMemberLocation(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberLocation(ctx context.Context, memberID uint) (*athena.MemberLocation, error)

	EmptyMemberShip(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberShip(ctx context.Context, memberID uint) (*athena.MemberShip, error)

	EmptyMemberOnline(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberOnline(ctx context.Context, memberID uint) (*athena.MemberOnline, error)
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	universe universe.Service

	location athena.MemberLocationRepository
}

const (
	serviceIdentifier = "Location Service"
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

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "EmptyMemberLocation",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterLocation, esi.ModWithCharacterID(member.ID))
	if err != nil {
		entry.WithError(err).Error("failed to fetch etag object")
		return nil, fmt.Errorf("failed to fetch etag object")
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {

			return etag, nil
		}

		petag = etag.Etag
	}

	location, etag, _, err := s.esi.GetCharacterLocation(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member location from ESI")
		return nil, fmt.Errorf("failed to fetch member location from ESI")
	}

	if etag.Etag == petag {
		return etag, nil
	}

	s.resolveLocationAttributes(ctx, member, location)

	existing, err := s.location.MemberLocation(ctx, member.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member location from DB")
		return nil, fmt.Errorf("failed to fetch member location from DB")
	}

	switch existing == nil || errors.Is(err, sql.ErrNoRows) {
	case true:
		location, err = s.location.CreateMemberLocation(ctx, member.ID, location)
		if err != nil {
			entry.WithError(err).Error("failed to create member location in database")
			return nil, fmt.Errorf("failed to create member location in database")
		}
	case false:
		location, err = s.location.UpdateMemberLocation(ctx, member.ID, location)
		if err != nil {
			entry.WithError(err).Error("failed to update member location in database")
			return nil, fmt.Errorf("failed to update member location in database")
		}
	}

	err = s.cache.SetMemberLocation(ctx, location.MemberID, location)
	if err != nil {
		entry.WithError(err).Error("failed to cache member location")
	}

	return etag, nil

}

func (s *service) MemberLocation(ctx context.Context, memberID uint) (*athena.MemberLocation, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "EmptyMemberLocation",
	})

	location, err := s.cache.MemberLocation(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member location from cache")
		return nil, fmt.Errorf("failed to fetch member location from cache")
	}

	if location != nil {
		return location, nil
	}

	location, err = s.location.MemberLocation(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member location from DB")
		return nil, fmt.Errorf("failed to fetch member location from DB")
	}

	if location != nil {
		err = s.cache.SetMemberLocation(ctx, memberID, location)
		if err != nil {
			entry.WithError(err).Error("failed to cache member location")
		}
	}

	return location, nil

}

func (s *service) resolveLocationAttributes(ctx context.Context, member *athena.Member, location *athena.MemberLocation) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "resolveLocationAttributes",
	})

	_, err := s.universe.SolarSystem(ctx, location.SolarSystemID)
	if err != nil {
		entry.WithError(err).WithFields(logrus.Fields{
			"solar_system_id": location.SolarSystemID,
		}).Error("failed to resolve solar system")
		return
	}

	if location.StationID.Valid {
		_, err = s.universe.Station(ctx, location.StationID.Uint)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"station_id": location.StationID.Uint,
			}).Error("failed to resolve station")
			return
		}
	}

	if location.StructureID.Valid {
		_, err := s.universe.Structure(ctx, member, location.StructureID.Uint64)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"member_id":    member.ID,
				"structure_id": location.StructureID.Uint64,
			}).Error("failed to resolve structure")
			return
		}
	}

}

func (s *service) EmptyMemberShip(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "EmptyMemberShip",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterShip, esi.ModWithCharacterID(member.ID))
	if err != nil {
		entry.WithError(err).Error("failed to fetch etag object")
		return nil, fmt.Errorf("failed to fetch etag object")
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {

			return etag, nil
		}

		petag = etag.Etag
	}

	ship, etag, _, err := s.esi.GetCharacterShip(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member ship from ESI")
		return nil, fmt.Errorf("failed to fetch member ship from ESI")
	}

	if etag.Etag == petag {
		return etag, nil
	}

	s.resolveShipAttributes(ctx, member, ship)

	existing, err := s.location.MemberOnline(ctx, member.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member ship from DB")
		return nil, fmt.Errorf("failed to fetch member ship from DB")
	}

	switch existing == nil || errors.Is(err, sql.ErrNoRows) {
	case true:
		ship, err = s.location.CreateMemberShip(ctx, member.ID, ship)
		if err != nil {
			entry.WithError(err).Error("failed to create member ship in database")
			return nil, fmt.Errorf("failed to create member ship in database")
		}
	case false:
		ship, err = s.location.UpdateMemberShip(ctx, member.ID, ship)
		if err != nil {
			entry.WithError(err).Error("failed to update member ship in database")
			return nil, fmt.Errorf("failed to update member ship in database")
		}
	}

	err = s.cache.SetMemberShip(ctx, ship.MemberID, ship)
	if err != nil {
		entry.WithError(err).Error("failed to cache member ship")
	}

	return etag, nil

}

func (s *service) MemberShip(ctx context.Context, memberID uint) (*athena.MemberShip, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberShip",
	})

	ship, err := s.cache.MemberShip(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member ship from cache")
		return nil, fmt.Errorf("failed to fetch member ship from cache")
	}

	if ship != nil {
		return ship, nil
	}

	ship, err = s.location.MemberShip(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member ship from DB")
		return nil, fmt.Errorf("failed to fetch member ship from DB")
	}

	if ship != nil {
		err = s.cache.SetMemberShip(ctx, memberID, ship)
		if err != nil {
			entry.WithError(err).Error("failed to cache member ship")
		}
	}

	return ship, nil

}

func (s *service) resolveShipAttributes(ctx context.Context, member *athena.Member, ship *athena.MemberShip) {

	_, err := s.universe.Type(ctx, ship.ShipTypeID)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"member_id":    member.ID,
			"service":      serviceIdentifier,
			"method":       "resolveShipAttributes",
			"ship_type_id": ship.ShipTypeID,
		}).Error("failed to resolve ship type id")
	}

}

func (s *service) EmptyMemberOnline(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "EmptyMemberOnline",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterOnline, esi.ModWithCharacterID(member.ID))
	if err != nil {
		entry.WithError(err).Error("failed to fetch etag object")
		return nil, fmt.Errorf("failed to fetch etag object")
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {

			return etag, nil
		}

		petag = etag.Etag
	}

	online, etag, _, err := s.esi.GetCharacterOnline(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member online from ESI")
		return nil, fmt.Errorf("failed to fetch member online from ESI")
	}

	if etag.Etag == petag {
		return etag, nil
	}

	existing, err := s.location.MemberOnline(ctx, member.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member online from DB")
		return nil, fmt.Errorf("failed to fetch member online from DB")
	}

	switch existing == nil || errors.Is(err, sql.ErrNoRows) {
	case true:
		online, err = s.location.CreateMemberOnline(ctx, member.ID, online)
		if err != nil {
			entry.WithError(err).Error("failed to update member online in database")
			return nil, fmt.Errorf("failed to update member online in database")
		}
	case false:
		online, err = s.location.UpdateMemberOnline(ctx, member.ID, online)
		if err != nil {
			entry.WithError(err).Error("failed to create member online in database")
			return nil, fmt.Errorf("failed to create member online in database")
		}
	}

	err = s.cache.SetMemberOnline(ctx, online.MemberID, online)
	if err != nil {
		entry.WithError(err).Error("failed to cache member online")
	}

	return etag, nil

}

func (s *service) MemberOnline(ctx context.Context, memberID uint) (*athena.MemberOnline, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberOnline",
	})

	online, err := s.cache.MemberOnline(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member online from cache")
		return nil, fmt.Errorf("failed to fetch member online from cache")
	}

	if online != nil {
		return online, nil
	}

	online, err = s.location.MemberOnline(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member online from DB")
		return nil, fmt.Errorf("failed to fetch member online from DB")
	}

	if online != nil {
		err = s.cache.SetMemberOnline(ctx, memberID, online)
		if err != nil {
			entry.WithError(err).Error("failed to cache member online")
		}
	}

	return online, nil

}
