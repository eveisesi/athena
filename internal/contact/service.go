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
	MemberContacts(ctx context.Context, member *athena.Member) ([]*athena.MemberContact, *athena.Etag, error)
	EmptyMemberContactLabels(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberContactLabels(ctx context.Context, member *athena.Member) ([]*athena.MemberContactLabel, *athena.Etag, error)
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

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContacts, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberContacts(ctx, member)

	return etag, err

}

func (s *service) MemberContacts(ctx context.Context, member *athena.Member) ([]*athena.MemberContact, *athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "MemberContacts",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContacts, esi.ModWithMember(member))
	if err != nil {
		entry.WithError(err).Error("failed to fetch etag object")
		return nil, nil, fmt.Errorf("failed to fetch etag object")
	}

	cached := true

	contacts, err := s.cache.MemberContacts(ctx, member.ID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contacts from cache")
		return nil, nil, fmt.Errorf("failed to fetch member contacts from cache")
	}

	if contacts == nil {
		cached = false
		contacts, err = s.contacts.MemberContacts(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			entry.WithError(err).Error("failed to fetch member contacts from DB")
			return nil, nil, fmt.Errorf("failed to fetch member contacts from DB")
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {

		if !cached {
			err = s.cache.SetMemberContacts(ctx, member.ID, contacts)
			if err != nil {
				entry.WithError(err).Error("failed to cache member contacts")
			}
		}

		return contacts, etag, nil
	}

	newContacts, etag, _, err := s.esi.GetCharacterContacts(ctx, member, make([]*athena.MemberContact, 0))
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contacts from ESI")
		return nil, nil, fmt.Errorf("failed to fetch member contacts from ESI")
	}

	if len(newContacts) > 0 {
		s.resolveContactAttributes(ctx, newContacts)
		contacts, err = s.diffAndUpdateContacts(ctx, member, contacts, newContacts)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to diff and update contacts")
		}
		if len(contacts) > 0 {
			err := s.cache.SetMemberContacts(ctx, member.ID, contacts)
			if err != nil {
				entry.WithError(err).Error("failed to cache member contacts")
			}
		}
	}

	return contacts, etag, nil

}

func (s *service) diffAndUpdateContacts(ctx context.Context, member *athena.Member, old []*athena.MemberContact, new []*athena.MemberContact) ([]*athena.MemberContact, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"service": serviceIdentifier,
		"method":  "diffAndUpdateContacts",
	})

	contactsToCreate := make([]*athena.MemberContact, 0)
	contactsToUpdate := make([]*athena.MemberContact, 0)
	contactsToDelete := make([]*athena.MemberContact, 0)

	oldContactMap := make(map[uint]*athena.MemberContact)
	for _, contact := range old {
		oldContactMap[contact.ContactID] = contact
	}

	for _, contact := range new {
		var ok bool
		// This is an unknown contact, so lets flag it to be created
		if _, ok = oldContactMap[contact.ContactID]; !ok {
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
			_, err := s.alliance.Alliance(ctx, contact.ContactID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve alliance contact type")
			}
		case "character":
			_, err := s.character.Character(ctx, contact.ContactID)
			if err != nil {
				entry.WithError(err).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve character contact type")
			}
		case "corporation":
			_, err := s.corporation.Corporation(ctx, contact.ContactID)
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

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContactLabels, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberContactLabels(ctx, member)

	return etag, err

}

func (s *service) MemberContactLabels(ctx context.Context, member *athena.Member) ([]*athena.MemberContactLabel, *athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "MemberContacts",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterContactLabels, esi.ModWithMember(member))
	if err != nil {
		entry.WithError(err).Error("failed to fetch etag object")
		return nil, nil, fmt.Errorf("failed to fetch etag object")
	}

	cached := true

	labels, err := s.cache.MemberContactLabels(ctx, member.ID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contact labels from cache")
		return nil, nil, fmt.Errorf("failed to fetch member contact labels from cache")
	}

	if labels == nil {
		cached = false
		labels, err = s.contacts.MemberContactLabels(ctx, member.ID)
		if err != nil {
			entry.WithError(err).Error("failed to fetch member contact labels from DB")
			return nil, nil, fmt.Errorf("failed to fetch member contact labels from DB")
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {

		if !cached && len(labels) > 0 {
			err = s.cache.SetMemberContactLabels(ctx, member.ID, labels)
			if err != nil {
				entry.WithError(err).Error("failed to cache member contact labels")
			}
		}

		return labels, etag, nil
	}

	newLabels, etag, _, err := s.esi.GetCharacterContactLabels(ctx, member, make([]*athena.MemberContactLabel, 0))
	if err != nil {
		entry.WithError(err).Error("failed to fetch member contact labels from ESI")
		return nil, nil, fmt.Errorf("failed to fetch member contact labels from ESI")
	}

	if len(newLabels) > 0 {
		labels, err = s.diffAndUpdateLabels(ctx, member, labels, newLabels)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to diff and update contacts")
		}

		if len(labels) > 0 {
			err = s.cache.SetMemberContactLabels(ctx, member.ID, labels)
			if err != nil {
				entry.WithError(err).Error("failed to cache member contacts")
			}
		}
	}

	return labels, etag, err

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

	var final = make([]*athena.MemberContactLabel, 0)
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
