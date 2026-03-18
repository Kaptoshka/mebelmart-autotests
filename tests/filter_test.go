package cases_test

import (
	"errors"
	"fmt"
	"testing"

	"autotests/pkg/suite"
	testdata "autotests/tests"
	testpages "autotests/tests/pages"

	"github.com/stretchr/testify/require"
)

func TestFilterProductAndCheckIsAvailable(t *testing.T) {
	t.Parallel()

	s := suite.New(t, "FilterProductAndCheckIsAvailable")

	if err := s.Setup(t.Name()); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	s.SetMeta(suite.TestMeta{
		Description: "Filter product by price range and check that product is available",
		Severity:    "normal",
		Feature:     "filters",
	})

	var testErr error
	defer s.Teardown(t.Name(), &testErr)

	catalog := testpages.NewCatalogPage(
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

	testErr = s.Step(
		"Click 'Price' filter container",
		func() error {
			return catalog.ClickFilterContainer(testdata.FilterName)
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		fmt.Sprintf(
			"Set price range (%d - %d)",
			testdata.MinBorder,
			testdata.MaxBorder,
		),
		func() error {
			return catalog.SetRangeSliderByDrag(
				testdata.FilterName,
				testdata.MinBorder,
				testdata.MaxBorder,
			)
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step("Click 'Apply filters' button", func() error {
		return catalog.ClickApplyButton()
	})

	testErr = s.Step("Wait for loading results", func() error {
		return catalog.WaitForResults()
	})
	require.NoError(t, testErr)

	testErr = s.Step(
		"Check that filter applied and results is not empty",
		func() error {
			count, err := catalog.GetResultsCount()
			if err != nil {
				return err
			}
			if count == 0 {
				return errors.New("no results found")
			}
			s.Log.Info("Results count", "count", count)
			return nil
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		fmt.Sprintf(
			"Find couch with price in range (%d - %d)",
			testdata.MinBorder,
			testdata.MaxBorder,
		),
		func() error {
			return catalog.FindProduct(testdata.ProductName)
		},
	)
	require.NoError(t, testErr)
}
