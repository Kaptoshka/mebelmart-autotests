package pages

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"autotests/pkg/elements"

	"github.com/playwright-community/playwright-go"
)

// BasePage is the base struct for all Page Objects.
type BasePage struct {
	Page    playwright.Page
	BaseURL string
	Timeout time.Duration
	Name    string
	Log     *slog.Logger
}

// New creates a new BasePage.
func New(
	page playwright.Page,
	baseURL string,
	timeout time.Duration,
	name string,
	log *slog.Logger,
) *BasePage {
	return &BasePage{
		Page:    page,
		BaseURL: baseURL,
		Timeout: timeout,
		Name:    name,
		Log:     log,
	}
}

// Navigate opens the page at the given path relative to BaseURL.
func (p *BasePage) Navigate(path string) error {
	url := p.BaseURL + path

	p.Log.Info("navigating to", "url", url)

	if _, err := p.Page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   new(float64(p.Timeout)),
	}); err != nil {
		return fmt.Errorf("[%s] navigation FAILED: %w", p.Name, err)
	}

	return nil
}

// WaitForURL waits until the current URL matches the expected value.
func (p *BasePage) WaitForURL(urlPattern string) error {
	p.Log.Info("waiting for URL", "pattern", urlPattern)

	if err := p.Page.WaitForURL(urlPattern, playwright.PageWaitForURLOptions{
		Timeout: new(float64(p.Timeout)),
	}); err != nil {
		return fmt.Errorf("[%s] URL did not match [%s]: %w", p.Name, urlPattern, err)
	}

	return nil
}

// GetTitle returns the current page title.
func (p *BasePage) GetTitle() (string, error) {
	title, err := p.Page.Title()
	if err != nil {
		return "", fmt.Errorf("[%s] could not get title: %w", p.Name, err)
	}

	return title, nil
}

// GetCurrentURL returns the current URL.
func (p *BasePage) GetCurrentURL() string {
	return p.Page.URL()
}

// CSS creates an element using a CSS selector.
func (p *BasePage) CSS(selector, description string) *elements.Element {
	return elements.NewCSS(p.Page, selector, description, p.Timeout, p.Log)
}

// XPath creates an element using an XPath expression.
func (p *BasePage) XPath(selector, description string) *elements.Element {
	return elements.NewXPath(p.Page, selector, description, p.Timeout, p.Log)
}

// WaitForNetworkIdle waits for network to be idle.
func (p *BasePage) WaitForNetworkIdle() error {
	p.Log.Debug("waiting for network idle")

	if err := p.Page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: new(float64(p.Timeout)),
	}); err != nil {
		return fmt.Errorf("[%s] network did not become idle: %w", p.Name, err)
	}

	return nil
}

// ExecuteScript runs JavaScript on the page.
func (p *BasePage) ExecuteScript(script string, args ...any) (any, error) {
	result, err := p.Page.Evaluate(script, args...)
	if err != nil {
		return nil, fmt.Errorf("[%s] script execution failed: %w", p.Name, err)
	}

	return result, nil
}

// ParseInt extracts and converts a string to an integer.
func (p *BasePage) ParseInt(s string) (int, error) {
	re := regexp.MustCompile(`[^0-9]`)
	cleanStr := re.ReplaceAllString(s, "")
	if cleanStr == "" {
		return 0, fmt.Errorf("string '%s' contains no digits", s)
	}
	return strconv.Atoi(cleanStr)
}
