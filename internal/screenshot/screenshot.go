package screenshot

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Service struct {
	outputDir string
	log       *slog.Logger
}

func New(outputDir string, log *slog.Logger) *Service {
	if err := os.MkdirAll(outputDir, 0o700); err != nil {
		log.Warn("could not create screenshot dir", "err", err)
	}
	return &Service{outputDir: outputDir, log: log}
}

func (s *Service) Capture(page playwright.Page, name string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s", sanitizeName(name), timestamp)
	fullPath := filepath.Join(s.outputDir, fileName)

	s.log.Info("taking screenshot", "path", fullPath)

	_, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     new(fullPath),
		FullPage: new(true),
	})
	if err != nil {
		return "", fmt.Errorf("screenshot failed: %w", err)
	}

	s.log.Info("screenshot saved", "path", fullPath)
	return fullPath, nil
}

func (s *Service) CaptureOnFailure(page playwright.Page, testName string) (string, error) {
	return s.Capture(page, fmt.Sprintf("FAIL_%s", testName))
}

func (s *Service) CaptureAsBites(page playwright.Page) ([]byte, error) {
	bytes, err := page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: new(true),
	})
	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}
	return bytes, nil
}

func sanitizeName(name string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|', ' ':
			return '_'
		default:
			return r
		}
	}, name)
}
