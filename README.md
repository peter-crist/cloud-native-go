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

## Retry
To demonstrate a Retry pattern, start the gRPC server with `make serve`.
In another terminal, run `make retry`.
If using the CLI tool, you can run `bin/client retry` and pass custom parameters to tweak the demo output.

Observe the log output from the server which showcases retrying given transient errors.

## Throttle
To demonstrate a Throttle pattern, start the gRPC server with `make serve`.
In another terminal, run `make throttle`.
If using the CLI tool, you can run `bin/client throttle` and pass custom parameters to tweak the demo output.

Observe the log output from the server which showcases throttling excessive requests.

## Timeout
To demonstrate a Timeout pattern, start the gRPC server with `make serve`.
In another terminal, run `make timeout`.
If using the CLI tool, you can run `bin/client timeout` and pass custom parameters to tweak the demo output.

Observe the log output from the server and the returned responses to observe that we can successfully wrap a timeout
around a given slow function, even if that function is out of our control and doesn't accept a context.