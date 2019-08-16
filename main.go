package main

import (
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mitchellh/colorstring"
	"os"
	"os/signal"
	"syscall"
)

var (
	client_key    = ""
	client_secret = ""
	access_token  = ""
	access_secret = ""
)

func printError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}

func init() {
	client_key = os.Getenv("CLIENT_KEY")
	client_secret = os.Getenv("CLIENT_SECRET")
	access_token = os.Getenv("ACCESS_TOKEN")
	access_secret = os.Getenv("ACCESS_SECRET")
}

func main() {
	if client_key == "" || client_secret == "" || access_token == "" || access_secret == "" {
		printError(errors.New("Please set the client key, secret, access token & secrent"))
	}

	config := oauth1.NewConfig(client_key, client_secret)
	token := oauth1.NewToken(access_token, access_secret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		url := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.IDStr)
		fmt.Printf("\n\n%v:\n%v", url, tweet.Text)

		media_url := ""
		if len(tweet.Entities.Media) > 0 {
			media_url = tweet.Entities.Media[0].URL
			fmt.Printf("\n- Media: %v", media_url)
		}
	}

	filters := &twitter.StreamFilterParams{Track: []string{"kitten", "cat", "puppy", "dogs"}}

	stream, err := client.Streams.Filter(filters)

	if err != nil {
		printError(err)
	}
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
