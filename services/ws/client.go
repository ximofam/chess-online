package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type HandleMsgFunc func(c *Client, msg *WsMessageRequest)

type Client struct {
	player        *Player
	conn          *websocket.Conn
	write         chan []byte
	handleMsgFunc HandleMsgFunc
	onClose       func(*Client)
	mu            sync.Mutex
	once          sync.Once
}

func NewClient(player *Player, conn *websocket.Conn, handleMsgFunc HandleMsgFunc, onClose func(*Client)) *Client {
	return &Client{
		player:        player,
		conn:          conn,
		write:         make(chan []byte, 256),
		handleMsgFunc: handleMsgFunc,
		onClose:       onClose,
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		c.onClose(c)
		c.conn.Close()
		close(c.write)
	})
}

func (c *Client) SetHandleFunc(handleMsgFunc HandleMsgFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handleMsgFunc = handleMsgFunc
}

func (c *Client) Send(msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Sprintf("Failed to json.Marshal: %v", err))
	}

	select {
	case c.write <- data:
	default:
		c.Close()
	}
}

func (c *Client) readPump() {
	defer c.Close()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket read message: %v", err)
			}
			return
		}

		var msg WsMessageRequest
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Printf("readPump failed to unmarshal: %v", err)
			continue
		}

		c.handleMsgFunc(c, &msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Close()
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.write:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Websocket write message: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type ClientManager struct {
	clients map[*Client]struct{}
	mu      sync.RWMutex
}

func (cm *ClientManager) register(c *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.clients[c] = struct{}{}
}

func (cm *ClientManager) remove(c *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.clients, c)
}

func (cm *ClientManager) broadcast(msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Sprintf("Failed to json.Marshal: %v", err))
	}

	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for c := range cm.clients {
		select {
		case c.write <- data:
		default:
			c.Close()
		}
	}
}
