package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type walletService interface {
	MemberWalletBalance(ctx context.Context, member *athena.Member) (*athena.MemberWalletBalance, error)
	SetMemberWalletBalance(ctx context.Context, member *athena.Member, balance *athena.MemberWalletBalance, optionFuncs ...OptionFunc) error

	MemberWalletTransactions(ctx context.Context, member *athena.Member) ([]*athena.MemberWalletTransaction, error)
	SetMemberWalletTransactions(ctx context.Context, member *athena.Member, transactions []*athena.MemberWalletTransaction, optionFuncs ...OptionFunc) error

	MemberWalletJournal(ctx context.Context, member *athena.Member) ([]*athena.MemberWalletJournal, error)
	SetMemberWalletJournal(ctx context.Context, member *athena.Member, entries []*athena.MemberWalletJournal, optionFuncs ...OptionFunc) error
}

const (
	keyMemberWalletBalance = "athena::member::%d::wallet::balance"

	keyMemberWalletTransactions       = "athena::member::%d::wallet::transactions::%d"
	keyMemberWalletTransactionIndexes = "athena::member::%d::wallet::transactions::indexes"
	limitMemberWalletTransactions     = 1000

	keyMemberWalletJournals       = "athena::member::%d::wallet::journals::%d"
	keyMemberWalletJournalIndexes = "athena::member::%d::wallet::journals::indexes"
	limitMemberWalletJournals     = 1000
)

func (s *service) MemberWalletBalance(ctx context.Context, member *athena.Member) (*athena.MemberWalletBalance, error) {

	key := fmt.Sprintf(keyMemberWalletBalance, member.ID)

	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var balance = new(athena.MemberWalletBalance)
	err = json.Unmarshal(result, balance)
	if err != nil {
		return nil, err
	}

	return balance, nil

}

func (s *service) SetMemberWalletBalance(ctx context.Context, member *athena.Member, balance *athena.MemberWalletBalance, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(balance)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	key := fmt.Sprintf(keyMemberWalletBalance, member.ID)

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberWalletTransactions(ctx context.Context, member *athena.Member) ([]*athena.MemberWalletTransaction, error) {

	indexKey := fmt.Sprintf(keyMemberWalletTransactionIndexes, member.ID)
	indexes, err := s.client.SMembers(ctx, indexKey).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch wallet transaction indexes")
	}

	transactions := make([]*athena.MemberWalletTransaction, 0, len(indexes)*limitMemberWalletTransactions)

	for _, indexKey := range indexes {
		members, err := s.client.SMembers(ctx, indexKey).Result()
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch members for key %s: %w", indexKey, err)
		}

		for _, member := range members {
			var transaction = new(athena.MemberWalletTransaction)
			err = json.Unmarshal([]byte(member), transaction)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarsahl member onto struct: %w", err)
			}

			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil

}

func (s *service) SetMemberWalletTransactions(ctx context.Context, member *athena.Member, transactions []*athena.MemberWalletTransaction, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)
	chunks := chunkTransactions(transactions, limitMemberWalletTransactions)

	indexKey := fmt.Sprintf(keyMemberWalletTransactionIndexes, member.ID)

	for i, chunk := range chunks {

		members := make([]string, 0, len(chunk))
		for _, transaction := range chunk {
			data, err := json.Marshal(transaction)
			if err != nil {
				return fmt.Errorf("failed to marshal transaction for cache: %w", err)
			}

			members = append(members, string(data))
		}

		key := fmt.Sprintf(keyMemberWalletTransactions, member.ID, i+1)

		_, err := s.client.SAdd(ctx, key, members).Result()
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to cache transactions for member %d: %w", member.ID, err)
		}

		_, err = s.client.Expire(ctx, key, options.expiry).Result()
		if err != nil {
			return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
		}

		// TODO: Figure out if we should bail since we can't cache the key
		// For right now do nothing if it errors
		_, _ = s.client.SAdd(ctx, indexKey, key).Result()

	}

	return nil

}

func (s *service) MemberWalletJournal(ctx context.Context, member *athena.Member) ([]*athena.MemberWalletJournal, error) {

	indexKey := fmt.Sprintf(keyMemberWalletJournalIndexes, member.ID)
	indexes, err := s.client.SMembers(ctx, indexKey).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch wallet journal indexes")
	}

	journals := make([]*athena.MemberWalletJournal, 0, len(indexes)*limitMemberWalletJournals)

	for _, indexKey := range indexes {
		members, err := s.client.SMembers(ctx, indexKey).Result()
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch members for key %s: %w", indexKey, err)
		}

		for _, member := range members {
			var journal = new(athena.MemberWalletJournal)
			err = json.Unmarshal([]byte(member), journal)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarsahl member onto struct: %w", err)
			}

			journals = append(journals, journal)
		}
	}

	return journals, nil

}

func (s *service) SetMemberWalletJournal(ctx context.Context, member *athena.Member, entries []*athena.MemberWalletJournal, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)
	chunks := chunkJournals(entries, limitMemberWalletJournals)

	indexKey := fmt.Sprintf(keyMemberWalletJournalIndexes, member.ID)

	for i, chunk := range chunks {

		members := make([]string, 0, len(chunk))
		for _, entry := range chunk {
			data, err := json.Marshal(entry)
			if err != nil {
				return fmt.Errorf("failed to marshal journal entry for cache: %w", err)
			}

			members = append(members, string(data))
		}

		key := fmt.Sprintf(keyMemberWalletJournals, member.ID, i+1)

		_, err := s.client.SAdd(ctx, key, members).Result()
		if err != nil {
			return fmt.Errorf("[Cache Layer] Failed to cache entrys for member %d: %w", member.ID, err)
		}

		_, err = s.client.Expire(ctx, key, options.expiry).Result()
		if err != nil {
			return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
		}

		// TODO: Figure out if we should bail since we can't cache the key
		// For right now do nothing if it errors
		_, _ = s.client.SAdd(ctx, indexKey, key).Result()

	}

	return nil

}

func chunkTransactions(slc []*athena.MemberWalletTransaction, size int) [][]*athena.MemberWalletTransaction {
	var slcLen = len(slc)
	var output = make([][]*athena.MemberWalletTransaction, 0, (slcLen/size)+1)

	for i := 0; i < slcLen; i += size {
		end := i + size

		if end > slcLen {
			end = slcLen
		}

		output = append(output, slc[i:end])
	}

	return output

}

func chunkJournals(slc []*athena.MemberWalletJournal, size int) [][]*athena.MemberWalletJournal {
	var slcLen = len(slc)
	var output = make([][]*athena.MemberWalletJournal, 0, (slcLen/size)+1)

	for i := 0; i < slcLen; i += size {
		end := i + size

		if end > slcLen {
			end = slcLen
		}

		output = append(output, slc[i:end])
	}

	return output

}
