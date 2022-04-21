package model

import "time"

// Audit ...
type Audit struct {
	CreatedDt *time.Time `bson:"crdt"`
	UpdatedDt *time.Time `bson:"updt"`
}
