package esi

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/athena"
	"github.com/eveisesi/athena/internal/cache"
	"github.com/eveisesi/athena/internal/etag"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

const (
	// ESI Timestamp Format
	ESI_EXPIRES_HEADER_FORMAT = "Mon, 02 Jan 2006 15:04:05 MST"
)

type (
	Service interface {
		etagInterface
		characterInterface
		clonesInterface
		contactsInterface
		locationInterface
		skillsInterface

		// Alliances
		GetAlliance(ctx context.Context, alliance *athena.Alliance) (*athena.Alliance, *http.Response, error)

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
		GetCategories(ctx context.Context, ids []uint) ([]uint, *http.Response, error)
		GetCategory(ctx context.Context, category *athena.Category) (*athena.Category, *http.Response, error)
		GetConstellation(ctx context.Context, constellation *athena.Constellation) (*athena.Constellation, *http.Response, error)
		GetFactions(ctx context.Context, factions []*athena.Faction) ([]*athena.Faction, *http.Response, error)
		GetGroup(ctx context.Context, group *athena.Group) (*athena.Group, *http.Response, error)
		GetRegions(ctx context.Context, ids []uint) ([]uint, *http.Response, error)
		GetRegion(ctx context.Context, region *athena.Region) (*athena.Region, *http.Response, error)
		GetRaces(ctx context.Context, races []*athena.Race) ([]*athena.Race, *http.Response, error)
		GetSolarSystem(ctx context.Context, solarSystem *athena.SolarSystem) (*athena.SolarSystem, *http.Response, error)
		GetStation(ctx context.Context, station *athena.Station) (*athena.Station, *http.Response, error)
		GetStructure(ctx context.Context, member *athena.Member, structure *athena.Structure) (*athena.Structure, *http.Response, error)
		GetType(ctx context.Context, item *athena.Type) (*athena.Type, *http.Response, error)
	}

	service struct {
		client *http.Client

		cache cache.Service
		etag  etag.Service

		ua        string
		endpoints endpointMap
	}
)

// NewService returns a default implementation of this service
func NewService(client *http.Client, cache cache.Service, etag etag.Service, uagent string) Service {

	s := &service{
		client: client,

		cache: cache,
		etag:  etag,

		ua: uagent,
	}

	s.buildEndpointMap()

	return s

}

type endpoint struct {
	Name     string
	FmtPath  string
	PathFunc pathFunc
	KeyFunc  keyFunc
}

var (
	// Alliances
	GetAlliance = &endpoint{Name: "GetAlliance", FmtPath: "/v3/alliances/%d/"}

	// Characters
	GetCharacter                   = &endpoint{Name: "GetCharacter", FmtPath: "/v4/characters/%d/"}
	GetCharacterCorporationHistory = &endpoint{Name: "GetCharacterCorporationHistory", FmtPath: "/v1/characters/%d/corporationhistory/"}

	// Skills
	GetCharacterAttributes = &endpoint{Name: "GetCharacterAttributes", FmtPath: "/v1/characters/%d/attributes/"}
	GetCharacterSkills     = &endpoint{Name: "GetCharacterSkills", FmtPath: "/v4/characters/%d/skills/"}
	GetCharacterSkillQueue = &endpoint{Name: "GetCharacterSkillQueue", FmtPath: "/v2/characters/%d/skillqueue/"}

	// Clones
	GetCharacterClones   = &endpoint{Name: "GetCharacterClones", FmtPath: "/v4/characters/%d/clones/"}
	GetCharacterImplants = &endpoint{Name: "GetCharacterImplants", FmtPath: "/v2/characters/%d/implants/"}

	// Contacts
	GetCharacterContacts      = &endpoint{Name: "GetCharacterContacts", FmtPath: "/v2/characters/%d/contacts/"}
	GetCharacterContactLabels = &endpoint{Name: "GetCharacterContactLabels", FmtPath: "/v1/characters/%d/contacts/labels/"}

	// Contracts
	GetCharacterContracts     = &endpoint{Name: "GetCharacterContracts", FmtPath: "/v1/characters/%d/contracts/"}
	GetCharacterContractItems = &endpoint{Name: "GetCharacterContractItems", FmtPath: "/v1/characters/%d/contracts/%d/items/"}
	GetCharacterContractBids  = &endpoint{Name: "GetCharacterContractBids", FmtPath: "/v1/characters/%d/contracts/%d/bids/"}

	// Fittings
	GetCharacterFittings = &endpoint{Name: "GetCharacterFittings", FmtPath: "/v2/characters/%d/fittings/"}

	// Locations
	GetCharacterLocation = &endpoint{Name: "GetCharacterLocation", FmtPath: "/v2/characters/%d/location/"}
	GetCharacterOnline   = &endpoint{Name: "GetCharacterOnline", FmtPath: "/v3/characters/%d/online/"}
	GetCharacterShip     = &endpoint{Name: "GetCharacterShip", FmtPath: "/v2/characters/%d/ship/"}

	// Corporations
	GetCorporation                = &endpoint{Name: "GetCorporation", FmtPath: "/v4/corporations/%d/"}
	GetCorporationAllianceHistory = &endpoint{Name: "GetCorporationAllianceHistory", FmtPath: "/v2/corporations/%d/alliancehistory/"}

	// Universe
	GetAncestries    = &endpoint{Name: "GetAncestries", FmtPath: "/v1/universe/ancestries/"}
	GetBloodlines    = &endpoint{Name: "GetBloodlines", FmtPath: "/v1/universe/bloodlines/"}
	GetCategories    = &endpoint{Name: "GetCategories", FmtPath: "/v1/universe/categories/"}
	GetCategory      = &endpoint{Name: "GetCategory", FmtPath: "/v1/universe/categories/%d/"}
	GetConstellation = &endpoint{Name: "GetConstellation", FmtPath: "/v1/universe/constellations/%d/"}
	GetFactions      = &endpoint{Name: "GetFactions", FmtPath: "/v2/universe/factions/"}
	GetGroup         = &endpoint{Name: "GetGroup", FmtPath: "/v1/universe/groups/%d/"}
	GetRaces         = &endpoint{Name: "GetRaces", FmtPath: "/v1/universe/races/"}
	GetRegions       = &endpoint{Name: "GetRegions", FmtPath: "/v1/universe/regions/"}
	GetRegion        = &endpoint{Name: "GetRegion", FmtPath: "/v1/universe/regions/%d/"}
	GetSolarSystem   = &endpoint{Name: "GetSolarSystem", FmtPath: "/v4/universe/systems/%d/"}
	GetStation       = &endpoint{Name: "GetStation", FmtPath: "/v2/universe/stations/%d/"}
	GetStructure     = &endpoint{Name: "GetStructure", FmtPath: "/v2/universe/structures/%d/"}
	GetType          = &endpoint{Name: "GetType", FmtPath: "/v3/universe/types/%d/"}
)

var AllEndpoints = []*endpoint{
	GetAlliance,
	GetCharacter,
	GetCharacterCorporationHistory,
	GetCharacterAttributes,
	GetCharacterClones,
	GetCharacterContacts,
	GetCharacterContactLabels,
	GetCharacterContracts,
	GetCharacterContractItems,
	GetCharacterContractBids,
	GetCharacterFittings,
	GetCharacterImplants,
	GetCharacterLocation,
	GetCharacterOnline,
	GetCharacterShip,
	GetCharacterSkills,
	GetCharacterSkillQueue,
	GetCorporation,
	GetCorporationAllianceHistory,
	GetAncestries,
	GetBloodlines,
	GetCategories,
	GetCategory,
	GetConstellation,
	GetFactions,
	GetGroup,
	GetRaces,
	GetRegions,
	GetRegion,
	GetSolarSystem,
	GetStation,
	GetStructure,
	GetType,
}

func (s *service) buildEndpointMap() {

	s.endpoints = endpointMap{
		GetAlliance.Name:                    s.newGetAllianceEndpoint(),
		GetCharacter.Name:                   s.newGetCharacterEndpoint(),
		GetCharacterCorporationHistory.Name: s.newGetCharacterCorporationHistoryEndpoint(),
		GetCharacterAttributes.Name:         s.newGetCharacterAttributesEndpoint(),
		GetCharacterClones.Name:             s.newGetCharacterClonesEndpoint(),
		GetCharacterContacts.Name:           s.newGetCharacterContactsEndpoint(),
		GetCharacterContactLabels.Name:      s.newGetCharacterContactLabelsEndpoint(),
		GetCharacterContracts.Name:          s.newGetCharacterContractsEndpoint(),
		GetCharacterContractItems.Name:      s.newGetCharacterContractItemsEndpoint(),
		GetCharacterContractBids.Name:       s.newGetCharacterContractBidsEndpoint(),
		GetCharacterFittings.Name:           s.newGetCharacterFittingsEndpoint(),
		GetCharacterImplants.Name:           s.newGetCharacterImplantsEndpoint(),
		GetCharacterLocation.Name:           s.newGetCharacterLocationEndpoint(),
		GetCharacterOnline.Name:             s.newGetCharacterOnlineEndpoint(),
		GetCharacterShip.Name:               s.newGetCharacterShipEndpoint(),
		GetCharacterSkills.Name:             s.newGetCharacterSkillsEndpoint(),
		GetCharacterSkillQueue.Name:         s.newGetCharacterSkillQueueEndpoint(),
		GetCorporation.Name:                 s.newGetCorporationEndpoint(),
		GetCorporationAllianceHistory.Name:  s.newGetCorporationAllianceHistoryEndpoint(),
		GetAncestries.Name:                  s.newGetAncestriesEndpoint(),
		GetBloodlines.Name:                  s.newGetBloodlinesEndpoint(),
		GetCategories.Name:                  s.newGetCategoriesEndpoint(),
		GetCategory.Name:                    s.newGetCategoryEndpoint(),
		GetConstellation.Name:               s.newGetConstellationEndpoint(),
		GetFactions.Name:                    s.newGetFactionsEndpoint(),
		GetGroup.Name:                       s.newGetGroupEndpoint(),
		GetRaces.Name:                       s.newGetRacesEndpoint(),
		GetRegions.Name:                     s.newGetRegionsEndpoint(),
		GetRegion.Name:                      s.newGetRegionEndpoint(),
		GetSolarSystem.Name:                 s.newGetSolarSystemEndpoint(),
		GetStation.Name:                     s.newGetStationEndpoint(),
		GetStructure.Name:                   s.newGetStructureEndpoint(),
		GetType.Name:                        s.newGetTypeEndpoint(),
	}

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
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}

		if response.StatusCode > http.StatusContinue && response.StatusCode < http.StatusInternalServerError && !options.retryOnError {
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

func (s *service) retrieveXPagesFromHeader(h http.Header) int {

	header := h.Get("X-Pages")
	if header == "" {
		return 0
	}

	pages, err := strconv.Atoi(header)
	if err != nil {
		return 0
	}

	return pages

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

func buildKey(s ...string) string {
	return strings.Join(s, "::")
}
