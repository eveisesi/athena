package athena

import (
	"context"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CloneRepository interface {
	cloneRepository
	implantRepository
}

type cloneRepository interface {
	MemberClones(ctx context.Context, id string) (*MemberClones, error)
	CreateMemberClones(ctx context.Context, location *MemberClones) (*MemberClones, error)
	UpdateMemberClones(ctx context.Context, id string, location *MemberClones) (*MemberClones, error)
	DeleteMemberClones(ctx context.Context, id string) (bool, error)
}

type implantRepository interface {
	MemberImplants(ctx context.Context, id string) (*MemberImplants, error)
	CreateMemberImplants(ctx context.Context, location *MemberImplants) (*MemberImplants, error)
	UpdateMemberImplants(ctx context.Context, id string, location *MemberImplants) (*MemberImplants, error)
	DeleteMemberImplants(ctx context.Context, id string) (bool, error)
}

type MemberClones struct {
	MemberID            primitive.ObjectID `bson:"member_id" json:"member_id"`
	HomeLocation        CloneHomeLocation  `bson:"home_location" json:"home_location"`
	JumpClones          []CloneJumpClone   `bson:"jump_clones" json:"jump_clones"`
	LastCloneJumpDate   null.Time          `bson:"last_clone_jump_data" json:"last_clone_jump_data"`
	LastCloneChangeDate null.Time          `bson:"last_clone_change_date" json:"last_clone_change_date"`
	Meta
}

type CloneHomeLocation struct {
	LocationID   int64  `bson:"location_id" json:"location_id"`
	LocationType string `bson:"location_type" json:"location_type"`
}

type CloneJumpClone struct {
	Implants     []int  `bson:"implants" json:"implants"`
	JumpCloneID  int    `bson:"jump_clone_id" json:"jump_clone_id"`
	LocationID   int64  `bson:"location_id" json:"location_id"`
	LocationType string `bson:"location_type" json:"location_type"`
}

type MemberImplants struct {
	MemberID primitive.ObjectID `bson:"member_id" json:"member_id"`
	Implants []*Type            `bson:"implants" json:"implants"`
	Raw      []int              `bson:"-" json:"-"`
	Meta
}
