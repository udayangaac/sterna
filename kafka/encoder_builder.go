// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import "github.com/Shopify/sarama"

// EncoderBuilder resposible for creating encoders.
type EncoderBuilder interface {
	// Build build the sarama.Encoders with the given subject and data.
	Build(subject string, data interface{}) sarama.Encoder
}
