package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberLocationRepository interface {
	locationRepository
	onlineRepository
	shipRepository
}

type locationRepository interface {
	MemberLocation(ctx context.Context, id string) (*MemberLocation, error)
	CreateMemberLocation(ctx context.Context, location *MemberLocation) (*MemberLocation, error)
	UpdateMemberLocation(ctx context.Context, id string, location *MemberLocation) (*MemberLocation, error)
	DeleteMemberLocation(ctx context.Context, id string) error
}

type onlineRepository interface {
	MemberOnline(ctx context.Context, id string) (*MemberOnline, error)
	CreateMemberOnline(ctx context.Context, onlien *MemberOnline) (*MemberOnline, error)
	UpdateMemberOnline(ctx context.Context, id string, onlien *MemberOnline) (*MemberOnline, error)
	DeleteMemberOnline(ctx context.Context, id string) error
}

type shipRepository interface {
	MemberShip(ctx context.Context, id string) (*MemberShip, error)
	CreateMemberShip(ctx context.Context, ship *MemberShip) (*MemberShip, error)
	UpdateMemberShip(ctx context.Context, id string, ship *MemberShip) (*MemberShip, error)
	DeleteMemberShip(ctx context.Context, id string) error
}

type MemberLocation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	MemberID      primitive.ObjectID `bson:"member_id" json:"member_id"`
	SolarSystemID uint               `bson:"solar_system_id" json:"solar_system_id"`
	StationID     null.Uint          `bson:"station_id,omitempty" json:"station_id,omitempty"`
	StructureID   null.Uint64        `bson:"structure_id,omitempty" json:"structure_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type MemberOnline struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	MemberID   primitive.ObjectID `bson:"member_id" json:"member_id"`
	LastLogin  null.Time          `bson:"last_login,omitempty" json:"last_login,omitempty"`
	LastLogout null.Time          `bson:"last_logout,omitempty" json:"last_logout,omitempty"`
	Logins     null.Uint          `bson:"logins,omitempty" json:"logins,omitempty"`
	Online     bool               `bson:"online,omitempty" json:"online,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type MemberShip struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	MemberID   primitive.ObjectID `bson:"member_id" json:"member_id"`
	ShipItemID uint64             `bson:"ship_item_id" json:"ship_item_id"`
	ShipName   string             `bson:"ship_name" json:"ship_name"`
	ShipTypeID uint               `bson:"ship_type_id" json:"ship_type_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
