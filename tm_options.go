package crontab

import "github.com/gohouse/crontab/adapter"

type Options struct {
	persist adapter.Persister
	logger  adapter.Logger
}
type OptionHandleFunc func(options *Options)

func Persist(ps adapter.Persister) OptionHandleFunc {
	return func(options *Options) {
		options.persist = ps
	}
}

func Logger(al adapter.Logger) OptionHandleFunc {
	return func(options *Options) {
		options.logger = al
	}
}
