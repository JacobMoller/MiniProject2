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

var latestMessageId = int32(0)
var lamport = int32(0)

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

	//Publish Join
	lamport++
	client.Publish(context.Background(), &protobuf.PublishRequest{Time: lamport, Type: "JOIN", Message: "JOIN MSG", From: name})

	go collectNewActivities(client, name)
	go enterToChat(client, name)
	time.Sleep(1000 * time.Second)
}

func enterToChat(client protobuf.ChittyChatClient, name string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		//Publish chat msg
		lamport++
		client.Publish(context.Background(), &protobuf.PublishRequest{Time: lamport, Type: "CHAT", Message: text, From: name})
	}
}

func collectNewActivities(client protobuf.ChittyChatClient, name string) {
	for {
		//Broadcast
		lamport++
		response, err := client.Broadcast(context.Background(), &protobuf.BroadcastRequest{Time: lamport, LatestMessageId: latestMessageId, From: name})
		if err != nil {
			log.Fatalf("could not get chat list: %v", err)
		}
		if response != nil {
			//Compare lamport here
			if response.Time > lamport {
				lamport = response.Time
			}
			lamport++
			for i := 0; i < len(response.Activities); i++ {
				log.Print("(" + fmt.Sprint(lamport) + ") " + response.Activities[i].From + ": " + response.Activities[i].Message)
				if response.Activities[i].Id > latestMessageId {
					latestMessageId = response.Activities[i].Id
				}
			}
		}
		time.Sleep(time.Second)
	}
}
