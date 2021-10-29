module websocket-client

go 1.17

require (
	github.com/gorilla/websocket v1.4.2
	google.golang.org/protobuf v1.27.1
)

replace websocket-client/wsmsg => ./wsmsg
