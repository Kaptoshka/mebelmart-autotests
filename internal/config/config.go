package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type BrowserType string

const (
	BrowserChromium BrowserType = "chromium"
	BrowserFirefox  BrowserType = "firefox"
	BrowserWebKit   BrowserType = "webkit"
)

const (
	defaultTimeoutMS      = 30000
	defaultSlowMoMS       = 0
	defaultViewportWidth  = 1920
	defaultViewportHeight = 1080

	DefaultLogDir    = "./artifacts/logs"
	DefaultTracesDir = "./artifacts/traces"
)

type Config struct {
	Browser         BrowserType
	Headless        bool
	Trace           bool
	BaseURL         string
	Timeout         time.Duration
	SlowMo          time.Duration
	AllureReportDir string
	LogLevel        string
	ViewportWidth   int
	ViewportHeight  int
}

func Load() *Config {
	_ = godotenv.Load(os.Getenv("ENV_FILE"))

	return &Config{
		Browser:         getBrowserType(),
		Headless:        getBool("HEADLESS", true),
		Trace:           getBool("TRACE", false),
		BaseURL:         getEnv("BASE_URL", "https://example.com"),
		Timeout:         getDuration("TIMEOUT_MS", defaultTimeoutMS),
		SlowMo:          getDuration("SLOW_MO_MS", defaultSlowMoMS),
		AllureReportDir: getEnv("ALLURE_RESULTS_DIR", "./allure-results"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		ViewportWidth:   getInt("VIEWPORT_WIDTH", defaultViewportWidth),
		ViewportHeight:  getInt("VIEWPORT_HEIGHT", defaultViewportHeight),
	}
}

func getBrowserType() BrowserType {
	b := getEnv("BROWSER", "chrome")
	switch BrowserType(b) {
	case BrowserChromium:
		return BrowserChromium
	case BrowserFirefox:
		return BrowserFirefox
	case BrowserWebKit:
		return BrowserWebKit
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
		return time.Duration(defaultMs)
	}
	ms, err := strconv.Atoi(v)
	if err != nil {
		return time.Duration(defaultMs)
	}
	return time.Duration(ms)
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
