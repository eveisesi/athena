package skill

import (
	"context"
	"database/sql"
	"errors"
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
)

type Service interface {
	EmptyMemberSkills(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberSkills(ctx context.Context, member *athena.Member) (*athena.MemberSkills, *athena.Etag, error)
	EmptyMemberSkillQueue(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberSkillQueue(ctx context.Context, member *athena.Member) ([]*athena.MemberSkillQueue, *athena.Etag, error)
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

func (s *service) EmptyMemberSkills(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterSkills, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("[Skills Service] Failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberSkills(ctx, member)

	return etag, err

}

func (s *service) MemberSkills(ctx context.Context, member *athena.Member) (*athena.MemberSkills, *athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterSkills, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch etag object: %w", err)
	}

	cached := true
	upsert := false

	properties, err := s.cache.MemberSkillProperties(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if properties == nil {
		cached = false
		properties, err = s.skills.MemberSkillProperties(ctx, member.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch member skill properties from db: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			upsert = true
			properties = &athena.MemberSkills{MemberID: member.ID}
			err = s.esi.ResetEtag(ctx, etag)
			if err != nil {
				return nil, nil, err
			}
		}

	}

	properties.Skills, err = s.cache.MemberSkills(ctx, member.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch member skills from db: %w", err)
	}

	if properties.Skills == nil || len(properties.Skills) == 0 {
		cached = false
		properties.Skills, err = s.skills.MemberSkills(ctx, member.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch member skills from db: %w", err)
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && len(properties.Skills) > 0 {

		if !cached {
			err = s.cache.SetMemberSkills(ctx, member.ID, properties.Skills)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}

			t := properties.Skills
			properties.Skills = nil
			err = s.cache.SetMemberSkillProperties(ctx, member.ID, properties)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
			properties.Skills = t
		}

		return properties, etag, nil
	}
	oldSkills := properties.Skills
	newProperties, etag, _, err := s.esi.GetCharacterSkills(ctx, member, properties)
	if err != nil {
		return nil, nil, fmt.Errorf("[Skills Service] Failed to fetch skills for member %d: %w", member.ID, err)
	}

	newSkills := newProperties.Skills

	s.resolveSkillAttributes(ctx, newSkills)
	properties.Skills = nil

	switch upsert {
	case true:
		properties, err = s.skills.CreateMemberSkillProperties(ctx, newProperties)
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to update skill properties for member %d: %w", member.ID, err)
		}
	case false:
		properties, err = s.skills.UpdateMemberSkillProperties(ctx, member.ID, newProperties)
		if err != nil {
			return nil, nil, fmt.Errorf("[Skills Service] Failed to update skill properties for member %d: %w", member.ID, err)
		}
	}
	_ = s.cache.SetMemberSkillProperties(ctx, member.ID, properties)

	s.resolveSkillAttributes(ctx, newSkills)
	skills, err := s.diffAndUpdateSkills(ctx, member, oldSkills, newSkills)
	if err != nil {
		return nil, nil, fmt.Errorf("[Skills Service] Failed to update skill properties for member %d: %w", member.ID, err)
	}

	properties.Skills = skills

	return properties, etag, nil

}

func (s *service) resolveSkillAttributes(ctx context.Context, skills []*athena.Skill) {

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

func (s *service) diffAndUpdateSkills(ctx context.Context, member *athena.Member, old []*athena.Skill, new []*athena.Skill) ([]*athena.Skill, error) {

	skillsToCreate := make([]*athena.Skill, 0)
	skillsToUpdate := make([]*athena.Skill, 0)

	oldSkillMap := make(map[uint]*athena.Skill)
	for _, skill := range old {
		oldSkillMap[skill.SkillID] = skill
	}

	for _, skill := range new {
		// Never seen this skill before, add it to the db
		if _, ok := oldSkillMap[skill.SkillID]; !ok {
			skillsToCreate = append(skillsToCreate, skill)

			// We've seen this skill before, check to see if the values are still the same
		} else if diff := deep.Equal(oldSkillMap[skill.SkillID], skill); len(diff) > 0 {
			skillsToUpdate = append(skillsToUpdate, skill)
		}
	}

	var final = make([]*athena.Skill, 0)
	if len(skillsToCreate) > 0 {
		createdSkills, err := s.skills.CreateMemberSkills(ctx, member.ID, skillsToCreate)
		if err != nil {
			return nil, err
		}
		final = append(final, createdSkills...)
	}

	if len(skillsToUpdate) > 0 {
		updatedSkill, err := s.skills.UpdateMemberSkills(ctx, member.ID, skillsToUpdate)
		if err != nil {
			return nil, err
		}
		final = append(final, updatedSkill...)
	}

	return final, nil

}

func (s *service) EmptyMemberSkillQueue(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterSkillQueue, esi.ModWithMember(member))
	if err != nil {
		return nil, fmt.Errorf("[Skills Service] Failed to fetch etag object: %w", err)
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) {
		return etag, nil
	}

	_, etag, err = s.MemberSkillQueue(ctx, member)

	return etag, err

}

func (s *service) MemberSkillQueue(ctx context.Context, member *athena.Member) ([]*athena.MemberSkillQueue, *athena.Etag, error) {

	etag, err := s.esi.Etag(ctx, esi.GetCharacterSkillQueue, esi.ModWithMember(member))
	if err != nil {
		return nil, nil, fmt.Errorf("[Skill Service] Failed to fetch etag object: %w", err)
	}

	cached := true

	positions, err := s.cache.MemberSkillQueue(ctx, member.ID)
	if err != nil {
		return nil, nil, err
	}

	if positions == nil {
		cached = false
		positions, err = s.skills.MemberSkillQueue(ctx, member.ID)
		if err != nil {
			return nil, nil, err
		}
	}

	if etag != nil && etag.CachedUntil.After(time.Now()) && len(positions) > 0 {
		if !cached {
			err = s.cache.SetMemberSkillQueue(ctx, member.ID, positions)
			if err != nil {
				newrelic.FromContext(ctx).NoticeError(err)
			}
		}

		return positions, etag, nil
	}

	newPositions, etag, _, err := s.esi.GetCharacterSkillQueue(ctx, member, make([]*athena.MemberSkillQueue, 0))
	if err != nil {
		return nil, nil, fmt.Errorf("[Skill Service] Failed to fetch skillQueue for member %d: %w", member.ID, err)
	}

	s.resolveSkillQueueAttributes(ctx, newPositions)
	positions, err = s.diffAndUpdateSkillQueue(ctx, member, positions, newPositions)
	if err != nil {
		return nil, nil, fmt.Errorf("[Skill Service] Failed to execute diffing options: %w", err)
	}

	if len(positions) > 0 {
		err = s.cache.SetMemberSkillQueue(ctx, member.ID, positions)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
		}
	}

	return positions, etag, nil

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
	positionsToDelete := make([]uint, 0)

	oldQueueMap := make(map[uint]*athena.MemberSkillQueue)
	for _, position := range old {
		oldQueueMap[position.QueuePosition] = position
	}

	for _, position := range new {
		var ok bool
		// This is an unknown position, so lets flag it to be created
		if _, ok = oldQueueMap[position.QueuePosition]; !ok {
			positionsToCreate = append(positionsToCreate, position)

			// We've seen this position before for this member, lets compare it to the existing position to see
			// if it needs to be updated
		} else if diff := deep.Equal(oldQueueMap[position.QueuePosition], position); len(diff) > 0 {
			positionsToUpdate = append(positionsToUpdate, position)
		}
	}

	newSkillQueueMap := make(map[uint]*athena.MemberSkillQueue)
	for _, position := range new {
		newSkillQueueMap[position.QueuePosition] = position
	}

	for _, position := range old {
		// This label is not in the list of new label, must've been deleted by the user in game
		if _, ok := newSkillQueueMap[position.QueuePosition]; !ok {
			positionsToDelete = append(positionsToDelete, position.QueuePosition)
		}
	}

	if len(positionsToDelete) > 0 {
		for _, position := range positionsToDelete {
			deleteOK, err := s.skills.DeleteMemberSkillQueuePosition(ctx, member.ID, position)
			if err != nil {
				return nil, fmt.Errorf("[Skill Service] Failed to delete member skill queue positions: %w", err)
			}

			if !deleteOK {
				return nil, fmt.Errorf("[Skill Service] Expected to delete %d documents, deleted none", len(positionsToDelete))
			}
		}
	}

	var final = make([]*athena.MemberSkillQueue, 0)
	if len(positionsToCreate) > 0 {
		createdPositions, err := s.skills.CreateMemberSkillQueue(ctx, member.ID, positionsToCreate)
		if err != nil {
			return nil, fmt.Errorf("[Skill Service] Failed to create member skill queue positions: %w", err)
		}
		final = append(final, createdPositions...)
	}

	if len(positionsToUpdate) > 0 {
		for _, position := range positionsToUpdate {
			updatedPosition, err := s.skills.UpdateMemberSkillQueue(ctx, member.ID, position)
			if err != nil {
				return nil, fmt.Errorf("[Skill Service] Failed to update member skill queue positions: %w", err)
			}
			final = append(final, updatedPosition)
		}
	}

	return final, nil

}
