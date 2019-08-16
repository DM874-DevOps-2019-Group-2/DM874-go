package main

//go:generate protoc -I pb --go_out=plugins=grpc:pb pb/pb.proto

import (
	"context"
	"fmt"
	"time"

	"github.com/ValdemarGr/DM874-go/pb"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
)

const topic = "testtopic3"

func produce() {
	time.Sleep(time.Second * 2)

	data := &pb.Message{
		SenderId: 69,
		Message:  "This is a cool message\nAnd Bye!",
	}

	as_ba, _ := proto.Marshal(data)

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	fmt.Print("Writing\n")

	_ = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("HEY"),
			Value: as_ba,
		},
	)
}

func main() {

	go produce()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       topic,
		Partition:   0,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     time.Millisecond * 100,
		StartOffset: kafka.LastOffset,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s\n", m.Offset, string(m.Key))

		newData := &pb.Message{}
		_ = proto.Unmarshal(m.Value, newData)

		fmt.Print(newData.Message)

		fmt.Print(newData.SenderId)
	}

	ce := r.Close()
	fmt.Print(ce)
}
