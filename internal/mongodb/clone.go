package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type memberCloneRepository struct {
	clones   *mongo.Collection
	implants *mongo.Collection
}

func NewCloneRepository(d *mongo.Database) (athena.CloneRepository, error) {

	clones := d.Collection("member_clones")
	clonesIdxModel := mongo.IndexModel{
		Keys: primitive.D{
			primitive.E{
				Key:   "member_id",
				Value: 1,
			},
		},
		Options: &options.IndexOptions{
			Name:   newString("member_clones_member_id_unique"),
			Unique: newBool(true),
		},
	}
	_, err := clones.Indexes().CreateOne(context.Background(), clonesIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository]: Failed to create index on clones collection: %w", err)
	}

	implants := d.Collection("member_implants")
	implantsIdxModel := mongo.IndexModel{
		Keys: primitive.D{
			primitive.E{
				Key:   "member_id",
				Value: 1,
			},
		},
		Options: &options.IndexOptions{
			Name:   newString("member_implants_member_id_unique"),
			Unique: newBool(true),
		},
	}
	_, err = implants.Indexes().CreateOne(context.Background(), implantsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository]: Failed to create index on implants collection: %w", err)
	}

	return &memberCloneRepository{
		clones:   clones,
		implants: implants,
	}, nil
}

func (r *memberCloneRepository) MemberClones(ctx context.Context, id string) (*athena.MemberClones, error) {

	clones := new(athena.MemberClones)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.clones.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(clones)

	return clones, err

}

func (r *memberCloneRepository) CreateMemberClones(ctx context.Context, clones *athena.MemberClones) (*athena.MemberClones, error) {

	clones.CreatedAt = time.Now()
	clones.UpdatedAt = time.Now()

	_, err := r.clones.InsertOne(ctx, clones)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to insert record into clones collection: %w", err)
	}

	return clones, nil

}

func (r *memberCloneRepository) UpdateMemberClones(ctx context.Context, id string, clones *athena.MemberClones) (*athena.MemberClones, error) {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	clones.MemberID = pid
	clones.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: clones}}

	_, err = r.clones.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Clone Repository] Failed to insert record into clones collection: %w", err)
	}

	return clones, err

}

func (r *memberCloneRepository) DeleteMemberClones(ctx context.Context, id string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	results, err := r.clones.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Clone Repository] Failed to delete record from clones collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberCloneRepository) MemberImplants(ctx context.Context, id string) (*athena.MemberImplants, error) {

	implants := new(athena.MemberImplants)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Implant Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.implants.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(implants)

	return implants, err

}

func (r *memberCloneRepository) CreateMemberImplants(ctx context.Context, implants *athena.MemberImplants) (*athena.MemberImplants, error) {

	implants.CreatedAt = time.Now()
	implants.UpdatedAt = time.Now()

	_, err := r.implants.InsertOne(ctx, implants)
	if err != nil {
		return nil, fmt.Errorf("[Implant Repository] Failed to insert record into implants collection: %w", err)
	}

	return implants, nil

}

func (r *memberCloneRepository) UpdateMemberImplants(ctx context.Context, id string, implants *athena.MemberImplants) (*athena.MemberImplants, error) {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Implant Repository] Failed to cast id to objectID: %w", err)
	}

	implants.MemberID = pid
	implants.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: implants}}

	_, err = r.implants.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Implant Repository] Failed to insert record into implants collection: %w", err)
	}

	return implants, err

}

func (r *memberCloneRepository) DeleteMemberImplants(ctx context.Context, id string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("[Implant Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	results, err := r.implants.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Implant Repository] Failed to delete record from implants collection: %w", err)
	}

	return results.DeletedCount > 0, err

}
