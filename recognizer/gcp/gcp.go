package gcp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/storage"
	"github.com/kaz/rusuden/recognizer"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

type (
	GcpRecognizer struct {
		bucket string
	}
)

func New(bucket string) recognizer.Recognizer {
	return &GcpRecognizer{bucket}
}

func (r *GcpRecognizer) Recognize(ctx context.Context, data []byte) (string, error) {
	storageCli, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient failed: %w", err)
	}

	obj := storageCli.Bucket(r.bucket).Object("rusuden/" + strconv.FormatInt(time.Now().Unix(), 36) + ".wav")

	w := obj.NewWriter(ctx)
	if _, err := w.Write(data); err != nil {
		return "", fmt.Errorf("w.Write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("w.Close failed: %w", err)
	}

	defer func() {
		if err := obj.Delete(ctx); err != nil {
			log.Printf("cleanup error: obj.Delete failed: %v\n", err)
		}
	}()

	speechCli, err := speech.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("speech.NewClient failed: %w", err)
	}

	op, err := speechCli.LongRunningRecognize(ctx, &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			LanguageCode: "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: "gs://" + obj.BucketName() + "/" + obj.ObjectName()},
		},
	})
	if err != nil {
		return "", fmt.Errorf("speechCli.LongRunningRecognize failed: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return "", fmt.Errorf("op.Wait failed: %w", err)
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			return alt.Transcript, nil
		}
	}
	return "", fmt.Errorf("no result")
}
