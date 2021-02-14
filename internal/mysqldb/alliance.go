package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type allianceRepository struct {
	db *sqlx.DB
}

func NewAllianceRepository(db *sql.DB) athena.AllianceRepository {
	return &allianceRepository{
		db: sqlx.NewDb(db, "mysql"),
	}
}

func (r *allianceRepository) Alliance(ctx context.Context, id uint) (*athena.Alliance, error) {

	alliances, err := r.Alliances(ctx, athena.NewEqualOperator("id", id), athena.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(alliances) == 0 {
		return nil, nil
	}

	return alliances[0], nil

}

func (r *allianceRepository) Alliances(ctx context.Context, operators ...*athena.Operator) ([]*athena.Alliance, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "ticker", "date_founded", "creator_id",
		"creator_corporation_id", "executor_corporation_id",
		"is_closed", "created_at", "updated_at",
	), operators...).From("alliances").ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Alliance Repository] Failed to generate query: %w", err)
	}

	var alliances = make([]*athena.Alliance, 0)
	err = r.db.SelectContext(ctx, &alliances, query, args...)

	return alliances, err

}

func (r *allianceRepository) CreateAlliance(ctx context.Context, alliance *athena.Alliance) (*athena.Alliance, error) {

	i := sq.Insert("alliances").
		Columns(
			"id", "name", "ticker", "date_founded",
			"creator_id", "creator_corporation_id", "executor_corporation_id",
			"is_closed", "created_at", "updated_at").
		Values(
			alliance.ID, alliance.Name, alliance.Ticker, alliance.DateFounded,
			alliance.CreatorID, alliance.CreatorCorporationID, alliance.ExecutorCorporationID,
			alliance.IsClosed, sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Alliance Repository] Failed to build SQL Query for Insert Statement: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Alliance Repository] Failed to insert records: %w", err)
	}

	return r.Alliance(ctx, alliance.ID)

}

func (r *allianceRepository) UpdateAlliance(ctx context.Context, id uint, alliance *athena.Alliance) (*athena.Alliance, error) {

	u := sq.Update("alliances").
		Set("executor_corporation_id", alliance.ExecutorCorporationID).
		Set("is_closed", alliance.IsClosed).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id})

	query, args, err := u.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Alliance Repository] Failed to insert records: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Alliance Repository] Failed to update records: %w", err)
	}

	return alliance, nil

}
