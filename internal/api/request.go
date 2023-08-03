package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 5 * time.Second}

// ErrHTTP occurs when a non-200 HTTP response is received
type ErrHTTP struct {
	StatusCode int
	Err        error
}

func (r ErrHTTP) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

// ErrResponseDecode occurs when a valid response is received, but it could not be decoded
type ErrResponseDecode error

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
	req.AddCookie(&http.Cookie{
		Name:  "rat",
		Value: c.token,
	})
	c.logger.Debug("GET " + url)
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = ErrHTTP{
			StatusCode: res.StatusCode,
			Err:        errors.New("got non-ok status code from API"),
		}
		return
	}
	defer func() {
		if errInner := res.Body.Close(); errInner != nil && err == nil {
			err = errInner
		}
	}()

	c.logger.Debug("decoding response")
	if err = json.NewDecoder(res.Body).Decode(output); err != nil {
		err = err.(ErrResponseDecode)
	}
	return
}

// apiGet creates a new authenticated (using RAT token) HTTP GET request to the specified seaofthieves.com endpoint
func (c *Client) apiGet(endpoint string, output any) error {
	return c.get(c.baseURL+endpoint, output)
}
