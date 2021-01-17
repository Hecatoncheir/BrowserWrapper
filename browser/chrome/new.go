package chrome

import (
	"context"
	"github.com/Hecatoncheir/BrowserWrapper/browser"

	"github.com/chromedp/chromedp"
)

func New(proxyDetails string) *Chrome {
	chromeContext, chromeContextCancel := context.WithCancel(context.Background())

	var options []chromedp.ExecAllocatorOption

	if proxyDetails != "" {
		options = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.DisableGPU,
			chromedp.ProxyServer(proxyDetails),
			chromedp.Headless,
			//chromedp.Flag("headless", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("restore-on-startup", false),
		)

	} else {
		options = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.DisableGPU,
			//chromedp.Headless,
			chromedp.Flag("headless", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("restore-on-startup", false),
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
		chromeContext:       chromeContext,
		chromeContextCancel: chromeContextCancel,
		Downloader:          NewDownloader(proxyDetails, chromeContext),
		Tabs:                map[int]browser.TabInterface{},
	}
}
