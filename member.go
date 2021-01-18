package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberRepository interface {
	Members(ctx context.Context, operators ...*Operator) ([]*Member, error)
	CreateMember(ctx context.Context, member *Member) (*Member, error)
	UpdateMember(ctx context.Context, id string, member *Member) (*Member, error)
	DeleteMember(ctx context.Context, id string) (bool, error)
}

type Member struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CharacterID       uint64             `bson:"character_id" json:"character_id"`
	MainID            null.Uint64        `bson:"main_id,omitempty" json:"main_id"`
	AccessToken       string             `bson:"access_token" json:"access_token"`
	RefreshToken      string             `bson:"refresh_token" json:"refresh_token"`
	Expires           time.Time          `bson:"expires" json:"expires"`
	OwnerHash         string             `bson:"owner_hash" json:"owner_hash"`
	Scopes            []MemberScope      `bson:"scopes,omitempty" json:"scopes,omitempty"`
	IsNew             bool               `bson:"is_new" json:"is_new"`
	Disabled          bool               `bson:"disabled" json:"disabled"`
	DisabledReason    null.String        `bson:"disabled_reason,omitempty" json:"disabled_reason"`
	DisabledTimestamp null.Time          `bson:"disabled_timestamp,omitempty" json:"disabled_timestamp"`
	LastLogin         time.Time          `bson:"last_login" json:"last_login"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

type MemberScope struct {
	Scope  Scope     `bson:"scope" json:"scope"`
	Expiry null.Time `bson:"expiry,omitempty" json:"expiry,omitempty"`
}

type MemberLogin struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
