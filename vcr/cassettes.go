package vcr

import (
	"net/http"
	"net/url"
)

type cassetteData struct {
	Request struct {
		Method string      `json:"method"`
		URL    string      `json:"url"`
		Header http.Header `json:"header"`
		Values url.Values  `json:"values"`
	} `json:"request"`

	Response struct {
		Status     string      `json:"status"`
		StatusCode int         `json:"status_code"`
		Header     http.Header `json:"header"`
		Body       string      `json:"body"`
	} `json:"response"`
}
