package clone

import (
	"context"
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
	MemberClones(ctx context.Context, member *athena.Member) (*athena.MemberClones, error)
	EmptyMemberImplants(ctx context.Context, member *athena.Member) error
	MemberImplants(ctx context.Context, member *athena.Member) (*athena.MemberImplants, error)
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

func (s *service) MemberClones(ctx context.Context, member *athena.Member) (*athena.MemberClones, error) {

	var upsert = "update"

	clones, err := s.cache.MemberClones(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if clones == nil {

		clones, err = s.clones.MemberClones(ctx, member.ID.Hex())
		if err != nil {
			upsert = "create"
			clones = &athena.MemberClones{
				MemberID: member.ID,
			}
		}

	}

	if clones.CachedUntil.After(time.Now()) {
		return clones, nil
	}

	clones, _, err = s.esi.GetCharacterClones(ctx, member, clones)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch clones for member")
		return nil, err
	}

	s.resolveCloneAttributes(ctx, member, clones)

	switch upsert {
	case "create":
		clones, err = s.clones.CreateMemberClones(ctx, clones)
		if err != nil {
			return nil, err
		}
	case "update":
		clones, err = s.clones.UpdateMemberClones(ctx, member.ID.Hex(), clones)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetMemberClones(ctx, clones.MemberID.Hex(), clones)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return clones, nil
}

func (s *service) resolveCloneAttributes(ctx context.Context, member *athena.Member, clones *athena.MemberClones) {

	switch clones.HomeLocation.LocationType {
	case "structure":
		_, err := s.universe.Structure(ctx, member, clones.HomeLocation.LocationID)
		if err != nil {
			s.logger.WithError(err).WithField("structure_id", clones.HomeLocation.LocationID).WithField("member_id", member.ID.Hex()).Error("failed to resolve structure id")
		}
	case "station":
		_, err := s.universe.Station(ctx, int(clones.HomeLocation.LocationID))
		if err != nil {
			s.logger.WithError(err).WithField("structure_id", clones.HomeLocation.LocationID).Error("failed to resolve station id")
		}
	}

	for _, jumpClone := range clones.JumpClones {

		switch jumpClone.LocationType {
		case "structure":
			_, err := s.universe.Structure(ctx, member, jumpClone.LocationID)
			if err != nil {
				s.logger.WithError(err).WithField("structure_id", jumpClone.LocationID).WithField("member_id", member.ID.Hex()).Error("failed to resolve structure id")
			}
		case "station":
			_, err := s.universe.Station(ctx, int(jumpClone.LocationID))
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

func (s *service) MemberImplants(ctx context.Context, member *athena.Member) (*athena.MemberImplants, error) {

	var upsert = "update"

	implants, err := s.cache.MemberImplants(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if implants == nil {

		implants, err = s.clones.MemberImplants(ctx, member.ID.Hex())
		if err != nil {
			upsert = "create"
			implants = &athena.MemberImplants{
				MemberID: member.ID,
			}
		}

	}

	if implants.CachedUntil.After(time.Now()) {
		return implants, nil
	}

	implants, _, err = s.esi.GetCharacterImplants(ctx, member, implants)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch implants for member")
		return nil, err
	}

	s.resolveImplantAttributes(ctx, implants)

	switch upsert {
	case "create":
		implants, err = s.clones.CreateMemberImplants(ctx, implants)
		if err != nil {
			return nil, err
		}
	case "update":
		implants, err = s.clones.UpdateMemberImplants(ctx, member.ID.Hex(), implants)
		if err != nil {
			return nil, err
		}
	}

	err = s.cache.SetMemberImplants(ctx, implants.MemberID.Hex(), implants)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
	}

	return implants, nil
}

func (s *service) resolveImplantAttributes(ctx context.Context, implants *athena.MemberImplants) {

	if len(implants.Raw) == 0 {
		return
	}

	implants.Implants = make([]*athena.Type, len(implants.Raw))

	for i, raw := range implants.Raw {

		implant, err := s.universe.Type(ctx, raw)
		if err != nil {
			s.logger.WithError(err).WithField("implant_type_id", implant).Error("failed to resolve implent type id")

		}

		implants.Implants[i] = implant
	}

	implants.Raw = nil

}
