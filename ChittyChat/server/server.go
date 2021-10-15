package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"Miniproject2/ChittyChat/protobuf"

	"google.golang.org/grpc"
)

var activityDB []Activity
var participants []string
var participantsAlive []string
var lamport = int32(0)

type Activity struct {
	Id      int32
	Time    int32
	Type    string
	Message string
	From    string
}

type server struct {
	protobuf.UnimplementedChittyChatServer
}

//the type of data we can handle
type Message struct {
	Name string `json:"name"`
}

func (s *server) Publish(ctx context.Context, in *protobuf.PublishRequest) (*protobuf.PublishReply, error) {
	//Compare lamport here
	if in.Time > lamport {
		lamport = in.Time
	}
	lamport++
	var id = int32(len(activityDB) + 1)
	activityDB = append(activityDB, Activity{id, in.Time, in.Type, in.Message, in.From})
	fmt.Println(in.Type)
	if in.Type == "JOIN" {
		participants = append(participants, in.From)
	}
	lamport++
	return &protobuf.PublishReply{}, nil
}

func (s *server) Broadcast(ctx context.Context, in *protobuf.BroadcastRequest) (*protobuf.BroadcastReply, error) {
	//Compare lamport here
	if in.Time > lamport {
		lamport = in.Time
	}
	lamport++
	var newActivities []*protobuf.Activity
	for i := 0; i < len(activityDB); i++ {
		activity := activityDB[i]
		activityForProtobuf := &protobuf.Activity{
			Id:      activity.Id,
			Time:    lamport,
			Type:    activity.Type,
			Message: activity.Message,
			From:    activity.From,
		}
		if in.LatestMessageId < activityDB[i].Id {
			newActivities = append(newActivities, activityForProtobuf)
		}
	}
	lamport++
	return &protobuf.BroadcastReply{Time: lamport, Activities: newActivities}, nil
}

func main() {
	//port :8080
	lis, err := net.Listen("tcp", ":8080")

	if err != nil { //error before listening
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //we create a new server
	protobuf.RegisterChittyChatServer(s, &server{})
	go test()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //error while listening
		log.Fatalf("failed to serve: %v", err)
	}
	time.Sleep(1000 * time.Second)
}

func test() {
	for {
		fmt.Println("du er sød")
		time.Sleep(time.Second)
	}
}

/*

	var listOfChatMessages = []*protobuf.Course{
		{From: "Jacob",Type:"message",Content:"Hej verden!"},
		{From: "Jeppe",Type:"message",Content:"Yo!"},
		{From: "Freja",Type:"message",Content:"Japan!"},
	}

*/
