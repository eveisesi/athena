package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eveisesi/athena"
	"github.com/go-redis/redis/v8"
)

type skillService interface {
	MemberSkillAttributes(ctx context.Context, memberID string) (*athena.MemberSkillAttributes, error)
	SetMemberSkillAttributes(ctx context.Context, memberID string, attributes *athena.MemberSkillAttributes, options ...OptionFunc) error
	MemberSkillQueue(ctx context.Context, memberID string) (*athena.MemberSkillQueue, error)
	SetMemberSkillQueue(ctx context.Context, memberID string, skillQueue *athena.MemberSkillQueue, options ...OptionFunc) error
	MemberSkills(ctx context.Context, memberID string) (*athena.MemberSkill, error)
	SetMemberSkills(ctx context.Context, memberID string, skills *athena.MemberSkill, options ...OptionFunc) error
}

const (
	keyMemberSkillAttributes = "athena::member::%s::skill::attributes"
	keyMemberSkillQueue      = "athena::member::%s::skill::queue"
	keyMemberSkills          = "athena::member::%s::skills"
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

func (s *service) MemberSkillQueue(ctx context.Context, memberID string) (*athena.MemberSkillQueue, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberSkillQueue, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var skillQueue = new(athena.MemberSkillQueue)

	err = json.Unmarshal(result, skillQueue)
	if err != nil {
		return nil, err
	}

	return skillQueue, nil

}

func (s *service) SetMemberSkillQueue(ctx context.Context, memberID string, skillQueue *athena.MemberSkillQueue, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(skillQueue)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberSkillQueue, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}

func (s *service) MemberSkills(ctx context.Context, memberID string) (*athena.MemberSkill, error) {

	result, err := s.client.Get(ctx, fmt.Sprintf(keyMemberSkills, memberID)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var skills = new(athena.MemberSkill)

	err = json.Unmarshal(result, skills)
	if err != nil {
		return nil, err
	}

	return skills, nil

}

func (s *service) SetMemberSkills(ctx context.Context, memberID string, skills *athena.MemberSkill, optionFuncs ...OptionFunc) error {

	options := applyOptionFuncs(nil, optionFuncs)

	data, err := json.Marshal(skills)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	_, err = s.client.Set(ctx, fmt.Sprintf(keyMemberSkills, memberID), data, options.expiry).Result()
	if err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil

}
