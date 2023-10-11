# godjango
Django Channels처럼 쓸 수 있는 Go Fiber Websocket 라이브러리 
godjango is a Go library for creating WebSocket-based real-time communication in your Go applications. It is inspired by Django Channels and provides an easy way to implement WebSocket functionally in your Go projects.

## Installation
To use godjango in your project, simply add it to your Go module using `go get`:
```sh
go get github.com/samkimpepper/godjango
```

# Getting Started
## Example Usage
Here's a simple example of how to use godjango in your Go application:
```go
import (
  w "awesomeProject/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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
```
