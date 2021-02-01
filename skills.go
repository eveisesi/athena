package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type MemberSkillRepository interface {
	memberAttributesRepository
	memberSkillQueueRepository
	memberSkillMetaRepository
	memberSkillsRepository
}

type memberAttributesRepository interface {
	MemberSkillAttributes(ctx context.Context, id uint) (*MemberSkillAttributes, error)
	CreateMemberSkillAttributes(ctx context.Context, location *MemberSkillAttributes) (*MemberSkillAttributes, error)
	UpdateMemberSkillAttributes(ctx context.Context, id uint, location *MemberSkillAttributes) (*MemberSkillAttributes, error)
	DeleteMemberSkillAttributes(ctx context.Context, id uint) (bool, error)
}

type memberSkillQueueRepository interface {
	MemberSkillQueue(ctx context.Context, memberID uint) ([]*MemberSkillQueue, error)
	CreateMemberSkillQueue(ctx context.Context, memberID uint, positions []*MemberSkillQueue) ([]*MemberSkillQueue, error)
	UpdateMemberSkillQueue(ctx context.Context, memberID uint, position *MemberSkillQueue) (*MemberSkillQueue, error)
	DeleteMemberSkillQueue(ctx context.Context, memberID uint, positions []*MemberSkillQueue) (bool, error)
}

type memberSkillMetaRepository interface {
	MemberSkillMeta(ctx context.Context, memberID uint) (*MemberSkillMeta, error)
	CreateMemberSkillMeta(ctx context.Context, meta *MemberSkillMeta) (*MemberSkillMeta, error)
	UpdateMemberSkillMeta(ctx context.Context, memberID uint, meta *MemberSkillMeta) (*MemberSkillMeta, error)
	DeleteMemberSkillMeta(ctx context.Context, memberID uint) (bool, error)
}

type memberSkillsRepository interface {
	MemberSkills(ctx context.Context, memberID uint) ([]*MemberSkill, error)
	CreateMemberSkills(ctx context.Context, memberID uint, skills []*MemberSkill) ([]*MemberSkill, error)
	UpdateMemberSkills(ctx context.Context, memberID uint, skills *MemberSkill) (*MemberSkill, error)
	DeleteMemberSkills(ctx context.Context, memberID uint) (bool, error)
}

type MemberSkillAttributes struct {
	MemberID                 uint      `db:"member_id" json:"member_id"`
	Charisma                 uint      `db:"charisma" json:"charisma"`
	Intelligence             uint      `db:"intelligence" json:"intelligence"`
	Memory                   uint      `db:"memory" json:"memory"`
	Perception               uint      `db:"perception" json:"perception"`
	Willpower                uint      `db:"willpower" json:"willpower"`
	BonusRemaps              null.Int  `db:"bonus_remaps,omitempty" json:"bonus_remaps,omitempty"`
	LastRemapDate            null.Time `db:"last_remap_date,omitempty" json:"last_remap_date,omitempty"`
	AccruedRemapCooldownDate null.Time `db:"accrued_remap_cooldown_date,omitempty" json:"accrued_remap_cooldown_date,omitempty"`
	Meta
}

type MemberSkillQueue struct {
	MemberID        uint      `db:"member_id" json:"member_id" deep:"-"`
	SkillID         uint      `db:"skill_id" json:"skill_id"`
	QueuePosition   uint      `db:"queue_position" json:"queue_position"`
	FinishedLevel   uint      `db:"finished_level" json:"finished_level"`
	TrainingStartSp null.Int  `db:"training_start_sp,omitempty" json:"training_start_sp,omitempty"`
	LevelStartSp    null.Int  `db:"level_start_sp,omitempty" json:"level_start_sp,omitempty"`
	LevelEndSp      null.Int  `db:"level_end_sp,omitempty" json:"level_end_sp,omitempty"`
	StartDate       null.Time `db:"start_date,omitempty" json:"start_date,omitempty"`
	FinishDate      null.Time `db:"finish_date,omitempty" json:"finish_date,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}

func (m *MemberSkillQueue) Valid() bool {
	return m.SkillID > 0 && m.QueuePosition > 0
}

type MemberSkillMeta struct {
	MemberID      uint           `db:"member_id" json:"member_id" deep:"-"`
	TotalSP       uint           `db:"total_sp" json:"total_sp"`
	Skills        []*MemberSkill `db:"-" json:"skills"`
	UnallocatedSP null.Int       `db:"unallocated_sp,omitempty" json:"unallocated_sp,omitempty"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at" deep:"-"`
}

func (m *MemberSkillMeta) Valid() bool {
	return m.TotalSP > 0
}

type MemberSkill struct {
	MemberID           uint      `db:"member_id" json:"member_id" deep:"-"`
	ActiveSkillLevel   uint      `db:"active_skill_level" json:"active_skill_level"`
	SkillID            uint      `db:"skill_id" json:"skill_id"`
	SkillpointsInSkill uint      `db:"skillpoints_in_skill" json:"skillpoints_in_skill"`
	TrainedSkillLevel  uint      `db:"trained_skill_level" json:"trained_skill_level"`
	CreatedAt          time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}

func (m *MemberSkill) Valid() bool {
	return m.SkillID > 0
}
