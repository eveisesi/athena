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
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service interface {
	EmptyMemberClones(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberClones(ctx context.Context, member *athena.Member) (*athena.MemberClones, *athena.Etag, error)
	EmptyMemberImplants(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberImplants(ctx context.Context, member *athena.Member) ([]*athena.MemberImplant, *athena.Etag, error)
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	universe universe.Service

	clones athena.CloneRepository
}

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

	etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberClones(ctx, member)

	return etag, err

}

func (s *service) MemberClones(ctx context.Context, member *athena.Member) (*athena.MemberClones, *athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}
	upsert := false
	cached := true

	clones, err := s.cache.MemberClones(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if clones == nil {
		cached = false
		clones, err = s.clones.MemberClones(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			upsert = true
			clones = &athena.MemberClones{MemberID: member.ID}
			err = s.esi.ResetEtag(ctx, etag)
			if err != nil {
				return nil, nil, err
			}
		}

	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && !upsert && clones != nil {

		if !cached {
			err = s.cache.SetMemberClones(ctx, member.ID, clones)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}

		}

		return clones, etag, nil

	}

	clones, etag, _, err = s.esi.GetCharacterClones(ctx, member, clones)
	if err != nil {
		return nil, nil, err
	}

	s.resolveCloneAttributes(ctx, member, clones)

	switch upsert {
	case true:
		clones, err = s.clones.CreateMemberClones(ctx, clones)
		if err != nil {
			return nil, nil, err
		}
	case false:
		clones, err = s.clones.UpdateMemberClones(ctx, clones)
		if err != nil {
			return nil, nil, err
		}
	}

	err = s.cache.SetMemberClones(ctx, clones.MemberID, clones)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return clones, etag, nil
}

func (s *service) resolveCloneAttributes(ctx context.Context, member *athena.Member, clones *athena.MemberClones) {

	clones.MemberID = member.ID

	homeClone := clones.HomeLocation
	switch homeClone.LocationType {
	case "structure":
		_, err := s.universe.Structure(ctx, member, homeClone.LocationID)
		if err != nil {
			s.logger.WithError(err).WithField("structure_id", homeClone.LocationID).WithField("member_id", member.ID).Error("failed to resolve structure id")
		}
	case "station":
		_, err := s.universe.Station(ctx, uint(homeClone.LocationID))
		if err != nil {
			s.logger.WithError(err).WithField("structure_id", homeClone.LocationID).Error("failed to resolve station id")
		}
	}

	for _, jumpClone := range clones.JumpClones {

		switch jumpClone.LocationType {
		case "structure":
			_, err := s.universe.Structure(ctx, member, jumpClone.LocationID)
			if err != nil {
				s.logger.WithError(err).WithField("structure_id", jumpClone.LocationID).WithField("member_id", member.ID).Error("failed to resolve structure id")
			}
		case "station":
			_, err := s.universe.Station(ctx, uint(jumpClone.LocationID))
			if err != nil {
				s.logger.WithError(err).WithField("station_id", jumpClone.LocationID).Error("failed to resolve station id")
			}
		}

		for _, implant := range jumpClone.Implants {
			_, err := s.universe.Type(ctx, implant)
			if err != nil {
				s.logger.WithError(err).WithField("implant_type_id", jumpClone.LocationID).Error("failed to resolve implent type id")
			}
		}

	}

}

func (s *service) EmptyMemberImplants(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberImplants(ctx, member)

	return etag, err

}

func (s *service) MemberImplants(ctx context.Context, member *athena.Member) ([]*athena.MemberImplant, *athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterImplants, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	cached := true

	implants, err := s.cache.MemberImplants(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if len(implants) == 0 {
		cached = false
		implants, err = s.clones.MemberImplants(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, err
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && len(implants) > 0 {

		if !cached {
			err = s.cache.SetMemberImplants(ctx, member.ID, implants)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return implants, etag, nil
	}

	newImplants, etag, _, err := s.esi.GetCharacterImplants(ctx, member, make([]uint, 0))
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch implants for member")
		return nil, nil, err
	}

	implants, err = s.resolveImplantAttributes(ctx, member, newImplants)
	if err != nil {
		return nil, nil, err
	}

	if len(implants) > 0 {
		err = s.cache.SetMemberImplants(ctx, member.ID, implants)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
		}
	}

	return implants, etag, nil

}

func (s *service) resolveImplantAttributes(ctx context.Context, member *athena.Member, new []uint) ([]*athena.MemberImplant, error) {

	implants := make([]*athena.MemberImplant, len(new))

	for i, raw := range new {

		implant, err := s.universe.Type(ctx, raw)
		if err != nil {
			err = fmt.Errorf("[Clones Service] Failed to resolve implent type id %d: %w", raw, err)
			newrelic.FromContext(ctx).AddAttribute("implant_type_id", raw)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		implants[i] = &athena.MemberImplant{
			MemberID:  member.ID,
			ImplantID: implant.ID,
		}
	}

	_, err := s.clones.DeleteMemberImplants(ctx, member.ID)
	if err != nil {
		err = fmt.Errorf("[Clone Service] Failed to delete member %d implants: %w", member.ID, err)
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	implants, err = s.clones.CreateMemberImplants(ctx, member.ID, implants)
	if err != nil {
		err = fmt.Errorf("[Clone Service] Failed to create implants for member %d: %w", member.ID, err)
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	return implants, nil

}
