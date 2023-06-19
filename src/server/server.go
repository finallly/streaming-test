package main

import (
	data "github.com/finallly/streaming-test/src/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type server struct {
	data.UnimplementedStreamServiceServer
}

func (s *server) StartStream(stream data.StreamService_StartStreamServer) error {
	for {
		message, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatal("error receiving message from stream", err.Error())
		}

		log.Printf("received message number: %d", message.GetId())
		time.Sleep(time.Millisecond * 500) //  тут эмулируем чтение сервером сообщений
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:50005")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	s := grpc.NewServer()
	data.RegisterStreamServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
