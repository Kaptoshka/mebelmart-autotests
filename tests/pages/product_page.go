package testpages

import (
	"fmt"
	"log/slog"
	"time"

	"autotests/pkg/pages"

	"github.com/playwright-community/playwright-go"
)

type ProductPage struct {
	*pages.BasePage
}

type ProductDetails struct {
	Name   string
	Filter int
}

func NewProductPage(
	page playwright.Page,
	baseURL string,
	timeout time.Duration,
	testLog *slog.Logger,
) *ProductPage {
	return &ProductPage{
		BasePage: pages.New(
			page,
			baseURL,
			timeout,
			"ProductPage",
			testLog,
		),
	}
}

func (p *ProductPage) OpenDetails() error {
	p.Log.Debug("Click button that opens details of product")
	return p.CSS(
		".product-tab #singleProdParamTab li:has-text('Характеристики')",
		"Open details of product",
	).Click()
}

func (p *ProductPage) GetDetails(
	name string,
	filter string,
) (*ProductDetails, error) {
	p.Log.Debug(
		"Get details of product",
		"name", name,
		"filter", filter,
	)

	productFilterLocator := p.CSS(
		".product-tab__block .table td",
		"Get product filter",
	).FilterByText(
		filter,
		"Get product filter value",
	)

	productFilterStr, err := productFilterLocator.FindCSS(
		"xpath=following-sibling::td[1]",
		"Get product filter value",
	).GetText()
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get product filter: %w",
			err,
		)
	}

	productFilter, err := p.ParseInt(productFilterStr)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot parse product filter to int: %w",
			err,
		)
	}

	return &ProductDetails{
		Name:   name,
		Filter: productFilter,
	}, nil
}
