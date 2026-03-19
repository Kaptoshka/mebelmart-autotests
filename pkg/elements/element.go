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
	e.log.Debug("Waiting for element to be visible", "element", e.description)

	if err := e.locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: new(float64(e.timeout)),
	}); err != nil {
		return fmt.Errorf("element [%s] not visible after %v: %w", e.description, e.timeout, err)
	}

	return nil
}

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

func (e *Element) Click() error {
	e.log.Debug("Clicking element", "element", e.description)

	if err := e.locator.Click(); err != nil {
		return fmt.Errorf("failed to click [%s]: %w", e.description, err)
	}

	e.log.Debug("Clicked element", "element", e.description)

	return nil
}

func (e *Element) Fill(text string) error {
	e.log.Debug("Filling element", "element", e.description, "text", text)

	if err := e.locator.Fill(text); err != nil {
		return fmt.Errorf("failed to fill [%s]: %w", e.description, err)
	}

	return nil
}

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

func (e *Element) GetText() (string, error) {
	e.log.Debug("Getting text from element", "element", e.description)

	text, err := e.locator.TextContent()
	if err != nil {
		return "", fmt.Errorf("failed to get text from [%s]: %w", e.description, err)
	}

	e.log.Debug("Got text from element", "element", e.description, "text", text)

	return text, nil
}

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

func (e *Element) Hover() error {
	e.log.Debug("hovering over element", "element", e.description)

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

func (e *Element) First(description string) *Element {
	return &Element{
		page:        e.page,
		locator:     e.locator.First(),
		description: fmt.Sprintf("%s [First]", e.description),
		timeout:     e.timeout,
		log:         e.log,
	}
}

func (e *Element) Nth(index int, description string) *Element {
	return &Element{
		page:        e.page,
		locator:     e.locator.Nth(index),
		description: fmt.Sprintf("%s[Index: %d]", e.description, index),
		timeout:     e.timeout,
		log:         e.log,
	}
}

func (e *Element) Count() (int, error) {
	e.log.Debug("Counting element", "element", e.description)

	count, err := e.locator.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count elements [%s]: %w", e.description, err)
	}

	return count, nil
}

func (e *Element) Blur() error {
	e.log.Debug("Removing focus from element", "element", e.description)

	if err := e.locator.Blur(); err != nil {
		return fmt.Errorf("failed to blur element [%s]: %w", e.description, err)
	}

	return nil
}

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
