package main

import (
	"errors"
	"fmt"
	"github.com/mitchellh/colorstring"
	"os"
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
}
