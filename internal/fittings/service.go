package fittings

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/glue"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Service interface{}

type service struct {
	logger *logrus.Logger

	redis       redis.Client
	esi         esi.Service
	alliance    alliance.Service
	character   character.Service
	corporation corporation.Service
	universe    universe.Service

	fittings athena.MemberFittingsRepository
}

const (
	serviceIdentifier = "Contract Service"
)

func (s *service) FetchMemberFittings(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterFittings, esi.ModWithCharacterID(member.ID))
	if err != nil {
		return nil, glue.FormatError(serviceIdentifier, "Failed to fetch etag object: %w", err)
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {
			return etag, nil
		}

		petag = etag.Etag
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberFittings",
	})

	etag, _, err = s.esi.HeadCharacterFittings(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to exec head request for member fittings from ESI")
		return nil, fmt.Errorf("failed to exec head request for member fittings from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, nil
	}

	fittings, _, _, err := s.esi.GetCharacterFittings(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member fittings from ESI")
		return nil, fmt.Errorf("failed to fetch member fittings from ESI")
	}

	s.resolveFittingAttributes(ctx, fittings)

	return etag, nil

}

func (s *service) resolveFittingAttributes(ctx context.Context, fittings []*athena.MemberFitting) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "resolveFittingAttributes",
	})

	for _, fitting := range fittings {
		entry := entry.WithField("ship_type_id", fitting.ShipTypeID)
		_, err := s.universe.Type(ctx, fitting.ShipTypeID)
		if err != nil {
			entry.WithError(err).Error("failed to resolve fitting ship type id to name")
			continue
		}

		for i, item := range fitting.Items {
			// We could take care of this later since it is technically outside the scope
			// of this function, but we're looping through it now
			fitting.Items[i].FittingID = fitting.FittingID

			entry := entry.WithField("item_type_id", item.TypeID)
			_, err := s.universe.Type(ctx, fitting.ShipTypeID)
			if err != nil {
				entry.WithError(err).Error("failed to resolve fitting item to name")
				continue
			}

		}

	}

}
