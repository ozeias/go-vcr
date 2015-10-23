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

// Forces all requests to be "http" (see below)
type httpsRewrite struct {
	transport   http.RoundTripper
	originalReq *http.Request
}

func (r *httpsRewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	r.originalReq = &http.Request{}
	// copy the values
	*r.originalReq = *req
	*r.originalReq.URL = *req.URL

	req.URL.Scheme = "http"
	return r.transport.RoundTrip(req)
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
	// Prepare for grossness.
	// Requests are made to the httptest server over http.
	// However, you may be actually making a https request for the cassette.
	// Even if that request is fine and writes back a response when this completes
	// the RoundTripper will have seen that the request was an https and attempt to decode it resulting in a "tls: oversized record received with length 24864" error.
	// So we'll copy the original request that is actually made and force the current request to not be https.
	var httpsTr *httpsRewrite

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cassette := verifyRequest(httpsTr.originalReq, filename)

		w.WriteHeader(cassette.Response.StatusCode)
		w.Header().Set("Content-Type", cassette.Response.Header.Get("Content-Type"))

		fmt.Fprintln(w, cassette.Response.Body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpsTr = &httpsRewrite{transport: tr}
	httpClient := &http.Client{Transport: httpsTr}
	return server, httpClient
}
