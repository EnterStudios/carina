package common

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

// Log prints formatted, colored logs to the console
var Log *consoleLogger

func init() {
	Log = &consoleLogger{Logger: &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.WarnLevel,
	}}
}

type consoleLogger struct {
	*logrus.Logger
	IsSilent bool
}

func (log *consoleLogger) SetDebug() {
	log.Level = logrus.DebugLevel
}

func (log *consoleLogger) SetSilent() {
	log.IsSilent = true
	log.Out = ioutil.Discard
}

// Dump does a deep debug dump of a variable
func (log *consoleLogger) Dump(a ...interface{}) {
	log.Debugln(a)
}

// WriteDebug prints debug information to stdout
func (log *consoleLogger) WriteDebug(format string, a ...interface{}) {
	log.Debugf(format, a...)
}

// WriteInfo prints text to stdout
func (log *consoleLogger) WriteInfo(format string, a ...interface{}) {
	log.Infof(format, a...)
}

// WriteWarning prints highlighted text to stdout
func (log *consoleLogger) WriteWarning(format string, a ...interface{}) {
	log.Warnf(format, a...)
}

// WriteError prints highlighted text and an error to stderr
func (log *consoleLogger) WriteError(format string, err error, a ...interface{}) {
	log.Errorf(format, a...)

	if err != nil {
		dumpper := spew.ConfigState{ContinueOnMethod: true}
		dump := dumpper.Sprintln(err)
		log.Error(dump)
	}
}