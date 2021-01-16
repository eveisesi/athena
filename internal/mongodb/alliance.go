package mongodb

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type allianceRepository struct {
	alliances *mongo.Collection
}

func NewAllianceRepository(d *mongo.Database) (athena.AllianceRepository, error) {
	alliances := d.Collection("alliances")

	return &allianceRepository{
		alliances,
	}, nil
}

func (r *allianceRepository) Alliances(ctx context.Context, operators ...*athena.Operator) ([]*athena.Alliance, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var alliances = make([]*athena.Alliance, 0)
	result, err := r.alliances.Find(ctx, filters, options)
	if err != nil {
		return alliances, err
	}

	err = result.All(ctx, &alliances)

	return alliances, err

}

func (r *allianceRepository) CreateAlliance(ctx context.Context, alliance *athena.Alliance) (*athena.Alliance, error) {
	result, err := r.alliances.InsertOne(ctx, alliance)
	if err != nil {
		return nil, err
	}

	alliance.ID = result.InsertedID.(primitive.ObjectID)

	return alliance, err
}

func (r *allianceRepository) UpdateAlliance(ctx context.Context, id string, alliance *athena.Alliance) (*athena.Alliance, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	alliance.ID = _id

	update := primitive.D{primitive.E{Key: "$set", Value: alliance}}

	_, err = r.alliances.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return alliance, err
}

func (r *allianceRepository) DeleteAlliance(ctx context.Context, id string) (bool, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	filters := BuildFilters(athena.NewEqualOperator("_id", _id))

	result, err := r.alliances.DeleteOne(ctx, filters)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, err

}
