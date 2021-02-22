package contact

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
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

type Service interface {
	EmptyMemberContacts(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberContacts(ctx context.Context, memberID, page uint) ([]*athena.MemberContact, error)
	EmptyMemberContactLabels(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error)
}

type service struct {
	logger *logrus.Logger

	cache       cache.Service
	esi         esi.Service
	alliance    alliance.Service
	character   character.Service
	corporation corporation.Service
	universe    universe.Service

	contacts athena.MemberContactRepository
}

const (
	serviceIdentifier = "Contact Service"
)

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, universe universe.Service, alliance alliance.Service, character character.Service, corporation corporation.Service, contacts athena.MemberContactRepository) Service {

	return &service{
		logger: logger,

		cache:       cache,
		esi:         esi,
		universe:    universe,
		alliance:    alliance,
		character:   character,
		corporation: corporation,

		contacts: contacts,
	}

}

func (s *service) EmptyMemberContacts(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContacts, esi.ModWithCharacterID(member.ID), esi.ModWithPage(1))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "FetchMemberContacts",
	})

	_, res, err := s.esi.HeadCharacterContacts(ctx, member.ID, 1, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to exec head request for member contacts from ESI")
		return nil, fmt.Errorf("failed to exec head request for member contacts from ESI")
	}

	pages := esi.RetrieveXPagesFromHeader(res.Header)

	for page := uint(1); page <= pages; page++ {
		entry := entry.WithField("page", page)

		etag, err := s.esi.Etag(ctx, esi.GetCharacterContacts, esi.ModWithCharacterID(member.ID), esi.ModWithPage(page))
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

		contacts, etag, _, err := s.esi.GetCharacterContacts(ctx, member.ID, page, member.AccessToken.String)
		if err != nil {
			entry.WithError(err).Error("failed to fetch member contacts from ESI")
			return nil, fmt.Errorf("failed to fetch member contacts from ESI")
		}

		if petag != "" && petag == etag.Etag {
			continue
		}

		existingContacts, err := s.contacts.MemberContacts(ctx, member.ID, athena.NewEqualOperator("source_page", page))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		// Need to add page or sourcePage property to Struct so that the ESI Page that the record was discovered on can be tracked.
		s.resolveContactAttributes(ctx, contacts)
		contacts, err = s.diffAndUpdateContacts(ctx, member, page, existingContacts, contacts)
		if err != nil {
			return nil, fmt.Errorf("failed to diff and update contacts")
		}

		if len(contacts) > 0 {
			err := s.cache.SetMemberContacts(ctx, member.ID, page, contacts)
			if err != nil {
				entry.WithError(err).Error("failed to cache member contacts")
			}
		}

	}

	return etag, err

}

func (s *service) MemberContacts(ctx context.Context, memberID, page uint) ([]*athena.MemberContact, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberContacts",
	})

	contacts, err := s.cache.MemberContacts(ctx, memberID, page)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contacts from cache")
		return nil, fmt.Errorf("failed to fetch member contacts from cache")
	}

	if len(contacts) > 0 {
		return contacts, nil
	}

	ops := make([]*athena.Operator, 0, 1)
	if page > 0 {
		ops = append(ops, athena.NewEqualOperator("source_page", page))
	}

	contacts, err = s.contacts.MemberContacts(ctx, memberID, ops...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member contacts from DB")
		return nil, fmt.Errorf("failed to fetch member contacts from DB")
	}

	if len(contacts) > 0 {
		err = s.cache.SetMemberContacts(ctx, memberID, page, contacts)
		if err != nil {
			entry.WithError(err).Error("failed to cache member contacts")
		}
	}

	return contacts, nil

}

func (s *service) diffAndUpdateContacts(ctx context.Context, member *athena.Member, page uint, old []*athena.MemberContact, new []*athena.MemberContact) ([]*athena.MemberContact, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "diffAndUpdateContacts",
	})

	// Apply Page to all structs in new array
	for _, contact := range new {
		contact.SourcePage = page
	}

	contactsToCreate := make([]*athena.MemberContact, 0)
	contactsToUpdate := make([]*athena.MemberContact, 0)
	contactsToDelete := make([]*athena.MemberContact, 0)

	oldContactMap := make(map[uint]*athena.MemberContact)
	for _, contact := range old {
		oldContactMap[contact.ContactID] = contact
	}

	for _, contact := range new {
		// This is an unknown contact, so lets flag it to be created
		if _, ok := oldContactMap[contact.ContactID]; !ok {
			contactsToCreate = append(contactsToCreate, contact)

			// We've seen this contact before for this member, lets compare it to the existing contact to see
			// if it needs to be updated
		} else if diff := deep.Equal(oldContactMap[contact.ContactID], contact); len(diff) > 0 {
			contactsToUpdate = append(contactsToUpdate, contact)
		}
	}

	newContactMap := make(map[uint]*athena.MemberContact)
	for _, contact := range new {
		newContactMap[contact.ContactID] = contact
	}

	for _, contact := range old {
		// This label is not in the list of new label, must've been deleted by the user in game
		if _, ok := newContactMap[contact.ContactID]; !ok {
			contactsToDelete = append(contactsToDelete, contact)
		}
	}

	if len(contactsToDelete) > 0 {
		_, err := s.contacts.DeleteMemberContacts(ctx, member.ID, contactsToDelete)
		if err != nil {
			entry.WithError(err).Error("failed to delete member contacts in the database")
			return nil, fmt.Errorf("failed to delete member contacts in the database")
		}

	}

	var final = make([]*athena.MemberContact, 0)
	if len(contactsToCreate) > 0 {
		createdContacts, err := s.contacts.CreateMemberContacts(ctx, member.ID, contactsToCreate)
		if err != nil {
			entry.WithError(err).Error("failed to create member contacts in the database")
			return nil, fmt.Errorf("failed to create member contacts in the database")
		}
		final = append(final, createdContacts...)
	}

	if len(contactsToUpdate) > 0 {
		for _, contact := range contactsToUpdate {
			updated, err := s.contacts.UpdateMemberContact(ctx, member.ID, contact)
			if err != nil {
				entry.WithError(err).Error("failed to update member contacts in the database")
				return nil, fmt.Errorf("failed to update member contacts in the database")
			}
			final = append(final, updated)
		}
	}

	return final, nil
}

func (s *service) resolveContactAttributes(ctx context.Context, contacts []*athena.MemberContact) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "resolveContactAttributes",
	})

	for _, contact := range contacts {
		switch contact.ContactType {
		case "alliance":
			_, err := s.alliance.FetchAlliance(ctx, contact.ContactID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve alliance contact type")
			}
		case "character":
			_, err := s.character.FetchCharacter(ctx, contact.ContactID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve character contact type")
			}
		case "corporation":
			_, err := s.corporation.FetchCorporation(ctx, contact.ContactID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve corporation contact type")
			}
		case "faction":
			_, err := s.universe.Faction(ctx, contact.ContactID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve faction contact type")
			}
		default:
			entry.WithFields(logrus.Fields{
				"contact_id":   contact.ContactID,
				"contact_type": contact.ContactType,
			}).Error("unknown contact type")
		}

	}

}

func (s *service) EmptyMemberContactLabels(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContactLabels, esi.ModWithCharacterID(member.ID))
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
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "EmptyMemberContactLabels",
	})

	labels, etag, _, err := s.esi.GetCharacterContactLabels(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contact labels from ESI")
		return nil, fmt.Errorf("failed to fetch member contact labels from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, nil
	}

	existing, err := s.contacts.MemberContactLabels(ctx, member.ID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contact labels from DB")
		return nil, fmt.Errorf("failed to fetch member contact labels from DB")
	}

	labels, err = s.diffAndUpdateLabels(ctx, member, labels, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to diff and update contacts")
	}

	if len(labels) > 0 {
		err = s.cache.SetMemberContactLabels(ctx, member.ID, labels)
		if err != nil {
			entry.WithError(err).Error("failed to cache member contacts")
		}

	}
	return etag, err

}

func (s *service) MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberContactLabels",
	})

	labels, err := s.cache.MemberContactLabels(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contact labels from cache")
		return nil, fmt.Errorf("failed to fetch member contact labels from cache")
	}

	if len(labels) > 0 {
		return labels, nil
	}

	labels, err = s.contacts.MemberContactLabels(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contact labels from DB")
		return nil, fmt.Errorf("failed to fetch member contact labels from DB")
	}

	if len(labels) > 0 {
		err = s.cache.SetMemberContactLabels(ctx, memberID, labels)
		if err != nil {
			entry.WithError(err).Error("failed to cache member contact labels")
		}
	}

	return labels, nil

}

func (s *service) diffAndUpdateLabels(ctx context.Context, member *athena.Member, old []*athena.MemberContactLabel, new []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "diffAndUpdateLabels",
	})

	labelsToCreate := make([]*athena.MemberContactLabel, 0)
	labelsToUpdate := make([]*athena.MemberContactLabel, 0)
	labelsToDelete := make([]*athena.MemberContactLabel, 0)

	oldLabelMap := make(map[uint64]*athena.MemberContactLabel)
	for _, label := range old {
		oldLabelMap[label.LabelID] = label
	}

	for _, label := range new {
		// Never seen this label before, add it to the db
		if _, ok := oldLabelMap[label.LabelID]; !ok {
			labelsToCreate = append(labelsToCreate, label)

			// We've seen this Label before, check to see if the values are still the same
		} else if diff := deep.Equal(oldLabelMap[label.LabelID], label); len(diff) > 0 {
			labelsToUpdate = append(labelsToUpdate, label)
		}
	}

	newLabelMap := make(map[uint64]*athena.MemberContactLabel)
	for _, label := range new {
		newLabelMap[label.LabelID] = label
	}

	for _, label := range old {
		if _, ok := newLabelMap[label.LabelID]; !ok {
			labelsToDelete = append(labelsToDelete, label)
		}
	}

	if len(labelsToDelete) > 0 {
		_, err := s.contacts.DeleteMemberContactLabels(ctx, member.ID, labelsToDelete)
		if err != nil {
			entry.WithError(err).Error("failed to delete member contact labels in the database")
			return nil, fmt.Errorf("failed to delete member contact labels in the database")
		}
	}

	var final = make([]*athena.MemberContactLabel, 0, len(labelsToCreate)+len(labelsToUpdate))
	if len(labelsToCreate) > 0 {
		createdLabels, err := s.contacts.CreateMemberContactLabels(ctx, member.ID, labelsToCreate)
		if err != nil {
			entry.WithError(err).Error("failed to create member contact labels in the database")
			return nil, fmt.Errorf("failed to create member contact labels in the database")
		}
		final = append(final, createdLabels...)
	}

	if len(labelsToUpdate) > 0 {
		updated, err := s.contacts.UpdateMemberContactLabel(ctx, member.ID, labelsToUpdate)
		if err != nil {
			entry.WithError(err).Error("failed to update member contact labels in the database")
			return nil, fmt.Errorf("failed to update member contact labels in the database")
		}
		final = append(final, updated...)
	}

	return final, nil

}
