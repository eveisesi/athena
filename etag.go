package athena

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EtagRepository interface {
	MemberEtag(ctx context.Context, endpointID Scope) (*MemberEtag, error)
	MemberEtags(ctx context.Context, memberID string) ([]*MemberEtag, error)
	CreateMemberEtag(ctx context.Context, location *MemberEtag) (*MemberEtag, error)
	UpdateMemberEtag(ctx context.Context, id string, location *MemberEtag) (*MemberEtag, error)
	DeleteMemberEtag(ctx context.Context, id string) (bool, error)
}

type MemberEtag struct {
	MemberID   primitive.ObjectID `bson:"member_id" json:"member_id"`
	EndpointID Scope              `bson:"endpoint_id" json:"endpoint_id"`
	Meta
}
