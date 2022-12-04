// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package log_test
package log_test

import (
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/udayangaac/sterna/log"
)

func mockTime() {
	zerolog.TimestampFunc = func() time.Time {
		mt := "Mon, 02 Jan 2006 15:04:05 -0700"
		t, _ := time.Parse(time.RFC1123Z, mt)
		return t
	}
}
func ExampleLogger() {
	mockTime()
	logger := log.NewZeroLogger(log.NewConfig())
	logger.Infof("%s", "hello world")
	// Output:
	// {"level":"info","time":"2006-01-02T15:04:05-07:00","file":"./example_test.go:26","message":"hello world"}
}

func ExampleLogger_With() {
	mockTime()
	logger := log.NewZeroLogger(log.NewConfig())

	logger = logger.With("name", "sterna")
	logger.Infof("%s", "hello world")
	// Output:
	// {"level":"info","name":"sterna","time":"2006-01-02T15:04:05-07:00","file":"./example_test.go:36","message":"hello world"}
}

func ExampleLogger_WithError() {
	mockTime()
	logger := log.NewZeroLogger(log.NewConfig())

	err := errors.New("seems we have an error here")
	logger = logger.With("name", "sterna")
	logger.WithError(err).Infof("%s", "hello world")
	// Output:
	// {"level":"info","name":"sterna","error":"seems we have an error here","time":"2006-01-02T15:04:05-07:00","file":"./example_test.go:47","message":"hello world"}

}
