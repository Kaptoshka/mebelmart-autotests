package cases_test

import (
	"fmt"
	"testing"

	"autotests/pkg/suite"
	testdata "autotests/tests"
	testpages "autotests/tests/pages"

	"github.com/stretchr/testify/require"
)

const (
	ProductFilter = "Ширина" // Ширина | Глубина
)

func TestCheckProductDetails(t *testing.T) {
	// t.Parallel()

	s := suite.New(t, "CheckProductDetails")

	if err := s.Setup(t.Name()); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	var testErr error
	defer s.Teardown(t.Name(), &testErr)

	s.SetMeta(suite.TestMeta{
		Description: "Check product card details and compare value with details on product page",
		Severity:    suite.SeverityNormal,
		Feature:     "checkout",
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

	var card *testpages.Product
	testErr = s.Step(
		fmt.Sprintf(
			"Get product card data '%s'",
			testdata.ProductName,
		),
		func() error {
			var err error
			card, err = catalog.GetProductCard(testdata.ProductName)
			if err != nil {
				return err
			}
			s.Log.Info("Card data",
				"name", card.Name,
				"width", card.Width,
				"url", card.URL,
			)
			return nil
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		"Navigate to product page",
		func() error {
			return productDetails.Navigate(card.URL)
		},
	)
	require.NoError(t, testErr)

	var details *testpages.Product
	testErr = s.Step(
		"Get product details",
		func() error {
			var err error
			details, err = productDetails.GetDetails()
			if err != nil {
				return err
			}
			s.Log.Info("Product details",
				"name", details.Name,
				"price", details.Price,
				"width", details.Width,
				"depth", details.Depth,
			)
			return nil
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		fmt.Sprintf(
			"Compare product by [%s]",
			ProductFilter,
		),
		func() error {
			if card.Width != details.Width {
				return fmt.Errorf(
					"'%s' mismatch: expected '%d', got '%d'",
					ProductFilter,
					card.Width,
					details.Width,
				)
			}
			return nil
		},
	)
	require.NoError(t, testErr)
}
