package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type mailService interface {
	MailHeader(ctx context.Context, mailID uint) (*athena.MailHeader, error)
	SetMailHeader(ctx context.Context, header *athena.MailHeader, optionFuncs ...OptionFunc) error

	MailHeaderRecipients(ctx context.Context, mailID uint) ([]*athena.MailRecipient, error)
	SetMailHeaderRecipients(ctx context.Context, mailID uint, recipients []*athena.MailRecipient, optionFuncs ...OptionFunc) error

	MemberMailHeaders(ctx context.Context, memberID, page uint) ([]*athena.MemberMailHeader, error)
	SetMemberMailHeaders(ctx context.Context, memberID, page uint, headers []*athena.MemberMailHeader, optionFuncs ...OptionFunc) error

	MemberMailLabels(ctx context.Context, memberID uint) (*athena.MemberMailLabels, error)
	SetMemberMailLabels(ctx context.Context, memberID uint, labels *athena.MemberMailLabels, optionFuncs ...OptionFunc) error

	MemberMailingLists(ctx context.Context, memberID uint) ([]*athena.MemberMailingList, error)
	SetMemberMailingLists(ctx context.Context, memberID uint, listID []*athena.MemberMailingList, optionFuncs ...OptionFunc) error

	MailingList(ctx context.Context, mailingListID uint) (*athena.MailingList, error)
	SetMailingList(ctx context.Context, list *athena.MailingList, optionFuncs ...OptionFunc) error
}

const (
	keyMailHeader         = "athena::mail::header::%d"              // *athena.MailHeader
	keyMailRecipients     = "athena::mail::recipients::%d"          // []*athena.MailRecipient
	keyMemberMailHeaders  = "athena::member::%d::mail::headers::%d" // []*athena.MemberMailHeaders (ties a header to a member)
	keyMemberMailLabels   = "athena::member::%d::mail::labels"      // *athena.MemberMailLabels
	keyMemberMailingLists = "athena::member::%d::mail::lists"       // []*athena.MemberMailingLists
	keyMailingList        = "athena::mail::list::%d"                // *athena.MailingLists
)

func (s *service) MailHeader(ctx context.Context, mailID uint) (*athena.MailHeader, error) {

	key := fmt.Sprintf(keyMailHeader, mailID)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch header from cache: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var header = new(athena.MailHeader)
	err = json.Unmarshal(result, header)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal mail header %d on struct: %w", mailID, err)
	}

	return header, nil

}

func (s *service) SetMailHeader(ctx context.Context, header *athena.MailHeader, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("[Cache Service] Failed to marshal header: %w", err)
	}

	key := fmt.Sprintf(keyMailHeader, header.MailID)
	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) MailHeaderRecipients(ctx context.Context, mailID uint) ([]*athena.MailRecipient, error) {

	key := fmt.Sprintf(keyMailRecipients, mailID)

	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("[Cache Layer] Failed to failed header from cache: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var recipients = make([]*athena.MailRecipient, 0)
	err = json.Unmarshal(result, &recipients)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal recipients on slice of struct: %w", err)
	}

	return recipients, nil

}

func (s *service) SetMailHeaderRecipients(ctx context.Context, mailID uint, recipients []*athena.MailRecipient, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(recipients)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal recipients: %w", err)
	}

	key := fmt.Sprintf(keyMailRecipients, mailID)
	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil
}

func (s *service) MemberMailHeaders(ctx context.Context, memberID, page uint) ([]*athena.MemberMailHeader, error) {

	key := fmt.Sprintf(keyMemberMailHeaders, memberID, page)
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch member mail headers from cache: %w", err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	var headers = make([]*athena.MemberMailHeader, 0, len(members))
	for _, member := range members {

		var header = new(athena.MemberMailHeader)
		err = json.Unmarshal([]byte(member), header)
		if err != nil {
			return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal set member onto struct: %w", err)
		}

		headers = append(headers, header)

	}

	return headers, nil

}

func (s *service) SetMemberMailHeaders(ctx context.Context, memberID, page uint, headers []*athena.MemberMailHeader, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	members := make([]string, 0, len(headers))
	for _, header := range headers {
		data, err := json.Marshal(header)
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to marshal member mail headers: %w", err)
		}

		members = append(members, string(data))

	}

	key := fmt.Sprintf(keyMemberMailHeaders, memberID, page)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) MemberMailLabels(ctx context.Context, memberID uint) (*athena.MemberMailLabels, error) {

	key := fmt.Sprintf(keyMemberMailLabels, memberID)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch member mail labels from cache: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var labels = new(athena.MemberMailLabels)
	err = json.Unmarshal(result, labels)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] failed to unmarshal member mail labels on struct: %w", err)
	}

	return labels, nil

}

func (s *service) SetMemberMailLabels(ctx context.Context, memberID uint, labels *athena.MemberMailLabels, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(labels)
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to marshal labels for member %d: %w", memberID, err)
	}

	key := fmt.Sprintf(keyMemberMailLabels, memberID)
	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) MailingList(ctx context.Context, mailingListID uint) (*athena.MailingList, error) {

	key := fmt.Sprintf(keyMailingList, mailingListID)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch mailing list from cache: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var list = new(athena.MailingList)
	err = json.Unmarshal(result, list)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal mail list %d on struct: %w", mailingListID, err)
	}

	return list, nil

}

func (s *service) SetMailingList(ctx context.Context, list *athena.MailingList, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("[Cache Service] Failed to marshal list: %w", err)
	}

	key := fmt.Sprintf(keyMailingList, list.MailingListID)
	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache for key %s: %w", key, err)
	}

	return nil

}

func (s *service) MemberMailingLists(ctx context.Context, memberID uint) ([]*athena.MemberMailingList, error) {

	key := fmt.Sprintf(keyMemberMailingLists, memberID)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch mailing list from cache: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	var mailingLists = make([]*athena.MemberMailingList, 0)
	err = json.Unmarshal(result, &mailingLists)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal mailingLists for member %d on struct: %w", memberID, err)
	}

	return mailingLists, nil

}

func (s *service) SetMemberMailingLists(ctx context.Context, memberID uint, lists []*athena.MemberMailingList, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(lists)
	if err != nil {
		return fmt.Errorf("[Cache Service] Failed to marshal member mailing lists for member %d: %w", memberID, err)
	}

	key := fmt.Sprintf(keyMemberMailingLists, memberID)
	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache for key %s: %w", key, err)
	}

	return nil

}
