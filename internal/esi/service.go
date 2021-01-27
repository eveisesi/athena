package esi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

const (
	// ESI Timestamp Format
	ESI_EXPIRES_HEADER_FORMAT = "Mon, 02 Jan 2006 15:04:05 MST"
)

type (
	Service interface {
		Etag(ctx context.Context, endpoint Endpoint)

		GenerateEndpointHash(endpoint Endpoint, obj interface{}) (bool, string)

		// Alliances
		GetAlliance(ctx context.Context, alliance *athena.Alliance) (*athena.Alliance, *http.Response, error)

		// Characters
		GetCharacter(ctx context.Context, character *athena.Character) (*athena.Character, *http.Response, error)
		GetCharacterAttributes(ctx context.Context, member *athena.Member, attributes *athena.MemberSkillAttributes) (*athena.MemberSkillAttributes, *http.Response, error)
		GetCharacterClones(ctx context.Context, member *athena.Member, clones *athena.MemberClones) (*athena.MemberClones, *http.Response, error)
		GetCharacterContacts(ctx context.Context, member *athena.Member, etag *athena.Etag, contacts []*athena.MemberContact) ([]*athena.MemberContact, *athena.Etag, *http.Response, error)
		GetCharacterContactLabels(ctx context.Context, member *athena.Member, etag *athena.Etag, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error)
		GetCharacterImplants(ctx context.Context, member *athena.Member, implants *athena.MemberImplants) (*athena.MemberImplants, *http.Response, error)
		GetCharacterLocation(ctx context.Context, member *athena.Member, location *athena.MemberLocation) (*athena.MemberLocation, *http.Response, error)
		GetCharacterOnline(ctx context.Context, member *athena.Member, online *athena.MemberOnline) (*athena.MemberOnline, *http.Response, error)
		GetCharacterShip(ctx context.Context, member *athena.Member, ship *athena.MemberShip) (*athena.MemberShip, *http.Response, error)
		GetCharacterSkills(ctx context.Context, member *athena.Member, etag *athena.Etag, meta *athena.MemberSkillMeta) (*athena.MemberSkillMeta, *athena.Etag, *http.Response, error)
		GetCharacterSkillQueue(ctx context.Context, member *athena.Member, etag *athena.Etag, queue []*athena.MemberSkillQueue) ([]*athena.MemberSkillQueue, *athena.Etag, *http.Response, error)

		// Corporations
		GetCorporation(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, *http.Response, error)

		// 		// Killmails
		// 		GetKillmailsKillmailIDKillmailHash(ctx context.Context, id uint, hash string) (*athena.Killmail, *http.Response, error)

		// 		// Market
		// 		HeadMarketsRegionIDTypes(ctx context.Context, regionID uint) *http.Response, error
		// 		GetMarketGroups(ctx context.Context) ([]int, *http.Response, error)
		// 		GetMarketGroupsMarketGroupID(ctx context.Context, id int) (*athena.MarketGroup, *http.Response, error)
		// 		GetMarketsRegionIDTypes(ctx context.Context, regionID uint, page null.String) ([]int, *http.Response, error)
		// 		GetMarketsRegionIDHistory(ctx context.Context, regionID uint, typeID uint) ([]*athena.HistoricalRecord, *http.Response, error)
		// 		GetMarketsPrices(ctx context.Context) ([]*athena.MarketPrices, *http.Response, error)

		// 		// Status
		// 		GetStatus(ctx context.Context) (*athena.ServerStatus, *http.Response, error)

		// Universe
		GetAncestries(ctx context.Context, ancestries []*athena.Ancestry) ([]*athena.Ancestry, *http.Response, error)
		GetBloodlines(ctx context.Context, bloodlines []*athena.Bloodline) ([]*athena.Bloodline, *http.Response, error)
		GetCategories(ctx context.Context, ids []int) ([]int, *http.Response, error)
		GetCategory(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error)
		GetConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error)
		GetFaction(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error)
		GetGroup(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error)
		GetRegions(ctx context.Context, ids []int) ([]int, *http.Response, error)
		GetRegion(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error)
		GetRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error)
		GetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error)
		GetStation(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error)
		GetStructure(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error)
		GetType(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error)
	}

	service struct {
		client    *http.Client
		cache     cache.Service
		ua        string
		endpoints endpointMap
	}
)

// NewService returns a default implementation of this service
func NewService(cache cache.Service, client *http.Client, uagent string) Service {

	s := &service{
		cache:  cache,
		client: client,
		ua:     uagent,
	}

	s.buildEndpointMap()

	return s

}

type Endpoint string

const (
	EndpointGetAlliance               Endpoint = "GetAlliance"
	EndpointGetCharacter              Endpoint = "GetCharacter"
	EndpointGetCharacterAttributes    Endpoint = "GetCharacterAttributes"
	EndpointGetCharacterClones        Endpoint = "GetCharacterClones"
	EndpointGetCharacterContacts      Endpoint = "GetCharacterContacts"
	EndpointGetCharacterContactLabels Endpoint = "GetCharacterContactLabels"
	EndpointGetCharacterImplants      Endpoint = "GetCharacterImplants"
	EndpointGetCharacterLocation      Endpoint = "GetCharacterLocation"
	EndpointGetCharacterOnline        Endpoint = "GetCharacterOnline"
	EndpointGetCharacterShip          Endpoint = "GetCharacterShip"
	EndpointGetCharacterSkills        Endpoint = "GetCharacterSkills"
	EndpointGetCharacterSkillQueue    Endpoint = "GetCharacterSkillQueue"
	EndpointGetCorporation            Endpoint = "GetCorporation"
	EndpointGetAncestries             Endpoint = "GetAncestries"
	EndpointGetBloodlines             Endpoint = "GetBloodlines"
	EndpointGetCategories             Endpoint = "GetCategories"
	EndpointGetCategory               Endpoint = "GetCategory"
	EndpointGetConstellation          Endpoint = "GetConstellation"
	EndpointGetFaction                Endpoint = "GetFaction"
	EndpointGetGroup                  Endpoint = "GetGroup"
	EndpointGetRaces                  Endpoint = "GetRaces"
	EndpointGetRegions                Endpoint = "GetRegions"
	EndpointGetRegion                 Endpoint = "GetRegion"
	EndpointGetSolarSystem            Endpoint = "GetSolarSystem"
	EndpointGetStation                Endpoint = "GetStation"
	EndpointGetStructure              Endpoint = "GetStructure"
	EndpointGetType                   Endpoint = "GetType"
)

var AllEndpoints = []Endpoint{
	EndpointGetAlliance,
	EndpointGetCharacter,
	EndpointGetCharacterAttributes,
	EndpointGetCharacterClones,
	EndpointGetCharacterContacts,
	EndpointGetCharacterContactLabels,
	EndpointGetCharacterImplants,
	EndpointGetCharacterLocation,
	EndpointGetCharacterOnline,
	EndpointGetCharacterShip,
	EndpointGetCharacterSkills,
	EndpointGetCharacterSkillQueue,
	EndpointGetCorporation,
	EndpointGetAncestries,
	EndpointGetBloodlines,
	EndpointGetCategories,
	EndpointGetCategory,
	EndpointGetConstellation,
	EndpointGetFaction,
	EndpointGetGroup,
	EndpointGetRaces,
	EndpointGetRegions,
	EndpointGetRegion,
	EndpointGetSolarSystem,
	EndpointGetStation,
	EndpointGetStructure,
	EndpointGetType,
}

func (e Endpoint) String() string {
	return string(e)
}

func (e Endpoint) Valid() bool {
	for _, v := range AllEndpoints {
		if e == v {
			return true
		}
	}

	return false
}

func (s *service) buildEndpointMap() {

	s.endpoints = endpointMap{
		EndpointGetAlliance:               s.resolveGetAllianceEndpoint,
		EndpointGetCharacter:              s.resolveGetCharacterEndpoint,
		EndpointGetCharacterAttributes:    s.resolveGetCharacterAttributes,
		EndpointGetCharacterClones:        s.resolveGetCharacterClonesEndpoint,
		EndpointGetCharacterContacts:      s.resolveGetCharacterContactsEndpoint,
		EndpointGetCharacterContactLabels: s.resolveGetCharacterContactLabelsEndpoint,
		EndpointGetCharacterImplants:      s.resolveGetCharacterImplantsEndpoint,
		EndpointGetCharacterLocation:      s.resolveGetCharacterLocationEndpoint,
		EndpointGetCharacterOnline:        s.resolveGetCharacterOnlineEndpoint,
		EndpointGetCharacterShip:          s.resolveGetCharacterShipEndpoint,
		EndpointGetCharacterSkills:        s.resolveGetCharacterSkills,
		EndpointGetCharacterSkillQueue:    s.resolveGetCharacterSkillQueue,
		EndpointGetCorporation:            s.resolveGetCorporationEndpoint,
		EndpointGetAncestries:             s.resolveUniverseAncestriesEndpoint,
		EndpointGetBloodlines:             s.resolveUniverseBloodlinesEndpoint,
		EndpointGetCategories:             s.resolveUniverseCategoriesEndpoint,
		EndpointGetCategory:               s.resolveUniverseCategoriesCategoryIDEndpoint,
		EndpointGetConstellation:          s.resolveGetConstellationEndpoint,
		EndpointGetFaction:                s.resolveUniverseFactionsEndpoint,
		EndpointGetGroup:                  s.resolveUniverseGroupsGroupIDEndpoint,
		EndpointGetRaces:                  s.resolveUniverseRacesEndpoint,
		EndpointGetRegions:                s.resolveGetRegionsEndpoint,
		EndpointGetRegion:                 s.resolveGetRegionEndpoint,
		EndpointGetSolarSystem:            s.resolveGetSolarSystemEndpoint,
		EndpointGetStation:                s.resolveGetStationEndpoint,
		EndpointGetStructure:              s.resolveGetStructureEndpoint,
		EndpointGetType:                   s.resolveGetTypeEndpoint,
	}

}

func (s *service) GenerateEndpointHash(endpoint Endpoint, obj interface{}) (bool, string) {

	if !endpoint.Valid() {
		return false, ""
	}

	if _, ok := s.endpoints[endpoint]; !ok {
		return false, ""
	}

	alg := sha256.New()
	n, _ := alg.Write([]byte(s.endpoints[endpoint](obj)))

	return n > 0, fmt.Sprintf("%s::%x", endpoint.String(), alg.Sum(nil))

}

// Request prepares and executes an http request to the EVE Swagger Interface OpenAPI
// and returns the response
func (s *service) request(
	ctx context.Context,
	optionFuncs ...OptionFunc,
) ([]byte, *http.Response, error) {

	options := s.opts(optionFuncs)

	uri := url.URL{
		Scheme:   "https",
		Host:     "esi.evetech.net",
		Path:     options.path,
		RawQuery: options.query.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, options.method, uri.String(), bytes.NewBuffer(options.body))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header = options.headers

	req = newrelic.RequestWithTransactionContext(req, newrelic.FromContext(ctx))

	response, err := s._exec(req, options)
	if err != nil {
		return nil, response, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrap(err, "error reading body")

		return nil, response, err
	}

	response.Body.Close()

	s.trackESICallStatusCode(ctx, response.StatusCode)

	s.retrieveErrorReset(ctx, response.Header)
	s.retrieveErrorCount(ctx, response.Header)

	return data, response, nil
}

func (s *service) _exec(req *http.Request, options *options) (response *http.Response, err error) {

	for i := 0; i < options.maxattempts; i++ {
		response, err = s.client.Do(req)
		if err != nil && !options.retryOnError {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}

		if response.StatusCode > http.StatusContinue && response.StatusCode < http.StatusInternalServerError {
			break
		}

		if !options.retryOnError {
			break
		}
	}

	return response, err

}

// retrieveExpiresHeader takes a map[string]string of the response headers, checks to see if the "Expires" key exists, and if it does, parses the timestamp and returns a time.Time. If duraction
// is greater than zero(0), then that number of minutes will be add to the expires time that is parsed from the header.
func (s *service) retrieveExpiresHeader(h http.Header, duration time.Duration) time.Time {

	oneHour := time.Now().Add(time.Minute * 60)

	header := h.Get("expires")
	if header == "" {
		return oneHour
	}

	expires, err := time.Parse(ESI_EXPIRES_HEADER_FORMAT, header)
	if err != nil {
		return oneHour
	}

	if duration > 0 {
		expires = expires.Add(duration)
	}

	return expires
}

// retrieveEtagHeader is a helper method that retrieves an Etag for the most recent request to
// ESI
func (s *service) retrieveEtagHeader(h http.Header) string {
	return h.Get("Etag")
}

// retrieveErrorCount is a helper method that retrieves the number of errors that this application
// has triggered and how many more can be triggered before potentially encountereding an HTTP Status 420
func (s *service) retrieveErrorCount(ctx context.Context, h http.Header) {
	// Default to a low count. This will cause the app to slow down
	// if the header is not present to set the actual value from the header
	var count int64 = 15
	strCount := h.Get("x-esi-error-limit-remain")
	if strCount != "" {
		i, err := strconv.ParseInt(strCount, 10, 64)
		if err == nil {
			count = i
		}
	}

	s.cache.SetESIErrCount(ctx, count)

}

// retrieveErrorReset is a helper method that retrieves the number of seconds until our Error Limit resets
func (s *service) retrieveErrorReset(ctx context.Context, h http.Header) {
	reset := h.Get("x-esi-error-limit-reset")
	if reset == "" {
		return
	}

	seconds, err := strconv.ParseUint(reset, 10, 32)
	if err != nil {
		return
	}

	s.cache.SetEsiErrorReset(ctx, time.Now().Add(time.Second*time.Duration(seconds)).Unix())
}

func (s *service) trackESICallStatusCode(ctx context.Context, code int) {
	s.cache.SetESITracking(ctx, code, time.Now().UnixNano())
}
