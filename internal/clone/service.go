package clone

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
	EmptyMemberClones(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberClones(ctx context.Context, memberID uint) (*athena.MemberClones, error)
	EmptyMemberImplants(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberImplants(ctx context.Context, memberID uint) ([]*athena.MemberImplant, error)
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	universe universe.Service

	clones athena.CloneRepository
}

const (
	serviceIdentifier = "Clone Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe universe.Service, clones athena.CloneRepository) Service {
	return &service{
		logger: logger,

		cache:    cache,
		esi:      esi,
		universe: universe,

		clones: clones,
	}
}

func (s *service) EmptyMemberClones(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "EmptyMemberClones",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, esi.ModWithCharacterID(member.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object")
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {
			return etag, nil
		}

		petag = etag.Etag
	}

	clones, etag, _, err := s.esi.GetCharacterClones(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member clones from ESI")
		return nil, fmt.Errorf("failed to fetch member clones from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, nil
	}

	s.resolveCloneAttributes(ctx, member, clones)

	existing, err := s.clones.MemberClones(ctx, member.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member clones from DB")
		return nil, fmt.Errorf("failed to fetch member clones from DB")
	}

	switch existing == nil || errors.Is(err, sql.ErrNoRows) {
	case true:
		clones, err = s.clones.CreateMemberClones(ctx, clones)
		if err != nil {
			entry.WithError(err).Error("failed to create member clones in DB")
			return nil, fmt.Errorf("failed to create member clones in DB")
		}
	case false:
		clones, err = s.clones.UpdateMemberClones(ctx, clones)
		if err != nil {
			entry.WithError(err).Error("failed to update member clones in DB")
			return nil, fmt.Errorf("failed to update member clones in DB")
		}
	}

	err = s.cache.SetMemberClones(ctx, clones.MemberID, clones)
	if err != nil {
		entry.WithError(err).Error("failed to cache member clones")
	}

	return etag, err

}

func (s *service) MemberClones(ctx context.Context, memberID uint) (*athena.MemberClones, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberClones",
	})

	clones, err := s.cache.MemberClones(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member clones from cache")
		return nil, fmt.Errorf("failed to fetch member clones from cache")
	}

	if clones != nil {
		return clones, nil
	}

	clones, err = s.clones.MemberClones(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member clones from DB")
		return nil, fmt.Errorf("failed to fetch member clones from DB")
	}

	if clones != nil {
		err = s.cache.SetMemberClones(ctx, memberID, clones)
		if err != nil {
			entry.WithError(err).Error("failed to cache member clones")
		}
	}
	return clones, nil

}

func (s *service) resolveCloneAttributes(ctx context.Context, member *athena.Member, clones *athena.MemberClones) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "resolveCloneAttributes",
	})

	clones.MemberID = member.ID

	homeClone := clones.HomeLocation
	switch homeClone.LocationType {
	case "structure":
		_, err := s.universe.Structure(ctx, member, homeClone.LocationID)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"location_id":   homeClone.LocationID,
				"location_type": homeClone.LocationType,
			}).Error("failed to resolve structure id")
		}
	case "station":
		_, err := s.universe.Station(ctx, uint(homeClone.LocationID))
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"location_id":   homeClone.LocationID,
				"location_type": homeClone.LocationType,
			}).Error("failed to resolve station id")
		}
	}

	for _, jumpClone := range clones.JumpClones {

		switch jumpClone.LocationType {
		case "structure":
			_, err := s.universe.Structure(ctx, member, jumpClone.LocationID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"location_id":   jumpClone.LocationID,
					"location_type": jumpClone.LocationType,
				}).Error("failed to resolve structure id")
			}
		case "station":
			_, err := s.universe.Station(ctx, uint(jumpClone.LocationID))
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"location_id":   jumpClone.LocationID,
					"location_type": jumpClone.LocationType,
				}).Error("failed to resolve station id")
			}
		}

		for _, implant := range jumpClone.Implants {
			_, err := s.universe.Type(ctx, uint(implant))
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"implant_type_id": implant,
				}).Error("failed to resolve station id")
			}
		}

	}

}

func (s *service) EmptyMemberImplants(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "EmptyMemberImplants",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterImplants, esi.ModWithCharacterID(member.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	newImplants, etag, _, err := s.esi.GetCharacterImplants(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member implants from ESI")
		return nil, fmt.Errorf("failed to fetch member implants from ESI")
	}

	implants, err := s.resolveImplantAttributes(ctx, member, newImplants)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve member implants")
	}

	if len(implants) > 0 {
		err = s.cache.SetMemberImplants(ctx, member.ID, implants)
		if err != nil {
			entry.WithError(err).Error("failed to cache member implants")
		}
	}

	return etag, err

}

func (s *service) MemberImplants(ctx context.Context, memberID uint) ([]*athena.MemberImplant, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberImplants",
	})

	implants, err := s.cache.MemberImplants(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member implants from cache")
		return nil, fmt.Errorf("failed to fetch member implants from cache")
	}

	if len(implants) > 0 {
		return implants, nil
	}

	implants, err = s.clones.MemberImplants(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member implants from DB")
		return nil, fmt.Errorf("failed to fetch member implants from DB")
	}

	if len(implants) > 0 {
		err = s.cache.SetMemberImplants(ctx, memberID, implants)
		if err != nil {
			entry.WithError(err).Error("failed to cache member implants")
		}
	}

	return implants, nil

}

func (s *service) resolveImplantAttributes(ctx context.Context, member *athena.Member, new []uint) ([]*athena.MemberImplant, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "resolveImplantAttributes",
	})

	implants := make([]*athena.MemberImplant, len(new))

	for i, raw := range new {

		implant, err := s.universe.Type(ctx, raw)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"implant_type_id": raw,
			}).Error("failed to resolve implant type id")
			continue
		}

		implants[i] = &athena.MemberImplant{
			MemberID:  member.ID,
			ImplantID: implant.ID,
		}

	}

	_, err := s.clones.DeleteMemberImplants(ctx, member.ID)
	if err != nil {
		entry.WithError(err).Error("failed to delete member implants")
		return implants, fmt.Errorf("failed to delete member implants")
	}

	if len(implants) > 0 {
		implants, err = s.clones.CreateMemberImplants(ctx, member.ID, implants)
		if err != nil {
			entry.WithError(err).Error("failed to create implants for member")
		}

	}

	return implants, err

}
