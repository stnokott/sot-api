package api

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	authURL    = "https://www.seaofthieves.com/de/login"
	targetHost = "www.seaofthieves.com"
)

// AuthResp contains the results of the login process.
// Token will be the empty string if Err is not nil.
type AuthResp struct {
	Token string
	Err   error
}

var chromedpOpts = append(chromedp.DefaultExecAllocatorOptions[:],
	chromedp.NoFirstRun,
	chromedp.NoDefaultBrowserCheck,
	chromedp.DisableGPU,
	chromedp.WindowSize(800, 600),
	chromedp.Flag("headless", false),
	chromedp.Flag("enable-automation", false),
	chromedp.Flag("disable-features", "ExtensionsToolbarMenu"),
)

// GetAuthFromBrowser opens a new browser window to get a new RAT token
func GetAuthFromBrowser() <-chan AuthResp {
	chCookie := make(chan AuthResp)

	go func() {
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), chromedpOpts...)
		defer cancel()

		taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		defer cancel()

		fnCheckLocation := func(ctx context.Context) error {
			var currentURL string
			for {
				if err := chromedp.Location(&currentURL).Do(ctx); err != nil {
					return err
				}

				url, err := url.Parse(currentURL)
				if err != nil {
					return err
				}
				if url.Host == targetHost {
					cookies, err := network.GetCookies().Do(ctx)
					if err != nil {
						return err
					}
					chCookie <- filterRatCookie(cookies)
					close(chCookie)
					return nil
				}

				time.Sleep(1 * time.Second)
			}
		}

		if err := chromedp.Run(taskCtx,
			chromedp.Navigate(authURL),
			chromedp.WaitReady("body"),
			chromedp.ActionFunc(fnCheckLocation),
		); err != nil {
			log.Fatalln(err)
		}
	}()

	return chCookie
}

func filterRatCookie(cookies []*network.Cookie) (r AuthResp) {
	for _, cookie := range cookies {
		if cookie.Name == "rat" {
			r.Token = cookie.Value
			return
		}
	}
	r.Err = errors.New("RAT cookie not found")
	return
}
