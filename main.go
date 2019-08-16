package main

import (
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
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

	stream := NewStream(client_key, client_secret, access_token, access_secret)

	stream.FilterKeyword("cats")
	stream.FilterKeyword("dogs")

	stream.OnTweet(func(tweet *twitter.Tweet) {
		url := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.IDStr)
		fmt.Printf("\n\n%v:\n%v", url, tweet.Text)
	})

	stream.Run()
	defer stream.Stop()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-ch)

	fmt.Println("Stopping Stream...")
}
