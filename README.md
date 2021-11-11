# Websocket server and testing client

## Websocket Server

### Usage

```
  -d    display debugging log
```

### Run client example

`$ go run .` 

`$ go run . -d` // showing logs


## Testing Websocket Client

### Usage

```
  -d    display debugging log
  -s string
        sessionid to join
  -u string
        username to use
  -url string
        websocket server url. default is ws://localhost:8080/ws
```

### Run client example

`$ go run . -s chatroom1 -u user1` // run with default URL

`$ go run . -s chatroom1 -u user1 -url ws://{your-hosting-server-URL}/ws` 

## Reference 

Websocket server is made based on the example code of `Gorilla websocket`

Posting link: https://golangdocs.com/golang-gorilla-websockets

GitHub: https://github.com/gorilla/websocket

chatting example: https://github.com/gorilla/websocket/tree/master/examples/chat
