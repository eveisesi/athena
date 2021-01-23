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

type memberSkillRepository struct {
	skills     *mongo.Collection
	skillQueue *mongo.Collection
	attributes *mongo.Collection
}

func NewMemberSkillRepository(d *mongo.Database) (athena.MemberSkillRepository, error) {

	var ctx = context.Background()

	skills := d.Collection("member_skills")
	skillIdxMods := []mongo.IndexModel{
		{
			Keys: bson.M{
				"member_id": 1,
			},
			Options: &options.IndexOptions{
				Name:   newString("member_skills_member_id_unique"),
				Unique: newBool(true),
			},
		},
		{
			Keys: bson.M{
				"skills.skill_id": 1,
			},
			Options: &options.IndexOptions{
				Name: newString("member_skills_skills_skill_id_idx"),
			},
		},
	}
	_, err := skills.Indexes().CreateMany(ctx, skillIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository]: Failed to create index on skills collection: %w", err)
	}

	skillQueue := d.Collection("member_skill_queue")
	skillQueueIdxMods := []mongo.IndexModel{
		{
			Keys: bson.M{
				"member_id": 1,
			},
			Options: &options.IndexOptions{
				Name:   newString("member_skillqueue_member_id_unique"),
				Unique: newBool(true),
			},
		},
		{
			Keys: bson.M{
				"skill_queue.skill_id":       1,
				"skill_queue.queue_position": 1,
			},
			Options: &options.IndexOptions{
				Name: newString("member_skillqueue_skill_queue_idx"),
			},
		},
	}
	_, err = skillQueue.Indexes().CreateMany(ctx, skillQueueIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository]: Failed to create index on skillQueue collection: %w", err)
	}

	attributes := d.Collection("member_attributes")
	attributeIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"member_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("member_attributes_member_id_unique"),
			Unique: newBool(true),
		},
	}
	_, err = attributes.Indexes().CreateOne(ctx, attributeIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository]: Failed to create index on attributes collection: %w", err)
	}

	return &memberSkillRepository{
		skills:     skills,
		skillQueue: skillQueue,
		attributes: attributes,
	}, nil

}

func (r *memberSkillRepository) MemberSkillAttributes(ctx context.Context, memberID string) (*athena.MemberSkillAttributes, error) {

	attributes := new(athena.MemberSkillAttributes)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.attributes.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(attributes)

	return attributes, err

}

func (r *memberSkillRepository) CreateMemberSkillAttributes(ctx context.Context, attributes *athena.MemberSkillAttributes) (*athena.MemberSkillAttributes, error) {

	attributes.CreatedAt = time.Now()
	attributes.UpdatedAt = time.Now()

	_, err := r.attributes.InsertOne(ctx, attributes)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert record into attributes collection: %w", err)
	}

	return attributes, nil

}

func (r *memberSkillRepository) UpdateMemberSkillAttributes(ctx context.Context, memberID string, attributes *athena.MemberSkillAttributes) (*athena.MemberSkillAttributes, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	attributes.MemberID = pid
	attributes.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: attributes}}

	_, err = r.attributes.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to insert record into attributes collection: %w", err)
	}

	return attributes, err

}

func (r *memberSkillRepository) DeleteMemberSkillAttributes(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	results, err := r.attributes.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to delete record from attributes collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberSkillRepository) MemberSkillQueue(ctx context.Context, memberID string) (*athena.MemberSkillQueue, error) {

	skillQueue := new(athena.MemberSkillQueue)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.skillQueue.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(skillQueue)

	return skillQueue, err

}

func (r *memberSkillRepository) CreateMemberSkillQueue(ctx context.Context, skillQueue *athena.MemberSkillQueue) (*athena.MemberSkillQueue, error) {

	skillQueue.CreatedAt = time.Now()
	skillQueue.UpdatedAt = time.Now()

	_, err := r.skillQueue.InsertOne(ctx, skillQueue)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert record into skillQueue collection: %w", err)
	}

	return skillQueue, nil

}

func (r *memberSkillRepository) UpdateMemberSkillQueue(ctx context.Context, memberID string, skillQueue *athena.MemberSkillQueue) (*athena.MemberSkillQueue, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	skillQueue.MemberID = pid
	skillQueue.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: skillQueue}}

	_, err = r.skillQueue.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to insert record into skillQueue collection: %w", err)
	}

	return skillQueue, err

}

func (r *memberSkillRepository) DeleteMemberSkillQueue(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	results, err := r.skillQueue.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to delete record from skillQueue collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberSkillRepository) MemberSkills(ctx context.Context, memberID string) (*athena.MemberSkill, error) {

	skills := new(athena.MemberSkill)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.skills.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(skills)

	return skills, err

}

func (r *memberSkillRepository) CreateMemberSkills(ctx context.Context, skills *athena.MemberSkill) (*athena.MemberSkill, error) {

	skills.CreatedAt = time.Now()
	skills.UpdatedAt = time.Now()

	_, err := r.skills.InsertOne(ctx, skills)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert record into skills collection: %w", err)
	}

	return skills, nil

}

func (r *memberSkillRepository) UpdateMemberSkills(ctx context.Context, memberID string, skills *athena.MemberSkill) (*athena.MemberSkill, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	skills.MemberID = pid
	skills.UpdatedAt = time.Now()

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: skills}}

	_, err = r.skills.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to insert record into skills collection: %w", err)
	}

	return skills, err

}

func (r *memberSkillRepository) DeleteMemberSkills(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	results, err := r.skills.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to delete record from skills collection: %w", err)
	}

	return results.DeletedCount > 0, err

}
