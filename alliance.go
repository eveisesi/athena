package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllianceRepository interface {
	Alliances(ctx context.Context, operators ...*Operator) ([]*Alliance, error)
	CreateAlliance(ctx context.Context, alliance *Alliance) (*Alliance, error)
	UpdateAlliance(ctx context.Context, id string, alliance *Alliance) (*Alliance, error)
}

// Alliance is an object representing the database table.
type Alliance struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AllianceID            uint               `bson:"alliance_id" json:"alliance_id"`
	Name                  string             `bson:"name" json:"name"`
	Ticker                string             `bson:"ticker" json:"ticker"`
	DateFounded           time.Time          `bson:"date_founded" json:"date_founded"`
	CreatorID             uint64             `bson:"creator_id" json:"creator_id"`
	CreatorCorporationID  uint               `bson:"creator_corporation_id" json:"creator_corporation_id"`
	ExecutorCorporationID uint               `bson:"executor_corporation_id" json:"executor_corporation_id"`
	FactionID             null.Uint          `bson:"faction_id,omitempty" json:"faction_id,omitempty"`
	IsClosed              bool               `bson:"is_closed" json:"is_closed"`

	Meta
}
