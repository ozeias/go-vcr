package vcr_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"testing"

	"github.com/ozeias/go-vcr/vcr"
)

const (
	VineAPI = "http://api.vineapp.com"
)

var (
	client *Client
)

type Vine struct {
	Code string `json:"code"`
	Data struct {
		Count   int64 `json:"count"`
		Records []struct {
			Description string `json:"description"`
			Likes       struct {
				Count int64 `json:"count"`
			} `json:"likes"`
			PermalinkURL string `json:"permalinkUrl"`
			ShareURL     string `json:"shareUrl"`
			Username     string `json:"username"`
			VideoURL     string `json:"videoUrl"`
		} `json:"records"`
		Size int `json:"size"`
	} `json:"data"`
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    *url.URL
}

func (c Client) getVine(id string) (vine *Vine, err error) {
	vine = &Vine{}
	urlStr := fmt.Sprintf("%v/timelines/posts/%v", c.BaseURL, id)
	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(vine)
	return vine, err
}

func setup() {
	baseURL, _ := url.Parse(VineAPI)

	client = &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    baseURL,
	}
}

func Test_UseCassette(t *testing.T) {
	setup()
	vineID := "1127994498119589888"

	server, httpClient := vcr.UseCassette("vine")
	client.HTTPClient = httpClient
	defer server.Close()

	vine, err := client.getVine(vineID)
	if err != nil {
		t.Errorf("getVine returned error: %s", err.Error())
	}

	if !vine.Success {
		t.Errorf("Expect response to be success and received %t", vine.Success)
	}

	if vine.Data.Records[0].Likes.Count != 343156 {
		t.Errorf("Vine %d contained incorrect count. Received: %d", vineID, vine.Data.Records[0].Likes.Count)
	}

	if vine.Data.Records[0].Description != "Do you think I'm cute? Yes or no?" {
		t.Errorf("Vine %d contained incorrect text. Received: %s", vineID, vine.Data.Records[0].Description)
	}
}
