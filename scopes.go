package athena

import "context"

type ScopeMap map[string]ScopeResolverFunc

type ScopeResolverFunc func(context.Context, *Member) error

type Scope string

// &scope=esi-location.read_location.v1+esi-location.read_online.v1

const (
	READ_LOCATION_V1 = "esi-location.read_location.v1"
	READ_ONLINE_V1   = "esi-location.read_online.v1"
	READ_SHIP_V1     = "esi-location.read_ship_type.v1"
)
