package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/milosgajdos/go-playht"
)

var (
	jobID   string
	outPath string
)

func init() {
	flag.StringVar(&jobID, "job-id", "", "TTS job ID")
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
	// * it reads PLAYHT_SECRET_KEY and PLAYHT_USER_ID env vars for auth
	// * uses playht.BaserURL and APIv2 to build the API endpoint URL
	// NOTE: you might need to adjust HTTP client parameters
	// so it does not time out during streaming.
	client := playht.NewClient()

	if err := client.GetTTSJobAudioStream(context.Background(), f, jobID); err != nil {
		log.Fatalf("failed streaming %v job into %s: %v", jobID, outPath, err)
	}

	log.Printf("successfully stored audio file in %s", outPath)
}
