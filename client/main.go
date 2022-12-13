package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

// https://deeeet.com/writing/2016/11/01/go-api-client/
type RLClient struct {
	URL         *url.URL
	Client      *http.Client
	RateLimiter *rate.Limiter
}

func NewRLClient(urlStr string, rateLimit *rate.Limiter) (*RLClient, error) {
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", urlStr)
	}
	return &RLClient{
		URL:         parsedURL,
		Client:      http.DefaultClient,
		RateLimiter: rateLimit,
	}, nil
}

func (c *RLClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.RateLimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	baseURL := "http://localhost:8080"
	c, err := NewRLClient(baseURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	var response Response

	for i := 0; i < 20; i++ {
		res, err := c.Client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(res.StatusCode)
			return
		}

		if err := decodeBody(res, &response); err != nil {
			fmt.Println(err.Error())
			return
		}
		time.Sleep(time.Millisecond * 100)
		fmt.Println(i+1, response, res.StatusCode, time.Now())
	}

}
