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
	// * uses playht.BaserURL and APIv2
	client := playht.NewClient()

	voices, err := client.GetVoices(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch voices: %v", err)
	}

	if len(voices) == 0 {
		log.Fatal("no voice found")
	}

	voice := voices[0].ID

	req := &playht.CreateTTSReq{
		Text:         input,
		Voice:        voice,
		Quality:      playht.Low,
		OutputFormat: playht.Mp3,
		Speed:        1.0,
		SampleRate:   24000,
		VoiceEngine:  playht.PlayHTv2,
	}

	job, err := client.CreateTTSJob(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to create a TTS job: %v", err)
	}

	log.Printf("successfully created a new TTS job: %#v", job)

	jobInfo, err := client.GetTTSJob(context.Background(), job.ID)
	if err != nil {
		log.Fatalf("failed getting %v job info: %v", job.ID, err)
	}

	log.Printf("successfully got job info: %#v", jobInfo)
}
