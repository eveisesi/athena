package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eveisesi/athena"
)

type walletInterface interface {
	GetCharacterWalletBalance(ctx context.Context, member *athena.Member) (float64, *athena.Etag, *http.Response, error)
	HeadCharacterWalletTransactions(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error)
	GetCharacterWalletTransactions(ctx context.Context, member *athena.Member, fromID uint64) ([]*athena.MemberWalletTransaction, *athena.Etag, *http.Response, error)
	HeadCharacterWalletJournals(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error)
	GetCharacterWalletJournals(ctx context.Context, member *athena.Member, page uint) ([]*athena.MemberWalletJournal, *athena.Etag, *http.Response, error)
}

func (s *service) GetCharacterWalletBalance(ctx context.Context, member *athena.Member) (float64, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletBalance]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return 0.00, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return 0.00, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return 0.00, etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", member.ID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return 0.00, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	var balance float64
	err = json.Unmarshal(b, &balance)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return 0.00, nil, nil, err
	}

	etag.Etag = s.retrieveEtagHeader(res.Header)
	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return 0.00, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return balance, etag, res, nil

}

func characterWalletBalanceKeyFunc(mods *modifiers) string {

	requireMember(mods)

	return buildKey(
		GetCharacterWalletBalance.String(),
		strconv.FormatUint(uint64(mods.member.ID), 10),
	)

}

func characterWalletBalancePathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterWalletBalance].Path, mods.member.ID)

}

func (s *service) HeadCharacterWalletTransactions(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletTransactions]

	mods := s.modifiers(ModWithMember(member), ModWithPage(page))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodHead),
		WithPath(path),
		WithPage(page),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", member.ID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterWalletTransactions(ctx context.Context, member *athena.Member, fromID uint64) ([]*athena.MemberWalletTransaction, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletTransactions]

	mods := s.modifiers(ModWithMember(member), ModWithFromID(fromID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)
	reqOpts := append(
		make([]OptionFunc, 0),
		WithMethod(http.MethodGet),
		WithPath(path),
		WithAuthorization(member.AccessToken),
	)
	if fromID > 0 {
		reqOpts = append(reqOpts, WithQuery("from_id", strconv.FormatUint(fromID, 10)))
	}

	b, res, err := s.request(
		ctx,
		reqOpts...,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	transactions := make([]*athena.MemberWalletTransaction, 0, 2500)

	if res.StatusCode >= http.StatusBadRequest {
		return transactions, etag, res, fmt.Errorf("failed to fetch contacts for character %d, received status code of %d", member.ID, res.StatusCode)
	}
	if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return transactions, etag, res, nil
	}

	err = json.Unmarshal(b, &transactions)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	return transactions, etag, nil, nil

}

func characterWalletTransactionsKeyFunc(mods *modifiers) string {

	requireMember(mods)

	param := append(make([]string, 0), GetCharacterWalletTransactions.String(), strconv.FormatUint(uint64(mods.member.ID), 10))

	if mods.from > 0 {
		param = append(param, strconv.FormatUint(mods.from, 10))
	}

	return buildKey(param...)

}

func characterWalletTransactionsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterWalletTransactions].Path, mods.member.ID)

}

func (s *service) HeadCharacterWalletJournals(ctx context.Context, member *athena.Member, page uint) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletJournal]

	mods := s.modifiers(ModWithMember(member), ModWithPage(page))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodHead),
		WithPath(path),
		WithPage(page),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("head request failed, received status code of %d", res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterWalletJournals(ctx context.Context, member *athena.Member, page uint) ([]*athena.MemberWalletJournal, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletJournal]

	mods := s.modifiers(ModWithMember(member), ModWithPage(page))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag),
		WithPage(page),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	journals := make([]*athena.MemberWalletJournal, 0, 2500)

	if res.StatusCode >= http.StatusBadRequest {
		return journals, etag, res, fmt.Errorf("failed to fetch wallet journal for character %d, received status code of %d", member.ID, res.StatusCode)
	}
	if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return journals, etag, res, nil

	}

	err = json.Unmarshal(b, &journals)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		return nil, nil, nil, err
	}

	etag.Etag = s.retrieveEtagHeader(res.Header)

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return journals, etag, res, nil

}

func characterWalletJournalKeyFunc(mods *modifiers) string {

	requireMember(mods)
	requirePage(mods)

	return buildKey(
		GetCharacterWalletJournal.String(),
		strconv.FormatUint(uint64(mods.member.ID), 10),
		strconv.FormatUint(uint64(mods.page), 10),
	)

}

func characterWalletJournalPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterWalletJournal].Path, mods.member.ID)

}
