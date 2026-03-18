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
	t.Parallel()

	s := suite.New(t, "CheckProductDetails")

	if err := s.Setup(t.Name()); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	var testErr error
	defer s.Teardown(t.Name(), &testErr)

	s.SetMeta(suite.TestMeta{
		Description: "Check product card details and compare value with details on product page",
		Severity:    "normal",
		Feature:     "checkout",
	})

	catalog := testpages.NewCatalogPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log,
	)
	productDetails := testpages.NewProductPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log,
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

	var card *testpages.ProductCard
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

	var details *testpages.ProductDetails
	testErr = s.Step(
		"Navigate to product page",
		func() error {
			var err error
			if err = productDetails.Navigate(card.URL); err != nil {
				return err
			}
			details, err = productDetails.GetDetails(card.Name, ProductFilter)
			if err != nil {
				return err
			}
			s.Log.Info("Product details",
				"name", details.Name,
				"filter", details.Filter,
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
			if card.Width != details.Filter {
				return fmt.Errorf(
					"'%s' mismatch: expected '%d', got '%d'",
					ProductFilter,
					card.Width,
					details.Filter,
				)
			}
			return nil
		},
	)
	require.NoError(t, testErr)
}
