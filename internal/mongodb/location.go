package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type memberLocationRepository struct {
	location *mongo.Collection
	online   *mongo.Collection
	ship     *mongo.Collection
}

func NewLocationRepository(d *mongo.Database) (athena.MemberLocationRepository, error) {

	location := d.Collection("member_location")
	locIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"member_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("location_member_id_unique"),
			Unique: newBool(true),
		},
	}
	_, err := location.Indexes().CreateOne(context.Background(), locIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository]: Failed to Create Index on Location Collection: %w", err)
	}

	online := d.Collection("member_online")
	onlineIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"member_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("online_member_id_unique"),
			Unique: newBool(true),
		},
	}
	_, err = online.Indexes().CreateOne(context.Background(), onlineIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository]: Failed to Create Index on Online Collection: %w", err)
	}

	ship := d.Collection("member_ship")
	shipIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"member_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("ship_member_id_unique"),
			Unique: newBool(true),
		},
	}
	_, err = ship.Indexes().CreateOne(context.Background(), shipIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository]: Failed to Create Index on Ship Collection: %w", err)
	}

	return &memberLocationRepository{
		location: location,
		online:   online,
		ship:     ship,
	}, nil

}

func (r *memberLocationRepository) MemberLocation(ctx context.Context, id string) (*athena.MemberLocation, error) {

	location := new(athena.MemberLocation)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.location.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(location)

	return location, err

}

func (r *memberLocationRepository) CreateMemberLocation(ctx context.Context, location *athena.MemberLocation) (*athena.MemberLocation, error) {

	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()

	result, err := r.location.InsertOne(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to insert record into location collection: %w", err)
	}

	location.ID = result.InsertedID.(primitive.ObjectID)

	return location, nil

}

func (r *memberLocationRepository) UpdateMemberLocation(ctx context.Context, id string, location *athena.MemberLocation) (*athena.MemberLocation, error) {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	location.ID = pid
	location.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: location}}

	_, err = r.location.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Location Repository] Failed to insert record into location collection: %w", err)
	}

	return location, err

}

func (r *memberLocationRepository) DeleteMemberLocation(ctx context.Context, id string) error {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "_id", Value: pid}}

	_, err = r.location.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Location Repository] Failed to delete record from ship collection: %w", err)
	}

	return err

}

func (r *memberLocationRepository) MemberOnline(ctx context.Context, id string) (*athena.MemberOnline, error) {

	online := new(athena.MemberOnline)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.online.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(online)

	return online, err

}

func (r *memberLocationRepository) CreateMemberOnline(ctx context.Context, online *athena.MemberOnline) (*athena.MemberOnline, error) {

	online.CreatedAt = time.Now()
	online.UpdatedAt = time.Now()

	result, err := r.online.InsertOne(ctx, online)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to insert record into online collection: %w", err)
	}

	online.ID = result.InsertedID.(primitive.ObjectID)

	return online, nil

}

func (r *memberLocationRepository) UpdateMemberOnline(ctx context.Context, id string, online *athena.MemberOnline) (*athena.MemberOnline, error) {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	online.ID = pid
	online.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: online}}

	_, err = r.online.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Location Repository] Failed to insert record into online collection: %w", err)
	}

	return online, err

}

func (r *memberLocationRepository) DeleteMemberOnline(ctx context.Context, id string) error {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "_id", Value: pid}}

	_, err = r.online.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Location Repository] Failed to delete record from ship collection: %w", err)
	}

	return err

}

func (r *memberLocationRepository) MemberShip(ctx context.Context, id string) (*athena.MemberShip, error) {

	ship := new(athena.MemberShip)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.ship.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(ship)

	return ship, err

}

func (r *memberLocationRepository) CreateMemberShip(ctx context.Context, ship *athena.MemberShip) (*athena.MemberShip, error) {

	ship.CreatedAt = time.Now()
	ship.UpdatedAt = time.Now()

	result, err := r.ship.InsertOne(ctx, ship)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to insert record into ship collection: %w", err)
	}

	ship.ID = result.InsertedID.(primitive.ObjectID)

	return ship, nil

}

func (r *memberLocationRepository) UpdateMemberShip(ctx context.Context, id string, ship *athena.MemberShip) (*athena.MemberShip, error) {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	ship.ID = pid
	ship.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: ship}}

	_, err = r.ship.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Location Repository] Failed to insert record into ship collection: %w", err)
	}

	return ship, err

}

func (r *memberLocationRepository) DeleteMemberShip(ctx context.Context, id string) error {

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("[Location Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "_id", Value: pid}}

	_, err = r.ship.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Location Repository] Failed to delete record from ship collection: %w", err)
	}

	return err

}
