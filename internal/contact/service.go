package contact

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/alliance"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/character"
	"github.com/eveisesi/athena/internal/corporation"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/go-test/deep"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	EmptyMemberContacts(ctx context.Context, member *athena.Member) error
	MemberContacts(ctx context.Context, member *athena.Member) ([]*athena.MemberContact, error)
	EmptyMemberContactLabels(ctx context.Context, member *athena.Member) error
	MemberContactLabels(ctx context.Context, member *athena.Member) ([]*athena.MemberContactLabel, error)
}

type service struct {
	logger *logrus.Logger

	cache       cache.Service
	esi         esi.Service
	etag        etag.Service
	alliance    alliance.Service
	character   character.Service
	corporation corporation.Service
	universe    universe.Service

	contacts athena.MemberContactRepository
}

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, etag etag.Service, universe universe.Service, alliance alliance.Service, character character.Service, corporation corporation.Service, contacts athena.MemberContactRepository) Service {

	return &service{
		logger: logger,

		cache:       cache,
		esi:         esi,
		etag:        etag,
		universe:    universe,
		alliance:    alliance,
		character:   character,
		corporation: corporation,

		contacts: contacts,
	}

}

func (s *service) EmptyMemberContacts(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberContacts(ctx, member)

	return err

}

func (s *service) MemberContacts(ctx context.Context, member *athena.Member) ([]*athena.MemberContact, error) {

	valid, etagID := s.esi.GenerateEndpointHash(esi.EndpointGetCharacterContacts, member)
	if !valid {
		return nil, fmt.Errorf("failed to generate valid etag hash")
	}

	cached := true

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	contacts, err := s.cache.MemberContacts(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if contacts == nil {
		cached = false

		contacts, err = s.contacts.MemberContacts(ctx, member.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	if etag.CachedUntil.After(time.Now()) && len(contacts) > 0 {

		if !cached {
			err = s.cache.SetMemberContacts(ctx, member.ID.Hex(), contacts)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return contacts, nil
	}

	newContacts, _, _, err := s.esi.GetCharacterContacts(ctx, member, etag, make([]*athena.MemberContact, 0))
	if err != nil {
		return nil, fmt.Errorf("[Contacts Service] Failed to fetch contacts for member %s: %w", member.ID.Hex(), err)
	}

	_, _ = s.etag.UpdateEtag(ctx, etag.EtagID, etag)

	if len(newContacts) > 0 {
		s.resolveContactAttributes(ctx, newContacts)
		contacts, err = s.diffAndUpdateContacts(ctx, member, contacts, newContacts)
		if len(contacts) > 0 {
			err := s.cache.SetMemberContacts(ctx, member.ID.Hex(), contacts)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}
	}

	return contacts, err

}

func (s *service) diffAndUpdateContacts(ctx context.Context, member *athena.Member, old []*athena.MemberContact, new []*athena.MemberContact) ([]*athena.MemberContact, error) {

	contactsToCreate := make([]*athena.MemberContact, 0)
	contactsToUpdate := make([]*athena.MemberContact, 0)
	contactsToDelete := make([]*athena.MemberContact, 0)

	oldContactMap := make(map[int]*athena.MemberContact)
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

	newContactMap := make(map[int]*athena.MemberContact)
	for _, contact := range new {
		newContactMap[contact.ContactID] = contact
	}

	for _, contact := range old {
		// This label is not in the list of new label, must've been deleted by the user in game
		if _, ok := newContactMap[contact.ContactID]; !ok {
			contactsToDelete = append(contactsToDelete, contact)
		}
	}

	var final = make([]*athena.MemberContact, 0)
	if len(contactsToCreate) > 0 {
		createdContacts, err := s.contacts.CreateMemberContacts(ctx, member.ID.Hex(), contactsToCreate)
		if err != nil {
			return nil, err
		}
		final = append(final, createdContacts...)
	}

	if len(contactsToUpdate) > 0 {
		for _, contact := range contactsToUpdate {
			updated, err := s.contacts.UpdateMemberContact(ctx, member.ID.Hex(), contact)
			if err != nil {
				return nil, err
			}
			final = append(final, updated)
		}
	}

	if len(contactsToDelete) > 0 {
		deleteOK, err := s.contacts.DeleteMemberContacts(ctx, member.ID.Hex(), contactsToDelete)
		if err != nil {
			return nil, err
		}

		if !deleteOK {
			return nil, fmt.Errorf("Expected to delete %d documents, deleted none", len(contactsToDelete))
		}
	}

	return final, nil
}

func (s *service) resolveContactAttributes(ctx context.Context, contacts []*athena.MemberContact) {

	for _, contact := range contacts {
		switch contact.ContactType {
		case "alliance":
			_, err := s.alliance.AllianceByAllianceID(ctx, uint(contact.ContactID), alliance.NewOptionFuncs())
			if err != nil {
				s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve alliance contact type")
				continue
			}
		case "character":
			_, err := s.character.CharacterByCharacterID(ctx, uint64(contact.ContactID), character.NewOptionFuncs())
			if err != nil {
				s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve character contact type")
				continue
			}
		case "corporation":
			_, err := s.corporation.CorporationByCorporationID(ctx, uint(contact.ContactID), corporation.NewOptionFuncs())
			if err != nil {
				s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve corporation contact type")
				continue
			}
		case "faction":
			_, err := s.universe.Faction(ctx, contact.ContactID)
			if err != nil {
				s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
					"contact_id":   contact.ContactID,
					"contact_type": contact.ContactType,
				}).Error("failed to resolve faction contact type")
				continue
			}
		default:
			s.logger.WithContext(ctx).WithFields(logrus.Fields{
				"contact_id":   contact.ContactID,
				"contact_type": contact.ContactType,
			}).Error("unknown contact type")
		}

	}

}

func (s *service) EmptyMemberContactLabels(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberContactLabels(ctx, member)

	return err

}

func (s *service) MemberContactLabels(ctx context.Context, member *athena.Member) ([]*athena.MemberContactLabel, error) {

	valid, etagID := s.esi.GenerateEndpointHash(esi.EndpointGetCharacterContactLabels, member)
	if !valid {
		return nil, fmt.Errorf("failed to generate valid etag hash")
	}

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	cached := true

	labels, err := s.cache.MemberContactLabels(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if labels == nil {
		cached = false
		labels, err = s.contacts.MemberContactLabels(ctx, member.ID.Hex())
		if err != nil {
			return nil, err
		}
	}

	if etag.CachedUntil.After(time.Now()) && len(labels) > 0 {
		if !cached {
			err = s.cache.SetMemberContactLabels(ctx, member.ID.Hex(), labels)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return labels, nil
	}

	newLabels, etag, _, err := s.esi.GetCharacterContactLabels(ctx, member, etag, make([]*athena.MemberContactLabel, 0))
	if err != nil {
		return nil, fmt.Errorf("[Contacts Service] Failed to fetch labels for member %s: %w", member.ID.Hex(), err)
	}

	_, _ = s.etag.UpdateEtag(ctx, etag.EtagID, etag)

	if len(newLabels) > 0 {
		labels, err = s.diffAndUpdateLabels(ctx, member, labels, newLabels)
		if err != nil {
			return nil, err
		}

		if len(labels) > 0 {
			err = s.cache.SetMemberContactLabels(ctx, member.ID.Hex(), labels)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}
	}

	return labels, err

}

func (s *service) diffAndUpdateLabels(ctx context.Context, member *athena.Member, old []*athena.MemberContactLabel, new []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, error) {

	labelsToCreate := make([]*athena.MemberContactLabel, 0)
	labelsToUpdate := make([]*athena.MemberContactLabel, 0)
	labelsToDelete := make([]*athena.MemberContactLabel, 0)

	oldLabelMap := make(map[int64]*athena.MemberContactLabel)
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

	newLabelMap := make(map[int64]*athena.MemberContactLabel)
	for _, label := range new {
		newLabelMap[label.LabelID] = label
	}

	for _, label := range old {
		if _, ok := newLabelMap[label.LabelID]; !ok {
			labelsToDelete = append(labelsToDelete, label)
		}
	}

	var final = make([]*athena.MemberContactLabel, 0)
	if len(labelsToCreate) > 0 {
		createdLabels, err := s.contacts.CreateMemberContactLabels(ctx, member.ID.Hex(), labelsToCreate)
		if err != nil {
			return nil, err
		}
		final = append(final, createdLabels...)
	}

	if len(labelsToUpdate) > 0 {
		for _, label := range labelsToUpdate {
			updated, err := s.contacts.UpdateMemberContactLabel(ctx, member.ID.Hex(), label)
			if err != nil {
				return nil, err
			}
			final = append(final, updated)
		}
	}

	if len(labelsToDelete) > 0 {
		deleteOK, err := s.contacts.DeleteMemberContactLabels(ctx, member.ID.Hex(), labelsToDelete)
		if err != nil {
			return nil, err
		}

		if !deleteOK {
			return nil, fmt.Errorf("Expected to delete %d documents, deleted none", len(labelsToDelete))
		}
	}

	return final, nil

}
