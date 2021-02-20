//go:generate go run github.com/vektah/dataloaden CharacterLoader uint *github.com/eveisesi/athena.Character
//go:generate go run github.com/vektah/dataloaden CharacterCorporationHistoryLoader uint []*github.com/eveisesi/athena.CharacterCorporationHistory
//go:generate go run github.com/vektah/dataloaden CorporationLoader uint *github.com/eveisesi/athena.Corporation
//go:generate go run github.com/vektah/dataloaden CorporationAllianceHistoryLoader uint *github.com/eveisesi/athena.CorporationAllianceHistory
//go:generate go run github.com/vektah/dataloaden AllianceLoader uint *github.com/eveisesi/athena.Alliance
//go:generate go run github.com/vektah/dataloaden AncestryLoader uint *github.com/eveisesi/athena.Ancestry
//go:generate go run github.com/vektah/dataloaden BloodlineLoader uint *github.com/eveisesi/athena.Bloodline
//go:generate go run github.com/vektah/dataloaden RaceLoader uint *github.com/eveisesi/athena.Race
//go:generate go run github.com/vektah/dataloaden FactionLoader uint *github.com/eveisesi/athena.Faction

//go:generate go run github.com/vektah/dataloaden RegionLoader uint *github.com/eveisesi/athena.Region
//go:generate go run github.com/vektah/dataloaden ConstellationLoader uint *github.com/eveisesi/athena.Constellation
//go:generate go run github.com/vektah/dataloaden SolarSystemLoader uint *github.com/eveisesi/athena.SolarSystem
//go:generate go run github.com/vektah/dataloaden StationLoader uint *github.com/eveisesi/athena.Station
//go:generate go run github.com/vektah/dataloaden StructureLoader uint *github.com/eveisesi/athena.Structure

//go:generate go run github.com/vektah/dataloaden CategoryLoader uint *github.com/eveisesi/athena.Category
//go:generate go run github.com/vektah/dataloaden GroupLoader uint *github.com/eveisesi/athena.Group
//go:generate go run github.com/vektah/dataloaden TypeLoader uint *github.com/eveisesi/athena.Type

package generated
