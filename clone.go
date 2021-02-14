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
	Implants     SliceUint `db:"implants" json:"implants"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type MemberImplant struct {
	MemberID  uint      `db:"member_id" json:"member_id"`
	ImplantID uint      `db:"implant_id" json:"implant_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type SliceUint []uint64

func (s *SliceUint) Scan(value interface{}) error {

	switch data := value.(type) {
	case []byte:
		err := json.Unmarshal(data, s)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s SliceUint) Value() (driver.Value, error) {

	var data []byte
	var err error
	if len(s) == 0 {
		data, err = json.Marshal([]interface{}{})
	} else {
		data, err = json.Marshal(s)
	}
	if err != nil {
		return nil, fmt.Errorf("[SliceUint] Failed to marshal slice of uints for storage in data store: %w", err)
	}

	return data, nil

}

func (s SliceUint) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint64(s))
}

func (s *SliceUint) UnmarshalJSON(value []byte) error {

	x := make([]uint64, 0)
	err := json.Unmarshal(value, &x)
	if err != nil {
		return err
	}
	a := SliceUint(x)
	*s = a
	return nil

}
