package contract

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/glue"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

type Service interface {
	FetchMemberContracts(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberContracts(ctx context.Context, memberID, page uint) ([]*athena.MemberContract, error)

	FetchMemberContractItems(ctx context.Context, member *athena.Member, contract *athena.MemberContract) (*athena.Etag, error)
	MemberContractItems(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractItem, error)

	FetchMemberContractBids(ctx context.Context, member *athena.Member, contract *athena.MemberContract) (*athena.Etag, error)
	MemberContractBids(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractBid, error)
}

type service struct {
	logger *logrus.Logger

	cache       cache.Service
	esi         esi.Service
	alliance    alliance.Service
	character   character.Service
	corporation corporation.Service
	universe    universe.Service

	contracts athena.MemberContractRepository
}

const (
	serviceIdentifier = "Contract Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe universe.Service, alliance alliance.Service, character character.Service, corporation corporation.Service, contracts athena.MemberContractRepository) Service {
	return &service{
		logger: logger,

		cache:       cache,
		esi:         esi,
		universe:    universe,
		alliance:    alliance,
		character:   character,
		corporation: corporation,

		contracts: contracts,
	}
}

func (s *service) FetchMemberContracts(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContracts, esi.ModWithCharacterID(member.ID))
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
		"method":    "FetchMemberContacts",
	})

	etag, res, err := s.esi.HeadCharacterContracts(ctx, member.ID, 1, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to exec head request for member contacts from ESI")
		return nil, fmt.Errorf("failed to exec head request for member contacts from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, nil
	}

	pages := esi.RetrieveXPagesFromHeader(res.Header)

	for page := uint(1); page <= pages; page++ {
		entry := entry.WithField("source_page", page)

		etag, err := s.esi.Etag(ctx, esi.GetCharacterContracts, esi.ModWithCharacterID(member.ID), esi.ModWithPage(page))
		if err != nil {
			entry.WithError(err).Error("failed to fetch page of member contacts from ESI")
			return nil, fmt.Errorf("failed to fetch page of member contacts from ESI")
		}

		var petag string
		if etag != nil {
			if etag.CachedUntil.After(time.Now()) {
				continue
			}

			petag = etag.Etag
		}

		contracts, etag, _, err := s.esi.GetCharacterContracts(ctx, member.ID, page, member.AccessToken.String)
		if err != nil {
			entry.WithError(err).Error("failed to fetch member contracts from ESI")
			return nil, fmt.Errorf("failed to fetch member contracts from ESI")
		}

		if petag != "" && petag == etag.Etag {
			continue
		}

		existingContracts, err := s.contracts.MemberContracts(ctx, member.ID, athena.NewInOperator("contract_id", contractIDsFromSliceContracts(contracts)))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		_ = s.diffAndUpdateContracts(ctx, member, page, existingContracts, contracts)

		if len(contracts) > 0 {
			err = s.cache.SetMemberContracts(ctx, member.ID, page, contracts)
			if err != nil {
				entry.WithError(err).Error("failed to cache member contracts ")
			}
		}
	}

	return etag, nil

}

func contractIDsFromSliceContracts(s []*athena.MemberContract) []uint {
	o := make([]uint, 0, len(s))
	for _, c := range s {
		o = append(o, c.ContractID)
	}
	return o
}

func (s *service) diffAndUpdateContracts(ctx context.Context, member *athena.Member, page uint, old, new []*athena.MemberContract) error {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "diffAndUpdateContracts",
	})

	for _, contract := range new {
		contract.SourcePage = page
	}

	contractsToCreate := make([]*athena.MemberContract, 0, len(new))
	contractsToUpdate := make([]*athena.MemberContract, 0, len(new))
	contractsWithBids := make([]*athena.MemberContract, 0, len(new))

	oldContractMap := make(map[uint]*athena.MemberContract)
	for _, contract := range old {
		oldContractMap[contract.ContractID] = contract
	}

	for _, contract := range new {
		// This is an unknown contact, so lets flag it to be created
		if _, ok := oldContractMap[contract.ContractID]; !ok {
			contractsToCreate = append(contractsToCreate, contract)

			// We've seen this contact before for this member, lets compare it to the existing contact to see
			// if it needs to be updated
		} else if diff := deep.Equal(oldContractMap[contract.ContractID], contract); len(diff) > 0 {
			contractsToUpdate = append(contractsToUpdate, contract)
		}
		if contract.Type == athena.ContractTypeAuction {
			// This is an Auction contract, we need to
			// check to see if the bids have been updated
			contractsWithBids = append(contractsWithBids, contract)
		}
	}

	if len(contractsToCreate) > 0 {
		s.resolveContractAttributes(ctx, member, contractsToCreate)

		_, err := s.contracts.CreateContracts(ctx, member.ID, contractsToCreate)
		if err != nil {
			entry.WithError(err).Error("failed to create member contracts in DB")
			return fmt.Errorf("failed to create member contracts in DB")
		}

		for _, contract := range contractsToCreate {
			_, _ = s.FetchMemberContractItems(ctx, member, contract)
		}
	}
	if len(contractsToUpdate) > 0 {
		for _, contract := range contractsToUpdate {
			_, err := s.contracts.UpdateContract(ctx, member.ID, contract)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"contract_id": contract.ContractID,
				}).Error("failed to update member contract in DB")
				return fmt.Errorf("failed to update member contract in DB")
			}
		}
	}
	if len(contractsWithBids) > 0 {
		for _, contract := range contractsWithBids {
			_, _ = s.FetchMemberContractBids(ctx, member, contract)
		}
	}

	return nil
}

func (s *service) resolveContractAttributes(ctx context.Context, member *athena.Member, contracts []*athena.MemberContract) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "resolveContractAttributes",
	})

	unknowns := make(map[uint]bool)

	for _, contract := range contracts {

		if contract.AssigneeID.Valid && contract.AssigneeID.Uint > 0 {
			idType := glue.ResolveIDTypeFromIDRange(uint64(contract.AssigneeID.Uint))
			if idType == glue.IDTypeUnknown {
				unknowns[contract.AssigneeID.Uint] = true
			}

			switch idType {
			case glue.IDTypeCharacter:
				contract.AssigneeType.SetValid(glue.IDTypeCharacter)
				_, err := s.character.Character(ctx, contract.AssigneeID.Uint)
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":   contract.ContractID,
						"assignee_type": contract.AssigneeType,
						"assignee_id":   contract.AssigneeID,
					}).Warn("failed to resolve assignee id to name")
				}
			case glue.IDTypeCorporation:
				contract.AssigneeType.SetValid(glue.IDTypeCorporation)
				_, err := s.corporation.Corporation(ctx, contract.AssigneeID.Uint)
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":   contract.ContractID,
						"assignee_type": contract.AssigneeType,
						"assignee_id":   contract.AssigneeID,
					}).Warn("failed to resolve assignee id to name")
				}
			}
		}

		if contract.AcceptorID.Valid && contract.AcceptorID.Uint > 0 {
			idType := glue.ResolveIDTypeFromIDRange(uint64(contract.AcceptorID.Uint))
			if idType == glue.IDTypeUnknown {
				unknowns[contract.AcceptorID.Uint] = true
			}

			switch idType {
			case glue.IDTypeCharacter:
				contract.AcceptorType.SetValid(glue.IDTypeCharacter)
				_, err := s.character.Character(ctx, contract.AcceptorID.Uint)
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":   contract.ContractID,
						"acceptor_type": contract.AcceptorType,
						"acceptor_id":   contract.AcceptorID,
					}).Warn("failed to resolve acceptor id to name")
				}
			case glue.IDTypeCorporation:
				contract.AcceptorType.SetValid(glue.IDTypeCorporation)
				_, err := s.corporation.Corporation(ctx, contract.AcceptorID.Uint)
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":   contract.ContractID,
						"acceptor_type": contract.AcceptorType,
						"acceptor_id":   contract.AcceptorID,
					}).Warn("failed to resolve acceptor id to name")
				}
			}
		}

		if contract.StartLocationID.Valid {
			idType := glue.ResolveIDTypeFromIDRange(contract.StartLocationID.Uint64)

			switch idType {
			case glue.IDTypeStation:
				contract.StartLocationType.SetValid(glue.IDTypeStation)
				_, err := s.universe.Station(ctx, uint(contract.StartLocationID.Uint64))
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":         contract.ContractID,
						"start_location_type": contract.StartLocationType.String,
						"start_location_id":   contract.StartLocationID.Uint64,
					}).Warn("failed to resolve start_location id to name")
				}
			case glue.IDTypeStructure:
				contract.StartLocationType.SetValid(glue.IDTypeStation)
				_, err := s.universe.Structure(ctx, member, contract.StartLocationID.Uint64)
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":         contract.ContractID,
						"start_location_type": contract.StartLocationType.String,
						"start_location_id":   contract.StartLocationID.Uint64,
					}).Warn("failed to resolve start_location id to name")
				}
			}
		}

		if contract.EndLocationID.Valid {
			idType := glue.ResolveIDTypeFromIDRange(contract.EndLocationID.Uint64)

			switch idType {
			case glue.IDTypeStation:
				contract.EndLocationType.SetValid(glue.IDTypeStation)
				_, err := s.universe.Station(ctx, uint(contract.EndLocationID.Uint64))
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":       contract.ContractID,
						"end_location_type": contract.EndLocationType.String,
						"end_location_id":   contract.EndLocationID.Uint64,
					}).Warn("failed to resolve end_location id to name")
				}
			case glue.IDTypeStructure:
				contract.EndLocationType.SetValid(glue.IDTypeStation)
				_, err := s.universe.Structure(ctx, member, contract.EndLocationID.Uint64)
				if err != nil {
					entry.WithError(err).WithFields(logrus.Fields{
						"contract_id":       contract.ContractID,
						"end_location_type": contract.StartLocationType.String,
						"end_location_id":   contract.EndLocationID.Uint64,
					}).Warn("failed to resolve end_location id to name")
				}
			}
		}

	}

	if len(unknowns) == 0 {
		return
	}

	unknownIDs := make([]uint, 0, len(unknowns))
	for id := range unknowns {
		unknownIDs = append(unknownIDs, id)
	}

	results, _, err := s.esi.PostUniverseNames(ctx, unknownIDs)
	if err != nil {
		entry.WithError(err).Error("failed to resolve ids to names with ESI")
		return
	}

	knowns := make(map[uint]string)
	for _, result := range results {
		knowns[result.ID] = string(result.Category)
	}

	for _, contract := range contracts {
		if contract.AcceptorID.Valid && contract.AcceptorType.String == glue.IDTypeUnknown {
			if category, ok := knowns[contract.AcceptorID.Uint]; ok {
				contract.AcceptorType.SetValid(category)
			}
		}
		if contract.AssigneeID.Valid && contract.AssigneeType.String == glue.IDTypeUnknown {
			if category, ok := knowns[contract.AssigneeID.Uint]; ok {
				contract.AssigneeType.SetValid(category)
			}
		}
	}

}

func (s *service) MemberContracts(ctx context.Context, memberID, page uint) ([]*athena.MemberContract, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service":     serviceIdentifier,
		"method":      "Characters",
		"member_id":   memberID,
		"source_page": page,
	})

	contracts, err := s.cache.MemberContracts(ctx, memberID, page)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contacts from cache")
		return nil, fmt.Errorf("failed to fetch member contacts from cache")
	}

	if len(contracts) > 0 {
		return contracts, nil
	}

	contracts, err = s.contracts.MemberContracts(ctx, memberID, athena.NewEqualOperator("source_page", page))
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contacts from DB")
		return nil, fmt.Errorf("failed to fetch member contacts from DB")
	}

	if len(contracts) > 0 {
		err = s.cache.SetMemberContracts(ctx, memberID, page, contracts)
		if err != nil {
			entry.WithError(err).Error("failed to cache member contacts")

		}
	}

	return contracts, nil

}

func (s *service) FetchMemberContractItems(ctx context.Context, member *athena.Member, contract *athena.MemberContract) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContractItems, esi.ModWithCharacterID(member.ID), esi.ModWithContractID(contract.ContractID), esi.ModWithPage(1))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {
			return etag, nil
		}

		petag = etag.Etag
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id":   member.ID,
		"contract_id": contract.ContractID,
		"service":     serviceIdentifier,
		"method":      "FetchMemberContractItems",
	})

	items, etag, _, err := s.esi.GetCharacterContractItems(ctx, member.ID, contract.ContractID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch contract items from ESI")
		return nil, fmt.Errorf("failed to fetch contract items from ESI")
	}

	if petag == etag.Etag {
		return etag, nil
	}

	for _, item := range items {
		entry := entry.WithFields(logrus.Fields{
			"record_id": item.RecordID,
			"type_id":   item.TypeID,
		})

		_, err := s.universe.Type(ctx, item.TypeID)
		if err != nil {
			entry.WithError(err).Error("failed to resolve contract item to type")
			return nil, fmt.Errorf("failed to resolve contract item to type")
		}

	}

	items, err = s.contracts.CreateMemberContractItems(ctx, member.ID, contract.ContractID, items)
	if err != nil {
		entry.WithError(err).Error("failed to insert member contract items into db")
		return nil, fmt.Errorf("failed to insert member contract items into db")
	}

	if len(items) > 0 {
		err = s.cache.SetMemberContractItems(ctx, member.ID, contract.ContractID, items)
		if err != nil {
			entry.WithError(err).Error("failed to insent member contract items into db")
		}
	}

	return etag, nil

}

func (s *service) MemberContractItems(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractItem, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberContractItems",
	})

	items, err := s.cache.MemberContractItems(ctx, memberID, contractID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contract items from cache")
		return nil, fmt.Errorf("failed to fetch member contract items from cache")
	}

	if len(items) > 0 {
		return items, nil
	}

	items, err = s.contracts.MemberContractItems(ctx, memberID, contractID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contract items from db")
		return nil, fmt.Errorf("failed to fetch member contract items from db")
	}

	if len(items) > 0 {
		err = s.cache.SetMemberContractItems(ctx, memberID, contractID, items)
		if err != nil {
			entry.WithError(err).Error("failed to fetch member contract items from db")
		}
	}

	return items, err
}

func (s *service) FetchMemberContractBids(ctx context.Context, member *athena.Member, contract *athena.MemberContract) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContractBids, esi.ModWithCharacterID(member.ID), esi.ModWithContractID(contract.ContractID), esi.ModWithPage(1))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {
			return etag, nil
		}

		petag = etag.Etag
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id":   member.ID,
		"contract_id": contract.ContractID,
		"service":     serviceIdentifier,
		"method":      "FetchMemberContractBids",
	})

	bids, etag, _, err := s.esi.GetCharacterContractBids(ctx, member.ID, contract.ContractID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch contract bids from ESI")
		return nil, fmt.Errorf("failed to fetch contract bids from ESI")
	}

	if petag == etag.Etag {
		return etag, nil
	}

	for _, bid := range bids {
		_, err := s.character.Character(ctx, bid.BidderID)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"bid_id":    bid.BidID,
				"bidder_id": bid.BidderID,
			}).Error("failed to resolve contract bidder id to character")
			return nil, fmt.Errorf("failed to resolve contract bidder id  to character")
		}
	}

	bids, err = s.contracts.CreateMemberContractBids(ctx, member.ID, contract.ContractID, bids)
	if err != nil {
		entry.WithError(err).Error("failed to insert member contract bids into db")
		return nil, fmt.Errorf("failed to insert member contract bids into db")
	}

	if len(bids) > 0 {
		err = s.cache.SetMemberContractBids(ctx, member.ID, contract.ContractID, bids)
		if err != nil {
			entry.WithError(err).Error("failed to insent member contract bids into db")
		}
	}

	return etag, nil

}

func (s *service) MemberContractBids(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractBid, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberContractBids",
	})

	bids, err := s.cache.MemberContractBids(ctx, memberID, contractID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contract bids from cache")
		return nil, fmt.Errorf("failed to fetch member contract bids from cache")
	}

	if len(bids) > 0 {
		return bids, nil
	}

	bids, err = s.contracts.MemberContractBids(ctx, memberID, contractID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contract bids from db")
		return nil, fmt.Errorf("failed to fetch member contract bids from db")
	}

	if len(bids) > 0 {
		err = s.cache.SetMemberContractBids(ctx, memberID, contractID, bids)
		if err != nil {
			entry.WithError(err).Error("failed to fetch member contract bids from db")
		}
	}

	return bids, err
}
