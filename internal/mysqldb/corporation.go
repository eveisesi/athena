package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type corporationRepository struct {
	db                   *sqlx.DB
	corporation, history string
}

func NewCorporationRepository(db *sql.DB) athena.CorporationRepository {
	return &corporationRepository{
		db:          sqlx.NewDb(db, "mysql"),
		corporation: "corporations",
		history:     "corporation_alliance_history",
	}
}

func (r *corporationRepository) Corporation(ctx context.Context, id uint) (*athena.Corporation, error) {

	corporations, err := r.Corporations(ctx, athena.NewEqualOperator("id", id), athena.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(corporations) == 0 {
		return nil, nil
	}

	return corporations[0], nil

}

func (r *corporationRepository) Corporations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Corporation, error) {

	query, args, err := BuildFilters(squirrel.Select(
		"id", "alliance_id", "ceo_id", "creator_id",
		"date_founded", "faction_id", "home_station_id", "member_count",
		"name", "shares", "tax_rate", "ticker",
		"url", "war_eligible", "created_at", "updated_at",
	).From(r.corporation), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to generate query: %w", err)
	}

	var corporations = make([]*athena.Corporation, 0)
	err = r.db.SelectContext(ctx, &corporations, query, args...)

	return corporations, err

}

func (r *corporationRepository) CreateCorporation(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, error) {

	i := squirrel.Insert("corporations").Columns(
		"id", "alliance_id", "ceo_id", "creator_id",
		"date_founded", "faction_id", "home_station_id", "member_count",
		"name", "shares", "tax_rate", "ticker",
		"url", "war_eligible", "created_at", "updated_at",
	).Values(
		corporation.ID, corporation.AllianceID, corporation.CeoID, corporation.CreatorID,
		corporation.DateFounded, corporation.FactionID, corporation.HomeStationID, corporation.MemberCount,
		corporation.Name, corporation.Shares, corporation.TaxRate, corporation.Ticker,
		corporation.URL, corporation.WarEligible, squirrel.Expr(`NOW()`), squirrel.Expr(`NOW()`),
	)

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to build SQL Query for Insert Statement: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to insert records: %w", err)
	}

	return r.Corporation(ctx, corporation.ID)

}

func (r *corporationRepository) UpdateCorporation(ctx context.Context, id uint, corporation *athena.Corporation) (*athena.Corporation, error) {

	u := squirrel.Update(r.corporation).
		Set("alliance_id", corporation.AllianceID).
		Set("ceo_id", corporation.CeoID).
		Set("faction_id", corporation.FactionID).
		Set("home_station_id", corporation.HomeStationID).
		Set("member_count", corporation.MemberCount).
		Set("shares", corporation.Shares).
		Set("tax_rate", corporation.TaxRate).
		Set("url", corporation.URL).
		Set("war_eligible", corporation.WarEligible).
		Set("updated_at", squirrel.Expr(`NOW()`)).
		Where(squirrel.Eq{"id": id})

	query, args, err := u.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to insert records: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to update records: %w", err)
	}

	return r.Corporation(ctx, corporation.ID)

}
func (r *corporationRepository) CorporationAllianceHistory(ctx context.Context, operators ...*athena.Operator) ([]*athena.CorporationAllianceHistory, error) {

	query, args, err := BuildFilters(sq.Select(
		"corporation_id",
		"alliance_id", "is_deleted",
		"record_id", "start_date",
		"created_at", "updated_at",
	).From(r.history), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to generate select query: %w", err)
	}

	var histories = make([]*athena.CorporationAllianceHistory, 0)
	err = r.db.SelectContext(ctx, &histories, query, args...)

	return histories, err

}

func (r *corporationRepository) CreateCorporationAllianceHistory(ctx context.Context, id uint, records []*athena.CorporationAllianceHistory) ([]*athena.CorporationAllianceHistory, error) {

	i := sq.Insert(r.history).Options("IGNORE").Columns(
		"corporation_id",
		"alliance_id", "is_deleted",
		"record_id", "start_date",
		"created_at", "updated_at",
	)
	for _, record := range records {
		i.Values(
			id,
			record.AllianceID, record.IsDeleteed,
			record.RecordID, record.StartDate,
			record.StartDate, sq.Expr(`NOW()`, sq.Expr(`NOW()`)),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to generate insert query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Corporation Repository] Failed to insert records: %w", err)
	}

	return r.CorporationAllianceHistory(ctx, athena.NewEqualOperator("corporation_id", id))

}
