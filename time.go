package fireblocksdk

import "time"

type ITimeProvider interface {
	Now() time.Time
}

type TimeProvider struct{}

func (tp *TimeProvider) Now() time.Time {
	return time.Now()
}
