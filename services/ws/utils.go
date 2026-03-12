package ws

import (
	"fmt"
	"time"
)

func notifyNewPlayerJoinToRoom(s *server, player *Player) {
	roomID := player.InRoom.roomID
	if roomID == 0 {
		return
	}

	chatMsg := ChatMessage{
		Text:   fmt.Sprintf("%s has joined the room", player.Username),
		SendAt: time.Now(),
	}

	room := s.getRoom(roomID)
	if room == nil {
		return
	}

	room.historyChatMsg = append(room.historyChatMsg, chatMsg)

	s.broadcastToRoom(roomID, WsMessageResponse{
		Type: ChessTypeChat,
		Data: chatMsg,
	})

	s.broadcastToRoom(roomID, WsMessageResponse{
		Type: ChessTypePlayerJoin,
		Data: map[string]any{
			"id":       player.ID,
			"username": player.Username,
			"role":     player.InRoom.role,
		},
	})
}
