package runtimetest

type RuntimeMock struct{}

var RuntimeMockFunc = func() string { return "linux" }

func (r RuntimeMock) GetOS() string {
	return RuntimeMockFunc()
}
