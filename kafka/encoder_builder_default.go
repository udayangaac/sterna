// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
)

type defaultEncoderBuilder struct{}

// DefaultEncoderBuilder creates instance of EncoderBuilder.
func DefaultEncoderBuilder() EncoderBuilder {
	return &defaultEncoderBuilder{}
}

// Build creates string encoder.
func (aeb *defaultEncoderBuilder) Build(_ string, data interface{}) sarama.Encoder {
	binaryMsg, _ := json.Marshal(data)
	se := sarama.StringEncoder(string(binaryMsg))
	return &se
}
