package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type MemberRepository interface {
	Member(ctx context.Context, id uint) (*Member, error)
	Members(ctx context.Context, operators ...*Operator) ([]*Member, error)
	CreateMember(ctx context.Context, member *Member) (*Member, error)
	UpdateMember(ctx context.Context, id uint, member *Member) (*Member, error)
	DeleteMember(ctx context.Context, id uint) (bool, error)
}

type Member struct {
	ID                uint         `db:"id" json:"id"`
	MainID            null.Uint    `db:"main_id,omitempty" json:"main_id"`
	AccessToken       null.String  `db:"access_token" json:"access_token"`
	RefreshToken      null.String  `db:"refresh_token" json:"refresh_token"`
	Expires           null.Time    `db:"expires," json:"expires"`
	OwnerHash         null.String  `db:"owner_hash" json:"owner_hash"`
	Scopes            MemberScopes `db:"scopes,omitempty" json:"scopes,omitempty"`
	IsNew             bool         `db:"-" json:"-"`
	Disabled          bool         `db:"disabled" json:"disabled"`
	DisabledReason    null.String  `db:"disabled_reason,omitempty" json:"disabled_reason"`
	DisabledTimestamp null.Time    `db:"disabled_timestamp,omitempty" json:"disabled_timestamp"`
	LastLogin         time.Time    `db:"last_login" json:"last_login"`
	CreatedAt         time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time    `db:"updated_at" json:"updated_at"`
}

type MemberLogin struct {
	ID        uint      `db:"_id,omitempty" json:"_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
