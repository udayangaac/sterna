// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package avro
package avro

import (
	"github.com/linkedin/goavro/v2"
)

type SchemaDetail struct {
	Subject string
	Version int
	Schema  string
	ID      int
	Codec   *goavro.Codec
}
type SchemaRegistry interface {
	GetSchema(int) (*goavro.Codec, error)
	GetSubjects() ([]string, error)
	GetVersions(string) ([]int, error)
	GetSchemaByVersion(string, int) (SchemaDetail, error)
	GetLatestSchema(string) (SchemaDetail, error)
	CreateSubject(string, *goavro.Codec) (int, error)
	IsSchemaRegistered(string, *goavro.Codec) (int, error)
	DeleteSubject(string) error
	DeleteVersion(string, int) error
}
