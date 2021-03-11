package entities

import "time"

// Metadata are generated by the store attached to stored items
type Metadata struct {
	Version     int
	Disabled    bool
	ExpireAt    time.Time
	CreatedAt   time.Time
	DeletedAt   time.Time
	DestroyedAt time.Time
}