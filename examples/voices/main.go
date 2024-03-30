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
	flag.StringVar(&input, "input", "", "input voice sample")
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

	log.Printf("Got %d voices", len(voices))

	clonedVoices, err := client.GetClonedVoices(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch cloned voices: %v", err)
	}

	log.Printf("Got %d cloned voices", len(clonedVoices))

	if input != "" {
		req := &playht.CloneVoiceFileRequest{
			SampleFile: input,
			VoiceName:  "foo-bar",
			MimeType:   "audio/x-m4a",
		}
		cloneResp, err := client.CreateInstantVoiceCloneFromFile(context.Background(), req)
		if err != nil {
			log.Fatalf("failed to clone voice from file: %v", err)
		}
		log.Printf("clone voice response: %v", cloneResp)

		del := &playht.DeleteClonedVoiceRequest{
			VoiceID: "s3://voice-cloning-zero-shot/3c3ab1f6-d42e-4d25-8660-4b31f7e1073e/foo-bar/manifest.json",
		}

		delResp, err := client.DeleteClonedVoice(context.Background(), del)
		if err != nil {
			log.Fatalf("failed to delete %s: %v", del.VoiceID, err)
		}
		log.Printf("voice %s successfully deleted: %v", del.VoiceID, delResp)
	}
}
