// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package log
package log

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

type callerHook struct {
	dir string
}

func (h callerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if _, file, line, ok := runtime.Caller(5); ok {
		rel, err := filepath.Rel(h.dir, file)
		if err != nil {
			rel = path.Base(file)
		}
		e.Str("file", fmt.Sprintf("./%v:%v", rel, line))
	}
}
