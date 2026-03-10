package ws

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ximofam/chess-online/services/auth"
)

type AutoIncID struct {
	counter uint
	mu      sync.Mutex
}

func (ai *AutoIncID) GetID() uint {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	ai.counter++
	return ai.counter
}

type server struct {
	clients map[uint]*Client

	lobby *Lobby

	gameID AutoIncID
	rooms  map[uint]*ChessRoom

	upgrader websocket.Upgrader

	mu sync.RWMutex
}

var Server = &server{
	clients: make(map[uint]*Client),
	lobby: &Lobby{
		players: make(map[*Player]struct{}),
	},
	rooms: make(map[uint]*ChessRoom),
	upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	},
}

func (s *server) register(c *Client) {
	s.broadcastToLobby(WsMessageResponse{
		Type: LobbyTypeUserOnline,
		Data: c.player,
	})

	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[c.player.ID] = c
}

func (s *server) remove(remove *Client) {
	s.mu.Lock()
	if c, ok := s.clients[remove.player.ID]; ok && c == remove {
		delete(s.clients, remove.player.ID)
	}
	s.mu.Unlock()

	s.leaveRoom(remove)
	s.lobby.remove(remove.player)

	s.broadcastToLobby(WsMessageResponse{
		Type: LobbyTypeUserLeave,
		Data: remove.player,
	})
}

func (s *server) broadcastToClients(playerIDs []uint, msg any) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, id := range playerIDs {
		if c, ok := s.clients[id]; ok {
			c.write <- msg
		}
	}
}

func (s *server) broadcastToLobby(msg any) {
	userIDs := s.lobby.getPlayerIDs()
	s.broadcastToClients(userIDs, msg)
}

func (s *server) broadcastToRoom(id uint, msg any) {
	s.mu.RLock()
	room, ok := s.rooms[id]
	s.mu.RUnlock()
	if !ok {
		return
	}

	userIDs := room.getPlayerIDs()
	s.broadcastToClients(userIDs, msg)
}

func (s *server) createRoom(name string, allowSpectate bool, maxSpectators int) *ChessRoom {
	gameID := s.gameID.GetID()

	s.mu.Lock()
	defer s.mu.Unlock()

	game := NewChessRoom(gameID, name, allowSpectate, maxSpectators)

	s.rooms[gameID] = game

	return game
}

func (s *server) getRoom(id uint) *ChessRoom {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.rooms[id]
}

func (s *server) deleteRoom(id uint) {
	s.mu.Lock()
	room, ok := s.rooms[id]
	if ok {
		if room.gameState != nil {
			room.gameState.Close()
		}

		playerIDs := room.getPlayerIDs()
		for _, id := range playerIDs {
			if c, ok := s.clients[id]; ok {
				s.lobby.register(c.player)
				c.SetHandleFunc(s.LobbyHandleFunc)
			}
		}
	}

	delete(s.rooms, id)
	s.mu.Unlock()

	s.broadcastToLobby(WsMessageResponse{
		Type: LobbyTypeRoomDelete,
		Data: map[string]any{
			"id": id,
		},
	})
}

func (s *server) leaveRoom(c *Client) {
	player := c.player
	InRoom := player.InRoom
	gameID := c.player.InRoom.roomID
	if gameID == 0 {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	game := s.getRoom(gameID)
	if game == nil {
		return
	}

	if InRoom.role == PlayerRoleWhite {
		s.broadcastToRoom(gameID, WsMessageResponse{
			Type: ChessTypeRoomDelete,
		})
		s.deleteRoom(gameID)
		return
	}

	game.remove(c.player)
	s.broadcastToRoom(gameID, WsMessageResponse{
		Type: ChessTypePlayerLeave,
		Data: map[string]any{
			"id":       player.ID,
			"username": player.Username,
			"role":     InRoom.role,
		},
	})
}

func (s *server) changeToChessHandleFunc(c *Client) {
	s.lobby.remove(c.player)
	c.SetHandleFunc(s.ChessHandleFunc)
}

func (s *server) changeToLobbyHandleFunc(c *Client) {
	s.lobby.register(c.player)
	c.SetHandleFunc(s.LobbyHandleFunc)
}

func (s *server) playChess(roomID uint) error {
	room := s.getRoom(roomID)
	if room == nil {
		return errors.New("room not found")
	}

	gameState := NewGameState(30 * time.Second)
	room.gameState = gameState

	go func() {
		s.broadcastToRoom(roomID, WsMessageResponse{
			Type: ChessTypeGameState,
			Data: toGameStateResponse(gameState),
		})

		gameover, ok := <-gameState.GameOver
		if !ok {
			return
		}

		s.broadcastToRoom(roomID, WsMessageResponse{
			Type: ChessTypeGameover,
			Data: GameResult{
				Result: gameover.String(),
				Method: gameState.Game.Method().String(),
			},
		})

		time.AfterFunc(10*time.Second, func() {
			s.deleteRoom(roomID)
		})
	}()

	return nil
}

func (s *server) ServeWS(c *gin.Context) {
	user := auth.GetUser(c)

	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(
		&Player{
			ID:       user.ID,
			Username: user.Username},
		conn,
		s.LobbyHandleFunc,
		func(c *Client) { s.remove(c) })

	s.register(client)
	s.lobby.register(client.player)

	go client.readPump()
	go client.writePump()
}
