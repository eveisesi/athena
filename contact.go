package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type MemberContactRepository interface {
	memberContactRepository
	memberContactLabelRepository
}

type memberContactRepository interface {
	MemberContact(ctx context.Context, memberID uint, contactID int) (*MemberContact, error)
	MemberContacts(ctx context.Context, memberID uint) ([]*MemberContact, error)
	CreateMemberContacts(ctx context.Context, memberID uint, contacts []*MemberContact) ([]*MemberContact, error)
	UpdateMemberContact(ctx context.Context, memberID uint, contact *MemberContact) (*MemberContact, error)
	DeleteMemberContacts(ctx context.Context, memberID uint, contacts []*MemberContact) (bool, error)
}

type memberContactLabelRepository interface {
	MemberContactLabels(ctx context.Context, memberID uint) ([]*MemberContactLabel, error)
	CreateMemberContactLabels(ctx context.Context, memberID uint, labels []*MemberContactLabel) ([]*MemberContactLabel, error)
	UpdateMemberContactLabel(ctx context.Context, memberID uint, label *MemberContactLabel) (*MemberContactLabel, error)
	DeleteMemberContactLabels(ctx context.Context, memberID uint, labels []*MemberContactLabel) (bool, error)
}

type MemberContact struct {
	MemberID    uint      `db:"member_id" json:"member_id" deep:"-"`
	ContactID   uint      `db:"contact_id" json:"contact_id"`
	ContactType string    `db:"contact_type" json:"contact_type"`
	IsBlocked   null.Bool `db:"is_blocked,omitempty" json:"is_blocked,omitempty"`
	IsWatched   null.Bool `db:"is_watched,omitempty" json:"is_watched,omitempty"`
	LabelIDs    []uint64  `db:"label_ids" json:"label_ids"`
	Standing    float64   `db:"standing" json:"standing"`
	CreatedAt   time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}

type MemberContactLabel struct {
	MemberID  uint      `db:"member_id" json:"member_id" deep:"-"`
	LabelID   uint64    `db:"label_id" json:"label_id"`
	LabelName string    `db:"label_name" json:"label_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}
