package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/eveisesi/athena"
)

type walletInterface interface {
	GetCharacterWalletBalance(ctx context.Context, member *athena.Member) (float64, *athena.Etag, *http.Response, error)
	GetCharacterWalletTransactions(ctx context.Context, member *athena.Member, transactions []*athena.MemberWalletTransaction) ([]*athena.MemberWalletTransaction, *athena.Etag, *http.Response, error)
	GetCharacterWalletJournals(ctx context.Context, member *athena.Member, journals []*athena.MemberWalletJournal) ([]*athena.MemberWalletJournal, *athena.Etag, *http.Response, error)
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

	var balance float64

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &balance)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
			return 0.00, nil, nil, err
		}

		etag.Etag = s.retrieveEtagHeader(res.Header)

	case sc >= http.StatusBadRequest:
		return balance, etag, res, fmt.Errorf("failed to fetch balance for character %d, received status code of %d", member.ID, sc)
	}

	etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
	_, err = s.etag.UpdateEtag(ctx, etag.EtagID, etag)
	if err != nil {
		return 0.00, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
	}

	return balance, etag, res, nil

}

func characterWalletBalanceKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	param := append(make([]string, 0), GetCharacterWalletBalance.String(), strconv.Itoa(int(mods.member.ID)))

	if mods.page != nil {
		param = append(param, strconv.Itoa(*mods.page))
	}

	return buildKey(param...)

}

func characterWalletBalancePathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterWalletBalance].Path, mods.member.ID),
	}

	return u.String()

}

func (s *service) GetCharacterWalletTransactions(ctx context.Context, member *athena.Member, transactions []*athena.MemberWalletTransaction) ([]*athena.MemberWalletTransaction, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletTransactions]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodHead),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return transactions, etag, res, fmt.Errorf("failed to fetch contacts for character %d, received status code of %d", member.ID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return transactions, etag, res, nil
	}

	fromID := uint64(0)

	for {

		pageTransactions := make([]*athena.MemberWalletTransaction, 0)

		path := endpoint.PathFunc(mods)

		reqOpts := append(make([]OptionFunc, 0), WithMethod(http.MethodGet), WithPath(path), WithAuthorization(member.AccessToken))
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

		switch sc := res.StatusCode; {
		case sc == http.StatusOK:
			err = json.Unmarshal(b, &pageTransactions)
			if err != nil {
				err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
				return nil, nil, nil, err
			}

			if len(pageTransactions) == 0 {
				goto EndOfLoop
			}

			transactions = append(transactions, pageTransactions...)

		case sc >= http.StatusBadRequest:
			return transactions, etag, res, fmt.Errorf("failed to fetch transactions for character %d, received status code of %d", member.ID, sc)
		}

		// Grab the transaction from the last ID. Subtract 1 and set it to From ID
		fromID = transactions[len(transactions)-1].TransactionID - 1
		time.Sleep(time.Second)
	}
EndOfLoop:

	return transactions, etag, nil, nil

}

func characterWalletTransactionsKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	param := append(make([]string, 0), GetCharacterWalletTransactions.String(), strconv.Itoa(int(mods.member.ID)))

	if mods.page != nil {
		param = append(param, strconv.Itoa(*mods.page))
	}

	return buildKey(param...)

}

func characterWalletTransactionsPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterWalletTransactions].Path, mods.member.ID),
	}

	return u.String()

}

func (s *service) GetCharacterWalletJournals(ctx context.Context, member *athena.Member, journals []*athena.MemberWalletJournal) ([]*athena.MemberWalletJournal, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterWalletJournal]

	pages := 1

	mods := s.modifiers(ModWithMember(member), ModWithPage(&pages))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodHead),
		WithPath(path),
		WithEtag(etag),
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return journals, etag, res, fmt.Errorf("failed to fetch contacts for character %d, received status code of %d", member.ID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return journals, etag, res, nil

	}

	pages = s.retrieveXPagesFromHeader(res.Header)
	if pages == 0 {
		return nil, nil, nil, fmt.Errorf("received 0 for X-Pages on request %s, expected number greater than 0", path)
	}

	for i := 1; i <= pages; i++ {

		pageJournal := make([]*athena.MemberWalletJournal, 0)

		mods := s.modifiers(ModWithMember(member), ModWithPage(&i))

		pageEtag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
		if err != nil {
			return nil, nil, nil, err
		}

		path := endpoint.PathFunc(mods)

		b, res, err := s.request(
			ctx,
			WithMethod(http.MethodGet),
			WithPath(path),
			WithEtag(pageEtag),
			WithPage(i),
			WithAuthorization(member.AccessToken),
		)
		if err != nil {
			return nil, nil, nil, err
		}

		switch sc := res.StatusCode; {
		case sc == http.StatusOK:
			err = json.Unmarshal(b, &pageJournal)
			if err != nil {
				err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
				return nil, nil, nil, err
			}

			journals = append(journals, pageJournal...)

			pageEtag.Etag = s.retrieveEtagHeader(res.Header)

		case sc >= http.StatusBadRequest:
			return journals, etag, res, fmt.Errorf("failed to fetch journals for character %d, received status code of %d", member.ID, sc)
		}

		pageEtag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err = s.etag.UpdateEtag(ctx, pageEtag.EtagID, pageEtag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", res.StatusCode, err)
		}

	}

	return journals, etag, res, nil

}

func characterWalletJournalKeyFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	param := append(make([]string, 0), GetCharacterWalletJournal.String(), strconv.Itoa(int(mods.member.ID)))

	if mods.page != nil {
		param = append(param, strconv.Itoa(*mods.page))
	}

	return buildKey(param...)

}

func characterWalletJournalPathFunc(mods *modifiers) string {

	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil for member instead")
	}

	u := url.URL{
		Path: fmt.Sprintf(endpoints[GetCharacterWalletJournal].Path, mods.member.ID),
	}

	return u.String()

}
