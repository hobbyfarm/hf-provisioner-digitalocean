package log

import "github.com/sirupsen/logrus"

var (
	Infof = func(format string, v ...interface{}) {
		logrus.Infof(format, v...)
	}

	Errorf = func(format string, v ...interface{}) {
		logrus.Errorf(format, v...)
	}

	Warnf = func(format string, v ...interface{}) {
		logrus.Warnf(format, v...)
	}

	Debugf = func(format string, v ...interface{}) {
		logrus.Debugf(format, v...)
	}

	Fatalf = func(format string, v ...interface{}) {
		logrus.Fatalf(format, v...)
	}
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// SetLogger overrides the default logger
func SetLogger(logger Logger) {
	Infof = logger.Infof
	Errorf = logger.Errorf
	Warnf = logger.Warnf
	Debugf = logger.Debugf
	Fatalf = logger.Fatalf
}

func SetLogLevel(level logrus.Level) {
	logrus.SetLevel(level)
}
