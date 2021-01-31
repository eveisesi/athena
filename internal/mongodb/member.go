package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type memberRepository struct {
	members *mongo.Collection
}

func NewMemberRepository(d *mongo.Database) (athena.MemberRepository, error) {
	members := d.Collection("members")

	return &memberRepository{
		members,
	}, nil
}

func (r *memberRepository) Member(ctx context.Context, id string) (*athena.Member, error) {

	var member = new(athena.Member)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Member Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.members.FindOne(ctx, primitive.D{primitive.E{Key: "_id", Value: pid}}).Decode(member)

	return member, err

}

func (r *memberRepository) Members(ctx context.Context, operators ...*athena.Operator) ([]*athena.Member, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var members = make([]*athena.Member, 0)
	result, err := r.members.Find(ctx, filters, options)
	if err != nil {
		return members, err
	}

	err = result.All(ctx, &members)

	return members, err

}

func (r *memberRepository) CreateMember(ctx context.Context, member *athena.Member) (*athena.Member, error) {

	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()

	result, err := r.members.InsertOne(ctx, member)
	if err != nil {
		return nil, err
	}

	member.ID = result.InsertedID.(primitive.ObjectID)

	return member, err
}

func (r *memberRepository) UpdateMember(ctx context.Context, id string, member *athena.Member) (*athena.Member, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	member.ID = _id
	member.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "_id", Value: _id}}
	update := primitive.D{primitive.E{Key: "$set", Value: member}}

	_, err = r.members.UpdateOne(ctx, filter, update)

	return member, err
}

func (r *memberRepository) DeleteMember(ctx context.Context, id string) (bool, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	filters := BuildFilters(athena.NewEqualOperator("_id", _id))

	_, err = r.members.DeleteOne(ctx, filters)
	if err != nil {
		return false, err
	}

	return err == nil, err

}
