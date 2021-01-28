package athena

import (
	"context"
	"time"
)

type EtagRepository interface {
	Etag(ctx context.Context, etagID string) (*Etag, error)
	Etags(ctx context.Context, operators ...*Operator) ([]*Etag, error)
	UpdateEtag(ctx context.Context, etagID string, etag *Etag) (*Etag, error)
	DeleteEtag(ctx context.Context, etagID string) (bool, error)
}

type Etag struct {
	EtagID      string    `bson:"etag_id" json:"etag_id"`
	Etag        string    `bson:"etag" json:"etag"`
	CachedUntil time.Time `bson:"cached_until" json:"cached_until"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
