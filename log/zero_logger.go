// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package log provides a lightweight logging library with multiple implementations.
package log

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type zeroLogger struct {
	log zerolog.Logger
	err error
}

// NewZeroLogger creates a new zero logger instance.
func NewZeroLogger(conf *Config) Logger {
	log := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Hook(callerHook{dir: conf.dir})
	log.Level(getLogLevel(conf.logLevel))
	return &zeroLogger{
		log: log,
	}
}

func (z zeroLogger) Tracef(format string, args ...interface{}) {
	write(z.log.Trace(), z.err, format, args...)
}

func (z zeroLogger) Debugf(format string, args ...interface{}) {
	write(z.log.Debug(), z.err, format, args...)
}

func (z zeroLogger) Infof(format string, args ...interface{}) {
	write(z.log.Info(), z.err, format, args...)
}

func (z zeroLogger) Warnf(format string, args ...interface{}) {
	write(z.log.Warn(), z.err, format, args...)
}

func (z zeroLogger) Errorf(format string, args ...interface{}) {
	write(z.log.Error(), z.err, format, args...)
}

func (z zeroLogger) Fatalf(format string, args ...interface{}) {
	write(z.log.Fatal(), z.err, format, args...)
}

func (z zeroLogger) WithError(err error) Logger {
	tmp := z
	tmp.err = err
	return tmp
}

func (z zeroLogger) With(keyvals ...interface{}) Logger {
	if len(keyvals)%2 != 0 {
		z.log.Fatal().Msgf("Invalid key value pairs: %s", keyvals)
	}
	length := len(keyvals)
	if length == 0 {
		return z
	}
	logger := z.log.With().Logger()
	for i := 0; i < length; i = i + 2 {
		key := fmt.Sprintf("%v", keyvals[i])
		if val, ok := (keyvals[i+1]).(string); ok {
			logger = logger.With().Str(key, val).Logger()
			continue
		}
		if val, ok := (keyvals[i+1]).([]string); ok {
			logger = logger.With().Strs(key, val).Logger()
			continue
		}
		val := fmt.Sprintf("%v", (keyvals[i+1]))
		logger = logger.With().Str(key, val).Logger()
	}
	tmp := z
	tmp.log = logger
	return &tmp
}

func write(event *zerolog.Event, err error, format string, args ...interface{}) {
	if err != nil {
		event.Err(err).Msgf(format, args...)
		return
	}
	event.Msgf(format, args...)
}

func getLogLevel(level Level) zerolog.Level {
	switch level {
	case Debug:
		return zerolog.DebugLevel
	case Info:
		return zerolog.InfoLevel
	case Warn:
		return zerolog.WarnLevel
	case Error:
		return zerolog.ErrorLevel
	case Fatal:
		return zerolog.FatalLevel
	case Panic:
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
