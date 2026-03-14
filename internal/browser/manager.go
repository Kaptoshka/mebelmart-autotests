package browser

import (
	"fmt"
	"log/slog"
	"os"

	"autotests/internal/config"

	"github.com/playwright-community/playwright-go"
)

type Manager struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	Page    playwright.Page
	cfg     *config.Config
	log     *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Manager {
	return &Manager{
		cfg: cfg,
		log: log,
	}
}

func (m *Manager) Launch() error {
	m.log.Info("Launching Playwright...")

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start Playwright: %w", err)
	}
	m.pw = pw

	browserType, err := m.getBrowserType()
	if err != nil {
		return err
	}

	m.log.Info("Starting browser", "browser", m.cfg.Browser, "headless", m.cfg.Headless)

	browser, err := browserType.Launch(m.getBrowserLaunchOptions())
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

	m.log.Info("Browser launched successfully")
	return nil
}

func (m *Manager) Close() {
	if m.context != nil {
		m.log.Debug("Closing browser context")
	}
	if m.browser != nil {
		m.log.Debug("Closing browser")
	}
	if m.pw != nil {
		m.log.Debug("Stopping playwright")
		m.pw.Stop()
	}
}

func (m *Manager) NavigateTo(url string) error {
	m.log.Info("Navigating to", "url", url)
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

func (m *Manager) getBrowserLaunchOptions() playwright.BrowserTypeLaunchOptions {
	slowMo := m.cfg.SlowMo.Milliseconds()
	opts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(m.cfg.Headless),
		SlowMo:   playwright.Float(float64(slowMo)),
	}

	switch m.cfg.Browser {
	case config.BrowserChromium:
		if execPath := os.Getenv("PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH"); execPath != "" {
			m.log.Info("using custom Chromium path", "path", execPath)
			opts.ExecutablePath = playwright.String(execPath)
		}
	case config.BrowserFirefox:
		if execPath := os.Getenv("PLAYWRIGHT_FIREFOX_EXECUTABLE_PATH"); execPath != "" {
			m.log.Info("using custom Firefox path", "path", execPath)
			opts.ExecutablePath = playwright.String(execPath)
		}
	}

	return opts
}
