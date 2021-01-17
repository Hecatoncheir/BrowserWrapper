package chrome

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type Tab struct {
	TabContext       context.Context
	TabContextCancel context.CancelFunc
}

func NewTabForChrome(chromeContext context.Context) (*Tab, error) {
	newContext, _ := chromedp.NewContext(chromeContext)

	err := chromedp.Run(newContext)
	if err != nil {
		return nil, err
	}

	tabContext := chromedp.FromContext(newContext)

	newTabContext, newTabContextCancel := chromedp.NewContext(
		chromeContext, chromedp.WithTargetID(tabContext.Target.TargetID),
	)

	tab := Tab{
		TabContext:       newTabContext,
		TabContextCancel: newTabContextCancel,
	}

	return &tab, nil
}

func (tab *Tab) GetHtmlOfPage(url string) (string, error) {
	return tab.WaitElementAndGetHtmlOfPage(url, "body")
}

func (tab *Tab) WaitElementAndGetHtmlOfPage(pageUrl, elementSelector string) (string, error) {
	var html string

	err := chromedp.Run(
		tab.TabContext,
		tab.RunWithTimeOut(&tab.TabContext, time.Second*30, chromedp.Tasks{
			chromedp.Navigate(pageUrl),
			chromedp.WaitReady(elementSelector),
			chromedp.OuterHTML("html", &html),
		}),
	)
	if err != nil {
		err := chromedp.Run(
			tab.TabContext,
			tab.RunWithTimeOut(&tab.TabContext, time.Second*30, chromedp.Tasks{
				chromedp.Navigate(pageUrl),
				chromedp.WaitReady(elementSelector),
				chromedp.OuterHTML("html", &html),
			}),
		)

		return "", err
	}

	return html, nil
}

func (tab *Tab) RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

func (tab *Tab) Close() error {
	err := chromedp.Cancel(tab.TabContext)
	if err != nil {
		return err
	}

	tab.TabContextCancel()
	return nil
}
