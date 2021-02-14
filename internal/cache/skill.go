package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type skillService interface {
	MemberAttributes(ctx context.Context, id uint) (*athena.MemberAttributes, error)
	SetMemberAttributes(ctx context.Context, id uint, attributes *athena.MemberAttributes, options ...OptionFunc) error
	MemberSkillQueue(ctx context.Context, id uint) ([]*athena.MemberSkillQueue, error)
	SetMemberSkillQueue(ctx context.Context, id uint, skillQueue []*athena.MemberSkillQueue, options ...OptionFunc) error
	MemberSkills(ctx context.Context, id uint) ([]*athena.Skill, error)
	SetMemberSkills(ctx context.Context, id uint, skills []*athena.Skill, options ...OptionFunc) error
	MemberSkillProperties(ctx context.Context, id uint) (*athena.MemberSkills, error)
	SetMemberSkillProperties(ctx context.Context, id uint, meta *athena.MemberSkills, optionFuncs ...OptionFunc) error
}

const (
	keyMemberAttributes      = "athena::member::%d::skill::attributes"
	keyMemberSkillQueue      = "athena::member::%d::skillqueue"
	keyMemberSkills          = "athena::member::%d::skills"
	keyMemberSkillProperties = "athena::member::%d::skill::properties"
)

func (s *service) MemberAttributes(ctx context.Context, id uint) (*athena.MemberAttributes, error) {

	key := fmt.Sprintf(keyMemberAttributes, id)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch data from cache for key %s: %w", key, err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var attributes = new(athena.MemberAttributes)
	err = json.Unmarshal(data, attributes)
	if err != nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal data for key %s on struct: %w", key, err)
	}

	return attributes, nil

}

func (s *service) SetMemberAttributes(ctx context.Context, id uint, attributes *athena.MemberAttributes, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	key := fmt.Sprintf(keyMemberAttributes, id)
	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberSkillQueue(ctx context.Context, id uint) ([]*athena.MemberSkillQueue, error) {

	key := fmt.Sprintf(keyMemberSkillQueue, id)
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	positions := make([]*athena.MemberSkillQueue, len(members))
	for i, member := range members {
		var position = new(athena.MemberSkillQueue)
		err = json.Unmarshal([]byte(member), position)
		if err != nil {
			return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
		}

		positions[i] = position
	}

	return positions, nil

}

func (s *service) SetMemberSkillQueue(ctx context.Context, id uint, positions []*athena.MemberSkillQueue, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	members := make([]string, 0, len(positions))
	for _, position := range positions {
		data, err := json.Marshal(position)
		if err != nil {
			return fmt.Errorf("failed to marshal skill queue position: %w", err)
		}

		members = append(members, string(data))
	}

	key := fmt.Sprintf(keyMemberSkillQueue, id)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to write to cache: %w", err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) MemberSkills(ctx context.Context, id uint) ([]*athena.Skill, error) {

	key := fmt.Sprintf(keyMemberSkills, id)
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("[Cache Layer] Failed to fetch results from cache for key %s: %w", key, err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	skills := make([]*athena.Skill, len(members))
	for i, member := range members {
		var skill = new(athena.Skill)
		err = json.Unmarshal([]byte(member), skill)
		if err != nil {
			return nil, fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
		}

		skills[i] = skill
	}

	return skills, nil

}

func (s *service) SetMemberSkills(ctx context.Context, id uint, skills []*athena.Skill, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	members := make([]string, 0, len(skills))
	for _, skill := range skills {
		data, err := json.Marshal(skill)
		if err != nil {
			return fmt.Errorf("failed to marshal skill queue position: %w", err)
		}

		members = append(members, string(data))
	}

	key := fmt.Sprintf(keyMemberSkills, id)
	_, err := s.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache skills for member %d: %w", id, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) MemberSkillProperties(ctx context.Context, id uint) (*athena.MemberSkills, error) {

	key := fmt.Sprintf(keyMemberSkillProperties, id)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var meta = new(athena.MemberSkills)

	err = json.Unmarshal(result, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil

}

func (s *service) SetMemberSkillProperties(ctx context.Context, id uint, meta *athena.MemberSkills, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	key := fmt.Sprintf(keyMemberSkillProperties, id)
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, key, data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
