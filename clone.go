package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type CloneRepository interface {
	jumpCloneRepository
	homeCloneRepository
	implantRepository
}

type jumpCloneRepository interface {
	MemberJumpClones(ctx context.Context, id uint) ([]*MemberJumpClone, error)
	CreateMemberJumpClones(ctx context.Context, id uint, clones []*MemberJumpClone) ([]*MemberJumpClone, error)
	UpdateMemberJumpClones(ctx context.Context, id uint, clones []*MemberJumpClone) ([]*MemberJumpClone, error)
	DeleteMemberJumpClones(ctx context.Context, id uint) (bool, error)
}

type homeCloneRepository interface {
	MemberHomeClone(ctx context.Context, id uint) (*MemberHomeClone, error)
	CreateMemberHomeClone(ctx context.Context, id uint, clone *MemberHomeClone) (*MemberHomeClone, error)
	UpdateMemberHomeClone(ctx context.Context, id uint, clone *MemberHomeClone) (*MemberHomeClone, error)
	DeleteMemberHomeClone(ctx context.Context, id uint) (bool, error)
}

type implantRepository interface {
	MemberImplants(ctx context.Context, id uint) ([]*MemberImplant, error)
	CreateMemberImplants(ctx context.Context, implants []*MemberImplant) ([]*MemberImplant, error)
	UpdateMemberImplants(ctx context.Context, id uint, implants []*MemberImplant) ([]*MemberImplant, error)
	DeleteMemberImplants(ctx context.Context, id uint) (bool, error)
}

type MemberHomeClone struct {
	MemberID              uint      `db:"member_id" json:"member_id"`
	LocationID            uint64    `db:"location_id" json:"location_id"`
	LocationType          string    `db:"location_type" json:"location_type"`
	LastCloneJumpDate     null.Time `db:"last_clone_jump_date" json:"last_clone_jump_date"`
	LastStationChangeDate null.Time `db:"last_station_change_date" json:"last_station_change_date"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`
}

type MemberJumpClone struct {
	MemberID     uint      `db:"member_id" json:"member_id"`
	JumpCloneID  uint      `db:"jump_clone_id" json:"jump_clone_id"`
	LocationID   uint64    `db:"location_id" json:"location_id"`
	LocationType string    `db:"location_type" json:"location_type"`
	Implants     []uint    `db:"implants" json:"implants"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type MemberImplant struct {
	MemberID  uint      `db:"member_id" json:"member_id"`
	TypeID    uint      `db:"type_id" json:"type_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
