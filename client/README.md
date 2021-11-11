# Websocket client for testing


## Usage

```
Usage of /tmp/go-build3015340734/b001/exe/websocket-client:
  -d    display debugging log
  -s string
        sessionid to join
  -u string
        username to use
  -url string
        websocket server url. default is ws://localhost:8080/ws
```

### Run client example

`$ go run . -s chatroom1 -u user1`

`$ go run . -s chatroom1 -u user1 -url ws://plantumliveserver.herokuapp.com/ws`