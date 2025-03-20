build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/hello.proto

	go build -o build/server server/main.go
	go build -o build/client-macos client/main.go
	GOOS=linux GOARCH=amd64 go build -o build/client-linux client/main.go
	GOOS=windows GOARCH=amd64 go build -o build/client-windows.exe client/main.go

run-server:
	./build/server

run-client:
	./build/client-macos --script=./myscript.sh

generate-java-sdk:
	protoc --java_out=./sdk/java proto/hello.proto

generate-python-sdk:
	protoc --python_out=./sdk/python proto/hello.proto