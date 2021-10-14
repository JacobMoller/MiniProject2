package main

import (
	"context"
	"log"
	"net"

	"Miniproject2/ChittyChat/protobuf"

	"google.golang.org/grpc"
)

var activityDB []Activity

type Activity struct {
	Time    string
	Type    string
	Message string
	From    string
	New     bool
}

type server struct {
	protobuf.UnimplementedChittyChatServer
}

//the type of data we can handle
type Message struct {
	Name string `json:"name"`
}

func (s *server) Publish(ctx context.Context, in *protobuf.PublishRequest) (*protobuf.PublishReply, error) {
	activityDB = append(activityDB, Activity{in.Time, in.Type, in.Message, in.From, true})
	return &protobuf.PublishReply{}, nil
}

func (s *server) Broadcast(ctx context.Context, in *protobuf.BroadcastRequest) (*protobuf.BroadcastReply, error) {
	var newActivities []*protobuf.Activity
	for i := 0; i < len(activityDB); i++ {
		activity := activityDB[i]
		activityForProtobuf := &protobuf.Activity{
			Time:    activity.Time,
			Type:    activity.Type,
			Message: activity.Message,
			From:    activity.From,
		}
		if activity.New {
			activityDB[i].New = false
			newActivities = append(newActivities, activityForProtobuf)
		}
	}
	return &protobuf.BroadcastReply{Activities: newActivities}, nil
}

func main() {
	//port :8080
	lis, err := net.Listen("tcp", ":8080")

	if err != nil { //error before listening
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //we create a new server
	protobuf.RegisterChittyChatServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //error while listening
		log.Fatalf("failed to serve: %v", err)
	}
}

/*

	var listOfChatMessages = []*protobuf.Course{
		{From: "Jacob",Type:"message",Content:"Hej verden!"},
		{From: "Jeppe",Type:"message",Content:"Yo!"},
		{From: "Freja",Type:"message",Content:"Japan!"},
	}

*/
