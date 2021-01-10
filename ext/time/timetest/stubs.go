package timetest

import (
	t "time"
)

var TimeAddMockFunc = func(d t.Duration) t.Time { return t.Now() }
var TimeNowMockFunc = func() t.Time { return t.Now() }
var TimeParseDurationMockFunc = func(s string) (t.Duration, error) { return t.Since(t.Now()), nil }
var TimeUnixMockFunc = func() int64 { return t.Now().Unix() }
var TimeUnixNanoMockFunc = func() int64 { return t.Now().UnixNano() }

type TimeMock struct{}

func (tm TimeMock) Add(d t.Duration) t.Time {
	return TimeAddMockFunc(d)
}

func (tm TimeMock) Now() t.Time {
	return TimeNowMockFunc()
}

func (tm TimeMock) ParseDuration(s string) (t.Duration, error) {
	return TimeParseDurationMockFunc(s)
}

func (tm TimeMock) Unix() int64 {
	return TimeUnixMockFunc()
}

func (tm TimeMock) UnixNano() int64 {
	return TimeUnixNanoMockFunc()
}
