// Sample speech-quickstart uses the Google Cloud Speech API to transcribe
// audio.
package main

import (
	"fmt"
	"log"

	// Imports the Google Cloud Speech API client package.
	"golang.org/x/net/context"

	"cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	opts := option.WithCredentialsFile("secret.json")

	// Creates a client.
	client, err := speech.NewClient(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	gcsURI := "gs://xxx"

	// Detects speech in the audio file.
	op, err := client.LongRunningRecognize(ctx, &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			LanguageCode:    "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		log.Fatalf("failed to recognize wait: %+v", err)
	}

	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}
