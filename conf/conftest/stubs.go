package conftest

import (
	"github.com/betasve/mstd/viper/vipertest"
	"time"
)

func StubGetAccessToken(v string) {
	vipertest.GetStringFunc = func(_ string) string {
		return v
	}
}

func StubGetRefreshToken(v string) {
	vipertest.GetStringFunc = func(_ string) string {
		return v
	}
}

func StubGetClientAccessTokenExpiry(t time.Time) {
	vipertest.GetInt64Func = func(_ string) int64 {
		return t.Unix()
	}
}

func StubGetClientRefreshTokenExpiry(t time.Time) {
	vipertest.GetInt64Func = func(_ string) int64 {
		return t.Unix()
	}
}
