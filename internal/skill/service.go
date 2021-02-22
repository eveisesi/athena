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
	"github.com/sirupsen/logrus"
)

type Service interface {
	FetchMemberSkills(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberSkillProperties(ctx context.Context, memberID uint) (*athena.MemberSkills, error)
	MemberSkills(ctx context.Context, memberID uint) ([]*athena.Skill, error)
	FetchMemberSkillQueue(ctx context.Context, member *athena.Member) (*athena.Etag, error)
	MemberSkillQueue(ctx context.Context, memberID uint) ([]*athena.MemberSkillQueue, error)
}

type service struct {
	logger *logrus.Logger

	cache    cache.Service
	esi      esi.Service
	etag     etag.Service
	universe universe.Service

	skills athena.MemberSkillRepository
}

const (
	serviceIdentifier = "Skill Service"
)

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

func (s *service) FetchMemberSkills(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "MemberSkills",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterSkills, esi.ModWithCharacterID(member.ID))
	if err != nil {
		entry.WithError(err).Error("failed to fetch etag object")
		return nil, fmt.Errorf("failed to fetch etag object")
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {
			return etag, nil
		}

		petag = etag.Etag
	}

	skillProperties, etag, _, err := s.esi.GetCharacterSkills(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch skills for member")
		return nil, fmt.Errorf("failed to fetch skills for member")
	}

	if petag != "" && etag.Etag == petag {
		return etag, nil
	}

	s.processSkills(ctx, member, skillProperties.Skills)
	skillProperties.MemberID = member.ID
	skillProperties.Skills = nil

	existing, err := s.skills.MemberSkillProperties(ctx, member.ID)

	switch existing == nil || errors.Is(err, sql.ErrNoRows) {
	case true:
		_, err := s.skills.CreateMemberSkillProperties(ctx, skillProperties)
		if err != nil {
			entry.WithError(err).Error("failed to create member skill properties in db")
			return nil, fmt.Errorf("failed to create member skill properties in db")
		}
	case false:
		_, err := s.skills.UpdateMemberSkillProperties(ctx, member.ID, skillProperties)
		if err != nil {
			entry.WithError(err).Error("failed to update member skill properties in db")
			return nil, fmt.Errorf("failed to update member skill properties in db")
		}
	}

	return etag, nil

}

func (s *service) processSkills(ctx context.Context, member *athena.Member, skills []*athena.Skill) error {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "processSkills",
	})

	s.resolveSkillAttributes(ctx, skills)

	existing, err := s.skills.MemberSkills(ctx, member.ID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch existing skills from DB")
		return fmt.Errorf("failed to fetch existing skills from DB")
	}

	skills, err = s.diffAndUpdateSkills(ctx, member, existing, skills)
	if err != nil {
		entry.WithError(err).Error("failed to diff and update skills")
		return fmt.Errorf("failed to diff and update skills")
	}

	if len(skills) > 0 {
		err = s.cache.SetMemberSkills(ctx, member.ID, skills)
		if err != nil {
			entry.WithError(err).Error("failed to cache member skills")
		}
	}

	return nil

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

func (s *service) MemberSkillProperties(ctx context.Context, memberID uint) (*athena.MemberSkills, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberSkillQueue",
	})

	properties, err := s.cache.MemberSkillProperties(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member skill properties from cache")
		return nil, fmt.Errorf("failed to fetch member skill properties from cache")
	}

	if properties != nil {
		return properties, nil
	}

	properties, err = s.skills.MemberSkillProperties(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member skill properties from DB")
		return nil, fmt.Errorf("failed to fetch member skill properties from DB")
	}

	if err == nil {
		err = s.cache.SetMemberSkillProperties(ctx, memberID, properties)
		if err != nil {
			entry.WithError(err).Error("failed to cache member skill properties")
		}
	}

	return properties, nil

}

func (s *service) MemberSkills(ctx context.Context, memberID uint) ([]*athena.Skill, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberSkills",
	})

	skills, err := s.cache.MemberSkills(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member skills from cache")
		return nil, fmt.Errorf("failed to fetch member skills from cache")
	}

	if len(skills) > 0 {
		return skills, nil
	}

	skills, err = s.skills.MemberSkills(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member skills from DB")
		return nil, fmt.Errorf("failed to fetch member skills from DB")
	}

	if len(skills) > 0 {
		err = s.cache.SetMemberSkills(ctx, memberID, skills)
		if err != nil {
			entry.WithError(err).Error("failed to cache member skills")
		}
	}

	return skills, nil

}

func (s *service) FetchMemberSkillQueue(ctx context.Context, member *athena.Member) (*athena.Etag, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": member.ID,
		"service":   serviceIdentifier,
		"method":    "MemberSkillQueue",
	})

	etag, err := s.esi.Etag(ctx, esi.GetCharacterSkillQueue, esi.ModWithCharacterID(member.ID))
	if err != nil {
		return nil, fmt.Errorf("[Skills Service] Failed to fetch etag object: %w", err)
	}

	var petag string
	if etag != nil {
		if etag.CachedUntil.After(time.Now()) {
			return etag, nil
		}

		petag = etag.Etag
	}

	newPositions, etag, _, err := s.esi.GetCharacterSkillQueue(ctx, member.ID, member.AccessToken.String)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member skill queue from ESI")
		return nil, fmt.Errorf("failed to fetch member skill queue from ESI")
	}

	if petag != "" && etag.Etag == petag {
		return etag, nil
	}

	existing, err := s.skills.MemberSkillQueue(ctx, member.ID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member skill queue from db")
		return nil, fmt.Errorf("failed to fetch member skill queue from db")
	}

	s.resolveSkillQueueAttributes(ctx, newPositions)
	positions, err := s.diffAndUpdateSkillQueue(ctx, member, existing, newPositions)
	if err != nil {
		return nil, fmt.Errorf("[Skill Service] Failed to execute diffing options: %w", err)
	}

	if len(positions) > 0 {
		err = s.cache.SetMemberSkillQueue(ctx, member.ID, positions)
		if err != nil {
			entry.WithError(err).Error("failed to cache member skill queue")
		}
	}

	return etag, nil

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

func (s *service) MemberSkillQueue(ctx context.Context, memberID uint) ([]*athena.MemberSkillQueue, error) {

	entry := s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"member_id": memberID,
		"service":   serviceIdentifier,
		"method":    "MemberSkillQueue",
	})

	positions, err := s.cache.MemberSkillQueue(ctx, memberID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch member skill queue from cache")
		return nil, fmt.Errorf("failed to fetch member skill queue from cache")
	}

	if len(positions) > 0 {
		return positions, nil
	}

	positions, err = s.skills.MemberSkillQueue(ctx, memberID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		entry.WithError(err).Error("failed to fetch member skill queue from DB")
		return nil, fmt.Errorf("failed to fetch member skill queue from DB")
	}

	if len(positions) > 0 {
		err = s.cache.SetMemberSkillQueue(ctx, memberID, positions)
		if err != nil {
			entry.WithError(err).Error("failed to cache member skill queue")
		}
	}

	return positions, nil

}
