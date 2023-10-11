package websocket

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"log"
	"sync"
)

type Channel struct {
	group map[string]map[*websocket.Conn]struct{}
	mu    sync.Mutex
}

func NewChannel() *Channel {
	return &Channel{
		group: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (c *Channel) AddGroup(groupName string, conn *websocket.Conn) {
	c.mu.Lock()
	if c.group[groupName] == nil {
		c.group[groupName] = make(map[*websocket.Conn]struct{})
	}
	c.group[groupName][conn] = struct{}{}
	c.mu.Unlock()
}

func (c *Channel) RemoveGroup(groupName string, conn *websocket.Conn) {
	c.mu.Lock()
	delete(c.group[groupName], conn)
	c.mu.Unlock()
}

func (c *Channel) Broadcast(groupName string, data interface{}) {
	c.mu.Lock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.mu.Unlock()
		log.Printf("Error when marshal data: %v", err)
		return
	}

	for conn := range c.group[groupName] {
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Printf("Error when broadcast data: %v", err)
			return
		}
	}
	c.mu.Unlock()
}

func (c *Channel) Send(conn *websocket.Conn, data map[string]interface{}) {
	c.mu.Lock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.mu.Unlock()
		log.Printf("Error when marshal data: %v", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, jsonData)
	c.mu.Unlock()
}

func (c *Channel) Receive(conn *websocket.Conn) (map[string]interface{}, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		return nil, err
	}

	return data, nil
}
