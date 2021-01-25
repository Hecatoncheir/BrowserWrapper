package chrome

import (
	"context"
	"github.com/Hecatoncheir/BrowserWrapper/browser"

	"github.com/chromedp/chromedp"
)

func New(proxyDetails string, isHeadless bool) *Chrome {
	chromeContext, chromeContextCancel := context.WithCancel(context.Background())

	var options []chromedp.ExecAllocatorOption

	options = append(chromedp.DefaultExecAllocatorOptions[:])
	options = append(options,
		chromedp.DisableGPU,
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("restore-on-startup", false),
	)

	if proxyDetails != "" {
		options = append(options,
			chromedp.ProxyServer(proxyDetails),
		)
	}

	if isHeadless {
		options = append(options,
			chromedp.Headless,
		)
	} else {
		options = append(options,
			chromedp.Flag("headless", false),
		)
	}

	chromeContext, chromeContextCancel = chromedp.NewExecAllocator(chromeContext, options...)
	chromeContext, chromeContextCancel = chromedp.NewContext(chromeContext)

	err := chromedp.Run(chromeContext)
	if err != nil {
		println("Chrome NewWithProxy error:")
		println(err)
	}

	return &Chrome{
		isClosed:            false,
		chromeContext:       chromeContext,
		chromeContextCancel: chromeContextCancel,
		Downloader:          NewDownloader(proxyDetails, chromeContext),
		Tabs:                map[int]browser.TabInterface{},
	}
}
