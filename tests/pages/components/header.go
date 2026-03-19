package components

import (
	"fmt"
	"log/slog"
	"time"

	"autotests/pkg/elements"

	"github.com/playwright-community/playwright-go"
)

type Header struct {
	page    playwright.Page
	Log     *slog.Logger
	timeout time.Duration

	searchInput *elements.Element
	searchBtn   *elements.Element
}

func NewHeader(
	page playwright.Page,
	timeout time.Duration,
	log *slog.Logger,
) *Header {
	return &Header{
		page:    page,
		Log:     log,
		timeout: timeout,
		searchInput: elements.NewCSS(
			page,
			"header .search .input-group input:visible",
			"Search Input",
			timeout,
			log,
		),
		searchBtn: elements.NewCSS(
			page,
			"header .search .input-group button.submit:visible",
			"Search Button",
			timeout,
			log,
		),
	}
}

func (h *Header) Search(query string) error {
	h.Log.Info("Searching for product via header", "query", query)

	if err := h.searchInput.WaitForVisible(); err != nil {
		return fmt.Errorf("search input not found: %w", err)
	}

	if err := h.searchInput.Fill(query); err != nil {
		return fmt.Errorf("failed to type search query: %w", err)
	}

	if err := h.searchInput.Press("Enter"); err != nil {
		return fmt.Errorf("failed to press Enter: %w", err)
	}

	return nil
}
