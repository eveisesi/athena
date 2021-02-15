package glue

const (
	IDTypeCharacter   = "character"
	IDTypeCorporation = "corporation"
	IDTypeAlliance    = "alliance"
	IDTypeSolarSystem = "solar_system"
	IDTypeStation     = "station"
	IDTypeStructure   = "structure"
	IDTypeUnknown     = "unknown"
)

// Incomplete implementation of https://github.com/esi/esi-docs/blob/master/docs/id_ranges.md
// May add more to it later
func ResolveIDTypeFromIDRange(id uint64) string {

	switch d := id; {
	case d >= 30000000 && d < 32000000:
		return IDTypeSolarSystem
	case d >= 60000000 && d < 64000000:
		return IDTypeStation
	case d >= 90000000 && d < 98000000:
		return IDTypeCharacter
	case d >= 98000000 && d < 99000000:
		return IDTypeCorporation
	case d >= 99000000 && d < 100000000:
		return IDTypeAlliance
	case d >= 100000000 && d < 2100000000: // I hate you CCP, this is BS, why did you do this.....
		return IDTypeUnknown
	case d >= 2100000000 && d < 1000000000000:
		return IDTypeCharacter
	case d >= 1000000000000 && d < 1020000000000: // This should be POC's and whatever else can be anchored in space
		return IDTypeUnknown
	case d >= 1020000000000: // This should be Upwell Structures, but I've been wrong before
		return IDTypeStructure
	default:
		return IDTypeUnknown
	}

}
