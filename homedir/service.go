package homedir

import "github.com/mitchellh/go-homedir"

var Client HomedirService = Homedir{}

type HomedirService interface {
	Dir() (string, error)
}

type Homedir struct{}

func (h Homedir) Dir() (string, error) {
	return homedir.Dir()
}
