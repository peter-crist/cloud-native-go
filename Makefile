build-client:
	go build -o bin/client ./client

serve:
	go run main.go

send:
	bin/client send "Hello World"

.PHONY: circuitbreaker debounce retry throttle
circuitbreaker: 
		bin/client cb

debounce:
	bin/client debounce

retry:
	bin/client retry

throttle:
	bin/client throttle

generate:
	protoc -I . ./proto/chat.proto --go_out=plugins=grpc:.
