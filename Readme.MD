# gRPC POC

This project demonstrates a gRPC server and client implementation in Go. The server provides a `HelloService` with unary, server streaming, and bidirectional streaming RPC methods. The client interacts with the server and can execute a script based on the server's response.

## Prerequisites

- Go 1.16 or later
- Protocol Buffers Compiler (`protoc`)
- gRPC and Protocol Buffers plugins for Go

## Installation

1. **Install Protocol Buffers Compiler**:
    ```sh
    brew install protobuf
    ```

2. **Install the Go gRPC and protobuf plugins**:
    ```sh
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

## Generating Code

Generate Go code from the `.proto` file:
```sh
protoc --go_out=. --go-grpc_out=. proto/hello.proto
```

## Building the Project

Use the `Makefile` to build the server and client:
```sh
make build
```

## Running the Server

Start the gRPC server:
```sh
make run-server
```

## Running the Client

Run the client with the script to execute:
```sh
make run-client
```

## Docker Setup

Build and run the services using Docker Compose:
```sh
docker-compose up --build
```

## Generating SDKs

Generate Java and Python SDKs from the `.proto` file:
```sh
make generate-java-sdk
make generate-python-sdk
```

## Project Structure

- `proto/hello.proto`: Protocol Buffers definition file.
- `server/main.go`: gRPC server implementation.
- `client/main.go`: gRPC client implementation.
- `Makefile`: Build and run commands.
- `docker-compose.yml`: Docker Compose configuration.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.