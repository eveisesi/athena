package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type MemberSkillRepository interface {
	memberAttributesRepository
	memberSkillsRepository
	memberSkillQueueRepository
}

type memberAttributesRepository interface {
	MemberAttributes(ctx context.Context, id uint) (*MemberAttributes, error)
	CreateMemberAttributes(ctx context.Context, attributes *MemberAttributes) (*MemberAttributes, error)
	UpdateMemberAttributes(ctx context.Context, id uint, attributes *MemberAttributes) error
	DeleteMemberAttributes(ctx context.Context, id uint) (bool, error)
}

type memberSkillsRepository interface {
	MemberSkillProperties(ctx context.Context, id uint) (*MemberSkills, error)
	CreateMemberSkillProperties(ctx context.Context, properties *MemberSkills) (*MemberSkills, error)
	UpdateMemberSkillProperties(ctx context.Context, id uint, properties *MemberSkills) (*MemberSkills, error)
	DeleteMemberSkillProperties(ctx context.Context, id uint) (bool, error)
	MemberSkills(ctx context.Context, id uint) ([]*Skill, error)
	CreateMemberSkills(ctx context.Context, id uint, skills []*Skill) ([]*Skill, error)
	UpdateMemberSkills(ctx context.Context, id uint, skills []*Skill) ([]*Skill, error)
	DeleteMemberSkills(ctx context.Context, id uint) (bool, error)
}

type memberSkillQueueRepository interface {
	MemberSkillQueue(ctx context.Context, id uint) ([]*MemberSkillQueue, error)
	MemberSkillQueuePosition(ctx context.Context, id uint, position uint) (*MemberSkillQueue, error)
	CreateMemberSkillQueue(ctx context.Context, id uint, positions []*MemberSkillQueue) ([]*MemberSkillQueue, error)
	UpdateMemberSkillQueue(ctx context.Context, id uint, position *MemberSkillQueue) (*MemberSkillQueue, error)
	DeleteMemberSkillQueue(ctx context.Context, id uint) (bool, error)
	DeleteMemberSkillQueuePosition(ctx context.Context, id, position uint) (bool, error)
}

type MemberAttributes struct {
	MemberID                 uint      `db:"member_id" json:"member_id"`
	Charisma                 uint      `db:"charisma" json:"charisma"`
	Intelligence             uint      `db:"intelligence" json:"intelligence"`
	Memory                   uint      `db:"memory" json:"memory"`
	Perception               uint      `db:"perception" json:"perception"`
	Willpower                uint      `db:"willpower" json:"willpower"`
	BonusRemaps              null.Int  `db:"bonus_remaps,omitempty" json:"bonus_remaps,omitempty"`
	LastRemapDate            null.Time `db:"last_remap_date,omitempty" json:"last_remap_date,omitempty"`
	AccruedRemapCooldownDate null.Time `db:"accrued_remap_cooldown_date,omitempty" json:"accrued_remap_cooldown_date,omitempty"`
	CreatedAt                time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt                time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}

type MemberSkillQueue struct {
	MemberID        uint      `db:"member_id" json:"member_id" deep:"-"`
	QueuePosition   uint      `db:"queue_position" json:"queue_position"`
	SkillID         uint      `db:"skill_id" json:"skill_id"`
	FinishedLevel   uint      `db:"finished_level" json:"finished_level"`
	TrainingStartSp null.Int  `db:"training_start_sp,omitempty" json:"training_start_sp,omitempty"`
	LevelStartSp    null.Int  `db:"level_start_sp,omitempty" json:"level_start_sp,omitempty"`
	LevelEndSp      null.Int  `db:"level_end_sp,omitempty" json:"level_end_sp,omitempty"`
	StartDate       null.Time `db:"start_date,omitempty" json:"start_date,omitempty"`
	FinishDate      null.Time `db:"finish_date,omitempty" json:"finish_date,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}

type MemberSkills struct {
	MemberID      uint      `db:"member_id" json:"member_id" deep:"-"`
	TotalSP       uint      `db:"total_sp" json:"total_sp"`
	UnallocatedSP null.Int  `db:"unallocated_sp,omitempty" json:"unallocated_sp,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at" deep:"-"`

	Skills []*Skill `json:"skills"`
}

type Skill struct {
	ActiveSkillLevel   uint      `db:"active_skill_level" json:"active_skill_level"`
	SkillID            uint      `db:"skill_id" json:"skill_id"`
	SkillpointsInSkill uint      `db:"skillpoints_in_skill" json:"skillpoints_in_skill"`
	TrainedSkillLevel  uint      `db:"trained_skill_level" json:"trained_skill_level"`
	CreatedAt          time.Time `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at" deep:"-"`
}
