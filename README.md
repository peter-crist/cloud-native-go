# cloud-native-go
Scratch repository for implementing cloud native Go patterns

Running the gRPC Server:
`make serve`

Building CLI client to communicate with the gRPC Server:
`make build-client`

Run a Hello World example:
`Ensure gRPC server is running and CLI binary is built`
`make send`

Generate protofile:
`make generate`

## Circuit Breaker
To demonstrate the CircuitBreaker pattern, start the gRPC server with `make serve`.
In another terminal, run `make circuitbreaker`.
If using the CLI tool, you can run `bin/client cb` and pass custom parameters to tweak the demo output.

Observe the log output from the server which showcases the Circuit Breaker in action!

## Debounce
To demonstrate a Debounce pattern, start the gRPC server with `make serve`.
In another terminal, run `make debounce`.
If using the CLI tool, you can run `bin/client debounce` and pass custom parameters to tweak the demo output.

Observe the log output from the server which showcases restricting clusters of requests.