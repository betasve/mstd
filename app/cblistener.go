package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var callbackFn func(string) string

func CallbackListen(callbackUrl string, cb func(string) string) {
	callbackFn = cb

	cb("test")

	http.HandleFunc(callbackUrl, responder)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Panic(err)
	}
}

func responder(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	code := strings.Join(values["code"], "")

	fmt.Fprint(
		w,
		"Successfully retrieved an authorization "+
			"code \nContinuing with authorization process",
	)

	if code == "" {
		return
	}

	callbackFn(code)
}
