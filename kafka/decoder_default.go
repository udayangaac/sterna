// Copyright (C) By Chamith Udayanga - All Rights Reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Chamith Udayanga <udayangaac@gmail.com>, February 2022

// Package kafka
package kafka

import (
	"github.com/Shopify/sarama"
)

// GetDefaultDecoder creates a default decoder (string decoder).
func GetDefaultDecoder() Decoder {
	return func(cm *sarama.ConsumerMessage) (key interface{}, value interface{}, err error) {
		value = string(cm.Value)
		key = string(cm.Key)
		return
	}
}
