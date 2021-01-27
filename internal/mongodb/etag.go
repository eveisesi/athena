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

type etagRepository struct {
	etag *mongo.Collection
}

func NewEtagRepository(d *mongo.Database) (athena.EtagRepository, error) {

	var ctx = context.Background()

	etag := d.Collection("etags")
	etagIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"etag_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("etags_etag_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err := etag.Indexes().CreateOne(ctx, etagIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository]: Failed to Create Index on Corporations Collection: %w", err)
	}

	return &etagRepository{
		etag: etag,
	}, nil

}

func (r *etagRepository) Etag(ctx context.Context, etagID string) (*athena.Etag, error) {

	etag := new(athena.Etag)

	filter := primitive.D{primitive.E{Key: "etag_id", Value: etagID}}

	err := r.etag.FindOne(ctx, filter).Decode(etag)

	return etag, err

}
func (r *etagRepository) Etags(ctx context.Context, operators ...*athena.Operator) ([]*athena.Etag, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	results, err := r.etag.Find(ctx, filters, options)
	if err != nil {
		return nil, err
	}

	var etags = make([]*athena.Etag, 0)

	return etags, results.All(ctx, &etags)

}

func (r *etagRepository) UpdateEtag(ctx context.Context, etagID string, etag *athena.Etag) (*athena.Etag, error) {

	etag.EtagID = etagID
	if etag.CreatedAt.IsZero() {
		etag.CreatedAt = time.Now()

	}
	if etag.UpdatedAt.IsZero() {
		etag.UpdatedAt = time.Now()
	}

	filter := primitive.D{primitive.E{Key: "etag_id", Value: etagID}}
	update := primitive.D{primitive.E{Key: "$set", Value: etag}}
	options := &options.UpdateOptions{
		Upsert: newBool(true),
	}

	_, err := r.etag.UpdateOne(ctx, filter, update, options)
	if err != nil {
		err = fmt.Errorf("[Etag Repository] Failed to insert record into etag collection: %w", err)
	}

	return etag, err

}

func (r *etagRepository) DeleteEtag(ctx context.Context, etagID string) (bool, error) {

	filter := primitive.D{primitive.E{Key: "etag_id", Value: etagID}}

	result, err := r.etag.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Etag Repository] Failed to delete records from the etag collection: %w", err)
	}

	return result.DeletedCount > 0, err

}
