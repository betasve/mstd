/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package app

import (
	"errors"
	"github.com/betasve/mstd/ext/exec"
	exectest "github.com/betasve/mstd/ext/exec/exectest"
	"github.com/betasve/mstd/ext/runtime"
	runtimetest "github.com/betasve/mstd/ext/runtime/runtimetest"
	"github.com/betasve/mstd/login"
	osexec "os/exec"
	"testing"
	"time"
)

func TestLoginNeededSuccess(test *testing.T) {
	dur, _ := time.ParseDuration("5m")
	creds = login.Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Now().Add(dur))

	result := creds.LoginNeeded()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestLoginNeededFailure(test *testing.T) {
	dur, _ := time.ParseDuration("5m")
	creds := login.Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Now().Add(-dur))

	result := creds.LoginNeeded()
	if !result {
		test.Errorf("\nexpected\ntrue\nbut got\n%v", result)
	}
}

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
