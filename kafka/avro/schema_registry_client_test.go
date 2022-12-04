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
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSchemaRegistryClient_Retries(t *testing.T) {
	count := 0
	response := []string{"test"}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.Header().Set("Content-Type", contentType)
		if count < 3 {
			http.Error(w, `{"error_code": 500, "message": "Error in the backend datastore"}`, 500)
		} else {
			str, _ := json.Marshal(response)
			fmt.Fprintf(w, string(str))
		}
	}))
	SchemaRegistryClient := NewSchemaRegistry([]string{mockServer.URL}, 2)
	subjects, err := SchemaRegistryClient.GetSubjects()
	if err != nil {
		t.Errorf("Found error %s", err)
	}
	if !reflect.DeepEqual(subjects, response) {
		t.Errorf("Subjects did not match expected %s, got %s", response, subjects)
	}
	expectedCallCount := 3
	if count != expectedCallCount {
		t.Errorf("Expected error count to be %d, got %d", expectedCallCount, count)
	}
}

func TestSchemaRegistryClient_Error(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error_code": 500, "message": "Error in the backend datastore"}`, 500)
	}))
	SchemaRegistryClient := NewSchemaRegistry([]string{mockServer.URL}, -1)
	_, err := SchemaRegistryClient.GetSubjects()
	expectedErr := Error{500, "Error in the backend datastore"}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error to be %s, got %s", expectedErr.Error(), err.Error())
	}
}
