package athena

import (
	"context"
	"time"
)

type MemberAssetsRepository interface {
	MemberAsset(ctx context.Context, memberID uint, itemID uint64) (*MemberAsset, error)
	MemberAssets(ctx context.Context, memberID uint, operators ...*Operator) ([]*MemberAsset, error)
	CreateMemberAssets(ctx context.Context, memberID uint, assets []*MemberAsset) ([]*MemberAsset, error)
	UpdateMemberAssets(ctx context.Context, memberID uint, itemID uint64, asset *MemberAsset) (*MemberAsset, error)
	DeleteMemberAssets(ctx context.Context, memberID uint, assets []*MemberAsset) (bool, error)
}

type MemberAsset struct {
	MemberID        uint              `db:"member_id" json:"member_id" deep:"-"`
	ItemID          uint64            `db:"item_id" json:"item_id" deep:"-"`
	TypeID          uint              `db:"type_id" json:"type_id" deep:"-"`
	LocationID      uint64            `db:"location_id" json:"location_id"`
	LocationFlag    AssetLocationFlag `db:"location_flag" json:"location_flag"`
	LocationType    AssetLocationType `db:"location_type" json:"location_type"`
	Quantity        int               `db:"quantity" json:"quantity"`
	IsBlueprintCopy bool              `db:"is_blueprint_copy" json:"is_blueprint_copy"`
	IsSingleton     bool              `db:"is_singleton" json:"is_singleton"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at" deep:"-"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at" deep:"-"`
}

type AssetLocationFlag string

const (
	AssetLocationFlagAssetSafety                         AssetLocationFlag = "AssetSafety"
	AssetLocationFlagAutoFit                             AssetLocationFlag = "AutoFit"
	AssetLocationFlagBoosterBay                          AssetLocationFlag = "BoosterBay"
	AssetLocationFlagCargo                               AssetLocationFlag = "Cargo"
	AssetLocationFlagCorpseBay                           AssetLocationFlag = "CorpseBay"
	AssetLocationFlagDeliveries                          AssetLocationFlag = "Deliveries"
	AssetLocationFlagDroneBay                            AssetLocationFlag = "DroneBay"
	AssetLocationFlagFighterBay                          AssetLocationFlag = "FighterBay"
	AssetLocationFlagFighterTube0                        AssetLocationFlag = "FighterTube0"
	AssetLocationFlagFighterTube1                        AssetLocationFlag = "FighterTube1"
	AssetLocationFlagFighterTube2                        AssetLocationFlag = "FighterTube2"
	AssetLocationFlagFighterTube3                        AssetLocationFlag = "FighterTube3"
	AssetLocationFlagFighterTube4                        AssetLocationFlag = "FighterTube4"
	AssetLocationFlagFleetHangar                         AssetLocationFlag = "FleetHangar"
	AssetLocationFlagFrigateEscapeBay                    AssetLocationFlag = "FrigateEscapeBay"
	AssetLocationFlagHangar                              AssetLocationFlag = "Hangar"
	AssetLocationFlagHangarAll                           AssetLocationFlag = "HangarAll"
	AssetLocationFlagHiSlot0                             AssetLocationFlag = "HiSlot0"
	AssetLocationFlagHiSlot1                             AssetLocationFlag = "HiSlot1"
	AssetLocationFlagHiSlot2                             AssetLocationFlag = "HiSlot2"
	AssetLocationFlagHiSlot3                             AssetLocationFlag = "HiSlot3"
	AssetLocationFlagHiSlot4                             AssetLocationFlag = "HiSlot4"
	AssetLocationFlagHiSlot5                             AssetLocationFlag = "HiSlot5"
	AssetLocationFlagHiSlot6                             AssetLocationFlag = "HiSlot6"
	AssetLocationFlagHiSlot7                             AssetLocationFlag = "HiSlot7"
	AssetLocationFlagHiddenModifiers                     AssetLocationFlag = "HiddenModifiers"
	AssetLocationFlagImplant                             AssetLocationFlag = "Implant"
	AssetLocationFlagLoSlot0                             AssetLocationFlag = "LoSlot0"
	AssetLocationFlagLoSlot1                             AssetLocationFlag = "LoSlot1"
	AssetLocationFlagLoSlot2                             AssetLocationFlag = "LoSlot2"
	AssetLocationFlagLoSlot3                             AssetLocationFlag = "LoSlot3"
	AssetLocationFlagLoSlot4                             AssetLocationFlag = "LoSlot4"
	AssetLocationFlagLoSlot5                             AssetLocationFlag = "LoSlot5"
	AssetLocationFlagLoSlot6                             AssetLocationFlag = "LoSlot6"
	AssetLocationFlagLoSlot7                             AssetLocationFlag = "LoSlot7"
	AssetLocationFlagLocked                              AssetLocationFlag = "Locked"
	AssetLocationFlagMedSlot0                            AssetLocationFlag = "MedSlot0"
	AssetLocationFlagMedSlot1                            AssetLocationFlag = "MedSlot1"
	AssetLocationFlagMedSlot2                            AssetLocationFlag = "MedSlot2"
	AssetLocationFlagMedSlot3                            AssetLocationFlag = "MedSlot3"
	AssetLocationFlagMedSlot4                            AssetLocationFlag = "MedSlot4"
	AssetLocationFlagMedSlot5                            AssetLocationFlag = "MedSlot5"
	AssetLocationFlagMedSlot6                            AssetLocationFlag = "MedSlot6"
	AssetLocationFlagMedSlot7                            AssetLocationFlag = "MedSlot7"
	AssetLocationFlagQuafeBay                            AssetLocationFlag = "QuafeBay"
	AssetLocationFlagRigSlot0                            AssetLocationFlag = "RigSlot0"
	AssetLocationFlagRigSlot1                            AssetLocationFlag = "RigSlot1"
	AssetLocationFlagRigSlot2                            AssetLocationFlag = "RigSlot2"
	AssetLocationFlagRigSlot3                            AssetLocationFlag = "RigSlot3"
	AssetLocationFlagRigSlot4                            AssetLocationFlag = "RigSlot4"
	AssetLocationFlagRigSlot5                            AssetLocationFlag = "RigSlot5"
	AssetLocationFlagRigSlot6                            AssetLocationFlag = "RigSlot6"
	AssetLocationFlagRigSlot7                            AssetLocationFlag = "RigSlot7"
	AssetLocationFlagShipHangar                          AssetLocationFlag = "ShipHangar"
	AssetLocationFlagSkill                               AssetLocationFlag = "Skill"
	AssetLocationFlagSpecializedAmmoHold                 AssetLocationFlag = "SpecializedAmmoHold"
	AssetLocationFlagSpecializedCommandCenterHold        AssetLocationFlag = "SpecializedCommandCenterHold"
	AssetLocationFlagSpecializedFuelBay                  AssetLocationFlag = "SpecializedFuelBay"
	AssetLocationFlagSpecializedGasHold                  AssetLocationFlag = "SpecializedGasHold"
	AssetLocationFlagSpecializedIndustrialShipHold       AssetLocationFlag = "SpecializedIndustrialShipHold"
	AssetLocationFlagSpecializedLargeShipHold            AssetLocationFlag = "SpecializedLargeShipHold"
	AssetLocationFlagSpecializedMaterialBay              AssetLocationFlag = "SpecializedMaterialBay"
	AssetLocationFlagSpecializedMediumShipHold           AssetLocationFlag = "SpecializedMediumShipHold"
	AssetLocationFlagSpecializedMineralHold              AssetLocationFlag = "SpecializedMineralHold"
	AssetLocationFlagSpecializedOreHold                  AssetLocationFlag = "SpecializedOreHold"
	AssetLocationFlagSpecializedPlanetaryCommoditiesHold AssetLocationFlag = "SpecializedPlanetaryCommoditiesHold"
	AssetLocationFlagSpecializedSalvageHold              AssetLocationFlag = "SpecializedSalvageHold"
	AssetLocationFlagSpecializedShipHold                 AssetLocationFlag = "SpecializedShipHold"
	AssetLocationFlagSpecializedSmallShipHold            AssetLocationFlag = "SpecializedSmallShipHold"
	AssetLocationFlagSubSystemBay                        AssetLocationFlag = "SubSystemBay"
	AssetLocationFlagSubSystemSlot0                      AssetLocationFlag = "SubSystemSlot0"
	AssetLocationFlagSubSystemSlot1                      AssetLocationFlag = "SubSystemSlot1"
	AssetLocationFlagSubSystemSlot2                      AssetLocationFlag = "SubSystemSlot2"
	AssetLocationFlagSubSystemSlot3                      AssetLocationFlag = "SubSystemSlot3"
	AssetLocationFlagSubSystemSlot4                      AssetLocationFlag = "SubSystemSlot4"
	AssetLocationFlagSubSystemSlot5                      AssetLocationFlag = "SubSystemSlot5"
	AssetLocationFlagSubSystemSlot6                      AssetLocationFlag = "SubSystemSlot6"
	AssetLocationFlagSubSystemSlot7                      AssetLocationFlag = "SubSystemSlot7"
	AssetLocationFlagUnlocked                            AssetLocationFlag = "Unlocked"
	AssetLocationFlagWardrobe                            AssetLocationFlag = "Wardrobe"
)

var AllAssetLocationFlags = []AssetLocationFlag{
	AssetLocationFlagAssetSafety, AssetLocationFlagAutoFit, AssetLocationFlagBoosterBay, AssetLocationFlagCargo,
	AssetLocationFlagCorpseBay, AssetLocationFlagDeliveries, AssetLocationFlagDroneBay, AssetLocationFlagFighterBay,
	AssetLocationFlagFighterTube0, AssetLocationFlagFighterTube1, AssetLocationFlagFighterTube2, AssetLocationFlagFighterTube3,
	AssetLocationFlagFighterTube4, AssetLocationFlagFleetHangar, AssetLocationFlagFrigateEscapeBay, AssetLocationFlagHangar,
	AssetLocationFlagHangarAll, AssetLocationFlagHiSlot0, AssetLocationFlagHiSlot1, AssetLocationFlagHiSlot2,
	AssetLocationFlagHiSlot3, AssetLocationFlagHiSlot4, AssetLocationFlagHiSlot5, AssetLocationFlagHiSlot6,
	AssetLocationFlagHiSlot7, AssetLocationFlagHiddenModifiers, AssetLocationFlagImplant, AssetLocationFlagLoSlot0,
	AssetLocationFlagLoSlot1, AssetLocationFlagLoSlot2, AssetLocationFlagLoSlot3, AssetLocationFlagLoSlot4,
	AssetLocationFlagLoSlot5, AssetLocationFlagLoSlot6, AssetLocationFlagLoSlot7, AssetLocationFlagLocked,
	AssetLocationFlagMedSlot0, AssetLocationFlagMedSlot1, AssetLocationFlagMedSlot2, AssetLocationFlagMedSlot3,
	AssetLocationFlagMedSlot4, AssetLocationFlagMedSlot5, AssetLocationFlagMedSlot6, AssetLocationFlagMedSlot7,
	AssetLocationFlagQuafeBay, AssetLocationFlagRigSlot0, AssetLocationFlagRigSlot1, AssetLocationFlagRigSlot2,
	AssetLocationFlagRigSlot3, AssetLocationFlagRigSlot4, AssetLocationFlagRigSlot5, AssetLocationFlagRigSlot6,
	AssetLocationFlagRigSlot7, AssetLocationFlagShipHangar, AssetLocationFlagSkill, AssetLocationFlagSpecializedAmmoHold,
	AssetLocationFlagSpecializedCommandCenterHold, AssetLocationFlagSpecializedFuelBay, AssetLocationFlagSpecializedGasHold, AssetLocationFlagSpecializedIndustrialShipHold,
	AssetLocationFlagSpecializedLargeShipHold, AssetLocationFlagSpecializedMaterialBay, AssetLocationFlagSpecializedMediumShipHold, AssetLocationFlagSpecializedMineralHold,
	AssetLocationFlagSpecializedOreHold, AssetLocationFlagSpecializedPlanetaryCommoditiesHold, AssetLocationFlagSpecializedSalvageHold, AssetLocationFlagSpecializedShipHold,
	AssetLocationFlagSpecializedSmallShipHold, AssetLocationFlagSubSystemBay, AssetLocationFlagSubSystemSlot0, AssetLocationFlagSubSystemSlot1,
	AssetLocationFlagSubSystemSlot2, AssetLocationFlagSubSystemSlot3, AssetLocationFlagSubSystemSlot4, AssetLocationFlagSubSystemSlot5,
	AssetLocationFlagSubSystemSlot6, AssetLocationFlagSubSystemSlot7, AssetLocationFlagUnlocked, AssetLocationFlagWardrobe,
}

func (c AssetLocationFlag) Valid() bool {
	for _, v := range AllAssetLocationFlags {
		if v == c {
			return true
		}
	}

	return true
}

func (s AssetLocationFlag) String() string {
	return string(s)
}

type AssetLocationType string

const (
	AssetLocationTypeStation     AssetLocationType = "station"
	AssetLocationTypeSolarSystem AssetLocationType = "solar_system"
	AssetLocationTypeItem        AssetLocationType = "item"
	AssetLocationTypeOther       AssetLocationType = "other"
)

var AllAssetLocationTypes = []AssetLocationType{
	AssetLocationTypeStation,
	AssetLocationTypeSolarSystem,
	AssetLocationTypeItem,
	AssetLocationTypeOther,
}

func (c AssetLocationType) Valid() bool {
	for _, v := range AllAssetLocationTypes {
		if v == c {
			return true
		}
	}

	return true
}

func (s AssetLocationType) String() string {
	return string(s)
}
