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
	return p.CSS(
		".page-favorite .product-card",
		"Get wishlist items count",
	).Count()
}

func (p *WishlistPage) Clear() error {
	p.Log.Debug("Clearing wishlist")

	if err := p.Navigate("/favorite"); err != nil {
		return fmt.Errorf(
			"cannot navigate to wishlist for clearing: %w",
			err,
		)
	}

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

		removeBtn := p.CSS(
			".product-card__favorite-delete",
			"Get remove from wishlist button",
		).First("First remove from wishlist button")

		if err = removeBtn.Click(); err != nil {
			return fmt.Errorf("cannot click remove from wishlist button: %w", err)
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

	productURL, err := p.CSS(
		".page-favorite .product-card__name a",
		"Get product url in wishlist",
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
