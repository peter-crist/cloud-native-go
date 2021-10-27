build-client:
	go build -o bin/client ./client

serve:
	go run main.go

helloworld:
	bin/client send "Hello World"

generate:
	protoc -I . ./proto/chat.proto --go_out=plugins=grpc:.
