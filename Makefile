build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/hello.proto

	go build -o build/server server/main.go
	go build -o build/client client/main.go

run-server:
	./build/server

run-client:
	./build/client --script=./myscript.sh