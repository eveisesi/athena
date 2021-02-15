package athena

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/volatiletech/null"
)

// IDEA: Use https://pkg.go.dev/golang.org/x/net/html to replace links in mails to internal resources with links to external sites if such sites exist.
// Evewho for Characters, Corporations, Alliances. Dotlan for Celestial information, etc

type MailRepository interface {
	mailRepository
	memberMailRepository
	mailRecipientRepository
	memberMailLabelsRepository
	memberMailingListRepository
	mailingListRepository
}

type mailRepository interface {
	MailHeader(ctx context.Context, mailID uint) (*MailHeader, error)
	CreateMailHeaders(ctx context.Context, headers []*MailHeader) ([]*MailHeader, error)
}

type mailRecipientRepository interface {
	MailRecipients(ctx context.Context, operators ...*Operator) ([]*MailRecipient, error)
	CreateMailRecipients(ctx context.Context, recipients []*MailRecipient) ([]*MailRecipient, error)
}

type memberMailRepository interface {
	MemberMailHeaders(ctx context.Context, operators ...*Operator) ([]*MemberMailHeader, error)
	CreateMemberMailHeaders(ctx context.Context, memberID uint, headers []*MemberMailHeader) ([]*MemberMailHeader, error)
	UpdateMemberMailHeaders(ctx context.Context, memberID uint, headers []*MemberMailHeader) ([]*MemberMailHeader, error)
}

type memberMailLabelsRepository interface {
	MemberMailLabels(ctx context.Context, memberID uint) (*MemberMailLabels, error)
	CreateMemberMailLabel(ctx context.Context, memberID uint, labels *MemberMailLabels) (*MemberMailLabels, error)
	UpdateMemberMailLabel(ctx context.Context, memberID uint, labels *MemberMailLabels) (*MemberMailLabels, error)
}

type memberMailingListRepository interface {
	MemberMailingLists(ctx context.Context, memberID uint) ([]*MemberMailingList, error)
	CreateMemberMailingLists(ctx context.Context, memberID uint, lists []*MemberMailingList) ([]*MemberMailingList, error)
	DeleteMemberMailingListsAll(ctx context.Context, memberID uint) (bool, error)
}

type mailingListRepository interface {
	MailingList(ctx context.Context, mailingListID uint) (*MailingList, error)
	MailingLists(ctx context.Context, operators ...*Operator) ([]*MailingList, error)
	CreateMailingLists(ctx context.Context, lists []*MailingList) ([]*MailingList, error)
	UpdateMailingList(ctx context.Context, mailingListID uint, list *MailingList) (*MailingList, error)
}

type MemberMailHeader struct {
	MemberID  uint      `db:"member_id" json:"member_id"`
	MailID    uint      `db:"mail_id" json:"mail_id"`
	Labels    SliceUint `db:"labels" json:"-"`
	IsRead    bool      `db:"is_read" json:"is_read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type MailHeader struct {
	MailID     uint             `db:"id" json:"mail_id"`
	Sender     null.Uint        `db:"sender_id" json:"from,omitempty"`
	SenderType null.String      `db:"sender_type" json:"from_type,omitempty"`
	Recipients []*MailRecipient `db:"-" json:"recipients"`
	Subject    null.String      `db:"subject" json:"subject"`
	Body       null.String      `db:"body" json:"body"`
	Timestamp  time.Time        `db:"sent" json:"timestamp"`
	CreatedAt  time.Time        `db:"created_at" json:"created_at"`
}

type MailRecipient struct {
	MailID        uint          `db:"mail_id" json:"mail_id"`
	RecipientID   uint          `db:"recipient_id" json:"recipient_id"`
	RecipientType RecipientType `db:"recipient_type" json:"recipient_type"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
}

type RecipientType string

const (
	RecipientTypeAlliance    RecipientType = "alliance"
	RecipientTypeCharacter   RecipientType = "character"
	RecipientTypeCorporation RecipientType = "corporation"
	RecipientTypeMailingList RecipientType = "mailing_list"
)

var AllRecipientTypes = []RecipientType{
	RecipientTypeAlliance, RecipientTypeCharacter,
	RecipientTypeCorporation, RecipientTypeMailingList,
}

func (r RecipientType) IsValid() bool {

	for _, rType := range AllRecipientTypes {
		if r == rType {
			return true
		}
	}

	return false

}

type MemberMailLabels struct {
	MemberID         uint       `db:"member_id" json:"member_id"`
	Labels           MailLabels `db:"labels,omitempty" json:"labels,omitempty"`
	TotalUnreadCount null.Int   `db:"total_unread_count,omitempty" json:"total_unread_count,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

type MailLabel struct {
	LabelID     null.Int    `db:"label_id,omitempty" json:"label_id,omitempty"`
	Name        null.String `db:"name,omitempty" json:"name,omitempty"`
	UnreadCount null.Int    `db:"unread_count,omitempty" json:"unread_count,omitempty"`
}

type MailLabels []MailLabel

func (l *MailLabels) Scan(value interface{}) error {

	switch data := value.(type) {
	case []byte:
		err := json.Unmarshal(data, l)
		if err != nil {
			return err
		}
	}

	return nil

}

func (l MailLabels) Value() (driver.Value, error) {

	data, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("[MailLabels] Failed to marshal slice of mail labels for storage in data store: %w", err)
	}

	if string(data) == "null" {
		data = []byte(`[]`)
	}

	return data, err

}

type MemberMailingList struct {
	MemberID      uint      `db:"member_id" json:"member_id"`
	MailingListID uint      `db:"mailing_list_id" json:"mailing_list_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type MailingList struct {
	MailingListID uint      `db:"mailing_list_id" json:"mailing_list_id"`
	Name          string    `db:"name" json:"name"`
	CreatedAt     time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}
