package main

import (
	"context"
	"flag"
	"log"

	"github.com/milosgajdos/go-playht"
)

var (
	input string
)

func init() {
	flag.StringVar(&input, "input", "what is life?", "input text sample")
}

func main() {
	flag.Parse()

	// Creates an API client with default options.
	// * it reads PLAYHT_SECRET_KEY and PLAYHT_USER_ID env vars
	// * uses playht.BaserURL and APIv2 to create API endpoint URL
	client := playht.NewClient()

	voices, err := client.GetVoices(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch voices: %v", err)
	}

	if len(voices) == 0 {
		log.Fatal("no voice found")
	}

	voice := voices[0].ID

	req := &playht.CreateTTSStreamReq{
		Text:         input,
		Voice:        voice,
		Quality:      playht.Low,
		OutputFormat: playht.Mp3,
		Speed:        1.0,
		SampleRate:   24000,
		VoiceEngine:  playht.PlayHTv2Turbo,
	}

	streamURL, err := client.TTSStreamURL(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to create stream URL: %v", err)
	}

	log.Printf("successfully retrieved stream URL: %#v", streamURL)
}
