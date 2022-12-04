// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

// ConsumerGroup the interface for manage consumer group
type ConsumerGroup interface {
	// Init initialize the consumer group.
	Init() (err error)
	// Run start the consumer group.
	Run() (err error)
}

func NewConsumerGroup(config Config) ConsumerGroup {
	// Validate configuration before create the consumer group instance.
	config.validate()
	return &consumerGroup{
		cfg: config,
	}
}

type consumerGroup struct {
	cfg       Config
	saramaCfg *sarama.Config
}

// Init initialize the consumer group.
func (c *consumerGroup) Init() (err error) {
	config := sarama.NewConfig()
	var version sarama.KafkaVersion

	if version, err = sarama.ParseKafkaVersion(string(c.cfg.Version)); err != nil {
		return
	}
	switch c.cfg.BalanceStrategy {
	case Sticky:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case RoundRobin:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case Range:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		err = fmt.Errorf("invalid strategy: %s", c.cfg.BalanceStrategy)
	}

	config.Version = version

	switch c.cfg.Offset {
	case Oldest:
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case Newest:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	sarama.Logger = c.getSaramaLogger()
	c.saramaCfg = config
	return
}

// Run start the consumer group.
func (c *consumerGroup) Run() (err error) {
	keepRunning := true

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(c.cfg.Brokers, c.cfg.Group, c.saramaCfg)
	if err != nil {
		cancel()
		return fmt.Errorf("error creating consumer group client: %s", err)
	}

	consumptionIsPaused := false
	cgh := getConsumerGroupHandler(c.cfg)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = client.Consume(ctx, c.cfg.Topics, cgh); err != nil {
			err = fmt.Errorf("error from consumer: %s", err)
			return
		}
		if ctx.Err() != nil {
			return
		}
		cgh.ready = make(chan bool)

	}()

	<-cgh.ready
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			c.toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		return err
	}
	return nil
}

func (c *consumerGroup) toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
	} else {
		client.PauseAll()
	}
	*isPaused = !*isPaused
}

func (c *consumerGroup) getSaramaLogger() sarama.StdLogger {
	return &customSaramLogger{logger: c.cfg.Logger}
}
