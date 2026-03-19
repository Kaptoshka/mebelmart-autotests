package testpages

import (
	"log/slog"
	"time"

	"autotests/pkg/pages"

	"github.com/playwright-community/playwright-go"
)

type SearchResultsPage struct {
	*pages.BasePage
}

func NewSearchResultsPage(
	page playwright.Page,
	baseURL string,
	timeout time.Duration,
	testLog *slog.Logger,
) *SearchResultsPage {
	return &SearchResultsPage{
		BasePage: pages.New(
			page,
			baseURL,
			timeout,
			"SearchResultsPage",
			testLog,
		),
	}
}

func (p *SearchResultsPage) CheckSearchResult(query string) (string, error) {
	p.Log.Debug("Checking search result", "query", query)
	return p.CSS(
		".content .product-card",
		"Check first search result",
	).First("First search result").FindCSS(
		".product-card__name a",
		"Product name",
	).GetText()
}
