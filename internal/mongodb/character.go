package mongodb

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type characterRepository struct {
	characters *mongo.Collection
}

func NewCharacterRepository(d *mongo.Database) (athena.CharacterRepository, error) {
	characters := d.Collection("characters")

	return &characterRepository{
		characters,
	}, nil
}

func (r *characterRepository) Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var characters = make([]*athena.Character, 0)
	result, err := r.characters.Find(ctx, filters, options)
	if err != nil {
		return characters, err
	}

	err = result.All(ctx, &characters)

	return characters, err

}

func (r *characterRepository) CreateCharacter(ctx context.Context, character *athena.Character) (*athena.Character, error) {
	result, err := r.characters.InsertOne(ctx, character)
	if err != nil {
		return nil, err
	}

	character.ID = result.InsertedID.(primitive.ObjectID)

	return character, err
}

func (r *characterRepository) UpdateCharacter(ctx context.Context, id string, character *athena.Character) (*athena.Character, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	character.ID = _id

	update := primitive.D{primitive.E{Key: "$set", Value: character}}

	_, err = r.characters.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return character, err
}

func (r *characterRepository) DeleteCharacter(ctx context.Context, id string) (bool, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("unable to cast %s to ObjectID: %w", id, err)
	}

	filters := BuildFilters(athena.NewEqualOperator("_id", _id))

	result, err := r.characters.DeleteOne(ctx, filters)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, err

}
