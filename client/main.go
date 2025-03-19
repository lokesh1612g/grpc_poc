package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_poc/proto"
	"io"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	ClientID string `json:"client_id"`
}

const configFile = "client_config.json"

func loadClientID() string {
	file, err := os.Open(configFile)
	if err == nil {
		defer file.Close()
		var config Config
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err == nil {
			return config.ClientID
		}
	}
	return ""
}

func saveClientID(clientID string) {
	file, err := os.Create(configFile)
	if err == nil {
		defer file.Close()
		config := Config{ClientID: clientID}
		encoder := json.NewEncoder(file)
		err = encoder.Encode(config)
	}
}

func executeScript(scriptPath string, message string) {
	log.Printf("Message: %s\n", message)
	cmd := exec.Command(scriptPath, message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Script execution failed: %v", err, output)
	}
	fmt.Printf("%s\n", output)
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)
	clientID := loadClientID()

	scriptPath := flag.String("script", "", "Path to script to execute on response")
	flag.Parse()

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()

	stream, err := client.SayHelloStream(context.Background(), &proto.Message{ClientId: clientID, Code: "START"})
	if err != nil {
		log.Fatalf("Error initiating stream: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove(configFile)
			log.Fatalf("Stream error: %v", err)
		}
		saveClientID(resp.ClientId)
		fmt.Println("ClientID:", resp.ClientId, "Code:", resp.Code, "Received from server:", resp.Message)
		if resp.Code == "COMPLETED" && *scriptPath != "" {
			executeScript(*scriptPath, resp.ClientId)
		}
	}
}
