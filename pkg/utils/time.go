package utils

import "time"

type TimeNower interface {
	Now() time.Time
}

type TimeNow struct{}

func (TimeNow) Now() time.Time {
	return time.Now()
}
