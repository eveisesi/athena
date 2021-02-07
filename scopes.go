package athena

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/volatiletech/null"
)

type ScopeMap map[Scope][]ScopeResolver

type ScopeResolver struct {
	Name string
	Func func(context.Context, *Member) (*Etag, error)
}

type Scope string

const (
	ReadLocationV1   Scope = "esi-location.read_location.v1"
	ReadOnlineV1     Scope = "esi-location.read_online.v1"
	ReadShipV1       Scope = "esi-location.read_ship_type.v1"
	ReadClonesV1     Scope = "esi-clones.read_clones.v1"
	ReadImplantsV1   Scope = "esi-clones.read_implants.v1"
	ReadContactsV1   Scope = "esi-characters.read_contacts.v1"
	ReadSkillQueueV1 Scope = "esi-skills.read_skillqueue.v1"
	ReadSkillsV1     Scope = "esi-skills.read_skills.v1"
)

var AllScopes = []Scope{
	ReadLocationV1,
	ReadOnlineV1,
	ReadShipV1,
	ReadClonesV1,
	ReadImplantsV1,
	ReadContactsV1,
	ReadSkillQueueV1,
	ReadSkillsV1,
}

func (s Scope) String() string {
	return string(s)
}

type MemberScope struct {
	Scope  Scope     `db:"scope" json:"scope"`
	Expiry null.Time `db:"expiry,omitempty" json:"expiry,omitempty"`
}

type MemberScopes []MemberScope

func (s *MemberScopes) Scan(value interface{}) error {

	switch data := value.(type) {
	case []byte:
		var scopes MemberScopes
		err := json.Unmarshal(data, &scopes)
		if err != nil {
			return err
		}

		*s = scopes
	}

	return nil
}

func (s MemberScopes) Value() (driver.Value, error) {
	var data []byte
	var err error
	if len(s) == 0 {
		data, err = json.Marshal([]interface{}{})
	} else {
		data, err = json.Marshal(s)
	}
	if err != nil {
		return nil, fmt.Errorf("[MemberScopes] Failed to marshal scope for storage in data store: %w", err)
	}

	return data, nil
}
