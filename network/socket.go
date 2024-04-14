package network

import (
	"golang-chat-server/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{ReadBufferSize: types.SocketBufferSize, WriteBufferSize: types.MessageBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}

type message struct {
	Name    string
	Message string
	Time    int64
}

type Room struct {
	Forward chan *message // 수신하는 메시지를 보관하는 값
	// 들어오는 메세지를 다른 클라이언트에게 전송합니다.

	Join  chan *Client // Socket이 연결되면 클라이언트를 추가합니다.
	Leave chan *Client // Socket이 연결이 끊기면 클라이언트를 제거합니다.

	Clients map[*Client]bool
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (c *Client) Read() {
	// 클라이언트가 메시지를 읽는 함수
	defer func() { c.Socket.Close() }()
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			panic(err)
		} else {
			msg.Time = time.Now().Unix()
			msg.Name = c.Name

			c.Room.Forward <- msg

		}
	}
}

func (c *Client) Write() {
	// 클라이언트가 메시지를 쓰는 함수
	defer func() { c.Socket.Close() }()
	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) RunInit() {
	// Room에 있는 모든 채널을 관리하는 함수
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			delete(r.Clients, client)
			close(client.Send)
		case msg := <-r.Forward:
			for client := range r.Clients {
				select {
				case client.Send <- msg:
				default:
					delete(r.Clients, client)
					close(client.Send)
				}
			}
		}
	}
}

func (r *Room) SocketServer(c *gin.Context) {

	socket, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	Client := &Client{
		Socket: socket,
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- Client

	defer func() { r.Leave <- Client }()

	go Client.Write()

	Client.Read()
}
