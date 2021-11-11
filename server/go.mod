module plantumlive-ws-server

go 1.17

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	google.golang.org/protobuf v1.27.1
)

require golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect

replace plantumlive-ws-server/wsmsg => ./wsmsg
