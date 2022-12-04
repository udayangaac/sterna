// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"encoding/binary"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/udayangaac/sterna/kafka/avro"
)

type avroEncoderBuilder struct {
	schemaStore avro.SchemaStore
}

// NewAvroEncoderBuilder create instance of EncoderBuilder.
func NewAvroEncoderBuilder(schemaStore avro.SchemaStore) EncoderBuilder {
	return &avroEncoderBuilder{schemaStore: schemaStore}
}

func (aeb *avroEncoderBuilder) buildBinayMessage(subject string, data interface{}) (binaryMsg []byte, err error) {
	var (
		detail     avro.SchemaDetail
		binaryData []byte
	)
	binaryMsg = append(binaryMsg, byte(0))
	binarySchemaId := make([]byte, 4)
	detail, err = aeb.schemaStore.GetSchemaBySubject(subject)
	if err != nil {
		return
	}
	binaryData, err = json.Marshal(data)
	if err != nil {
		return
	}
	binary.BigEndian.PutUint32(binarySchemaId, uint32(detail.Version))
	binaryMsg = append(binaryMsg, binarySchemaId...)
	binaryMsg = append(binaryMsg, binaryData...)
	return
}

// Build creates sarama.Encoder for given subject and data.
func (aeb *avroEncoderBuilder) Build(subject string, data interface{}) sarama.Encoder {
	binaryData, err := aeb.buildBinayMessage(subject, data)
	ae := avroEncoder{
		binaryData: binaryData,
		err:        err,
	}
	return &ae
}

// avroEncoder implemtation of sarama.Encoder
type avroEncoder struct {
	binaryData []byte
	err        error
}

func (a *avroEncoder) Encode() (binaryMessge []byte, err error) {
	return a.binaryData, a.err
}

func (a *avroEncoder) Length() int {
	return 5 + len(a.binaryData)
}
