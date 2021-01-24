package athena

import "time"

type Meta struct {
	NotModifiedCount uint      `bson:"not_modified_count" json:"not_modified_count"`
	UpdatePriority   uint      `bson:"update_priority" json:"update_priority"`
	UpdateError      int64     `bson:"update_error" json:"update_error"`
	Etag             string    `bson:"etag" json:"etag"`
	CachedUntil      time.Time `bson:"cached_until" json:"cached_until"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}
