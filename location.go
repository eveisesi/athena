package athena

import (
	"context"

	"github.com/volatiletech/null"
)

type MemberLocationRepository interface {
	locationRepository
	onlineRepository
	shipRepository
}

type locationRepository interface {
	MemberLocation(ctx context.Context, id uint) (*MemberLocation, error)
	CreateMemberLocation(ctx context.Context, location *MemberLocation) (*MemberLocation, error)
	UpdateMemberLocation(ctx context.Context, id uint, location *MemberLocation) (*MemberLocation, error)
	DeleteMemberLocation(ctx context.Context, id uint) (bool, error)
}

type onlineRepository interface {
	MemberOnline(ctx context.Context, id uint) (*MemberOnline, error)
	CreateMemberOnline(ctx context.Context, onlien *MemberOnline) (*MemberOnline, error)
	UpdateMemberOnline(ctx context.Context, id uint, onlien *MemberOnline) (*MemberOnline, error)
	DeleteMemberOnline(ctx context.Context, id uint) (bool, error)
}

type shipRepository interface {
	MemberShip(ctx context.Context, id uint) (*MemberShip, error)
	CreateMemberShip(ctx context.Context, ship *MemberShip) (*MemberShip, error)
	UpdateMemberShip(ctx context.Context, id uint, ship *MemberShip) (*MemberShip, error)
	DeleteMemberShip(ctx context.Context, id uint) (bool, error)
}

type MemberLocation struct {
	MemberID      uint        `db:"member_id" json:"member_id"`
	SolarSystemID uint        `db:"solar_system_id" json:"solar_system_id"`
	StationID     null.Uint   `db:"station_id,omitempty" json:"station_id,omitempty"`
	StructureID   null.Uint64 `db:"structure_id,omitempty" json:"structure_id,omitempty"`
	Meta
}

type MemberOnline struct {
	MemberID   uint      `db:"member_id" json:"member_id"`
	LastLogin  null.Time `db:"last_login,omitempty" json:"last_login,omitempty"`
	LastLogout null.Time `db:"last_logout,omitempty" json:"last_logout,omitempty"`
	Logins     null.Uint `db:"logins,omitempty" json:"logins,omitempty"`
	Online     bool      `db:"online,omitempty" json:"online,omitempty"`
	Meta
}

type MemberShip struct {
	MemberID   uint   `db:"member_id" json:"member_id"`
	ShipItemID uint64 `db:"ship_item_id" json:"ship_item_id"`
	ShipName   string `db:"ship_name" json:"ship_name"`
	ShipTypeID uint   `db:"ship_type_id" json:"ship_type_id"`
	Meta
}
