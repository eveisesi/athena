package athena

import "context"

type ScopeMap map[string][]ScopeResolver

type ScopeResolver struct {
	Name string
	Func func(context.Context, *Member) error
}

const (
	ReadLocationV1 = "esi-location.read_location.v1"
	ReadOnlineV1   = "esi-location.read_online.v1"
	ReadShipV1     = "esi-location.read_ship_type.v1"
	ReadClonesV1   = "esi-clones.read_clones.v1"
	ReadImplants   = "esi-clones.read_implants.v1"
)
