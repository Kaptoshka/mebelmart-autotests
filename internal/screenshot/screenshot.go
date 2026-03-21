package screenshot

import (
	"fmt"
	"log/slog"

	"github.com/playwright-community/playwright-go"
)

type Service struct {
	log *slog.Logger
}

// New creates a new screenshot service.
func New(log *slog.Logger) *Service {
	return &Service{log: log}
}

// CaptureAsBites captures a screenshot of the page as bytes.
func (s *Service) CaptureAsBites(page playwright.Page) ([]byte, error) {
	bytes, err := page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: new(true),
	})
	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}
	return bytes, nil
}
