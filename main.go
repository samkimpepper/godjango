package main

import (
	w "awesomeProject/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()

	app.Get("/ws/:group", websocket.New(func(conn *websocket.Conn) {
		groupName := conn.Params("group")

		channel := w.NewChannel()
		channel.AddGroup(groupName, conn)
		log.Printf("New user connected to group %s", groupName)

		for {
			var (
				data map[string]interface{}
				err  error
			)
			if data, err = channel.Receive(conn); err != nil {
				log.Printf("Error when receive data: %v", err)
				break
			}
			log.Printf("Receive data: %v", data)
			channel.Broadcast(groupName, data)
		}

		msg := map[string]interface{}{
			"type":       "notify",
			"content":    "New user connected",
			"todoItemID": -1,
		}
		channel.Broadcast(groupName, msg)

		channel.RemoveGroup(groupName, conn)
	}))

	app.Listen(":3000")
}
