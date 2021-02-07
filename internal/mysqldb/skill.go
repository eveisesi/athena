package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type skillRepository struct {
	db *sqlx.DB
	skillProperties,
	skills,
	skillqueue,
	attributes string
}

func NewSkillRepository(db *sql.DB) athena.MemberSkillRepository {

	return &skillRepository{
		db:              sqlx.NewDb(db, "mysql"),
		attributes:      "member_attributes",
		skillProperties: "member_skill_properties",
		skills:          "member_skills",
		skillqueue:      "member_skillqueue",
	}

}

func (r *skillRepository) MemberSkillProperties(ctx context.Context, id uint) (*athena.MemberSkills, error) {

	query, args, err := sq.Select(
		"member_id", "total_sp", "unallocated_sp", "created_at", "updated_at",
	).From(r.skillProperties).Where(sq.Eq{"member_id": id}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	var properties = new(athena.MemberSkills)
	err = r.db.GetContext(ctx, properties, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to fetch member skill properties: %w", err)
	}

	return properties, nil

}

func (r *skillRepository) CreateMemberSkillProperties(ctx context.Context, properties *athena.MemberSkills) (*athena.MemberSkills, error) {

	query, args, err := sq.Insert(r.skillProperties).Columns(
		"member_id", "total_sp", "unallocated_sp", "created_at", "updated_at",
	).Values(
		properties.MemberID, properties.TotalSP, properties.UnallocatedSP,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to create member skill properties: %w", err)
	}

	return r.MemberSkillProperties(ctx, properties.MemberID)

}

func (r *skillRepository) UpdateMemberSkillProperties(ctx context.Context, id uint, properties *athena.MemberSkills) (*athena.MemberSkills, error) {

	query, args, err := sq.Update(r.skillProperties).
		Set("total_sp", properties.TotalSP).
		Set("unallocated_sp", properties.UnallocatedSP).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to update member skill properties: %w", err)
	}

	return r.MemberSkillProperties(ctx, id)

}

func (r *skillRepository) DeleteMemberSkillProperties(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.skillProperties).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to delete member skill properties: %w", err)
	}

	return true, nil

}

func (r *skillRepository) MemberSkills(ctx context.Context, id uint) ([]*athena.Skill, error) {

	query, args, err := sq.Select(
		"active_skill_level", "skill_id",
		"skillpoints_in_skill", "trained_skill_level",
		"created_at", "updated_at",
	).From(r.skills).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	var skills = make([]*athena.Skill, 0)
	err = r.db.SelectContext(ctx, &skills, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to fetch member skills: %w", err)
	}

	return skills, nil

}

func (r *skillRepository) CreateMemberSkills(ctx context.Context, id uint, skills []*athena.Skill) ([]*athena.Skill, error) {

	i := sq.Insert(r.skills).Columns(
		"member_id",
		"active_skill_level", "skill_id",
		"skillpoints_in_skill", "trained_skill_level",
		"created_at", "updated_at",
	)
	for _, skill := range skills {
		i = i.Values(
			id,
			skill.ActiveSkillLevel, skill.SkillID,
			skill.SkillpointsInSkill, skill.TrainedSkillLevel,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert member skills: %w", err)
	}

	return r.MemberSkills(ctx, id)

}

func (r *skillRepository) UpdateMemberSkills(ctx context.Context, id uint, skills []*athena.Skill) ([]*athena.Skill, error) {

	for _, skill := range skills {

		query, args, err := sq.Update(r.skills).
			Set("active_skill_level", skill.ActiveSkillLevel).
			Set("skillpoints_in_skill", skill.SkillpointsInSkill).
			Set("trained_skill_level", skill.TrainedSkillLevel).
			Set("updated_at", sq.Expr(`NOW()`)).
			Where(sq.Eq{"member_id": id, "skill_id": skill.SkillID}).ToSql()
		if err != nil {
			return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("[Skill Repository] Failed to update member skills: %w", err)
		}

	}

	return r.MemberSkills(ctx, id)

}

func (r *skillRepository) DeleteMemberSkills(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.skills).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to delete member skills: %w", err)
	}

	return true, nil

}

func (r *skillRepository) MemberSkillQueue(ctx context.Context, id uint) ([]*athena.MemberSkillQueue, error) {

	query, args, err := sq.Select(
		"member_id", "queue_position", "skill_id",
		"finished_level", "training_start_sp", "level_start_sp",
		"level_end_sp", "start_date", "finish_date",
		"created_at", "updated_at",
	).From(r.skillqueue).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	var skills = make([]*athena.MemberSkillQueue, 0)
	err = r.db.SelectContext(ctx, &skills, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to fetch member skillqueue: %w", err)
	}

	return skills, nil

}
func (r *skillRepository) MemberSkillQueuePosition(ctx context.Context, id uint, position uint) (*athena.MemberSkillQueue, error) {

	query, args, err := sq.Select(
		"member_id", "queue_position", "skill_id",
		"finished_level", "training_start_sp", "level_start_sp",
		"level_end_sp", "start_date", "finish_date",
		"created_at", "updated_at",
	).From(r.skillqueue).Where(sq.Eq{"member_id": id, "queue_position": position}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	var skill = new(athena.MemberSkillQueue)
	err = r.db.GetContext(ctx, skill, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to fetch member skillqueue: %w", err)
	}

	return skill, nil

}

func (r *skillRepository) CreateMemberSkillQueue(ctx context.Context, id uint, positions []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, error) {

	i := sq.Insert(r.skillqueue).Columns(
		"member_id", "queue_position", "skill_id",
		"finished_level", "training_start_sp", "level_start_sp",
		"level_end_sp", "start_date", "finish_date",
		"created_at", "updated_at",
	)
	for _, position := range positions {
		i = i.Values(
			id, position.QueuePosition, position.SkillID,
			position.FinishedLevel, position.TrainingStartSp, position.LevelStartSp,
			position.LevelEndSp, position.StartDate, position.FinishDate,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert member skillqueue: %w", err)
	}

	return r.MemberSkillQueue(ctx, id)

}

func (r *skillRepository) UpdateMemberSkillQueue(ctx context.Context, id uint, position *athena.MemberSkillQueue) (*athena.MemberSkillQueue, error) {

	query, args, err := sq.Update(r.skillqueue).
		Set("skill_id", position.SkillID).
		Set("finished_level", position.FinishedLevel).
		Set("training_start_sp", position.TrainingStartSp).
		Set("level_start_sp", position.LevelStartSp).
		Set("level_end_sp", position.LevelEndSp).
		Set("start_date", position.StartDate).
		Set("finish_date", position.FinishDate).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": id, "queue_position": position.QueuePosition}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert member skillqueue: %w", err)
	}

	return r.MemberSkillQueuePosition(ctx, id, position.QueuePosition)

}

func (r *skillRepository) DeleteMemberSkillQueue(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.skillqueue).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to delete member skillqueue: %w", err)
	}

	return true, nil

}

func (r *skillRepository) DeleteMemberSkillQueuePosition(ctx context.Context, id, position uint) (bool, error) {

	query, args, err := sq.Delete(r.skillqueue).Where(sq.Eq{"member_id": id, "queue_position": position}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to delete member skillqueue: %w", err)
	}

	return true, nil

}

func (r *skillRepository) MemberAttributes(ctx context.Context, id uint) (*athena.MemberAttributes, error) {

	query, args, err := sq.Select(
		"member_id", "charisma", "intelligence",
		"memory", "perception", "willpower",
		"bonus_remaps", "last_remap_date", "accrued_remap_cooldown_date",
		"created_at", "updated_at",
	).From(r.attributes).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	var attributes = new(athena.MemberAttributes)
	err = r.db.GetContext(ctx, &attributes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to fetch member attributes: %w", err)
	}

	return attributes, nil

}

func (r *skillRepository) CreateMemberAttributes(ctx context.Context, attributes *athena.MemberAttributes) (*athena.MemberAttributes, error) {

	i := sq.Insert(r.attributes).Columns(
		"member_id", "charisma", "intelligence",
		"memory", "perception", "willpower",
		"bonus_remaps", "last_remap_date", "accrued_remap_cooldown_date",
		"created_at", "updated_at",
	).Values(
		attributes.MemberID, attributes.Charisma, attributes.Intelligence,
		attributes.Memory, attributes.Perception, attributes.Willpower,
		attributes.BonusRemaps, attributes.LastRemapDate, attributes.AccruedRemapCooldownDate,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	)

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Skill Repository] Failed to insert member attributes: %w", err)
	}

	return r.MemberAttributes(ctx, attributes.MemberID)

}

func (r *skillRepository) UpdateMemberAttributes(ctx context.Context, id uint, attributes *athena.MemberAttributes) error {

	query, args, err := sq.Update(r.attributes).
		Set("charisma", attributes.Charisma).
		Set("intelligence", attributes.Intelligence).
		Set("memory", attributes.Memory).
		Set("perception", attributes.Perception).
		Set("willpower", attributes.Willpower).
		Set("bonus_remaps", attributes.BonusRemaps).
		Set("last_remap_date", attributes.LastRemapDate).
		Set("accrued_remap_cooldown_date", attributes.AccruedRemapCooldownDate).
		Set("updated_at", sq.Expr(`NOW()`)).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("[Skill Repository] Failed to update member attributes: %w", err)
	}

	return nil

}

func (r *skillRepository) DeleteMemberAttributes(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.attributes).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Skill Repository] Failed to delete member attributes: %w", err)
	}

	return true, nil

}
