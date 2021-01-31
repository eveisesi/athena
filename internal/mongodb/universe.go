package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type universeRepository struct {
	ancestries     *mongo.Collection
	bloodlines     *mongo.Collection
	categories     *mongo.Collection
	constellations *mongo.Collection
	factions       *mongo.Collection
	groups         *mongo.Collection
	items          *mongo.Collection
	races          *mongo.Collection
	regions        *mongo.Collection
	solarSystems   *mongo.Collection
	stations       *mongo.Collection
	structures     *mongo.Collection
}

func NewUniverseRepository(d *mongo.Database) (athena.UniverseRepository, error) {

	var ctx = context.Background()
	var err error

	ancestries := d.Collection("ancestries")
	ancestriesIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("ancestries_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = ancestries.Indexes().CreateOne(ctx, ancestriesIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on ancestries collection: %w", err)
	}

	bloodlines := d.Collection("bloodlines")
	bloodlinesIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"bloodline_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("bloodlines_bloodline_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = bloodlines.Indexes().CreateOne(ctx, bloodlinesIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on bloodlines collection: %w", err)
	}

	categories := d.Collection("categories")
	categoriesIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"category_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("categories_category_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = categories.Indexes().CreateOne(ctx, categoriesIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on categories collection: %w", err)
	}

	constellations := d.Collection("constellations")
	constellationsIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"constellation_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("constellations_constellation_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = constellations.Indexes().CreateOne(ctx, constellationsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on constellations collection: %w", err)
	}

	factions := d.Collection("factions")
	factionsIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"faction_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("factions_faction_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = factions.Indexes().CreateOne(ctx, factionsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on factions collection: %w", err)
	}

	groups := d.Collection("groups")
	groupsIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"group_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("groups_group_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = groups.Indexes().CreateOne(ctx, groupsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on groups collection: %w", err)
	}

	items := d.Collection("items")
	itemsIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"type_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("items_type_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = items.Indexes().CreateOne(ctx, itemsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on items collection: %w", err)
	}

	races := d.Collection("races")
	racesIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"race_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("races_race_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = races.Indexes().CreateOne(ctx, racesIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on races collection: %w", err)
	}

	regions := d.Collection("regions")
	regionsIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"region_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("regions_region_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = regions.Indexes().CreateOne(ctx, regionsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on regions collection: %w", err)
	}

	solarSystems := d.Collection("solarSystems")
	solarSystemsIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"system_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("solar_system_system_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = solarSystems.Indexes().CreateOne(ctx, solarSystemsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on solarSystems collection: %w", err)
	}

	stations := d.Collection("stations")
	stationsIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"station_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("station_station_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = stations.Indexes().CreateOne(ctx, stationsIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on stations collection: %w", err)
	}

	structures := d.Collection("structures")
	structuresIdxModel := mongo.IndexModel{
		Keys: bson.M{
			"station_id": 1,
		},
		Options: &options.IndexOptions{
			Name:   newString("structures_structure_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err = structures.Indexes().CreateOne(ctx, structuresIdxModel)
	if err != nil {
		return nil, fmt.Errorf("[Universe Repository]: Failed to create index on structures collection: %w", err)
	}

	return &universeRepository{
		ancestries:     ancestries,
		bloodlines:     bloodlines,
		categories:     categories,
		constellations: constellations,
		factions:       factions,
		groups:         groups,
		items:          items,
		races:          races,
		regions:        regions,
		solarSystems:   solarSystems,
		stations:       stations,
		structures:     structures,
	}, nil

}

func (r *universeRepository) Ancestry(ctx context.Context, id int) (*athena.Ancestry, error) {

	ancestry := new(athena.Ancestry)

	err := r.ancestries.FindOne(ctx, primitive.D{primitive.E{Key: "ancestry_id", Value: id}}).Decode(ancestry)

	return ancestry, err

}

func (r *universeRepository) Ancestries(ctx context.Context, operators ...*athena.Operator) ([]*athena.Ancestry, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var ancestries = make([]*athena.Ancestry, 0)
	result, err := r.ancestries.Find(ctx, filters, options)
	if err != nil {
		return ancestries, err
	}

	err = result.All(ctx, &ancestries)

	return ancestries, err

}

func (r *universeRepository) CreateAncestry(ctx context.Context, ancestry *athena.Ancestry) (*athena.Ancestry, error) {

	ancestry.CreatedAt = time.Now()
	ancestry.UpdatedAt = time.Now()

	_, err := r.ancestries.InsertOne(ctx, ancestry)
	if err != nil {
		return nil, err
	}

	return ancestry, err

}

func (r *universeRepository) UpdateAncestry(ctx context.Context, id int, ancestry *athena.Ancestry) (*athena.Ancestry, error) {

	ancestry.AncestryID = id
	ancestry.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: ancestry}}

	_, err := r.ancestries.UpdateOne(ctx, primitive.D{primitive.E{Key: "ancestry_id", Value: id}}, update)

	return ancestry, err

}

func (r *universeRepository) DeleteAncestry(ctx context.Context, id int) (bool, error) {

	_, err := r.ancestries.DeleteOne(ctx, primitive.D{primitive.E{Key: "ancestry_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Bloodline(ctx context.Context, id int) (*athena.Bloodline, error) {

	bloodline := new(athena.Bloodline)

	err := r.bloodlines.FindOne(ctx, primitive.D{primitive.E{Key: "bloodline_id", Value: id}}).Decode(bloodline)

	return bloodline, err

}

func (r *universeRepository) Bloodlines(ctx context.Context, operators ...*athena.Operator) ([]*athena.Bloodline, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var bloodlines = make([]*athena.Bloodline, 0)
	result, err := r.bloodlines.Find(ctx, filters, options)
	if err != nil {
		return bloodlines, err
	}

	err = result.All(ctx, &bloodlines)

	return bloodlines, err

}

func (r *universeRepository) CreateBloodline(ctx context.Context, bloodline *athena.Bloodline) (*athena.Bloodline, error) {

	bloodline.CreatedAt = time.Now()
	bloodline.UpdatedAt = time.Now()

	_, err := r.bloodlines.InsertOne(ctx, bloodline)
	if err != nil {
		return nil, err
	}

	return bloodline, err

}

func (r *universeRepository) UpdateBloodline(ctx context.Context, id int, bloodline *athena.Bloodline) (*athena.Bloodline, error) {

	bloodline.BloodlineID = id
	bloodline.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: bloodline}}

	_, err := r.bloodlines.UpdateOne(ctx, primitive.D{primitive.E{Key: "bloodline_id", Value: id}}, update)

	return bloodline, err

}

func (r *universeRepository) DeleteBloodline(ctx context.Context, id int) (bool, error) {

	_, err := r.bloodlines.DeleteOne(ctx, primitive.D{primitive.E{Key: "bloodline_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Category(ctx context.Context, id int) (*athena.Category, error) {

	category := new(athena.Category)

	err := r.categories.FindOne(ctx, primitive.D{primitive.E{Key: "category_id", Value: id}}).Decode(category)

	return category, err

}

func (r *universeRepository) Categories(ctx context.Context, operators ...*athena.Operator) ([]*athena.Category, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var categories = make([]*athena.Category, 0)
	result, err := r.categories.Find(ctx, filters, options)
	if err != nil {
		return categories, err
	}

	err = result.All(ctx, &categories)

	return categories, err

}

func (r *universeRepository) CreateCategory(ctx context.Context, group *athena.Category) (*athena.Category, error) {

	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	_, err := r.categories.InsertOne(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, err

}

func (r *universeRepository) UpdateCategory(ctx context.Context, id int, category *athena.Category) (*athena.Category, error) {

	category.CategoryID = id
	category.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: category}}

	_, err := r.categories.UpdateOne(ctx, primitive.D{primitive.E{Key: "category_id", Value: id}}, update)

	return category, err

}

func (r *universeRepository) DeleteCategory(ctx context.Context, id int) (bool, error) {

	_, err := r.categories.DeleteOne(ctx, primitive.D{primitive.E{Key: "category_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Constellation(ctx context.Context, id int) (*athena.Constellation, error) {

	constellation := new(athena.Constellation)

	err := r.constellations.FindOne(ctx, primitive.D{primitive.E{Key: "constellation_id", Value: id}}).Decode(constellation)

	return constellation, err

}

func (r *universeRepository) Constellations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Constellation, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var constellations = make([]*athena.Constellation, 0)
	result, err := r.constellations.Find(ctx, filters, options)
	if err != nil {
		return constellations, err
	}

	err = result.All(ctx, &constellations)

	return constellations, err

}

func (r *universeRepository) CreateConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, error) {

	constellation.CreatedAt = time.Now()
	constellation.UpdatedAt = time.Now()

	_, err := r.constellations.InsertOne(ctx, constellation)
	if err != nil {
		return nil, err
	}

	return constellation, err

}

func (r *universeRepository) UpdateConstellation(ctx context.Context, id int, constellation *athena.Constellation) (*athena.Constellation, error) {

	constellation.ConstellationID = id
	constellation.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: constellation}}

	_, err := r.constellations.UpdateOne(ctx, primitive.D{primitive.E{Key: "constellation_id", Value: id}}, update)

	return constellation, err

}

func (r *universeRepository) DeleteConstellation(ctx context.Context, id int) (bool, error) {

	_, err := r.constellations.DeleteOne(ctx, primitive.D{primitive.E{Key: "constellation_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Faction(ctx context.Context, id int) (*athena.Faction, error) {

	faction := new(athena.Faction)

	err := r.factions.FindOne(ctx, primitive.D{primitive.E{Key: "faction_id", Value: id}}).Decode(faction)

	return faction, err

}

func (r *universeRepository) Factions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Faction, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var factions = make([]*athena.Faction, 0)
	result, err := r.factions.Find(ctx, filters, options)
	if err != nil {
		return factions, err
	}

	err = result.All(ctx, &factions)

	return factions, err

}

func (r *universeRepository) CreateFaction(ctx context.Context, faction *athena.Faction) (*athena.Faction, error) {

	faction.CreatedAt = time.Now()
	faction.UpdatedAt = time.Now()

	_, err := r.factions.InsertOne(ctx, faction)
	if err != nil {
		return nil, err
	}

	return faction, err

}

func (r *universeRepository) UpdateFaction(ctx context.Context, id int, faction *athena.Faction) (*athena.Faction, error) {

	faction.FactionID = id
	faction.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: faction}}

	_, err := r.factions.UpdateOne(ctx, primitive.D{primitive.E{Key: "faction_id", Value: id}}, update)

	return faction, err

}

func (r *universeRepository) DeleteFaction(ctx context.Context, id int) (bool, error) {

	_, err := r.factions.DeleteOne(ctx, primitive.D{primitive.E{Key: "faction_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Group(ctx context.Context, id int) (*athena.Group, error) {

	group := new(athena.Group)

	err := r.groups.FindOne(ctx, primitive.D{primitive.E{Key: "group_id", Value: id}}).Decode(group)

	return group, err

}

func (r *universeRepository) Groups(ctx context.Context, operators ...*athena.Operator) ([]*athena.Group, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var groups = make([]*athena.Group, 0)
	result, err := r.groups.Find(ctx, filters, options)
	if err != nil {
		return groups, err
	}

	err = result.All(ctx, &groups)

	return groups, err

}

func (r *universeRepository) CreateGroup(ctx context.Context, group *athena.Group) (*athena.Group, error) {

	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	_, err := r.groups.InsertOne(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, err

}

func (r *universeRepository) UpdateGroup(ctx context.Context, id int, group *athena.Group) (*athena.Group, error) {

	group.GroupID = id
	group.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: group}}

	_, err := r.groups.UpdateOne(ctx, primitive.D{primitive.E{Key: "group_id", Value: id}}, update)

	return group, err

}

func (r *universeRepository) DeleteGroup(ctx context.Context, id int) (bool, error) {

	_, err := r.groups.DeleteOne(ctx, primitive.D{primitive.E{Key: "group_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Race(ctx context.Context, id int) (*athena.Race, error) {

	race := new(athena.Race)

	err := r.races.FindOne(ctx, primitive.D{primitive.E{Key: "race_id", Value: id}}).Decode(race)

	return race, err

}

func (r *universeRepository) Races(ctx context.Context, operators ...*athena.Operator) ([]*athena.Race, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var races = make([]*athena.Race, 0)
	result, err := r.races.Find(ctx, filters, options)
	if err != nil {
		return races, err
	}

	err = result.All(ctx, &races)

	return races, err

}

func (r *universeRepository) CreateRace(ctx context.Context, race *athena.Race) (*athena.Race, error) {

	race.CreatedAt = time.Now()
	race.UpdatedAt = time.Now()

	_, err := r.races.InsertOne(ctx, race)
	if err != nil {
		return nil, err
	}

	return race, err

}

func (r *universeRepository) UpdateRace(ctx context.Context, id int, race *athena.Race) (*athena.Race, error) {

	race.RaceID = id
	race.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: race}}

	_, err := r.races.UpdateOne(ctx, primitive.D{primitive.E{Key: "race_id", Value: id}}, update)

	return race, err

}

func (r *universeRepository) DeleteRace(ctx context.Context, id int) (bool, error) {

	_, err := r.races.DeleteOne(ctx, primitive.D{primitive.E{Key: "race_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Region(ctx context.Context, id int) (*athena.Region, error) {

	region := new(athena.Region)

	err := r.regions.FindOne(ctx, primitive.D{primitive.E{Key: "region_id", Value: id}}).Decode(region)

	return region, err

}

func (r *universeRepository) Regions(ctx context.Context, operators ...*athena.Operator) ([]*athena.Region, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var regions = make([]*athena.Region, 0)
	result, err := r.regions.Find(ctx, filters, options)
	if err != nil {
		return regions, err
	}

	err = result.All(ctx, &regions)

	return regions, err

}

func (r *universeRepository) CreateRegion(ctx context.Context, region *athena.Region) (*athena.Region, error) {

	region.CreatedAt = time.Now()
	region.UpdatedAt = time.Now()

	_, err := r.regions.InsertOne(ctx, region)
	if err != nil {
		return nil, err
	}

	return region, err

}

func (r *universeRepository) UpdateRegion(ctx context.Context, id int, region *athena.Region) (*athena.Region, error) {

	region.RegionID = id
	region.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: region}}

	_, err := r.regions.UpdateOne(ctx, primitive.D{primitive.E{Key: "region_id", Value: id}}, update)

	return region, err

}

func (r *universeRepository) DeleteRegion(ctx context.Context, id int) (bool, error) {

	_, err := r.regions.DeleteOne(ctx, primitive.D{primitive.E{Key: "region_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) SolarSystem(ctx context.Context, id int) (*athena.SolarSystem, error) {

	solarSystem := new(athena.SolarSystem)

	err := r.solarSystems.FindOne(ctx, primitive.D{primitive.E{Key: "system_id", Value: id}}).Decode(solarSystem)

	return solarSystem, err

}

func (r *universeRepository) SolarSystems(ctx context.Context, operators ...*athena.Operator) ([]*athena.SolarSystem, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var solarSystems = make([]*athena.SolarSystem, 0)
	result, err := r.solarSystems.Find(ctx, filters, options)
	if err != nil {
		return solarSystems, err
	}

	err = result.All(ctx, &solarSystems)

	return solarSystems, err

}

func (r *universeRepository) CreateSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, error) {

	solarSystem.CreatedAt = time.Now()
	solarSystem.UpdatedAt = time.Now()

	_, err := r.solarSystems.InsertOne(ctx, solarSystem)
	if err != nil {
		return nil, err
	}

	return solarSystem, err

}

func (r *universeRepository) UpdateSolarSystem(ctx context.Context, id int, solarSystem *athena.SolarSystem) (*athena.SolarSystem, error) {

	solarSystem.SystemID = id
	solarSystem.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: solarSystem}}

	_, err := r.solarSystems.UpdateOne(ctx, primitive.D{primitive.E{Key: "system_id", Value: id}}, update)

	return solarSystem, err

}

func (r *universeRepository) DeleteSolarSystem(ctx context.Context, id int) (bool, error) {

	_, err := r.solarSystems.DeleteOne(ctx, primitive.D{primitive.E{Key: "system_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Station(ctx context.Context, id int) (*athena.Station, error) {

	station := new(athena.Station)

	err := r.stations.FindOne(ctx, primitive.D{primitive.E{Key: "station_id", Value: id}}).Decode(station)

	return station, err

}

func (r *universeRepository) Stations(ctx context.Context, operators ...*athena.Operator) ([]*athena.Station, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var stations = make([]*athena.Station, 0)
	result, err := r.stations.Find(ctx, filters, options)
	if err != nil {
		return stations, err
	}

	err = result.All(ctx, &stations)

	return stations, err

}

func (r *universeRepository) CreateStation(ctx context.Context, station *athena.Station) (*athena.Station, error) {

	station.CreatedAt = time.Now()
	station.UpdatedAt = time.Now()

	_, err := r.stations.InsertOne(ctx, station)
	if err != nil {
		return nil, err
	}

	return station, err

}

func (r *universeRepository) UpdateStation(ctx context.Context, id int, station *athena.Station) (*athena.Station, error) {

	station.SystemID = id
	station.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: station}}

	_, err := r.stations.UpdateOne(ctx, primitive.D{primitive.E{Key: "station_id", Value: id}}, update)

	return station, err

}

func (r *universeRepository) DeleteStation(ctx context.Context, id int) (bool, error) {

	_, err := r.stations.DeleteOne(ctx, primitive.D{primitive.E{Key: "station_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Structure(ctx context.Context, id int64) (*athena.Structure, error) {

	solarSystem := new(athena.Structure)

	err := r.solarSystems.FindOne(ctx, primitive.D{primitive.E{Key: "structure_id", Value: id}}).Decode(solarSystem)

	return solarSystem, err

}

func (r *universeRepository) Structures(ctx context.Context, operators ...*athena.Operator) ([]*athena.Structure, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var structures = make([]*athena.Structure, 0)
	result, err := r.structures.Find(ctx, filters, options)
	if err != nil {
		return structures, err
	}

	err = result.All(ctx, &structures)

	return structures, err

}

func (r *universeRepository) CreateStructure(ctx context.Context, structure *athena.Structure) (*athena.Structure, error) {

	structure.CreatedAt = time.Now()
	structure.UpdatedAt = time.Now()

	_, err := r.structures.InsertOne(ctx, structure)
	if err != nil {
		return nil, err
	}

	return structure, err

}

func (r *universeRepository) UpdateStructure(ctx context.Context, id int64, structure *athena.Structure) (*athena.Structure, error) {

	structure.StructureID = id
	structure.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: structure}}

	_, err := r.structures.UpdateOne(ctx, primitive.D{primitive.E{Key: "structure_id", Value: id}}, update)

	return structure, err

}

func (r *universeRepository) DeleteStructure(ctx context.Context, id int64) (bool, error) {

	_, err := r.structures.DeleteOne(ctx, primitive.D{primitive.E{Key: "structure_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}

func (r *universeRepository) Type(ctx context.Context, id int) (*athena.Type, error) {

	item := new(athena.Type)

	err := r.items.FindOne(ctx, primitive.D{primitive.E{Key: "type_id", Value: id}}).Decode(item)

	return item, err

}

func (r *universeRepository) Types(ctx context.Context, operators ...*athena.Operator) ([]*athena.Type, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var items = make([]*athena.Type, 0)
	result, err := r.items.Find(ctx, filters, options)
	if err != nil {
		return items, err
	}

	err = result.All(ctx, &items)

	return items, err

}

func (r *universeRepository) CreateType(ctx context.Context, item *athena.Type) (*athena.Type, error) {

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	_, err := r.items.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, err

}

func (r *universeRepository) UpdateType(ctx context.Context, id int, item *athena.Type) (*athena.Type, error) {

	item.TypeID = id
	item.UpdatedAt = time.Now()

	update := primitive.D{primitive.E{Key: "$set", Value: item}}

	_, err := r.items.UpdateOne(ctx, primitive.D{primitive.E{Key: "type_id", Value: id}}, update)

	return item, err

}

func (r *universeRepository) DeleteType(ctx context.Context, id int) (bool, error) {

	_, err := r.items.DeleteOne(ctx, primitive.D{primitive.E{Key: "type_id", Value: id}})
	if err != nil {
		return false, err
	}

	return err == nil, err

}
