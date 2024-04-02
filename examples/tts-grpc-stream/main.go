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

	if len(outPath) == 0 {
		log.Fatal("you must specify output file path")
	}

	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed creating file %s: %v", outPath, err)
	}
	defer f.Close()

	// Creates an API client with default options.
	// * it reads PLAYHT_SECRET_KEY and PLAYHT_USER_ID env vars
	// * uses playht.BaserURL and APIv2 to create API endpoint URL
	client := playht.NewClient()

	lease, err := client.CreateLease(context.Background(), &playht.CreateLeaseReq{})
	if err != nil {
		log.Fatalf("failed getting a lease: %v", err)
	}

	log.Printf("successfully got PlayHT lease; expires: %s", lease.Expires())
}
