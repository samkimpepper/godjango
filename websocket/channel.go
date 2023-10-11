package websocket

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
	"sync"
)

type Channel struct {
	group map[string]map[*websocket.Conn]struct{}
	users map[*websocket.Conn]string
	mu    sync.Mutex
}

func NewChannel() *Channel {
	return &Channel{
		group: make(map[string]map[*websocket.Conn]struct{}),
		users: make(map[*websocket.Conn]string),
	}
}

func (c *Channel) AddGroup(groupName string, conn *websocket.Conn, userID string) {
	c.mu.Lock()
	if c.group[groupName] == nil {
		c.group[groupName] = make(map[*websocket.Conn]struct{})
	}
	c.group[groupName][conn] = struct{}{}
	if userID != "" {
		c.users[conn] = userID
	}
	c.mu.Unlock()
}

func (c *Channel) RemoveGroup(groupName string, conn *websocket.Conn) {
	c.mu.Lock()
	delete(c.group[groupName], conn)
	delete(c.users, conn)
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

func (c *Channel) Send(conn *websocket.Conn, data interface{}) {
	c.mu.Lock()

	var (
		result []byte
		err    error
	)
	switch d := data.(type) {
	case map[string]interface{}:
		result, err = json.Marshal(d)
		if err != nil {
			c.mu.Unlock()
			log.Printf("Error when marshal data: %v", err)
			return
		}
	case fiber.Map:
		dataMap := make(map[string]interface{})
		for k, v := range d {
			dataMap[k] = v
		}

		result, err = json.Marshal(dataMap)
		if err != nil {
			c.mu.Unlock()
			log.Printf("Error when marshal data: %v", err)
			return
		}
	}

	conn.WriteMessage(websocket.TextMessage, result)
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

func (c *Channel) GetUsers() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	userIDs := make([]string, 0, len(c.users))
	for _, userID := range c.users {
		userIDs = append(userIDs, userID)
	}

	return userIDs
}
