package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-playht"
)

func main() {
	// Creates an API client with default options.
	// * it reads PLAYHT_SECRET_KEY and PLAYHT_USER_ID env vars
	// * uses playht.BaserURL and APIv2
	client := playht.NewClient()

	voices, err := client.GetVoices(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch voices: %v", err)
	}

	log.Printf("Got %d voices", len(voices))

	clonedVoices, err := client.GetClonedVoices(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch cloned voices: %v", err)
	}

	log.Printf("Got %d cloned voices", len(clonedVoices))
}
