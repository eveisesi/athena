package athena

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/volatiletech/null"
)

type CloneRepository interface {
	MemberClones(ctx context.Context, id uint) (*MemberClones, error)
	CreateMemberClones(ctx context.Context, clones *MemberClones) (*MemberClones, error)
	UpdateMemberClones(ctx context.Context, clones *MemberClones) (*MemberClones, error)
	DeleteMemberClones(ctx context.Context, id uint) (bool, error)

	MemberImplants(ctx context.Context, id uint) ([]*MemberImplant, error)
	CreateMemberImplants(ctx context.Context, id uint, implants []*MemberImplant) ([]*MemberImplant, error)
	DeleteMemberImplants(ctx context.Context, id uint) (bool, error)
}

type MemberClones struct {
	MemberID              uint                `db:"member_id" json:"member_id"`
	HomeLocation          *MemberHomeLocation `json:"home_location"`
	JumpClones            []*MemberJumpClone  `json:"jump_clones"`
	LastCloneJumpDate     null.Time           `db:"last_clone_jump_date" json:"last_clone_jump_date"`
	LastStationChangeDate null.Time           `db:"last_station_change_date" json:"last_station_change_date"`
	CreatedAt             time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time           `db:"updated_at" json:"updated_at"`
}

type MemberHomeLocation struct {
	LocationID   uint64    `db:"location_id" json:"location_id"`
	LocationType string    `db:"location_type" json:"location_type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type MemberJumpClone struct {
	JumpCloneID  uint      `db:"jump_clone_id" json:"jump_clone_id"`
	LocationID   uint64    `db:"location_id" json:"location_id"`
	LocationType string    `db:"location_type" json:"location_type"`
	Implants     UintSlice `db:"implants" json:"implants"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type MemberImplant struct {
	MemberID  uint      `db:"member_id" json:"member_id"`
	ImplantID uint      `db:"implant_id" json:"implant_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UintSlice []uint

func (s *UintSlice) Scan(value interface{}) error {

	switch data := value.(type) {
	case []byte:
		var values UintSlice
		err := json.Unmarshal(data, &values)
		if err != nil {
			return err
		}

		*s = values
	}

	return nil

}

func (s UintSlice) Value() (driver.Value, error) {

	var data []byte
	var err error
	if len(s) == 0 {
		data, err = json.Marshal([]interface{}{})
	} else {
		data, err = json.Marshal(s)
	}
	if err != nil {
		return nil, fmt.Errorf("[UintSlice] Failed to marshal slice of string for storage in data store: %w", err)
	}

	return data, nil

}
