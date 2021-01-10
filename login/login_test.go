package login

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	// "github.com/betasve/mstd/exec"
	// exectest "github.com/betasve/mstd/exec/exectest"
	// "github.com/betasve/mstd/log"
	// logtest "github.com/betasve/mstd/log/logtest"
	// "github.com/betasve/mstd/runtime"
	// runtimetest "github.com/betasve/mstd/runtime/runtimetest"
	"github.com/betasve/mstd/time"
	"net/url"
	// timetest "github.com/betasve/mstd/time/timetest"
	// osexec "os/exec"
	httpService "github.com/betasve/mstd/http/httptest"
	"net/http"
	"testing"
	t "time"
)

var fiveMins = func() t.Duration {
	d, _ := time.Client.ParseDuration("5m")

	return d
}()

func TestPerformLoginWhenValidAccessToken(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	err := creds.PerformLogin()

	if err != nil {
		test.Errorf("\nexpected no errors\nbut got\n%s", err)
	}
}

func TestPerformLoginWhenValidRefreshToken(test *testing.T) {
	var refreshFuncCalled bool

	httpClient = &httpService.ClientMock{}
	refreshCalledFn := func(a *AuthData) error { refreshFuncCalled = true; return nil }

	creds := Creds{}
	creds.SetRefreshToken("acc_token")
	creds.SetRefreshTokenExpiresAt(time.Client.Now().Add(fiveMins))
	creds.SetLoginDataCallbackFn(refreshCalledFn)

	err := creds.PerformLogin()

	if err != nil {
		test.Errorf("\nexpected no errors\nbut got\n%s", err)
	}

	if !refreshFuncCalled {
		test.Error("\nexpected\nto call for refresh token\nbut\ndid not")
	}
}

func TestPerformLoginWhenNoTokens(test *testing.T) {
	creds := Creds{}
	var handlerFuncCalled, listenAndServeFuncCalled bool

	httpClient = &httpService.ClientMock{}
	httpService.HandlerStubFn =
		func(in string, handler func(w http.ResponseWriter, r *http.Request)) {
			handlerFuncCalled = true
		}

	httpService.ListenAndServeStubFn =
		func(addr string, handler http.Handler) error {
			listenAndServeFuncCalled = true
			return nil
		}

	creds.SetLoginUrlHandlerFn(func(in string) error { return nil })
	creds.SetAuthCallbackPath("/test")

	err := creds.PerformLogin()

	if err != nil {
		test.Errorf("\nexpected no errors\nbut got\n%s", err)
	}

	if !handlerFuncCalled {
		test.Error("\nexpected to have called HandleFunc but it did not")
	}

	if !listenAndServeFuncCalled {
		test.Error("\nexpected to have called ListenAndServe but it did not")
	}
}

func TestGetAccessTokenSuccess(test *testing.T) {
	var response *AuthData

	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))
	creds.SetLoginDataCallbackFn(func(a *AuthData) error { response = a; return nil })
	authKey := "testAuthKey"

	err := creds.getAccessToken(authKey)

	if err != nil {
		test.Errorf("\nexpected no errors\nbut got\n%s", err)
	}

	if len(response.TokenType) == 0 ||
		len(response.Scope) == 0 ||
		response.ExpiresIn == 0 ||
		response.ExtExpiresIn == 0 ||
		len(response.AccessToken) == 0 ||
		len(response.RefreshToken) == 0 {
		test.Error("\nexpected all AuthData values to be populated\nbut some were not\n")
	}
}

func TestGetAccessTokenRequestFailure(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	creds.SetLoginDataCallbackFn(func(a *AuthData) error { return nil })
	authKey := "testAuthKey"
	expectedErr := errors.New("problematic request")

	httpService.MockFn = func(req *http.Request) (*http.Response, error) {
		return nil, expectedErr
	}

	err := creds.getAccessToken(authKey)

	if err != expectedErr {
		test.Errorf("\nexpected %s\nbut got\n%s", expectedErr, err)
	}
}

func TestGetAccessTokenCallbackFailure(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	expectedErr := errors.New("problematic callback")
	creds.SetLoginDataCallbackFn(func(a *AuthData) error { return expectedErr })
	authKey := "testAuthKey"

	httpService.MockFn = httpService.DefaultMockFn
	err := creds.getAccessToken(authKey)

	if err != expectedErr {
		test.Errorf("\nexpected %s\nbut got\n%s", expectedErr, err)
	}
}

func TestGetRefreshTokenSuccess(test *testing.T) {
	var response *AuthData

	creds := Creds{}
	creds.SetAccessToken("refresh_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	creds.SetLoginDataCallbackFn(func(a *AuthData) error { response = a; return nil })

	err := creds.getRefreshToken()

	if err != nil {
		test.Errorf("\nexpected no errors\nbut got\n%s", err)
	}

	if len(response.TokenType) == 0 ||
		len(response.Scope) == 0 ||
		response.ExpiresIn == 0 ||
		response.ExtExpiresIn == 0 ||
		len(response.AccessToken) == 0 ||
		len(response.RefreshToken) == 0 {
		test.Error("\nexpected all AuthData values to be populated\nbut some were not\n")
	}
}

func TestGetRefreshTokenRequestFailure(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("refresh_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	creds.SetLoginDataCallbackFn(func(a *AuthData) error { return nil })
	expectedErr := errors.New("problematic request")

	httpService.MockFn = func(req *http.Request) (*http.Response, error) {
		return nil, expectedErr
	}

	err := creds.getRefreshToken()

	if err != expectedErr {
		test.Errorf("\nexpected %s\nbut got\n%s", expectedErr, err)
	}
}

func TestGetRefreshTokenCallbackFailure(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("refresh_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	expectedErr := errors.New("problematic callback")
	creds.SetLoginDataCallbackFn(func(a *AuthData) error { return expectedErr })

	httpService.MockFn = httpService.DefaultMockFn
	err := creds.getRefreshToken()

	if err != expectedErr {
		test.Errorf("\nexpected %s\nbut got\n%s", expectedErr, err)
	}
}

func TestProcessTokenRequestSuccess(test *testing.T) {
	request, _ := http.NewRequest("POST", "localhost/test", nil)
	var response *AuthData
	var expectedRequest *http.Request

	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))
	creds.SetLoginDataCallbackFn(func(a *AuthData) error { response = a; return nil })

	httpService.MockFn = func(r *http.Request) (*http.Response, error) {
		expectedRequest = r
		res := &http.Response{}
		res.Body = httpService.StubbedBody()
		return res, nil
	}

	err := creds.processTokenRequest(request)

	if err != nil {
		test.Errorf("\nexpected no errors\nbut got\n%s", err)
	}

	if request != expectedRequest {
		test.Errorf("\nexpected %v\nbut got\n%v", expectedRequest, request)
	}

	if len(response.TokenType) == 0 ||
		len(response.Scope) == 0 ||
		response.ExpiresIn == 0 ||
		response.ExtExpiresIn == 0 ||
		len(response.AccessToken) == 0 ||
		len(response.RefreshToken) == 0 {
		test.Error("\nexpected all AuthData values to be populated\nbut some were not\n")
	}
}

func TestProcessTokenRequestFailure(test *testing.T) {
	request, _ := http.NewRequest("POST", "localhost/test", nil)

	expectedErr := errors.New("problematic callback")

	creds := Creds{}
	creds.SetLoginDataCallbackFn(func(a *AuthData) error { return expectedErr })

	httpService.MockFn = func(r *http.Request) (*http.Response, error) {
		return nil, expectedErr
	}

	err := creds.processTokenRequest(request)

	if err != expectedErr {
		test.Errorf("\nexpected %s\nbut got\n%s", expectedErr, err)
	}
}

func TestAlreadyLoggedInWithAccessTokenSucess(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.alreadyLoggedIn()
	if !result {
		test.Errorf("\nexpected\ntrue\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenNoAccessToken(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenExpiredAccessToken(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(-fiveMins))

	result := creds.alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInWithRefreshTokenSucess(test *testing.T) {
	creds := Creds{}
	creds.SetRefreshToken("ref_token")
	creds.SetRefreshTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.alreadyLoggedIn()
	if !result {
		test.Errorf("\nexpected\ntrue\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenNoRefreshToken(test *testing.T) {
	creds := Creds{}
	creds.SetRefreshToken("")
	creds.SetRefreshTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestAlreadyLoggedInFailureWhenExpiredRefreshToken(test *testing.T) {
	creds := Creds{}
	creds.SetRefreshToken("ref_token")
	creds.SetRefreshTokenExpiresAt(time.Client.Now().Add(-fiveMins))

	result := creds.alreadyLoggedIn()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestIsAccessTokenValidSuccess(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.isAccessTokenValid()
	if !result {
		test.Errorf("\nexpected\ntrue\nbut got\n%v", result)
	}
}

func TestIsAccessTokenValidFailureForNoToken(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.isAccessTokenValid()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestIsAccessTokenValidFailureForExpiredToken(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(-fiveMins))

	result := creds.isAccessTokenValid()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestIsRefreshTokenValidSuccess(test *testing.T) {
	creds := Creds{}
	creds.SetRefreshToken("refresh_token")
	creds.SetRefreshTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.isRefreshTokenValid()
	if !result {
		test.Errorf("\nexpected\ntrue\nbut got\n%v", result)
	}
}

func TestIsRefreshTokenValidFailureForNoToken(test *testing.T) {
	creds := Creds{}
	creds.SetRefreshToken("")
	creds.SetRefreshTokenExpiresAt(time.Client.Now().Add(fiveMins))

	result := creds.isRefreshTokenValid()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestIsRefreshTokenValidFailureForExpiredToken(test *testing.T) {
	creds := Creds{}
	creds.SetRefreshToken("acc_token")
	creds.SetRefreshTokenExpiresAt(time.Client.Now().Add(-fiveMins))

	result := creds.isRefreshTokenValid()
	if result {
		test.Errorf("\nexpected\nfalse\nbut got\n%v", result)
	}
}

func TestRefreshTokenIfNeededWhenNeeded(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(-fiveMins))
	creds.SetLoginDataCallbackFn(func(a *AuthData) error { return nil })
	httpService.MockFn = httpService.DefaultMockFn

	err := creds.refreshTokenIfNeeded()

	if err != nil {
		test.Errorf("\nexpected\nno errors\nbut got\n%s", err)
	}
}

func TestRefreshTokenIfNeededWhenNotNeeded(test *testing.T) {
	creds := Creds{}
	creds.SetAccessToken("acc_token")
	creds.SetAccessTokenExpiresAt(time.Client.Now().Add(fiveMins))
	creds.SetLoginDataCallbackFn(func(a *AuthData) error { return nil })
	httpService.MockFn = httpService.DefaultMockFn

	err := creds.refreshTokenIfNeeded()

	if err != nil {
		test.Errorf("\nexpected\nno errors\nbut got\n%s", err)
	}
}

func TestBuildRequestBodyForAuthTokenSuccess(test *testing.T) {
	authKey := "testAuthKey"
	creds := Creds{}
	creds.SetId("testId")
	creds.SetSecret("testSecret")
	creds.SetPermissions("perms")
	creds.SetAuthCallbackHost("http://localhost")
	creds.SetAuthCallbackPath("/authorize")

	result := creds.buildRequestBodyForAuthToken(authKey)

	if result["client_id"][0] != "testId" ||
		result["scope"][0] != "perms" ||
		result["code"][0] != authKey ||
		result["redirect_uri"][0] != "http://localhost/authorize" ||
		result["grant_type"][0] != "authorization_code" ||
		result["client_secret"][0] != "testSecret" {
		test.Errorf("\nexpected\n%v\nbut got\n%v", creds, result)
	}
}

func TestBuildRequestBodyForRefreshTokenSuccess(test *testing.T) {
	creds := Creds{}
	creds.SetId("testId")
	creds.SetSecret("testSecret")
	creds.SetRefreshToken("refreshToken")

	result := creds.buildRequestBodyForRefreshToken()

	if result["client_id"][0] != "testId" ||
		result["grant_type"][0] != "refresh_token" ||
		result["client_secret"][0] != "testSecret" {
		test.Errorf("\nexpected\n%v\nbut got\n%v", creds, result)
	}
}

func TestBuildRequestObjectWithEncodedParamsSuccess(test *testing.T) {
	testUrl := "http://localhost/authorize"
	testData := url.Values{}
	testData.Set("some", "value")

	result, err := buildRequestObjectWithEncodedParams(testUrl, testData.Encode())

	result.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	result.Header.Add("Content-Length", strconv.Itoa(len(testData.Encode())))

	if err != nil {
		test.Errorf("expected:\nno errors\nbut got\n%s", err)
	}

	if result.Method != "POST" {
		test.Errorf("expected:\nPOST\nbut got\n%s", result.Method)
	}

	if result.URL.String() != testUrl {
		test.Errorf("expected:\n%s\nbut got\n%s", testUrl, result.URL)
	}

	if result.Header["Content-Type"][0] != "application/x-www-form-urlencoded" {
		test.Errorf("expected:\napplication/x-www-form-urlencoded\nbut got\n%s", result.Header)
	}

	body, _ := ioutil.ReadAll(result.Body)
	if string(body) != "some=value" {
		test.Errorf("expected:\nsome=value\nbut got\n%s", string(body))
	}
}

func TestSendRequestSuccess(test *testing.T) {
	testBody := "some=data"
	testUrl := "http://localhost"

	req, _ := http.NewRequest(
		"POST",
		testUrl,
		strings.NewReader(testBody),
	)

	var httpDoCalled bool

	httpService.MockFn = func(r *http.Request) (*http.Response, error) {
		httpDoCalled = true
		res := &http.Response{}
		res.Body = httpService.StubbedBody()
		return res, nil
	}

	_, err := sendRequest(req)

	if err != nil {
		test.Errorf("expected:\nno errors\nbut got\n%s", err)
	}

	if !httpDoCalled {
		test.Error("expected:\nto send request\nbut did not")
	}
}

func TestSendRequestFailure(test *testing.T) {
	testBody := "some=data"
	testUrl := "http://localhost"

	req, _ := http.NewRequest(
		"POST",
		testUrl,
		strings.NewReader(testBody),
	)

	var httpDoCalled bool
	expectedErr := errors.New("error on request")

	httpService.MockFn = func(r *http.Request) (*http.Response, error) {
		httpDoCalled = true
		return nil, expectedErr
	}

	_, err := sendRequest(req)

	if err != expectedErr {
		test.Errorf("expected:\n%s\nbut got\n%s", expectedErr, err)
	}

	if !httpDoCalled {
		test.Error("expected:\nto send request\nbut did not")
	}
}

func TestCallbackListenSuccess(test *testing.T) {
	var calledListenAndServe bool
	var calledCallbackUrl string
	var calledResponderFn func(w http.ResponseWriter, r *http.Request)
	var callbackFnCalledWith string

	testUrl := "/authorize"
	expectedCallbackFn := func(s string) error { callbackFnCalledWith = s; return nil }

	httpService.HandlerStubFn = func(
		url string,
		responder func(w http.ResponseWriter, r *http.Request),
	) {
		calledCallbackUrl = url
		calledResponderFn = responder
	}

	httpService.ListenAndServeStubFn = func(s string, h http.Handler) error {
		calledListenAndServe = true
		return nil
	}

	err := CallbackListen(testUrl, expectedCallbackFn)

	if err != nil {
		test.Errorf("expected:\nno errors\nbut got\n%s", err)
	}

	if _ = callbackFn("some"); callbackFnCalledWith != "some" {
		test.Error("expected:\nto have assigned the callback function\nbut\nit did not")
	}

	if calledCallbackUrl != testUrl {
		test.Error("expected:\nto handle callbackUrl\nbut\nit did not")
	}

	if calledResponderFn == nil {
		test.Error("expected:\nto have assigned responder function\nbut\nit did not")
	}

	if !calledListenAndServe {
		test.Error("expected:\nto have called ListenAndServe\nbut\nit did not")
	}
}

func TestResponderSuccess(test *testing.T) {
	var calledCallbackFnWith string
	code := "testCode"

	callbackFn = func(s string) error {
		calledCallbackFnWith = s
		return nil
	}

	req, _ := http.NewRequest("GET", "localhost", nil)
	q := req.URL.Query()
	q.Add("code", code)

	req.URL.RawQuery = q.Encode()

	res := httpService.StubbedResponseWriter()

	responder(res, req)

	if calledCallbackFnWith != code {
		test.Errorf("expected:\nto have called callbackFn with %s\nbut\nit did not", code)
	}
}

func TestResponderFailureWithErrorInCallbackFn(test *testing.T) {
	callbackFn = func(s string) error {
		return errors.New("error in callbackFn")
	}

	req, _ := http.NewRequest("GET", "localhost", nil)
	q := req.URL.Query()
	q.Add("code", "test")

	req.URL.RawQuery = q.Encode()

	res := httpService.StubbedResponseWriter()

	responder(res, req)

	resStr := fmt.Sprint(res)
	if !strings.Contains(resStr, "error in callbackFn") {
		test.Error("expected:\nto have error in response\nbut\nit did not")
	}
}

func TestResponderFailureWithEmptyQuery(test *testing.T) {
	var calledCallbackFn bool

	callbackFn = func(s string) error {
		calledCallbackFn = true
		return nil
	}

	req, _ := http.NewRequest("GET", "localhost", nil)
	res := httpService.StubbedResponseWriter()

	responder(res, req)

	if calledCallbackFn {
		test.Error("expected:\nto not have called callbackFn\nbut\nit did")
	}
}

func TestPrepareLoginUrl(test *testing.T) {
	creds := Creds{}
	creds.SetId("client-id")
	creds.SetPermissions("Some.Permissions")
	creds.SetAuthCallbackHost("http://localhost:3000")
	creds.SetAuthCallbackPath("/callback")

	baseRequestUrl = "http://test.com"
	authRequestPath = "/authorize"

	result := creds.prepareLoginUrl()
	expectedResult := "http://test.com/authorize?client_id=client-id" +
		"&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback" +
		"&response_mode=query&response_type=code&scope=Some.Permissions"

	if result != expectedResult {
		test.Errorf("\nexpected\n%s\nbut got\n%s", expectedResult, result)
	}
}

//
// // TODO: Create test helpers for common parts in different
// // variants of tests for this function
// func TestOpenLoginUrlForLinuxSuccess(test *testing.T) {
// 	runtime.Client = runtimetest.RuntimeMock{}
// 	exec.CmdClient = &exectest.CmdMock{}
// 	exec.RunClient = &exectest.CmdRunMock{}
// 	log.Client = logtest.LoggerServiceMock{}
//
// 	testUrl := "http://localhost/open"
// 	testCmd := osexec.Command("echo", testUrl)
//
// 	var ranCmd bool
// 	var ranRunCmd bool
//
// 	runtimetest.RuntimeMockFunc = func() string { return "linux" }
//
// 	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
// 		ranCmd = true
// 		return exectest.InitCmdMock(testCmd)
// 	}
// 	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
// 		ranRunCmd = true
// 		return nil
// 	}
//
// 	var logResult interface{}
// 	logtest.FatalMock = func(s ...interface{}) { logResult = s }
//
// 	// openLoginUrl(testUrl)
// 	// if !ranCmd && !ranRunCmd {
// 	// 	test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
// 	// }
// 	//
// 	// if logResult != nil {
// 	// 	test.Errorf(
// 	// 		"\nexpected\nto run without errors\nbut was\n%s",
// 	// 		fmt.Sprintf("%v", logResult),
// 	// 	)
// 	// }
// }
//
// func TestOpenLoginUrlForLinuxFailure(test *testing.T) {
// 	runtime.Client = runtimetest.RuntimeMock{}
// 	exec.CmdClient = &exectest.CmdMock{}
// 	exec.RunClient = &exectest.CmdRunMock{}
// 	log.Client = logtest.LoggerServiceMock{}
//
// 	testUrl := "http://localhost/open"
// 	testCmd := osexec.Command("echo", testUrl)
//
// 	var ranCmd bool
// 	var ranRunCmd bool
//
// 	expectedErr := errors.New("could not run")
// 	runtimetest.RuntimeMockFunc = func() string { return "linux" }
//
// 	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
// 		ranCmd = true
// 		return exectest.InitCmdMock(testCmd)
// 	}
//
// 	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
// 		ranRunCmd = true
// 		return expectedErr
// 	}
//
// 	var logResult error
// 	logtest.FatalMock = func(s ...interface{}) { logResult = expectedErr }
//
// 	// openLoginUrl(testUrl)
//
// 	if !ranCmd && !ranRunCmd {
// 		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
// 	}
//
// 	if logResult != expectedErr {
// 		test.Errorf(
// 			"\nexpected\nto run with error\n%s\nbut was\n%s",
// 			expectedErr.Error(),
// 			fmt.Sprintf("%v", logResult),
// 		)
// 	}
// }
//
// func TestOpenLoginUrlForWindowsSuccess(test *testing.T) {
// 	runtime.Client = runtimetest.RuntimeMock{}
// 	exec.CmdClient = &exectest.CmdMock{}
// 	exec.RunClient = &exectest.CmdRunMock{}
// 	log.Client = logtest.LoggerServiceMock{}
//
// 	testUrl := "http://localhost/open"
// 	testCmd := osexec.Command("echo", testUrl)
//
// 	var ranCmd bool
// 	var ranRunCmd bool
//
// 	runtimetest.RuntimeMockFunc = func() string { return "windows" }
//
// 	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
// 		ranCmd = true
// 		return exectest.InitCmdMock(testCmd)
// 	}
// 	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
// 		ranRunCmd = true
// 		return nil
// 	}
//
// 	var logResult interface{}
// 	logtest.FatalMock = func(s ...interface{}) { logResult = s }
//
// 	// openLoginUrl(testUrl)
// 	if !ranCmd && !ranRunCmd {
// 		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
// 	}
//
// 	if logResult != nil {
// 		test.Errorf(
// 			"\nexpected\nto run without errors\nbut was\n%s",
// 			fmt.Sprintf("%v", logResult),
// 		)
// 	}
// }
//
// func TestOpenLoginUrlForWindowsFailure(test *testing.T) {
// 	runtime.Client = runtimetest.RuntimeMock{}
// 	exec.CmdClient = &exectest.CmdMock{}
// 	exec.RunClient = &exectest.CmdRunMock{}
// 	log.Client = logtest.LoggerServiceMock{}
//
// 	testUrl := "http://localhost/open"
// 	testCmd := osexec.Command("echo", testUrl)
//
// 	var ranCmd bool
// 	var ranRunCmd bool
//
// 	expectedErr := errors.New("could not run")
// 	runtimetest.RuntimeMockFunc = func() string { return "windows" }
//
// 	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
// 		ranCmd = true
// 		return exectest.InitCmdMock(testCmd)
// 	}
//
// 	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
// 		ranRunCmd = true
// 		return expectedErr
// 	}
//
// 	var logResult error
// 	logtest.FatalMock = func(s ...interface{}) { logResult = expectedErr }
//
// 	// openLoginUrl(testUrl)
//
// 	if !ranCmd && !ranRunCmd {
// 		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
// 	}
//
// 	if logResult != expectedErr {
// 		test.Errorf(
// 			"\nexpected\nto run with error\n%s\nbut was\n%s",
// 			expectedErr.Error(),
// 			fmt.Sprintf("%v", logResult),
// 		)
// 	}
// }
//
// func TestOpenLoginUrlForMacOSSuccess(test *testing.T) {
// 	runtime.Client = runtimetest.RuntimeMock{}
// 	exec.CmdClient = &exectest.CmdMock{}
// 	exec.RunClient = &exectest.CmdRunMock{}
// 	log.Client = logtest.LoggerServiceMock{}
//
// 	testUrl := "http://localhost/open"
// 	testCmd := osexec.Command("echo", testUrl)
//
// 	var ranCmd bool
// 	var ranRunCmd bool
//
// 	runtimetest.RuntimeMockFunc = func() string { return "darwin" }
//
// 	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
// 		ranCmd = true
// 		return exectest.InitCmdMock(testCmd)
// 	}
// 	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
// 		ranRunCmd = true
// 		return nil
// 	}
//
// 	var logResult interface{}
// 	logtest.FatalMock = func(s ...interface{}) { logResult = s }
//
// 	// openLoginUrl(testUrl)
// 	if !ranCmd && !ranRunCmd {
// 		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
// 	}
//
// 	if logResult != nil {
// 		test.Errorf(
// 			"\nexpected\nto run without errors\nbut was\n%s",
// 			fmt.Sprintf("%v", logResult),
// 		)
// 	}
// }
//
// func TestOpenLoginUrlForMacOSFailure(test *testing.T) {
// 	runtime.Client = runtimetest.RuntimeMock{}
// 	exec.CmdClient = &exectest.CmdMock{}
// 	exec.RunClient = &exectest.CmdRunMock{}
// 	log.Client = logtest.LoggerServiceMock{}
//
// 	testUrl := "http://localhost/open"
// 	testCmd := osexec.Command("echo", testUrl)
//
// 	var ranCmd bool
// 	var ranRunCmd bool
//
// 	expectedErr := errors.New("could not run")
// 	runtimetest.RuntimeMockFunc = func() string { return "darwin" }
//
// 	exectest.CommandMockFunc = func(name string, arg ...string) exec.RunService {
// 		ranCmd = true
// 		return exectest.InitCmdMock(testCmd)
// 	}
//
// 	exectest.CommandRunMockFunc = func(c *exectest.CmdRunMock) error {
// 		ranRunCmd = true
// 		return expectedErr
// 	}
//
// 	var logResult error
// 	logtest.FatalMock = func(s ...interface{}) { logResult = expectedErr }
//
// 	// openLoginUrl(testUrl)
//
// 	if !ranCmd && !ranRunCmd {
// 		test.Errorf("\nexpected\n%s\nto run but it did not", testCmd)
// 	}
//
// 	if logResult != expectedErr {
// 		test.Errorf(
// 			"\nexpected\nto run with error\n%s\nbut was\n%s",
// 			expectedErr.Error(),
// 			fmt.Sprintf("%v", logResult),
// 		)
// 	}
// }
//
// func TestOpenLoginUrlForUnknownSuccess(test *testing.T) {
// 	runtime.Client = runtimetest.RuntimeMock{}
// 	log.Client = logtest.LoggerServiceMock{}
//
// 	runtimetest.RuntimeMockFunc = func() string { return "unknown" }
//
// 	expectedLog := "visit some url"
// 	var logResult interface{}
// 	logtest.PrintfMock = func(s string, v ...interface{}) { logResult = expectedLog }
//
// 	testUrl := "http://localhost/open"
// 	// openLoginUrl(testUrl)
//
// 	if logResult == nil {
// 		test.Errorf(
// 			"\nexpected to log\n\"%s\"\nbut nothing was logged\n",
// 			expectedLog,
// 		)
// 	}
// }
