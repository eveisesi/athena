package universe

// import (
// 	"context"
// 	"sync"

// 	"github.com/eveisesi/athena"
// 	"github.com/eveisesi/athena/internal/cache"
// 	"github.com/sirupsen/logrus"
// 	"github.com/vbauerster/mpb"
// 	"github.com/vbauerster/mpb/decor"
// )

// var (
// 	systemTotal, planetTotal,
// 	beltTotal, moonTotal int64
// 	mux = new(sync.Mutex)
// )

// func (s *service) initializeMap(ctx context.Context) {

// 	entry := s.logger.WithContext(ctx)

// 	regionIDs, _, err := s.esi.GetRegions(ctx, make([]uint, 0, 110))
// 	if err != nil {
// 		entry.WithError(err).Error("failed to fetch list of region id from ESI")
// 		return
// 	}

// 	wg := &sync.WaitGroup{}

// 	p := mpb.New(mpb.WithWaitGroup(wg))

// 	regionsName := "Regions"
// 	regionsBar := p.AddBar(int64(len(regionIDs)),
// 		mpb.PrependDecorators(
// 			// simple name decorator
// 			decor.Name(regionsName, decor.WC{W: len(regionsName) + 1, C: decor.DidentRight}),
// 			// decor.DSyncWidth bit enables column width synchronization
// 			decor.CountersNoUnit("%d / %d"),
// 		),
// 		mpb.AppendDecorators(decor.Percentage()),
// 	)

// 	for _, regionID := range regionIDs {

// 		systemQueue, systemBar := s.initializeMapWorkers(wg, p, entry)

// 		regionEntry := s.logger.WithField("region_id", regionID)

// 		region, _, err := s.esi.GetRegion(ctx, &athena.Region{ID: regionID})
// 		if err != nil {
// 			regionEntry.WithError(err).Error("failed to fetch region from ESI")
// 			continue
// 		}

// 		constellationIDs := region.ConstellationIDs

// 		region, err = s.universe.CreateRegion(ctx, region)
// 		if err != nil {
// 			regionEntry.WithError(err).Error("failed to create region in DB")
// 			continue
// 		}

// 		constellationsName := "Constellations"
// 		constellationsBar := p.AddBar(int64(len(constellationIDs)), mpb.BarRemoveOnComplete(),
// 			mpb.PrependDecorators(
// 				// simple name decorator
// 				decor.Name(constellationsName, decor.WC{W: len(constellationsName) + 1, C: decor.DidentRight}),
// 				// decor.DSyncWidth bit enables column width synchronization
// 				decor.CountersNoUnit("%d / %d"),
// 			),
// 			mpb.AppendDecorators(decor.Percentage()),
// 		)

// 		for _, constellationID := range constellationIDs {

// 			constellationsEntry := regionEntry.WithField("constellation_id", constellationID)

// 			constellation, _, err := s.esi.GetConstellation(ctx, &athena.Constellation{ID: constellationID})
// 			if err != nil {
// 				constellationsEntry.WithError(err).Error("failed to fetch constellation from ESI")
// 				return
// 			}

// 			systemIDs := constellation.SystemIDs

// 			constellation, err = s.universe.CreateConstellation(ctx, constellation)
// 			if err != nil {
// 				constellationsEntry.WithError(err).Error("failed to create constellation in DB")
// 				return
// 			}

// 			s.updateSystemTotal(int64(len(systemIDs)), systemBar)
// 			for _, systemID := range systemIDs {
// 				wg.Add(1)
// 				systemQueue <- systemID
// 			}

// 			err = s.cache.SetConstellation(ctx, constellation, cache.ExpiryHours(0))
// 			if err != nil {
// 				constellationsEntry.WithError(err).Error("failed to cache constellation")
// 				return
// 			}

// 			wg.Wait()

// 			constellationsBar.Increment()

// 		}

// 		err = s.cache.SetRegion(ctx, region, cache.ExpiryHours(0))
// 		if err != nil {
// 			regionEntry.WithError(err).Error("failed to cache region")
// 		}

// 		regionsBar.Increment()

// 		systemTotal = 0

// 		planetTotal = 0
// 		moonTotal = 0
// 		beltTotal = 0

// 	}

// }

// func (s *service) initializeMapWorkers(
// 	wg *sync.WaitGroup, p *mpb.Progress, logEntry *logrus.Entry,
// ) (
// 	systemQueue chan uint, systemBar *mpb.Bar,
// ) {

// 	systemQueue = make(chan uint, 10)
// 	planetQueue := make(chan uint, 10)
// 	moonQueue := make(chan uint, 10)
// 	beltQueue := make(chan uint, 10)

// 	systemName := "Systems Imported"
// 	systemBar = p.AddBar(
// 		0,
// 		mpb.PrependDecorators(
// 			// simple name decorator
// 			decor.Name(systemName, decor.WC{W: len(systemName) + 1, C: decor.DidentRight}),
// 			// decor.DSyncWidth bit enables column width synchronization
// 			decor.CountersNoUnit("%d / %d"),
// 		),
// 		mpb.AppendDecorators(decor.Percentage()),
// 	)

// 	planetName := "Planets Imported"
// 	planetBar := p.AddBar(
// 		0,
// 		mpb.PrependDecorators(
// 			// simple name decorator
// 			decor.Name(planetName, decor.WC{W: len(planetName) + 1, C: decor.DidentRight}),
// 			// decor.DSyncWidth bit enables column width synchronization
// 			decor.CountersNoUnit("%d / %d"),
// 		),
// 		mpb.AppendDecorators(decor.Percentage()),
// 	)

// 	beltName := "Asteroid Belts Imported"
// 	beltBar := p.AddBar(
// 		0,
// 		mpb.PrependDecorators(
// 			// simple name decorator
// 			decor.Name(beltName, decor.WC{W: len(beltName) + 1, C: decor.DidentRight}),
// 			// decor.DSyncWidth bit enables column width synchronization
// 			decor.CountersNoUnit("%d / %d"),
// 		),
// 		mpb.AppendDecorators(decor.Percentage()),
// 	)

// 	moonName := "Moons Imported"
// 	moonBar := p.AddBar(
// 		0,
// 		mpb.PrependDecorators(
// 			// simple name decorator
// 			decor.Name(moonName, decor.WC{W: len(moonName) + 1, C: decor.DidentRight}),
// 			// decor.DSyncWidth bit enables column width synchronization
// 			decor.CountersNoUnit("%d / %d"),
// 		),
// 		mpb.AppendDecorators(decor.Percentage()),
// 	)

// 	for w := 1; w <= 5; w++ {
// 		go s.systemWorker(systemQueue, planetQueue, moonQueue, beltQueue, systemBar, planetBar, moonBar, beltBar, wg, logEntry)
// 		go s.planetWorker(planetQueue, wg, planetBar, logEntry)
// 		go s.beltWorker(beltQueue, wg, beltBar, logEntry)
// 		go s.moonWorker(moonQueue, wg, moonBar, logEntry)
// 	}

// 	return

// }

// func (s *service) updateSystemTotal(n int64, bar *mpb.Bar) {
// 	s.logger.WithField("n", n).Info("updateSystemTotal")
// 	mux.Lock()
// 	newTotal := systemTotal + n
// 	bar.SetTotal(newTotal, false)
// 	systemTotal = newTotal
// 	mux.Unlock()
// }

// func (s *service) updatePlanetTotal(n int64, bar *mpb.Bar) {
// 	s.logger.WithField("n", n).Info("updatePlanetTotal")
// 	mux.Lock()
// 	newTotal := planetTotal + n
// 	bar.SetTotal(newTotal, false)
// 	planetTotal = newTotal
// 	mux.Unlock()
// }

// func (s *service) updateBeltTotal(n int64, bar *mpb.Bar) {
// 	s.logger.WithField("n", n).Info("updateBeltTotal")
// 	mux.Lock()
// 	newTotal := beltTotal + n
// 	bar.SetTotal(newTotal, false)
// 	beltTotal = newTotal
// 	mux.Unlock()
// }

// func (s *service) updateMoonTotal(n int64, bar *mpb.Bar) {
// 	s.logger.WithField("n", n).Info("updateMoonTotal")
// 	mux.Lock()
// 	newTotal := moonTotal + n
// 	bar.SetTotal(newTotal, false)
// 	moonTotal = newTotal
// 	mux.Unlock()
// }

// func (s *service) systemWorker(
// 	systemQueue, planetQueue, moonQueue, beltQueue chan uint,
// 	systemBar, planetBar, moonBar, beltBar *mpb.Bar,
// 	wg *sync.WaitGroup, logEntry *logrus.Entry) {

// 	for systemID := range systemQueue {

// 		var ctx = context.Background()

// 		systemEntry := logEntry.WithField("system_id", systemID)
// 		systemEntry.Warn()

// 		system, _, err := s.esi.GetSolarSystem(ctx, &athena.SolarSystem{ID: systemID})
// 		if err != nil {
// 			systemEntry.WithError(err).Error("failed to fetch system from ESI")
// 			wg.Done()
// 			systemBar.Increment()
// 			continue
// 		}

// 		planets := system.Planets

// 		system, err = s.universe.CreateSolarSystem(ctx, system)
// 		if err != nil {
// 			systemEntry.WithError(err).Error("failed to create solar system in database")
// 			wg.Done()
// 			systemBar.Increment()
// 			continue
// 		}

// 		if len(planets) == 0 {
// 			wg.Done()
// 			systemBar.Increment()
// 			continue
// 		}

// 		s.updatePlanetTotal(int64(len(planets)), planetBar)
// 		for _, planet := range planets {
// 			if len(planet.MoonIDs) > 0 {
// 				s.updateMoonTotal(int64(len(planet.MoonIDs)), moonBar)
// 				for _, moonID := range planet.MoonIDs {
// 					wg.Add(1)
// 					moonQueue <- moonID
// 				}
// 			}

// 			if len(planet.BeltIDs) > 0 {
// 				s.updateBeltTotal(int64(len(planet.BeltIDs)), beltBar)
// 				for _, beltID := range planet.BeltIDs {
// 					wg.Add(1)
// 					beltQueue <- beltID
// 				}
// 			}

// 			wg.Add(1)
// 			planetQueue <- planet.ID
// 		}

// 		system.Planets = nil

// 		err = s.cache.SetSolarSystem(ctx, system)
// 		if err != nil {
// 			systemEntry.WithError(err).Error("failed to cache solar system")
// 		}

// 		wg.Done()
// 		systemBar.Increment()

// 	}

// }

// func (s *service) planetWorker(buffPlanetQ chan uint, wg *sync.WaitGroup, bar *mpb.Bar, logEntry *logrus.Entry) {

// 	for planetID := range buffPlanetQ {

// 		var ctx = context.Background()

// 		planetEntry := logEntry.WithContext(ctx).WithField("planet_id", planetID)

// 		planet, _, err := s.esi.GetPlanet(ctx, &athena.Planet{ID: planetID})
// 		if err != nil {
// 			planetEntry.WithError(err).Error("failed to fetch planet from ESI")
// 			wg.Done()
// 			bar.Increment()
// 			continue
// 		}

// 		planet, err = s.universe.CreatePlanet(ctx, planet)
// 		if err != nil {
// 			planetEntry.WithError(err).Error("failed to create planet in db")
// 			wg.Done()
// 			bar.Increment()
// 			continue
// 		}

// 		err = s.cache.SetPlanet(ctx, planet)
// 		if err != nil {
// 			planetEntry.WithError(err).Error("failed to cache planet")
// 		}

// 		wg.Done()
// 		bar.Increment()

// 	}

// }

// func (s *service) moonWorker(buffMoonQ chan uint, wg *sync.WaitGroup, bar *mpb.Bar, logEntry *logrus.Entry) {

// 	for moonID := range buffMoonQ {

// 		var ctx = context.Background()

// 		moonEntry := logEntry.WithContext(ctx).WithField("moon_id", moonID)

// 		moon, _, err := s.esi.GetMoon(ctx, &athena.Moon{ID: moonID})
// 		if err != nil {
// 			moonEntry.WithError(err).Error("failed to fetch moon from ESI")
// 			wg.Done()
// 			bar.Increment()
// 			continue
// 		}

// 		moon, err = s.universe.CreateMoon(ctx, moon)
// 		if err != nil {
// 			moonEntry.WithError(err).Error("failed to create moon in db")
// 			wg.Done()
// 			bar.Increment()
// 			continue
// 		}

// 		err = s.cache.SetMoon(ctx, moon)
// 		if err != nil {
// 			moonEntry.WithError(err).Error("failed to cache moon")
// 		}

// 		wg.Done()
// 		bar.Increment()

// 	}

// }

// func (s *service) beltWorker(buffBeltQ chan uint, wg *sync.WaitGroup, bar *mpb.Bar, logEntry *logrus.Entry) {

// 	for beltID := range buffBeltQ {

// 		var ctx = context.Background()

// 		beltEntry := logEntry.WithContext(ctx).WithField("belt_id", beltID)

// 		belt, _, err := s.esi.GetAsteroidBelt(ctx, &athena.AsteroidBelt{ID: beltID})
// 		if err != nil {
// 			beltEntry.WithError(err).Error("failed to fetch belt from ESI")
// 			wg.Done()
// 			bar.Increment()
// 			continue
// 		}

// 		belt, err = s.universe.CreateAsteroidBelt(ctx, belt)
// 		if err != nil {
// 			beltEntry.WithError(err).Error("failed to create belt in db")
// 			wg.Done()
// 			bar.Increment()
// 			continue
// 		}

// 		err = s.cache.SetAsteroidBelt(ctx, belt)
// 		if err != nil {
// 			beltEntry.WithError(err).Error("failed to cache belt")
// 		}

// 		wg.Done()
// 		bar.Increment()
// 	}

// }
