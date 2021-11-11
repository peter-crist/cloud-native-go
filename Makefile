build-client:
	go build -o bin/client ./client

serve:
	go run main.go

send:
	bin/client send "Hello World"

circuitbreaker:
	bin/client cb

debounce:
	bin/client debounce

generate:
	protoc -I . ./proto/chat.proto --go_out=plugins=grpc:.
