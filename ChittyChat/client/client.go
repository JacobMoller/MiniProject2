package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"Miniproject2/ChittyChat/protobuf"

	"google.golang.org/grpc"
)

//var int latestTimestamp

var possibleResponses = []string{"hej", "dav", "dejligt vejr", "godnat"}

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

	client.Publish(context.Background(), &protobuf.PublishRequest{Time: "nu", Type: "JOIN", Message: "JOIN MSG", From: name})

	for {
		response3, err2 := client.Broadcast(context.Background(), &protobuf.BroadcastRequest{LatestMessageTimestamp: "hewwo"})
		if err2 != nil {
			log.Fatalf("could not get chat list: %v", err)
		}
		for i := 0; i < len(response3.Activities); i++ {
			log.Print(response3.Activities[i])
			//if response3.Activities[i].Time > latestTimestamp {
			//update the latest timestamp
			//}
		}
		fmt.Println("listening again")
		time.Sleep(time.Second)
		//wait x seconds
	}
}
