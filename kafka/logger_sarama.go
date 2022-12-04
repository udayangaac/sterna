// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"fmt"
	"strings"

	"github.com/udayangaac/sterna/log"
)

const saramaPrefix = "sarama library"

// customSaramLogger custom logger implementation for saram.
type customSaramLogger struct {
	logger log.Logger
}

func (c *customSaramLogger) Print(v ...interface{}) {
	c.logger.Infof(c.getLogFormat(saramaPrefix, v), v...)
}
func (c *customSaramLogger) Printf(format string, v ...interface{}) {
	c.logger.Infof("prefix=%v"+format, append([]interface{}{saramaPrefix}, v...)...)
}
func (c *customSaramLogger) Println(v ...interface{}) {
	c.logger.Infof(c.getLogFormat(saramaPrefix, v), v...)
}

func (c *customSaramLogger) getLogFormat(prefix string, v ...interface{}) (format string) {
	format = fmt.Sprintf("prefix=%v", prefix)
	l := len(v)
	if l == 0 {
		return
	}
	format = format + " %v"
	format = format + strings.Repeat(", %v", l-1)
	return
}
