package asset

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/glue"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

type Service interface {
	EmptyMemberAssets(ctx context.Context, member *athena.Member) (*athena.Etag, error)
}

type service struct {
	logger *logrus.Logger

	cache cache.Service
	esi   esi.Service

	universe universe.Service

	assets athena.MemberAssetsRepository
}

const (
	serviceIdentifier = "Asset Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe universe.Service, assets athena.MemberAssetsRepository) Service {
	return &service{
		logger: logger,

		cache: cache,
		esi:   esi,

		universe: universe,

		assets: assets,
	}
}

func (s *service) EmptyMemberAssets(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterAssets, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated etag")
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	return s.FetchMemberAssets(ctx, member, etag)

}

func (s *service) FetchMemberAssets(ctx context.Context, member *athena.Member, etag *athena.Etag) (*athena.Etag, error) {

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberAssets",
	})

	newAssets, _, err := s.esi.GetCharacterAssets(ctx, member, make([]*athena.MemberAsset, 0))
	if err != nil {
		entry.WithError(err).Error("failed to fetch member assets from ESI")
		return nil, fmt.Errorf("failed to fetch member assets from ESI")
	}

	if len(newAssets) > 0 {
		s.resolveMemberAssetsAttributes(ctx, member, newAssets)

		oldAssets, err := s.assets.MemberAssets(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			entry.WithError(err).Error("failed to fetch existing member assets from DB")
			return nil, fmt.Errorf("failed to fetch existing member assets from DB")
		}

		err = s.diffAndCreateOrUpdateAssets(ctx, member, oldAssets, newAssets)
		if err != nil {
			entry.WithError(err).Error("failed to diff and create or update assets")
			return nil, fmt.Errorf("failed to diff and create or update assets")
		}

	}

	etag, err = s.esi.Etag(ctx, esi.GetCharacterAssets, esi.ModWithMember(member))
	if err != nil {
		entry.WithError(err).Error("failed to fetch updated etag")
		return nil, fmt.Errorf("failed to fetch updated etag")
	}

	return etag, nil

}

func (s *service) diffAndCreateOrUpdateAssets(ctx context.Context, member *athena.Member, oldAssets, newAssets []*athena.MemberAsset) error {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "diffAndCreateOrUpdateAssets",
	})

	assetsToCreate := make([]*athena.MemberAsset, 0, len(newAssets))
	assetsToUpdate := make([]*athena.MemberAsset, 0, len(oldAssets))
	assetsToDelete := make([]*athena.MemberAsset, 0, len(oldAssets))

	// Build a Map of Asset Item IDs to Items for the old assets
	mapOldAssets := make(map[uint64]*athena.MemberAsset)
	for _, asset := range oldAssets {
		mapOldAssets[asset.ItemID] = asset
	}

	for _, asset := range newAssets {
		if _, ok := mapOldAssets[asset.ItemID]; !ok {
			assetsToCreate = append(assetsToCreate, asset)
		} else if diff := deep.Equal(mapOldAssets[asset.ItemID], asset); len(diff) > 0 {
			assetsToUpdate = append(assetsToUpdate, asset)
		}
	}

	newAssetMap := make(map[uint64]*athena.MemberAsset)
	for _, asset := range newAssets {
		newAssetMap[asset.ItemID] = asset
	}

	for _, asset := range oldAssets {
		if _, ok := newAssetMap[asset.ItemID]; !ok {
			assetsToDelete = append(assetsToDelete, asset)
		}
	}

	if len(assetsToDelete) > 0 {
		_, err := s.assets.DeleteMemberAssets(ctx, member.ID, assetsToDelete)
		if err != nil {
			entry.WithError(err).Error("failed to delete member assets")
			return fmt.Errorf("failed to delete member assets")
		}
	}
	if len(assetsToUpdate) > 0 {
		for _, asset := range assetsToUpdate {
			entry := entry.WithField("item_id", asset.ItemID)
			_, err := s.assets.UpdateMemberAssets(ctx, member.ID, asset.ItemID, asset)
			if err != nil {
				entry.WithError(err).Error("failed to update member assets")
				return fmt.Errorf("failed to update member assets")
			}
		}
	}
	if len(assetsToCreate) > 0 {
		_, err := s.assets.CreateMemberAssets(ctx, member.ID, assetsToCreate)
		if err != nil {
			entry.WithError(err).Error("failed to create member assets")
			return fmt.Errorf("failed to create member assets")
		}

	}

	return nil
}

func (s *service) resolveMemberAssetsAttributes(ctx context.Context, member *athena.Member, assets []*athena.MemberAsset) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "resolveMemberAssetsAttributes",
	})

	mapTypeIDs := make(map[uint]bool)
	mapLocationStationIDs := make(map[uint]bool)
	mapLocationSystemIDs := make(map[uint]bool)
	mapLocationStructureIDs := make(map[uint64]bool)

	// Loop through assets once and get a unique list of ids
	for _, asset := range assets {

		if _, ok := mapTypeIDs[asset.TypeID]; !ok {
			mapTypeIDs[asset.TypeID] = true
		}

		locationType := glue.ResolveIDTypeFromIDRange(asset.LocationID)
		switch locationType {
		case glue.IDTypeStation:
			if _, ok := mapLocationStationIDs[uint(asset.LocationID)]; !ok {
				mapLocationStationIDs[uint(asset.LocationID)] = true
			}
		case glue.IDTypeSolarSystem:
			if _, ok := mapLocationSystemIDs[uint(asset.LocationID)]; !ok {
				mapLocationSystemIDs[uint(asset.LocationID)] = true
			}
		case glue.IDTypeStructure:
			if _, ok := mapLocationStructureIDs[asset.LocationID]; !ok {
				mapLocationStructureIDs[asset.LocationID] = true
			}
		}

	}

	for k := range mapLocationStructureIDs {
		_, err := s.universe.Structure(ctx, member, k)
		if err != nil {
			entry.WithError(err).WithField("structure_id", k).Error("failed to resolve structure id")
		}
	}

	for k := range mapLocationStationIDs {
		_, err := s.universe.Station(ctx, k)
		if err != nil {
			entry.WithError(err).WithField("station_id", k).Error("failed to resolve station id")
		}
	}

	for k := range mapLocationSystemIDs {
		_, err := s.universe.SolarSystem(ctx, k)
		if err != nil {
			entry.WithError(err).WithField("system_id", k).Error("failed to resolve solar system id")
		}
	}

}
