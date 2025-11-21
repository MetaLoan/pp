package handler

import (
	"log"
	"net/http"

	"github.com/MetaLoan/pp/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket:", err)
		return
	}

	// Register connection to Market Engine
	core.GlobalMarket.Register(conn)

	// Keep connection alive and handle incoming messages (if any)
	// For now, we just read and discard to detect disconnection
	go func() {
		defer func() {
			core.GlobalMarket.Unregister(conn)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
