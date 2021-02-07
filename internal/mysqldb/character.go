package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type characterRepository struct {
	db *sqlx.DB
}

func NewCharacterRepository(db *sql.DB) athena.CharacterRepository {
	return &characterRepository{
		db: sqlx.NewDb(db, "mysql"),
	}
}

func (r *characterRepository) Character(ctx context.Context, id uint) (*athena.Character, error) {

	characters, err := r.Characters(ctx, athena.NewEqualOperator("id", id), athena.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(characters) == 0 {
		return nil, nil
	}

	return characters[0], nil

}

func (r *characterRepository) Characters(ctx context.Context, operators ...*athena.Operator) ([]*athena.Character, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "corporation_id", "alliance_id", "faction_id",
		"security_status", "gender", "birthday", "title", "ancestry_id",
		"bloodline_id", "race_id", "created_at", "updated_at",
	), operators...).From("characters").ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Character Repository] Failed to generate query: %w", err)
	}

	var characters = make([]*athena.Character, 0)
	err = r.db.SelectContext(ctx, &characters, query, args...)

	return characters, err

}

func (r *characterRepository) CreateCharacter(ctx context.Context, character *athena.Character) (*athena.Character, error) {

	now := time.Now()
	character.CreatedAt = now
	character.UpdatedAt = now

	i := sq.Insert("characters").Columns(
		"id", "name", "corporation_id", "gender", "birthday", "bloodline_id", "race_id",
		"security_status", "title", "ancestry_id", "alliance_id", "faction_id", "created_at", "updated_at",
	).Values(
		character.ID, character.Name, character.CorporationID, character.Gender, character.Birthday,
		character.BloodlineID, character.RaceID, character.SecurityStatus, character.Title, character.AncestryID,
		character.AllianceID, character.FactionID, character.CreatedAt, character.UpdatedAt,
	)

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Character Repository] Failed to build SQL Query for Insert Statement: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Character Repository] Failed to insert records: %w", err)
	}

	return character, nil

}

func (r *characterRepository) UpdateCharacter(ctx context.Context, id uint, character *athena.Character) (*athena.Character, error) {

	query, args, err := sq.Update("characters").
		Set("corporation_id", character.CorporationID).
		Set("alliance_id", character.AllianceID).
		Set("faction_id", character.FactionID).
		Set("security_status", character.SecurityStatus).
		Set("title", character.Title).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to generate update query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to update records: %w", err)
	}

	return r.Character(ctx, character.ID)

}

func (r *characterRepository) CharacterCorporationHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CharacterCorporationHistory, error) {

	query, args, err := BuildFilters(sq.Select(
		"character_id", "record_id", "corporation_id", "is_deleted",
		"start_date", "created_at", "updated_at"), operators...).
		From("character_corporation_history").ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Character Repository] Failed to generate select query: %w", err)
	}

	var histories = make([]*athena.CharacterCorporationHistory, 0)
	err = r.db.SelectContext(ctx, &histories, query, args...)

	return histories, err

}

func (r *characterRepository) CreateCharacterCorporationHistory(ctx context.Context, id uint, records []*athena.CharacterCorporationHistory) ([]*athena.CharacterCorporationHistory, error) {

	i := sq.Insert("character_corporation_history").Columns(
		"character_id", "record_id", "corporation_id", "is_deleted", "start_date",
		"created_at", "updated_at",
	)
	for _, record := range records {
		i.Values(
			id, record.RecordID, record.CorporationID, record.IsDeleted,
			record.StartDate, sq.Expr(`NOW()`, sq.Expr(`NOW()`)),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to generate insert query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Character Repository] Failed to insert records: %w", err)
	}

	return r.CharacterCorporationHistory(ctx, athena.NewEqualOperator("character_id", id))

}

func (r *characterRepository) DeleteCharacterCorporationHistory(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete("character_corporation_history").Where(sq.Eq{"character_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Corporation Repository] Failed to generate delete query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("[Character Repository] Failed to insert records: %w", err)
	}

	return err == nil, err

}
