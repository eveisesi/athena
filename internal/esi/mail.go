package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/eveisesi/athena"
	"github.com/volatiletech/null"
)

type mailInterface interface {
	GetCharacterMailHeaders(ctx context.Context, member *athena.Member) ([]*MailHeader, *athena.Etag, *http.Response, error)
	GetCharacterMailHeader(ctx context.Context, member *athena.Member, header *MailHeader) (*MailHeader, *athena.Etag, *http.Response, error)
	GetCharacterMailLists(ctx context.Context, member *athena.Member, lists []*athena.MailingList) ([]*athena.MailingList, *athena.Etag, *http.Response, error)
	GetCharacterMailLabels(ctx context.Context, member *athena.Member, labels *athena.MemberMailLabels) (*athena.MemberMailLabels, *athena.Etag, *http.Response, error)
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

func (s *service) GetCharacterMailHeaders(ctx context.Context, member *athena.Member) ([]*MailHeader, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailHeaders]

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

	headers := make([]*MailHeader, 0)

	if res.StatusCode >= http.StatusBadRequest {
		return headers, etag, res, fmt.Errorf("failed to exec mail headers head request for character %d, received status code of %d", member.ID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
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
			WithAuthorization(member.AccessToken),
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
			} else {
				fmt.Println(pageHeaders[len(pageHeaders)-1].Timestamp.Format(time.RFC3339))
			}

		case sc >= http.StatusBadRequest:
			return headers, etag, res, fmt.Errorf("Failed to fetch mail headers from character  %d, received status code of %d", member.ID, sc)
		}

		fmt.Printf("FromID: %d\n", fromID)
		spew.Config.MaxDepth = 2
		for _, header := range headers {
			fmt.Printf("header ID: %d\n", header.MailID.Uint)
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

	requireMember(mods)

	return buildKey(GetCharacterMailHeaders.String(), strconv.FormatUint(uint64(mods.member.ID), 10))

}

func characterMailsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailHeaders].Path, mods.member.ID)

}

func (s *service) GetCharacterMailHeader(ctx context.Context, member *athena.Member, header *MailHeader) (*MailHeader, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailHeader]

	mods := s.modifiers(ModWithMember(member), ModWithMailHeader(header))

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
		WithAuthorization(member.AccessToken),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return header, etag, res, fmt.Errorf("failed to exec mail header head request for character %d with path %s, received status code of %d", member.ID, path, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return header, etag, res, nil
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &header)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		}

	case sc >= http.StatusBadRequest:
		return header, etag, res, fmt.Errorf("Failed to fetch mail header %d from character  %d, received status code of %d", header.MailID.Uint, member.ID, sc)
	}

	return header, etag, res, nil

}

func characterMailKeyFunc(mods *modifiers) string {

	requireMember(mods)
	requireMailHeader(mods)

	return buildKey(GetCharacterMailHeader.String(), strconv.FormatUint(uint64(mods.member.ID), 10), strconv.FormatUint(uint64(mods.header.MailID.Uint), 10))

}

func characterMailPathFunc(mods *modifiers) string {

	requireMember(mods)
	requireMailHeader(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailHeader].Path, mods.member.ID, mods.header.MailID.Uint)

}

func (s *service) GetCharacterMailLists(ctx context.Context, member *athena.Member, lists []*athena.MailingList) ([]*athena.MailingList, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailLists]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.PathFunc(mods))
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return lists, etag, res, fmt.Errorf("failed to exec mailing list head request for character %d, received status code of %d", member.ID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return lists, etag, res, nil
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &lists)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		}
	case sc >= http.StatusBadRequest:
		return lists, etag, res, fmt.Errorf("Failed to fetch mailing lists for character %d, received status code of %d", member.ID, sc)
	}

	return lists, etag, res, nil

}

func characterMailListsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailLists].Path, mods.member.ID)

}

func characterMailListsKeyFunc(mods *modifiers) string {

	requireMember(mods)

	return buildKey(GetCharacterMailLists.String(), strconv.FormatUint(uint64(mods.member.ID), 10))

}

func (s *service) GetCharacterMailLabels(ctx context.Context, member *athena.Member, labels *athena.MemberMailLabels) (*athena.MemberMailLabels, *athena.Etag, *http.Response, error) {

	endpoint := endpoints[GetCharacterMailLabels]

	mods := s.modifiers(ModWithMember(member))

	etag, err := s.etag.Etag(ctx, endpoint.PathFunc(mods))
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return labels, etag, res, fmt.Errorf("failed to exec mail labels head request for character %d, received status code of %d", member.ID, res.StatusCode)
	} else if res.StatusCode == http.StatusNotModified {
		etag.Etag = s.retrieveEtagHeader(res.Header)
		etag.CachedUntil = s.retrieveExpiresHeader(res.Header, 0)
		_, err := s.etag.UpdateEtag(ctx, etag.EtagID, etag)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to update etag after receiving %d: %w", http.StatusNotModified, err)
		}

		return labels, etag, res, nil
	}

	switch sc := res.StatusCode; {
	case sc == http.StatusOK:
		err = json.Unmarshal(b, &labels)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unable to unmarshal response body on request %s: %w", path, err)
		}
	case sc >= http.StatusBadRequest:
		return labels, etag, res, fmt.Errorf("Failed to fetch mailing labels for character %d, received status code of %d", member.ID, sc)
	}

	return labels, etag, res, nil

}

func characterMailLabelsPathFunc(mods *modifiers) string {

	requireMember(mods)

	return buildKey(GetCharacterMailLabels.String(), strconv.FormatUint(uint64(mods.member.ID), 10))

}

func characterMailLabelsKeyFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailLabels].Path, mods.member.ID)

}
