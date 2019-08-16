package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/dghubble/go-twitter/twitter"
)

type Tweet map[string]string

// NewProducer creates producers with the given brokers.
func NewProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}

// Converts a tweet into a message for kafka.
func Tweet2Message(tweet *twitter.Tweet) *sarama.ProducerMessage {
	url := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.IDStr)

	data := map[string]string{
		"id":     tweet.IDStr,
		"author": tweet.User.ScreenName,
		"text":   tweet.Text,
		"url":    url,
		"lang":   tweet.Lang,
	}

	jsonData, _ := json.Marshal(data)

	msg := &sarama.ProducerMessage{
		Topic:     "tweets",
		Partition: -1,
		Value:     sarama.ByteEncoder(jsonData),
	}

	return msg
}
