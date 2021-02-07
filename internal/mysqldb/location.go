package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type memberLocationRepository struct {
	db *sqlx.DB
	location,
	online,
	ship string
}

func NewMemberLocationRepository(db *sql.DB) athena.MemberLocationRepository {
	return &memberLocationRepository{
		db:       sqlx.NewDb(db, "mysql"),
		location: "member_location",
		online:   "member_online",
		ship:     "member_ship",
	}
}

func (r *memberLocationRepository) MemberLocation(ctx context.Context, memberID uint) (*athena.MemberLocation, error) {

	query, args, err := sq.Select(
		"member_id",
		"solar_system_id", "station_id", "structure_id",
		"created_at", "updated_at",
	).From(r.location).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var location = new(athena.MemberLocation)
	err = r.db.GetContext(ctx, location, query, args...)

	return location, err

}

func (r *memberLocationRepository) CreateMemberLocation(ctx context.Context, memberID uint, location *athena.MemberLocation) (*athena.MemberLocation, error) {

	query, args, err := sq.Insert(r.location).Columns(
		"member_id",
		"solar_system_id", "station_id", "structure_id",
		"created_at", "updated_at",
	).Values(
		memberID,
		location.SolarSystemID, location.StationID, location.StructureID,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberLocation(ctx, memberID)

}

func (r *memberLocationRepository) UpdateMemberLocation(ctx context.Context, memberID uint, location *athena.MemberLocation) (*athena.MemberLocation, error) {

	query, args, err := sq.Update(r.location).
		Set("solar_system_id", location.SolarSystemID).
		Set("station_id", location.StationID).
		Set("structure_id", location.StructureID).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberLocation(ctx, memberID)

}

func (r *memberLocationRepository) MemberOnline(ctx context.Context, memberID uint) (*athena.MemberOnline, error) {

	query, args, err := sq.Select(
		"member_id",
		"last_login", "last_logout",
		"logins", "online",
		"created_at",
		"updated_at",
	).From(r.online).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var online = new(athena.MemberOnline)
	err = r.db.GetContext(ctx, online, query, args...)

	return online, err

}

func (r *memberLocationRepository) CreateMemberOnline(ctx context.Context, memberID uint, online *athena.MemberOnline) (*athena.MemberOnline, error) {

	query, args, err := sq.Insert(r.online).Columns(
		"member_id",
		"last_login", "last_logout",
		"logins", "online",
		"created_at",
		"updated_at",
	).Values(
		memberID,
		online.LastLogin, online.LastLogout, online.Logins, online.Online,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberOnline(ctx, memberID)

}

func (r *memberLocationRepository) UpdateMemberOnline(ctx context.Context, memberID uint, online *athena.MemberOnline) (*athena.MemberOnline, error) {

	query, args, err := sq.Update(r.online).
		Set("last_login", online.LastLogin).
		Set("last_logout", online.LastLogout).
		Set("logins", online.Logins).
		Set("online", online.Online).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberOnline(ctx, memberID)

}

func (r *memberLocationRepository) MemberShip(ctx context.Context, memberID uint) (*athena.MemberShip, error) {

	query, args, err := sq.Select(
		"member_id",
		"ship_item_id", "ship_name", "ship_type_id",
		"created_at", "updated_at",
	).From(r.ship).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var ship = new(athena.MemberShip)
	err = r.db.GetContext(ctx, ship, query, args...)

	return ship, err

}

func (r *memberLocationRepository) CreateMemberShip(ctx context.Context, memberID uint, ship *athena.MemberShip) (*athena.MemberShip, error) {

	query, args, err := sq.Insert(r.ship).Columns(
		"member_id",
		"ship_item_id", "ship_name", "ship_type_id",
		"created_at", "updated_at",
	).Values(
		memberID,
		ship.ShipItemID, ship.ShipName, ship.ShipTypeID,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberShip(ctx, memberID)

}

func (r *memberLocationRepository) UpdateMemberShip(ctx context.Context, memberID uint, ship *athena.MemberShip) (*athena.MemberShip, error) {

	query, args, err := sq.Update(r.ship).
		Set("ship_item_id", ship.ShipItemID).
		Set("ship_name", ship.ShipName).
		Set("ship_type_id", ship.ShipTypeID).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberShip(ctx, memberID)

}
