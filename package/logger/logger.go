package logger

import (
	"pcstakehometest/utilities"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
)

type LogRus struct {
	*logrus.Entry
}

func NewLogRus() *LogRus {
	logger := logrus.New()

	newEntry := logrus.NewEntry(logger)
	return &LogRus{newEntry}
}

func (l *LogRus) Request() *LogRus {
	formatter := runtime.Formatter{
		File:         true,
		Package:      true,
		BaseNameOnly: true,
		Line:         true,
		ChildFormatter: &logrus.JSONFormatter{
			DataKey:     utilities.RandomString(20),
			PrettyPrint: false,
		}}
	l.Logger.SetFormatter(&formatter)
	return l
}
