package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"Miniproject2/ChittyChat/protobuf"

	"google.golang.org/grpc"
)

func main() {
	log.Print("Welcome to ChittyChat. Please enter a username:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	name := strings.Replace(text, "\n", "", 1)

	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock()) //maybe it has to be: localhost:8080
	if err != nil {                                                        //error can not establish connection
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := protobuf.NewChittyChatClient(conn)

	response, err := client.Publish(context.Background(), &protobuf.PublishRequest{Time: "nu", Type: "JOIN", Message: "JOIN MSG", From: name})
	if err != nil {
		log.Fatalf("could not get chat list: %v", err)
	}

	response2, err2 := client.Publish(context.Background(), &protobuf.PublishRequest{Time: "nu", Type: "CHAT", Message: "Japan!", From: name})
	if err2 != nil {
		log.Fatalf("could not get chat list: %v", err)
	}

	log.Print(response.GetMessage())
	log.Print(response2.GetMessage())
}
