package chrome

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Downloader struct {
	Client  *http.Client
	Context context.Context
}

func NewDownloader(proxy string, chromeContext context.Context) *Downloader {
	httpClient := prepareHttpClient(proxy)
	downloader := Downloader{
		Client:  httpClient,
		Context: chromeContext,
	}

	return &downloader
}

func (downloader *Downloader) Download(url string) ([]byte, error) {
	tab, err := NewTabForChrome(downloader.Context)
	if err != nil {
		return nil, err
	}

	done := make(chan bool)

	var requestId network.RequestID

	chromedp.ListenTarget(tab.TabContext, func(ev interface{}) {

		switch ev := ev.(type) {

		case *network.EventRequestWillBeSent:
			req := ev.Request
			if req.URL == url {
				requestId = ev.RequestID
			}

		case *network.EventLoadingFinished:
			if ev.RequestID == requestId {
				close(done)
			}
		}
	})

	err = chromedp.Run(tab.TabContext, chromedp.Tasks{
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorAllow).WithDownloadPath(os.TempDir()),
		chromedp.Navigate(url),
	})
	if err != nil {
		return nil, err
	}

	<-done

	var bytes []byte
	err = chromedp.Run(tab.TabContext, chromedp.ActionFunc(func(cxt context.Context) error {
		bytes, err = network.GetResponseBody(requestId).Do(cxt)
		return err
	}))

	err = tab.Close()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (downloader *Downloader) DownloadExternal(url string) ([]byte, error) {

	response, err := downloader.Client.Get(url)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func prepareHttpClient(proxy string) *http.Client {
	var client http.Client

	if proxy != "" {
		proxyUrl, _ := url.Parse(proxy)
		client = http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			//Timeout: 30 * time.Second,
		}
	} else {
		client = http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
			Timeout: 15 * time.Second,
		}
	}

	return &client
}
