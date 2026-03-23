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

func (p *ProductPage) GetDetails() (*Product, error) {
	p.Log.Debug("Getting details of product")

	productName, err := p.CSS(
		".page-product h1.text-center",
		"Product name",
	).GetText()
	if err != nil {
		return nil, fmt.Errorf("cannot get product name: %w", err)
	}

	priceStr, err := p.CSS(
		".page-product__now-price span.productPrice",
		"Product price",
	).GetText()
	if err != nil {
		return nil, fmt.Errorf("cannot get product price: %w", err)
	}

	price, err := p.ParseInt(priceStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse product price to int: %w", err)
	}

	width, err := p.GetParam("Ширина")
	if err != nil {
		return nil, fmt.Errorf("cannot get product width: %w", err)
	}

	depth, err := p.GetParam("Глубина")
	if err != nil {
		return nil, fmt.Errorf("cannot get product depth: %w", err)
	}

	return &Product{
		Name:  productName,
		Price: price,
		Width: width,
		Depth: depth,
	}, nil
}

func (p *ProductPage) GetParam(filter string) (int, error) {
	xpath := fmt.Sprintf(
		`//div[contains(@class,'product-tab__block')]//td[contains(.,'%s')]`+
			`/following-sibling::td[1]`,
		filter,
	)

	valueLocator := p.XPath(
		xpath,
		fmt.Sprintf("Filter value for [%s]", filter),
	)

	productFilterStr, err := valueLocator.GetText()
	if err != nil {
		return 0, fmt.Errorf(
			"cannot get product filter: %w",
			err,
		)
	}

	productFilter, err := p.ParseInt(productFilterStr)
	if err != nil {
		return 0, fmt.Errorf(
			"cannot parse product filter to int: %w",
			err,
		)
	}

	return productFilter, nil
}

func (p *ProductPage) AddToCart() error {
	p.Log.Debug("Adding product to cart")

	p.Page.On("dialog", func(dialog playwright.Dialog) {
		p.Log.Debug("Auto-accepting dialog", "message", dialog.Message())
		_ = dialog.Accept()
	})

	addToCartBtn := p.CSS(
		".page-product__buy a.btn.btnToCart:has-text('В корзину')",
		"Add product to cart button",
	)

	if err := addToCartBtn.Click(); err != nil {
		return fmt.Errorf("failed to click add to cart button: %w", err)
	}

	if err := p.Page.WaitForURL("**/cart**", playwright.PageWaitForURLOptions{
		Timeout: new(float64(p.Timeout)),
	}); err != nil {
		return fmt.Errorf("failed to navigate to cart after adding product: %w", err)
	}

	p.Log.Debug("Navigated to cart", "url", p.Page.URL())
	return nil
}

func (p *ProductPage) GoToCartViaDialog(dialog playwright.Dialog) error {
	p.Log.Debug("Accepting dialog to navigate to cart")

	if err := dialog.Accept(); err != nil {
		return fmt.Errorf("failed to accept dialog: %w", err)
	}

	if err := p.Page.WaitForURL("**/cart**"); err != nil {
		return fmt.Errorf("failed to navigate to cart after accepting dialog: %w", err)
	}

	return nil
}
