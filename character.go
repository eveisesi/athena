package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharacterRepository interface {
	Characters(ctx context.Context, operators ...*Operator) ([]*Character, error)
	CreateCharacter(ctx context.Context, character *Character) (*Character, error)
	UpdateCharacter(ctx context.Context, id string, character *Character) (*Character, error)
}

type Character struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CharacterID    uint64             `bson:"character_id" json:"character_id"`
	Name           string             `bson:"name" json:"name"`
	CorporationID  uint               `bson:"corporation_id" json:"corporation_id"`
	AllianceID     null.Uint          `bson:"alliance_id,omitempty" json:"alliance_id,omitempty"`
	FactionID      null.Uint          `bson:"faction_id,omitempty" json:"faction_id,omitempty"`
	SecurityStatus null.Float64       `bson:"security_status,omitempty" json:"security_status,omitempty"`
	Gender         string             `bson:"gender" json:"gender"`
	Birthday       time.Time          `bson:"birthday" json:"birthday"`
	Title          string             `bson:"title" json:"title"`
	AncestryID     null.Uint          `bson:"ancestry_id,omitempty" json:"ancestry_id,omitempty"`
	BloodlineID    uint               `bson:"bloodline_id" json:"bloodline_id"`
	RaceID         uint               `bson:"race_id" json:"race_id"`

	Meta
}
