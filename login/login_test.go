package login

import (
	"github.com/betasve/mstd/conf"
	"github.com/betasve/mstd/time"
	timetest "github.com/betasve/mstd/time/timetest"
	"testing"
	t "time"
)

func TestAlreadyLoggedInWithAccessTokenSucess(test *testing.T) {
	time.Client = timetest.TimeMock{}
	now := time.Client.Now()

	d, err := time.Client.ParseDuration("5m")
	if err != nil {
		test.Error(err)
	}

	conf.CurrentState = conf.State{
		AccessToken:          "acc_token",
		AccessTokenExpiresAt: now.Add(d),
	}

	timetest.TimeNowMockFunc = func() t.Time { return now }

	result := alreadyLoggedIn()
	if !result {
		test.Errorf("\nexpected\ntrue\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenNoAccessToken(test *testing.T) {
	time.Client = timetest.TimeMock{}
	now := time.Client.Now()

	d, err := time.Client.ParseDuration("5m")
	if err != nil {
		test.Error(err)
	}

	conf.CurrentState = conf.State{
		AccessToken:          "",
		AccessTokenExpiresAt: now.Add(d),
	}

	timetest.TimeNowMockFunc = func() t.Time { return now }

	result := alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenExpiredAccessToken(test *testing.T) {
	time.Client = timetest.TimeMock{}
	now := time.Client.Now()

	d, err := time.Client.ParseDuration("5m")
	if err != nil {
		test.Error(err)
	}

	conf.CurrentState = conf.State{
		AccessToken:          "acc_token",
		AccessTokenExpiresAt: now.Add(-d),
	}

	timetest.TimeNowMockFunc = func() t.Time { return now }

	result := alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInWithRefreshTokenSucess(test *testing.T) {
	time.Client = timetest.TimeMock{}
	now := time.Client.Now()

	d, err := time.Client.ParseDuration("5m")
	if err != nil {
		test.Error(err)
	}

	conf.CurrentState = conf.State{
		RefreshToken:          "ref_token",
		RefreshTokenExpiresAt: now.Add(d),
	}

	timetest.TimeNowMockFunc = func() t.Time { return now }

	result := alreadyLoggedIn()
	if !result {
		test.Errorf("\nexpected\ntrue\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenNoRefreshToken(test *testing.T) {
	time.Client = timetest.TimeMock{}
	now := time.Client.Now()

	d, err := time.Client.ParseDuration("5m")
	if err != nil {
		test.Error(err)
	}

	conf.CurrentState = conf.State{
		RefreshToken:          "",
		RefreshTokenExpiresAt: now.Add(d),
	}

	timetest.TimeNowMockFunc = func() t.Time { return now }

	result := alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenExpiredRefreshToken(test *testing.T) {
	time.Client = timetest.TimeMock{}
	now := time.Client.Now()

	d, err := time.Client.ParseDuration("5m")
	if err != nil {
		test.Error(err)
	}

	conf.CurrentState = conf.State{
		RefreshToken:          "ref_token",
		RefreshTokenExpiresAt: now.Add(-d),
	}

	timetest.TimeNowMockFunc = func() t.Time { return now }

	result := alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}
