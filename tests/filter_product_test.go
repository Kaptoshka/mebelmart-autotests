package cases_test

import (
	"errors"
	"fmt"
	"testing"

	"autotests/pkg/suite"
	testpages "autotests/tests/pages"

	"github.com/stretchr/testify/require"
)

const (
	catalogURL  = "/myagkaya_mebel_v_saratove/divanyi_v_saratove"
	filterName  = "Цена"
	productName = "Диван Чебурашка"

	minBorder = 9315
	maxBorder = 10815
)

func TestFilterProductAndCheckIsAvailable(t *testing.T) {
	t.Parallel()

	s := suite.New(t, "FilterProductAndCheckIsAvailable")

	if err := s.Setup(t.Name()); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	var testErr error
	defer s.Teardown(t.Name(), &testErr)

	s.Reporter.AddLabel("feature", "Filters")
	s.Reporter.AddLabel("severity", "normal")

	catalog := testpages.NewCatalogPage(s.Browser.Page, s.Config.BaseURL, s.Config.Timeout, s.Log)

	testErr = s.Step("Navigate to couches page", func() error {
		return catalog.Navigate(catalogURL)
	})
	require.NoError(t, testErr)

	testErr = s.Step("Wait for network idle", func() error {
		return catalog.WaitForNetworkIdle()
	})
	require.NoError(t, testErr)

	testErr = s.Step("Click filter container for price", func() error {
		return catalog.ClickFilterContainer(filterName)
	})
	require.NoError(t, testErr)

	testErr = s.Step(fmt.Sprintf("Set range %d - %d", minBorder, maxBorder), func() error {
		return catalog.SetRangeSliderByDrag(filterName, minBorder, maxBorder)
	})
	require.NoError(t, testErr)

	testErr = s.Step("Click apply button", func() error {
		return catalog.ClickApplyButton()
	})

	testErr = s.Step("Wait for loading results", func() error {
		return catalog.WaitForResults()
	})
	require.NoError(t, testErr)

	testErr = s.Step("Check that results is not empty", func() error {
		count, err := catalog.GetResultsCount()
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("no results found")
		}
		s.Log.Info("Results count", "count", count)
		return nil
	})
	require.NoError(t, testErr)

	testErr = s.Step(
		"Find couch with price in range",
		func() error {
			return catalog.FindProduct(productName)
		},
	)
	require.NoError(t, testErr)
}
