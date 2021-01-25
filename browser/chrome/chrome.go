package chrome

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/Hecatoncheir/BrowserWrapper/browser"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Chrome struct {
	isClosed            bool
	chromeContext       context.Context
	chromeContextCancel context.CancelFunc

	tabsMutex sync.Mutex
	Tabs      map[int]browser.TabInterface

	Downloader *Downloader
}

func (chrome *Chrome) GetPage(url string) (*goquery.Document, error) {
	return chrome.WaitElementAndGetPage(url, "body")
}

func (chrome *Chrome) WaitElementAndGetPage(pageUrl, elementSelector string) (*goquery.Document, error) {
	html, err := chrome.WaitElementAndGetHtmlOfPage(pageUrl, elementSelector)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (chrome *Chrome) Download(url string) ([]byte, error) {
	bytes, err := chrome.Downloader.Download(url)
	if err != nil || len(bytes) == 0 {
		return chrome.Downloader.Download(url)
	}
	return bytes, err
}

func (chrome *Chrome) GetHtmlOfPage(url string) (string, error) {
	return chrome.WaitElementAndGetHtmlOfPage(url, "body")
}

func (chrome *Chrome) GetTab(tabNumber int) (tab browser.TabInterface, err error) {
	tab = chrome.Tabs[tabNumber]
	if tab != nil {
		return tab, nil
	} else {
		return nil, errors.New("not tab found")
	}
}

func (chrome *Chrome) OpenTab() (tab browser.TabInterface, tabNumber int, err error) {
	chrome.tabsMutex.Lock()

	tab, err = NewTabForChrome(chrome.chromeContext)
	if err != nil {
		return nil, 0, err
	}

	tabNumber = len(chrome.Tabs) + 1
	chrome.Tabs[tabNumber] = tab

	chrome.tabsMutex.Unlock()
	return tab, tabNumber, nil
}

func (chrome *Chrome) CloseTab(tabNumber int) error {
	chrome.tabsMutex.Lock()

	tab, err := chrome.GetTab(tabNumber)
	if err != nil {
		return err
	}

	err = tab.Close()
	if err != nil {
		return err
	}

	delete(chrome.Tabs, tabNumber)
	chrome.tabsMutex.Unlock()

	return nil
}

func (chrome *Chrome) WaitElementAndGetHtmlOfPage(pageUrl, elementSelector string) (string, error) {

	tab, tabNumber, err := chrome.OpenTab()
	if err != nil {
		return "", err
	}

	html, err := tab.WaitElementAndGetHtmlOfPage(pageUrl, elementSelector)
	if err != nil {
		return "", err
	}

	err = chrome.CloseTab(tabNumber)
	if err != nil {
		return "", err
	}

	return html, nil
}

func (chrome *Chrome) DownloadFileBySelector(pageUrl, selector, attribute string) (bytes []byte, fileName, fileExtension, fileUrl string, err error) {
	page, err := chrome.WaitElementAndGetPage(pageUrl, selector)
	if err != nil {
		return nil, "", "", "", err
	}

	fileSrc, isExists := page.Find(selector).Attr(attribute)
	if !isExists {
		err := browser.ElementBySelectorDoesNotExist
		return nil, "", "", "", err
	}

	bytes, err = chrome.Download(fileSrc)
	if err != nil {
		return nil, "", "", "", err
	}

	fileSrcSegments := strings.Split(fileSrc, "/")
	fileName = fileSrcSegments[len(fileSrcSegments)-1]

	fileWithExtension := strings.Split(fileName, ".")
	fileExtension = fileWithExtension[len(fileWithExtension)-1]

	return bytes, fileName, fileExtension, fileSrc, nil
}

func (chrome *Chrome) IsClosed() bool {
	return chrome.isClosed
}

func (chrome *Chrome) Close() error {
	if chrome.isClosed {
		return nil
	}

	for tabNumber := range chrome.Tabs {
		err := chrome.CloseTab(tabNumber)
		if err != nil {
			return err
		}
	}

	err := chromedp.Cancel(chrome.chromeContext)
	if err != nil {
		return err
	}

	chrome.chromeContextCancel()

	chrome.isClosed = true

	return nil
}
