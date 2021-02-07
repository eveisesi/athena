package athena

import (
	"context"
	"time"
)

type EtagRepository interface {
	Etag(ctx context.Context, etagID string) (*Etag, error)
	Etags(ctx context.Context, operators ...*Operator) ([]*Etag, error)
	InsertEtag(ctx context.Context, etag *Etag) (*Etag, error)
	UpdateEtag(ctx context.Context, etagID string, etag *Etag) (*Etag, error)
	DeleteEtag(ctx context.Context, id string) (bool, error)
}

type Etag struct {
	ID          uint      `db:"id" json:"id"`
	EtagID      string    `db:"etag_id" json:"etag_id"`
	Etag        string    `db:"etag" json:"etag"`
	CachedUntil time.Time `db:"cached_until" json:"cached_until"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
