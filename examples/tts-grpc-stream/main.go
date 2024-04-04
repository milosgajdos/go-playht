package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"

	"github.com/milosgajdos/go-playht"
	pb "github.com/milosgajdos/go-playht/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, playht.GrpcAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		log.Fatalf("failed creating gRPC connection: %v", err)
	}
	defer conn.Close()

	// Creates an API client with default options.
	// * it reads PLAYHT_SECRET_KEY and PLAYHT_USER_ID env vars
	// * uses playht.BaserURL and APIv2 to create API endpoint URL
	client := playht.NewClient(playht.WithGRPCClient(conn))

	lease, err := client.CreateLease(ctx, &playht.CreateLeaseReq{})
	if err != nil {
		log.Fatalf("failed getting a lease: %v", err)
	}

	log.Printf("successfully got PlayHT lease; expires: %s", lease.Expires())

	voices, err := client.GetVoices(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch voices: %v", err)
	}

	if len(voices) == 0 {
		log.Fatal("no voice found")
	}

	voice := voices[0].ID

	log.Printf("using voice: %s", voice)

	req := &pb.TtsRequest{
		Lease: lease.Data,
		Params: &pb.TtsParams{
			Text:       []string{input},
			Voice:      voice,
			Quality:    playht.ToPbQuality(playht.Low),
			Format:     playht.ToPbFormat(playht.Mp3),
			Speed:      playht.Float32Ptr(1.0),
			SampleRate: playht.Int32Ptr(24000),
		},
	}

	if err := client.TTSGrpcStream(ctx, f, req); err != nil {
		log.Fatalf("gRPC stream failed: %v", err)
	}

	log.Printf("successfully stream into %v", outPath)
}
