package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type contactService interface {
	MemberContacts(ctx context.Context, memberID, page uint) ([]*athena.MemberContact, error)
	SetMemberContacts(ctx context.Context, memberID, page uint, contacts []*athena.MemberContact) error
	MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error)
	SetMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel) error
}

const (
	keyMemberContacts      = "athena::member::%d::contacts::%d"
	keyMemberContactPages  = "athena::member::%d::contact::pages"
	keyMemberContactLabels = "athena::member::%d::contact::labels"
)

func (s *service) MemberContacts(ctx context.Context, memberID, page uint) ([]*athena.MemberContact, error) {

	key := fmt.Sprintf(keyMemberContacts, memberID, page)
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch set members for key %s: %w", key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	contacts := make([]*athena.MemberContact, 0, len(members))
	for _, member := range members {
		var contact = new(athena.MemberContact)
		err = json.Unmarshal([]byte(member), contact)
		if err != nil {
			return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil

}

func (s *service) SetMemberContacts(ctx context.Context, memberID, page uint, contacts []*athena.MemberContact) error {

	// Build the interface to send to redis
	members := make([]string, 0, len(contacts))
	for _, contact := range contacts {
		data, err := json.Marshal(contact)
		if err != nil {
			return fmt.Errorf("failed to marshal skill queue position: %w", err)
		}

		members = append(members, string(data))
	}

	key := fmt.Sprintf(keyMemberContacts, memberID, page)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache contacts for member %d: %w", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	pageKey := fmt.Sprintf(keyMemberContactPages, memberID)
	_, err = s.client.SAdd(ctx, pageKey, []string{key}).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to push key to contact key set: %w", err)
	}

	_, err = s.client.Expire(ctx, pageKey, 0).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error) {

	key := fmt.Sprintf(keyMemberContactLabels, memberID)
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch set members for key %s: %w", key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	labels := make([]*athena.MemberContactLabel, 0, len(members))
	for _, member := range members {
		var label = new(athena.MemberContactLabel)
		err = json.Unmarshal([]byte(member), label)
		if err != nil {
			err = fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		labels = append(labels, label)
	}

	return labels, nil

}

func (s *service) SetMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel) error {

	// Build the interface to send to redis
	members := make([]string, 0, len(labels))
	for _, label := range labels {
		b, err := json.Marshal(label)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		members = append(members, string(b))
	}

	key := fmt.Sprintf(keyMemberContactLabels, memberID)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache labels for member %d: %w", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}
