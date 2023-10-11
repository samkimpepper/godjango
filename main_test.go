package main

import (
	w "awesomeProject/websocket"
	fwebsocket "github.com/fasthttp/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestWebsocketServer(t *testing.T) {
	app := setup()
	defer app.Shutdown()

	conn, res, err := fwebsocket.DefaultDialer.Dial("ws://localhost:3000/ws/test", nil)
	defer conn.Close()
	assert.NoError(t, err)
	assert.Equal(t, 101, res.StatusCode)
	assert.Equal(t, "websocket", res.Header.Get("Upgrade"))

	var msg map[string]interface{}
	err = conn.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, "notify", msg["type"])
}

func setup() *fiber.App {
	handler := websocket.New(func(conn *websocket.Conn) {
		groupName := conn.Params("group")

		channel := w.NewChannel()
		channel.AddGroup(groupName, conn)

		msg := map[string]interface{}{
			"type":    "notify",
			"content": "New user connected",
		}

		channel.Send(conn, msg)
	})

	app := fiber.New()

	app.Get("/ws/:group", handler)

	go app.Listen(":3000")

	readyCh := make(chan struct{})

	go func() {
		for {
			conn, err := net.Dial("tcp", "localhost:3000")
			if err != nil {
				log.Println(err)
				continue
			}

			if conn != nil {
				readyCh <- struct{}{}
				conn.Close()
				break
			}
		}
	}()

	<-readyCh

	return app
}
