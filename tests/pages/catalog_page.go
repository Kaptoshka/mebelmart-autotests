package testpages

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"autotests/pkg/elements"
	"autotests/pkg/pages"

	"github.com/playwright-community/playwright-go"
)

type ProductCard struct {
	Name  string
	Width int
	Depth int
	URL   string
}

type CatalogPage struct {
	*pages.BasePage

	productCards *elements.Element
}

func NewCatalogPage(
	page playwright.Page,
	baseURL string,
	timeout time.Duration,
	testLog *slog.Logger,
) *CatalogPage {
	return &CatalogPage{
		BasePage: pages.New(
			page,
			baseURL,
			timeout,
			"CatalogPage",
			testLog,
		),
		productCards: elements.NewCSS(
			page,
			".content .container .product-card:not(.owl-carousel .product-card)",
			"Product Cards List",
			timeout,
			testLog,
		),
	}
}

func (p *CatalogPage) ClickFilterContainer(filterName string) error {
	p.Log.Debug("Clicking filter container", "filter", filterName)
	return p.CSS(
		fmt.Sprintf("div.filter__title a:has-text('%s')", filterName),
		fmt.Sprintf("Filter title link [%s]", filterName),
	).Click()
}

func (p *CatalogPage) SetRangeSliderByDrag(
	filterName string,
	from int,
	to int,
) error {
	p.Log.Debug(
		"Setting range filter",
		"filter", filterName,
		"from", from,
		"to", to,
	)

	filterContainer := p.Page.Locator(".filter__item").Filter(
		playwright.LocatorFilterOptions{
			Has: p.Page.Locator(".filter__title").Filter(
				playwright.LocatorFilterOptions{
					HasText: filterName,
				},
			),
		},
	)

	minHandle := filterContainer.Locator(
		".slider-handle.min-slider-handle",
	)
	maxHandle := filterContainer.Locator(
		".slider-handle.max-slider-handle",
	)
	track := filterContainer.Locator(".slider-track")

	absMinStr, err := minHandle.GetAttribute("aria-valuemin")
	if err != nil {
		return fmt.Errorf("cannot read aria-valuemin: %w", err)
	}

	absMaxStr, err := maxHandle.GetAttribute("aria-valuemax")
	if err != nil {
		return fmt.Errorf("cannot read aria-valuemax: %w", err)
	}

	absMin, err := p.ParseInt(absMinStr)
	if err != nil {
		return fmt.Errorf("cannot parse aria-valuemin: %w", err)
	}

	absMax, err := p.ParseInt(absMaxStr)
	if err != nil {
		return fmt.Errorf("cannot parse aria-valuemax: %w", err)
	}

	if from < absMin || to > absMax {
		return fmt.Errorf(
			"invalid range [%d, %d] for [%s]: allowed [%d, %d]",
			from, to, filterName, absMin, absMax,
		)
	}

	trackBox, err := track.BoundingBox()
	if err != nil {
		return fmt.Errorf("cannot get track bounding box: %w", err)
	}

	rangeSize := float64(absMax - absMin)
	fromX := trackBox.X + (float64(from-absMin)/rangeSize)*trackBox.Width
	toX := trackBox.X + (float64(to-absMin)/rangeSize)*trackBox.Width

	const centerDivisor = 2

	centerY := trackBox.Y + trackBox.Height/centerDivisor

	if err = p.dragHandleTo(minHandle, fromX, centerY); err != nil {
		return fmt.Errorf("failed to drag min handle: %w", err)
	}

	if err = p.dragHandleTo(maxHandle, toX, centerY); err != nil {
		return fmt.Errorf("failed to drag max handle: %w", err)
	}

	p.Log.Info("Range filter set", "filter", filterName, "from", from, "to", to)
	return nil
}

func (p *CatalogPage) dragHandleTo(
	handle playwright.Locator,
	targetX float64,
	targetY float64,
) error {
	const centerDivisor = 2

	box, err := handle.BoundingBox()
	if err != nil {
		return fmt.Errorf("cannot get handle bounding box: %w", err)
	}

	startX := box.X + box.Width/centerDivisor
	startY := box.Y + box.Height/centerDivisor

	mouse := p.Page.Mouse()

	if err = mouse.Move(startX, startY); err != nil {
		return err
	}

	if err = mouse.Down(); err != nil {
		return err
	}

	steps := 10

	if err = p.Page.Mouse().Move(targetX, targetY, playwright.MouseMoveOptions{
		Steps: new(steps),
	}); err != nil {
		return err
	}

	return mouse.Up()
}

func (p *CatalogPage) ClickApplyButton() error {
	p.Log.Debug("Clicking apply button")
	return p.CSS(
		".filter__link div.btn:has-text('Применить фильтр')",
		"Apply button",
	).Click()
}

func (p *CatalogPage) WaitForResults() error {
	p.Log.Debug("Waiting for results to update")

	if err := p.productCards.First("First product card").WaitForVisible(); err != nil {
		return fmt.Errorf("no product cards appeared after filter: %w", err)
	}

	return nil
}

func (p *CatalogPage) GetResultsCount() (int, error) {
	p.Log.Debug("Getting results count")
	return p.productCards.Count()
}

func (p *CatalogPage) ClickSortButton(sortName string) error {
	p.Log.Debug("Clicking sort button", "SortName", sortName)
	return p.CSS(
		fmt.Sprintf(
			".sorting-bar .sorting-bar__text b:has-text('%s')",
			sortName,
		),
		"Sort button",
	).Click()
}

func (p *CatalogPage) FindProduct(
	name string,
) error {
	p.Log.Debug("Checking product visibility", "name", name)

	product := p.productCards.FilterByText(name, "Card for "+name)

	count, err := product.Count()
	if err != nil || count == 0 {
		return fmt.Errorf("failed to find product '%s': %w", name, err)
	}

	p.Log.Debug("Product is visible", "name", name, "count", count)
	return nil
}

func (p *CatalogPage) GetProductCard(name string) (*ProductCard, error) {
	p.Log.Debug("Getting product card", "name", name)

	card := p.productCards.FilterByText(name, "Card for "+name).First(
		"First product card",
	)

	cardName, err := card.FindCSS(
		".product-card__name",
		"Product card name",
	).GetText()
	if err != nil {
		return nil, fmt.Errorf("cannot get card name: %w", err)
	}

	cardName = strings.TrimSpace(cardName)
	paramSelector := ".text-center small:has-text('%s')"

	cardWidthStr, err := p.GetParamCSS(
		*card,
		fmt.Sprintf(paramSelector, "Ширина"),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get card width: %w", err)
	}

	width, err := p.ParseInt(cardWidthStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse card width to int: %w", err)
	}

	cardDepthStr, err := p.GetParamCSS(
		*card,
		fmt.Sprintf(paramSelector, "Глубина"),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get card depth: %w", err)
	}

	depth, err := p.ParseInt(cardDepthStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse card depth to int: %w", err)
	}

	url, err := p.GetURL(card, ".product-card__name a")
	if err != nil {
		return nil, fmt.Errorf("cannot get card URL: %w", err)
	}

	return &ProductCard{
		Name:  cardName,
		Width: width,
		Depth: depth,
		URL:   url,
	}, nil
}

func (p *CatalogPage) GetParamCSS(
	elem elements.Element,
	selector string,
) (string, error) {
	return elem.FindCSS(selector, "Parameter value retrieving").GetText()
}

func (p *CatalogPage) GetURL(
	elem *elements.Element,
	selector string,
) (string, error) {
	return elem.FindCSS(
		selector,
		"Product card URL retrieving",
	).GetAttribute("href")
}

func (p *CatalogPage) GetProductCardURL(name string) (string, error) {
	return p.GetURL(
		p.productCards.FilterByText(
			name,
			"Product card with name "+name,
		).First("Product card URL"),
		".product-card__name a",
	)
}

func (p *CatalogPage) AddToWishlist(name string) error {
	p.Log.Debug("Click button that adds product to wishlist", "name", name)
	return p.productCards.FilterByText(
		name,
		"Product card with name "+name,
	).First(
		"First product card",
	).FindCSS(
		".product-card__favorites .favorite-icon",
		"Favorite icon",
	).Click()
}

func (p *CatalogPage) OpenWishlist() error {
	p.Log.Debug("Click button that opens wishlist")

	return p.CSS(
		"header .container .favorite-informer:visible",
		"Wishlist link",
	).Click()
}

func (p *CatalogPage) IsActiveIcon(name string) error {
	p.Log.Debug(
		"Checking if favorite icon for product [%s] is active",
		"product",
		name,
	)

	productCard := p.productCards.FilterByText(
		name,
		"Card for "+name,
	)

	activeIcon := productCard.FindCSS(
		".product-card__favorites .favorite-icon.active",
		"Active favorite icon",
	)

	err := activeIcon.WaitForVisible()
	if err != nil {
		return fmt.Errorf("favorite icon for product [%s] is not active: %w", name, err)
	}

	return nil
}
