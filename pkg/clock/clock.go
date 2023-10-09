package clock

import "time"

type Time interface {
	Now() time.Time
	Today() time.Time
}
type RealClock struct {
}

func NewTime() *RealClock {
	return &RealClock{}
}

func (rc *RealClock) Now() time.Time {
	return time.Now()
}

func (rc *RealClock) Today() time.Time {
	now := rc.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
