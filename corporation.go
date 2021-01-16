package athena

import (
	"context"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CorporationRespository interface {
	Corporations(ctx context.Context, operators ...*Operator) ([]*Corporation, error)
	CreateCorporation(ctx context.Context, corporation *Corporation) error
	UpdateCorporation(ctx context.Context, id uint, corporation *Corporation) error
}

type Corporation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CorporationID uint               `bson:"corporation_id" json:"corporation_id"`

	AllianceID    null.Uint   `bson:"alliance_id,omitempty" json:"alliance_id,omitempty"`
	CeoID         uint        `bson:"ceo_id" json:"ceo_id"`
	CreatorID     uint        `bson:"creator_id" json:"creator_id"`
	DateFounded   null.Time   `bson:"date_founded,omitempty" json:"date_founded,omitempty"`
	FactionID     null.Uint   `bson:"faction_id,omitempty" json:"faction_id,omitempty"`
	HomeStationID null.Uint   `bson:"home_station_id,omitempty" json:"home_station_id,omitempty"`
	MemberCount   uint        `bson:"member_count" json:"member_count"`
	Name          string      `bson:"name" json:"name"`
	Shares        uint64      `bson:"shares,omitempty" json:"shares,omitempty"`
	TaxRate       float32     `bson:"tax_rate" json:"tax_rate"`
	Ticker        string      `bson:"ticker" json:"ticker"`
	URL           null.String `bson:"url,omitempty" json:"url,omitempty"`
	WarEligible   bool        `bson:"war_eligible" json:"war_eligible"`

	Meta
}
