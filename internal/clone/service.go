package clone

import (
	"context"
	"database/sql"
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
	EmptyMemberClones(ctx context.Context, member *athena.Member) error
	MemberClones(ctx context.Context, member *athena.Member) (*athena.MemberHomeClone, []*athena.MemberJumpClone, error)
	EmptyMemberImplants(ctx context.Context, member *athena.Member) error
	MemberImplants(ctx context.Context, member *athena.Member) ([]*athena.MemberImplant, error)
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

func (s *service) EmptyMemberClones(ctx context.Context, member *athena.Member) error {
	_, err := s.MemberImplants(ctx, member)
	return err
}

func (s *service) MemberClones(ctx context.Context, member *athena.Member) (*athena.MemberHomeClone, []*athena.MemberJumpClone, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	cached := true

	clone, err := s.cache.MemberHomeClone(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	clones, err := s.cache.MemberJumpClones(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if clone == nil || clones == nil || len(clones) == 0 {
		cached = false
		clone, err = s.clones.MemberHomeClone(ctx, member.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, err
		}

		if err == sql.ErrNoRows {
			clone = &athena.MemberHomeClone{MemberID: member.ID}
		}

		clones, err = s.clones.MemberJumpClones(ctx, member.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, err
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && len(clones) > 0 && clone != nil {

		if !cached {
			err = s.cache.SetMemberHomeClone(ctx, member.ID, clone)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}

			err = s.cache.SetMemberJumpClones(ctx, member.ID, clones)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return clone, clones, nil
	}

	_, _, _, err = s.esi.GetCharacterClones(ctx, member, clone, clones)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch clones for member")
		return clone, clones, err
	}

	// s.resolveCloneAttributes(ctx, member, clone, clones)

	// switch upsert {
	// case "create":
	// 	clones, err = s.clones.CreateMemberClones(ctx, clones)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// case "update":
	// 	clones, err = s.clones.UpdateMemberClones(ctx, member.ID, clones)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// err = s.cache.SetMemberClones(ctx, clones.MemberID, clones)
	// if err != nil {
	// 	newrelic.FromContext(ctx).NoticeError(err)
	// }

	return clone, clones, nil
}

func (s *service) resolveCloneAttributes(ctx context.Context, member *athena.Member, clone *athena.MemberHomeClone, clones []*athena.MemberJumpClone) {

	switch clone.LocationType {
	case "structure":
		_, err := s.universe.Structure(ctx, member, clone.LocationID)
		if err != nil {
			s.logger.WithError(err).WithField("structure_id", clone.LocationID).WithField("member_id", member.ID).Error("failed to resolve structure id")
		}
	case "station":
		_, err := s.universe.Station(ctx, uint(clone.LocationID))
		if err != nil {
			s.logger.WithError(err).WithField("structure_id", clone.LocationID).Error("failed to resolve station id")
		}
	}

	for _, jumpClone := range clones {

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

func (s *service) EmptyMemberImplants(ctx context.Context, member *athena.Member) error {
	_, err := s.MemberImplants(ctx, member)
	return err
}

func (s *service) MemberImplants(ctx context.Context, member *athena.Member) ([]*athena.MemberImplant, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterImplants, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	cached := true

	implants, err := s.cache.MemberImplants(ctx, member.ID)
	if err != nil {
		return nil, err
	}

	if len(implants) == 0 {
		cached = false
		implants, err = s.clones.MemberImplants(ctx, member.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && len(implants) > 0 {

		if !cached {
			err = s.cache.SetMemberImplants(ctx, member.ID, implants)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return implants, nil
	}

	newImplants, _, err := s.esi.GetCharacterImplants(ctx, member, make([]uint, 0))
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch implants for member")
		return nil, err
	}

	implants, err = s.resolveImplantAttributes(ctx, member, newImplants)
	if err != nil {
		return nil, err
	}

	if len(implants) > 0 {
		err = s.cache.SetMemberImplants(ctx, member.ID, implants)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
		}
	}

	return implants, nil
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
			MemberID: member.ID,
			TypeID:   implant.ID,
		}
	}

	_, err := s.clones.DeleteMemberImplants(ctx, member.ID)
	if err != nil {
		err = fmt.Errorf("[Clone Service] Failed to delete member %d implants: %w", member.ID, err)
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	implants, err = s.clones.CreateMemberImplants(ctx, implants)
	if err != nil {
		err = fmt.Errorf("[Clone Service] Failed to create implants for member %d: %w", member.ID, err)
		newrelic.FromContext(ctx).NoticeError(err)
		return nil, err
	}

	return implants, nil

}
