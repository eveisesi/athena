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
	skillMeta  *mongo.Collection
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
				"skill_id":  1,
			},
			Options: &options.IndexOptions{
				Name:   newString("member_skills_member_id_skill_id_unique"),
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

	skillMeta := d.Collection("member_skill_meta")
	skillMetaIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"member_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("member_skill_meta_member_id_unique"),
			Unique: newBool(true),
		},
	}
	_, err = skillMeta.Indexes().CreateOne(ctx, skillMetaIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository]: Failed to create index on skill meta collection: %w", err)
	}

	skillQueue := d.Collection("member_skillqueue")
	skillQueueIdxMods := []mongo.IndexModel{
		{
			Keys: bson.M{
				"member_id":      1,
				"queue_position": 1,
			},
			Options: &options.IndexOptions{
				Name: newString("member_skillqueue_member_id_queue_position_idx"),
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
		skillMeta:  skillMeta,
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

func (r *memberSkillRepository) MemberSkillQueue(ctx context.Context, memberID string) ([]*athena.MemberSkillQueue, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	result, err := r.skillQueue.Find(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}})
	if err != nil {
		return nil, fmt.Errorf("Failed : %w", err)
	}

	var skillQueue = make([]*athena.MemberSkillQueue, 0)

	return skillQueue, result.All(ctx, &skillQueue)

}

func (r *memberSkillRepository) CreateMemberSkillQueue(ctx context.Context, memberID string, skillQueue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	now := time.Now()
	documents := make([]interface{}, len(skillQueue))
	for i, entry := range skillQueue {
		entry.MemberID = pid
		entry.CreatedAt = now
		entry.UpdatedAt = now
		documents[i] = entry
	}

	_, err = r.skillQueue.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert record into member_skillqueue collection: %w", err)
	}

	return skillQueue, nil

}

func (r *memberSkillRepository) UpdateMemberSkillQueue(ctx context.Context, memberID string, position *athena.MemberSkillQueue) (*athena.MemberSkillQueue, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	now := time.Now()
	position.MemberID = pid
	position.UpdatedAt = now
	if position.CreatedAt.IsZero() {
		position.CreatedAt = now
	}

	filter := primitive.D{
		primitive.E{Key: "member_id", Value: pid},
		primitive.E{Key: "queue_position", Value: position.QueuePosition},
	}
	update := primitive.D{primitive.E{Key: "$set", Value: position}}

	_, err = r.skillQueue.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to update record in the member_skillqueue collection: %w", err)
	}

	return position, err

}

func (r *memberSkillRepository) DeleteMemberSkillQueue(ctx context.Context, memberID string, entries []*athena.MemberSkillQueue) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	if len(entries) > 0 {
		queuePositions := make([]int, len(entries))
		for i, entry := range entries {
			queuePositions[i] = entry.QueuePosition
		}

		if len(queuePositions) > 0 {
			filter = append(filter, primitive.E{
				Key: "queue_position",
				Value: primitive.D{
					primitive.E{
						Key:   "$in",
						Value: queuePositions,
					},
				},
			})
		}
	}

	results, err := r.skillQueue.DeleteMany(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to delete record from skillQueue collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberSkillRepository) MemberSkills(ctx context.Context, memberID string) ([]*athena.MemberSkill, error) {

	skills := make([]*athena.MemberSkill, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id %s to objectID: %w", memberID, err)
	}

	results, err := r.skills.Find(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}})
	if err != nil {
		return nil, fmt.Errorf("[Skills Repository] Failed to fetch skills for member %s: %w", memberID, err)
	}

	err = results.All(ctx, &skills)

	return skills, err

}

func (r *memberSkillRepository) CreateMemberSkills(ctx context.Context, memberID string, skills []*athena.MemberSkill) ([]*athena.MemberSkill, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	now := time.Now()
	insert := make([]interface{}, len(skills))
	for i, skill := range skills {
		skill.MemberID = pid
		skill.CreatedAt = now
		skill.UpdatedAt = now
		insert[i] = skill
	}

	_, err = r.skills.InsertMany(ctx, insert)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert record into skills collection: %w", err)
	}

	return skills, nil

}

func (r *memberSkillRepository) UpdateMemberSkills(ctx context.Context, memberID string, skill *athena.MemberSkill) (*athena.MemberSkill, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	skill.MemberID = pid
	skill.UpdatedAt = time.Now()
	if skill.CreatedAt.IsZero() {
		skill.CreatedAt = time.Now()
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "skill_id", Value: skill.SkillID}}
	update := primitive.D{primitive.E{Key: "$set", Value: skill}}

	_, err = r.skills.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to update records into skills collection: %w", err)
	}

	return skill, err

}

func (r *memberSkillRepository) DeleteMemberSkills(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	results, err := r.skills.DeleteMany(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to delete record from skills collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberSkillRepository) MemberSkillMeta(ctx context.Context, memberID string) (*athena.MemberSkillMeta, error) {

	meta := new(athena.MemberSkillMeta)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.skillMeta.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}}).Decode(meta)

	return meta, err

}

func (r *memberSkillRepository) CreateMemberSkillMeta(ctx context.Context, meta *athena.MemberSkillMeta) (*athena.MemberSkillMeta, error) {

	now := time.Now()
	meta.CreatedAt = now
	meta.UpdatedAt = now

	_, err := r.skillMeta.InsertOne(ctx, meta)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert skill meta into the skill meta collection: %w", err)
	}

	return meta, nil

}

func (r *memberSkillRepository) UpdateMemberSkillMeta(ctx context.Context, memberID string, meta *athena.MemberSkillMeta) (*athena.MemberSkillMeta, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	meta.MemberID = pid
	if meta.CreatedAt.IsZero() {
		meta.CreatedAt = time.Now()

	}
	if meta.UpdatedAt.IsZero() {
		meta.UpdatedAt = time.Now()
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}
	update := primitive.D{primitive.E{Key: "$set", Value: meta}}
	options := &options.UpdateOptions{
		Upsert: newBool(true),
	}

	_, err = r.skillMeta.UpdateOne(ctx, filter, update, options)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to insert record into skillMeta collection: %w", err)
	}

	return meta, err

}

func (r *memberSkillRepository) DeleteMemberSkillMeta(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}}

	results, err := r.skillMeta.DeleteOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("[Skill Repository] Failed to delete record from attributes collection: %w", err)
	}

	return results.DeletedCount > 0, err

}
