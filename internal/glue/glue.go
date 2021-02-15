package glue

const (
	IDTypeCharacter   = "character"
	IDTypeCorporation = "corporation"
	IDTypeAlliance    = "alliance"
	IDTypeStation     = "station"
	IDTypeStructure   = "structure"
	IDTypeUnknown     = "unknown"
)

// Incomplete implementation of https://github.com/esi/esi-docs/blob/master/docs/id_ranges.md
// May add more to it later
func ResolveIDTypeFromIDRange(id uint64) string {

	switch d := id; {
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
	case d >= 1000000000000: // This unfortunately also includes POC's, so yeah....have I said that i hate CCP yet?
		return IDTypeStructure
	default:
		return IDTypeUnknown
	}

}
