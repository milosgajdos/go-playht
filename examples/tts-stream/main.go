package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/milosgajdos/go-playht"
)

var (
	input   string
	outPath string
)

func init() {
	flag.StringVar(&input, "input", "what is life?", "input text sample")
	flag.StringVar(&outPath, "out", "", "Output file path")
}

func main() {
	flag.Parse()

	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed creating file %s: %v", outPath, err)
	}
	defer f.Close()

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

	if err := client.TTSStream(context.Background(), f, req); err != nil {
		log.Fatalf("failed to stream into %s: %v", outPath, err)
	}

	log.Printf("successfully streamed audio into: %s", outPath)
}
