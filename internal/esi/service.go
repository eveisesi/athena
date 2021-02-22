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
		assetInterface
		allianceInterface
		characterInterface
		clonesInterface
		contactsInterface
		contractInterface
		corporationInterface
		etagInterface
		fittingInterface
		locationInterface
		mailInterface
		skillsInterface
		walletInterface
		universeInterface
	}

	service struct {
		client *http.Client

		cache cache.Service
		etag  etag.Service

		ua string
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

type (
	endpoint struct {
		Path     string
		PathFunc pathFunc
		KeyFunc  keyFunc
	}
	endpointID string
)

func (e endpointID) String() string {
	return string(e)
}

const (
	GetAlliance                    endpointID = "GetAlliance"
	GetCharacter                   endpointID = "GetCharacter"
	GetCharacterCorporationHistory endpointID = "GetCharacterCorporationHistory"
	GetCharacterAttributes         endpointID = "GetCharacterAttributes"
	GetCharacterAssets             endpointID = "GetCharacterAssets"
	GetCharacterSkills             endpointID = "GetCharacterSkills"
	GetCharacterSkillQueue         endpointID = "GetCharacterSkillQueue"
	GetCharacterClones             endpointID = "GetCharacterClones"
	GetCharacterImplants           endpointID = "GetCharacterImplants"
	GetCharacterContacts           endpointID = "GetCharacterContacts"
	GetCharacterContactLabels      endpointID = "GetCharacterContactLabels"
	GetCharacterContracts          endpointID = "GetCharacterContracts"
	GetCharacterContractItems      endpointID = "GetCharacterContractItems"
	GetCharacterContractBids       endpointID = "GetCharacterContractBids"
	GetCharacterFittings           endpointID = "GetCharacterFittings"
	GetCharacterLocation           endpointID = "GetCharacterLocation"
	GetCharacterMailHeaders        endpointID = "GetCharacterMailHeaders"
	GetCharacterMailHeader         endpointID = "GetCharacterMailHeader"
	GetCharacterMailLabels         endpointID = "GetCharacterMailLabels"
	GetCharacterMailLists          endpointID = "GetCharacterMailLists"
	GetCharacterOnline             endpointID = "GetCharacterOnline"
	GetCharacterShip               endpointID = "GetCharacterShip"
	GetCharacterWalletBalance      endpointID = "GetCharacterWalletBalance"
	GetCharacterWalletTransactions endpointID = "GetCharacterWalletTransactions"
	GetCharacterWalletJournal      endpointID = "GetCharacterWalletJournal"
	GetCorporation                 endpointID = "GetCorporation"
	GetCorporationAllianceHistory  endpointID = "GetCorporationAllianceHistory"
	GetAncestries                  endpointID = "GetAncestries"
	GetAsteroidBelt                endpointID = "GetAsteroidBelt"
	GetBloodlines                  endpointID = "GetBloodlines"
	GetCategories                  endpointID = "GetCategories"
	GetCategory                    endpointID = "GetCategory"
	GetConstellation               endpointID = "GetConstellation"
	GetFactions                    endpointID = "GetFactions"
	GetGroup                       endpointID = "GetGroup"
	GetMoon                        endpointID = "GetMoon"
	GetPlanet                      endpointID = "GetPlanet"
	GetRaces                       endpointID = "GetRaces"
	GetRegions                     endpointID = "GetRegions"
	GetRegion                      endpointID = "GetRegion"
	GetSolarSystem                 endpointID = "GetSolarSystem"
	GetStation                     endpointID = "GetStation"
	GetStructure                   endpointID = "GetStructure"
	GetType                        endpointID = "GetType"
	PostUniverseNames              endpointID = "PostUniverseNames"
)

var (
	endpoints endpointMap
)

func (s *service) buildEndpointMap() {

	endpoints = endpointMap{
		GetAlliance: &endpoint{
			Path:     "/v3/alliances/%d/",
			KeyFunc:  allianceKeyFunc,
			PathFunc: alliancePathFunc,
		},
		GetCharacter: &endpoint{
			Path:     "/v4/characters/%d/",
			KeyFunc:  characterKeyFunc,
			PathFunc: characterPathFunc,
		},
		GetCharacterCorporationHistory: &endpoint{
			Path:     "/v1/characters/%d/corporationhistory/",
			PathFunc: characterCorporationHistoryPathFunc,
			KeyFunc:  characterCorporationHistoryKeyFunc,
		},
		GetCharacterAssets: &endpoint{
			Path:     "/v5/characters/%d/assets/",
			PathFunc: characterAssetsPathFunc,
			KeyFunc:  characterAssetsKeyFunc,
		},
		GetCharacterAttributes: &endpoint{
			Path:     "/v1/characters/%d/attributes/",
			KeyFunc:  characterAttributesKeyFunc,
			PathFunc: characterAttributesPathFunc,
		},
		GetCharacterClones: &endpoint{
			Path:     "/v4/characters/%d/clones/",
			KeyFunc:  characterClonesKeyFunc,
			PathFunc: characterClonesPathFunc,
		},
		GetCharacterContacts: &endpoint{
			Path:     "/v2/characters/%d/contacts/",
			KeyFunc:  characterContactsKeyFunc,
			PathFunc: characterContactsPathFunc,
		},
		GetCharacterContactLabels: &endpoint{
			Path:     "/v1/characters/%d/contacts/labels/",
			KeyFunc:  characterContactLabelsKeyFunc,
			PathFunc: characterContactLabelsPathFunc,
		},
		GetCharacterContracts: &endpoint{
			Path:     "/v1/characters/%d/contracts/",
			KeyFunc:  characterContractsKeyFunc,
			PathFunc: characterContractsPathFunc,
		},
		GetCharacterContractItems: &endpoint{
			Path:     "/v1/characters/%d/contracts/%d/items/",
			KeyFunc:  characterContractItemsKeyFunc,
			PathFunc: characterContractItemsPathFunc,
		},
		GetCharacterContractBids: &endpoint{
			Path:     "/v1/characters/%d/contracts/%d/bids/",
			KeyFunc:  characterContractBidsKeyFunc,
			PathFunc: characterContractBidsPathFunc,
		},
		GetCharacterFittings: &endpoint{
			Path:     "/v2/characters/%d/fittings/",
			KeyFunc:  characterFittingsKeyFunc,
			PathFunc: characterFittingsPathFunc,
		},
		GetCharacterImplants: &endpoint{
			Path:     "/v2/characters/%d/implants/",
			KeyFunc:  characterImplantsKeyFunc,
			PathFunc: characterImplantsPathFunc,
		},
		GetCharacterMailHeaders: &endpoint{
			Path:     "/v1/characters/%d/mail/",
			KeyFunc:  characterMailsKeyFunc,
			PathFunc: characterMailsPathFunc,
		},
		GetCharacterMailHeader: &endpoint{
			Path:     "/v1/characters/%d/mail/%d/",
			KeyFunc:  characterMailKeyFunc,
			PathFunc: characterMailPathFunc,
		},
		GetCharacterMailLists: &endpoint{
			Path:     "/v1/characters/%d/mail/lists/",
			KeyFunc:  characterMailListsKeyFunc,
			PathFunc: characterMailListsPathFunc,
		},
		GetCharacterMailLabels: &endpoint{
			Path:     "/v3/characters/%d/mail/labels/",
			KeyFunc:  characterMailLabelsKeyFunc,
			PathFunc: characterMailLabelsPathFunc,
		},
		GetCharacterLocation: &endpoint{
			Path:     "/v2/characters/%d/location/",
			KeyFunc:  characterLocationsKeyFunc,
			PathFunc: characterLocationsPathFunc,
		},
		GetCharacterOnline: &endpoint{
			Path:     "/v3/characters/%d/online/",
			KeyFunc:  characterOnlinesKeyFunc,
			PathFunc: characterOnlinesPathFunc,
		},
		GetCharacterShip: &endpoint{
			Path:     "/v2/characters/%d/ship/",
			KeyFunc:  characterShipsKeyFunc,
			PathFunc: characterShipsPathFunc,
		},
		GetCharacterSkills: &endpoint{
			Path:     "/v4/characters/%d/skills/",
			KeyFunc:  characterSkillsKeyFunc,
			PathFunc: characterSkillsPathFunc,
		},
		GetCharacterSkillQueue: &endpoint{
			Path:     "/v2/characters/%d/skillqueue/",
			KeyFunc:  characterSkillQueueKeyFunc,
			PathFunc: characterSkillQueuePathFunc,
		},
		GetCharacterWalletBalance: &endpoint{
			Path:     "/v1/characters/%d/wallet/",
			PathFunc: characterWalletBalancePathFunc,
			KeyFunc:  characterWalletBalanceKeyFunc,
		},
		GetCharacterWalletTransactions: &endpoint{
			Path:     "/v1/characters/%d/wallet/transactions",
			KeyFunc:  characterWalletTransactionsKeyFunc,
			PathFunc: characterWalletTransactionsPathFunc,
		},
		GetCharacterWalletJournal: &endpoint{
			Path:     "/v6/characters/%d/wallet/journal/",
			KeyFunc:  characterWalletJournalKeyFunc,
			PathFunc: characterWalletJournalPathFunc,
		},

		// Corporations
		GetCorporation: &endpoint{
			Path:     "/v4/corporations/%d/",
			KeyFunc:  corporationKeyFunc,
			PathFunc: corporationPathFunc,
		},
		GetCorporationAllianceHistory: &endpoint{
			Path:     "/v2/corporations/%d/alliancehistory/",
			KeyFunc:  corporationAllianceHistoryKeyFunc,
			PathFunc: corporationAllianceHistoryPathFunc,
		},

		// Universe
		GetAncestries: &endpoint{
			Path:     "/v1/universe/ancestries/",
			KeyFunc:  ancestriesKeyFunc,
			PathFunc: ancestriesPathFunc,
		},
		GetBloodlines: &endpoint{
			Path:     "/v1/universe/bloodlines/",
			KeyFunc:  bloodlinesKeyFunc,
			PathFunc: bloodlinesPathFunc,
		},
		GetCategories: &endpoint{
			Path:     "/v1/universe/categories/",
			KeyFunc:  categoriesKeyFunc,
			PathFunc: categoriesPathFunc,
		},
		GetCategory: &endpoint{
			Path:     "/v1/universe/categories/%d/",
			KeyFunc:  categoryKeyFunc,
			PathFunc: categoryPathFunc,
		},
		GetConstellation: &endpoint{
			Path:     "/v1/universe/constellations/%d/",
			KeyFunc:  constellationKeyFunc,
			PathFunc: constellationPathFunc,
		},
		GetFactions: &endpoint{
			Path:     "/v2/universe/factions/",
			KeyFunc:  factionsKeyFunc,
			PathFunc: factionsPathFunc,
		},
		GetGroup: &endpoint{
			Path:     "/v1/universe/groups/%d/",
			KeyFunc:  groupKeyFunc,
			PathFunc: groupPathFunc,
		},
		GetRaces: &endpoint{
			Path:     "/v1/universe/races/",
			KeyFunc:  racesKeyFunc,
			PathFunc: racesPathFunc,
		},
		GetRegions: &endpoint{
			Path:     "/v1/universe/regions/",
			KeyFunc:  regionsKeyFunc,
			PathFunc: regionsPathFunc,
		},
		GetRegion: &endpoint{
			Path:     "/v1/universe/regions/%d/",
			KeyFunc:  regionKeyFunc,
			PathFunc: regionPathFunc,
		},
		GetSolarSystem: &endpoint{
			Path:     "/v4/universe/systems/%d/",
			KeyFunc:  solarSystemKeyFunc,
			PathFunc: solarSystemPathFunc,
		},
		GetStation: &endpoint{
			Path:     "/v2/universe/stations/%d/",
			KeyFunc:  stationKeyFunc,
			PathFunc: stationPathFunc,
		},
		GetStructure: &endpoint{
			Path:     "/v2/universe/structures/%d/",
			KeyFunc:  structureKeyFunc,
			PathFunc: structurePathFunc,
		},
		GetType: &endpoint{
			Path:     "/v3/universe/types/%d/",
			KeyFunc:  typeKeyFunc,
			PathFunc: typePathFunc,
		},

		// Etags are not cached on non Get endpoints, hence the lack of the KeyFunc
		PostUniverseNames: &endpoint{
			Path:     "/v3/universe/names/",
			PathFunc: func(_ *modifiers) string { return "/v3/universe/names/" },
		},
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

	reset := RetrieveErrorReset(response.Header)
	s.storeESIErrorReset(ctx, reset)
	count := RetrieveErrorCount(response.Header)
	s.storeESIErrorCount(ctx, count)

	if count <= 20 {
		time.Sleep(time.Second * time.Duration(reset))
	}

	return data, response, nil
}

func (s *service) _exec(req *http.Request, options *options) (response *http.Response, err error) {

	for i := 0; i < options.maxattempts; i++ {
		response, err = s.client.Do(req)
		if err != nil && (!options.retryOnError || i != options.maxattempts-1) {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}

		if response != nil && response.StatusCode > http.StatusContinue && response.StatusCode < http.StatusInternalServerError {
			break
		}

		time.Sleep(time.Second)

	}

	return response, err

}

// retrieveExpiresHeader takes a map[string]string of the response headers, checks to see if the "Expires" key exists, and if it does, parses the timestamp and returns a time.Time. If duraction
// is greater than zero(0), then that number of minutes will be add to the expires time that is parsed from the header.
func RetrieveExpiresHeader(h http.Header, duration time.Duration) time.Time {
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

func RetrieveXPagesFromHeader(h http.Header) uint {

	header := h.Get("X-Pages")
	if header == "" {
		return 0
	}

	pages, err := strconv.Atoi(header)
	if err != nil {
		return 0
	}

	return uint(pages)

}

// RetrieveEtagHeader is a helper method that retrieves an
// Etag for the most recent request to ESI
func RetrieveEtagHeader(h http.Header) string {
	return h.Get("Etag")
}

// RetrieveErrorCount is a helper method that retrieves the number of errors that this application
// has triggered and how many more can be triggered before potentially encountereding an HTTP Status 420
func RetrieveErrorCount(h http.Header) uint64 {
	// Default to a low count. This will cause the app to slow down
	// if the header is not present to set the actual value from the header
	var count uint64 = 20
	strCount := h.Get("x-esi-error-limit-remain")
	if strCount != "" {
		i, err := strconv.ParseUint(strCount, 10, 64)
		if err == nil {
			count = i
		}
	}

	return count

}

func (s *service) storeESIErrorCount(ctx context.Context, count uint64) {
	s.cache.SetESIErrCount(ctx, count)
}

// RetrieveErrorReset is a helper method that retrieves the number of seconds until our Error Limit resets
func RetrieveErrorReset(h http.Header) uint64 {
	reset := h.Get("x-esi-error-limit-reset")
	if reset == "" {
		return 0
	}

	seconds, err := strconv.ParseUint(reset, 10, 64)
	if err != nil {
		return 0
	}

	return seconds

}

func (s *service) storeESIErrorReset(ctx context.Context, seconds uint64) {
	s.cache.SetEsiErrorReset(ctx, uint64(time.Now().Add(time.Second*time.Duration(seconds)).Unix()))
}

func (s *service) trackESICallStatusCode(ctx context.Context, code int) {
	s.cache.SetESITracking(ctx, code, uint64(time.Now().UnixNano()))
}

func buildKey(s ...string) string {
	return strings.Join(s, "::")
}
