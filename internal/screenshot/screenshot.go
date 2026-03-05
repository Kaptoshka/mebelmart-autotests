package screenshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"autotests/internal/logger"

	"github.com/playwright-community/playwright-go"
)

type Service struct {
	outputDir string
}

func New(outputDir string) *Service {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		logger.Warn("could not create screenshot dir", "err", err)
	}
	return &Service{outputDir: outputDir}
}

func (s *Service) Capture(page playwright.Page, name string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s", sanitizeName(name), timestamp)
	fullPath := filepath.Join(s.outputDir, fileName)

	logger.Info("taking screenshot", "path", fullPath)

	_, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(fullPath),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("screenshot failed: %w", err)
	}

	logger.Info("screenshot saved", "path", fullPath)
	return fullPath, nil
}

func (s *Service) CaptureOnFailure(page playwright.Page, testName string) (string, error) {
	return s.Capture(page, fmt.Sprintf("FAIL_%s", testName))
}

func (s *Service) CaptureAsBites(page playwright.Page) ([]byte, error) {
	bytes, err := page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}
	return bytes, nil
}

func sanitizeName(name string) string {
	result := make([]byte, len(name))
	for i, c := range name {
		switch c {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|', ' ':
			result[i] = '_'
		default:
			result[i] = byte(c)
		}
	}
	return string(result)
}
