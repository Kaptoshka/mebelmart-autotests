package reporter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	uuidG "github.com/google/uuid"
)

type Status string

const (
	StatusPassed  Status = "passed"
	StatusFailed  Status = "failed"
	StatusBroken  Status = "broken"
	StatusSkipped Status = "skipped"
)

type Attachment struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

type Step struct {
	Name        string       `json:"name"`
	Status      Status       `json:"status"`
	Start       int64        `json:"start"`
	Stop        int64        `json:"stop"`
	Attachments []Attachment `json:"attachments"`
	Steps       []Step       `json:"steps"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type StatusDetails struct {
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

type TestResult struct {
	UUID          string         `json:"uuid"`
	HistoryID     string         `json:"historyId"`
	FullName      string         `json:"fullName"`
	Name          string         `json:"name"`
	Status        Status         `json:"status"`
	Start         int64          `json:"start"`
	Stop          int64          `json:"stop"`
	Description   string         `json:"description,omitempty"`
	Steps         []Step         `json:"steps"`
	Attachments   []Attachment   `json:"attachments"`
	Labels        []Label        `json:"labels"`
	StatusDetails *StatusDetails `json:"statusDetails.omitempty"`
}

type AllureReporter struct {
	outputDir string
	result    *TestResult
	stepStack []*Step
	startTime int64
	log       *slog.Logger
}

func New(outputDir, testName, suiteName string, log *slog.Logger) *AllureReporter {
	if err := os.MkdirAll(outputDir, 0o700); err != nil {
		log.Warn("could not create allure results dir", "err", err)
	}

	uuid, err := uuidG.NewV7()
	if err != nil {
		log.Warn("could not generate uuid", "err", err)
	}

	uuidStr := uuid.String()
	now := time.Now().UnixMilli()

	return &AllureReporter{
		outputDir: outputDir,
		startTime: now,
		log:       log,
		result: &TestResult{
			UUID:      uuidStr,
			HistoryID: testName,
			FullName:  fmt.Sprintf("%s.%s", suiteName, testName),
			Name:      testName,
			Status:    StatusPassed,
			Start:     now,
			Labels: []Label{
				{Name: "suite", Value: suiteName},
				{Name: "framework", Value: "playwright-go"},
				{Name: "language", Value: "golang"},
			},
		},
	}
}

func (r *AllureReporter) StartStep(name string) {
	r.log.Info(name)
	step := &Step{
		Name:   name,
		Status: StatusPassed,
		Start:  time.Now().UnixMilli(),
	}

	if len(r.stepStack) > 0 {
		parent := r.stepStack[len(r.stepStack)-1]
		parent.Steps = append(parent.Steps, *step)
		r.stepStack = append(r.stepStack, &parent.Steps[len(parent.Steps)-1])
	} else {
		r.result.Steps = append(r.result.Steps, *step)
		r.stepStack = append(r.stepStack, &r.result.Steps[len(r.result.Steps)-1])
	}
}

func (r *AllureReporter) StopStep(status Status) {
	if len(r.stepStack) == 0 {
		return
	}
	step := r.stepStack[len(r.stepStack)-1]
	step.Status = status
	step.Stop = time.Now().UnixMilli()
	r.stepStack = r.stepStack[:len(r.stepStack)-1]
}

func (r *AllureReporter) AddScreenshot(screenshotBytes []byte, name string) error {
	uuid, err := uuidG.NewV7()
	if err != nil {
		r.log.Warn("could not generate uuid", "err", err)
	}
	filename := fmt.Sprintf("%s-%d", uuid.String(), time.Now().UnixMilli())
	destPath := filepath.Join(r.outputDir, filename)

	if err := os.WriteFile(destPath, screenshotBytes, 0o600); err != nil {
		return fmt.Errorf("failed to save screenshot attachment: %w", err)
	}

	attachment := Attachment{
		Name:   name,
		Source: filename,
		Type:   "image/png",
	}

	if len(r.stepStack) > 0 {
		step := r.stepStack[len(r.stepStack)-1]
		step.Attachments = append(step.Attachments, attachment)
	} else {
		r.result.Attachments = append(r.result.Attachments, attachment)
	}

	r.log.Info("screenshot attached to allure report", "name", name)
	return nil
}

func (r *AllureReporter) SetFailed(err error) {
	r.result.Status = StatusFailed
	r.result.StatusDetails = &StatusDetails{
		Message: err.Error(),
	}

	for _, step := range r.stepStack {
		step.Status = StatusFailed
		step.Stop = time.Now().UnixMilli()
	}
}

func (r *AllureReporter) SetBroken(err error) {
	r.result.Status = StatusBroken
	r.result.StatusDetails = &StatusDetails{
		Message: err.Error(),
	}
}

func (r *AllureReporter) AddLabel(name, value string) {
	r.result.Labels = append(r.result.Labels, Label{
		Name:  name,
		Value: value,
	})
}

func (r *AllureReporter) SetDescription(desc string) {
	r.result.Description = desc
}

func (r *AllureReporter) Finalize() error {
	r.result.Stop = time.Now().UnixMilli()

	filename := fmt.Sprintf("%s-result.json", r.result.UUID)
	path := filepath.Join(r.outputDir, filename)

	data, err := json.MarshalIndent(r.result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal allure report: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write allure result: %w", err)
	}

	r.log.Info("allure result written", "file", filename, "status", r.result.Status)
	return nil
}
