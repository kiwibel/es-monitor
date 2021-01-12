package main

import "testing"

func TestPoller(t *testing.T) {
	want := "Exiting OK"
	if got := pollClusters(); got != want {
		t.Errorf("pollClusters() = %q, want %q", got, want)
	}
}
