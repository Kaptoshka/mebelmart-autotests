package suite

import (
	"fmt"
	"log/slog"
	"testing"

	"autotests/internal/browser"
	"autotests/internal/config"
	"autotests/internal/logger"
	"autotests/internal/reporter"
	"autotests/internal/screenshot"
)

type TestSuite struct {
	T          *testing.T
	Config     *config.Config
	Browser    *browser.Manager
	Screenshot *screenshot.Service
	Reporter   *reporter.AllureReporter
	SuiteName  string
	Log        *slog.Logger
}

type TestMeta struct {
	Description string
	Severity    string
	Feature     string
}

func New(t *testing.T, suiteName string) *TestSuite {
	cfg := config.Load()

	return &TestSuite{
		T:         t,
		Config:    cfg,
		SuiteName: suiteName,
	}
}

func (s *TestSuite) Setup(testName string) error {
	s.Log = logger.ForTest(s.T)
	s.Screenshot = screenshot.New(s.Log)
	s.Reporter = reporter.New(s.Config.AllureReportDir, testName, s.SuiteName, s.Log)

	s.Browser = browser.New(s.Config, s.Log)
	if err := s.Browser.Launch(); err != nil {
		s.Reporter.SetBroken(err)
		_ = s.Reporter.Finalize()
		return fmt.Errorf("browser setup failed: %w", err)
	}

	s.Log.Info("test setup complete", "test", testName)
	return nil
}

func (s *TestSuite) Teardown(testName string, testErr *error) {
	if testErr != nil && *testErr != nil {
		s.Log.Warn("test FAILED -- capturing screenshot", "test", testName)
		if bytes, err := s.Screenshot.CaptureAsBites(s.Browser.Page); err == nil {
			_ = s.Reporter.AddScreenshot(bytes, fmt.Sprintf("Failure: %s", testName))
		} else {
			s.Log.Warn("failed to capture screenshot", "err", err)
		}
		s.Reporter.SetFailed(*testErr)
	}

	if err := s.Reporter.Finalize(); err != nil {
		s.Log.Warn("could not finalize Allure report", "err", err)
	}

	err := s.Browser.Close()
	if err != nil {
		s.Log.Warn("could not close browser", "err", err)
	}

	s.Log.Info("test teardown complete", "test", testName)
}

func (s *TestSuite) Step(name string, fn func() error) error {
	s.Reporter.StartStep(name)

	if err := fn(); err != nil {
		s.Reporter.StopStep(reporter.StatusFailed)
		return fmt.Errorf("step [%s] failed: %w", name, err)
	}

	s.Reporter.StopStep(reporter.StatusPassed)
	return nil
}

func (s *TestSuite) NavigateTo(url string) error {
	return s.Browser.NavigateTo(url)
}

type label struct {
	key   string
	value string
}

func (s *TestSuite) SetMeta(meta TestMeta) {
	if meta.Description != "" {
		s.Reporter.SetDescription(meta.Description)
	}

	labels := []label{
		{"severity", meta.Severity},
		{"feature", meta.Feature},
	}

	for _, l := range labels {
		if l.value != "" {
			s.Reporter.AddLabel(l.key, l.value)
		}
	}
}
