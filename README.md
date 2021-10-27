# cloud-native-go
Scratch repository for implementing cloud native Go patterns

Running the gRPC Server:
`make serve`

Building CLI client to communicate with the gRPC Server:
`make build-client`

Run a Hello World example:
`Ensure gRPC server is running and CLI binary is built`
`make helloworld`

Generate protofile:
`make generate`

# ToDo
- Serve via Docker
- Potentially run serve as CLI command so it can be customized for various pattern examples
- Start implementing patterns