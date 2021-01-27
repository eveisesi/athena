package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirkon/go-format"
)

type contractService interface {
	MemberContracts(ctx context.Context, memberID string, page int) ([]*athena.MemberContract, error)
	SetMemberContracts(ctx context.Context, memberID string, page int, contracts []*athena.MemberContract, optionFunc ...OptionFunc) error
}

const (
	// Segmented sets by page. Each set will hold up to 1000 contracts
	keyMemberContracts = "athena::member::${memberID}::contracts::${pageID}"
)

const (
	// The follow const are strings meant to be passed to fmt.Errorf. They may have
	// format args included in the string
	errMaxNumContractsExceeded    = "[Cache Layer] Max number of contracts is limited to %d"
	errFailedToCacheMembers       = "[Cache Layer] Failed to cache set members for key %s: %w"
	errFailedToCachePage          = "[Cache Layer] Failed to cache page %d of %s for member %s: %w"
	errFailedToUnmarshalSetMember = "[Cache Layer] Failed to unmarshal member of set %s: %w"
	errFailedToSetExpiry          = "[Cache Layer] Failed to set expiry on key %s: %w"
)

func (s *service) MemberContracts(ctx context.Context, memberID string, page int) ([]*athena.MemberContract, error) {

	key := format.Formatm(keyMemberContracts, format.Values{
		"memberID": memberID,
		"page":     page,
	})

	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf(errFailedToCacheMembers, key, err)
	}
	if len(members) == 0 {
		return nil, nil
	}

	contracts := make([]*athena.MemberContract, len(members))
	for i, member := range members {
		var contact = new(athena.MemberContract)
		err = json.Unmarshal([]byte(member), contact)
		if err != nil {
			err = fmt.Errorf(errFailedToUnmarshalSetMember, key, err)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		contracts[i] = contact
	}

	return contracts, nil

}

func (s *service) SetMemberContracts(ctx context.Context, memberID string, page int, contracts []*athena.MemberContract, optionFuncs ...OptionFunc) error {

	if len(contracts) > 1000 {
		return fmt.Errorf(errMaxNumContractsExceeded, 1000)
	}

	members := make([]interface{}, len(contracts))
	for i, contract := range contracts {
		members[i] = contract
	}

	options := applyOptionFuncs(nil, optionFuncs)

	key := format.Formatm(keyMemberContracts, format.Values{
		"memberID": memberID,
		"page":     page,
	})

	_, err := s.client.SAdd(ctx, key, members...).Result()
	if err != nil {
		return fmt.Errorf(errFailedToCachePage, page, "contracts", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf(errFailedToSetExpiry, key, err)
	}

	return nil

}
