package skill

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
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

type Service interface {
	EmptyMemberSkills(ctx context.Context, member *athena.Member) error
	MemberSkills(ctx context.Context, member *athena.Member) (*athena.MemberSkillMeta, []*athena.MemberSkill, error)
	EmptyMemberSkillQueue(ctx context.Context, member *athena.Member) error
	MemberSkillQueue(ctx context.Context, member *athena.Member) ([]*athena.MemberSkillQueue, error)
}

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

	var newSkills []*athena.MemberSkill
	if newMeta.Valid() {
		newSkills = newMeta.Skills
		meta, err = s.skills.UpdateMemberSkillMeta(ctx, member.ID.Hex(), newMeta)
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to update skill meta for member %s: %w", member.ID.Hex(), err)
		}
		meta.Skills = nil
		_ = s.cache.SetMemberSkillMeta(ctx, member.ID.Hex(), meta)

	}

	if len(newSkills) > 0 {
		s.resolveSkillAttributes(ctx, newSkills)
		skills, err = s.diffAndUpdateSkills(ctx, member, skills, newSkills)
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to update skill meta for member %s: %w", member.ID.Hex(), err)
		}
	}

	return meta, skills, nil
}

func (s *service) resolveSkillAttributes(ctx context.Context, skills []*athena.MemberSkill) {

	for _, skill := range skills {
		_, err := s.universe.Type(ctx, skill.SkillID)
		if err != nil {
			s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
				"skill_id": skill.SkillID,
			}).Error("failed to resolve skill type")
			continue
		}
	}

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
		for _, skill := range skillsToUpdate {
			updatedSkill, err := s.skills.UpdateMemberSkills(ctx, member.ID.Hex(), skill)
			if err != nil {
				return nil, err
			}
			final = append(final, updatedSkill)
		}
	}

	return final, nil

}

func (s *service) EmptyMemberSkillQueue(ctx context.Context, member *athena.Member) error {

	_, err := s.MemberSkillQueue(ctx, member)

	return err

}

func (s *service) MemberSkillQueue(ctx context.Context, member *athena.Member) ([]*athena.MemberSkillQueue, error) {

	valid, etagID := s.esi.GenerateEndpointHash(esi.EndpointGetCharacterSkillQueue, member)
	if !valid {
		return nil, fmt.Errorf("[Skill Service] Failed to generate valid etag hash")
	}

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, fmt.Errorf("[Skill Service] Failed to fetch etag object: %w", err)
	}

	cached := true

	positions, err := s.cache.MemberSkillQueue(ctx, member.ID.Hex())
	if err != nil {
		return nil, err
	}

	if positions == nil {
		cached = false
		positions, err = s.skills.MemberSkillQueue(ctx, member.ID.Hex())
		if err != nil {
			return nil, err
		}
	}

	if etag.CachedUntil.After(time.Now()) && len(positions) > 0 {
		if !cached {
			err = s.cache.SetMemberSkillQueue(ctx, member.ID.Hex(), positions)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return positions, nil
	}

	newPositions, etag, _, err := s.esi.GetCharacterSkillQueue(ctx, member, etag, make([]*athena.MemberSkillQueue, 0))
	if err != nil {
		return nil, fmt.Errorf("[Skill Service] Failed to fetch skillQueue for member %s: %w", member.ID.Hex(), err)
	}

	_, _ = s.etag.UpdateEtag(ctx, etag.EtagID, etag)

	if len(newPositions) > 0 {
		s.resolveSkillQueueAttributes(ctx, newPositions)
		positions, err = s.diffAndUpdateSkillQueue(ctx, member, positions, newPositions)
		if err != nil {
			return nil, fmt.Errorf("[Skill Service] Failed to execute diffing options: %w", err)
		}

		if len(positions) > 0 {
			err = s.cache.SetMemberSkillQueue(ctx, member.ID.Hex(), positions)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}
	}

	return positions, nil

}

func (s *service) resolveSkillQueueAttributes(ctx context.Context, positions []*athena.MemberSkillQueue) {

	for _, position := range positions {
		_, err := s.universe.Type(ctx, position.SkillID)
		if err != nil {
			s.logger.WithError(err).WithContext(ctx).WithFields(logrus.Fields{
				"skill_id": position.SkillID,
			}).Error("failed to resolve skill type")
			continue
		}
	}

}

func (s *service) diffAndUpdateSkillQueue(ctx context.Context, member *athena.Member, old []*athena.MemberSkillQueue, new []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, error) {

	positionsToCreate := make([]*athena.MemberSkillQueue, 0)
	positionsToUpdate := make([]*athena.MemberSkillQueue, 0)
	positionsToDelete := make([]*athena.MemberSkillQueue, 0)

	oldContactMap := make(map[int]*athena.MemberSkillQueue)
	for _, position := range old {
		oldContactMap[position.QueuePosition] = position
	}

	for _, position := range new {
		var ok bool
		// This is an unknown position, so lets flag it to be created
		if _, ok = oldContactMap[position.QueuePosition]; !ok {
			positionsToCreate = append(positionsToCreate, position)

			// We've seen this position before for this member, lets compare it to the existing position to see
			// if it needs to be updated
		} else if diff := deep.Equal(oldContactMap[position.QueuePosition], position); len(diff) > 0 {
			spew.Dump(position.QueuePosition, position)
			positionsToUpdate = append(positionsToUpdate, position)
		}
	}

	newContactMap := make(map[int]*athena.MemberSkillQueue)
	for _, position := range new {
		newContactMap[position.QueuePosition] = position
	}

	for _, position := range old {
		// This label is not in the list of new label, must've been deleted by the user in game
		if _, ok := newContactMap[position.QueuePosition]; !ok {
			positionsToDelete = append(positionsToDelete, position)
		}
	}

	var final = make([]*athena.MemberSkillQueue, 0)
	if len(positionsToCreate) > 0 {
		createdPositions, err := s.skills.CreateMemberSkillQueue(ctx, member.ID.Hex(), positionsToCreate)
		if err != nil {
			return nil, fmt.Errorf("[Skill Service] Failed to create member skill queue positions: %w", err)
		}
		final = append(final, createdPositions...)
	}

	if len(positionsToUpdate) > 0 {
		for _, position := range positionsToUpdate {
			updatedPosition, err := s.skills.UpdateMemberSkillQueue(ctx, member.ID.Hex(), position)
			if err != nil {
				return nil, fmt.Errorf("[Skill Service] Failed to update member skill queue positions: %w", err)
			}
			final = append(final, updatedPosition)
		}
	}

	if len(positionsToDelete) > 0 {
		deleteOK, err := s.skills.DeleteMemberSkillQueue(ctx, member.ID.Hex(), positionsToDelete)
		if err != nil {
			return nil, fmt.Errorf("[Skill Service] Failed to delete member skill queue positions: %w", err)
		}

		if !deleteOK {
			return nil, fmt.Errorf("[Skill Service] Expected to delete %d documents, deleted none", len(positionsToDelete))
		}
	}

	return final, nil

}
