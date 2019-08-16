package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Stream represents the twitter stream client containing the filter params,
// and the actual twitter client
type Stream struct {
	client  *twitter.Client
	stream  *twitter.Stream
	filters *twitter.StreamFilterParams
	demux   twitter.SwitchDemux
}

// NewStream creates a new stream object.
func NewStream(clientKey string, clientSecret string, accessKey string, accessToken string) *Stream {
	config := oauth1.NewConfig(client_key, client_secret)
	token := oauth1.NewToken(access_token, access_secret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	stream := &Stream{
		client: client,
		demux:  twitter.NewSwitchDemux(),
	}

	return stream
}

// FilterKeyword filters by the given keyword.
func (stream *Stream) FilterKeyword(keyword string) {
	if stream.filters == nil {
		stream.filters = &twitter.StreamFilterParams{
			Track: []string{keyword},
		}
		return
	}

	if stream.filters.Track == nil {
		stream.filters.Track = []string{keyword}
		return
	}

	stream.filters.Track = append(stream.filters.Track, keyword)
}

// OnTweetHandler registers a new handler function for when a tweet is posted.
func (stream *Stream) OnTweetHandler(handler func(tweet *twitter.Tweet)) {
	stream.demux.Tweet = handler
}

// Run creates a new twitter stream object and starts a go routine to listen to
// new tweets or events.
func (stream *Stream) Run() {
	stream.stream, _ = stream.client.Streams.Filter(stream.filters)
	go stream.demux.HandleChan(stream.stream.Messages);
}

// Stop stops the twitter stream.
func (stream *Stream) Stop() {
	stream.stream.Stop()
}
