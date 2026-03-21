package elements

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/playwright-community/playwright-go"
)

// LocatorType defines how to find an element.
type LocatorType string

const (
	CSS   LocatorType = "css"
	XPath LocatorType = "xpath"
)

// Element wraps a Playwright locator with explicit waits and logging.
type Element struct {
	page        playwright.Page
	locator     playwright.Locator
	description string
	timeout     time.Duration
	log         *slog.Logger
}

// NewCSS creates an element using CSS selector.
func NewCSS(
	page playwright.Page,
	selector string,
	description string,
	timeout time.Duration,
	log *slog.Logger,
) *Element {
	return newElement(page, selector, description, CSS, timeout, log)
}

// NewXPath creates an element using XPath locator.
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

// WaitForVisible explicitly waits until the element is visible.
func (e *Element) WaitForVisible() error {
	e.log.Debug("Waiting for element to be visible", "element", e.description)

	if err := e.locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: new(float64(e.timeout)),
	}); err != nil {
		return fmt.Errorf("element [%s] not visible after %v: %w", e.description, e.timeout, err)
	}

	return nil
}

// WaitForHidden waits until the element is hidden or detached.
func (e *Element) WaitForHidden() error {
	e.log.Debug("Waiting for element to be hidden", "element", e.description)

	if err := e.locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateHidden,
		Timeout: new(float64(e.timeout)),
	}); err != nil {
		return fmt.Errorf("element [%s] not visible after %v: %w", e.description, e.timeout, err)
	}

	return nil
}

// Click clicks the element.
func (e *Element) Click() error {
	e.log.Debug("Clicking element", "element", e.description)

	if err := e.locator.Click(); err != nil {
		return fmt.Errorf("failed to click [%s]: %w", e.description, err)
	}

	e.log.Debug("Clicked element", "element", e.description)

	return nil
}

// Fill fills input element with text.
func (e *Element) Fill(text string) error {
	e.log.Debug("Filling element", "element", e.description, "text", text)

	if err := e.locator.Fill(text); err != nil {
		return fmt.Errorf("failed to fill [%s]: %w", e.description, err)
	}

	return nil
}

// Clear clears an input field.
func (e *Element) Clear() error {
	e.log.Debug("Clearing element", "element", e.description)

	if err := e.WaitForVisible(); err != nil {
		return err
	}

	if err := e.locator.Clear(); err != nil {
		return fmt.Errorf("failed to clear [%s]: %w", e.description, err)
	}

	return nil
}

// GetText returns visible text of the element.
func (e *Element) GetText() (string, error) {
	e.log.Debug("Getting text from element", "element", e.description)

	text, err := e.locator.TextContent()
	if err != nil {
		return "", fmt.Errorf("failed to get text from [%s]: %w", e.description, err)
	}

	e.log.Debug("Got text from element", "element", e.description, "text", text)

	return text, nil
}

// GetAttribute returns the value of an attribute.
func (e *Element) GetAttribute(attr string) (string, error) {
	e.log.Debug(
		"Getting attribute",
		"element", e.description,
		"attribute", attr,
	)

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

// IsVisible returns true if element is visible without waiting.
func (e *Element) IsVisible() (bool, error) {
	visible, err := e.locator.IsVisible()
	if err != nil {
		return false, fmt.Errorf("failed to check visibility of [%s]: %w", e.description, err)
	}

	return visible, nil
}

// IsEnabled returns true if element is enabled.
func (e *Element) IsEnabled() (bool, error) {
	enabled, err := e.locator.IsEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to check enabled state of [%s]: %w", e.description, err)
	}

	return enabled, nil
}

// SelectOption selects an option in a dropdown.
func (e *Element) SelectOption(value string) error {
	e.log.Debug("Selecting option", "value", value, "element", e.description)

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

// Hover hovers over the element.
func (e *Element) Hover() error {
	e.log.Debug("hovering over element", "element", e.description)

	if err := e.locator.Hover(); err != nil {
		return fmt.Errorf("failed to hover over [%s]: %w", e.description, err)
	}

	return nil
}

// ScrollIntoView scrolls element into the viewport.
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

// FilterByText filters elements by text.
func (e *Element) FilterByText(text string, description string) *Element {
	e.log.Debug(
		"Filtering element by text",
		"element", e.description,
		"text", text,
	)

	return &Element{
		page: e.page,
		locator: e.locator.Filter(playwright.LocatorFilterOptions{
			HasText: text,
		}),
		description: fmt.Sprintf("%s [Text: %s]", e.description, text),
		timeout:     e.timeout,
		log:         e.log,
	}
}

// FindCSS finds a sub-element using CSS selector.
func (e *Element) FindCSS(subSelector string, description string) *Element {
	e.log.Debug(
		"Finding sub-element by CSS",
		"parent", e.description,
		"child", description,
	)

	return &Element{
		page:        e.page,
		locator:     e.locator.Locator(subSelector),
		description: fmt.Sprintf("%s -> %s", e.description, description),
		timeout:     e.timeout,
		log:         e.log,
	}
}

// FindXPath finds a sub-element using XPath locator.
func (e *Element) FindXPath(xpath string, description string) *Element {
	e.log.Debug(
		"Finding sub-element by XPath",
		"parent", e.description,
		"child", description,
	)

	return &Element{
		page:        e.page,
		locator:     e.locator.Locator("xpath=" + xpath),
		description: fmt.Sprintf("%s -> %s", e.description, description),
		timeout:     e.timeout,
		log:         e.log,
	}
}

// First returns the first element of the locator.
func (e *Element) First(description string) *Element {
	return &Element{
		page:        e.page,
		locator:     e.locator.First(),
		description: fmt.Sprintf("%s [First]", e.description),
		timeout:     e.timeout,
		log:         e.log,
	}
}

// Nth returns the nth element of the locator.
func (e *Element) Nth(index int, description string) *Element {
	return &Element{
		page:        e.page,
		locator:     e.locator.Nth(index),
		description: fmt.Sprintf("%s[Index: %d]", e.description, index),
		timeout:     e.timeout,
		log:         e.log,
	}
}

// Count returns the number of elements matched by the locator.
func (e *Element) Count() (int, error) {
	e.log.Debug("Counting element", "element", e.description)

	count, err := e.locator.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count elements [%s]: %w", e.description, err)
	}

	return count, nil
}

// Blur removes focus from the element.
func (e *Element) Blur() error {
	e.log.Debug("Removing focus from element", "element", e.description)

	if err := e.locator.Blur(); err != nil {
		return fmt.Errorf("failed to blur element [%s]: %w", e.description, err)
	}

	return nil
}

// Press simulates a key press on the element.
func (e *Element) Press(key string) error {
	e.log.Debug(
		"Pressing key on element",
		"element", e.description,
		"key", key,
	)

	if err := e.locator.Press(key); err != nil {
		return fmt.Errorf("failed to press key [%s] on [%s]: %w", key, e.description, err)
	}

	return nil
}

// GetBoundingBox returns the bounding box of the element.
func (e *Element) GetBoundingBox() (*playwright.Rect, error) {
	e.log.Debug(
		"Getting bounding box",
		"element",
		e.description,
	)

	if err := e.WaitForVisible(); err != nil {
		return nil, fmt.Errorf(
			"cannot get bounding box for [%s]: %w",
			e.description,
			err,
		)
	}

	box, err := e.locator.BoundingBox()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get bounding box for [%s]: %w",
			e.description,
			err,
		)
	}

	if box == nil {
		return nil, fmt.Errorf(
			"bounding box for [%s] is nil (element is not visible or detached)",
			e.description,
		)
	}

	return box, nil
}
