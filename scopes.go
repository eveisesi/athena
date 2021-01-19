package athena

import "context"

type ScopeMap map[Scope][]ScopeResolver

type ScopeResolver struct {
	Name string
	Func func(context.Context, *Member) error
}

type Scope string

// &scope=esi-location.read_location.v1+esi-location.read_online.v1

const (
	ReadLocationV1 Scope = "esi-location.read_location.v1"
	ReadOnlineV1   Scope = "esi-location.read_online.v1"
	ReadShipV1     Scope = "esi-location.read_ship_type.v1"
)

func (s Scope) String() string {
	return string(s)
}
