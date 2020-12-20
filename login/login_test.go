package login

import (
	"errors"
	"fmt"
	"github.com/betasve/mstd/conf"
	"github.com/betasve/mstd/exec"
	exectest "github.com/betasve/mstd/exec/exectest"
	"github.com/betasve/mstd/log"
	logtest "github.com/betasve/mstd/log/logtest"
	"github.com/betasve/mstd/runtime"
	runtimetest "github.com/betasve/mstd/runtime/runtimetest"
	"github.com/betasve/mstd/time"
	timetest "github.com/betasve/mstd/time/timetest"
	osexec "os/exec"
	"testing"
	t "time"
)

func TestPerformFailWhenAlreadyLoggedIn(test *testing.T) {
	log.Client = logtest.LoggerServiceMock{}
	var result interface{}
	msg := "logged in"
	logtest.PrintlnMock = func(in ...interface{}) { result = msg }
	d, err := time.Client.ParseDuration("5m")
	if err != nil {
		test.Error(err)
	}

	conf.CurrentState = conf.State{
		AccessToken:          "acc_token",
		AccessTokenExpiresAt: time.Client.Now().Add(d),
	}

	Perform()

	if fmt.Sprintf("%v", result) != msg {
		test.Errorf("\nexpected to log\n%s\nbut got\n%s", msg, fmt.Sprintf("%v", result))
	}
}

func TestPrepareLoginUrl(test *testing.T) {
	conf.CurrentState = conf.State{
		ClientId:         "client-id",
		Permissions:      "Some.Permissions",
		AuthCallbackHost: "http://localhost:3000",
		AuthCallbackPath: "/callback",
	}

	baseRequestUrl = "http://test.com"
	authRequestPath = "/authorize"

	result := prepareLoginUrl()
	expectedResult := "http://test.com/authorize?client_id=client-id" +
		"&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback" +
		"&response_mode=query&response_type=code&scope=Some.Permissions"

	if result != expectedResult {
		test.Errorf("\nexpected\n%s\nbut got\n%s", expectedResult, result)
	}
}

// TODO: Create test helpers for common parts in different
// variants of tests for this function
func TestOpenLoginUrlForLinuxSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}
	log.Client = logtest.LoggerServiceMock{}

	testUrl := "http://localhost/open"
	testCmd := osexec.Command("echo", testUrl)

	var ranCmd bool
	var ranRunCmd bool

	runtimetest.RuntimeMockFunc = func() string { return "linux" }

	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
		ranCmd = true
		return exectest.InitCmdMock(testCmd)
	}
	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
		ranRunCmd = true
		return nil
	}

	var logResult interface{}
	logtest.FatalMock = func(s ...interface{}) { logResult = s }

	openLoginUrl(testUrl)
	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if logResult != nil {
		test.Errorf(
			"\nexpected\nto run without errors\nbut was\n%s",
			fmt.Sprintf("%v", logResult),
		)
	}
}

func TestOpenLoginUrlForLinuxFailure(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}
	log.Client = logtest.LoggerServiceMock{}

	testUrl := "http://localhost/open"
	testCmd := osexec.Command("echo", testUrl)

	var ranCmd bool
	var ranRunCmd bool

	expectedErr := errors.New("could not run")
	runtimetest.RuntimeMockFunc = func() string { return "linux" }

	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
		ranCmd = true
		return exectest.InitCmdMock(testCmd)
	}

	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
		ranRunCmd = true
		return expectedErr
	}

	var logResult error
	logtest.FatalMock = func(s ...interface{}) { logResult = expectedErr }

	openLoginUrl(testUrl)

	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if logResult != expectedErr {
		test.Errorf(
			"\nexpected\nto run with error\n%s\nbut was\n%s",
			expectedErr.Error(),
			fmt.Sprintf("%v", logResult),
		)
	}
}

func TestOpenLoginUrlForWindowsSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}
	log.Client = logtest.LoggerServiceMock{}

	testUrl := "http://localhost/open"
	testCmd := osexec.Command("echo", testUrl)

	var ranCmd bool
	var ranRunCmd bool

	runtimetest.RuntimeMockFunc = func() string { return "windows" }

	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
		ranCmd = true
		return exectest.InitCmdMock(testCmd)
	}
	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
		ranRunCmd = true
		return nil
	}

	var logResult interface{}
	logtest.FatalMock = func(s ...interface{}) { logResult = s }

	openLoginUrl(testUrl)
	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if logResult != nil {
		test.Errorf(
			"\nexpected\nto run without errors\nbut was\n%s",
			fmt.Sprintf("%v", logResult),
		)
	}
}

func TestOpenLoginUrlForWindowsFailure(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}
	log.Client = logtest.LoggerServiceMock{}

	testUrl := "http://localhost/open"
	testCmd := osexec.Command("echo", testUrl)

	var ranCmd bool
	var ranRunCmd bool

	expectedErr := errors.New("could not run")
	runtimetest.RuntimeMockFunc = func() string { return "windows" }

	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
		ranCmd = true
		return exectest.InitCmdMock(testCmd)
	}

	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
		ranRunCmd = true
		return expectedErr
	}

	var logResult error
	logtest.FatalMock = func(s ...interface{}) { logResult = expectedErr }

	openLoginUrl(testUrl)

	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if logResult != expectedErr {
		test.Errorf(
			"\nexpected\nto run with error\n%s\nbut was\n%s",
			expectedErr.Error(),
			fmt.Sprintf("%v", logResult),
		)
	}
}

func TestOpenLoginUrlForMacOSSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}
	log.Client = logtest.LoggerServiceMock{}

	testUrl := "http://localhost/open"
	testCmd := osexec.Command("echo", testUrl)

	var ranCmd bool
	var ranRunCmd bool

	runtimetest.RuntimeMockFunc = func() string { return "darwin" }

	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
		ranCmd = true
		return exectest.InitCmdMock(testCmd)
	}
	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
		ranRunCmd = true
		return nil
	}

	var logResult interface{}
	logtest.FatalMock = func(s ...interface{}) { logResult = s }

	openLoginUrl(testUrl)
	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if logResult != nil {
		test.Errorf(
			"\nexpected\nto run without errors\nbut was\n%s",
			fmt.Sprintf("%v", logResult),
		)
	}
}

func TestOpenLoginUrlForMacOSFailure(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}
	log.Client = logtest.LoggerServiceMock{}

	testUrl := "http://localhost/open"
	testCmd := osexec.Command("echo", testUrl)

	var ranCmd bool
	var ranRunCmd bool

	expectedErr := errors.New("could not run")
	runtimetest.RuntimeMockFunc = func() string { return "darwin" }

	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
		ranCmd = true
		return exectest.InitCmdMock(testCmd)
	}

	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
		ranRunCmd = true
		return expectedErr
	}

	var logResult error
	logtest.FatalMock = func(s ...interface{}) { logResult = expectedErr }

	openLoginUrl(testUrl)

	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if logResult != expectedErr {
		test.Errorf(
			"\nexpected\nto run with error\n%s\nbut was\n%s",
			expectedErr.Error(),
			fmt.Sprintf("%v", logResult),
		)
	}
}

func TestOpenLoginUrlForUnknownSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	log.Client = logtest.LoggerServiceMock{}

	runtimetest.RuntimeMockFunc = func() string { return "unknown" }

	expectedLog := "visit some url"
	var logResult interface{}
	logtest.PrintfMock = func(s string, v ...interface{}) { logResult = expectedLog }

	testUrl := "http://localhost/open"
	openLoginUrl(testUrl)

	if logResult == nil {
		test.Errorf(
			"\nexpected to log\n\"%s\"\nbut nothing was logged\n",
			expectedLog,
		)
	}
}

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
