package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type CharacterRepository interface {
	characterRepository
	characterHistoryRepository
}

type characterRepository interface {
	Character(ctx context.Context, id uint) (*Character, error)
	Characters(ctx context.Context, operators ...*Operator) ([]*Character, error)
	CreateCharacter(ctx context.Context, character *Character) (*Character, error)
	UpdateCharacter(ctx context.Context, id uint, character *Character) (*Character, error)
}

type characterHistoryRepository interface {
	CharacterCorporationHistory(ctx context.Context, operators ...*Operator) ([]*CharacterCorporationHistory, error)
	CreateCharacterCorporationHistory(ctx context.Context, id uint, records []*CharacterCorporationHistory) ([]*CharacterCorporationHistory, error)
}

type Character struct {
	ID             uint         `db:"id" json:"id"`
	Name           string       `db:"name" json:"name"`
	CorporationID  uint         `db:"corporation_id" json:"corporation_id"`
	AllianceID     null.Uint    `db:"alliance_id,omitempty" json:"alliance_id,omitempty"`
	FactionID      null.Uint    `db:"faction_id,omitempty" json:"faction_id,omitempty"`
	SecurityStatus null.Float64 `db:"security_status,omitempty" json:"security_status,omitempty"`
	Gender         string       `db:"gender" json:"gender"`
	Birthday       time.Time    `db:"birthday" json:"birthday"`
	Title          null.String  `db:"title,omitempty" json:"title,omitempty"`
	AncestryID     null.Uint    `db:"ancestry_id,omitempty" json:"ancestry_id,omitempty"`
	BloodlineID    uint         `db:"bloodline_id" json:"bloodline_id"`
	RaceID         uint         `db:"race_id" json:"race_id"`
	CreatedAt      time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at" json:"updated_at"`
}

type CharacterCorporationHistory struct {
	CharacterID   uint      `db:"character_id" json:"character_id"`
	RecordID      uint64    `db:"record_id" json:"record_id"`
	CorporationID uint      `db:"corporation_id" json:"corporation_id"`
	IsDeleted     bool      `db:"is_deleted" json:"is_deleted"`
	StartDate     time.Time `db:"start_date" json:"start_date"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
