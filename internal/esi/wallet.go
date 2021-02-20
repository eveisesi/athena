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
	HeadCharacterWalletBalance(ctx context.Context, characterID uint, token string) (*athena.Etag, *http.Response, error)
	GetCharacterWalletBalance(ctx context.Context, characterID uint, token string) (float64, *athena.Etag, *http.Response, error)
	HeadCharacterWalletTransactions(ctx context.Context, characterID uint, from uint64, token string) (*athena.Etag, *http.Response, error)
	GetCharacterWalletTransactions(ctx context.Context, characterID uint, fromID uint64, token string) ([]*athena.MemberWalletTransaction, *athena.Etag, *http.Response, error)
	HeadCharacterWalletJournals(ctx context.Context, characterID, page uint, token string) (*athena.Etag, *http.Response, error)
	GetCharacterWalletJournals(ctx context.Context, characterID, page uint, token string) ([]*athena.MemberWalletJournal, *athena.Etag, *http.Response, error)
}

func (s *service) HeadCharacterWalletBalance(ctx context.Context, characterID uint, token string) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletBalance]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("failed to make head request to contracts for character %d, received status code of %d", characterID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterWalletBalance(ctx context.Context, characterID uint, token string) (float64, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletBalance]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return 0.00, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
		WithAuthorization(token),
	)
	if err != nil {
		return 0.00, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return 0.00, etag, res, fmt.Errorf("failed to fetch contracts for character %d, received status code of %d", characterID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
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

	etag.Etag = RetrieveEtagHeader(res.Header)
	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return 0.00, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return balance, etag, res, nil

}

func characterWalletBalanceKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(
		GetCharacterWalletBalance.String(),
		strconv.FormatUint(uint64(mods.characterID), 10),
	)

}

func characterWalletBalancePathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterWalletBalance].Path, mods.characterID)

}

func (s *service) HeadCharacterWalletTransactions(ctx context.Context, characterID uint, fromID uint64, token string) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletTransactions]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithFromID(fromID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, err
	}

	path := endpoint.PathFunc(mods)
	reqOpts := append(
		make([]OptionFunc, 0),
		WithMethod(http.MethodHead),
		WithPath(path),
		WithAuthorization(token),
	)
	if fromID > 0 {
		reqOpts = append(reqOpts, WithQuery("from_id", strconv.FormatUint(fromID, 10)))
	}

	_, res, err := s.request(
		ctx,
		reqOpts...,
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("failed to exec head request to character wallet transactions for character %d, received status code of %d", characterID, res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterWalletTransactions(ctx context.Context, characterID uint, fromID uint64, token string) ([]*athena.MemberWalletTransaction, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletTransactions]

	modFuncs := append(make([]modifierFunc, 0, 2), ModWithCharacterID(characterID))
	if fromID > 0 {
		modFuncs = append(modFuncs, ModWithFromID(fromID))
	}

	mods := s.modifiers(modFuncs...)

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)
	reqOpts := append(
		make([]OptionFunc, 0, 6),
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
		WithAuthorization(token),
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

	var transactions = make([]*athena.MemberWalletTransaction, 0, 2500)

	if res.StatusCode >= http.StatusBadRequest {
		return transactions, etag, res, fmt.Errorf("failed to fetch contacts for character %d, received status code of %d", characterID, res.StatusCode)
	}

	etag.Etag = RetrieveEtagHeader(res.Header)
	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
	}

	if res.StatusCode == http.StatusNotModified {
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

	requireCharacterID(mods)

	param := append(make([]string, 0), GetCharacterWalletTransactions.String(), strconv.FormatUint(uint64(mods.characterID), 10))

	if mods.from > 0 {
		param = append(param, strconv.FormatUint(mods.from, 10))
	}

	return buildKey(param...)

}

func characterWalletTransactionsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterWalletTransactions].Path, mods.characterID)

}

func (s *service) HeadCharacterWalletJournals(ctx context.Context, characterID, page uint, token string) (*athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletJournal]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithPage(page))

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
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return etag, res, fmt.Errorf("head request failed, received status code of %d", res.StatusCode)
	}

	if res.StatusCode == http.StatusNotModified {
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}
	}

	return etag, res, nil

}

func (s *service) GetCharacterWalletJournals(ctx context.Context, characterID, page uint, token string) ([]*athena.MemberWalletJournal, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletJournal]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithPage(page))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	b, res, err := s.request(
		ctx,
		WithMethod(http.MethodGet),
		WithPath(path),
		WithEtag(etag.Etag),
		WithPage(page),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	journals := make([]*athena.MemberWalletJournal, 0, 2500)

	if res.StatusCode >= http.StatusBadRequest {
		return journals, etag, res, fmt.Errorf("failed to fetch wallet journal for character %d, received status code of %d", characterID, res.StatusCode)
	}
	if res.StatusCode == http.StatusNotModified {
		etag.Etag = RetrieveEtagHeader(res.Header)
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
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

	etag.Etag = RetrieveEtagHeader(res.Header)

	etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return journals, etag, res, nil

}

func characterWalletJournalKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)
	requirePage(mods)

	return buildKey(
		GetCharacterWalletJournal.String(),
		strconv.FormatUint(uint64(mods.characterID), 10),
		strconv.FormatUint(uint64(mods.page), 10),
	)

}

func characterWalletJournalPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterWalletJournal].Path, mods.characterID)

}
