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
}

type mailRepository interface {
	MailHeader(ctx context.Context, mailID uint) (*MailHeader, error)
	CreateMailHeaders(ctx context.Context, headers []*MailHeader) ([]*MailHeader, error)
}

type mailRecipientRepository interface {
	MailRecipients(ctx context.Context, operators ...*Operator) ([]*MailRecipient, error)
	CreateMailRecipients(ctx context.Context, mailID int, recipients []*MailRecipient) ([]*MailRecipient, error)
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

type MemberMailHeader struct {
	MemberID  uint      `db:"member_id" json:"member_id"`
	MailID    uint      `db:"mail_id" json:"mail_id"`
	Labels    SliceUint `db:"labels" json:"-"`
	IsRead    bool      `db:"is_read" json:"is_read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type MailHeader struct {
	MailID     uint             `db:"mail_id" json:"mail_id"`
	From       null.Int         `db:"from,omitempty" json:"from,omitempty"`
	Labels     []uint           `db:"-" json:"labels"`
	Recipients []*MailRecipient `db:"-" json:"recipients"`
	Subject    null.String      `db:"subject" json:"subject"`
	Timestamp  time.Time        `db:"timestamp" json:"timestamp"`
	CreatedAt  time.Time        `db:"created_at" json:"created_at"`
}

type MailRecipient struct {
	MemberID      uint          `db:"member_id" json:"member_id"`
	MailID        uint          `db:"mail_id" json:"mail_id"`
	RecipientID   uint          `db:"recipient_id" json:"recipient_id"`
	RecipientType RecipientType `db:"recipient_type" json:"recipient_type"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updated_at"`
}

type RecipientType string

const (
	RecipientTypeAlliance    RecipientType = "alliance"
	RecipientTypeCharacter   RecipientType = "character"
	RecipientTypeCorporation RecipientType = "corporation"
	RecipientTypeMailingList RecipientType = "mailing_list"
)

type MemberMailLabels struct {
	MemberID         uint         `db:"member_id" json:"member_id"`
	Labels           []*MailLabel `db:"labels,omitempty" json:"labels,omitempty"`
	TotalUnreadCount null.Int     `db:"total_unread_count,omitempty" json:"total_unread_count,omitempty"`
	CreatedAt        time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time    `db:"updated_at" json:"updated_at"`
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
		return nil, fmt.Errorf("[SliceUint] Failed to marshal slice of mail labels for storage in data store: %w", err)
	}

	if string(data) == "null" {
		data = []byte(`[]`)
	}

	return data, err

}
