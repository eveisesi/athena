package wallet

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
	"github.com/eveisesi/athena/internal/universe"
	"github.com/sirupsen/logrus"
)

type Service interface {
	EmptyMemberBalance(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	EmptyMembetWalletTransactions(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	EmptyMemberWalletJournals(ctx context.Context, member *athena.Member) (*athena.Etag, error)
}

type service struct {
	logger *logrus.Logger

	cache       cache.Service
	esi         esi.Service
	universe    universe.Service
	alliance    alliance.Service
	corporation corporation.Service
	character   character.Service

	wallet athena.MemberWalletRepository
}

const (
	serviceIdentifier = "Wallet Service"
)

func NewService(
	logger *logrus.Logger,

	cache cache.Service, esi esi.Service, universe universe.Service,
	alliance alliance.Service, corporation corporation.Service, character character.Service,

	wallet athena.MemberWalletRepository,
) Service {
	return &service{
		logger:      logger,
		cache:       cache,
		esi:         esi,
		universe:    universe,
		alliance:    alliance,
		corporation: corporation,
		character:   character,
		wallet:      wallet,
	}
}

func (s *service) EmptyMemberBalance(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterWalletBalance, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	return s.FetchMemberBalance(ctx, member, etag)

}

func (s *service) MemberBalance(ctx context.Context, member *athena.Member) (*athena.MemberWalletBalance, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "MemberBalance",
	})

	balance, err := s.cache.MemberWalletBalance(ctx, member)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member wallet balance from cache")
		return nil, fmt.Errorf("failed to fetch member wallet balance from cache")
	}

	if balance == nil {
		balance, err = s.wallet.MemberWalletBalance(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			entry.WithError(err).Error("failed to fetch member wallet balance from db")
			return nil, fmt.Errorf("failed to fetch member wallet balance from db")
		}

		if balance == nil || errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no balance could be found for the specified character")
		}

		err = s.cache.SetMemberWalletBalance(ctx, member, balance)
		if err != nil {
			entry.WithError(err).Error("failed to cache member wallet balance")
		}

	}

	return balance, nil

}

func (s *service) FetchMemberBalance(ctx context.Context, member *athena.Member, etag *athena.Etag) (*athena.Etag, error) {

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberBalance",
	})

	rawBalance, etag, _, err := s.esi.GetCharacterWalletBalance(ctx, member)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member balance from ESI")
		return nil, fmt.Errorf("failed to fetch member balance from ESI")
	}

	balance, err := s.wallet.MemberWalletBalance(ctx, member.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member balance from DB")
		return nil, fmt.Errorf("failed to fetch member balance from DB")
	}

	exists := true

	if balance == nil || errors.Is(err, sql.ErrNoRows) {
		exists = false
	}

	switch exists {
	case true:
		balance, err = s.wallet.UpdateMemberWalletBalance(ctx, member.ID, rawBalance)
		if err != nil {
			entry.WithError(err).Error("failed to update member wallet balance in database")
			return nil, fmt.Errorf("failed to update member wallet balance in database")
		}
	case false:
		balance, err = s.wallet.CreateMemberWalletBalance(ctx, member.ID, rawBalance)
		if err != nil {
			entry.WithError(err).Error("failed to create member wallet balance in database")
			return nil, fmt.Errorf("failed to create member wallet balance in database")
		}
	}

	err = s.cache.SetMemberWalletBalance(ctx, member, balance)
	if err != nil {
		entry.WithError(err).Error("failed to create member wallet balance in database")
	}

	return etag, err

}

func (s *service) EmptyMembetWalletTransactions(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterWalletTransactions, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	return s.FetchMemberWalletTransaction(ctx, member, etag)

}

func (s *service) FetchMemberWalletTransaction(ctx context.Context, member *athena.Member, etag *athena.Etag) (*athena.Etag, error) {

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberWalletTransaction",
	})

	transactions, etag, _, err := s.esi.GetCharacterWalletTransactions(ctx, member, make([]*athena.MemberWalletTransaction, 0))
	if err != nil {
		entry.WithError(err).Error("failed to fetch member wallet transactions from ESI")
		return nil, fmt.Errorf("failed to fetch member wallet transactions from ESI")
	}

	if len(transactions) > 0 {

		s.resolveMemberWalletTransactionAttributes(ctx, member, transactions)
		_, err = s.wallet.CreateMemberWalletTransactions(ctx, member.ID, transactions)
		if err != nil {
			entry.WithError(err).Error("failed to create transaction in db")
			return nil, fmt.Errorf("failed to create transaction in db")
		}

	}

	return etag, err
}

func (s *service) resolveMemberWalletTransactionAttributes(ctx context.Context, member *athena.Member, transactions []*athena.MemberWalletTransaction) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "resolveMemberWalletTransactionAttributes",
	})

	unknowns := make(map[uint]bool)

	for _, transaction := range transactions {

		clientType := resolveIDTypeFromIDRange(uint64(transaction.ClientID))
		if clientType == IDTypeUnknown {
			if _, ok := unknowns[transaction.ClientID]; !ok {
				unknowns[transaction.ClientID] = true
			}
		}

		switch clientType {
		case IDTypeCharacter:
			transaction.ClientType = athena.ClientTypeCharacter
			_, err := s.character.Character(ctx, transaction.ClientID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"transaction_id": transaction.TransactionID,
					"client_type":    transaction.ClientType.String(),
					"client_id":      transaction.ClientID,
				}).Warn("failed to resolve client id to name")
			}
		case IDTypeCorporation:
			transaction.ClientType = athena.ClientTypeCorporation
			_, err := s.corporation.Corporation(ctx, transaction.ClientID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"transaction_id": transaction.TransactionID,
					"client_type":    transaction.ClientType.String(),
					"client_id":      transaction.ClientID,
				}).Warn("failed to resolve client id to name")
			}
		default:
			transaction.ClientType = IDTypeUnknown
		}

		locationType := resolveIDTypeFromIDRange(transaction.LocationID)

		switch locationType {
		case IDTypeStation:
			transaction.LocationType = athena.LocationTypeStation
			_, err := s.universe.Station(ctx, uint(transaction.LocationID))
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"transaction_id": transaction.TransactionID,
					"location_type":  transaction.LocationType.String(),
					"location_id":    transaction.LocationID,
				}).Warn("failed to resolve location id to name")
			}
		case IDTypeStructure:
			transaction.LocationType = athena.LocationTypeStructure
			_, err := s.universe.Structure(ctx, member, transaction.LocationID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"transaction_id": transaction.TransactionID,
					"location_type":  transaction.LocationType.String(),
					"location_id":    transaction.LocationID,
				}).Warn("failed to resolve location id to name")
			}
		default:
			transaction.LocationType = athena.LocationTypeUnknown
			entry.WithFields(logrus.Fields{
				"transaction_id": transaction.TransactionID,
				"type":           transaction.LocationType.String(),
			}).Error("unknown location type encountered")
		}

		_, err := s.universe.Type(ctx, transaction.TypeID)
		if err != nil {
			entry.WithError(err).WithFields(logrus.Fields{
				"transaction_id": transaction.TransactionID,
				"type_id":        transaction.TypeID,
			}).Warn("failed to resolve type id to name")
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

	for _, transaction := range transactions {
		switch transaction.ClientType {
		case IDTypeUnknown:
			if category, ok := knowns[transaction.ClientID]; ok {
				transaction.ClientType = athena.ClientType(category)
			}
		}
	}

}

func (s *service) EmptyMemberWalletJournals(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterWalletJournal, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	return s.FetchMemberWalletJournal(ctx, member, etag)

}

func (s *service) FetchMemberWalletJournal(ctx context.Context, member *athena.Member, etag *athena.Etag) (*athena.Etag, error) {

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberWalletJournal",
	})

	entries, etag, _, err := s.esi.GetCharacterWalletJournals(ctx, member, make([]*athena.MemberWalletJournal, 0))
	if err != nil {
		entry.WithError(err).Error("failed to fetch member wallet journals from ESI")
		return nil, fmt.Errorf("failed to fetch member wallet journals from ESI")
	}

	s.logger.WithField("count", len(entries)).Info("num of entries received")

	if len(entries) > 0 {
		s.resolveMemberWalletJournalEntries(ctx, member, entries)
		_, err = s.wallet.CreateMemberWalletJournals(ctx, member.ID, entries)
		if err != nil {
			entry.WithError(err).Error("failed to create entries in db")
			return nil, fmt.Errorf("failed to create entries in db")
		}
	}

	return etag, err

}

func (s *service) resolveMemberWalletJournalEntries(ctx context.Context, member *athena.Member, entries []*athena.MemberWalletJournal) {

	for _, entry := range entries {
		if entry.ContextID.Valid && entry.ContextType.Valid {
			s.resolveContextIDType(ctx, member, entry.ContextID.Uint64, entry.ContextType.ContextIDType)
		}
	}

}

func (s *service) resolveContextIDType(ctx context.Context, member *athena.Member, id uint64, idtype athena.ContextIDType) {
	entry := s.logger.WithFields(logrus.Fields{
		"context_id":   id,
		"context_type": idtype,
	})
	switch idtype {
	case athena.ContextIDTypeStructureID:
		_, err := s.universe.Structure(ctx, member, id)
		if err != nil {
			entry.WithError(err).Error("failed to resolve id")
		}

	case athena.ContextIDTypeStationID:
		_, err := s.universe.Station(ctx, uint(id))
		if err != nil {
			entry.WithError(err).Error("failed to resolve id")
		}

		// case athena.ContextIDTypeMarketTransactionID:
	case athena.ContextIDTypeCharacterID:
		_, err := s.character.Character(ctx, uint(id))
		if err != nil {
			entry.WithError(err).Error("failed to resolve id")
		}

	case athena.ContextIDTypeCorporationID:
		_, err := s.corporation.Corporation(ctx, uint(id))
		if err != nil {
			entry.WithError(err).Error("failed to resolve id")
		}
	case athena.ContextIDTypeAllianceID:
		_, err := s.alliance.Alliance(ctx, uint(id))
		if err != nil {
			entry.WithError(err).Error("failed to resolve id")
		}
	case athena.ContextIDTypeSystemID:
		_, err := s.universe.SolarSystem(ctx, uint(id))
		if err != nil {
			entry.WithError(err).Error("failed to resolve id")
		}
	case athena.ContextIDTypeTypeID:
		_, err := s.universe.Type(ctx, uint(id))
		if err != nil {
			entry.WithError(err).Error("failed to resolve id")
		}
	}
}

const (
	IDTypeCharacter   = "character"
	IDTypeCorporation = "corporation"
	IDTypeAlliance    = "alliance"
	IDTypeStation     = "station"
	IDTypeStructure   = "structure"
	IDTypeUnknown     = "unknown"
)

func resolveIDTypeFromIDRange(id uint64) string {

	switch d := id; {
	case d >= 60000000 && d < 64000000:
		return IDTypeStation
	case d >= 90000000 && d < 98000000:
		return IDTypeCharacter
	case d >= 98000000 && d < 99000000:
		return IDTypeCorporation
	case d >= 99000000 && d < 100000000:
		return IDTypeAlliance
	case d >= 100000000 && d < 2100000000:
		return IDTypeUnknown
	case d >= 2100000000 && d < 1000000000000:
		return IDTypeCharacter
	case d >= 1000000000000:
		return IDTypeStructure
	default:
		return IDTypeUnknown
	}

}
