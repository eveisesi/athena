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

type allianceRepository struct {
	alliances *mongo.Collection
}

func NewAllianceRepository(d *mongo.Database) (athena.AllianceRepository, error) {
	alliances := d.Collection("alliances")
	allianceIndexModel := mongo.IndexModel{
		Keys: primitive.D{
			primitive.E{
				Key:   "alliance_id",
				Value: 1,
			},
		},
		Options: &options.IndexOptions{
			Name:   newString("alliances_alliance_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err := alliances.Indexes().CreateOne(context.Background(), allianceIndexModel)
	if err != nil {
		return nil, fmt.Errorf("[Alliance Repository]: Failed to Create Index on Alliances Collection: %w", err)
	}

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

	alliance.CreatedAt = time.Now()
	alliance.UpdatedAt = time.Now()

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
	alliance.UpdatedAt = time.Now()

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
