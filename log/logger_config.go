// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package log provides a lightweight logging library with multiple implementations.
package log

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
	Fatal Level = "fatel"
	Panic Level = "panic"
)

type Config struct {
	logLevel Level
	dir      string
}

func NewConfig() *Config {
	return &Config{
		logLevel: Info,
	}
}

func (c *Config) WithLogLevel(level Level) {
	c.logLevel = level
}

func (c *Config) ProjectDir(dir string) {
	c.dir = dir
}
