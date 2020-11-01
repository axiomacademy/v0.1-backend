package db

import (
	"github.com/jackc/pgtype"
	"time"
)

// Utility function to take two time.Time(s) and construct a tstzrange type in pgtype
func getTstzrange(startTime time.Time, endTime time.Time) pgtype.Tstzrange {
	// Convert to timestamp
	upper := pgtype.Timestamptz{}
	lower := pgtype.Timestamptz{}
	upper.Set(endTime)
	lower.Set(startTime)

	// Create period
	return pgtype.Tstzrange{
		Upper:     upper,
		Lower:     lower,
		UpperType: pgtype.Inclusive,
		LowerType: pgtype.Exclusive,
		Status:    pgtype.Present,
	}
}
