// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package avro
package avro

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d - %s", e.ErrorCode, e.Message)
}

func newError(resp *http.Response) *Error {
	err := &Error{}
	parsingErr := json.NewDecoder(resp.Body).Decode(&err)
	if parsingErr != nil {
		return &Error{resp.StatusCode, "Unrecognized error found"}
	}
	return err
}
