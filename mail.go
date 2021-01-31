package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDEA: Use https://pkg.go.dev/golang.org/x/net/html to replace links in mails to internal resources with links to external sites if such sites exist.
// Evewho for Characters, Corporations, Alliances. Dotlan for Celestial information, etc

type MemberMailRepository interface {
	memberMailRepository
	memberMailRecipientRepository
	memberMailLabelsRepository
}

type memberMailRepository interface {
	MemberMail(ctx context.Context, memberID string, mailID int) (*MemberMail, error)
	MemberMails(ctx context.Context, memberID string, operators ...*Operator) ([]*MemberMail, error)
	CreateMemberMail(ctx context.Context, memberID string, mail *MemberMail) (*MemberMail, error)
	UpdateMemberMail(ctx context.Context, memberID string, mailID int, mail *MemberMail) (*MemberMail, error)
	DeleteMemberMail(ctx context.Context, memberID string, mailID int) (bool, error)
	DeleteMemberMails(ctx context.Context, memberID string) (bool, error)
}

type memberMailRecipientRepository interface {
	MemberMailRecipients(ctx context.Context, memberID string, mailID int) ([]*MemberMailRecipient, error)
	CreateMemberMailRecipients(ctx context.Context, memberID string, mailID int, recipients []*MemberMailRecipient) ([]*MemberMailRecipient, error)
	DeleteMemberMailRecipients(ctx context.Context, memberID string, mailID int) (bool, error)
	DeletememberMailRecipientsAll(ctx context.Context, memberID string) (bool, error)
}

type memberMailLabelsRepository interface {
	MemberMailLabels(ctx context.Context, memberID string) ([]*MemberMailLabels, error)
	CreateMemberMailLabel(ctx context.Context, memberID string, labels *MemberMailLabels) (*MemberMailLabels, error)
	UpdateMemberMailLabel(ctx context.Context, memberID string, labels *MemberMailLabels) (*MemberMailLabels, error)
	DeleteMemberMailLabels(ctx context.Context, memberID string) (bool, error)
}

type MemberMail struct {
	MemberID   primitive.ObjectID    `bson:"member_id" json:"member_id"`
	From       null.Int              `bson:"from,omitempty" json:"from,omitempty"`
	IsRead     bool                  `bson:"is_read" json:"is_read"`
	Labels     []int                 `bson:"labels" json:"labels"`
	MailID     null.Int              `bson:"mail_id" json:"mail_id"`
	Recipients []MemberMailRecipient `bson:"recipients" json:"recipients"`
	Subject    null.String           `bson:"subject" json:"subject"`
	Timestamp  time.Time             `bson:"timestamp" json:"timestamp"`
	CreatedAt  time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time             `bson:"updated_at" json:"updated_at"`
}

type MemberMailRecipient struct {
	MemberID      primitive.ObjectID `bson:"member_id" json:"member_id"`
	MailID        int                `bson:"mail_id" json:"mail_id"`
	RecipientID   int                `bson:"recipient_id" json:"recipient_id"`
	RecipientType RecipientType      `bson:"recipient_type" json:"recipient_type"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type RecipientType string

const (
	RecipientTypeAlliance    RecipientType = "alliance"
	RecipientTypeCharacter   RecipientType = "character"
	RecipientTypeCorporation RecipientType = "corporation"
	RecipientTypeMailingList RecipientType = "mailing_list"
)

type MemberMailLabels struct {
	MemberID primitive.ObjectID `bson:"member_id" json:"member_id"`
	Labels   []struct {
		LabelID     null.Int    `bson:"label_id,omitempty" json:"label_id,omitempty"`
		Name        null.String `bson:"name,omitempty" json:"name,omitempty"`
		UnreadCount null.Int    `bson:"unread_count,omitempty" json:"unread_count,omitempty"`
	} `bson:"labels,omitempty" json:"labels,omitempty"`
	TotalUnreadCount null.Int  `bson:"total_unread_count,omitempty" json:"total_unread_count,omitempty"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}
