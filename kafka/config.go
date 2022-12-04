// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"github.com/udayangaac/sterna/log"
)

const (
	Version_2_1_1 Version = "2.1.1"

	Sticky     BalanceStrategy = "sticky"
	RoundRobin BalanceStrategy = "roundrobin"
	Range      BalanceStrategy = "range"

	Newest Offset = "newest"
	Oldest Offset = "oldest"
)

// Config General configurations.
type Config struct {
	Brokers              []string
	Group                string
	Version              Version
	Topics               []string
	BalanceStrategy      BalanceStrategy
	Offset               Offset
	Logger               log.Logger
	LogLevel             log.Level
	EncoderBuilder       EncoderBuilder
	Decoder              Decoder
	ConsumerCallback     ConsumerCallback
	ConsumerErrorHandler ConsumerErrorHandler
}

// validate verify the configurations and set default values.
func (c *Config) validate() {
	if c.Logger == nil {
		ll := log.Info
		if c.LogLevel != "" {
			ll = c.LogLevel
		}
		conf := log.NewConfig()
		conf.WithLogLevel(ll)
		c.Logger = log.NewZeroLogger(conf)
	}
	if c.ConsumerCallback == nil {
		c.Logger.Fatalf("Please define a consumer callback in avro.Config")
	}
	if c.EncoderBuilder == nil {
		c.Logger.Fatalf("Please define an Encoder builder")
	}
	if c.ConsumerErrorHandler == nil {
		c.ConsumerErrorHandler = func(err error) (commitMsg bool) {
			c.Logger.WithError(err).Errorf("Unable read the message")
			return true
		}
	}
}
