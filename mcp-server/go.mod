module github.com/example/mcp-server

go 1.23.0

toolchain go1.24.2

require (
	github.com/lucasnoah/rescue-titan-proto v0.0.0
	google.golang.org/grpc v1.73.0
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

// Use local proto files
replace github.com/lucasnoah/rescue-titan-proto => ../../proto/main
