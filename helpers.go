package playht

import (
	"slices"

	pb "github.com/milosgajdos/go-playht/proto"
)

// MakeGrpcStreamRequest creates a new gRPC stream request from lease and req.
// NOTE: gRPC doesn't provide VoiceEngine and Emotion configuration at the moment.
func MakeGrpcStreamRequest(lease []byte, req *CreateTTSStreamReq) *pb.TtsRequest {
	ttsReq := &pb.TtsRequest{
		Lease: slices.Clone(lease),
		Params: &pb.TtsParams{
			Text:          []string{req.Text},
			Voice:         req.Voice,
			Quality:       ToPbQuality(req.Quality),
			Format:        ToPbFormat(req.OutputFormat),
			SampleRate:    Int32Ptr(req.SampleRate),
			Speed:         Float32Ptr(req.Speed),
			Seed:          Int32Ptr(req.Seed),
			Temperature:   Float32Ptr(req.Temperature),
			StyleGuidance: Float32Ptr(req.StyleGuidance),
			VoiceGuidance: Float32Ptr(req.VoiceGuidance),
			TextGuidance:  Float32Ptr(req.TextGuidance),
		},
	}

	return ttsReq
}

// ToPbQuality converts Quality to its proto representation.
func ToPbQuality(q Quality) *pb.Quality {
	switch q {
	case Draft:
		return qualityPtr(pb.Quality_QUALITY_DRAFT)
	case Low:
		return qualityPtr(pb.Quality_QUALITY_LOW)
	case Medium:
		return qualityPtr(pb.Quality_QUALITY_MEDIUM)
	case High:
		return qualityPtr(pb.Quality_QUALITY_HIGH)
	case Premium:
		return qualityPtr(pb.Quality_QUALITY_PREMIUM)
	default:
		return nil
	}
}

// ToPbFormat converts OutputFormat to its proto representation.
func ToPbFormat(f OutputFormat) *pb.Format {
	switch f {
	case Mp3:
		return formatPtr(pb.Format_FORMAT_MP3)
	case Wav:
		return formatPtr(pb.Format_FORMAT_WAV)
	case Ogg:
		return formatPtr(pb.Format_FORMAT_OGG)
	case Flac:
		return formatPtr(pb.Format_FORMAT_FLAC)
	case Mulaw:
		return formatPtr(pb.Format_FORMAT_MULAW)
	default:
		return nil
	}
}

func qualityPtr(q pb.Quality) *pb.Quality {
	return &q
}

func formatPtr(f pb.Format) *pb.Format {
	return &f
}

// Int32Ptr returns pointer to i.
func Int32Ptr(i int32) *int32 {
	return &i
}

// Float32Ptr returns pointer to f.
func Float32Ptr(f float32) *float32 {
	return &f
}
