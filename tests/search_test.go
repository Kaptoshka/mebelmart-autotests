package cases_test

import (
	"fmt"
	"strings"
	"testing"

	"autotests/pkg/pages"
	"autotests/pkg/suite"
	testdata "autotests/tests"
	testpages "autotests/tests/pages"
	components "autotests/tests/pages/components"

	"github.com/stretchr/testify/require"
)

func TestSearchProduct(t *testing.T) {
	t.Parallel()

	s := suite.New(t, "SearchProduct")

	if err := s.Setup(t.Name()); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	s.SetMeta(suite.TestMeta{
		Description: fmt.Sprintf(
			"Search product by name [%s]",
			testdata.ProductName,
		),
		Severity: suite.SeverityNormal,
		Feature:  "search",
	})

	var testErr error
	defer s.Teardown(t.Name(), &testErr)

	basePage := pages.New(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		"BasePage",
		s.Log.With("page", "BasePage"),
	)

	searchResultsPage := testpages.NewSearchResultsPage(
		s.Browser.Page,
		s.Config.BaseURL,
		s.Config.Timeout,
		s.Log.With("page", "SearchResultsPage"),
	)

	header := components.NewHeader(
		s.Browser.Page,
		s.Config.Timeout,
		s.Log.With("component", "Header"),
	)

	testErr = s.Step(
		"Navigate to Home Page",
		func() error {
			return basePage.Navigate("/")
		},
	)
	require.NoError(t, testErr)

	testErr = s.Step(
		fmt.Sprintf(
			"Search product by query [%s] and press 'Enter'",
			testdata.SearchQuery,
		),
		func() error {
			return header.Search(testdata.SearchQuery)
		},
	)
	require.NoError(t, testErr)

	var cardName string
	testErr = s.Step(
		fmt.Sprintf(
			"Check search result for [%s]",
			testdata.SearchQuery,
		),
		func() error {
			var err error
			cardName, err = searchResultsPage.CheckSearchResult(testdata.SearchQuery)
			if err != nil {
				return fmt.Errorf("failed to check search result: %w", err)
			}

			if !strings.Contains(
				strings.TrimSpace(cardName),
				strings.TrimSpace(testdata.SearchQuery),
			) {
				return fmt.Errorf(
					"search result does not contain search query [%s]",
					testdata.SearchQuery,
				)
			}

			return nil
		},
	)
	require.NoError(t, testErr)
}
