package network

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type Network struct {
	engine *gin.Engine
}

func NewServer() *Network {
	n := &Network{
		engine: gin.New(),
	}

	n.engine.Use(gin.Logger())
	n.engine.Use(gin.Recovery())
	n.engine.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	})) // Enable CORS

	pprof.Register(n.engine)
	r := NewRoom()
	go r.RunInit()

	n.engine.GET("/room", r.SocketServer)

	return n
}

func (n *Network) StartServer() {
	log.Println("Starting server on port 8080")
	n.engine.Run(":8080")

}
