package main

import (
	"errors"
	"fmt"
	schemaregistry "github.com/Landoop/schema-registry"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	producer *TweetProducer
)

func printError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}

func init() {
	// Default kafka settings.
	viper.SetDefault("KAFKA_BROKERS", "localhost:9092")
	viper.SetDefault("SCHEMA_REGISTRY", schemaregistry.DefaultURL)
	viper.SetDefault("TWEET_TOPIC", "tweets")

	// Bind all values to viper.
	viper.BindEnv("kafka_brokers")
	viper.BindEnv("schema_registry")
	viper.BindEnv("tweet_topic")

	viper.BindEnv("client_key")
	viper.BindEnv("client_secret")
	viper.BindEnv("access_token")
	viper.BindEnv("access_secret")

	// Twitter stream settings.
	viper.BindEnv("keywords")

	// Allow to set the config via config file.
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/twitter_kafka_connect")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Ignore config file not found, since env vars can be set.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			printError(err)
		}
	}
}

func main() {
	var err error

	brokers := strings.Split(viper.GetString("kafka_brokers"), ",")

	producer, err = NewTweetProducer(brokers, viper.GetString("schema_registry"), viper.GetString("tweet_topic"))

	if err != nil {
		printError(err)
	}

	clientKey := viper.GetString("client_key")
	clientSecret := viper.GetString("client_secret")
	accessToken := viper.GetString("access_token")
	accessSecret := viper.GetString("access_secret")

	if clientKey == "" || clientSecret == "" || accessToken == "" || accessSecret == "" {
		printError(errors.New("Please set the client key, secret, access token & secrent."))
	}

	stream := NewStream(clientKey, clientSecret, accessToken, accessSecret)

	keywords := viper.GetStringSlice("keywords")


	for _, keyword := range keywords {
		stream.TrackKeyword(keyword)
	}

	stream.OnTweetHandler(func(tweet *twitter.Tweet) {
		url := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.IDStr)
		fmt.Printf("\n\n%v:\n%v", url, tweet.Text)

		err := producer.Post(tweet)

		if err != nil {
			fmt.Println(colorstring.Color("[red] Error while posting tweet: " + err.Error()))
		}
	})

	fmt.Printf("Starting ingesting twitter stream for the keywords: %v", strings.Join(keywords, ", "))

	stream.Run()
	defer stream.Stop()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-ch)

	fmt.Println("Stopping Stream...")
}
