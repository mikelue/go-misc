package slf4go_logrus

import (
	lr "github.com/sirupsen/logrus"
)

// Name of default logger
const DEFAULT_LOGGER = "<DEFAULT>"

// Configurations for logrous, which maps name to an instance of "*logrus.Logger"
type LogrousConfig map[string]*lr.Logger

// Text configuration for Logrous
const (
	levelTrace = "trace"
	levelDebug = "debug"
	levelInfo = "info"
	levelWarn = "warn"
	levelError = "error"
	levelPanic = "panic"
	levelFatal = "fatal"

	outputStdout = "stdout"

	formatter_json = "json"
	formatter_text = "text"
)
