package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type memberFittingRepository struct {
	db              *sqlx.DB
	fittings, items string
}

func NewFittingRepository(db *sql.DB) athena.MemberFittingsRepository {
	return &memberFittingRepository{
		db:       sqlx.NewDb(db, "mysql"),
		fittings: "member_fittings",
		items:    "member_fitting_items",
	}
}

func (r *memberFittingRepository) MemberFitting(ctx context.Context, memberID, fittingID uint) (*athena.MemberFitting, error) {

	query, args, err := sq.Select(
		"member_id", "fitting_id", "ship_type_id",
		"name", "description", "created_at", "updated_at",
	).From(r.fittings).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	var fitting = new(athena.MemberFitting)
	err = r.db.GetContext(ctx, fitting, query, args...)

	return fitting, err

}

func (r *memberFittingRepository) MemberFittings(ctx context.Context, memberID uint, operators ...*athena.Operator) ([]*athena.MemberFitting, error) {

	query, args, err := sq.Select(
		"member_id", "fitting_id", "ship_type_id",
		"name", "description", "created_at", "updated_at",
	).From(r.fittings).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	var fittings = make([]*athena.MemberFitting, 0)
	err = r.db.SelectContext(ctx, &fittings, query, args...)

	return fittings, err

}

func (r *memberFittingRepository) CreateMemberFittings(ctx context.Context, memberID uint, fittings []*athena.MemberFitting) ([]*athena.MemberFitting, error) {

	i := sq.Insert(r.fittings).Columns(
		"member_id", "fitting_id", "ship_type_id",
		"name", "description", "created_at", "updated_at",
	)
	for _, fitting := range fittings {
		i.Values(
			fitting.MemberID,
			fitting.FittingID, fitting.ShipTypeID,
			fitting.Name, fitting.Description,
			fitting.CreatedAt, fitting.UpdatedAt,
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to insert records: %w", err)
	}

	return r.MemberFittings(ctx, memberID)

}

func (r *memberFittingRepository) UpdateMemberFitting(ctx context.Context, memberID, fittingID uint, fitting *athena.MemberFitting) (*athena.MemberFitting, error) {

	query, args, err := sq.Update(r.fittings).
		Set("name", fitting.Name).
		Set("description", fitting.Description).
		Where(sq.Eq{"member_id": memberID, "fitting_id": fittingID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to insert records: %w", err)
	}

	return r.MemberFitting(ctx, memberID, fitting.FittingID)

}

func (r *memberFittingRepository) DeleteMemberFitting(ctx context.Context, memberID, fittingID uint) (bool, error) {

	query, args, err := sq.Delete(r.fittings).Where(sq.Eq{"member_id": memberID, "fitting_id": fittingID}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to insert records: %w", err)
	}

	return true, nil

}

func (r *memberFittingRepository) DeleteMemberFittings(ctx context.Context, memberID uint) (bool, error) {

	query, args, err := sq.Delete(r.fittings).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to insert records: %w", err)
	}

	return true, nil

}

func (r *memberFittingRepository) MemberFittingItems(ctx context.Context, memberID, fittingID uint) ([]*athena.MemberFittingItem, error) {

	query, args, err := sq.Select(
		"member_id",
		"fitting_id", "type_id", "quantity",
		"flag", "created_at", "updated_at",
	).From(r.items).Where(sq.Eq{"member_id": memberID, "fitting_id": fittingID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	var items = make([]*athena.MemberFittingItem, 0)
	err = r.db.SelectContext(ctx, &items, query, args...)

	return items, err

}

func (r *memberFittingRepository) CreateMemberFittingItems(ctx context.Context, memberID, fittingID uint, items []*athena.MemberFittingItem) ([]*athena.MemberFittingItem, error) {

	i := sq.Insert(r.fittings).Columns(
		"member_id",
		"fitting_id", "type_id", "quantity",
		"flag", "created_at", "updated_at",
	)
	for _, item := range items {
		i.Values(
			memberID,
			fittingID, item.TypeID,
			item.Quantity, item.Flag,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Fitting Repository] Failed to insert records: %w", err)
	}

	return r.MemberFittingItems(ctx, memberID, fittingID)

}

func (r *memberFittingRepository) DeleteMemberFittingItems(ctx context.Context, memberID, fittingID uint) (bool, error) {

	query, args, err := sq.Delete(r.items).Where(sq.Eq{"member_id": memberID, "fitting_id": fittingID}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to insert records: %w", err)
	}

	return true, nil

}

func (r *memberFittingRepository) DeleteMemberFittingItemsAll(ctx context.Context, memberID uint) (bool, error) {

	query, args, err := sq.Delete(r.items).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Fitting Repository] Failed to insert records: %w", err)
	}

	return true, nil

}
