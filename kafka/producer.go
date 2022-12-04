// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// Producer kafka message producer interface.
type Producer interface {
	// Produce produce the kafka message to the given topic.
	Produce(topic string, schema string, key interface{}, value interface{}) (partition int32, offset int64, err error)
}

type producer struct {
	syncProd sarama.SyncProducer
	cfg      Config
}

// NewProducer creates a producer instance.
func NewProducer(cfg Config) Producer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_1_0
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionNone
	config.Producer.MaxMessageBytes = 10000000
	config.Producer.Retry.Max = 10
	config.Producer.Retry.Backoff = 1000 * time.Millisecond
	p, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		cfg.Logger.WithError(err).Fatalf("Unable to create producer.")
	}
	return &producer{
		syncProd: p,
		cfg:      cfg,
	}
}

// Produce produce the kafka message to the given topic.
func (p *producer) Produce(topic string, schema string, key interface{}, value interface{}) (partition int32, offset int64, err error) {
	keyStr, ok := value.(string)
	if !ok {
		err = fmt.Errorf("key should be string. got: %s", value)
		p.cfg.Logger.WithError(err).Errorf("Invalid key")
		return
	}
	valueEncoder := p.cfg.EncoderBuilder.Build(schema, value)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(keyStr),
		Value: valueEncoder,
	}
	return p.syncProd.SendMessage(msg)
}

func (ap *producer) Close() {
	ap.syncProd.Close()
}
