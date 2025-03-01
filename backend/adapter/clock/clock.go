package clock

import (
	"time"
)

type IClock interface {
	Now() time.Time
}

type Clock struct{}

func New() IClock {
	return Clock{}
}

func (c Clock) Now() time.Time {
	return time.Now()
}
