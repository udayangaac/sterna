// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"encoding/binary"

	"github.com/Shopify/sarama"
	"github.com/udayangaac/sterna/kafka/avro"
)

// GetAvroDecoder creates avro decoder with avro.SchemaStore.
func GetAvroDecoder(ss avro.SchemaStore) Decoder {
	return func(cm *sarama.ConsumerMessage) (key, value interface{}, err error) {
		key = string(cm.Key)
		schemaId := binary.BigEndian.Uint32(cm.Value[1:5])
		detail := avro.SchemaDetail{}
		detail, err = ss.GetSchemaByID(int(schemaId))
		if err != nil {
			return
		}
		var native interface{}
		native, _, err = detail.Codec.NativeFromBinary(cm.Value[5:])
		if err != nil {
			return
		}

		var textual []byte
		textual, err = detail.Codec.TextualFromNative(nil, native)
		if err != nil {
			return
		}
		value = string(textual)
		return
	}
}
