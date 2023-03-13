package timestamp

import "time"

type TimestampImpl struct {
	Timestamp func() time.Time
}

func (r *TimestampImpl) MockResponse(mock func() time.Time) {
	r.Timestamp = mock
}

func (r *TimestampImpl) Now() time.Time {
	return r.Timestamp()
}
