package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	Logger struct {
		Message any
		Log     *logrus.Logger
	}
)

func Logging(msg string, args ...interface{}) Logger {
	log := Logger{}
	log.Log = logrus.New()
	log.Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.Message = fmt.Sprintf(msg, args...)
	return log
}

func (l Logger) Info() {
	l.Log.Info(l.Message)
}

func (l Logger) Warn() {
	l.Log.Warn(l.Message)
}

func (l Logger) Error() {
	l.Log.Error(l.Message)
}
