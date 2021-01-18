# BrowserWrapper

```go
type Interface interface {
	GetPage(url string) (*goquery.Document, error)
	WaitElementAndGetPage(pageUrl, elementSelector string) (*goquery.Document, error)
	GetHtmlOfPage(url string) (string, error)
	WaitElementAndGetHtmlOfPage(pageUrl, elementSelector string) (string, error)
	Close() error

	GetTab(tabNumber int) (tab TabInterface, err error)
	OpenTab() (tab TabInterface, tabNumber int, err error)
	CloseTab(tabNumber int) error

	Download(url string) ([]byte, error)
	DownloadFileBySelector(pageUrl, selector, attribute string) (bytes []byte, fileName, fileExtension, fileUrl string, err error)
}
```


```go
type TabInterface interface {
	GetHtmlOfPage(url string) (string, error)
	WaitElementAndGetHtmlOfPage(pageUrl, elementSelector string) (string, error)
	Close() error
}
```
