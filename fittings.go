package athena

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberFittingsRepository interface {
	memberFittingRepository
	memberFittingItemRepository
}

type memberFittingRepository interface {
	MemberFittings(ctx context.Context, memberID string, operators ...*Operator) ([]*MemberFitting, error)
	CreateMemberFittings(ctx context.Context, memberID string, fitting []*MemberFitting) ([]*MemberFitting, error)
	UpdateMemberFitting(ctx context.Context, memberID string, fittingID uint, fitting *MemberFitting) (*MemberFitting, error)
	DeleteMemberFitting(ctx context.Context, memberID string, fittingID uint) (bool, error)
	DeleteMemberFittings(ctx context.Context, memberID string) (bool, error)
}

type memberFittingItemRepository interface {
	MemberFittingItems(ctx context.Context, memberID string, fittingID uint) ([]*MemberFittingItem, error)
	CreateMemberFittingItems(ctx context.Context, memberID string, fittingID uint, items []*MemberFittingItem) ([]*MemberFittingItem, error)
	DeleteMemberFittingItems(ctx context.Context, memberID string, fittingID uint) (bool, error)
}

type MemberFitting struct {
	MemberID    primitive.ObjectID   `bson:"member_id" json:"member_id"`
	FittingID   uint                 `bson:"fitting_id" json:"fitting_id"`
	ShipTypeID  uint                 `bson:"ship_type_id" json:"ship_type_id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Items       []*MemberFittingItem `bson:"-" json:"items,omitempty"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type MemberFittingItem struct {
	MemberID  primitive.ObjectID `bson:"member_id" json:"member_id"`
	FittingID uint               `bson:"fitting_id" json:"fitting_id"`
	TypeID    uint               `bson:"type_id" json:"type_id"`
	Quantity  uint               `bson:"quantity" json:"quantity"`
	Flag      FittingItemFlag    `bson:"flag" json:"flag"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type FittingItemFlag string

const (
	FittingItemFlagCargo          FittingItemFlag = "Cargo"
	FittingItemFlagDroneBay       FittingItemFlag = "DroneBay"
	FittingItemFlagFighterBay     FittingItemFlag = "FighterBay"
	FittingItemFlagHiSlot0        FittingItemFlag = "HiSlot0"
	FittingItemFlagHiSlot1        FittingItemFlag = "HiSlot1"
	FittingItemFlagHiSlot2        FittingItemFlag = "HiSlot2"
	FittingItemFlagHiSlot3        FittingItemFlag = "HiSlot3"
	FittingItemFlagHiSlot4        FittingItemFlag = "HiSlot4"
	FittingItemFlagHiSlot5        FittingItemFlag = "HiSlot5"
	FittingItemFlagHiSlot6        FittingItemFlag = "HiSlot6"
	FittingItemFlagHiSlot7        FittingItemFlag = "HiSlot7"
	FittingItemFlagInvalid        FittingItemFlag = "Invalid"
	FittingItemFlagLoSlot0        FittingItemFlag = "LoSlot0"
	FittingItemFlagLoSlot1        FittingItemFlag = "LoSlot1"
	FittingItemFlagLoSlot2        FittingItemFlag = "LoSlot2"
	FittingItemFlagLoSlot3        FittingItemFlag = "LoSlot3"
	FittingItemFlagLoSlot4        FittingItemFlag = "LoSlot4"
	FittingItemFlagLoSlot5        FittingItemFlag = "LoSlot5"
	FittingItemFlagLoSlot6        FittingItemFlag = "LoSlot6"
	FittingItemFlagLoSlot7        FittingItemFlag = "LoSlot7"
	FittingItemFlagMedSlot0       FittingItemFlag = "MedSlot0"
	FittingItemFlagMedSlot1       FittingItemFlag = "MedSlot1"
	FittingItemFlagMedSlot2       FittingItemFlag = "MedSlot2"
	FittingItemFlagMedSlot3       FittingItemFlag = "MedSlot3"
	FittingItemFlagMedSlot4       FittingItemFlag = "MedSlot4"
	FittingItemFlagMedSlot5       FittingItemFlag = "MedSlot5"
	FittingItemFlagMedSlot6       FittingItemFlag = "MedSlot6"
	FittingItemFlagMedSlot7       FittingItemFlag = "MedSlot7"
	FittingItemFlagRigSlot0       FittingItemFlag = "RigSlot0"
	FittingItemFlagRigSlot1       FittingItemFlag = "RigSlot1"
	FittingItemFlagRigSlot2       FittingItemFlag = "RigSlot2"
	FittingItemFlagServiceSlot0   FittingItemFlag = "ServiceSlot0"
	FittingItemFlagServiceSlot1   FittingItemFlag = "ServiceSlot1"
	FittingItemFlagServiceSlot2   FittingItemFlag = "ServiceSlot2"
	FittingItemFlagServiceSlot3   FittingItemFlag = "ServiceSlot3"
	FittingItemFlagServiceSlot4   FittingItemFlag = "ServiceSlot4"
	FittingItemFlagServiceSlot5   FittingItemFlag = "ServiceSlot5"
	FittingItemFlagServiceSlot6   FittingItemFlag = "ServiceSlot6"
	FittingItemFlagServiceSlot7   FittingItemFlag = "ServiceSlot7"
	FittingItemFlagSubSystemSlot0 FittingItemFlag = "SubSystemSlot0"
	FittingItemFlagSubSystemSlot1 FittingItemFlag = "SubSystemSlot1"
	FittingItemFlagSubSystemSlot2 FittingItemFlag = "SubSystemSlot2"
	FittingItemFlagSubSystemSlot3 FittingItemFlag = "SubSystemSlot3"
)

var AllFittingItemFlags = []FittingItemFlag{
	FittingItemFlagCargo,
	FittingItemFlagDroneBay,
	FittingItemFlagFighterBay,
	FittingItemFlagHiSlot0,
	FittingItemFlagHiSlot1,
	FittingItemFlagHiSlot2,
	FittingItemFlagHiSlot3,
	FittingItemFlagHiSlot4,
	FittingItemFlagHiSlot5,
	FittingItemFlagHiSlot6,
	FittingItemFlagHiSlot7,
	FittingItemFlagInvalid,
	FittingItemFlagLoSlot0,
	FittingItemFlagLoSlot1,
	FittingItemFlagLoSlot2,
	FittingItemFlagLoSlot3,
	FittingItemFlagLoSlot4,
	FittingItemFlagLoSlot5,
	FittingItemFlagLoSlot6,
	FittingItemFlagLoSlot7,
	FittingItemFlagMedSlot0,
	FittingItemFlagMedSlot1,
	FittingItemFlagMedSlot2,
	FittingItemFlagMedSlot3,
	FittingItemFlagMedSlot4,
	FittingItemFlagMedSlot5,
	FittingItemFlagMedSlot6,
	FittingItemFlagMedSlot7,
	FittingItemFlagRigSlot0,
	FittingItemFlagRigSlot1,
	FittingItemFlagRigSlot2,
	FittingItemFlagServiceSlot0,
	FittingItemFlagServiceSlot1,
	FittingItemFlagServiceSlot2,
	FittingItemFlagServiceSlot3,
	FittingItemFlagServiceSlot4,
	FittingItemFlagServiceSlot5,
	FittingItemFlagServiceSlot6,
	FittingItemFlagServiceSlot7,
	FittingItemFlagSubSystemSlot0,
	FittingItemFlagSubSystemSlot1,
	FittingItemFlagSubSystemSlot2,
	FittingItemFlagSubSystemSlot3,
}

func (i FittingItemFlag) Valid() bool {
	for _, v := range AllFittingItemFlags {
		if i == v {
			return true
		}
	}

	return false
}

func (i FittingItemFlag) String() string {
	return string(i)
}
