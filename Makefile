build-client:
	go build -o bin/client ./client

serve:
	go run main.go

send:
	bin/client send "Hello World"

.PHONY: circuitbreaker debounce retry throttle timeout fanin fanout
circuitbreaker: 
		bin/client cb

debounce:
	bin/client debounce

retry:
	bin/client retry

throttle:
	bin/client throttle

timeout:
	bin/client timeout

fanin:
	bin/client fanin

fanout:
	bin/client fanout

generate:
	protoc -I . ./proto/chat.proto --go_out=plugins=grpc:.
