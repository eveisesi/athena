package mysqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type memberRepository struct {
	db    *sqlx.DB
	table string
}

func NewMemberRepository(db *sql.DB) athena.MemberRepository {
	return &memberRepository{
		db:    sqlx.NewDb(db, "mysql"),
		table: "members",
	}
}

func (r *memberRepository) Member(ctx context.Context, id uint) (*athena.Member, error) {

	members, err := r.Members(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(members) != 1 {
		return nil, nil
	}

	return members[0], nil

}

func (r *memberRepository) Members(ctx context.Context, operators ...*athena.Operator) ([]*athena.Member, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "main_id", "access_token", "refresh_token", "expires",
		"owner_hash", "scopes", "disabled", "disabled_reason", "disabled_timestamp",
		"last_login", "created_at", "updated_at",
	).From(r.table), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Member Repository] Failed to generate sql: %w", err)
	}

	var members = make([]*athena.Member, 0)
	err = r.db.SelectContext(ctx, &members, query, args...)

	return members, err

}

func (r *memberRepository) CreateMember(ctx context.Context, member *athena.Member) (*athena.Member, error) {

	query, args, err := squirrel.Insert(r.table).Columns(
		"id", "main_id", "access_token", "refresh_token", "expires",
		"owner_hash", "scopes", "disabled", "disabled_reason", "disabled_timestamp",
		"last_login", "created_at", "updated_at",
	).Values(
		member.ID,
		member.MainID,
		member.AccessToken,
		member.RefreshToken,
		member.Expires,
		member.OwnerHash,
		member.Scopes,
		member.Disabled,
		member.DisabledReason,
		member.DisabledTimestamp,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Member Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Member Repository] Failed to insert record: %w", err)
	}

	return r.Member(ctx, member.ID)

}

func (r *memberRepository) UpdateMember(ctx context.Context, id uint, member *athena.Member) (*athena.Member, error) {

	query, args, err := sq.Update(r.table).
		Set("id", member.ID).
		Set("main_id", member.MainID).
		Set("access_token", member.AccessToken).
		Set("refresh_token", member.RefreshToken).
		Set("expires", member.Expires).
		Set("owner_hash", member.OwnerHash).
		Set("scopes", member.Scopes).
		Set("disabled", member.Disabled).
		Set("disabled_reason", member.DisabledReason).
		Set("disabled_timestamp", member.DisabledTimestamp).
		Set("last_login", member.LastLogin).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Member Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Member Repository] Failed to update record: %w", err)
	}

	return r.Member(ctx, id)

}

func (r *memberRepository) DeleteMember(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.table).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Member Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err == nil, err

}
