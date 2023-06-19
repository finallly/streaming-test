package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	data "github.com/finallly/streaming-test/src/proto"
)

func main() {
	rand.Seed(time.Now().Unix())

	conn, err := grpc.Dial("localhost:50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := data.NewStreamServiceClient(conn)

	stream, err := client.StartStream(context.Background())

	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	index := int32(0)

	message := data.Message{Content: &data.Message_Word{Word: randomSequence(819200)}}

	for {
		err := stream.Send(&data.Stream{
			Id: index,
			Message: []*data.Message{
				&message,
			},
		})
		log.Printf("message with id %d sent", index)

		if err == io.EOF {
			return
		}
		index++
	}
}

func randomSequence(n int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
