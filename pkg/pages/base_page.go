package pages

import (
	"fmt"
	"log/slog"
	"time"

	"autotests/pkg/elements"

	"github.com/playwright-community/playwright-go"
)

type BasePage struct {
	Page    playwright.Page
	BaseURL string
	Timeout time.Duration
	Name    string
	log     *slog.Logger
}

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
		log:     log,
	}
}

func (p *BasePage) Navigate(path string) error {
	url := p.BaseURL + path

	p.log.Info("navigating to", "url", url)

	if _, err := p.Page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(float64(p.Timeout)),
	}); err != nil {
		return fmt.Errorf("[%s] navigation FAILED: %w", p.Name, err)
	}

	return nil
}

func (p *BasePage) WaitForURL(urlPattern string) error {
	p.log.Info("waiting for URL", "pattern", urlPattern)

	if err := p.Page.WaitForURL(urlPattern, playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(float64(p.Timeout)),
	}); err != nil {
		return fmt.Errorf("[%s] URL did not match [%s]: %w", p.Name, urlPattern, err)
	}

	return nil
}

func (p *BasePage) GetTitle() (string, error) {
	title, err := p.Page.Title()
	if err != nil {
		return "", fmt.Errorf("[%s] could not get title: %w", p.Name, err)
	}

	return title, nil
}

func (p *BasePage) GetCurrentURL() string {
	return p.Page.URL()
}

func (p *BasePage) CSS(selector, description string) *elements.Element {
	return elements.NewCSS(p.Page, selector, description, p.Timeout, p.log)
}

func (p *BasePage) XPath(selector, description string) *elements.Element {
	return elements.NewXPath(p.Page, selector, description, p.Timeout, p.log)
}

func (p *BasePage) WaitForNetworkIdle() error {
	p.log.Debug("waiting for network idle")

	if err := p.Page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: playwright.Float(float64(p.Timeout)),
	}); err != nil {
		return fmt.Errorf("[%s] network did not become idle: %w", p.Name, err)
	}

	return nil
}

func (p *BasePage) ExecuteScript(script string, args ...any) (any, error) {
	result, err := p.Page.Evaluate(script, args...)
	if err != nil {
		return nil, fmt.Errorf("[%s] script execution failed: %w", p.Name, err)
	}

	return result, nil
}
