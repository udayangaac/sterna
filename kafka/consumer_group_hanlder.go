// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"github.com/Shopify/sarama"
)

// ConsumerGroupHandler implementation for ConsumerGroupHandler.
type ConsumerGroupHandler struct {
	cfg   Config
	ready chan bool
}

func getConsumerGroupHandler(cfg Config) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		ready: make(chan bool),
		cfg:   cfg,
	}
}

// Setup setup the consumer group session.
func (c *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup cleanup the consumer group session.
func (c *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim decode messages and call the consumer callback function configured.
func (c *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.cfg.Logger.Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		key, value, err := c.cfg.Decoder(message)
		if err != nil {
			c.cfg.Logger.Errorf("Unable to encode the message. value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		}
		err = c.cfg.ConsumerCallback(key, value)
		if err != nil {
			// If any errors while consuming the message.
			if c.cfg.ConsumerErrorHandler(err) {
				session.MarkMessage(message, "")
			}
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}
