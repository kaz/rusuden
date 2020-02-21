package gcp

import (
	"context"
	"fmt"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/kaz/rusuden/recognizer"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

type (
	GcpRecognizer struct{}
)

func New() recognizer.Recognizer {
	return &GcpRecognizer{}
}

func (r *GcpRecognizer) Recognize(ctx context.Context, data []byte) (string, error) {
	client, err := speech.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("speech.NewClient failed: %w", err)
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			LanguageCode: "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		return "", fmt.Errorf("client.Recognize failed: %w", err)
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			return alt.Transcript, nil
		}
	}
	return "", fmt.Errorf("no result")
}
