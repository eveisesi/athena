package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/eveisesi/athena"
	"github.com/volatiletech/null"
)

type mailInterface interface {
	GetCharacterMailHeaders(ctx context.Context, characterID uint, token string) ([]*MailHeader, *athena.Etag, *http.Response, error)
	GetCharacterMailHeader(ctx context.Context, characterID, mailID uint, token string) (*MailHeader, *athena.Etag, *http.Response, error)
	GetCharacterMailLists(ctx context.Context, characterID uint, token string) ([]*athena.MailingList, *athena.Etag, *http.Response, error)
	GetCharacterMailLabels(ctx context.Context, characterID uint, token string) (*athena.MemberMailLabels, *athena.Etag, *http.Response, error)
}

type MailHeader struct {
	Body null.String `json:"body"`

	From       null.Uint `json:"from,omitempty"`
	IsRead     bool      `json:"is_read"`
	Labels     []uint64  `json:"labels"`
	MailID     null.Uint `json:"mail_id"`
	Recipients []struct {
		RecipientID   uint                 `json:"recipient_id"`
		RecipientType athena.RecipientType `json:"recipient_type"`
	} `json:"recipients"`
	Subject   null.String `json:"subject"`
	Timestamp time.Time   `json:"timestamp"`
}

func (s *service) GetCharacterMailHeaders(ctx context.Context, characterID uint, token string) ([]*MailHeader, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailHeaders]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.KeyFunc(mods))
	if err != nil {
		return nil, nil, nil, err
	}

	path := endpoint.PathFunc(mods)

	_, res, err := s.request(
		ctx,
		WithMethod(http.MethodHead),
		WithPath(path),
		WithEtag(etag.Etag),
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	var headers = make([]*MailHeader, 0)

	if res.StatusCode >= http.StatusBadRequest {
		return headers, etag, res, fmt.Errorf("failed to exec mail headers head request for character %d, received status code of %d", characterID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = RetrieveEtagHeader(res.Header)
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return headers, etag, res, nil
	}

	fromID := uint64(0)

	ageLimit := time.Now().AddDate(0, -3, 0)

MailHeaderLoop:
	for {

		pageHeaders := make([]*MailHeader, 0, 50)

		pageReqOpts := append(
			make([]OptionFunc, 0),
			WithMethod(http.MethodGet),
			WithPath(path),
			WithAuthorization(token),
		)
		if fromID > 0 {
			pageReqOpts = append(pageReqOpts, WithQuery("last_mail_id", strconv.FormatUint(fromID, 10)))
		}

		b, res, err := s.request(
			ctx,
			pageReqOpts...,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		switch sc := res.StatusCode; {
		case sc == http.StatusOK:
			err = json.Unmarshal(b, &pageHeaders)
			if err != nil {
				err = fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
				return nil, nil, nil, err
			}

			headers = append(headers, pageHeaders...)

			// If the last mail is more than three month old, break out of this loop
			if pageHeaders[len(pageHeaders)-1].Timestamp.Before(ageLimit) {
				break MailHeaderLoop
			}

		case sc >= http.StatusBadRequest:
			return headers, etag, res, fmt.Errorf("Failed to fetch mail headers from character  %d, received status code of %d", characterID, sc)
		}

		if len(pageHeaders) < 50 {
			break
		}

		if !headers[len(headers)-1].MailID.Valid {
			break
		}

		fromID = uint64(headers[len(headers)-1].MailID.Uint)

		time.Sleep(time.Second)
	}

	return headers, etag, res, nil

}

func characterMailsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterMailHeaders.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func characterMailsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailHeaders].Path, mods.characterID)

}

func (s *service) GetCharacterMailHeader(ctx context.Context, characterID, mailID uint, token string) (*MailHeader, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailHeader]

	mods := s.modifiers(ModWithCharacterID(characterID), ModWithMailID(mailID))

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
		WithAuthorization(token),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	var header = new(MailHeader)

	if res.StatusCode >= http.StatusBadRequest {
		return header, etag, res, fmt.Errorf("failed to exec mail header head request for character %d with path %s, received status code of %d", characterID, path, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = RetrieveEtagHeader(res.Header)
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return header, etag, res, nil
	}

	err = json.Unmarshal(b, &header)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return header, etag, res, nil

}

func characterMailKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)
	requireMailID(mods)

	return buildKey(GetCharacterMailHeader.String(), strconv.FormatUint(uint64(mods.characterID), 10), strconv.FormatUint(uint64(mods.mailID), 10))

}

func characterMailPathFunc(mods *modifiers) string {

	requireCharacterID(mods)
	requireMailID(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailHeader].Path, mods.characterID, mods.mailID)

}

func (s *service) GetCharacterMailLists(ctx context.Context, characterID uint, token string) ([]*athena.MailingList, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailLists]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.PathFunc(mods))
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}

	var lists = make([]*athena.MailingList, 0, 250)

	if res.StatusCode >= http.StatusBadRequest {
		return lists, etag, res, fmt.Errorf("failed to exec mailing list head request for character %d, received status code of %d", characterID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = RetrieveEtagHeader(res.Header)
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return lists, etag, res, nil
	}

	err = json.Unmarshal(b, &lists)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return lists, etag, res, nil

}

func characterMailListsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailLists].Path, mods.characterID)

}

func characterMailListsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterMailLists.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}

func (s *service) GetCharacterMailLabels(ctx context.Context, characterID uint, token string) (*athena.MemberMailLabels, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailLabels]

	mods := s.modifiers(ModWithCharacterID(characterID))

	etag, err := s.etag.Etag(ctx, endpoint.PathFunc(mods))
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}

	var labels = new(athena.MemberMailLabels)

	if res.StatusCode >= http.StatusBadRequest {
		return labels, etag, res, fmt.Errorf("failed to exec mail labels head request for character %d, received status code of %d", characterID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = RetrieveEtagHeader(res.Header)
		etag.CachedUntil = RetrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return labels, etag, res, nil
	}

	err = json.Unmarshal(b, labels)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
	}

	return labels, etag, res, nil

}

func characterMailLabelsPathFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailLabels].Path, mods.characterID)

}

func characterMailLabelsKeyFunc(mods *modifiers) string {

	requireCharacterID(mods)

	return buildKey(GetCharacterMailLabels.String(), strconv.FormatUint(uint64(mods.characterID), 10))

}
