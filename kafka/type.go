// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/udayangaac/sterna/kafka/avro"
)

// Version kafka version.
type Version string

// BalanceStrategy consumer group balancing strategy.
type BalanceStrategy string

// Offset message offset of the consumer group.
type Offset string

// ConsumerCallback hanler function for the consumer.
type ConsumerCallback func(key, value interface{}) (err error)

// ConsumerErrorHandler hanlde the error returning from the ConsumerCallback.
type ConsumerErrorHandler func(err error) (commitMsg bool)

// Decoder decode the consumer message to according to the given implementation
type Decoder func(consumerMessage *sarama.ConsumerMessage) (key, value interface{}, err error)

type Encoder func(schemaStore avro.SchemaStore, schema string, data []byte) sarama.Encoder
