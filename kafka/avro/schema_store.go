// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package avro
package avro

import (
	"fmt"
	"sync"
)

type SchemaVersion int

const (
	LatestVersion SchemaVersion = -1
)

type Schema struct {
	Subject string
	Version SchemaVersion
}

type SchemaStore interface {
	GetSchemaBySubject(subject string) (detail SchemaDetail, err error)
	GetSchemaByID(id int) (detail SchemaDetail, err error)
}

type schemaStore struct {
	ids           map[string]int
	schemaDetails map[int]SchemaDetail
	mu            sync.RWMutex
}

func NewSchemaStore(schemaReg SchemaRegistry, schemas []Schema) (SchemaStore, error) {
	ss := schemaStore{}
	ss.schemaDetails = make(map[int]SchemaDetail)
	ss.ids = make(map[string]int)

	var (
		schemaDetail SchemaDetail
		err          error
	)
	for _, schema := range schemas {
		if schema.Version == LatestVersion {
			schemaDetail, err = schemaReg.GetLatestSchema(schema.Subject)
		} else {
			schemaDetail, err = schemaReg.GetSchemaByVersion(schema.Subject, int(schema.Version))
		}
		if err != nil {
			return nil, err
		}
		ss.ids[schema.Subject] = schemaDetail.ID
		ss.schemaDetails[schemaDetail.ID] = schemaDetail
	}
	return &ss, nil
}

func (ss *schemaStore) GetSchemaByID(id int) (detail SchemaDetail, err error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	ok := false
	detail, ok = ss.schemaDetails[id]
	if !ok {
		err = fmt.Errorf("schema id %v was not added", id)
		return
	}
	return
}

func (ss *schemaStore) GetSchemaBySubject(subject string) (detail SchemaDetail, err error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	ok := false
	id, ok := ss.ids[subject]
	if !ok {
		err = fmt.Errorf("schema %s was not added", subject)
		return
	}
	detail, ok = ss.schemaDetails[id]
	if !ok {
		err = fmt.Errorf("schema %s was not added", subject)
		return
	}
	return
}
