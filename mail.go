package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

// IDEA: Use https://pkg.go.dev/golang.org/x/net/html to replace links in mails to internal resources with links to external sites if such sites exist.
// Evewho for Characters, Corporations, Alliances. Dotlan for Celestial information, etc

type MemberMailRepository interface {
	memberMailRepository
	memberMailRecipientRepository
	memberMailLabelsRepository
}

type memberMailRepository interface {
	MemberMail(ctx context.Context, memberID uint, mailID int) (*MemberMail, error)
	MemberMails(ctx context.Context, memberID uint, operators ...*Operator) ([]*MemberMail, error)
	CreateMemberMail(ctx context.Context, memberID uint, mail *MemberMail) (*MemberMail, error)
	UpdateMemberMail(ctx context.Context, memberID uint, mailID int, mail *MemberMail) (*MemberMail, error)
	DeleteMemberMail(ctx context.Context, memberID uint, mailID int) (bool, error)
	DeleteMemberMails(ctx context.Context, memberID uint) (bool, error)
}

type memberMailRecipientRepository interface {
	MemberMailRecipients(ctx context.Context, memberID uint, mailID int) ([]*MemberMailRecipient, error)
	CreateMemberMailRecipients(ctx context.Context, memberID uint, mailID int, recipients []*MemberMailRecipient) ([]*MemberMailRecipient, error)
	DeleteMemberMailRecipients(ctx context.Context, memberID uint, mailID int) (bool, error)
	DeletememberMailRecipientsAll(ctx context.Context, memberID uint) (bool, error)
}

type memberMailLabelsRepository interface {
	MemberMailLabels(ctx context.Context, memberID uint) ([]*MemberMailLabels, error)
	CreateMemberMailLabel(ctx context.Context, memberID uint, labels *MemberMailLabels) (*MemberMailLabels, error)
	UpdateMemberMailLabel(ctx context.Context, memberID uint, labels *MemberMailLabels) (*MemberMailLabels, error)
	DeleteMemberMailLabels(ctx context.Context, memberID uint) (bool, error)
}

type MemberMail struct {
	ID         uint                  `db:"id" json:"id"`
	From       null.Int              `db:"from,omitempty" json:"from,omitempty"`
	IsRead     bool                  `db:"is_read" json:"is_read"`
	Labels     []int                 `db:"labels" json:"labels"`
	MailID     null.Int              `db:"mail_id" json:"mail_id"`
	Recipients []MemberMailRecipient `db:"-" json:"recipients"`
	Subject    null.String           `db:"subject" json:"subject"`
	Timestamp  time.Time             `db:"timestamp" json:"timestamp"`
	CreatedAt  time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time             `db:"updated_at" json:"updated_at"`
}

type MemberMailRecipient struct {
	MemberID      uint          `db:"member_id" json:"member_id"`
	MailID        int           `db:"mail_id" json:"mail_id"`
	RecipientID   int           `db:"recipient_id" json:"recipient_id"`
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
	MemberID uint `db:"member_id" json:"member_id"`
	Labels   []struct {
		LabelID     null.Int    `db:"label_id,omitempty" json:"label_id,omitempty"`
		Name        null.String `db:"name,omitempty" json:"name,omitempty"`
		UnreadCount null.Int    `db:"unread_count,omitempty" json:"unread_count,omitempty"`
	} `db:"labels,omitempty" json:"labels,omitempty"`
	TotalUnreadCount null.Int  `db:"total_unread_count,omitempty" json:"total_unread_count,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
