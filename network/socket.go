package network

import "golang.org/x/net/websocket"

type Message struct {
	Name    string
	Message string
	Time    int64
}

type Room struct {
	Forward chan *Message // 수신하는 메시지를 보관하는 값
	// 들어오는 메세지를 다른 클라이언트에게 전송합니다.

	Join  chan *Client
	Leave chan *Client

	Clients map[*Client]bool
}

type Client struct {
	Send   chan Message
	Room   *Room
	Name   string
	Socket websocket.Conn
}
