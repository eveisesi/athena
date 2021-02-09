package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type migrationRepository struct {
	db    *sqlx.DB
	table string
	scols []string
	icols []string
}

func NewMigrationRepository(db *sql.DB) athena.MigrationRepository {
	return &migrationRepository{
		db:    sqlx.NewDb(db, "mysql"),
		table: "migrations",
		scols: []string{"id", "name", "created_at"},
		icols: []string{"name", "created_at"},
	}
}

const createTableQuery = `
	CREATE TABLE IF NOT EXISTS %s (          
		id INT UNSIGNED NOT NULL AUTO_INCREMENT, 
		name VARCHAR(255) NOT NULL,                  
		created_at TIMESTAMP NOT NULL,               
		PRIMARY KEY (id) USING BTREE,                
		UNIQUE INDEX migrations_name_unique_idx (name)   
	) COLLATE = 'utf8mb4_unicode_ci' ENGINE = INNODB;
`

const checkTableExistsQuery = `
	SELECT 
		COUNT(*)
	FROM information_schema.tables
	WHERE 
		table_schema = ? AND table_name = '%s'
	LIMIT 1;
`

func (r *migrationRepository) InitializeMigrationsTable(ctx context.Context, dbname string) error {

	var count int
	err := r.db.GetContext(ctx, &count, fmt.Sprintf(checkTableExistsQuery, r.table), dbname)
	if err != nil {
		return fmt.Errorf("Failed to check if migrations table exists: %w", err)
	}

	if count == 0 {
		_, err = r.db.ExecContext(ctx, fmt.Sprintf(createTableQuery, r.table))
		if err != nil {
			return fmt.Errorf("Failed to create migrations table: %w", err)
		}
	}

	return nil

}

func (r *migrationRepository) Migration(c context.Context, migration string) (*athena.Migration, error) {

	q, a, e := sq.Select(r.scols...).From(r.table).Where(sq.Eq{"name": migration}).ToSql()
	if e != nil {
		return nil, e
	}

	var m = new(athena.Migration)
	e = r.db.GetContext(c, m, q, a...)

	return m, e
}

func (r *migrationRepository) Migrations(ctx context.Context) ([]*athena.Migration, error) {

	q, a, e := sq.Select(r.scols...).From(r.table).ToSql()
	if e != nil {
		return nil, e
	}

	var ms = make([]*athena.Migration, 0)

	return ms, r.db.SelectContext(ctx, &ms, q, a...)

}

func (r *migrationRepository) CreateMigration(c context.Context, migration string) (*athena.Migration, error) {

	q, a, e := sq.Insert(r.table).Columns(r.icols...).Values(migration, sq.Expr(`NOW()`)).ToSql()
	if e != nil {
		return nil, e
	}

	_, e = r.db.ExecContext(c, q, a...)
	if e != nil {
		return nil, e
	}

	return r.Migration(c, migration)

}

func (r *migrationRepository) DeleteMigration(c context.Context, migration string) (bool, error) {

	q, a, e := sq.Delete(r.table).Where(sq.Eq{"name": migration}).ToSql()
	if e != nil {
		return false, e
	}

	_, e = r.db.ExecContext(c, q, a...)

	return e == nil, e

}
