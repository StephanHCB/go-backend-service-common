package repository

import "time"

const TimestampAcornName = "timestamp"

type Timestamp interface {
	IsTimestamp() bool

	Now() time.Time
	MockResponse(mock func() time.Time)
}
