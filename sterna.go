// Copyright 2021 The Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sterna consist of all functions of microservice framework.
package sterna

// Sterna sterna consist of configuration details.
type Sterna struct {
}

// NewSterna creates a new instance of Sterna with default configurations.
func NewSterna() *Sterna {
	return &Sterna{}
}

// Start starts the application.
func (s *Sterna) Start() {}
