package athena

import "time"

type Meta struct {
	NotModifiedCount uint      `db:"not_modified_count" json:"not_modified_count"`
	UpdatePriority   uint      `db:"update_priority" json:"update_priority"`
	UpdateError      int64     `db:"update_error" json:"update_error"`
	Etag             string    `db:"etag" json:"etag"`
	CachedUntil      time.Time `db:"cached_until" json:"cached_until"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
