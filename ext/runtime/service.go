package runtime

import "runtime"

type RuntimeService interface {
	GetOS() string
}

type Runtime struct{}

var Client RuntimeService = Runtime{}

func (r Runtime) GetOS() string {
	return runtime.GOOS
}
