package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirkon/go-format"
)

type contactService interface {
	MemberContacts(ctx context.Context, memberID uint) ([]*athena.MemberContact, error)
	SetMemberContacts(ctx context.Context, memberID uint, contacts []*athena.MemberContact, optionFuncs ...OptionFunc) error
	MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error)
	SetMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel, optionFuncs ...OptionFunc) error
}

const (
	keyMemberContacts      = "athena::member::${id}::contacts"
	keyMemberContactLabels = "athena::member::${id}::contact::labels"
)

func (s *service) MemberContacts(ctx context.Context, memberID uint) ([]*athena.MemberContact, error) {

	key := format.Formatm(keyMemberContacts, format.Values{
		"id": memberID,
	})
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch set members for key %s: %w", key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	contacts := make([]*athena.MemberContact, len(members))
	for i, member := range members {
		var contact = new(athena.MemberContact)
		err = json.Unmarshal([]byte(member), contact)
		if err != nil {
			err = fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		contacts[i] = contact
	}

	return contacts, nil

}

func (s *service) SetMemberContacts(ctx context.Context, memberID uint, contacts []*athena.MemberContact, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	// Build the interface to send to redis
	members := make([]string, 0, len(contacts))
	for _, contact := range contacts {
		data, err := json.Marshal(contact)
		if err != nil {
			return fmt.Errorf("failed to marshal skill queue position: %w", err)
		}

		members = append(members, string(data))
	}

	// Send members to redis
	key := format.Formatm(keyMemberContacts, format.Values{
		"id": memberID,
	})
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		// spew.Dump("Original:", contacts, "Members: ", members)
		return fmt.Errorf("[Cache Layer] Failed to cache contacts for member %d: %w", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error) {

	key := format.Formatm(keyMemberContactLabels, format.Values{
		"id": memberID,
	})
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch set members for key %s: %w", key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	labels := make([]*athena.MemberContactLabel, len(members))
	for i, member := range members {
		var label = new(athena.MemberContactLabel)
		err = json.Unmarshal([]byte(member), label)
		if err != nil {
			err = fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		labels[i] = label
	}

	return labels, nil

}

func (s *service) SetMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	// Build the interface to send to redis
	members := make([]interface{}, len(labels))
	for i, label := range labels {
		b, err := json.Marshal(label)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		members[i] = b
	}

	// Send members to redis
	key := format.Formatm(keyMemberContactLabels, format.Values{
		"id": memberID,
	})
	_, err := s.client.SAdd(ctx, key, members...).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache labels for member %d: %w", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}
