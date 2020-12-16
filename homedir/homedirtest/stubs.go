package homedirtest

type HomedirServiceMock struct{}

var HomedirErr error
var HomedirPath string

func (h HomedirServiceMock) Dir() (string, error) {
	HomedirPath = "/home/user/"
	return HomedirPath, HomedirErr
}
