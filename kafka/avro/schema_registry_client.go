// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package avro
package avro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/linkedin/goavro/v2"
)

const (
	schemaByID       = "/schemas/ids/%d"
	subjects         = "/subjects"
	subjectVersions  = "/subjects/%s/versions"
	deleteSubject    = "/subjects/%s"
	subjectByVersion = "/subjects/%s/versions/%s"

	latestVersion = "latest"

	contentType = "application/vnd.schemaregistry.v1+json"

	defaultTimeout = 2 * time.Second
)

type schemaRegistryClient struct {
	SchemaRegistryBaseUrls []string
	httpClient             *http.Client
	retries                int
}

type schemaResponse struct {
	Schema string `json:"schema"`
}

type schemaVersionResponse struct {
	Subject string `json:"subject"`
	Version int    `json:"version"`
	Schema  string `json:"schema"`
	ID      int    `json:"id"`
}

type idResponse struct {
	ID int `json:"id"`
}

func NewSchemaRegistry(urls []string, retries int) SchemaRegistry {
	if retries < 0 {
		retries = len(urls)
	}
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return &schemaRegistryClient{
		SchemaRegistryBaseUrls: urls,
		retries:                retries,
		httpClient:             client,
	}
}

func (client *schemaRegistryClient) GetSchema(id int) (*goavro.Codec, error) {
	resp, err := client.httpCall("GET", fmt.Sprintf(schemaByID, id), nil)
	if nil != err {
		return nil, err
	}
	schema, err := parseSchema(resp)
	if nil != err {
		return nil, err
	}
	return goavro.NewCodec(schema.Schema)
}

func (client *schemaRegistryClient) GetSubjects() ([]string, error) {
	resp, err := client.httpCall("GET", subjects, nil)
	if nil != err {
		return []string{}, err
	}
	var result = []string{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (client *schemaRegistryClient) GetVersions(subject string) ([]int, error) {
	resp, err := client.httpCall("GET", fmt.Sprintf(subjectVersions, subject), nil)
	if nil != err {
		return []int{}, err
	}
	var result = []int{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (client *schemaRegistryClient) GetSchemaByVersion(subject string, version int) (SchemaDetail, error) {
	return client.getSchemaByVersionInternal(subject, fmt.Sprintf("%d", version))
}

func (client *schemaRegistryClient) GetLatestSchema(subject string) (SchemaDetail, error) {
	return client.getSchemaByVersionInternal(subject, latestVersion)
}

func (client *schemaRegistryClient) CreateSubject(subject string, codec *goavro.Codec) (int, error) {
	schema := schemaResponse{codec.Schema()}
	json, err := json.Marshal(schema)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(json)
	resp, err := client.httpCall("POST", fmt.Sprintf(subjectVersions, subject), payload)
	if err != nil {
		return 0, err
	}
	return parseID(resp)
}

func (client *schemaRegistryClient) IsSchemaRegistered(subject string, codec *goavro.Codec) (int, error) {
	schema := schemaResponse{codec.Schema()}
	json, err := json.Marshal(schema)
	if err != nil {
		return 0, err
	}
	payload := bytes.NewBuffer(json)
	resp, err := client.httpCall("POST", fmt.Sprintf(deleteSubject, subject), payload)
	if err != nil {
		return 0, err
	}
	return parseID(resp)
}

func (client *schemaRegistryClient) DeleteSubject(subject string) error {
	_, err := client.httpCall("DELETE", fmt.Sprintf(deleteSubject, subject), nil)
	return err
}

func (client *schemaRegistryClient) DeleteVersion(subject string, version int) error {
	_, err := client.httpCall("DELETE", fmt.Sprintf(subjectByVersion, subject, fmt.Sprintf("%d", version)), nil)
	return err
}

func (client *schemaRegistryClient) getSchemaByVersionInternal(subject string, version string) (SchemaDetail, error) {
	resp, err := client.httpCall("GET", fmt.Sprintf(subjectByVersion, subject, version), nil)
	if nil != err {
		return SchemaDetail{}, err
	}
	var schema = new(schemaVersionResponse)
	err = json.Unmarshal(resp, &schema)
	if nil != err {
		return SchemaDetail{}, err
	}
	var codec = new(goavro.Codec)
	codec, err = goavro.NewCodec(schema.Schema)
	if err != nil {
		return SchemaDetail{}, err
	}
	return SchemaDetail{
		ID:      schema.ID,
		Subject: schema.Subject,
		Version: schema.Version,
		Codec:   codec,
	}, nil
}

// Calls all base urls.
func (client *schemaRegistryClient) httpCall(method, uri string, payload io.Reader) ([]byte, error) {
	nServers := len(client.SchemaRegistryBaseUrls)
	offset := rand.Intn(nServers)
	for i := 0; ; i++ {
		url := fmt.Sprintf("%s%s", client.SchemaRegistryBaseUrls[(i+offset)%nServers], uri)
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
		resp, err := client.httpClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if i < client.retries && (err != nil || retriable(resp)) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !okStatus(resp) {
			return nil, newError(resp)
		}
		return ioutil.ReadAll(resp.Body)
	}
}

func retriable(resp *http.Response) bool {
	return resp.StatusCode >= 500 && resp.StatusCode < 600
}

func okStatus(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func parseSchema(str []byte) (*schemaResponse, error) {
	var schema = new(schemaResponse)
	err := json.Unmarshal(str, &schema)
	return schema, err
}

func parseID(str []byte) (int, error) {
	var id = new(idResponse)
	err := json.Unmarshal(str, &id)
	return id.ID, err
}
