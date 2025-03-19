package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc_poc/proto"
)

type server struct {
	proto.UnimplementedHelloServiceServer
	clients  map[string]bool
	mutex    sync.Mutex
	messages map[string]string
	updates  map[string]time.Time
}

func checkClient(s *server, clientID string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.clients[clientID]; ok {
		return clientID, nil
	} else if clientID == "" {
		newID := fmt.Sprintf("client-%d", rand.Int())
		s.clients[newID] = true
		return newID, nil
	}
	return "", errors.New("Invalid client ID")
}

func (s *server) SayHello(ctx context.Context, req *proto.Message) (*proto.Message, error) {
	id, err := checkClient(s, req.ClientId)
	if err != nil {
		return nil, err
	}
	if req.Code == "PING" {
		return &proto.Message{ClientId: id, Code: "PONG", Message: "Hello, " + id}, nil
	} else {
		return &proto.Message{ClientId: id, Code: "UNKNOWN", Message: "Only understand PING"}, nil
	}
}

func (s *server) SayHelloStream(req *proto.Message, stream proto.HelloService_SayHelloStreamServer) error {
	_ = stream.RecvMsg(&req)
	id, err := checkClient(s, req.ClientId)
	if err != nil {
		return err
	}
	stream.Send(&proto.Message{ClientId: id, Code: "STARTING", Message: "Starting work"})

	startTime := time.Now()
	for {
		elapsed := time.Since(startTime)
		code := "STATUS"
		if elapsed > 30*time.Second {
			code = "COMPLETED"
			stream.Send(&proto.Message{ClientId: id, Code: code, Message: "Please execute script"})
			break
		}
		stream.Send(&proto.Message{ClientId: id, Code: code, Message: "In Progress..."})
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (s *server) BiDiHello(stream proto.HelloService_BiDiHelloServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		id, err := checkClient(s, req.ClientId)
		if err != nil {
			return err
		}
		stream.Send(&proto.Message{ClientId: id, Code: "STARTING", Message: "Starting work"})

		startTime := time.Now()
		for {
			elapsed := time.Since(startTime)
			code := "STATUS"
			if elapsed > 30*time.Second {
				code = "COMPLETED"
				stream.Send(&proto.Message{ClientId: id, Code: code, Message: "Please execute script"})
				break
			}
			stream.Send(&proto.Message{ClientId: id, Code: code, Message: "In Progress..."})
			time.Sleep(5 * time.Second)
		}
	}
}

func (s *server) PollHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Polling request received...")
	clientID := r.URL.Query().Get("client_id")
	id, err := checkClient(s, clientID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s.mutex.Lock()
	updateTime, timeExists := s.updates[id]
	s.mutex.Unlock()

	resp := map[string]interface{}{
		"client_id": id,
		"code":      "STATUS",
		"message":   "In Progress...",
	}

	if !timeExists {
		resp = map[string]interface{}{
			"client_id": id,
			"code":      "STARTING",
			"message":   "Starting work",
		}
		updateTime = time.Now()
		s.updates[id] = updateTime
	}

	elapsed := time.Since(updateTime)

	if elapsed > 30*time.Second {
		resp = map[string]interface{}{
			"client_id": id,
			"code":      "COMPLETED",
			"message":   "Please execute script",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := &server{
		clients:  make(map[string]bool),
		updates:  make(map[string]time.Time),
		messages: make(map[string]string),
	}
	grpcServer := grpc.NewServer()
	proto.RegisterHelloServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	// Start REST Polling Server
	http.HandleFunc("/api/poll", srv.PollHandler)
	go func() {
		log.Println("REST Polling Server is running on port 8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	fmt.Println("gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
