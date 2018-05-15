package main

import (
	"cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"log"
	"io"
	"fmt"
	"os"
	"context"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 300 * time.Second)
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding: speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 48000,
					LanguageCode: "ja-JP",
				},
				InterimResults: true,
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	f := os.Stdin

	log.Println("start recognition")
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("Could not close stream: %v", err)
				}
				return
			}
			if err != nil {
				log.Printf("Could not read from input: %v", err)
				continue
			}
			if err = stream.Send(&speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
					AudioContent: buf[:n],
				},
			}); err != nil {
				log.Printf("Could not send audio: %v", err)
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("EOF.")
			break
		}
		if err != nil {
			log.Fatalf("Cannot stream results: %v", err)
		}
		if err := resp.Error; err != nil {
			log.Fatalf("Could not recognize: %v", err)
		}
		for _, result := range resp.Results {
			for _, alt := range result.GetAlternatives() {
				log.Println(alt.GetTranscript())
			}
		}
	}
}
