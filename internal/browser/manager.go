package browser

import (
	"fmt"

	"autotests/internal/config"
	"autotests/internal/logger"

	"github.com/playwright-community/playwright-go"
)

type Manager struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	Page    playwright.Page
	cfg     *config.Config
}

func New(cfg *config.Config) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

func (m *Manager) Launch() error {
	logger.Info("Launching Playwright...")

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start Playwright: %w", err)
	}
	m.pw = pw

	browserType, err := m.getBrowserType()
	if err != nil {
		return err
	}

	logger.Info("Starting browser", "browser", m.cfg.Browser, "headless", m.cfg.Headless)

	slowMo := m.cfg.SlowMo.Milliseconds()
	browser, err := browserType.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(m.cfg.Headless),
		SlowMo:   playwright.Float(float64(slowMo)),
	})
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}
	m.browser = browser

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport: &playwright.Size{
			Width:  m.cfg.ViewportWidth,
			Height: m.cfg.ViewportHeight,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create browser context: %w", err)
	}
	m.context = context

	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}
	m.Page = page

	page.SetDefaultNavigationTimeout(float64(m.cfg.Timeout.Milliseconds()))
	page.SetDefaultTimeout(float64(m.cfg.Timeout.Milliseconds()))

	logger.Info("Browser launched successfully")
	return nil
}

func (m *Manager) Close() {
	if m.context != nil {
		logger.Debug("Closing browser context")
	}
	if m.browser != nil {
		logger.Debug("Closing browser")
	}
	if m.pw != nil {
		logger.Debug("Stopping playwright")
		m.pw.Stop()
	}
}

func (m *Manager) NavigateTo(url string) error {
	logger.Info("Navigating to", "url", url)
	if _, err := m.Page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		return fmt.Errorf("navigation to %s failed: %w", url, err)
	}
	return nil
}

func (m *Manager) getBrowserType() (playwright.BrowserType, error) {
	switch m.cfg.Browser {
	case config.BrowserFirefox:
		return m.pw.Firefox, nil
	case config.BrowserChromium:
		return m.pw.Chromium, nil
	default:
		return nil, fmt.Errorf("unsupported browser: %s", m.cfg.Browser)
	}
}
