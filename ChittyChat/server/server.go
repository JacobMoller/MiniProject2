package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"Miniproject2/ChittyChat/protobuf"

	"google.golang.org/grpc"
)

var activityDB []Activity
var participants []string
var receivedBroadcasts [300]string
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
	var responseMessage = ""
	if in.Type == "JOIN" {
		participants = append(participants, in.From)
		responseMessage = "Participant " + in.From + " joined Chitty-Chat at Lamport time " + strconv.Itoa(int(lamport))
	} else {
		responseMessage = "(Lamport time " + strconv.Itoa(int(lamport)) + ") " + in.From + ": " + in.Message
	}
	activityDB = append(activityDB, Activity{id, in.Time, in.Type, responseMessage, in.From})
	lamport++
	addToRecievedBroadcast(in.From)
	return &protobuf.PublishReply{}, nil
}

func (s *server) Broadcast(ctx context.Context, in *protobuf.BroadcastRequest) (*protobuf.BroadcastReply, error) {
	//Compare lamport here
	addToRecievedBroadcast(in.From)
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

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //error while listening
		log.Fatalf("failed to serve: %v", err)
	}
	time.Sleep(1000 * time.Second)
}

var arrayCount = 0
var checkForDeadCount = 0

func addToRecievedBroadcast(username string) {
	if arrayCount >= len(participants)*3 {
		arrayCount = 0
	}
	receivedBroadcasts[arrayCount] = username
	checkForDeadCount++
	if checkForDeadCount >= len(participants)*3 {
		CheckIfOneOrMoreIsDead()
		checkForDeadCount = 0
	}
	arrayCount++
	fmt.Print("[")
	for i := 0; i < len(participants)*3; i++ {
		fmt.Print(receivedBroadcasts[i] + ",")
	}
	fmt.Println("]")
}

var participantsAlive []string

func CheckIfOneOrMoreIsDead() {
	participantsAlive = nil
	for i := 0; i < len(participants); i++ {
		var appeared bool = false
		for j := 0; j < len(participants)*3; j++ {
			if participants[i] == receivedBroadcasts[j] {
				appeared = true
			}
		}
		if appeared {
			participantsAlive = append(participantsAlive, participants[i])
		} else {
			lamport++
			var responseMessage = "Participant " + participants[i] + " left Chitty-Chat at Lamport time " + strconv.Itoa(int(lamport))
			activityDB = append(activityDB, Activity{int32(len(activityDB) + 1), lamport, "LEAVE", responseMessage, participants[i]})
		}
	}
	participants = participantsAlive
}
