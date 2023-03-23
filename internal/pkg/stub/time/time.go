package time

import "time"

type Clock interface {
	Now() time.Time
}

type stdClock struct{}

func New() Clock {
	return stdClock{}
}

func (stdClock) Now() time.Time {
	return time.Now()
}
