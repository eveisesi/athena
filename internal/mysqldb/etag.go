package mysqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type etagRepository struct {
	db    *sqlx.DB
	table string
}

func NewEtagRepository(db *sql.DB) athena.EtagRepository {
	return &etagRepository{
		db:    sqlx.NewDb(db, "driver"),
		table: "etags",
	}
}

func (r *etagRepository) Etag(ctx context.Context, etagID string) (*athena.Etag, error) {

	etags, err := r.Etags(ctx, athena.NewEqualOperator("etag_id", etagID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(etags) != 1 {
		return nil, nil
	}

	return etags[0], nil

}

func (r *etagRepository) Etags(ctx context.Context, operators ...*athena.Operator) ([]*athena.Etag, error) {

	query, args, err := BuildFilters(sq.Select(
		"etag_id", "etag", "cached_until", "created_at", "updated_at",
	).From(r.table), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Etag Repository] Failed to generate sql: %w", err)
	}

	var etags = make([]*athena.Etag, 0)
	err = r.db.SelectContext(ctx, &etags, query, args...)

	return etags, err

}

func (r *etagRepository) InsertEtag(ctx context.Context, etag *athena.Etag) (*athena.Etag, error) {

	query := `
		INSERT INTO etags (
			etag_id,etag,cached_until,created_at,updated_at
		) VALUES (?, ?, ?, NOW(),NOW()) AS alias 
		ON DUPLICATE KEY UPDATE 
			etag=alias.etag,
			cached_until=alias.cached_until,
			updated_at=alias.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, etag.EtagID, etag.Etag, etag.CachedUntil)
	if err != nil {
		return nil, fmt.Errorf("[Etag Repository] Failed to insert record: %w", err)
	}

	return r.Etag(ctx, etag.EtagID)

}

func (r *etagRepository) UpdateEtag(ctx context.Context, id string, etag *athena.Etag) (*athena.Etag, error) {

	query, args, err := sq.Update(r.table).
		Set("etag", etag.Etag).
		Set("cached_until", etag.CachedUntil).
		Where(sq.Eq{"etag_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Etag Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Etag Repository] Failed to update record: %w", err)
	}

	return r.Etag(ctx, id)

}

func (r *etagRepository) DeleteEtag(ctx context.Context, id string) (bool, error) {

	query, args, err := sq.Delete(r.table).Where(sq.Eq{"etag_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Etag Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}
