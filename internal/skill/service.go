package skill

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/esi"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/eveisesi/athena/internal/universe"
	"github.com/go-test/deep"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface{}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	etag     etag.Service
	universe universe.Service

	skills athena.MemberSkillRepository
}

func NewService(logger *logrus.Logger, cache cache.Service, esi esi.Service, etag etag.Service, universe universe.Service, skills athena.MemberSkillRepository) Service {
	return &service{
		logger: logger,

		cache:    cache,
		esi:      esi,
		etag:     etag,
		universe: universe,

		skills: skills,
	}
}

func (s *service) EmptyMemberSkills(ctx context.Context, member *athena.Member) error {

	_, _, err := s.MemberSkills(ctx, member)

	return err

}

func (s *service) MemberSkills(ctx context.Context, member *athena.Member) (*athena.MemberSkillMeta, []*athena.MemberSkill, error) {

	valid, etagID := s.esi.GenerateEndpointHash(esi.EndpointGetCharacterSkills, member)
	if !valid {
		return nil, nil, fmt.Errorf("[Skills Service] Failed to generate valid etag hash")
	}

	cached := true

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch etag object: %w", err)
	}

	meta, err := s.cache.MemberSkillMeta(ctx, member.ID.Hex())
	if err != nil {
		return nil, nil, err
	}

	skills, err := s.cache.MemberSkills(ctx, member.ID.Hex())
	if err != nil {
		return nil, nil, err
	}

	if meta == nil || skills == nil || len(skills) == 0 {
		cached = false
		meta, err = s.skills.MemberSkillMeta(ctx, member.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch member skill meta from db: %w", err)
		}

		if err == mongo.ErrNoDocuments {
			meta = &athena.MemberSkillMeta{MemberID: member.ID}
		}

		skills, err = s.skills.MemberSkills(ctx, member.ID.Hex())
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch member skills from db: %w", err)
		}

	}

	if etag.CachedUntil.After(time.Now()) && len(skills) > 0 && meta != nil {

		if !cached {
			err = s.cache.SetMemberSkillMeta(ctx, member.ID.Hex(), meta)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}

			err = s.cache.SetMemberSkills(ctx, member.ID.Hex(), skills)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return meta, skills, nil
	}

	newMeta, etag, _, err := s.esi.GetCharacterSkills(ctx, member, etag, &athena.MemberSkillMeta{})
	if err != nil {
		return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch skills for member %s: %w", member.ID.Hex(), err)
	}

	_, _ = s.etag.UpdateEtag(ctx, etagID, etag)

	if newMeta.Valid() {
		meta, err = s.skills.UpdateMemberSkillMeta(ctx, member.ID.Hex(), newMeta)
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to update skill meta for member %s: %w", member.ID.Hex(), err)
		}

		_ = s.cache.SetMemberSkillMeta(ctx, member.ID.Hex(), meta)

	}

	if len(newMeta.Skills) > 0 {
		skills, err = s.diffAndUpdateSkills(ctx, member, skills, newMeta.Skills)
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to update skill meta for member %s: %w", member.ID.Hex(), err)
		}
	}

	return meta, skills, nil
}

func (s *service) diffAndUpdateSkills(ctx context.Context, member *athena.Member, old []*athena.MemberSkill, new []*athena.MemberSkill) ([]*athena.MemberSkill, error) {

	skillsToCreate := make([]*athena.MemberSkill, 0)
	skillsToUpdate := make([]*athena.MemberSkill, 0)

	oldLabelMap := make(map[int]*athena.MemberSkill)
	for _, skill := range old {
		oldLabelMap[skill.SkillID] = skill
	}

	for _, skill := range new {
		// Never seen this skill before, add it to the db
		if _, ok := oldLabelMap[skill.SkillID]; !ok {
			skillsToCreate = append(skillsToCreate, skill)

			// We've seen this skill before, check to see if the values are still the same
		} else if diff := deep.Equal(oldLabelMap[skill.SkillID], skill); len(diff) > 0 {
			skillsToUpdate = append(skillsToUpdate, skill)
		}
	}

	var final = make([]*athena.MemberSkill, 0)
	if len(skillsToCreate) > 0 {
		createdSkills, err := s.skills.CreateMemberSkills(ctx, member.ID.Hex(), skillsToCreate)
		if err != nil {
			return nil, err
		}
		final = append(final, createdSkills...)
	}

	if len(skillsToUpdate) > 0 {
		updatedSkills, err := s.skills.UpdateMemberSkills(ctx, member.ID.Hex(), skillsToUpdate)
		if err != nil {
			return nil, err
		}
		final = append(final, updatedSkills...)
	}

	return final, nil

}

func (s *service) MemberSkillQueue(ctx context.Context, member *athena.Member) ([]*athena.MemberSkillQueue, error) {

	valid, etagID := s.esi.GenerateEndpointHash(esi.EndpointGetCharacterSkillQueue, member)
	if !valid {
		return nil, fmt.Errorf("failed to generate valid etag hash")
	}

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch etag object: %w", err)
	}

	cached := true

	skillQueue, err := s.cache.MemberSkillQueue(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if skillQueue == nil {
		cached = false
		skillQueue, err = s.skills.MemberSkillQueue(ctx, member.ID.Hex())
		if err != nil {
			return nil, err
		}
	}

	if etag.CachedUntil.After(time.Now()) && len(skillQueue) > 0 {
		if !cached {
			err = s.cache.SetMemberContactLabels(ctx, member.ID.Hex(), labels)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return labels, nil
	}

}
