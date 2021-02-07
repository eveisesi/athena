package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type cloneRepository struct {
	db *sqlx.DB
	clone,
	cloneMeta,
	clones,
	implants string
}

func NewCloneRepository(db *sql.DB) athena.CloneRepository {
	return &cloneRepository{
		db:        sqlx.NewDb(db, "mysql"),
		cloneMeta: "member_clone_meta",
		clone:     "member_home_clone",
		clones:    "member_jump_clones",
		implants:  "member_implants",
	}
}

func (r *cloneRepository) MemberClones(ctx context.Context, id uint) (*athena.MemberClones, error) {

	clones, err := r.memberCloneMeta(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to fetch clone meta: %w", err)
	}

	clones.HomeLocation, err = r.memberHomeClone(ctx, clones.MemberID)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to fetch home clone: %w", err)
	}

	clones.JumpClones, err = r.memberJumpClones(ctx, clones.MemberID)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to fetch jump clones: %w", err)
	}

	return clones, err

}

func (r *cloneRepository) CreateMemberClones(ctx context.Context, clones *athena.MemberClones) (*athena.MemberClones, error) {

	var err error
	_, err = r.createMemberCloneMeta(ctx, clones.MemberID, clones)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to create member clone meta record: %w", err)
	}

	_, err = r.createMemberHomeClone(ctx, clones.MemberID, clones.HomeLocation)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to create member home clone record: %w", err)
	}

	_, err = r.createMemberJumpClones(ctx, clones.MemberID, clones.JumpClones)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to create member jump clones record: %w", err)
	}

	return r.MemberClones(ctx, clones.MemberID)

}

func (r *cloneRepository) UpdateMemberClones(ctx context.Context, clones *athena.MemberClones) (*athena.MemberClones, error) {

	var err error
	_, err = r.updateMemberCloneMeta(ctx, clones.MemberID, clones)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to update member clone meta record: %w", err)
	}

	_, err = r.updateMemberHomeClone(ctx, clones.MemberID, clones.HomeLocation)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to update member home clone record: %w", err)
	}

	_, err = r.updateMemberJumpClones(ctx, clones.MemberID, clones.JumpClones)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to update member jump clones record: %w", err)
	}

	return r.MemberClones(ctx, clones.MemberID)

}

func (r *cloneRepository) DeleteMemberClones(ctx context.Context, id uint) (bool, error) {

	var err error
	_, err = r.deleteMemberJumpClones(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[Clone Repository] Failed to delete member jump clones record: %w", err)
	}

	_, err = r.deleteMemberHomeClone(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[Clone Repository] Failed to delete member home clone record: %w", err)
	}

	_, err = r.deleteMemberCloneMeta(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[Clone Repository] Failed to delete member clone meta record: %w", err)
	}

	return err == nil, err

}

func (r *cloneRepository) memberCloneMeta(ctx context.Context, id uint) (*athena.MemberClones, error) {

	query, args, err := sq.Select(
		"member_id", "last_clone_jump_date", "last_station_change_date",
		"created_at", "updated_at",
	).From(r.cloneMeta).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	var clones = new(athena.MemberClones)
	err = r.db.GetContext(ctx, clones, query, args...)

	return clones, err

}

func (r *cloneRepository) createMemberCloneMeta(ctx context.Context, id uint, meta *athena.MemberClones) (*athena.MemberClones, error) {

	i := sq.Insert(r.cloneMeta).
		Columns(
			"member_id", "last_clone_jump_date", "last_station_change_date",
			"created_at", "updated_at",
		).Values(id, meta.LastCloneJumpDate, meta.LastStationChangeDate,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	)

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.memberCloneMeta(ctx, id)

}

func (r *cloneRepository) updateMemberCloneMeta(ctx context.Context, id uint, meta *athena.MemberClones) (*athena.MemberClones, error) {

	u := sq.Update(r.cloneMeta).
		Set("last_clone_jump_date", meta.LastCloneJumpDate).
		Set("last_station_change_date", meta.LastStationChangeDate).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": id})

	query, args, err := u.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to update records: %w", err)
	}

	return r.memberCloneMeta(ctx, id)

}

func (r *cloneRepository) deleteMemberCloneMeta(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.cloneMeta).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to update records: %w", err)
	}

	return err == nil, err

}

func (r *cloneRepository) memberHomeClone(ctx context.Context, id uint) (*athena.MemberHomeLocation, error) {

	query, args, err := sq.Select(
		"location_id", "location_type", "created_at", "updated_at",
	).From(r.clone).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	var location = new(athena.MemberHomeLocation)
	err = r.db.GetContext(ctx, location, query, args...)

	return location, err

}

func (r *cloneRepository) createMemberHomeClone(ctx context.Context, id uint, clone *athena.MemberHomeLocation) (*athena.MemberHomeLocation, error) {

	i := sq.Insert(r.clone).
		Columns(
			"member_id",
			"location_id", "location_type",
			"created_at", "updated_at",
		).
		Values(
			id,
			clone.LocationID, clone.LocationType,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.memberHomeClone(ctx, id)

}

func (r *cloneRepository) updateMemberHomeClone(ctx context.Context, id uint, clone *athena.MemberHomeLocation) (*athena.MemberHomeLocation, error) {

	u := sq.Update(r.clones).
		Set("location_id", clone.LocationID).
		Set("location_type", clone.LocationType).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": id})

	query, args, err := u.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to update records: %w", err)
	}

	return r.memberHomeClone(ctx, id)

}

func (r *cloneRepository) deleteMemberHomeClone(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.clone).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to update records: %w", err)
	}

	return err == nil, err

}

func (r *cloneRepository) memberJumpClones(ctx context.Context, id uint) ([]*athena.MemberJumpClone, error) {

	query, args, err := sq.Select(
		"jump_clone_id", "location_id", "location_type", "implants", "created_at", "updated_at",
	).From(r.clones).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	var clones = make([]*athena.MemberJumpClone, 0)
	err = r.db.SelectContext(ctx, &clones, query, args...)

	return clones, err

}

func (r *cloneRepository) createMemberJumpClones(ctx context.Context, id uint, clones []*athena.MemberJumpClone) ([]*athena.MemberJumpClone, error) {

	i := sq.Insert(r.clones).
		Columns(
			"member_id", "jump_clone_id", "location_id",
			"location_type", "implants", "created_at", "updated_at",
		)
	for _, clone := range clones {
		i = i.Values(id, clone.JumpCloneID, clone.LocationID, clone.LocationType,
			clone.Implants, sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.memberJumpClones(ctx, id)

}

func (r *cloneRepository) updateMemberJumpClones(ctx context.Context, id uint, clones []*athena.MemberJumpClone) ([]*athena.MemberJumpClone, error) {

	for _, clone := range clones {
		u := sq.Update(r.clones).
			Set("location_id", clone.LocationID).
			Set("location_type", clone.LocationType).
			Set("implants", clone.Implants).
			Set("updated_at", sq.Expr(`NOW()`)).
			Where(sq.Eq{"member_id": id, "jump_clone_id": clone.JumpCloneID})

		query, args, err := u.ToSql()
		if err != nil {
			return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("[Clones Repository] Failed to update records: %w", err)
		}
	}

	return r.memberJumpClones(ctx, id)

}

func (r *cloneRepository) deleteMemberJumpClones(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.clones).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to update records: %w", err)
	}

	return err == nil, err

}

func (r *cloneRepository) MemberImplants(ctx context.Context, id uint) ([]*athena.MemberImplant, error) {

	query, args, err := sq.Select(
		"member_id", "implant_id", "created_at",
	).From(r.implants).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	var implants = make([]*athena.MemberImplant, 0)
	err = r.db.SelectContext(ctx, &implants, query, args...)

	return implants, err

}

func (r *cloneRepository) CreateMemberImplants(ctx context.Context, id uint, implants []*athena.MemberImplant) ([]*athena.MemberImplant, error) {

	i := sq.Insert(r.implants).
		Columns(
			"member_id", "implant_id", "created_at",
		)
	for _, implant := range implants {
		i = i.Values(implant.MemberID, implant.ImplantID, sq.Expr(`NOW()`), sq.Expr(`NOW()`))
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberImplants(ctx, id)

}

func (r *cloneRepository) DeleteMemberImplants(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.implants).Where(sq.Eq{"member_id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}
