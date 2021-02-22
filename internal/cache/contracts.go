package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
)

type contractService interface {
	MemberContracts(ctx context.Context, memberID, page uint) ([]*athena.MemberContract, error)
	SetMemberContracts(ctx context.Context, memberID, page uint, contracts []*athena.MemberContract) error

	MemberContractItems(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractItem, error)
	SetMemberContractItems(ctx context.Context, memberID, contractID uint, bids []*athena.MemberContractItem) error

	MemberContractBids(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractBid, error)
	SetMemberContractBids(ctx context.Context, memberID, contractID uint, items []*athena.MemberContractBid) error
}

const (
	// Segmented sets by page. Each set will hold up to 1000 contracts
	keyMemberContracts     = "athena::member::%d::contracts::%d"
	keyMemberContractBids  = "athena::member::%d::contracts::%d::bids"
	keyMemberContractItems = "athena::member::%d::contracts::%d::items"
)

const (
	// The follow const are strings meant to be passed to fmt.Errorf. They may have
	// format args included in the string
	errMaxNumContractsExceeded    = "[Cache Layer] Max number of contracts is limited to %d"
	errFailedToCacheMembers       = "[Cache Layer] Failed to cache set members for key %s: %w"
	errFailedToCachePage          = "[Cache Layer] Failed to cache page %d of %s for member %d: %w"
	errFailedToUnmarshalSetMember = "[Cache Layer] Failed to unmarshal member of set %s: %w"
	errFailedToSetExpiry          = "[Cache Layer] Failed to set expiry on key %s: %w"
)

func (s *service) MemberContracts(ctx context.Context, memberID, page uint) ([]*athena.MemberContract, error) {

	key := fmt.Sprintf(keyMemberContracts, memberID, page)

	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errFailedToCacheMembers, key, err)
	}
	if len(members) == 0 {
		return nil, nil
	}

	contracts := make([]*athena.MemberContract, 0, len(members))
	for _, member := range members {

		var contract = new(athena.MemberContract)
		err = json.Unmarshal([]byte(member), contract)
		if err != nil {
			return nil, fmt.Errorf(errFailedToUnmarshalSetMember, key, err)
		}

		contracts = append(contracts, contract)

	}

	return contracts, nil

}

func (s *service) SetMemberContracts(ctx context.Context, memberID, page uint, contracts []*athena.MemberContract) error {

	if len(contracts) > 1000 {
		return fmt.Errorf(errMaxNumContractsExceeded, 1000)
	}

	members := make([]string, 0, len(contracts))
	for _, contract := range contracts {
		data, err := json.Marshal(contract)
		if err != nil {
			return fmt.Errorf("failed to marshal contract: %w", err)
		}

		members = append(members, string(data))
	}

	key := fmt.Sprintf(keyMemberContracts, memberID, page)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf(errFailedToCachePage, page, "contracts", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf(errFailedToSetExpiry, key, err)
	}

	return nil

}

func (s *service) MemberContractItems(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractItem, error) {

	key := fmt.Sprintf(keyMemberContractItems, memberID, contractID)

	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errFailedToCacheMembers, key, err)
	}
	if len(members) == 0 {
		return nil, nil
	}

	items := make([]*athena.MemberContractItem, 0, len(members))
	for _, member := range members {
		var item = new(athena.MemberContractItem)
		err = json.Unmarshal([]byte(member), item)
		if err != nil {
			return nil, fmt.Errorf(errFailedToUnmarshalSetMember, key, err)
		}

		items = append(items, item)
	}

	return items, nil

}

func (s *service) SetMemberContractItems(ctx context.Context, memberID, contractID uint, items []*athena.MemberContractItem) error {

	members := make([]string, 0, len(items))
	for _, item := range items {
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal skill queue position: %w", err)
		}

		members = append(members, string(data))
	}

	key := fmt.Sprintf(keyMemberContractItems, memberID, contractID)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf(errFailedToCacheMembers, key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf(errFailedToSetExpiry, key, err)
	}

	return nil

}

func (s *service) MemberContractBids(ctx context.Context, memberID, contractID uint) ([]*athena.MemberContractBid, error) {

	key := fmt.Sprintf(keyMemberContractBids, memberID, contractID)

	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errFailedToCacheMembers, key, err)
	}
	if len(members) == 0 {
		return nil, nil
	}

	bids := make([]*athena.MemberContractBid, 0, len(members))
	for _, member := range members {
		var bid = new(athena.MemberContractBid)
		err = json.Unmarshal([]byte(member), bid)
		if err != nil {
			return nil, fmt.Errorf(errFailedToUnmarshalSetMember, key, err)
		}

		bids = append(bids, bid)
	}

	return bids, nil

}

func (s *service) SetMemberContractBids(ctx context.Context, memberID, contractID uint, bids []*athena.MemberContractBid) error {

	members := make([]interface{}, len(bids))
	for i, bid := range bids {
		members[i] = bid
	}

	key := fmt.Sprintf(keyMemberContractBids, memberID, contractID)
	_, err := s.client.SAdd(ctx, key, members...).Result()
	if err != nil {
		return fmt.Errorf(errFailedToCacheMembers, key, err)
	}

	_, err = s.client.Expire(ctx, key, time.Hour).Result()
	if err != nil {
		return fmt.Errorf(errFailedToSetExpiry, key, err)
	}

	return nil

}
