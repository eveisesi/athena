package esi

import (
	"github.com/eveisesi/athena"
)

type (
	modifiers struct {
		alliance      *athena.Alliance
		category      *athena.Category
		character     *athena.Character
		constellation *athena.Constellation
		corporation   *athena.Corporation
		group         *athena.Group
		item          *athena.Type
		member        *athena.Member
		page          *int
		region        *athena.Region
		station       *athena.Station
		solarSystem   *athena.SolarSystem
		structure     *athena.Structure
	}

	modifierFunc func(mod *modifiers) *modifiers

	pathFunc func(mod *modifiers) string
	keyFunc  func(mod *modifiers) string

	endpointMap map[string]*endpoint
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

func ModWithMember(member *athena.Member) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.member = member
		return mod
	}
}

func ModWithAlliance(alliance *athena.Alliance) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.alliance = alliance
		return mod
	}
}

func ModWithCategory(category *athena.Category) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.category = category
		return mod
	}
}

func ModWithCharacter(character *athena.Character) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.character = character
		return mod
	}
}

func ModWithCorporation(corporation *athena.Corporation) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.corporation = corporation
		return mod
	}
}

func ModWithRegion(region *athena.Region) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.region = region
		return mod
	}
}

func ModWithConstellation(constellation *athena.Constellation) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.constellation = constellation
		return mod
	}
}

func ModWithGroup(group *athena.Group) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.group = group
		return mod
	}
}

func ModWithItem(item *athena.Type) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.item = item
		return mod
	}
}

func ModWithSystem(system *athena.SolarSystem) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.solarSystem = system
		return mod
	}
}

func ModWithStation(station *athena.Station) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.station = station
		return mod
	}
}

func ModWithStructure(structure *athena.Structure) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.structure = structure
		return mod
	}
}
