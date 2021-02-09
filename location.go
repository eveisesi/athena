package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type MemberLocationRepository interface {
	locationRepository
	onlineRepository
	shipRepository
}

type locationRepository interface {
	MemberLocation(ctx context.Context, memberID uint) (*MemberLocation, error)
	CreateMemberLocation(ctx context.Context, memberID uint, location *MemberLocation) (*MemberLocation, error)
	UpdateMemberLocation(ctx context.Context, memberID uint, location *MemberLocation) (*MemberLocation, error)
}

type onlineRepository interface {
	MemberOnline(ctx context.Context, memberID uint) (*MemberOnline, error)
	CreateMemberOnline(ctx context.Context, memberID uint, online *MemberOnline) (*MemberOnline, error)
	UpdateMemberOnline(ctx context.Context, memberID uint, online *MemberOnline) (*MemberOnline, error)
}

type shipRepository interface {
	MemberShip(ctx context.Context, memberID uint) (*MemberShip, error)
	CreateMemberShip(ctx context.Context, memberID uint, ship *MemberShip) (*MemberShip, error)
	UpdateMemberShip(ctx context.Context, memberID uint, ship *MemberShip) (*MemberShip, error)
}

type MemberLocation struct {
	MemberID      uint        `db:"member_id" json:"member_id"`
	SolarSystemID uint        `db:"solar_system_id" json:"solar_system_id"`
	StationID     null.Uint   `db:"station_id,omitempty" json:"station_id,omitempty"`
	StructureID   null.Uint64 `db:"structure_id,omitempty" json:"structure_id,omitempty"`
	CreatedAt     time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at"`
}

type MemberOnline struct {
	MemberID   uint      `db:"member_id" json:"member_id"`
	LastLogin  null.Time `db:"last_login,omitempty" json:"last_login,omitempty"`
	LastLogout null.Time `db:"last_logout,omitempty" json:"last_logout,omitempty"`
	Logins     uint      `db:"logins,omitempty" json:"logins,omitempty"`
	Online     bool      `db:"online,omitempty" json:"online,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type MemberShip struct {
	MemberID   uint      `db:"member_id" json:"member_id"`
	ShipItemID uint64    `db:"ship_item_id" json:"ship_item_id"`
	ShipName   string    `db:"ship_name" json:"ship_name"`
	ShipTypeID uint      `db:"ship_type_id" json:"ship_type_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
