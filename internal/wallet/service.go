package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/glue"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/sirupsen/logrus"
)

type Service interface {
	// Fetch Member Balance fetches the provided characters balance from ESI and stores it in the repository
	FetchMemberBalance(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberBalance(ctx context.Context, member *athena.Member) (*athena.MemberWalletBalance, error)

	// Fetch Member Wallet Transactions fetches the provided characters transactions from ESI and stores them in the repository
	FetchMemberWalletTransactions(ctx context.Context, member *athena.Member) (*athena.Etag, error)

	// Fetch Member Wallet Journals fetches the provided characters jounral entries from ESI and stores them in the repository
	FetchMemberWalletJournals(ctx context.Context, member *athena.Member) (*athena.Etag, error)
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

func (s *service) FetchMemberBalance(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterWalletBalance, esi.ModWithCharacterID(member.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberBalance",
	})

	rawBalance, etag, _, err := s.esi.GetCharacterWalletBalance(ctx, member.ID, member.AccessToken.String)
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

func (s *service) FetchMemberWalletTransactions(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterWalletTransactions, esi.ModWithCharacterID(member.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberWalletTransaction",
	})

	_, _, err = s.esi.HeadCharacterWalletTransactions(ctx, member.ID, 0, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to exec head request for member wallet transactions from ESI")
		return nil, fmt.Errorf("failed to exec head request for member wallet transactions from ESI")
	}

	from := uint64(0)
	for {

		entry := entry.WithField("from", from)

		ptransactions, _, _, err := s.esi.GetCharacterWalletTransactions(ctx, member.ID, from, member.AccessToken.String)
		if err != nil {
			entry.WithError(err).Error("failed to fetch member wallet transactions from ESI")
			return nil, fmt.Errorf("failed to fetch member wallet transactions from ESI")
		}

		if len(ptransactions) > 0 {

			s.resolveMemberWalletTransactionAttributes(ctx, member, ptransactions)
			_, err = s.wallet.CreateMemberWalletTransactions(ctx, member.ID, ptransactions)
			if err != nil {
				entry.WithError(err).Error("failed to create transaction in db")
				return nil, fmt.Errorf("failed to create transaction in db")
			}

		}

		lptrans := len(ptransactions)

		if lptrans < 2500 {
			break
		}

		from = ptransactions[lptrans-1].TransactionID - 1

	}

	etag, err = s.esi.Etag(ctx, esi.GetCharacterWalletTransactions, esi.ModWithCharacterID(member.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
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

		clientType := glue.ResolveIDTypeFromIDRange(uint64(transaction.ClientID))
		if clientType == glue.IDTypeUnknown {
			if _, ok := unknowns[transaction.ClientID]; !ok {
				unknowns[transaction.ClientID] = true
			}
		}

		switch clientType {
		case glue.IDTypeCharacter:
			transaction.ClientType = athena.ClientTypeCharacter
			_, err := s.character.Character(ctx, transaction.ClientID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"transaction_id": transaction.TransactionID,
					"client_type":    transaction.ClientType.String(),
					"client_id":      transaction.ClientID,
				}).Warn("failed to resolve client id to name")
			}
		case glue.IDTypeCorporation:
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
			transaction.ClientType = glue.IDTypeUnknown
		}

		locationType := glue.ResolveIDTypeFromIDRange(transaction.LocationID)

		switch locationType {
		case glue.IDTypeStation:
			transaction.LocationType = athena.LocationTypeStation
			_, err := s.universe.Station(ctx, uint(transaction.LocationID))
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"transaction_id": transaction.TransactionID,
					"location_type":  transaction.LocationType.String(),
					"location_id":    transaction.LocationID,
				}).Warn("failed to resolve location id to name")
			}
		case glue.IDTypeStructure:
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
		case glue.IDTypeUnknown:
			if category, ok := knowns[transaction.ClientID]; ok {
				transaction.ClientType = athena.ClientType(category)
			}
		}
	}

}

func (s *service) FetchMemberWalletJournals(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterWalletJournal, esi.ModWithCharacterID(member.ID))
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

	etag, res, err := s.esi.HeadCharacterWalletJournals(ctx, member.ID, 1, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to exec head request for member wallet journals from ESI")
		return nil, fmt.Errorf("failed to exec head request for member wallet journals from ESI")
	}

	pages := esi.RetrieveXPagesFromHeader(res.Header)
	s.logger.WithField("pages", pages).Info("num of pages")

	for page := uint(1); page <= pages; page++ {
		entry := entry.WithField("page", page)

		entries, _, _, err := s.esi.GetCharacterWalletJournals(ctx, member.ID, page, member.AccessToken.String)
		if err != nil {
			entry.WithError(err).Error("failed to fetch member wallet journals from ESI")
			return nil, fmt.Errorf("failed to fetch member wallet journals from ESI")
		}

		if len(entries) > 0 {
			s.resolveMemberWalletJournalEntries(ctx, member, entries)
			_, err = s.wallet.CreateMemberWalletJournals(ctx, member.ID, entries)
			if err != nil {
				entry.WithError(err).Error("failed to create entries in db")
				return nil, fmt.Errorf("failed to create entries in db")
			}
		}

		time.Sleep(time.Millisecond * 100)

	}

	return etag, err

}

func (s *service) resolveMemberWalletJournalEntries(ctx context.Context, member *athena.Member, entries []*athena.MemberWalletJournal) {
	logEntry := s.logger.WithFields(logrus.Fields{
		"service":   serviceIdentifier,
		"method":    "resolveMemberWalletJournalEntries",
		"member_id": member.ID,
	})

	for _, entry := range entries {

		logEntry := logEntry.WithField("journal_id", entry.JournalID)

		if entry.ContextID.Valid && entry.ContextType.Valid {
			s.resolveContextIDType(ctx, member, entry.ContextID.Uint64, entry.ContextType.ContextIDType)
		}

		if entry.FirstPartyID.Valid && entry.SecondPartyID.Valid {
			s.resolvePartyID(ctx, entry)
		}

		if entry.TaxReceiverID.Valid {
			_, err := s.corporation.Corporation(ctx, entry.TaxReceiverID.Uint)
			if err != nil {
				logEntry.WithError(err).Error("failed to resolve corporation id on journal for tax receiver id")
			}
		}

	}

}

func (s *service) resolvePartyID(ctx context.Context, journal *athena.MemberWalletJournal) {
	entry := s.logger.WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "resolvePartyID",
	})

	results, res, err := s.esi.PostUniverseNames(ctx, []uint{journal.FirstPartyID.Uint, journal.SecondPartyID.Uint})
	if err != nil && res == nil {
		entry.WithError(err).Error("failed to execute request to post universe names")
		return
	}

	if res.StatusCode == http.StatusBadRequest {
		entry.Error("request to post universe names failed to with status code 400")
		return
	}

	if len(results) != 2 {
		entry.Error("request to post universe names returned unexpected number of results")
	}

	for _, result := range results {
		if result.ID == journal.FirstPartyID.Uint {
			journal.FirstPartyType.SetValid(string(result.Category))
		}
		if result.ID == journal.SecondPartyID.Uint {
			journal.SecondPartyType.SetValid(string(result.Category))
		}
	}

}

func (s *service) resolveContextIDType(ctx context.Context, member *athena.Member, id uint64, idtype athena.ContextIDType) {
	entry := s.logger.WithFields(logrus.Fields{
		"service":      serviceIdentifier,
		"method":       "resolveContextIDType",
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
