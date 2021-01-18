package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type corporationRepository struct {
	corporations *mongo.Collection
}

func NewCorporationRepository(d *mongo.Database) (athena.CorporationRepository, error) {
	corporations := d.Collection("corporations")
	corporationIndexModel := mongo.IndexModel{
		Keys: bson.M{
			"corporation_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("corporations_corporation_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err := corporations.Indexes().CreateOne(context.Background(), corporationIndexModel)
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository]: Failed to Create Index on Corporations Collection: %w", err)
	}

	return &corporationRepository{
		corporations,
	}, nil
}

func (r *corporationRepository) Corporations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Corporation, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var corporations = make([]*athena.Corporation, 0)
	result, err := r.corporations.Find(ctx, filters, options)
	if err != nil {
		return corporations, err
	}

	err = result.All(ctx, &corporations)

	return corporations, err

}

func (r *corporationRepository) CreateCorporation(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, error) {

	corporation.CreatedAt = time.Now()
	corporation.UpdatedAt = time.Now()

	result, err := r.corporations.InsertOne(ctx, corporation)
	if err != nil {
		return nil, err
	}

	corporation.ID = result.InsertedID.(primitive.ObjectID)

	return corporation, err
}

func (r *corporationRepository) UpdateCorporation(ctx context.Context, id string, corporation *athena.Corporation) (*athena.Corporation, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	corporation.ID = _id
	corporation.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: corporation}}

	_, err = r.corporations.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return corporation, err
}

func (r *corporationRepository) DeleteCorporation(ctx context.Context, id string) (bool, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	filters := BuildFilters(athena.NewEqualOperator("_id", _id))

	result, err := r.corporations.DeleteOne(ctx, filters)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, err

}
