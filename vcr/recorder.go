package vcr

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func createCassetteFile(filename string) string {
	if !exists("fixtures") {
		if os.MkdirAll("fixtures", 0755) != nil {
			panic("Unable to create directory for tagfile!")
		}
	}

	if !exists(filename) {
		file, err := os.Create(filename)
		if err != nil {
			panic("Unable to create tag file!")
		}
		defer file.Close()
	}

	return filename
}

func requestHandler(r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header = r.Header
	return client.Do(req)
}

func useCassette(filename string) cassetteData {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	cassette := cassetteData{}
	json.Unmarshal(file, &cassette)

	return cassette
}

func newCassette(r *http.Request, resp *http.Response, filename string) cassetteData {
	cassette := cassetteData{}

	cassette.Request.Method = r.Method
	cassette.Request.URL = r.URL.String()
	cassette.Request.Header = r.Header
	cassette.Request.Values = r.PostForm

	cassette.Response.Header = resp.Header
	cassette.Response.Status = resp.Status
	cassette.Response.StatusCode = resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	cassette.Response.Body = string(body)

	data, err := json.Marshal(cassette)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filename, data, 0644)

	return cassette
}

func recordCassette(r *http.Request, filename string) cassetteData {
	createCassetteFile(filename)

	resp, err := requestHandler(r)
	if err != nil {
		log.Fatal(err)
	}

	cassette := newCassette(r, resp, filename)
	return cassette
}
