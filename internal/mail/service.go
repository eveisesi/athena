package mail

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
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

type Service interface {
	EmptyMemberMailHeaders(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	EmptyMemberMailingLists(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	// EmptyMemberMailLabels(ctx context.Context, member *athena.Member) (*athena.Etag, error)
}

type service struct {
	logger *logrus.Logger

	cache cache.Service
	esi   esi.Service

	alliance    alliance.Service
	corporation corporation.Service
	character   character.Service

	mail athena.MailRepository
}

const (
	serviceIdentifier = "Mail Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, character character.Service, alliance alliance.Service, corporation corporation.Service, mail athena.MailRepository) Service {
	return &service{
		logger: logger,

		cache: cache,
		esi:   esi,

		alliance:    alliance,
		corporation: corporation,
		character:   character,

		mail: mail,
	}
}

func (s *service) EmptyMemberMailHeaders(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterMailHeaders, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	return s.FetchMemberMailHeaders(ctx, member, etag)

}

func (s *service) FetchMemberMailHeaders(ctx context.Context, member *athena.Member, etag *athena.Etag) (*athena.Etag, error) {

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberMailHeaders",
	})

	// We need to make sure that mailing lists are update to date before processing headers
	// We can call EmptyMemberMailingLists to accomplish this
	_, err := s.EmptyMemberMailingLists(ctx, member)
	if err != nil {
		entry.WithError(err).Error("failed to prefetch member mailing lists")
		return etag, fmt.Errorf("failed to prefetch member mailing lists")
	}

	newHeaders, etag, _, err := s.esi.GetCharacterMailHeaders(ctx, member)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member mail headers from ESI")
		return nil, fmt.Errorf("failed to fetch member mail headers from ESI")
	}

	if len(newHeaders) == 0 {
		return etag, nil
	}

	memberHeaders, err := s.mail.MemberMailHeaders(ctx, athena.NewEqualOperator("member_id", member.ID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member mail headers from DB")
		return nil, fmt.Errorf("failed to fetch member mail headers from DB")
	}
	// processHeaders takes care of creating/updating the headers in the db for us.
	s.processHeaders(ctx, member, memberHeaders, newHeaders)

	return etag, err

}

func (s *service) processHeaders(ctx context.Context, member *athena.Member, old []*athena.MemberMailHeader, newHeader []*esi.MailHeader) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "processHeaders",
	})

	var oldMap = make(map[uint]*athena.MemberMailHeader)
	for _, header := range old {
		oldMap[header.MailID] = header
	}

	var headersToCreate = make([]*esi.MailHeader, 0, len(newHeader))
	var headersToUpdate = make([]*esi.MailHeader, 0, len(newHeader))

	for _, header := range newHeader {
		if _, ok := oldMap[header.MailID.Uint]; !ok {
			headersToCreate = append(headersToCreate, header)
		} else {
			oldHeader := oldMap[header.MailID.Uint]
			if oldHeader.IsRead != header.IsRead {
				headersToUpdate = append(headersToUpdate, header)
			}
			if SliceUint64Equal([]uint64(oldHeader.Labels), header.Labels) {
				headersToUpdate = append(headersToUpdate, header)
			}
		}
	}

	if len(headersToCreate) > 0 {
		mailHeaders := make([]*athena.MailHeader, 0, len(headersToCreate))
		// The multiple of 52 in the capacity is the max number of recipients per mail
		// according to ESI Specification
		mailRecipients := make([]*athena.MailRecipient, 0, len(headersToCreate)*52)
		memberMailHeaders := make([]*athena.MemberMailHeader, 0, len(headersToCreate))

		// Recipients are an immutable property, so we only need
		// to process them on new headers
		for _, header := range headersToCreate {

			entry := entry.WithField("header_id", header.MailID.Uint)

			// Fetch the body of the header
			singleMailHeader, _, _, err := s.esi.GetCharacterMailHeader(ctx, member, header)
			if err != nil {
				entry.WithError(err).Error("failed to fetch indiviudal header from ESI")
				continue
			}

			// According to the API Spec, it is possible for the header to not return who, or what
			// the sender id is, so check to see if it is valid
			if header.From.Valid {
				// ESI Docs has a list of ID Ranges, this function switches
				// over a simplified version of that range
				senderType := glue.ResolveIDTypeFromIDRange(uint64(header.From.Uint))
				entry := entry.WithFields(logrus.Fields{
					"sender_type": senderType,
				})

				mailHeader := &athena.MailHeader{
					MailID:    header.MailID.Uint,
					Sender:    header.From,
					Body:      singleMailHeader.Body,
					Subject:   header.Subject,
					Timestamp: header.Timestamp,
				}

			SenderSwitch:
				switch senderType {
				case glue.IDTypeCorporation:
					mailHeader.SenderType.SetValid(glue.IDTypeCorporation)
					_, err := s.corporation.Corporation(ctx, header.From.Uint)
					if err != nil {
						entry.WithField("corporation_id", header.From.Uint).WithError(err).Error("failed to resolve corporation id to corporation")
					}
				case glue.IDTypeCharacter:

					mailHeader.SenderType.SetValid(glue.IDTypeCharacter)
					_, err := s.character.Character(ctx, header.From.Uint)
					if err != nil {
						entry.WithField("character_id", header.From.Uint).WithError(err).Error("failed to resolve character id to character id")
					}
				case glue.IDTypeUnknown:

					entry.Info("attempting to resolve unknown sender type")

					// Since we can't determine the ID Type, lets check to see if it is a mailing list id
					list, err := s.MailingList(ctx, header.From.Uint)
					if list != nil && err == nil {
						mailHeader.SenderType.SetValid("mailing_list")
						entry.Info("sender is a mailing list")
						break SenderSwitch
					}
					// Ok so this is not a mailing list that we know about, so now we need to POST the ID to /universe/names. If that 404's this is a mailing list ID that we haven't discovered yet.
					results, res, err := s.esi.PostUniverseNames(ctx, []uint{header.From.Uint})
					// If res is nil and err is not nil
					// Then there was an error executing the request
					// Log an error and break out of the switch
					if err != nil && res == nil {
						entry.WithError(err).Error("failed to execute request to post universe names")
						break SenderSwitch
					}

					// This is one of the few, if not the only please
					// that we give the service knowledge of the status code
					// that was returned from the result. Im not sure how
					// else to handle this
					if res.StatusCode == http.StatusBadRequest {
						// If the request was a bad request, this is PROBABLY
						// a mailing list that we don't know about, so just
						// list the mailHeader.SenderType invalid on the struct
						// It will be inserted as null
						entry.Error("request to post universe names failed with status code 400")
						break SenderSwitch
					}

					var result = new(esi.PostUniverseNamesOK)
					if len(results) == 1 {
						result = results[0]
					}
					// ResultCategorySwitch:
					switch cat := result.Category; {
					case cat == esi.CategoryAlliance:
						mailHeader.SenderType.SetValid(string(cat))
						_, err := s.alliance.Alliance(ctx, result.ID)
						if err != nil {
							entry.WithError(err).Error("failed to resolve alliance id")
							break SenderSwitch
						}
					case cat == esi.CategoryCorporation:
						mailHeader.SenderType.SetValid(string(cat))
						_, err := s.corporation.Corporation(ctx, result.ID)
						if err != nil {
							entry.WithError(err).Error("failed to resolve corporation id")
							break SenderSwitch
						}
					case cat == esi.CategoryCharacter:
						mailHeader.SenderType.SetValid(string(cat))
						_, err := s.character.Character(ctx, result.ID)
						if err != nil {
							entry.WithError(err).Error("failed to resolve corporation id")
							break SenderSwitch
						}
					default:
						mailHeader.SenderType.SetValid("unknown")
					}
				}

				// Now we need to resolve the recipients
				for _, recipient := range header.Recipients {
					if !recipient.RecipientType.IsValid() {
						// Im not sure how this would happen, but Im testing for it anyway
						continue
					}

					switch recipient.RecipientType {
					case athena.RecipientTypeAlliance:
						_, err := s.alliance.Alliance(ctx, recipient.RecipientID)
						if err != nil {
							entry.WithError(err).Error("failed to resolve alliance id")
						}
					case athena.RecipientTypeCharacter:
						_, err := s.character.Character(ctx, recipient.RecipientID)
						if err != nil {
							entry.WithError(err).Error("failed to resolve character id")
						}
					case athena.RecipientTypeCorporation:
						_, err := s.corporation.Corporation(ctx, recipient.RecipientID)
						if err != nil {
							entry.WithError(err).Error("failed to resolve corporation id")
						}
					}

					mailRecipients = append(mailRecipients, &athena.MailRecipient{
						MailID:        mailHeader.MailID,
						RecipientID:   recipient.RecipientID,
						RecipientType: recipient.RecipientType,
					})

				}

				mailHeaders = append(mailHeaders, mailHeader)
				memberMailHeaders = append(memberMailHeaders, &athena.MemberMailHeader{
					MemberID: member.ID,
					MailID:   mailHeader.MailID,
					IsRead:   header.IsRead,
					Labels:   athena.SliceUint(header.Labels),
				})

				// We're potentially making a lot of requests to ESI in this loop
				//  (assuming some of this shit isn't cached in redis already and
				// we don't already have it in our db)
				// Lets be nice and sleep for a few hunderd milliseconds
				time.Sleep(time.Millisecond * 200)

			}

		}

		// Out Slice of mailHeaders and memberMailHeaders should be equal, lets test for that here
		if len(mailHeaders) != len(memberMailHeaders) {
			// We aren't going to exit anything, but we do need to log an error
			entry.WithFields(logrus.Fields{
				"lenMailHeaders":       len(mailHeaders),
				"lenMemberMailHeaders": len(memberMailHeaders),
			}).Error("Mail Headers and Member Mail Header are not of equal length")
		}

		_, err := s.mail.CreateMailHeaders(ctx, mailHeaders)
		if err != nil {
			entry.WithError(err).Error("failed to create mail headers in DB")
			return
		}

		_, err = s.mail.CreateMailRecipients(ctx, mailRecipients)
		if err != nil {
			entry.WithError(err).Error("failed to create mail header recipients in DB")
			return
		}

		_, err = s.mail.CreateMemberMailHeaders(ctx, member.ID, memberMailHeaders)
		if err != nil {
			entry.WithError(err).Error("failed to create member mail headers in DB")
			return
		}

	}
	if len(headersToUpdate) > 0 {
		memberMailHeaders := make([]*athena.MemberMailHeader, 0, len(headersToUpdate))
		for _, header := range headersToUpdate {
			memberMailHeaders = append(memberMailHeaders, &athena.MemberMailHeader{
				MemberID: member.ID,
				MailID:   header.MailID.Uint,
				IsRead:   header.IsRead,
				Labels:   athena.SliceUint(header.Labels),
			})
		}

		_, err := s.mail.UpdateMemberMailHeaders(ctx, member.ID, memberMailHeaders)
		if err != nil {
			entry.WithError(err).Error("failed to update member mail header in DB")
		}

	}

}

func SliceUint64Equal(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func (s *service) EmptyMemberMailingLists(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterMailLists, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	return s.FetchMemberMailingLists(ctx, member, etag)

}

func (s *service) MailingList(ctx context.Context, mailingListID uint) (*athena.MailingList, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"mailing_list_id": mailingListID,
		"service":         serviceIdentifier,
		"method":          "MailingList",
	})

	// Check the cache first
	list, err := s.cache.MailingList(ctx, mailingListID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch mailing list from cache")
		return nil, fmt.Errorf("failed to fetch mailing list from cache")
	}

	if list != nil {
		return list, nil
	}

	list, err = s.mail.MailingList(ctx, mailingListID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch mailing list from DB")
		return nil, fmt.Errorf("failed to fetch mailing list from DB")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return list, nil

}

func (s *service) FetchMemberMailingLists(ctx context.Context, member *athena.Member, etag *athena.Etag) (*athena.Etag, error) {

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberMailingLists",
	})

	lists, etag, _, err := s.esi.GetCharacterMailLists(ctx, member, make([]*athena.MailingList, 0))
	if err != nil {
		entry.WithError(err).Error("failed to fetch member mailing lists from ESI")
		return nil, fmt.Errorf("failed to fetch member mailing lists from ESI")
	}

	if len(lists) > 0 {
		// Fetch the lists we know about
		memberMailingLists, err := s.mail.MemberMailingLists(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			entry.WithError(err).Error("failed to members existing mailing lists from db")
			return nil, fmt.Errorf("failed to members existing mailing lists from db")
		}

		mailingListIDs := make([]interface{}, 0, len(memberMailingLists))
		for _, list := range memberMailingLists {
			mailingListIDs = append(mailingListIDs, list.MailingListID)
		}

		mailingLists, err := s.mail.MailingLists(ctx, athena.NewInOperator("mailing_list_id", mailingListIDs...))
		if err != nil {
			entry.WithError(err).Error("failed to resolve member mailing list ids to mailing lists")
			return nil, fmt.Errorf("failed to resolve member mailing list ids to mailing lists")
		}

		s.processMailingLists(ctx, member, mailingLists, lists)

	}

	return etag, nil

}

func (s *service) processMailingLists(ctx context.Context, member *athena.Member, old []*athena.MailingList, new []*athena.MailingList) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service":   serviceIdentifier,
		"method":    "processMailingLists",
		"member_id": member.ID,
	})

	oldMap := make(map[uint]*athena.MailingList)
	for _, list := range old {
		oldMap[list.MailingListID] = list
	}

	listsToCreate := make([]*athena.MailingList, 0, len(new))
	listsToUpdate := make([]*athena.MailingList, 0, len(new))
	for _, list := range new {
		if _, ok := oldMap[list.MailingListID]; !ok {
			listsToCreate = append(listsToCreate, list)
		} else if diff := deep.Equal(oldMap[list.MailingListID], list); len(diff) > 0 {
			listsToUpdate = append(listsToUpdate, list)
		}
	}

	if len(listsToCreate) > 0 {

		createdLists, err := s.mail.CreateMailingLists(ctx, listsToCreate)
		if err != nil {
			entry.WithError(err).Error("Failed to create mailing lists")
			return
		}

		memberLists := make([]*athena.MemberMailingList, 0, len(createdLists))
		for _, list := range createdLists {
			memberLists = append(memberLists, &athena.MemberMailingList{
				MemberID:      member.ID,
				MailingListID: list.MailingListID,
			})
		}

		if len(memberLists) > 0 {
			_, err = s.mail.CreateMemberMailingLists(ctx, member.ID, memberLists)
			if err != nil {
				entry.WithError(err).Error("Failed to create member mailing lists")
				return
			}
		}

	}

	if len(listsToUpdate) > 0 {
		for _, list := range listsToUpdate {
			entry := entry.WithField("mailing_list_id", list.MailingListID)

			_, err := s.mail.UpdateMailingList(ctx, list.MailingListID, list)
			if err != nil {
				entry.WithError(err).Error("failed to update mailing list")
			}
		}
	}

}
