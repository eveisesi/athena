package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberContactRepository interface {
	memberContactRepository
	memberContactLabelRepository
}

type memberContactRepository interface {
	MemberContact(ctx context.Context, memberID string, contactID int) (*MemberContact, error)
	MemberContacts(ctx context.Context, memberID string) ([]*MemberContact, error)

	CreateMemberContacts(ctx context.Context, memberID string, contacts []*MemberContact) ([]*MemberContact, error)

	UpdateMemberContact(ctx context.Context, memberID string, contact *MemberContact) (*MemberContact, error)
	UpdateMemberContacts(ctx context.Context, memberID string, contacts []*MemberContact) ([]*MemberContact, error)

	DeleteMemberContact(ctx context.Context, memberID string, contactID int) (bool, error)
	DeleteMemberContacts(ctx context.Context, memberID string, contacts []*MemberContact) (bool, error)
}

type memberContactLabelRepository interface {
	MemberContactLabels(ctx context.Context, memberID string) ([]*MemberContactLabel, error)
	CreateMemberContactLabels(ctx context.Context, memberID string, labels []*MemberContactLabel) ([]*MemberContactLabel, error)
	UpdateMemberContactLabels(ctx context.Context, memberID string, labels []*MemberContactLabel) ([]*MemberContactLabel, error)
	DeleteMemberContactLabels(ctx context.Context, memberID string, labels []*MemberContactLabel) (bool, error)
}

type MemberContact struct {
	MemberID    primitive.ObjectID `bson:"member_id" json:"member_id" deep:"-"`
	ContactID   int                `bson:"contact_id" json:"contact_id"`
	ContactType string             `bson:"contact_type" json:"contact_type"`
	IsBlocked   null.Bool          `bson:"is_blocked,omitempty" json:"is_blocked,omitempty"`
	IsWatched   null.Bool          `bson:"is_watched,omitempty" json:"is_watched,omitempty"`
	LabelIDs    []int64            `bson:"label_ids" json:"label_ids"`
	Standing    float64            `bson:"standing" json:"standing"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at" deep:"-"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at" deep:"-"`
}

type MemberContactLabel struct {
	MemberID  primitive.ObjectID `bson:"member_id" json:"member_id" deep:"-"`
	LabelID   int64              `bson:"label_id" json:"label_id"`
	LabelName string             `bson:"label_name" json:"label_name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at" deep:"-"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at" deep:"-"`
}
