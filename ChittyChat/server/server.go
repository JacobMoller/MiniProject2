package main

import (
	"context"
	"log"
	"net"

	"Miniproject2/ChittyChat/protobuf"

	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedChittyChatServer
}

//the type of data we can handle
type Message struct {
	Name string `json:"name"`
}

func (s *server) Publish(ctx context.Context, in *protobuf.PublishRequest) (*protobuf.PublishReply, error) {
	if in.Type == "JOIN" {
		return &protobuf.PublishReply{Message: "Participant " + in.From + " joined Chitty-Chat at Lamport time L"}, nil
	} else {
		return &protobuf.PublishReply{Message: in.From + ": " + in.Message}, nil
	}
	//implement actual logic of getting the courses
	//.HelloReply{Message: fmt.Sprint(courses[num].Name) + ";" + fmt.Sprint(courses[num].StudentSatisfactionRating)}, nil
}

/*func (s *server) Broadcast(ctx context.Context) (*protobuf.GetCourseListReply, error) {
	//implement actual logic of getting the courses

	var listOfChatMessages = []*protobuf.Course{
		{From: "Jacob", Type: "message", Content: "Hej verden!"},
		{From: "Jeppe", Type: "message", Content: "Yo!"},
		{From: "Freja", Type: "message", Content: "Japan!"},
	}
	return &protobuf.broadcastReply{listOfChatMessages}, nil //.HelloReply{Message: fmt.Sprint(courses[num].Name) + ";" + fmt.Sprint(courses[num].StudentSatisfactionRating)}, nil
}*/

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
