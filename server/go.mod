module plantumlive-ws-server

go 1.17

require (
	github.com/daixiang0/gci v0.2.9
	github.com/golang/mock v1.6.0
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	mvdan.cc/gofumpt v0.1.1
)

replace plantumlive-ws-server/wsmsg => ./wsmsg
