package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type contactService interface {
	MemberContacts(ctx context.Context, memberID uint, page int) ([]*athena.MemberContact, error)
	SetMemberContacts(ctx context.Context, memberID uint, page int, contacts []*athena.MemberContact, optionFuncs ...OptionFunc) error
	MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error)
	SetMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel, optionFuncs ...OptionFunc) error
}

const (
	keyMemberContacts      = "athena::member::%d::contacts::%d"
	keyMemberContactPages  = "athena::member::%d::contact::pages"
	keyMemberContactLabels = "athena::member::%d::contact::labels"
)

func (s *service) MemberContacts(ctx context.Context, memberID uint, page int) ([]*athena.MemberContact, error) {

	pages := make([]string, 0, 11)

	if page < 0 {
		pages, err := s.client.SMembers(ctx, keyMemberContactPages).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, err
		}

		if len(pages) == 0 {
			return nil, nil
		}

	} else {
		pages = append(pages, fmt.Sprintf(keyMemberContacts, memberID, page))
	}
	contacts := make([]*athena.MemberContact, 0, len(pages)*1024)
	for _, page := range pages {

		members, err := s.client.SMembers(ctx, page).Result()
		if err != nil {
			return nil, fmt.Errorf("[Cache Layer] Failed to fetch set members for key %s: %w", page, err)
		}

		if len(members) == 0 {
			_, err = s.client.Del(ctx, page).Result()
			if err != nil {
				return nil, err
			}

			continue
		}

		for _, member := range members {
			var contact = new(athena.MemberContact)
			err = json.Unmarshal([]byte(member), contact)
			if err != nil {
				err = fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", page, err)
				newrelic.FromContext(ctx).NoticeError(err)
				continue
			}

			contacts = append(contacts, contact)
		}

	}

	return contacts, nil

}

func (s *service) SetMemberContacts(ctx context.Context, memberID uint, page int, contacts []*athena.MemberContact, optionFuncs ...OptionFunc) error {

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

	key := fmt.Sprintf(keyMemberContacts, memberID, page)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache contacts for member %d: %w", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
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

func (s *service) SetMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

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

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}
