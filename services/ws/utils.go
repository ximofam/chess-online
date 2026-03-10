package ws

func notifyNewPlayerJoinToRoom(s *server, player *Player) {
	roomID := player.InRoom.roomID
	if roomID == 0 {
		return
	}

	s.broadcastToRoom(roomID, WsMessageResponse{
		Type: ChessTypePlayerJoin,
		Data: map[string]any{
			"id":       player.ID,
			"username": player.Username,
			"role":     player.InRoom.role,
		},
	})
}
