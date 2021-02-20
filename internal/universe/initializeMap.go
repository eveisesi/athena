package universe

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/eveisesi/athena/internal"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func (s *service) InitializeUniverse(options ...OptionFunc) error {

	o := s.options(options...)

	var wg = &sync.WaitGroup{}
	var ctx = context.Background()
	var bar = &mpb.Bar{}
	pOpts := make([]mpb.ProgressOption, 0)
	pOpts = append(pOpts, mpb.WithWaitGroup(wg))
	if o.disableProgress {
		pOpts = append(pOpts, mpb.WithOutput(ioutil.Discard))
	}
	p := mpb.New(pOpts...)

	if o.chr {
		races, _, _, err := s.esi.GetRaces(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch races from ESI: %w", err)
		}

		name := "Races"
		bar = p.AddBar(int64(len(races)),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		for _, race := range races {
			entry := s.logger.WithField("race_id", race.ID)
			_, err = s.universe.CreateRace(ctx, race)
			if err != nil {
				entry.WithError(err).Error("failed to insert race into db")
				continue
			}

			err = s.cache.SetRace(ctx, race)
			if err != nil {
				entry.WithError(err).Error("failed to cache race")
			}

			bar.Increment()
			time.Sleep(time.Millisecond * 500)
		}

		bloodlines, _, _, err := s.esi.GetBloodlines(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch bloodlines from ESI: %w", err)
		}

		name = "Bloodlines"
		bar = p.AddBar(int64(len(bloodlines)),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		for _, bloodline := range bloodlines {
			entry := s.logger.WithField("bloodline_id", bloodline.ID)
			_, err = s.universe.CreateBloodline(ctx, bloodline)
			if err != nil {
				entry.WithError(err).Error("failed to create bloodline in DB")
				continue
			}

			err = s.cache.SetBloodline(ctx, bloodline)
			if err != nil {
				s.logger.WithError(err).WithField("bloodline_id", bloodline.ID).Error("failed to cache bloodline")
			}

			bar.Increment()

			time.Sleep(time.Millisecond * 50)
		}

		ancestries, _, _, err := s.esi.GetAncestries(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch ancestries from ESI: %w", err)
		}

		name = "Ancestries"
		bar = p.AddBar(int64(len(ancestries)),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		for _, ancestry := range ancestries {
			entry := s.logger.WithField("ancestry_id", ancestry.ID)
			_, err = s.universe.CreateAncestry(ctx, ancestry)
			if err != nil {
				entry.WithError(err).Error("failed to create ancestry in DB")
				continue
			}

			err = s.cache.SetAncestry(ctx, ancestry)
			if err != nil {
				entry.WithError(err).Error("failed to cache ancestry")
			}

			bar.Increment()
			time.Sleep(time.Millisecond * 50)
		}

		factions, _, _, err := s.esi.GetFactions(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch factions from ESI: %w", err)
		}

		name = "Factions"
		bar = p.AddBar(int64(len(factions)),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		for _, faction := range factions {
			entry := s.logger.WithField("faction_id", faction.ID)
			_, err = s.universe.CreateFaction(ctx, faction)
			if err != nil {
				entry.WithError(err).Error("failed to create faction in DB")
				continue
			}

			err = s.cache.SetFaction(ctx, faction)
			if err != nil {
				entry.WithError(err).Error("failed to cache faction")
			}

			bar.Increment()

			time.Sleep(time.Millisecond * 50)
		}

	}

	if o.inv {
		categoryIDs, _, _, err := s.esi.GetCategories(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch category IDs from ESI: %w", err)
		}

		categoriesName := "Categories"
		categoriesBar := p.AddBar(int64(len(categoryIDs)),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(categoriesName, decor.WC{W: len(categoriesName) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		for _, categoryID := range categoryIDs {
			categoryEntry := s.logger.WithField("category_id", categoryID)

			category, _, _, err := s.esi.GetCategory(ctx, categoryID)
			if err != nil {
				categoryEntry.WithError(err).Error("failed to fetch category from ESI")
				continue
			}

			_, err = s.universe.CreateCategory(ctx, category)
			if err != nil {
				categoryEntry.WithError(err).Error("failed to create category in DB")
				continue
			}

			groupsName := fmt.Sprintf("Groups of Category %d", categoryID)
			groupsBar := p.AddBar(int64(len(category.Groups)), mpb.BarRemoveOnComplete(),
				mpb.PrependDecorators(
					// simple name decorator
					decor.Name(groupsName, decor.WC{W: len(groupsName) + 1, C: decor.DidentRight}),
					// decor.DSyncWidth bit enables column width synchronization
					decor.CountersNoUnit("%d / %d"),
				),
				mpb.AppendDecorators(decor.Percentage()),
			)

			buffQ := make(chan uint, 20)

			for w := 1; w <= 10; w++ {
				go s.groupWorker(w, buffQ, wg, p, groupsBar, categoryEntry)
			}

			for j := 0; j < len(category.Groups); j++ {
				wg.Add(1)
				buffQ <- category.Groups[j]
			}

			close(buffQ)

			wg.Wait()

			groupsBar.SetTotal(int64(len(category.Groups)), true)

			category.Groups = nil

			err = s.cache.SetCategory(ctx, category)
			if err != nil {
				categoryEntry.WithError(err).Error("failed to cache category")
			}

			categoriesBar.Increment()

		}

		categoriesBar.SetTotal(int64(len(categoryIDs)), true)
	}

	if o.loc {
		regionIDs, _, _, err := s.esi.GetRegions(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch region id for ESI: %w", err)
		}

		regionsName := "Regions"
		regionsBar := p.AddBar(int64(len(regionIDs)),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(regionsName, decor.WC{W: len(regionsName) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		for _, regionID := range regionIDs {

			regionEntry := s.logger.WithField("region_id", regionID)

			region, _, _, err := s.esi.GetRegion(ctx, regionID)
			if err != nil {
				regionEntry.WithError(err).Error("failed to fetch region from ESI")
				continue
			}

			_, err = s.universe.CreateRegion(ctx, region)
			if err != nil {
				regionEntry.WithError(err).Error("failed to create region in DB")
				continue
			}

			constellationsName := "Constellations"
			constellationsBar := p.AddBar(int64(len(region.ConstellationIDs)), mpb.BarRemoveOnComplete(),
				mpb.PrependDecorators(
					// simple name decorator
					decor.Name(constellationsName, decor.WC{W: len(constellationsName) + 1, C: decor.DidentRight}),
					// decor.DSyncWidth bit enables column width synchronization
					decor.CountersNoUnit("%d / %d"),
				),
				mpb.AppendDecorators(decor.Percentage()),
			)

			for _, constellationID := range region.ConstellationIDs {
				wg.Add(1)
				go func(constellationID uint) {
					defer wg.Done()
					defer constellationsBar.Increment()

					constellationsEntry := regionEntry.WithField("constellation_id", constellationID)

					constellation, _, _, err := s.esi.GetConstellation(ctx, constellationID)
					if err != nil {
						constellationsEntry.WithError(err).Error("failed to fetch constellation from ESI")
						return
					}

					_, err = s.universe.CreateConstellation(ctx, constellation)
					if err != nil {
						constellationsEntry.WithError(err).Error("failed to create constellation in DB")
						return
					}

					systemsName := fmt.Sprintf("Systems for Constellation %d", constellationID)
					systemsBar := p.AddBar(int64(len(constellation.SystemIDs)), mpb.BarRemoveOnComplete(),
						mpb.PrependDecorators(
							// simple name decorator
							decor.Name(systemsName, decor.WC{W: len(systemsName) + 1, C: decor.DidentRight}),
							// decor.DSyncWidth bit enables column width synchronization
							decor.CountersNoUnit("%d / %d"),
						),
						mpb.AppendDecorators(decor.Percentage()),
					)

					for _, systemID := range constellation.SystemIDs {
						systemEntry := constellationsEntry.WithField("type_id", systemID)
						system, _, _, err := s.esi.GetSolarSystem(ctx, systemID)
						if err != nil {
							systemEntry.WithError(err).Error("failed to fetch system from ESI")
							continue
						}

						_, err = s.universe.CreateSolarSystem(ctx, system)
						if err != nil {
							systemEntry.WithError(err).Error("failed to create system from DB")
							continue
						}

						err = s.cache.SetSolarSystem(ctx, system)
						if err != nil {
							systemEntry.WithError(err).Error("failed to cache system")
						}

						systemsBar.Increment()

					}

					systemsBar.SetTotal(int64(len(constellation.SystemIDs)), true)

					constellation.SystemIDs = nil

					err = s.cache.SetConstellation(ctx, constellation)
					if err != nil {
						constellationsEntry.WithError(err).Error("failed to cache constellation")
						return
					}
				}(constellationID)
			}
			wg.Wait()
			constellationsBar.SetTotal(int64(len(region.ConstellationIDs)), true)

			region.ConstellationIDs = nil

			err = s.cache.SetRegion(ctx, region)
			if err != nil {
				regionEntry.WithError(err).Error("failed to cache region")
			}

			regionsBar.Increment()
		}
	}

	return nil
}

func (s *service) groupWorker(wid int, buffQ chan uint, wg *sync.WaitGroup, progress *mpb.Progress, bar *mpb.Bar, logEntry *logrus.Entry) {

	for groupID := range buffQ {
		var ctx = context.Background()

		groupEntry := logEntry.WithField("group_id", groupID)

		group, _, _, err := s.esi.GetGroup(ctx, groupID)
		if err != nil {
			groupEntry.WithError(err).Error("failed to fetch group from ESI")
			return
		}

		_, err = s.universe.CreateGroup(ctx, group)
		if err != nil {
			groupEntry.WithError(err).Error("failed to create group from DB")
			return
		}

		if len(group.Types) == 0 {
			wg.Done()
			bar.Increment()
			continue
		}

		typesName := fmt.Sprintf("Types for Group %d", groupID)
		typesBar := progress.AddBar(int64(len(group.Types)), mpb.BarRemoveOnComplete(),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(typesName, decor.WC{W: len(typesName) + 1, C: decor.DidentRight}),
				// decor.DSyncWidth bit enables column width synchronization
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)

		buffQQ := make(chan uint, 20)

		wgi := &sync.WaitGroup{}
		numWorkers := 5
		for w := 0; w <= numWorkers; w++ {
			go s.typeWorker(w, buffQQ, wgi, typesBar, groupEntry)
		}

		size := math.Ceil(float64(len(group.Types)) / float64(numWorkers))

		chunks := internal.ChunkSliceUints(group.Types, int(size))

		for _, chunk := range chunks {
			for j := 0; j < len(chunk); j++ {
				wgi.Add(1)
				buffQQ <- chunk[j]
			}
		}

		// for j := 0; j < len(group.Types); j++ {
		// 	wgi.Add(1)
		// 	buffQQ <- group.Types[j]
		// }

		close(buffQQ)

		wgi.Wait()

		typesBar.SetTotal(int64(len(group.Types)), true)

		group.Types = nil

		err = s.cache.SetGroup(ctx, group)
		if err != nil {
			groupEntry.WithError(err).Error("failed to cache category")
			return
		}
		wg.Done()
		bar.Increment()
	}
}

func (s *service) typeWorker(wid int, buffQQ chan uint, wg *sync.WaitGroup, bar *mpb.Bar, groupEntry *logrus.Entry) {

	for typeID := range buffQQ {
		var ctx = context.Background()

		typeEntry := groupEntry.WithField("type_id", typeID)
		item, _, _, err := s.esi.GetType(ctx, typeID)
		if err != nil {
			typeEntry.WithError(err).Error("failed to fetch type from ESI")
			continue
		}

		_, err = s.universe.CreateType(ctx, item)
		if err != nil {
			typeEntry.WithError(err).Error("failed to create type from DB")
			continue
		}

		err = s.cache.SetType(ctx, item)
		if err != nil {
			typeEntry.WithError(err).Error("failed to cache type")
		}

		bar.Increment()
		wg.Done()
	}

}
