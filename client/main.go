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

	ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(90 * time.Second)
		done <- true
	}()
	i := 1
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			HTTPRequest(c, req, t, i)
			i += 1
		}
	}

	// for i := 0; i < 90; i++ {
	// 	// res, err := c.Client.Do(req)
	// 	// if err != nil {
	// 	// 	fmt.Println(err.Error())
	// 	// 	fmt.Println(res.StatusCode)
	// 	// 	return
	// 	// }

	// 	// if err := decodeBody(res, &response); err != nil {
	// 	// 	fmt.Println(err.Error())
	// 	// 	return
	// 	// }
	// 	go HTTPRequest(c, req, start, i)

	// 	// if i == 29 || i == 59 {
	// 	// 	time.Sleep(time.Second * 30)
	// 	// }
	// 	time.Sleep(time.Second * 1)
	// }
}

func HTTPRequest(c *RLClient, req *http.Request, t time.Time, i int) {
	var response Response
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
	fmt.Println(i, response, res.StatusCode, t)
}
