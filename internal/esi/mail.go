package esi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/eveisesi/athena"
	"github.com/volatiletech/null"
)

type MailHeader struct {
	From       null.Uint `json:"from,omitempty"`
	IsRead     bool      `json:"is_read"`
	Labels     []uint    `json:"labels"`
	MailID     null.Uint `json:"mail_id"`
	Recipients []struct {
		RecipientID   uint   `json:"recipient_id"`
		RecipientType string `json:"recipient_type"`
	} `json:"recipients"`
	Subject   null.String `json:"subject"`
	Timestamp time.Time   `json:"timestamp"`
}

func (s *service) GetCharacterMail(ctx context.Context, member *athena.Member, mail []*athena.MemberMailHeader) ([]*athena.MemberMailHeader, *athena.Etag, *http.Response, error) {
	return nil, nil, nil, nil
}

func characterMailKeyFunc(mods *modifiers) string {

	requireMember(mods)

	return buildKey(GetCharacterMailHeader.String(), strconv.FormatUint(uint64(mods.member.ID), 10))

}

func characterMailPathFunc(mods *modifiers) string {

	requireMember(mods)

	return fmt.Sprintf(endpoints[GetCharacterMailHeaders].Path, mods.member.ID)

}
