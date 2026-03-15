package elements

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/playwright-community/playwright-go"
)

type LocatorType string

const (
	CSS   LocatorType = "css"
	XPath LocatorType = "xpath"
)

type Element struct {
	page        playwright.Page
	locator     playwright.Locator
	description string
	timeout     time.Duration
	log         *slog.Logger
}

func NewCSS(
	page playwright.Page,
	selector string,
	description string,
	timeout time.Duration,
	log *slog.Logger,
) *Element {
	return newElement(page, selector, description, CSS, timeout, log)
}

func NewXPath(
	page playwright.Page,
	xpath string,
	description string,
	timeout time.Duration,
	log *slog.Logger,
) *Element {
	return newElement(page, "xpath="+xpath, description, XPath, timeout, log)
}

func newElement(
	page playwright.Page,
	selector string,
	description string,
	lt LocatorType,
	timeout time.Duration,
	log *slog.Logger,
) *Element {
	log.Debug("Creating element", "element", description, "type", lt, "selector", selector)

	return &Element{
		page:        page,
		locator:     page.Locator(selector),
		description: description,
		timeout:     timeout,
		log:         log,
	}
}

func (e *Element) WaitForVisible() error {
	e.log.Debug("waiting for element to be visible", "element", e.description)

	if err := e.locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: new(float64(e.timeout.Milliseconds())),
	}); err != nil {
		return fmt.Errorf("element [%s] not visible after %v: %w", e.description, e.timeout, err)
	}

	return nil
}

func (e *Element) WaitForHidden() error {
	e.log.Debug("waiting for element to be hidden", "element", e.description)

	if err := e.locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateHidden,
		Timeout: new(float64(e.timeout.Milliseconds())),
	}); err != nil {
		return fmt.Errorf("element [%s] not visible after %v: %w", e.description, e.timeout, err)
	}

	return nil
}

func (e *Element) Click() error {
	e.log.Debug("clicking element", "element", e.description)

	if err := e.WaitForVisible(); err != nil {
		return err
	}

	if err := e.locator.Click(); err != nil {
		return fmt.Errorf("failed to click [%s]: %w", e.description, err)
	}

	e.log.Debug("clicked element", "element", e.description)

	return nil
}

func (e *Element) Fill(text string) error {
	e.log.Debug("filling element", "element", e.description, "text", text)

	if err := e.WaitForVisible(); err != nil {
		return err
	}

	if err := e.locator.Fill(text); err != nil {
		return fmt.Errorf("failed to fill [%s]: %w", e.description, err)
	}

	return nil
}

func (e *Element) Clear() error {
	e.log.Debug("clearing element", "element", e.description)

	if err := e.WaitForVisible(); err != nil {
		return err
	}

	if err := e.locator.Clear(); err != nil {
		return fmt.Errorf("failed to clear [%s]: %w", e.description, err)
	}

	return nil
}

func (e *Element) GetText() (string, error) {
	if err := e.WaitForVisible(); err != nil {
		return "", err
	}

	text, err := e.locator.TextContent()
	if err != nil {
		return "", fmt.Errorf("failed to get text from [%s]: %w", e.description, err)
	}

	e.log.Debug("got text from element", "element", e.description, "text", text)

	return text, nil
}

func (e *Element) GetAttribute(attr string) (string, error) {
	if err := e.WaitForVisible(); err != nil {
		return "", err
	}

	value, err := e.locator.GetAttribute(attr)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get attribute [%s] from [%s]: %w",
			attr,
			e.description,
			err,
		)
	}

	return value, nil
}

func (e *Element) IsVisible() (bool, error) {
	visible, err := e.locator.IsVisible()
	if err != nil {
		return false, fmt.Errorf("failed to check visibility of [%s]: %w", e.description, err)
	}

	return visible, nil
}

func (e *Element) IsEnabled() (bool, error) {
	enabled, err := e.locator.IsEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to check enabled state of [%s]: %w", e.description, err)
	}

	return enabled, nil
}

func (e *Element) SelectOption(value string) error {
	e.log.Debug("selecting option", "value", value, "element", e.description)
	if err := e.WaitForVisible(); err != nil {
		return err
	}

	_, err := e.locator.SelectOption(playwright.SelectOptionValues{
		Values: &[]string{value},
	})
	if err != nil {
		return fmt.Errorf(
			"failed to select option [%s] in [%s]: %w",
			value,
			e.description,
			err,
		)
	}

	return nil
}

func (e *Element) Hover() error {
	e.log.Debug("hovering over element", "element", e.description)

	if err := e.WaitForVisible(); err != nil {
		return err
	}

	if err := e.locator.Hover(); err != nil {
		return fmt.Errorf("failed to hover over [%s]: %w", e.description, err)
	}

	return nil
}

func (e *Element) ScrollIntoView() error {
	e.log.Debug("scrolling element into view", "element", e.description)

	if err := e.locator.ScrollIntoViewIfNeeded(); err != nil {
		return fmt.Errorf(
			"failed to scroll [%s] into view: %w",
			e.description,
			err,
		)
	}

	return nil
}
