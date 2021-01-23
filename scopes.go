package athena

import "context"

type ScopeMap map[Scope][]ScopeResolver

type ScopeResolver struct {
	Name string
	Func func(context.Context, *Member) error
}

type Scope string

const (
	ReadLocationV1 Scope = "esi-location.read_location.v1"
	ReadOnlineV1   Scope = "esi-location.read_online.v1"
	ReadShipV1     Scope = "esi-location.read_ship_type.v1"
	ReadClonesV1   Scope = "esi-clones.read_clones.v1"
	ReadImplants   Scope = "esi-clones.read_implants.v1"
)

func (s Scope) String() string {
	return string(s)
}
