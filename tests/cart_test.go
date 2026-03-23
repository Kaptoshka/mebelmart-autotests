package cases_test

import (
	"fmt"
	"testing"

	"autotests/pkg/suite"
	testdata "autotests/tests"
	testpages "autotests/tests/pages"

	"github.com/stretchr/testify/require"
)

func TestAddProductToCart(t *testing.T) {
	// t.Parallel()

	s := suite.New(t, "AddProductToCart")

	if err := s.Setup(t.Name()); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	s.SetMeta(suite.TestMeta{
		Description: "Add product to cart and check that it is in cart and price is correct",
		Severity:    suite.SeverityNormal,
		Feature:     "cart",
	})

	catalog := testpages.NewCatalogPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log.With("page", "CatalogPage"),
	)

	productDetails := testpages.NewProductPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log.With("page", "ProductPage"),
	)

	cart := testpages.NewCartPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log.With("page", "CartPage"),
	)

	var testErr error
	defer s.Teardown(t.Name(), &testErr)

	defer func() {
		if err := cart.Clear(); err != nil {
			s.Log.Warn("Failed to clear cart", "err", err)
		}
	}()

	testErr = s.Step(
		fmt.Sprintf(
			"Navigate to catalog page [%s]",
			testdata.CatalogURL,
		),
		func() error {
			return catalog.Navigate(testdata.CatalogURL)
		},
	)
	require.NoError(t, testErr)

	var catalogProduct *testpages.Product
	testErr = s.Step(
		fmt.Sprintf(
			"Get price of product [%s]",
			testdata.ProductName,
		),
		func() error {
			var err error
			catalogProduct, err = catalog.GetProductCard(testdata.ProductName)
			if err != nil {
				return err
			}
			s.Log.Debug("Product card data:",
				"name", catalogProduct.Name,
				"price", catalogProduct.Price,
				"url", catalogProduct.URL,
			)
			return nil
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		"Navigate to product details page",
		func() error {
			return productDetails.Navigate(catalogProduct.URL)
		},
	)
	require.NoError(t, testErr)

	var detailsProduct *testpages.Product
	testErr = s.Step(
		"Get product details",
		func() error {
			var err error
			detailsProduct, err = productDetails.GetDetails()
			return err
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		"Check that product price is correct",
		func() error {
			if detailsProduct.Price != catalogProduct.Price {
				return fmt.Errorf(
					"product price is not correct. Expected: %d, Actual: %d",
					catalogProduct.Price,
					detailsProduct.Price,
				)
			}
			return nil
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		fmt.Sprintf(
			"Add product [%s] to cart and navigate to cart page",
			detailsProduct.Name,
		),
		func() error {
			return productDetails.AddToCart()
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		"Check that product is in cart and price is correct",
		func() error {
			return cart.CheckCartForAndValidatePrice(detailsProduct)
		},
	)
	require.NoError(t, testErr)
}
