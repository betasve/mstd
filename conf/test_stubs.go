package conf

type HomedirServiceMock struct{}

var homedirErr error
var homedirPath string

func (h HomedirServiceMock) Dir() (string, error) {
	homedirPath = "/home/user/"
	return homedirPath, homedirErr
}
