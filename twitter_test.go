package main

import "testing"

func TestStream_TrackKeyword(t *testing.T) {
	track := Stream{}

	track.TrackKeyword("dog")
	track.TrackKeyword("cat")

	if track.filters == nil {
		t.Error("Filters were not added.")
	}

	if track.filters.Track[0] != "dog" {
		t.Error( "Dog filter was not added")
	}

	if track.filters.Track[1] != "cat" {
		t.Error( "Cat filter was not added")
	}
}
