package esi

type (
	modifiers struct {
		allianceID      uint
		categoryID      uint
		characterID     uint
		constellationID uint
		contractID      uint
		corporationID   uint
		from            uint64
		groupID         uint
		mailID          uint
		itemID          uint
		lastMailID      uint64
		page            uint
		regionID        uint
		stationID       uint
		solarSystemID   uint
		structureID     uint64
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

func ModWithPage(page uint) modifierFunc {
	return func(mods *modifiers) *modifiers {
		mods.page = page
		return mods
	}
}

func requirePage(mods *modifiers) {
	if mods.page == 0 {
		panic("page modifier should be greater than zero for this request")
	}
}

func ModWithFromID(from uint64) modifierFunc {
	return func(mods *modifiers) *modifiers {
		mods.from = from
		return mods
	}
}

func ModWithLastMailID(lastMailID uint64) modifierFunc {
	return func(mods *modifiers) *modifiers {
		mods.lastMailID = lastMailID
		return mods
	}
}

func ModWithMailID(mailID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.mailID = mailID
		return mod
	}
}

func requireMailID(mods *modifiers) {
	if mods.mailID == 0 {
		panic("modifier mailID should be greater than 0")
	}
}

func ModWithContractID(contractID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.contractID = contractID
		return mod
	}
}

func requireContractID(mods *modifiers) {
	if mods.contractID == 0 {
		panic("modifier allianceID should be greater than 0")
	}
}

func ModWithAllianceID(allianceID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.allianceID = allianceID
		return mod
	}
}

func requireAllianceID(mods *modifiers) {
	if mods.allianceID == 0 {
		panic("modifier allianceID should be greater than 0")
	}
}

func ModWithCategoryID(categoryID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.categoryID = categoryID
		return mod
	}
}

func requireCategoryID(mods *modifiers) {
	if mods.categoryID == 0 {
		panic("modifier categoryID should be greater than 0")
	}
}

func ModWithCharacterID(characterID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.characterID = characterID
		return mod
	}
}

func requireCharacterID(mods *modifiers) {
	if mods.characterID == 0 {
		panic("modifier characterID should be greater than 0")
	}
}

func ModWithCorporationID(corporationID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.corporationID = corporationID
		return mod
	}
}

func requireCorporationID(mods *modifiers) {
	if mods.corporationID == 0 {
		panic("modifier corporationID should be greater than 0")
	}
}

func ModWithRegionID(regionID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.regionID = regionID
		return mod
	}
}

func requireRegionID(mods *modifiers) {
	if mods.regionID == 0 {
		panic("modifier regionID should be greater than 0")
	}
}

func ModWithConstellationID(constellationID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.constellationID = constellationID
		return mod
	}
}

func requireConstellationID(mods *modifiers) {
	if mods.constellationID == 0 {
		panic("modifier constellationID should be greater than 0")
	}
}

func ModWithGroupID(groupID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.groupID = groupID
		return mod
	}
}

func requireGroupID(mods *modifiers) {
	if mods.groupID == 0 {
		panic("modifier groupID should be greater than 0")
	}
}

func ModWithItemID(itemID uint) modifierFunc {
	return func(mods *modifiers) *modifiers {
		mods.itemID = itemID
		return mods
	}
}

func requireItemID(mods *modifiers) {
	if mods.itemID == 0 {
		panic("modifier itemID should be greater than 0")
	}
}

func ModWithSystemID(solarSystemID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.solarSystemID = solarSystemID
		return mod
	}
}

func requireSystemID(mods *modifiers) {
	if mods.solarSystemID == 0 {
		panic("modifier solarSystemID should be greater than 0")
	}
}

func ModWithStationID(stationID uint) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.stationID = stationID
		return mod
	}
}

func requireStationID(mods *modifiers) {
	if mods.stationID == 0 {
		panic("modifier stationID should be greater than 0")
	}
}

func ModWithStructureID(structureID uint64) modifierFunc {
	return func(mod *modifiers) *modifiers {
		mod.structureID = structureID
		return mod
	}
}

func requireStructureID(mods *modifiers) {
	if mods.structureID == 0 {
		panic("modifier structureID should be greater than 0")
	}
}
