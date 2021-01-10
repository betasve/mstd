package app

import (
	"errors"
	"github.com/betasve/mstd/exec"
	exectest "github.com/betasve/mstd/exec/exectest"
	"github.com/betasve/mstd/runtime"
	runtimetest "github.com/betasve/mstd/runtime/runtimetest"
	osexec "os/exec"
	"testing"
)

// TODO: Create test helpers for common parts in different
// variants of tests for this function
func TestOpenLoginUrlForLinuxSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}

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

	err := openLoginUrl(testUrl)
	if err != nil {
		test.Errorf("\nexpected\nno errors\nbut got\n%s", err)
	}

	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}
}

func TestOpenLoginUrlForLinuxFailure(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}

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

	err := openLoginUrl(testUrl)

	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if err != expectedErr {
		test.Errorf(
			"\nexpected\nto run with error\n%s\nbut was\n%s",
			expectedErr,
			err,
		)
	}
}

func TestOpenLoginUrlForWindowsSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}

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

	err := openLoginUrl(testUrl)
	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if err != nil {
		test.Errorf(
			"\nexpected\nto run without errors\nbut was\n%s",
			err,
		)
	}
}

func TestOpenLoginUrlForWindowsFailure(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}

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

	err := openLoginUrl(testUrl)

	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if err != expectedErr {
		test.Errorf(
			"\nexpected\nto run with error\n%s\nbut was\n%s",
			expectedErr,
			err,
		)
	}
}

func TestOpenLoginUrlForMacOSSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}

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

	err := openLoginUrl(testUrl)
	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if err != nil {
		test.Errorf(
			"\nexpected\nto run without errors\nbut was\n%s",
			err,
		)
	}
}

func TestOpenLoginUrlForMacOSFailure(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	exec.CmdClient = &exectest.CmdMock{}
	exec.RunClient = &exectest.CmdRunMock{}

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

	err := openLoginUrl(testUrl)

	if !ranCmd && !ranRunCmd {
		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
	}

	if err != expectedErr {
		test.Errorf(
			"\nexpected\nto run with error\n%s\nbut was\n%s",
			expectedErr,
			err,
		)
	}
}

func TestOpenLoginUrlForUnknownSuccess(test *testing.T) {
	runtime.Client = runtimetest.RuntimeMock{}
	runtimetest.RuntimeMockFunc = func() string { return "unknown" }

	testUrl := "http://localhost/open"
	err := openLoginUrl(testUrl)

	if err == nil {
		test.Errorf(
			"\nexpected\nno erors\nbut got\n%s",
			err,
		)
	}
}
