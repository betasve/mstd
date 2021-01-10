package time

import (
	t "time"
)

type TimeService interface {
	Add(d t.Duration) t.Time
	Now() t.Time
	ParseDuration(s string) (t.Duration, error)
	Unix() int64
	UnixNano() int64
}

type Time struct {
	time t.Time
}

var Client TimeService = Time{}

func (tm Time) Add(d t.Duration) t.Time {
	return tm.time.Add(d)
}

func (tm Time) Now() t.Time {
	tm.time = t.Now()
	return tm.time
}

func (tm Time) ParseDuration(s string) (t.Duration, error) {
	return t.ParseDuration(s)
}

func (tm Time) Unix() int64 {
	return tm.time.Unix()
}

func (tm Time) UnixNano() int64 {
	return tm.time.UnixNano()
}
