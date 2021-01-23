package athena

import (
	"context"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberSkillRepository interface {
	memberAttributesRepository
	memberSkillQueueRepository
	memberSkillRepository
}

type memberAttributesRepository interface {
	MemberSkillAttributes(ctx context.Context, id string) (*MemberSkillAttributes, error)
	CreateMemberSkillAttributes(ctx context.Context, location *MemberSkillAttributes) (*MemberSkillAttributes, error)
	UpdateMemberSkillAttributes(ctx context.Context, id string, location *MemberSkillAttributes) (*MemberSkillAttributes, error)
	DeleteMemberSkillAttributes(ctx context.Context, id string) (bool, error)
}

type memberSkillQueueRepository interface {
	MemberSkillQueue(ctx context.Context, memberID string) (*MemberSkillQueue, error)
	CreateMemberSkillQueue(ctx context.Context, skillQueue *MemberSkillQueue) (*MemberSkillQueue, error)
	UpdateMemberSkillQueue(ctx context.Context, memberID string, skillQueue *MemberSkillQueue) (*MemberSkillQueue, error)
	DeleteMemberSkillQueue(ctx context.Context, memberID string) (bool, error)
}

type memberSkillRepository interface {
	MemberSkills(ctx context.Context, memberID string) (*MemberSkill, error)
	CreateMemberSkills(ctx context.Context, skillQueue *MemberSkill) (*MemberSkill, error)
	UpdateMemberSkills(ctx context.Context, memberID string, skillQueue *MemberSkill) (*MemberSkill, error)
	DeleteMemberSkills(ctx context.Context, memberID string) (bool, error)
}

type MemberSkillAttributes struct {
	MemberID                 primitive.ObjectID `bson:"member_id" json:"member_id"`
	Charisma                 int                `bson:"charisma" json:"charisma"`
	Intelligence             int                `bson:"intelligence" json:"intelligence"`
	Memory                   int                `bson:"memory" json:"memory"`
	Perception               int                `bson:"perception" json:"perception"`
	Willpower                int                `bson:"willpower" json:"willpower"`
	BonusRemaps              null.Int           `bson:"bonus_remaps,omitempty" json:"bonus_remaps,omitempty"`
	LastRemapDate            null.Time          `bson:"last_remap_date,omitempty" json:"last_remap_date,omitempty"`
	AccruedRemapCooldownDate null.Time          `bson:"accrued_remap_cooldown_date,omitempty" json:"accrued_remap_cooldown_date,omitempty"`
	Meta
}

type MemberSkillQueue struct {
	MemberID   primitive.ObjectID `bson:"member_id" json:"member_id"`
	SkillQueue []SkillQueueItem   `bson:"skill_queue"`
	Meta
}

type SkillQueueItem struct {
	SkillID         int       `bson:"skill_id" json:"skill_id"`
	QueuePosition   int       `bson:"queue_position" json:"queue_position"`
	FinishedLevel   int       `bson:"finished_level" json:"finished_level"`
	TrainingStartSp null.Int  `bson:"training_start_sp,omitempty" json:"training_start_sp,omitempty"`
	LevelStartSp    null.Int  `bson:"level_start_sp,omitempty" json:"level_start_sp,omitempty"`
	LevelEndSp      null.Int  `bson:"level_end_sp,omitempty" json:"level_end_sp,omitempty"`
	StartDate       null.Time `bson:"start_date,omitempty" json:"start_date,omitempty"`
	FinishDate      null.Time `bson:"finish_date,omitempty" json:"finish_date,omitempty"`
}

type MemberSkill struct {
	MemberID      primitive.ObjectID `bson:"member_id" json:"member_id"`
	Skills        []SkillItem        `bson:"skills" json:"skills"`
	TotalSP       int64              `bson:"total_sp" json:"total_sp"`
	UnallocatedSP null.Int           `bson:"unallocated_sp,omitempty" json:"unallocated_sp,omitempty"`
	Meta
}

type SkillItem struct {
	ActiveSkillLevel   int `bson:"active_skill_level" json:"active_skill_level"`
	SkillID            int `bson:"skill_id" json:"skill_id"`
	SkillpointsInSkill int `bson:"skillpoints_in_skill" json:"skillpoints_in_skill"`
	TrainedSkillLevel  int `bson:"trained_skill_level" json:"trained_skill_level"`
}
