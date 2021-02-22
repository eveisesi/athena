package fittings

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/glue"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/go-redis/redis/v8"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

type Service interface{}

type service struct {
	logger *logrus.Logger

	cache cache

	esi      esi.Service
	universe universe.Service

	fittings athena.MemberFittingsRepository
}

const (
	serviceIdentifier = "Contract Service"
)

func New(logger *logrus.Logger, redis *redis.Client, esi esi.Service, universe universe.Service, fittings athena.MemberFittingsRepository) Service {
	return &service{
		logger: logger,

		cache:    newCacher(redis),
		esi:      esi,
		universe: universe,

		fittings: fittings,
	}
}

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

	if len(fittings) == 0 {
		return etag, nil
	}

	s.resolveFittingAttributes(ctx, fittings)

	existing, err := s.fittings.MemberFittings(ctx, member.ID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member fittings from DB")
		return nil, fmt.Errorf("failed to fetch member fittings from DB")
	}

	fittings, err = s.processFittings(ctx, member, existing, fittings)
	if err != nil {
		entry.WithError(err).Error("failed to process fittings")
		return nil, fmt.Errorf("failed to process fittings")
	}

	if len(fittings) > 0 {
		// Cache Fittings
	}

	return etag, nil

}

func (s *service) processFittings(ctx context.Context, member *athena.Member, old, new []*athena.MemberFitting) ([]*athena.MemberFitting, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "processFittings",
	})

	fittingsToCreate := make([]*athena.MemberFitting, 0, len(new))
	fittingsToUpdate := make([]*athena.MemberFitting, 0, len(new))
	fittingsToDelete := make([]*athena.MemberFitting, 0, len(new))

	oldFittingMap := make(map[uint]*athena.MemberFitting)
	for _, fitting := range old {
		oldFittingMap[fitting.FittingID] = fitting
	}

	for _, fitting := range new {
		if _, ok := oldFittingMap[fitting.FittingID]; !ok {
			fittingsToCreate = append(fittingsToCreate, fitting)
			// Deep support recursion. This means it SHOULD dig into the item array and compare those as well?
		} else if diff := deep.Equal(oldFittingMap[fitting.FittingID], fitting); len(diff) > 0 {
			fittingsToUpdate = append(fittingsToUpdate, fitting)
		}
	}

	newFittingsMap := make(map[uint]*athena.MemberFitting)
	for _, fitting := range new {
		newFittingsMap[fitting.FittingID] = fitting
	}

	for _, fitting := range old {
		if _, ok := newFittingsMap[fitting.FittingID]; !ok {
			fittingsToDelete = append(fittingsToDelete, fitting)
		}
	}

	if len(fittingsToDelete) > 0 {
		for _, fitting := range fittingsToDelete {
			_, err := s.fittings.DeleteMemberFitting(ctx, member.ID, fitting.FittingID)
			if err != nil {
				entry.WithField("fitting_id", fitting.FittingID).WithError(err).Error("failed to remove fitting")
			}
		}
	}

	if len(fittingsToCreate) > 0 {
		_, err := s.fittings.CreateMemberFittings(ctx, member.ID, fittingsToCreate)
		if err != nil {
			entry.WithError(err).Error("failed to create member fittings in db")
			return nil, fmt.Errorf("failed to create member fittings in db")
		}
	}

	if len(fittingsToUpdate) > 0 {
		for _, fitting := range fittingsToUpdate {
			entry := entry.WithField("fitting_id", fitting.FittingID)
			_, err := s.fittings.UpdateMemberFitting(ctx, member.ID, fitting.FittingID, fitting)
			if err != nil {
				entry.WithError(err).Error("failed to remove fitting")
				continue
			}

			// err = s.cache.MemberFittings(ctx, )

			_, err = s.fittings.DeleteMemberFittingItems(ctx, member.ID, fitting.FittingID)
			if err != nil {
				entry.WithError(err).Error("failed to drop fitting items for fit")
				continue
			}

			_, err = s.fittings.CreateMemberFittingItems(ctx, member.ID, fitting.FittingID, fitting.Items)
			if err != nil {
				entry.WithError(err).Error("failed to create fitting items for fit")
			}
		}
	}

	return nil, nil

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
