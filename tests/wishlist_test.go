package cases_test

import (
	"fmt"
	"strings"
	"testing"

	"autotests/pkg/suite"
	testdata "autotests/tests"
	testpages "autotests/tests/pages"

	"github.com/stretchr/testify/require"
)

func TestAddProductToWishlist(t *testing.T) {
	t.Parallel()

	s := suite.New(t, "AddProductToWishlist")

	if err := s.Setup(t.Name()); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	catalogPage := testpages.NewCatalogPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log.With("page", "CatalogPage"),
	)
	wishlistPage := testpages.NewWishlistPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log.With("page", "WishlistPage"),
	)

	var testErr error
	defer func() {
		_ = wishlistPage.Clear()
		s.Teardown(t.Name(), &testErr)
	}()

	s.SetMeta(suite.TestMeta{
		Description: "Add product to wishlist, check that product is added to wishlist",
		Severity:    "normal",
		Feature:     "wishlist",
	})

	testErr = s.Step(
		fmt.Sprintf(
			"Navigate to catalog page [%s]",
			testdata.CatalogURL,
		),
		func() error {
			return catalogPage.Navigate(testdata.CatalogURL)
		},
	)

	var productCatalogURL string
	testErr = s.Step(
		fmt.Sprintf(
			"Click favorite icon for adding product [%s] to wishlist.",
			testdata.ProductName,
		),
		func() error {
			var err error
			productCatalogURL, err = catalogPage.GetProductCardURL(testdata.ProductName)
			if err != nil {
				return fmt.Errorf("failed to get product card URL: %w", err)
			}
			return catalogPage.AddToWishlist(testdata.ProductName)
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		"Check that favorite icon is active",
		func() error {
			return catalogPage.IsActiveIcon(testdata.ProductName)
		},
	)

	testErr = s.Step(
		"Navigate to wishlist page",
		func() error {
			return catalogPage.OpenWishlist()
		},
	)
	require.NoError(t, testErr)

	var productNameWishlist string
	var err error
	testErr = s.Step(
		"Find product that we added to wishlist",
		func() error {
			productNameWishlist, err = wishlistPage.FindProductURL(testdata.ProductName)
			if err != nil {
				return err
			}

			if !strings.Contains(
				strings.TrimSpace(productNameWishlist),
				strings.TrimSpace(productCatalogURL),
			) {
				return fmt.Errorf("product [%s] not found in wishlist", testdata.ProductName)
			}

			return nil
		},
	)
	require.NoError(t, testErr)
}
