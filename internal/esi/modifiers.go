package esi

import "github.com/eveisesi/athena"

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

	ModifierFunc func(mod *modifiers) *modifiers

	endpointMap map[Endpoint]func(modFuncs ...ModifierFunc) string
)

func (s *service) modifiers(modFuncs []ModifierFunc) *modifiers {

	mod := &modifiers{}
	for _, modFunc := range modFuncs {
		mod = modFunc(mod)
	}

	return mod

}

func ModWithPage(page *int) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.page = page
		return mod
	}
}

func ModWithMember(member *athena.Member) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.member = member
		return mod
	}
}

func ModWithAlliance(alliance *athena.Alliance) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.alliance = alliance
		return mod
	}
}

func ModWithCategory(category *athena.Category) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.category = category
		return mod
	}
}

func ModWithCharacter(character *athena.Character) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.character = character
		return mod
	}
}

func ModWithCorporation(corporation *athena.Corporation) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.corporation = corporation
		return mod
	}
}

func ModWithConstellation(constellation *athena.Constellation) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.constellation = constellation
		return mod
	}
}

func ModWithGroup(group *athena.Group) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.group = group
		return mod
	}
}

func ModWithItem(item *athena.Type) ModifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.item = item
		return mod
	}
}
