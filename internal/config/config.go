package config

import (
	"os"
	"strconv"
	"time"
)

type BrowserType string

const (
	BrowserChromium BrowserType = "chromium"
	BrowserFirefox  BrowserType = "firefox"
)

type Config struct {
	Browser         BrowserType
	Headless        bool
	BaseURL         string
	Timeout         time.Duration
	SlowMo          time.Duration
	ScreenshotsDir  string
	AllureReportDir string
	LogLevel        string
	ViewportWidth   int
	ViewportHeight  int
}

func Load() *Config {
	return &Config{
		Browser:         getBrowserType(),
		Headless:        getBool("HEADLESS", true),
		BaseURL:         getEnv("BASE_URL", "https://example.com"),
		Timeout:         getDuration("TIMEOUT_MS", 30000),
		SlowMo:          getDuration("SLOW_MO_MS", 0),
		ScreenshotsDir:  getEnv("SCREENSHOTS_DIR", "./screenshots"),
		AllureReportDir: getEnv("ALLURE_RESULTS_DIR", "./allure-results"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		ViewportWidth:   getInt("VIEWPORT_WIDTH", 1920),
		ViewportHeight:  getInt("VIEWPORT_HEIGHT", 1080),
	}
}

func getBrowserType() BrowserType {
	b := getEnv("BROWSER", "chrome")
	switch BrowserType(b) {
	case BrowserFirefox:
		return BrowserFirefox
	default:
		return BrowserChromium
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}

func getDuration(key string, defaultMs int) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return time.Duration(defaultMs) * time.Millisecond
	}
	ms, err := strconv.Atoi(v)
	if err != nil {
		return time.Duration(defaultMs) * time.Millisecond
	}
	return time.Duration(ms) * time.Millisecond
}

func getInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return i
}
