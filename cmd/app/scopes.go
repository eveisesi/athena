package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/clone"
	"github.com/eveisesi/athena/internal/contact"
	"github.com/eveisesi/athena/internal/location"
	"github.com/eveisesi/athena/internal/skill"
)

func buildScopeMap(
	location location.Service,
	clone clone.Service,
	contact contact.Service,
	skill skill.Service,
) athena.ScopeMap {

	scopeMap := make(athena.ScopeMap, 10)
	scopeMap[athena.ReadLocationV1] = []athena.ScopeResolver{
		{
			Name: "MemberLocation",
			Func: location.EmptyMemberLocation,
		},
	}

	scopeMap[athena.ReadOnlineV1] = []athena.ScopeResolver{
		{
			Name: "MemberOnline",
			Func: location.EmptyMemberOnline,
		},
	}

	scopeMap[athena.ReadShipV1] = []athena.ScopeResolver{
		{
			Name: "MemberShip",
			Func: location.EmptyMemberShip,
		},
	}

	scopeMap[athena.ReadClonesV1] = []athena.ScopeResolver{
		{
			Name: "MemberClones",
			Func: clone.EmptyMemberClones,
		},
	}

	scopeMap[athena.ReadImplantsV1] = []athena.ScopeResolver{
		{
			Name: "MemberImplants",
			Func: clone.EmptyMemberImplants,
		},
	}

	scopeMap[athena.ReadContactsV1] = []athena.ScopeResolver{
		{
			Name: "MemberContacts",
			Func: contact.EmptyMemberContacts,
		},
		{
			Name: "MemberContactLabels",
			Func: contact.EmptyMemberContactLabels,
		},
	}

	scopeMap[athena.ReadSkillQueueV1] = []athena.ScopeResolver{
		{
			Name: "MemberSkillQueue",
			Func: skill.EmptyMemberSkillQueue,
		},
	}

	scopeMap[athena.ReadSkillsV1] = []athena.ScopeResolver{
		{
			Name: "MemberSkills",
			Func: skill.EmptyMemberSkills,
		},
	}

	return scopeMap

}

func listScopes() {
	scopes := athena.AllScopes
	list := make([]string, 0)
	for _, v := range scopes {
		list = append(list, strconv.Quote(v.String()))
	}

	fmt.Println(strings.Join(list, ","))
}
