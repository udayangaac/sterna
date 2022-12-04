// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package avro
package avro

import (
	"sync"

	"github.com/linkedin/goavro/v2"
)

type cachedSchemaRegistryClient struct {
	nativeClient SchemaRegistry
	codecs       map[int]*goavro.Codec
	ids          map[string]int
	mu           sync.RWMutex
}

func NewCachedSchemaRegistry(urls []string, retries int) SchemaRegistry {
	return &cachedSchemaRegistryClient{
		nativeClient: NewSchemaRegistry(urls, retries),
		codecs:       make(map[int]*goavro.Codec),
		ids:          make(map[string]int),
	}
}

func (client *cachedSchemaRegistryClient) GetSchema(id int) (*goavro.Codec, error) {
	client.mu.RLock()
	cachedResult, ok := client.codecs[id]
	client.mu.RUnlock()
	if ok {
		return cachedResult, nil
	}
	result, err := client.nativeClient.GetSchema(id)
	if err != nil {
		return nil, err
	}
	client.mu.Lock()
	client.codecs[id] = result
	client.mu.Unlock()
	return result, nil
}

func (client *cachedSchemaRegistryClient) GetSubjects() ([]string, error) {
	return client.nativeClient.GetSubjects()
}

func (client *cachedSchemaRegistryClient) GetVersions(subject string) ([]int, error) {
	return client.nativeClient.GetVersions(subject)
}

func (client *cachedSchemaRegistryClient) GetSchemaByVersion(subject string, version int) (SchemaDetail, error) {
	return client.nativeClient.GetSchemaByVersion(subject, version)
}

func (client *cachedSchemaRegistryClient) GetLatestSchema(subject string) (SchemaDetail, error) {
	return client.nativeClient.GetLatestSchema(subject)
}

func (client *cachedSchemaRegistryClient) CreateSubject(subject string, codec *goavro.Codec) (int, error) {
	schemaJson := codec.Schema()
	client.mu.RLock()
	cachedResult, ok := client.ids[schemaJson]
	client.mu.RUnlock()
	if ok {
		return cachedResult, nil
	}

	id, err := client.nativeClient.CreateSubject(subject, codec)
	if err != nil {
		return 0, err
	}
	client.mu.Lock()
	client.ids[schemaJson] = id
	client.mu.Unlock()
	return id, nil
}

func (client *cachedSchemaRegistryClient) IsSchemaRegistered(subject string, codec *goavro.Codec) (int, error) {
	return client.nativeClient.IsSchemaRegistered(subject, codec)
}

func (client *cachedSchemaRegistryClient) DeleteSubject(subject string) error {
	return client.nativeClient.DeleteSubject(subject)
}

func (client *cachedSchemaRegistryClient) DeleteVersion(subject string, version int) error {
	return client.nativeClient.DeleteVersion(subject, version)
}
