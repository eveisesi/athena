package esi

import (
	"github.com/eveisesi/athena"
)

type (
	modifiers struct {
		alliance      *athena.Alliance
		asteroidBelt  *athena.AsteroidBelt
		category      *athena.Category
		character     *athena.Character
		constellation *athena.Constellation
		contract      *athena.MemberContract
		corporation   *athena.Corporation
		group         *athena.Group
		item          *athena.Type
		member        *athena.Member
		moon          *athena.Moon
		page          *int
		planet        *athena.Planet
		region        *athena.Region
		station       *athena.Station
		solarSystem   *athena.SolarSystem
		structure     *athena.Structure
	}

	modifierFunc func(mod *modifiers) *modifiers

	pathFunc func(mod *modifiers) string
	keyFunc  func(mod *modifiers) string

	endpointMap map[endpointID]*endpoint
)

func (s *service) modifiers(modFuncs ...modifierFunc) *modifiers {

	mod := &modifiers{}
	for _, modFunc := range modFuncs {
		mod = modFunc(mod)
	}

	return mod

}

func ModWithPage(page *int) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.page = page
		return mod
	}
}

func ModWithAsteroidBelt(belt *athena.AsteroidBelt) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.asteroidBelt = belt
		return mod
	}
}

func requireAsteriodBelt(mods *modifiers) {
	if mods.asteroidBelt == nil {
		panic("expected type *athena.AsteroidBelt to be provided, received nil instead")
	}
}

func ModWithMoon(moon *athena.Moon) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.moon = moon
		return mod
	}
}

func requireMoon(mods *modifiers) {
	if mods.moon == nil {
		panic("expected type *athena.Moon to be provided, received nil instead")
	}
}

func ModWithContract(contract *athena.MemberContract) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.contract = contract
		return mod
	}
}

func requireContract(mods *modifiers) {
	if mods.contract == nil {
		panic("expected type *athena.MemberContract to be provided, received nil instead")
	}
}

func ModWithMember(member *athena.Member) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.member = member
		return mod
	}
}

func requireMember(mods *modifiers) {
	if mods.member == nil {
		panic("expected type *athena.Member to be provided, received nil instead")
	}
}

func ModWithAlliance(alliance *athena.Alliance) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.alliance = alliance
		return mod
	}
}

func requireAlliance(mods *modifiers) {
	if mods.alliance == nil {
		panic("expected type *athena.Alliance to be provided, received nil instead")
	}
}

func ModWithCategory(category *athena.Category) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.category = category
		return mod
	}
}

func requireCategory(mods *modifiers) {
	if mods.category == nil {
		panic("expected type *athena.Category to be provided, received nil instead")
	}
}

func ModWithCharacter(character *athena.Character) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.character = character
		return mod
	}
}

func requireCharacter(mods *modifiers) {
	if mods.character == nil {
		panic("expected type *athena.Character to be provided, received nil instead")
	}
}

func ModWithCorporation(corporation *athena.Corporation) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.corporation = corporation
		return mod
	}
}

func requireCorporation(mods *modifiers) {
	if mods.corporation == nil {
		panic("expected type *athena.Corporation to be provided, received nil instead")
	}
}

func ModWithPlanet(planet *athena.Planet) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.planet = planet
		return mod
	}
}

func requirePlanet(mods *modifiers) {
	if mods.planet == nil {
		panic("expected type *athena.Planet to be provided, received nil instead")
	}
}

func ModWithRegion(region *athena.Region) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.region = region
		return mod
	}
}

func requireRegion(mods *modifiers) {
	if mods.region == nil {
		panic("expected type *athena.Region to be provided, received nil instead")
	}
}

func ModWithConstellation(constellation *athena.Constellation) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.constellation = constellation
		return mod
	}
}

func requireConstellation(mods *modifiers) {
	if mods.constellation == nil {
		panic("expected type *athena.Constellation to be provided, received nil instead")
	}
}

func ModWithGroup(group *athena.Group) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.group = group
		return mod
	}
}

func requireGroup(mods *modifiers) {
	if mods.group == nil {
		panic("expected type *athena.Group to be provided, received nil instead")
	}
}

func ModWithItem(item *athena.Type) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.item = item
		return mod
	}
}

func requireItem(mods *modifiers) {
	if mods.item == nil {
		panic("expected type *athena.Item to be provided, received nil instead")
	}
}

func ModWithSystem(system *athena.SolarSystem) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.solarSystem = system
		return mod
	}
}

func requireSystem(mods *modifiers) {
	if mods.solarSystem == nil {
		panic("expected type *athena.System to be provided, received nil instead")
	}
}

func ModWithStation(station *athena.Station) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.station = station
		return mod
	}
}

func requireStation(mods *modifiers) {
	if mods.station == nil {
		panic("expected type *athena.station to be provided, received nil instead")
	}
}

func ModWithStructure(structure *athena.Structure) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.structure = structure
		return mod
	}
}

func requireStructure(mods *modifiers) {
	if mods.structure == nil {
		panic("expected type *athena.Structure to be provided, received nil instead")
	}
}
