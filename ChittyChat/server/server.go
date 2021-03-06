package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"Miniproject2/ChittyChat/protobuf"

	"sync"

	"google.golang.org/grpc"
)

var activityDB []Activity
var participants []string
var receivedBroadcasts [300]string
var l Lamport

type Activity struct {
	Id      int32
	Time    int32
	Type    string
	Message string
	From    string
}

type Lamport struct {
	mu      sync.Mutex
	lamport int32
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
	if in.Time > l.getLamport() {
		l.replaceLamport(in.Time)
	}
	l.Inc()
	var id = int32(len(activityDB) + 1)
	var responseMessage = ""
	if in.Type == "JOIN" {
		participants = append(participants, in.From)
		responseMessage = "Participant " + in.From + " joined Chitty-Chat at Lamport time " + strconv.Itoa(int(l.getLamport()))
	} else {
		responseMessage = "(Lamport time " + strconv.Itoa(int(l.getLamport())) + ") " + in.From + ": " + in.Message
	}
	activityDB = append(activityDB, Activity{id, in.Time, in.Type, responseMessage, in.From})
	l.Inc()
	addToRecievedBroadcast(in.From)
	return &protobuf.PublishReply{}, nil
}

func (s *server) Broadcast(ctx context.Context, in *protobuf.BroadcastRequest) (*protobuf.BroadcastReply, error) {
	//Compare lamport here
	addToRecievedBroadcast(in.From)
	if in.Time > l.getLamport() {
		l.replaceLamport(in.Time)
	}
	l.Inc()
	var newActivities []*protobuf.Activity
	for i := 0; i < len(activityDB); i++ {
		activity := activityDB[i]
		activityForProtobuf := &protobuf.Activity{
			Id:      activity.Id,
			Time:    l.getLamport(),
			Type:    activity.Type,
			Message: activity.Message,
			From:    activity.From,
		}
		if in.LatestMessageId < activityDB[i].Id {
			newActivities = append(newActivities, activityForProtobuf)
		}
	}
	l.Inc()
	return &protobuf.BroadcastReply{Time: l.getLamport(), Activities: newActivities}, nil
}

func main() {
	l = Lamport{}
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
			l.Inc()
			var responseMessage = "Participant " + participants[i] + " left Chitty-Chat at Lamport time " + strconv.Itoa(int(l.getLamport()))
			activityDB = append(activityDB, Activity{int32(len(activityDB) + 1), l.getLamport(), "LEAVE", responseMessage, participants[i]})
		}
	}
	participants = participantsAlive
}

func (l *Lamport) Inc() {
	l.mu.Lock()
	l.lamport++
	l.mu.Unlock()
}

func (l *Lamport) getLamport() int32 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.lamport
}
func (l *Lamport) replaceLamport(newLamport int32) {
	l.mu.Lock()
	l.lamport = newLamport
	defer l.mu.Unlock()
}
