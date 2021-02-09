package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const migrationDir = ".config/migrations/"

func migrateUpCommand(c *cli.Context) error {

	basics := basics("migrate")

	successes := 0
	strSteps := c.Args().Get(0)
	steps := 0
	if strSteps != "" {
		steps, _ = strconv.Atoi(strSteps)
	}

	files, err := filepath.Glob(fmt.Sprintf("%s*.up.sql", migrationDir))
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to read migrations")
	}

	sort.Strings(files)

	var ctx = context.Background()

	err = basics.repositories.migration.InitializeMigrationsTable(ctx, basics.cfg.MySQL.DB)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize migrations table")
	}

	basics.logger.Info("migrations table initialize successfully")

	printMostRecentMigration(ctx, basics)

	for _, file := range files {

		// Clean the migration file name up so we can query the migration table
		name := strings.TrimPrefix(file, migrationDir)
		name = strings.TrimSuffix(name, ".up.sql")

		_, err := basics.repositories.migration.Migration(ctx, name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			basics.logger.WithError(err).Fatal("failed to check if migration has been executed")
		}

		if err == nil {
			continue
		}

		entry := basics.logger.WithField("name", name)

		handle, err := os.Open(file)
		if err != nil {
			entry.WithError(err).Fatal("failed to open file for migration")
		}

		data, err := ioutil.ReadAll(handle)
		if err != nil {
			entry.WithError(err).Fatal("failed to read file for migration")
		}

		if len(data) == 0 {
			entry.WithError(err).Fatal("empty migration file detected, halting execution")
		}

		query := string(data)

		_, err = basics.db.ExecContext(ctx, query)
		if err != nil {
			entry.WithError(err).Fatal("failed to execute migration")
		}

		_, err = basics.repositories.migration.CreateMigration(ctx, name)
		if err != nil {
			entry.WithError(err).Fatal("failed to log migration execute in migration table")
		}

		entry.Info("migration executed successfully")
		successes++
		if successes >= steps && steps > 0 {
			break
		}

	}

	return nil

}

func migrateDownCommand(c *cli.Context) error {

	var ctx = context.Background()

	basics := basics("migrate")

	strSteps := c.Args().First()
	steps := 0
	if strSteps != "" {
		steps, _ = strconv.Atoi(strSteps)
	}

	if steps == 0 {
		basics.logger.Fatal("Destructive down command requires explict negating of steps. To run all down migration, please pass -1 or pass N > 0")
	}

	successes := 0

	migrations, err := basics.repositories.migration.Migrations(ctx)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to fetch migration from database")
	}

	printMostRecentMigration(ctx, basics)

	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]

		fileName := fmt.Sprintf("%s%s.down.sql", migrationDir, migration.Name)

		entry := basics.logger.WithFields(logrus.Fields{
			"name":     migration.Name,
			"fileName": fileName,
		})
		file, err := os.Open(fileName)
		if err != nil {
			entry.WithError(err).Fatal("failed to open migration file")
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			entry.WithError(err).Fatal("failed to read migration file")
		}

		if len(data) == 0 {
			entry.WithError(err).Fatal("empty migration file detected, halting execution")
		}

		query := string(data)

		_, err = basics.db.ExecContext(ctx, query)
		if err != nil {
			entry.WithError(err).Fatal("failed to execute query")
		}

		_, err = basics.repositories.migration.DeleteMigration(ctx, migration.Name)
		if err != nil {
			entry.WithError(err).Fatal("failed to remove migration from migrations table")
		}

		entry.Info("migration executed successfully")
		successes++
		if successes >= steps && steps > 0 {
			break
		}

	}

	return nil

}

func migrateCreateCommand(c *cli.Context) error {

	name := c.String("name")
	if name == "" {
		return fmt.Errorf("name is required, received empty string")
	}

	basics := basics("migrations")

	now := time.Now()

	filename := fmt.Sprintf("%s%s_%s.%%s.sql", migrationDir, now.Format("20060102150405"), name)
	up := fmt.Sprintf(filename, "up")
	entry := basics.logger.WithFields(logrus.Fields{
		"name": name,
	})
	_, err := os.Create(up)
	if err != nil {
		entry.WithField("up", up).WithError(err).Fatal("failed to create up file")
	}

	entry.WithField("up", up).Info("migration created successfully")
	down := fmt.Sprintf(filename, "down")
	_, err = os.Create(down)
	if err != nil {
		entry.WithField("down", down).WithError(err).Fatal("failed to create down file")
	}

	entry.WithField("down", down).Info("migration created successfully")

	return nil

}

func migrateCurrentCommand(c *cli.Context) error {

	basics := basics("migrate")
	var ctx = context.Background()

	printMostRecentMigration(ctx, basics)
	return nil
}

func printMostRecentMigration(ctx context.Context, basics *app) {

	migrations, err := basics.repositories.migration.Migrations(ctx)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to fetch migration from database")
	}

	if len(migrations) == 0 {
		basics.logger.Info("no migrations bave been run")
		return
	}

	basics.logger.WithField("current", migrations[len(migrations)-1].Name).Info("current migration")

}
