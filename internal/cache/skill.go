package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type skillService interface {
	MemberSkillAttributes(ctx context.Context, memberID string) (*athena.MemberSkillAttributes, error)
	SetMemberSkillAttributes(ctx context.Context, memberID string, attributes *athena.MemberSkillAttributes, options ...OptionFunc) error
	MemberSkillQueue(ctx context.Context, memberID string) ([]*athena.MemberSkillQueue, error)
	SetMemberSkillQueue(ctx context.Context, memberID string, skillQueue []*athena.MemberSkillQueue, options ...OptionFunc) error
	MemberSkills(ctx context.Context, memberID string) ([]*athena.MemberSkill, error)
	SetMemberSkills(ctx context.Context, memberID string, skills []*athena.MemberSkill, options ...OptionFunc) error
	MemberSkillMeta(ctx context.Context, memberID string) (*athena.MemberSkillMeta, error)
	SetMemberSkillMeta(ctx context.Context, memberID string, meta *athena.MemberSkillMeta, optionFuncs ...OptionFunc) error
}

const (
	keyMemberSkillAttributes = "athena::member::%s::skill::attributes"
	keyMemberSkillQueue      = "athena::member::%s::skillqueue"
	keyMemberSkills          = "athena::member::%s::skills"
	keyMemberSkillMeta       = "athena::member::%s::skill::meta"
)

func (s *service) MemberSkillAttributes(ctx context.Context, memberID string) (*athena.MemberSkillAttributes, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberSkillAttributes, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var attributes = new(athena.MemberSkillAttributes)

	err = json.Unmarshal(result, attributes)
	if err != nil {
		return nil, err
	}

	return attributes, nil

}

func (s *service) SetMemberSkillAttributes(ctx context.Context, memberID string, attributes *athena.MemberSkillAttributes, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberSkillAttributes, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberSkillQueue(ctx context.Context, memberID string) ([]*athena.MemberSkillQueue, error) {

	key := fmt.Sprintf(keyMemberSkillQueue, memberID)
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
			err = fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		positions[i] = position
	}

	return positions, nil

}

func (s *service) SetMemberSkillQueue(ctx context.Context, memberID string, positions []*athena.MemberSkillQueue, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	members := make([]interface{}, len(positions))
	for i, position := range positions {
		members[i] = position
	}

	key := fmt.Sprintf(keyMemberSkillQueue, memberID)
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

func (s *service) MemberSkills(ctx context.Context, memberID string) ([]*athena.MemberSkill, error) {

	key := fmt.Sprintf(keyMemberSkills, memberID)
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	skills := make([]*athena.MemberSkill, len(members))
	for i, member := range members {
		var skill = new(athena.MemberSkill)
		err = json.Unmarshal([]byte(member), skill)
		if err != nil {
			err = fmt.Errorf("[Cache Layer] Failed to unmarshal set member for key %s onto struct: %w", key, err)
			newrelic.FromContext(ctx).NoticeError(err)
			continue
		}

		skills[i] = skill
	}

	return skills, nil

}

func (s *service) SetMemberSkills(ctx context.Context, memberID string, skills []*athena.MemberSkill, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	members := make([]interface{}, len(skills))
	for i, skill := range skills {
		members[i] = skill
	}

	key := fmt.Sprintf(keyMemberSkills, memberID)
	_, err := s.client.SAdd(ctx, key, members...).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Failed to cache skills for member %s: %w", memberID, err)
	}

	_, err = s.client.Expire(ctx, key, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("[Cache Layer] Field to set expiry on key %s: %w", key, err)
	}

	return nil

}

func (s *service) MemberSkillMeta(ctx context.Context, memberID string) (*athena.MemberSkillMeta, error) {

	key := fmt.Sprintf(keyMemberSkillMeta, memberID)
	result, err := s.client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var meta = new(athena.MemberSkillMeta)

	err = json.Unmarshal(result, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil

}

func (s *service) SetMemberSkillMeta(ctx context.Context, memberID string, meta *athena.MemberSkillMeta, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberSkillMeta, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
