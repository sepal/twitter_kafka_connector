package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	schemaregistry "github.com/Landoop/schema-registry"
	"github.com/Shopify/sarama"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/linkedin/goavro"
	"time"
)

// Tweet represents a twitter post in kafka.
type Tweet map[string]string

// TweetProducer creates Tweet messages in kafka.
type TweetProducer struct {
	producer       sarama.SyncProducer
	schemaRegistry *schemaregistry.Client
	codec          *goavro.Codec
	topic          string
	schema         schemaregistry.Schema
}

// valueSchemaName Returns the subject name for a value schema
func valueSchemaName(topic string) string {
	return topic + "-value"
}

// NewProducer creates producers with the given brokers.
func NewTweetProducer(brokers []string, schemaRegistry string, topic string) (*TweetProducer, error) {
	// Create a new sync kafka producer.
	config := sarama.NewConfig()
	// @todo: Allow users to set the kafka versions.
	config.Version = sarama.V2_3_0_0
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer(brokers, config)

	// Create a new schema registry client.
	r, err := schemaregistry.NewClient(schemaRegistry)

	if err != nil {
		return nil, err
	}

	// Get the latest schema version for this given topic.
	// @todo: Give the option to pin the schema id.
	schema, err := r.GetLatestSchema(valueSchemaName(topic))

	if err != nil {
		return nil, err
	}

	// Create a new avro codec for the schema we fetched.
	codec, err := goavro.NewCodec(schema.Schema)

	if err != nil {
		return nil, err
	}

	producer := TweetProducer{
		producer:       p,
		schemaRegistry: r,
		topic:          topic,
		schema:         schema,
		codec:          codec,
	}

	return &producer, err
}

// Post a tweet to kafka.
func (p *TweetProducer) Post(tweet *twitter.Tweet) error {
	url := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.IDStr)

	// Twitter provides the dates in the format
	// "created_at": "Wed Oct 10 20:19:24 +0000 2018", which seems like ruby
	// format. See:
	// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/tweet-object.html#tweet-dictionary
	t, err := time.Parse(time.RubyDate, tweet.CreatedAt)

	if err != nil {
		return err
	}

	// Build the data we want to post, and convert it first to json, then to the
	// the required avro format.
	data := map[string]interface{}{
		"id":      tweet.ID,
		"timestamp": t.Unix(),
		"author":  tweet.User.ScreenName,
		"text":    tweet.Text,
		"url":     url,
		"lang":    tweet.Lang,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}
	native, _, err := p.codec.NativeFromTextual(jsonData)

	if err != nil {
		return err
	}

	binData, err := p.codec.BinaryFromNative(nil, native)

	if err != nil {
		return err
	}

	// Kafka Avro uses 4 bytes for the schema id. For details check
	// https://docs.confluent.io/current/schema-registry/serializer-formatter.html#wire-format
	binSchemaID := make([]byte, 4)
	binary.BigEndian.PutUint32(binSchemaID, uint32(p.schema.ID))

	// Construct the message expected by kafka avro.
	var binValue []byte
	// The first byte is a magic byte.
	binValue = append(binValue, byte(0))
	// The next 4 bytes is reserved for the schema id.
	binValue = append(binValue, binSchemaID...)
	// The rest is the actual data.
	binValue = append(binValue, binData...)

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(binValue),
	}

	_, _, err = p.producer.SendMessage(msg)

	return err
}
