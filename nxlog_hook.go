package nxlog

import (
	logrus "github.com/sirupsen/logrus"
	"sync"
)

// Hook to send logs to a service compatible with the NXlog API.
type Hook struct {
	Level     logrus.Level
	Formatter logrus.Formatter
	writer    Writer
	mutex     sync.RWMutex
}

// NewNxlogHook creates a hook to be added to an instance of logrus logger.
func NewNxlogHook(protocol string, endpoint string, options interface{}) (*Hook, error) {
	writer, err := NewWriter(protocol, endpoint, options)
	if err != nil {
		return nil, err
	}

	hook := &Hook{
		Level:     logrus.TraceLevel,
		Formatter: &logrus.JSONFormatter{},
		writer:    *writer,
	}
	return hook, nil
}

// Fire is triggered by a log event. The event might be altered by another hook.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	hook.mutex.RLock()
	defer hook.mutex.RUnlock()

	message, err := hook.getMessage(entry)
	if err != nil {
		return err
	}

	_, err = hook.writer.Write([]byte(message), true)
	if err != nil {
		hook.writer.Close()
		return err
	}
	return nil
}

// Levels returns the available logging levels.
func (hook *Hook) Levels() []logrus.Level {
	levels := []logrus.Level{}
	for _, level := range logrus.AllLevels {
		if level <= hook.Level {
			levels = append(levels, level)
		}
	}
	return levels
}

func (hook *Hook) getMessage(entry *logrus.Entry) (string, error) {
	if hook.Formatter != nil {
		serialized, err := hook.Formatter.Format(entry)
		if err == nil {
			return string(serialized), err
		}
	}
	return entry.String()
}
