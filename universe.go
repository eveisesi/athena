package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type UniverseRepository interface {
	ancestryRepository
	bloodlineRepository
	categoryRepository
	constellationRepository
	factionRepository
	groupRepository
	raceRepository
	regionRepository
	solarSystemRepository
	stationRepository
	structureRepository
	typeRepository
}

type ancestryRepository interface {
	Ancestry(ctx context.Context, id int) (*Ancestry, error)
	Ancestries(ctx context.Context, operators ...*Operator) ([]*Ancestry, error)
	CreateAncestry(ctx context.Context, ancestry *Ancestry) (*Ancestry, error)
	UpdateAncestry(ctx context.Context, id int, ancestry *Ancestry) (*Ancestry, error)
	DeleteAncestry(ctx context.Context, id int) (bool, error)
}

type bloodlineRepository interface {
	Bloodline(ctx context.Context, id int) (*Bloodline, error)
	Bloodlines(ctx context.Context, operators ...*Operator) ([]*Bloodline, error)
	CreateBloodline(ctx context.Context, bloodline *Bloodline) (*Bloodline, error)
	UpdateBloodline(ctx context.Context, id int, bloodline *Bloodline) (*Bloodline, error)
	DeleteBloodline(ctx context.Context, id int) (bool, error)
}

type categoryRepository interface {
	Category(ctx context.Context, id int) (*Category, error)
	Categories(ctx context.Context, operators ...*Operator) ([]*Category, error)
	CreateCategory(ctx context.Context, group *Category) (*Category, error)
	UpdateCategory(ctx context.Context, id int, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, id int) (bool, error)
}

type constellationRepository interface {
	Constellation(ctx context.Context, id int) (*Constellation, error)
	Constellations(ctx context.Context, operators ...*Operator) ([]*Constellation, error)
	CreateConstellation(ctx context.Context, constellation *Constellation) (*Constellation, error)
	UpdateConstellation(ctx context.Context, id int, constellation *Constellation) (*Constellation, error)
	DeleteConstellation(ctx context.Context, id int) (bool, error)
}

type factionRepository interface {
	Faction(ctx context.Context, id int) (*Faction, error)
	Factions(ctx context.Context, operators ...*Operator) ([]*Faction, error)
	CreateFaction(ctx context.Context, faction *Faction) (*Faction, error)
	UpdateFaction(ctx context.Context, id int, faction *Faction) (*Faction, error)
	DeleteFaction(ctx context.Context, id int) (bool, error)
}

type groupRepository interface {
	Group(ctx context.Context, id int) (*Group, error)
	Groups(ctx context.Context, operators ...*Operator) ([]*Group, error)
	CreateGroup(ctx context.Context, group *Group) (*Group, error)
	UpdateGroup(ctx context.Context, id int, group *Group) (*Group, error)
	DeleteGroup(ctx context.Context, id int) (bool, error)
}

type raceRepository interface {
	Race(ctx context.Context, id int) (*Race, error)
	Races(ctx context.Context, operators ...*Operator) ([]*Race, error)
	CreateRace(ctx context.Context, race *Race) (*Race, error)
	UpdateRace(ctx context.Context, id int, race *Race) (*Race, error)
	DeleteRace(ctx context.Context, id int) (bool, error)
}

type regionRepository interface {
	Region(ctx context.Context, id int) (*Region, error)
	Regions(ctx context.Context, operators ...*Operator) ([]*Region, error)
	CreateRegion(ctx context.Context, region *Region) (*Region, error)
	UpdateRegion(ctx context.Context, id int, region *Region) (*Region, error)
	DeleteRegion(ctx context.Context, id int) (bool, error)
}

type solarSystemRepository interface {
	SolarSystem(ctx context.Context, id int) (*SolarSystem, error)
	SolarSystems(ctx context.Context, operators ...*Operator) ([]*SolarSystem, error)
	CreateSolarSystem(ctx context.Context, solarSystem *SolarSystem) (*SolarSystem, error)
	UpdateSolarSystem(ctx context.Context, id int, solarSystem *SolarSystem) (*SolarSystem, error)
	DeleteSolarSystem(ctx context.Context, id int) (bool, error)
}

type stationRepository interface {
	Station(ctx context.Context, id int) (*Station, error)
	Stations(ctx context.Context, operators ...*Operator) ([]*Station, error)
	CreateStation(ctx context.Context, station *Station) (*Station, error)
	UpdateStation(ctx context.Context, id int, solarSystem *Station) (*Station, error)
	DeleteStation(ctx context.Context, id int) (bool, error)
}

type structureRepository interface {
	Structure(ctx context.Context, id int64) (*Structure, error)
	Structures(ctx context.Context, operators ...*Operator) ([]*Structure, error)
	CreateStructure(ctx context.Context, solarSystem *Structure) (*Structure, error)
	UpdateStructure(ctx context.Context, id int64, struture *Structure) (*Structure, error)
	DeleteStructure(ctx context.Context, id int64) (bool, error)
}

type typeRepository interface {
	Type(ctx context.Context, id int) (*Type, error)
	Types(ctx context.Context, operators ...*Operator) ([]*Type, error)
	CreateType(ctx context.Context, item *Type) (*Type, error)
	UpdateType(ctx context.Context, id int, item *Type) (*Type, error)
	DeleteType(ctx context.Context, id int) (bool, error)
}

type Ancestry struct {
	AncestryID  int       `bson:"id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	BloodlineID int       `bson:"bloodline_id" json:"bloodline_id"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type Bloodline struct {
	BloodlineID   int       `bson:"bloodline_id" json:"bloodline_id"`
	Name          string    `bson:"name" json:"name"`
	RaceID        int       `bson:"race_id" json:"race_id"`
	CorporationID int       `bson:"corporation_id" json:"corporation_id"`
	ShipTypeID    int       `bson:"ship_type_id" json:"ship_type_id"`
	Charisma      int       `bson:"charisma" json:"charisma"`
	Intelligence  int       `bson:"intelligence" json:"intelligence"`
	Memory        int       `bson:"memory" json:"memory"`
	Perception    int       `bson:"perception" json:"perception"`
	Willpower     int       `bson:"willpower" json:"willpower"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type Category struct {
	CategoryID int       `bson:"category_id" json:"category_id"`
	Name       string    `bson:"name" json:"name"`
	Published  bool      `bson:"published" json:"published"`
	Groups     []int     `bson:"-" json:"groups,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type Constellation struct {
	ConstellationID int       `bson:"constellation_id" json:"constellation_id"`
	Name            string    `bson:"name" json:"name"`
	RegionID        int       `bson:"region_id" json:"region_id"`
	SystemIDs       []int     `bson:"-" json:"systems,omitempty"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}

type Faction struct {
	FactionID            int       `bson:"faction_id" json:"faction_id"`
	Name                 string    `bson:"name" json:"name"`
	IsUnique             bool      `bson:"is_unique" json:"is_unique"`
	SizeFactor           float64   `bson:"size_factor" json:"size_factor"`
	StationCount         int       `bson:"station_count" json:"station_count"`
	StationSystemCount   int       `bson:"station_system_count" json:"station_system_count"`
	CorporationID        null.Int  `bson:"corporation_id,omitempty" json:"corporation_id,omitempty"`
	MilitiaCorporationID null.Int  `bson:"militia_corporation_id,omitempty" json:"militia_corporation_id,omitempty"`
	SolarSystemID        null.Int  `bson:"solar_system_id,omitempty" json:"solar_system_id,omitempty"`
	CreatedAt            time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time `bson:"updated_at" json:"updated_at"`
}

type Group struct {
	GroupID    int       `bson:"group_id" json:"group_id"`
	Name       string    `bson:"name" json:"name"`
	Published  bool      `bson:"published" json:"published"`
	CategoryID int       `bson:"category_id" json:"category_id"`
	Types      []int     `bson:"-" json:"types,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type Race struct {
	RaceID     int       `bson:"race_id" json:"race_id"`
	Name       string    `bson:"name" json:"name"`
	AllianceID int       `bson:"alliance_id" json:"alliance_id"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type Region struct {
	RegionID         int       `bson:"region_id" json:"region_id"`
	Name             string    `bson:"name" json:"name"`
	ConstellationIDs []int     `bson:"-" json:"constellations,omitempty"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}

type SolarSystem struct {
	SystemID        int         `bson:"system_id" json:"system_id"`
	Name            string      `bson:"name" json:"name"`
	ConstellationID int         `bson:"constellation_id" json:"constellation_id"`
	SecurityStatus  float64     `bson:"security_status" json:"security_status"`
	StarID          null.Int    `bson:"star_id,omitempty" json:"star_id,omitempty"`
	SecurityClass   null.String `bson:"security_class,omitempty" json:"security_class,omitempty"`
	CreatedAt       time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `bson:"updated_at" json:"updated_at"`
}

type Station struct {
	MaxDockableShipVolue     float64   `bson:"max_dockable_ship_volume" json:"max_dockable_ship_volume"`
	Name                     string    `bson:"name" json:"name"`
	OfficeRentalCost         float64   `bson:"office_rental_cost" json:"office_rental_cost"`
	OwnerCorporationID       null.Int  `bson:"owner_corporation_id" json:"owner"`
	RaceID                   null.Int  `bson:"race_id" json:"race_id"`
	ReprocessingEfficiency   float64   `bson:"reprocessing_efficiency" json:"reprocessing_efficiency"`
	ReprocessingStationsTake float64   `bson:"reprocessing_stations_take" json:"reprocessing_stations_take"`
	Services                 []string  `bson:"services" json:"services"`
	StationID                int       `bson:"station_id" json:"station_id"`
	SystemID                 int       `bson:"system_id" json:"system_id"`
	TypeID                   int       `bson:"type_id" json:"type_id"`
	CreatedAt                time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time `bson:"updated_at" json:"updated_at"`
}

type Structure struct {
	StructureID   int64     `bson:"structure_id" json:"structure_id"`
	Name          string    `bson:"name" json:"name"`
	OwnerID       int       `bson:"owner_id" json:"owner_id"`
	SolarSystemID int       `bson:"solar_system_id" json:"solar_system_id"`
	TypeID        null.Int  `bson:"type_id" json:"type_id"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type Type struct {
	TypeID         int          `bson:"type_id" json:"type_id"`
	Name           string       `bson:"name" json:"name"`
	GroupID        int          `bson:"group_id" json:"group_id"`
	Published      bool         `bson:"published" json:"published"`
	Capacity       null.Float64 `bson:"capacity,omitempty" json:"capacity,omitempty"`
	MarketGroupID  null.Int     `bson:"market_group_id,omitempty" json:"market_group_id,omitempty"`
	Mass           null.Float64 `bson:"mass,omitempty" json:"mass,omitempty"`
	PackagedVolume null.Float64 `bson:"packaged_volume,omitempty" json:"packaged_volume,omitempty"`
	PortionSize    null.Int     `bson:"portion_size,omitempty" json:"portion_size,omitempty"`
	Radius         null.Float64 `bson:"radius,omitempty" json:"radius,omitempty"`
	Volume         null.Float64 `bson:"volume,omitempty" json:"volume,omitempty"`
	CreatedAt      time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time    `bson:"updated_at" json:"updated_at"`
}
