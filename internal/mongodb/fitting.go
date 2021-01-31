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

type memberFittingRepository struct {
	fittings *mongo.Collection
	items    *mongo.Collection
}

func NewMemberFittingRepository(ctx context.Context, d *mongo.Database) (athena.MemberFittingsRepository, error) {

	fittings := d.Collection("member_fittings")
	fittingsIdxMods := []mongo.IndexModel{
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "member_id",
					Value: 1,
				},
				primitive.E{
					Key:   "fitting_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name:   newString("member_fittings_member_id_fitting_id_unique"),
				Unique: newBool(true),
			},
		},
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "ship_type_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name: newString("member_fitting_ship_type_id_idx"),
			},
		},
	}

	_, err := fittings.Indexes().CreateMany(ctx, fittingsIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to create index on fitting repository: %w", err)
	}

	items := d.Collection("member_fitting_items")
	itemsIdxMods := []mongo.IndexModel{
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "member_id",
					Value: 1,
				},
				primitive.E{
					Key:   "fitting_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name: newString("member_fitting_item_type_id_idx"),
			},
		},
	}

	_, err = items.Indexes().CreateMany(ctx, itemsIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to create index on fitting repository: %w", err)
	}

	return &memberFittingRepository{
		fittings: fittings,
		items:    items,
	}, nil

}

func (r *memberFittingRepository) MemberFittings(ctx context.Context, memberID string, operators ...*athena.Operator) ([]*athena.MemberFitting, error) {

	fittings := make([]*athena.MemberFitting, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	operators = append(operators, athena.NewEqualOperator("member_id", pid))

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	cursor, err := r.fittings.Find(ctx, filters, options)
	if err != nil {
		return nil, err
	}

	return fittings, cursor.All(ctx, &fittings)

}

func (r *memberFittingRepository) CreateMemberFittings(ctx context.Context, memberID string, fittings []*athena.MemberFitting) ([]*athena.MemberFitting, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	now := time.Now()
	documents := make([]interface{}, len(fittings))
	for i, fitting := range fittings {
		fitting.MemberID = pid
		fitting.CreatedAt = now
		fitting.UpdatedAt = now

		documents[i] = fitting
	}

	_, err = r.fittings.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Member Fitting Repository] Failed to insert record into the member_fittings collection: %w", err)
	}

	return fittings, nil

}

func (r *memberFittingRepository) UpdateMemberFitting(ctx context.Context, memberID string, fittingID uint, fitting *athena.MemberFitting) (*athena.MemberFitting, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "fitting_id", Value: fitting.FittingID}}
	update := primitive.D{primitive.E{Key: "$set", Value: fitting}}

	_, err = r.fittings.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("[Member Fitting Repository] Failed to update record in the member_fittings collection: %w", err)
	}

	return fitting, nil

}

func (r *memberFittingRepository) DeleteMemberFitting(ctx context.Context, memberID string, fittingID uint) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "fitting_id", Value: fittingID}}

	_, err = r.fittings.DeleteOne(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to delete a record from the member_fittings collection: %w", err)
	}

	return err == nil, err

}

func (r *memberFittingRepository) DeleteMemberFittings(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	_, err = r.fittings.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to delete records from the member_fittings collection: %w", err)
	}

	return err == nil, err

}

func (r *memberFittingRepository) MemberFittingItems(ctx context.Context, memberID string, fittingID uint) ([]*athena.MemberFittingItem, error) {

	items := make([]*athena.MemberFittingItem, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return items, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "fitting_id", Value: fittingID}}

	cursor, err := r.items.Find(ctx, filter)
	if err != nil {
		return items, err
	}

	return items, cursor.All(ctx, items)

}

func (r *memberFittingRepository) CreateMemberFittingItems(ctx context.Context, memberID string, fittingID uint, items []*athena.MemberFittingItem) ([]*athena.MemberFittingItem, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	now := time.Now()
	documents := make([]interface{}, len(items))
	for i, fitting := range items {
		fitting.MemberID = pid
		fitting.FittingID = fittingID
		fitting.CreatedAt = now
		fitting.UpdatedAt = now

		documents[i] = fitting
	}

	_, err = r.items.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Member Fitting Repository] Failed to insert record into the member_firtting_items collection: %w", err)
	}

	return items, nil

}

func (r *memberFittingRepository) DeleteMemberFittingItems(ctx context.Context, memberID string, fittingID uint) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "fitting_id", Value: fittingID}}

	_, err = r.items.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to delete a record from the member_fittings collection: %w", err)
	}

	return err == nil, err

}

func (r *memberFittingRepository) DeleteMemberFittingItemsAll(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	_, err = r.items.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Member Fitting Repository] Failed to delete a record from the member_fittings collection: %w", err)
	}

	return err == nil, err

}
