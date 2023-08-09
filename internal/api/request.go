package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ErrUnauthorized occurs when the API response indicates an invalid or expired RAT token
type ErrUnauthorized struct {
	Err error
}

func (e ErrUnauthorized) Error() string {
	return "RAT token invalid or expired"
}

func (e ErrUnauthorized) Unwrap() error {
	return e.Err
}

var client = &http.Client{
	Timeout: 5 * time.Second,
	CheckRedirect: func(req *http.Request, _ []*http.Request) error {
		// disallow redirects
		return http.ErrUseLastResponse
	},
}

// ErrHTTP occurs when a non-200 HTTP response is received
type ErrHTTP struct {
	StatusCode int
}

func (e ErrHTTP) Error() string {
	return fmt.Sprintf("status %d", e.StatusCode)
}

// ErrResponseDecode occurs when a valid response is received, but it could not be decoded
type ErrResponseDecode struct {
	Err error
}

func (e ErrResponseDecode) Error() string {
	return "Could not decode API response"
}

func (e ErrResponseDecode) Unwrap() error {
	return e.Err
}

func (c *Client) get(url string, output any) (err error) {
	var req *http.Request
	req, err = http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "github.com/stnokott/sot-api / v0.0.0")
	req.Header.Set("Referer", c.httpReferer)
	if c.token == "" {
		err = ErrUnauthorized{Err: errors.New("no RAT token provided")}
		return
	}
	req.AddCookie(&http.Cookie{
		Name:  "rat",
		Value: c.token,
	})
	c.logger.Debug("GET " + url)
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = handleInvalidResponse(resp)
		return
	}
	defer func() {
		if errInner := resp.Body.Close(); errInner != nil && err == nil {
			err = errInner
		}
	}()

	c.logger.Debug("decoding response")
	if err = json.NewDecoder(resp.Body).Decode(output); err != nil {
		err = ErrResponseDecode{Err: err}
	}
	return
}

func handleInvalidResponse(resp *http.Response) error {
	if resp.StatusCode == 302 {
		if location := resp.Header.Get("Location"); location == "/logout" {
			return ErrUnauthorized{
				Err: errors.New("redirected to logout page, please update your RAT token"),
			}
		}
	}
	return ErrHTTP{StatusCode: resp.StatusCode}
}

// apiGet creates a new authenticated (using RAT token) HTTP GET request to the specified seaofthieves.com endpoint
func (c *Client) apiGet(endpoint string, output any) error {
	return c.get(c.baseURL+endpoint, output)
}
