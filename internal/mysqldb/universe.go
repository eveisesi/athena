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

type universeRepository struct {
	db *sqlx.DB
	ancestries, bloodlines, categories, constellations,
	factions, groups, races, regions,
	solarSystems, stations, structures, types string
}

func NewUniverseRepository(db *sql.DB) athena.UniverseRepository {
	return &universeRepository{
		db:             sqlx.NewDb(db, "mysql"),
		ancestries:     "ancestries",
		bloodlines:     "bloodlines",
		categories:     "type_categories",
		constellations: "constellations",
		factions:       "factions",
		groups:         "type_groups",
		races:          "races",
		regions:        "regions",
		solarSystems:   "solar_systems",
		stations:       "stations",
		structures:     "structures",
		types:          "types",
	}
}

func (r *universeRepository) Ancestry(ctx context.Context, id uint) (*athena.Ancestry, error) {

	ancestries, err := r.Ancestries(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(ancestries) != 1 {
		return nil, nil
	}

	return ancestries[0], nil

}

func (r *universeRepository) Ancestries(ctx context.Context, operators ...*athena.Operator) ([]*athena.Ancestry, error) {

	query, args, err := BuildFilters(
		sq.Select(
			"id", "name", "bloodline_id",
			"created_at", "updated_at",
		).From(r.ancestries), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var ancestries = make([]*athena.Ancestry, 0)
	err = r.db.SelectContext(ctx, &ancestries, query, args...)

	return ancestries, err

}

func (r *universeRepository) CreateAncestry(ctx context.Context, ancestry *athena.Ancestry) (*athena.Ancestry, error) {

	query, args, err := sq.Insert(r.ancestries).Columns(
		"id", "name", "bloodline_id",
		"created_at", "updated_at",
	).Values(
		ancestry.ID, ancestry.Name, ancestry.BloodlineID,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to insert record: %w", err)
	}

	return r.Ancestry(ctx, ancestry.ID)

}

func (r *universeRepository) UpdateAncestry(ctx context.Context, id uint, ancestry *athena.Ancestry) (*athena.Ancestry, error) {

	query, args, err := sq.Update(r.ancestries).
		Set("name", ancestry.Name).
		Set("bloodline_id", ancestry.BloodlineID).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Ancestry(ctx, id)

}

func (r *universeRepository) DeleteAncestry(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.ancestries).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Bloodline(ctx context.Context, id uint) (*athena.Bloodline, error) {

	bloodlines, err := r.Bloodlines(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(bloodlines) != 1 {
		return nil, nil
	}

	return bloodlines[0], nil

}

func (r *universeRepository) Bloodlines(ctx context.Context, operators ...*athena.Operator) ([]*athena.Bloodline, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "race_id", "corporation_id",
		"ship_type_id", "charisma", "intelligence", "memory",
		"perception", "willpower", "created_at", "updated_at",
	).From(r.bloodlines), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var bloodlines = make([]*athena.Bloodline, 0)
	err = r.db.SelectContext(ctx, &bloodlines, query, args...)

	return bloodlines, err

}

func (r *universeRepository) CreateBloodline(ctx context.Context, bloodline *athena.Bloodline) (*athena.Bloodline, error) {

	query, args, err := sq.Insert(r.bloodlines).Columns(
		"id", "name", "race_id", "corporation_id",
		"ship_type_id", "charisma", "intelligence", "memory",
		"perception", "willpower", "created_at", "updated_at",
	).Values(
		bloodline.ID, bloodline.Name, bloodline.RaceID, bloodline.CorporationID,
		bloodline.ShipTypeID, bloodline.Charisma, bloodline.Intelligence,
		bloodline.Memory, bloodline.Perception, bloodline.Willpower,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to insert record: %w", err)
	}

	return r.Bloodline(ctx, bloodline.ID)

}

func (r *universeRepository) UpdateBloodline(ctx context.Context, id uint, bloodline *athena.Bloodline) (*athena.Bloodline, error) {

	query, args, err := sq.Update(r.bloodlines).
		Set("name", bloodline.Name).
		Set("race_id", bloodline.RaceID).
		Set("corporation_id", bloodline.CorporationID).
		Set("ship_type_id", bloodline.ShipTypeID).
		Set("charisma", bloodline.Charisma).
		Set("intelligence", bloodline.Intelligence).
		Set("memory", bloodline.Memory).
		Set("perception", bloodline.Perception).
		Set("willpower", bloodline.Willpower).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Bloodline(ctx, id)

}

func (r *universeRepository) DeleteBloodline(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.bloodlines).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Category(ctx context.Context, id uint) (*athena.Category, error) {

	categories, err := r.Categories(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(categories) != 1 {
		return nil, nil
	}

	return categories[0], nil

}

func (r *universeRepository) Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "published",
		"created_at", "updated_at",
	).From(r.categories), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var categories = make([]*athena.Category, 0)
	err = r.db.SelectContext(ctx, &categories, query, args...)

	return categories, err

}

func (r *universeRepository) CreateCategory(ctx context.Context, category *athena.Category) (*athena.Category, error) {

	query, args, err := sq.Insert(r.categories).Columns(
		"id", "name", "published",
		"created_at", "updated_at",
	).Values(
		category.ID, category.Name, category.Published,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to insert record: %w", err)
	}

	return r.Category(ctx, category.ID)

}

func (r *universeRepository) UpdateCategory(ctx context.Context, id uint, category *athena.Category) (*athena.Category, error) {

	query, args, err := sq.Update(r.categories).
		Set("name", category.Name).
		Set("published", category.Published).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Category(ctx, id)

}

func (r *universeRepository) DeleteCategory(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.categories).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Constellation(ctx context.Context, id uint) (*athena.Constellation, error) {

	constellations, err := r.Constellations(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(constellations) != 1 {
		return nil, nil
	}

	return constellations[0], nil

}

func (r *universeRepository) Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "region_id",
		"created_at", "updated_at",
	).From(r.constellations), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var constellations = make([]*athena.Constellation, 0)
	err = r.db.SelectContext(ctx, &constellations, query, args...)

	return constellations, err

}

func (r *universeRepository) CreateConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, error) {

	query, args, err := sq.Insert(r.constellations).Columns(
		"id", "name", "region_id",
		"created_at", "updated_at",
	).Values(
		constellation.ID, constellation.Name, constellation.RegionID,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to insert record: %w", err)
	}

	return r.Constellation(ctx, constellation.ID)

}

func (r *universeRepository) UpdateConstellation(ctx context.Context, id uint, constellation *athena.Constellation) (*athena.Constellation, error) {

	query, args, err := sq.Update(r.constellations).
		Set("name", constellation.Name).
		Set("region_id", constellation.RegionID).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Constellation(ctx, id)

}

func (r *universeRepository) DeleteConstellation(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.constellations).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Faction(ctx context.Context, id uint) (*athena.Faction, error) {

	factions, err := r.Factions(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(factions) != 1 {
		return nil, nil
	}

	return factions[0], nil

}

func (r *universeRepository) Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "is_unique", "size_factor",
		"station_count", "station_system_count", "corporation_id", "militia_corporation_id",
		"solar_system_id", "created_at", "updated_at",
	).From(r.factions), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var factions = make([]*athena.Faction, 0)
	err = r.db.SelectContext(ctx, &factions, query, args...)

	return factions, err

}

func (r *universeRepository) CreateFaction(ctx context.Context, faction *athena.Faction) (*athena.Faction, error) {

	query, args, err := sq.Insert(r.factions).Columns(
		"id", "name", "is_unique", "size_factor",
		"station_count", "station_system_count", "corporation_id", "militia_corporation_id",
		"solar_system_id", "created_at", "updated_at",
	).Values(
		faction.ID, faction.Name, faction.IsUnique, faction.SizeFactor,
		faction.StationCount, faction.StationSystemCount, faction.CorporationID,
		faction.MilitiaCorporationID, faction.SolarSystemID,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Faction(ctx, faction.ID)

}

func (r *universeRepository) UpdateFaction(ctx context.Context, id uint, faction *athena.Faction) (*athena.Faction, error) {

	query, args, err := sq.Update(r.factions).
		Set("name", faction.Name).
		Set("is_unique", faction.IsUnique).
		Set("size_factor", faction.SizeFactor).
		Set("station_count", faction.StationCount).
		Set("station_system_count", faction.StationSystemCount).
		Set("corporation_id", faction.CorporationID).
		Set("militia_corporation_id", faction.MilitiaCorporationID).
		Set("solar_system_id", faction.SolarSystemID).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Faction(ctx, id)

}

func (r *universeRepository) DeleteFaction(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.factions).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Group(ctx context.Context, id uint) (*athena.Group, error) {

	groups, err := r.Groups(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(groups) != 1 {
		return nil, nil
	}

	return groups[0], nil

}

func (r *universeRepository) Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "published",
		"category_id", "created_at", "updated_at",
	).From(r.groups), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var groups = make([]*athena.Group, 0)
	err = r.db.SelectContext(ctx, &groups, query, args...)

	return groups, err

}

func (r *universeRepository) CreateGroup(ctx context.Context, group *athena.Group) (*athena.Group, error) {

	query, args, err := sq.Insert(r.groups).Columns(
		"id", "name",
		"published", "category_id",
		"created_at", "updated_at",
	).Values(
		group.ID, group.Name,
		group.Published, group.CategoryID,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to insert record: %w", err)
	}

	return r.Group(ctx, group.ID)

}

func (r *universeRepository) UpdateGroup(ctx context.Context, id uint, group *athena.Group) (*athena.Group, error) {

	query, args, err := sq.Update(r.groups).
		Set("name", group.Name).
		Set("published", group.Published).
		Set("category_id", group.CategoryID).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Group(ctx, id)

}

func (r *universeRepository) DeleteGroup(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.groups).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Race(ctx context.Context, id uint) (*athena.Race, error) {

	races, err := r.Races(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(races) != 1 {
		return nil, nil
	}

	return races[0], nil

}

func (r *universeRepository) Races(ctx context.Context, operators ...*athena.Operator) ([]*athena.Race, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name",
		"created_at", "updated_at",
	).From(r.races), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var races = make([]*athena.Race, 0)
	err = r.db.SelectContext(ctx, &races, query, args...)

	return races, err

}

func (r *universeRepository) CreateRace(ctx context.Context, race *athena.Race) (*athena.Race, error) {

	query, args, err := sq.Insert(r.races).Columns(
		"id", "name",
		"created_at", "updated_at",
	).Values(
		race.ID, race.Name,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to insert record: %w", err)
	}

	return r.Race(ctx, race.ID)

}

func (r *universeRepository) UpdateRace(ctx context.Context, id uint, race *athena.Race) (*athena.Race, error) {

	query, args, err := sq.Update(r.races).
		Set("name", race.Name).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Race(ctx, id)

}

func (r *universeRepository) DeleteRace(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.races).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Region(ctx context.Context, id uint) (*athena.Region, error) {

	regions, err := r.Regions(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(regions) != 1 {
		return nil, nil
	}

	return regions[0], nil

}

func (r *universeRepository) Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name",
		"created_at", "updated_at",
	).From(r.regions), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var regions = make([]*athena.Region, 0)
	err = r.db.SelectContext(ctx, &regions, query, args...)

	return regions, err

}

func (r *universeRepository) CreateRegion(ctx context.Context, region *athena.Region) (*athena.Region, error) {

	query, args, err := sq.Insert(r.regions).Columns(
		"id", "name",
		"created_at", "updated_at",
	).Values(
		region.ID, region.Name,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Region(ctx, region.ID)

}

func (r *universeRepository) UpdateRegion(ctx context.Context, id uint, region *athena.Region) (*athena.Region, error) {

	query, args, err := sq.Update(r.regions).
		Set("name", region.Name).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Region(ctx, id)

}

func (r *universeRepository) DeleteRegion(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.regions).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) SolarSystem(ctx context.Context, id uint) (*athena.SolarSystem, error) {

	systems, err := r.SolarSystems(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(systems) != 1 {
		return nil, nil
	}

	return systems[0], nil

}

func (r *universeRepository) SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "constellation_id",
		"security_status", "star_id", "security_class",
		"created_at", "updated_at",
	).From(r.solarSystems), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var solarSystems = make([]*athena.SolarSystem, 0)
	err = r.db.SelectContext(ctx, &solarSystems, query, args...)

	return solarSystems, err

}

func (r *universeRepository) CreateSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, error) {

	query, args, err := sq.Insert(r.solarSystems).Columns(
		"id", "name", "constellation_id",
		"security_status", "star_id", "security_class",
		"created_at", "updated_at",
	).Values(
		solarSystem.ID, solarSystem.Name,
		solarSystem.ConstellationID, solarSystem.SecurityStatus,
		solarSystem.StarID, solarSystem.SecurityClass,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.SolarSystem(ctx, solarSystem.ID)

}

func (r *universeRepository) UpdateSolarSystem(ctx context.Context, id uint, solarSystem *athena.SolarSystem) (*athena.SolarSystem, error) {

	query, args, err := sq.Update(r.solarSystems).
		Set("name", solarSystem.Name).
		Set("constellation_id", solarSystem.ConstellationID).
		Set("security_status", solarSystem.SecurityStatus).
		Set("star_id", solarSystem.StarID).
		Set("security_class", solarSystem.SecurityClass).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.SolarSystem(ctx, id)

}

func (r *universeRepository) DeleteSolarSystem(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.solarSystems).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Station(ctx context.Context, id uint) (*athena.Station, error) {

	stations, err := r.Stations(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(stations) != 1 {
		return nil, nil
	}

	return stations[0], nil

}

func (r *universeRepository) Stations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Station, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "system_id", "type_id",
		"race_id", "owner_corporation_id", "max_dockable_ship_volume", "office_rental_cost",
		"reprocessing_efficiency", "reprocessing_stations_take",
		"created_at", "updated_at",
	).From(r.stations), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var stations = make([]*athena.Station, 0)
	err = r.db.SelectContext(ctx, &stations, query, args...)

	return stations, err

}

func (r *universeRepository) CreateStation(ctx context.Context, station *athena.Station) (*athena.Station, error) {

	query, args, err := sq.Insert(r.stations).Columns(
		"id", "name", "system_id", "type_id",
		"race_id", "owner_corporation_id", "max_dockable_ship_volume", "office_rental_cost",
		"reprocessing_efficiency", "reprocessing_stations_take",
		"created_at", "updated_at",
	).Values(
		station.ID,
		station.Name,
		station.SystemID,
		station.TypeID,
		station.RaceID,
		station.OwnerCorporationID,
		station.MaxDockableShipVolume,
		station.OfficeRentalCost,
		station.ReprocessingEfficiency,
		station.ReprocessingStationsTake,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Station(ctx, station.ID)

}

func (r *universeRepository) UpdateStation(ctx context.Context, id uint, station *athena.Station) (*athena.Station, error) {

	query, args, err := sq.Update(r.stations).
		Set("name", station.Name).
		Set("system_id", station.SystemID).
		Set("type_id", station.TypeID).
		Set("race_id", station.RaceID).
		Set("owner_corporation_id", station.OwnerCorporationID).
		Set("max_dockable_ship_volume", station.MaxDockableShipVolume).
		Set("office_rental_cost", station.OfficeRentalCost).
		Set("reprocessing_efficiency", station.ReprocessingEfficiency).
		Set("reprocessing_stations_take", station.ReprocessingStationsTake).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Station(ctx, id)

}

func (r *universeRepository) DeleteStation(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.stations).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Structure(ctx context.Context, id uint64) (*athena.Structure, error) {

	structures, err := r.Structures(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(structures) != 1 {
		return nil, nil
	}

	return structures[0], nil

}

func (r *universeRepository) Structures(ctx context.Context, operators ...*athena.Operator) ([]*athena.Structure, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "owner_id",
		"solar_system_id", "type_id",
		"created_at", "updated_at",
	).From(r.structures), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var structures = make([]*athena.Structure, 0)
	err = r.db.SelectContext(ctx, &structures, query, args...)

	return structures, err

}

func (r *universeRepository) CreateStructure(ctx context.Context, structure *athena.Structure) (*athena.Structure, error) {

	query, args, err := sq.Insert(r.structures).Columns(
		"id", "name", "owner_id",
		"solar_system_id", "type_id",
		"created_at", "updated_at",
	).Values(
		structure.ID,
		structure.Name,
		structure.OwnerID,
		structure.SolarSystemID,
		structure.TypeID,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Structure(ctx, structure.ID)

}

func (r *universeRepository) UpdateStructure(ctx context.Context, id uint64, structure *athena.Structure) (*athena.Structure, error) {

	query, args, err := sq.Update(r.structures).
		Set("name", structure.Name).
		Set("owner_id", structure.OwnerID).
		Set("solar_system_id", structure.SolarSystemID).
		Set("type_id", structure.TypeID).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Structure(ctx, id)

}

func (r *universeRepository) DeleteStructure(ctx context.Context, id uint64) (bool, error) {

	query, args, err := sq.Delete(r.structures).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}

func (r *universeRepository) Type(ctx context.Context, id uint) (*athena.Type, error) {

	types, err := r.Types(ctx, athena.NewEqualOperator("id", id))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if len(types) != 1 {
		return nil, nil
	}

	return types[0], nil

}

func (r *universeRepository) Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "name", "group_id", "published", "capacity",
		"market_group_id", "mass", "packaged_volume", "portion_size", "radius",
		"volume", "created_at", "updated_at",
	).From(r.types), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	var types = make([]*athena.Type, 0)
	err = r.db.SelectContext(ctx, &types, query, args...)

	return types, err

}

func (r *universeRepository) CreateType(ctx context.Context, item *athena.Type) (*athena.Type, error) {

	query, args, err := sq.Insert(r.types).Columns(
		"id", "name", "group_id", "published", "capacity",
		"market_group_id", "mass", "packaged_volume", "portion_size", "radius",
		"volume", "created_at", "updated_at",
	).Values(
		item.ID, item.Name, item.GroupID, item.Published,
		item.Capacity, item.MarketGroupID, item.Mass, item.PackagedVolume,
		item.PortionSize, item.Radius, item.Volume,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Type(ctx, item.ID)

}

func (r *universeRepository) UpdateType(ctx context.Context, id uint, item *athena.Type) (*athena.Type, error) {

	query, args, err := sq.Update(r.types).
		Set("name", item.Name).
		Set("group_id", item.GroupID).
		Set("published", item.Published).
		Set("capacity", item.Capacity).
		Set("market_group_id", item.MarketGroupID).
		Set("mass", item.Mass).
		Set("packaged_volume", item.PackagedVolume).
		Set("portion_size", item.PortionSize).
		Set("radius", item.Radius).
		Set("volume", item.Volume).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository] Failed to update record: %w", err)
	}

	return r.Type(ctx, id)

}

func (r *universeRepository) DeleteType(ctx context.Context, id uint) (bool, error) {

	query, args, err := sq.Delete(r.types).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Universe Repository] Failed to generate sql query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err == nil, err

}
