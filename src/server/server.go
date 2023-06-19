package main

import (
	data "github.com/finallly/streaming-test/src/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"reflect"
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

		log.Printf("received message number: %d, message size: %d", message.GetId(), getSize(message.GetMessage()))
		//time.Sleep(time.Millisecond * 300) - тут эмулируем чтение сервером сообщений
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

func getSize(v interface{}) int {
	size := int(reflect.TypeOf(v).Size())
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			size += getSize(s.Index(i).Interface())
		}
	case reflect.Map:
		s := reflect.ValueOf(v)
		keys := s.MapKeys()
		size += int(float64(len(keys)) * 10.79)
		for i := range keys {
			size += getSize(keys[i].Interface()) + getSize(s.MapIndex(keys[i]).Interface())
		}
	case reflect.String:
		size += reflect.ValueOf(v).Len()
	case reflect.Struct:
		s := reflect.ValueOf(v)
		for i := 0; i < s.NumField(); i++ {
			if s.Field(i).CanInterface() {
				size += getSize(s.Field(i).Interface())
			}
		}
	}
	return size
}
