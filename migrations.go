package athena

import (
	"context"
	"time"
)

type MigrationRepository interface {
	InitializeMigrationsTable(ctx context.Context, dbname string) error
	Migration(ctx context.Context, migration string) (*Migration, error)
	Migrations(ctx context.Context) ([]*Migration, error)
	CreateMigration(ctx context.Context, migration string) (*Migration, error)
	DeleteMigration(ctx context.Context, migration string) (bool, error)
}

type Migration struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
