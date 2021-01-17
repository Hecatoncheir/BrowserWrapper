package browser

type TabInterface interface {
	GetHtmlOfPage(url string) (string, error)
	WaitElementAndGetHtmlOfPage(pageUrl, elementSelector string) (string, error)
	Close() error
}
