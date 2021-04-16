/*
This package contains driver(base on logrus) for usage of "slf4go".

slf4go - https://github.com/go-eden/slf4go

logrus - https://github.com/sirupsen/logrus

Use driver

You can use "UseLogrus.xxx()" to register the driver to "slf4go"

  UseLogrus.Default()

Customized loggers

You can constructs your own "*logrus.Logger" with named mapping.

  UseLogrus.WithConfig(LogrousConfig{
    DEFAULT_LOGGER: yourDefaultLogger,
    "log.name.1": yourLogger1,
    "log.name.2": yourLogger2,
  })

Default level

The default level is "slf4go.InfoLevel".
*/
package slf4go_logrus

import (
	l4 "github.com/go-eden/slf4go"
	lr "github.com/sirupsen/logrus"
)

// Method space to use logrus as driver of slf4go
const UseLogrus IUseLogrus = 0

type IUseLogrus int
// Constructs default setting
//
// Log level of default logger: "slf4go.InfoLevel"
func (self IUseLogrus) Default() {
	self.WithConfig(LogrousConfig{})
}
func (IUseLogrus) WithConfig(config LogrousConfig) {
	/**
	 * Sets-up default logger
	 */
	defaultLogger, ok := config[DEFAULT_LOGGER]
	if !ok {
		defaultLogger = lr.New()
		defaultLogger.SetLevel(lr.InfoLevel)
		config[DEFAULT_LOGGER] = defaultLogger
	}
	// :~)

	/**
	 * Sets-up cache of loggers(and levels)
	 */
	newDriverInstance := newLogrusDriver()
	newDriverInstance.defaultLevel = levelMappingFromLogrus[config[DEFAULT_LOGGER].Level]
	for name, logger := range config {
		newDriverInstance.loggerMap[name] = logger
		newDriverInstance.loggerLevelMap[name] = levelMappingFromLogrus[logger.Level]
	}
	// :~)

	l4.SetDriver(newDriverInstance)
}

//
type logrusDriver struct {
	defaultLevel l4.Level

	loggerMap map[string]*lr.Logger
	loggerLevelMap map[string]l4.Level
}
func newLogrusDriver() *logrusDriver {
	return &logrusDriver{
		defaultLevel: l4.InfoLevel,

		loggerMap: make(map[string]*lr.Logger),
		loggerLevelMap: make(map[string]l4.Level),
	}
}

func (self *logrusDriver) Name() string {
	return "slf4go-logrus"
}
func (self *logrusDriver) Print(l *l4.Log) {
	/**
	 * Sets-up fields with "logger" field
	 */
	fields := lr.Fields{}
	if l.Fields != nil {
		fields = lr.Fields(l.Fields)
	}

	// Assigns the logName
	fields["logName"] = l.Logger
	// :~)

	/**
	 * Gets logger from cache or use default one
	 */
	logger, _ := self.getLogger(DEFAULT_LOGGER)
	if cachedLogger, ok := self.getLogger(l.Logger); ok {
		logger = cachedLogger
	}
	// :~)

	entry := lr.NewEntry(logger)
	entry = entry.WithFields(fields)

	switch l.Level {
	case l4.TraceLevel:
		if l.Format == nil {
			entry.Trace(l.Args...)
		} else {
			entry.Tracef(*l.Format, l.Args...)
		}
	case l4.DebugLevel:
		if l.Format == nil {
			entry.Debug(l.Args...)
		} else {
			entry.Debugf(*l.Format, l.Args...)
		}
	case l4.InfoLevel:
		if l.Format == nil {
			entry.Info(l.Args...)
		} else {
			entry.Infof(*l.Format, l.Args...)
		}
	case l4.WarnLevel:
		if l.Format == nil {
			entry.Warn(l.Args...)
		} else {
			entry.Warnf(*l.Format, l.Args...)
		}
	case l4.ErrorLevel:
		if l.Format == nil {
			entry.Error(l.Args...)
		} else {
			entry.Errorf(*l.Format, l.Args...)
		}
	case l4.PanicLevel:
		if l.Format == nil {
			entry.Panic(l.Args...)
		} else {
			entry.Panicf(*l.Format, l.Args...)
		}
	case l4.FatalLevel:
		if l.Format == nil {
			entry.Fatal(l.Args...)
		} else {
			entry.Fatalf(*l.Format, l.Args...)
		}
	}
}

func (self *logrusDriver) GetLevel(name string) l4.Level {
	level, ok := self.getLevel(name)
	if ok {
		return level
	}

	return self.defaultLevel
}

func (self *logrusDriver) getLevel(name string) (l4.Level, bool) {
	level, ok := self.loggerLevelMap[name]
	return level, ok
}
func (self *logrusDriver) getLogger(name string) (*lr.Logger, bool) {
	logger, ok := self.loggerMap[name]
	return logger, ok
}

var levelMappingFromLogrus = map[lr.Level]l4.Level {
	lr.TraceLevel: l4.TraceLevel,
	lr.DebugLevel: l4.DebugLevel,
	lr.InfoLevel: l4.InfoLevel,
	lr.WarnLevel: l4.WarnLevel,
	lr.ErrorLevel: l4.ErrorLevel,
	lr.PanicLevel: l4.PanicLevel,
	lr.FatalLevel: l4.FatalLevel,
}
