package testpages

import (
	"fmt"
	"log/slog"
	"time"

	"autotests/pkg/pages"

	"github.com/playwright-community/playwright-go"
)

type WishlistPage struct {
	*pages.BasePage
}

func NewWishlistPage(
	page playwright.Page,
	baseURL string,
	timeout time.Duration,
	testLog *slog.Logger,
) *WishlistPage {
	return &WishlistPage{
		BasePage: pages.New(
			page,
			baseURL,
			timeout,
			"WishlistPage",
			testLog.With("page", "WishlistPage"),
		),
	}
}

func (p *WishlistPage) GetItemsCount() (int, error) {
	return p.Page.Locator(".page-favorite .product-card").Count()
}

func (p *WishlistPage) Clear() error {
	p.Log.Debug("Clearing wishlist")

	for {
		count, err := p.GetItemsCount()
		if err != nil {
			return fmt.Errorf("cannot count wishlist items: %w", err)
		}
		if count == 0 {
			p.Log.Info("Wishlist is empty")
			return nil
		}

		p.Log.Debug("Items remaining", "count", count)

		removeBtn := p.Page.Locator(
			".product-card__favorite-delete",
		).First()

		if err = removeBtn.Click(); err != nil {
			return fmt.Errorf("cannot click remove button: %w", err)
		}

		if err = p.WaitForNetworkIdle(); err != nil {
			return fmt.Errorf("page did not update after remove: %w", err)
		}
	}
}

func (p *WishlistPage) IsEmpty() (bool, error) {
	count, err := p.GetItemsCount()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (p *WishlistPage) FindProductURL(name string) (string, error) {
	p.Log.Debug("Checking product in wishlist", "name", name)

	productURL, err := p.Page.Locator(
		".page-favorite .product-card__name a",
	).GetAttribute("href")
	if err != nil {
		p.Log.Error(
			"failed to find product url in wishlist",
			"name",
			name,
			"error",
			err,
		)
		return "", fmt.Errorf("failed to find product url in wishlist: %w", err)
	}

	return productURL, nil
}
