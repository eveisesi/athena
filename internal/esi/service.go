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
		GetAlliancesAllianceID(ctx context.Context, alliance *athena.Alliance) (*athena.Alliance, *http.Response, error)

		// Characters
		GetCharactersCharacterID(ctx context.Context, character *athena.Character) (*athena.Character, *http.Response, error)
		GetCharactersCharacterIDClones(ctx context.Context, member *athena.Member, clones *athena.MemberClones) (*athena.MemberClones, *http.Response, error)
		GetCharactersCharacterIDContacts(ctx context.Context, member *athena.Member, etag *athena.Etag, contacts []*athena.MemberContact) ([]*athena.MemberContact, *athena.Etag, *http.Response, error)
		GetCharactersCharacterIDContactLabels(ctx context.Context, member *athena.Member, etag *athena.Etag, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, *athena.Etag, *http.Response, error)
		GetCharactersCharacterIDImplants(ctx context.Context, member *athena.Member, implants *athena.MemberImplants) (*athena.MemberImplants, *http.Response, error)
		GetCharactersCharacterIDLocation(ctx context.Context, member *athena.Member, location *athena.MemberLocation) (*athena.MemberLocation, *http.Response, error)
		GetCharactersCharacterIDOnline(ctx context.Context, member *athena.Member, online *athena.MemberOnline) (*athena.MemberOnline, *http.Response, error)
		GetCharactersCharacterIDShip(ctx context.Context, member *athena.Member, ship *athena.MemberShip) (*athena.MemberShip, *http.Response, error)

		// Corporations
		GetCorporationsCorporationID(ctx context.Context, corporation *athena.Corporation) (*athena.Corporation, *http.Response, error)

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
		GetUniverseAncestries(ctx context.Context, ancestries []*athena.Ancestry) ([]*athena.Ancestry, *http.Response, error)
		GetUniverseBloodlines(ctx context.Context, bloodlines []*athena.Bloodline) ([]*athena.Bloodline, *http.Response, error)
		GetUniverseCategories(ctx context.Context, ids []int) ([]int, *http.Response, error)
		GetUniverseCategoriesCategoryID(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error)
		GetUniverseConstellationsConstellationID(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error)
		GetUniverseFactions(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error)
		GetUniverseGroupsGroupID(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error)
		GetUniverseRegions(ctx context.Context, ids []int) ([]int, *http.Response, error)
		GetUniverseRegionsRegionID(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error)
		GetUniverseRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error)
		GetUniverseSolarSystemsSolarSystemID(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error)
		GetUniverseStationsStationID(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error)
		GetUniverseStructuresStructureID(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error)
		GetUniverseTypesTypeID(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error)
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
	EndpointGetAlliancesAllianceID                   Endpoint = "GetAlliancesAllianceID"
	EndpointGetCharactersCharacterID                 Endpoint = "GetCharactersCharacterID"
	EndpointGetCharactersCharacterIDClones           Endpoint = "GetCharactersCharacterIDClones"
	EndpointGetCharactersCharacterIDContacts         Endpoint = "GetCharactersCharacterIDContacts"
	EndpointGetCharactersCharacterIDContactLabels    Endpoint = "GetCharactersCharacterIDContactLabels"
	EndpointGetCharactersCharacterIDImplants         Endpoint = "GetCharactersCharacterIDImplants"
	EndpointGetCharactersCharacterIDLocation         Endpoint = "GetCharactersCharacterIDLocation"
	EndpointGetCharactersCharacterIDOnline           Endpoint = "GetCharactersCharacterIDOnline"
	EndpointGetCharactersCharacterIDShip             Endpoint = "GetCharactersCharacterIDShip"
	EndpointGetCorporationsCorporationID             Endpoint = "GetCorporationsCorporationID"
	EndpointGetUniverseAncestries                    Endpoint = "GetUniverseAncestries"
	EndpointGetUniverseBloodlines                    Endpoint = "GetUniverseBloodlines"
	EndpointGetUniverseCategories                    Endpoint = "GetUniverseCategories"
	EndpointGetUniverseCategoriesCategoryID          Endpoint = "GetUniverseCategoriesCategoryID"
	EndpointGetUniverseConstellationsConstellationID Endpoint = "GetUniverseConstellationsConstellationID"
	EndpointGetUniverseFactions                      Endpoint = "GetUniverseFactions"
	EndpointGetUniverseGroupsGroupID                 Endpoint = "GetUniverseGroupsGroupID"
	EndpointGetUniverseRaces                         Endpoint = "GetUniverseRaces"
	EndpointGetUniverseRegions                       Endpoint = "GetUniverseRegions"
	EndpointGetUniverseRegionsRegionID               Endpoint = "GetUniverseRegionsRegionID"
	EndpointGetUniverseSolarSystemsSolarSystemID     Endpoint = "GetUniverseSolarSystemsSolarSystemID"
	EndpointGetUniverseStationsStationID             Endpoint = "GetUniverseStationsStationID"
	EndpointGetUniverseStructuresStructureID         Endpoint = "GetUniverseStructuresStructureID"
	EndpointGetUniverseTypesTypeID                   Endpoint = "GetUniverseTypesTypeID"
)

var AllEndpoints = []Endpoint{
	EndpointGetAlliancesAllianceID,
	EndpointGetCharactersCharacterID,
	EndpointGetCharactersCharacterIDClones,
	EndpointGetCharactersCharacterIDContacts,
	EndpointGetCharactersCharacterIDContactLabels,
	EndpointGetCharactersCharacterIDImplants,
	EndpointGetCharactersCharacterIDLocation,
	EndpointGetCharactersCharacterIDOnline,
	EndpointGetCharactersCharacterIDShip,
	EndpointGetCorporationsCorporationID,
	EndpointGetUniverseAncestries,
	EndpointGetUniverseBloodlines,
	EndpointGetUniverseCategories,
	EndpointGetUniverseCategoriesCategoryID,
	EndpointGetUniverseConstellationsConstellationID,
	EndpointGetUniverseFactions,
	EndpointGetUniverseGroupsGroupID,
	EndpointGetUniverseRaces,
	EndpointGetUniverseRegions,
	EndpointGetUniverseRegionsRegionID,
	EndpointGetUniverseSolarSystemsSolarSystemID,
	EndpointGetUniverseStationsStationID,
	EndpointGetUniverseStructuresStructureID,
	EndpointGetUniverseTypesTypeID,
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
		EndpointGetAlliancesAllianceID:                   s.resolveGetAlliancesAllianceIDEndpoint,
		EndpointGetCharactersCharacterID:                 s.resolveGetCharactersCharacterIDEndpoint,
		EndpointGetCharactersCharacterIDClones:           s.resolveGetCharactersCharacterIDClonesEndpoint,
		EndpointGetCharactersCharacterIDContacts:         s.resolveGetCharactersCharacterIDContactsEndpoint,
		EndpointGetCharactersCharacterIDContactLabels:    s.resolveGetCharactersCharacterIDContactLabelsEndpoint,
		EndpointGetCharactersCharacterIDImplants:         s.resolveGetCharactersCharacterIDImplantsEndpoint,
		EndpointGetCharactersCharacterIDLocation:         s.resolveGetCharactersCharacterIDLocationEndpoint,
		EndpointGetCharactersCharacterIDOnline:           s.resolveGetCharactersCharacterIDOnlineEndpoint,
		EndpointGetCharactersCharacterIDShip:             s.resolveGetCharactersCharacterIDShipEndpoint,
		EndpointGetCorporationsCorporationID:             s.resolveGetCorporationsCorporationIDEndpoint,
		EndpointGetUniverseAncestries:                    s.resolveUniverseAncestriesEndpoint,
		EndpointGetUniverseBloodlines:                    s.resolveUniverseBloodlinesEndpoint,
		EndpointGetUniverseCategories:                    s.resolveUniverseCategoriesEndpoint,
		EndpointGetUniverseCategoriesCategoryID:          s.resolveUniverseCategoriesCategoryIDEndpoint,
		EndpointGetUniverseConstellationsConstellationID: s.resolveGetUniverseConstellationsConstellationIDEndpoint,
		EndpointGetUniverseFactions:                      s.resolveUniverseFactionsEndpoint,
		EndpointGetUniverseGroupsGroupID:                 s.resolveUniverseGroupsGroupIDEndpoint,
		EndpointGetUniverseRaces:                         s.resolveUniverseRacesEndpoint,
		EndpointGetUniverseRegions:                       s.resolveGetUniverseRegionsEndpoint,
		EndpointGetUniverseRegionsRegionID:               s.resolveGetUniverseRegionsRegionIDEndpoint,
		EndpointGetUniverseSolarSystemsSolarSystemID:     s.resolveGetUniverseSolarSystemsSolarSystemIDEndpoint,
		EndpointGetUniverseStationsStationID:             s.resolveGetUniverseStationsStationIDEndpoint,
		EndpointGetUniverseStructuresStructureID:         s.resolveGetUniverseStructuresStructureIDEndpoint,
		EndpointGetUniverseTypesTypeID:                   s.resolveGetUniverseTypesTypeIDEndpoint,
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
