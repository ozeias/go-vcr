// Package vcr provides the interface and implementation to record your test suite's HTTP interactions and replay them during future test runs for fast, deterministic, accurate tests.
package vcr

import (
	"fmt"

	"net/http"
	"net/http/httptest"
	"net/url"
)

const (
	libraryVersion = "0.1"
)

func verifyRequest(r *http.Request, filename string) cassetteData {
	filename = fmt.Sprintf("fixtures/%v.json", filename)

	if exists(filename) {
		return useCassette(filename)
	}

	return recordCassette(r, filename)
}

/*
UseCassette is responsible to mock a http.Client,
and server setup(httptest.Server) to respond with a specific
status code and body recorded in the cassette.

The client must close the server when finished with it:
  ...
  server, httpClient := vcr.UseCassette("vine")
  client.HTTPClient = httpClient
  defer server.Close()
  // ...
*/
func UseCassette(filename string) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cassette := verifyRequest(r, filename)

		w.WriteHeader(cassette.Response.StatusCode)
		w.Header().Set("Content-Type", cassette.Response.Header.Get("Content-Type"))

		fmt.Fprintln(w, cassette.Response.Body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: tr}
	return server, httpClient
}
