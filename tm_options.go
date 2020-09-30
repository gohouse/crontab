package crontab

import (
	"github.com/sirupsen/logrus"
)

type Options struct {
	logger  *logrus.Logger
}
type OptionHandleFunc func(options *Options)

func Logger(al *logrus.Logger) OptionHandleFunc {
	if al == nil {
		al = logrus.New()
	}
	return func(options *Options) {
		options.logger = al
	}
}
