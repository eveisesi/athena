package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
	"github.com/sirkon/go-format"
)

type corporationService interface {
	Corporation(ctx context.Context, id uint) (*athena.Corporation, error)
	SetCorporation(ctx context.Context, corporation *athena.Corporation, optionFuncs ...OptionFunc) error
	CorporationAllianceHistory(ctx context.Context, id uint) ([]*athena.CorporationAllianceHistory, error)
	SetCorporationAllianceHistory(ctx context.Context, id uint, history []*athena.CorporationAllianceHistory, optionFuncs ...OptionFunc) error
	// Corporations(ctx context.Context, operators []*athena.Operator) ([]*athena.Corporation, error)
	// SetCorporations(ctx context.Context, operators []*athena.Operator, corporations []*athena.Corporation, optionFuncs ...OptionFunc) error
}

const (
	keyCorporation                = "athena::corporation::%d"
	keyCorporationAllianceHistory = "athena::corporation::%d::history"
	// keyCorporations = "athena::corporations::%s"
)

func (s *service) Corporation(ctx context.Context, id uint) (*athena.Corporation, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyCorporation, id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var corporation = new(athena.Corporation)

	err = json.Unmarshal([]byte(result), &corporation)
	if err != nil {
		return nil, err
	}

	return corporation, nil

}

func (s *service) SetCorporation(ctx context.Context, corporation *athena.Corporation, optionFuncs ...OptionFunc) error {

	data, err := json.Marshal(corporation)
	if err != nil {
		return fmt.Errorf("failed to marshal corporation: %w", err)
	}

	options := applyOptionFuncs(nil, optionFuncs)

	_, err = s.client.Set(ctx, fmt.Sprintf(keyCorporation, corporation.ID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (s *service) CorporationAllianceHistory(ctx context.Context, id uint) ([]*athena.CorporationAllianceHistory, error) {

	key := format.Formatm(keyCorporationAllianceHistory, format.Values{
		"id": id,
	})

	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errKeyNotFound, key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	records := make([]*athena.CorporationAllianceHistory, len(members))
	for i, member := range members {
		var record = new(athena.CorporationAllianceHistory)
		err = json.Unmarshal([]byte(member), record)
		if err != nil {
			return nil, fmt.Errorf(errSetUnmarshalFailed, key, err)
		}

		records[i] = record

	}

	return records, nil

}

func (s *service) SetCorporationAllianceHistory(ctx context.Context, id uint, records []*athena.CorporationAllianceHistory, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	members := make([]string, 0, len(records))
	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal skill queue position: %w", err)
		}

		members = append(members, string(data))
	}

	key := format.Formatm(keyCharacterCorporationHistory, format.Values{
		"id": id,
	})
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache alliance history for corporation %d: %w", id, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}

// func (s *service) Corporations(ctx context.Context, operators []*athena.Operator) ([]*athena.Corporation, error) {

// 	data, err := json.Marshal(operators)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal operators: %w", err)
// 	}

// 	h := sha256.New()
// 	_, _ = h.Write(data)
// 	bs := h.Sum(nil)

// 	result, err := s.client.Get(ctx, fmt.Sprintf(keyCorporations, fmt.Sprintf("%x", bs))).Result()
// 	if err != nil && err != redis.Nil {
// 		return nil, err
// 	}

// 	if len(result) > 0 {
// 		var corporations = make([]*athena.Corporation, 0)

// 		err = json.Unmarshal([]byte(result), &corporations)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return corporations, nil
// 	}

// 	return nil, nil

// }

// func (s *service) SetCorporations(ctx context.Context, operators []*athena.Operator, corporations []*athena.Corporation, optionFuncs ...OptionFunc) error {

// 	data, err := json.Marshal(operators)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal operators: %w", err)
// 	}

// 	h := sha256.New()
// 	_, _ = h.Write(data)
// 	bs := h.Sum(nil)

// 	data, err = json.Marshal(corporations)
// 	if err != nil {
// 		return fmt.Errorf("Failed to marsahl payload: %w", err)
// 	}

// 	options := applyOptionFuncs(nil, optionFuncs)

// 	_, err = s.client.Set(ctx, fmt.Sprintf(keyCorporations, fmt.Sprintf("%x", bs)), data, options.expiry).Result()
// 	if err != nil {
// 		return fmt.Errorf("failed to write to cache: %w", err)
// 	}

// 	return nil

// }
