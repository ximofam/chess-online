package ws

import (
	"encoding/json"
	"log"
)

const (
	LobbyEventInfoList   = "INFO_LIST"
	LobbyEventRoomJoin   = "ROOM_JOIN"
	LobbyEventRoomCreate = "ROOM_CREATE"
)

const (
	LobbyTypeInfoList       = "INFO_LIST"
	LobbyTypeRoomCreate     = "ROOM_CREATE"
	LobbyTypeRoomDelete     = "ROOM_DELETE"
	LobbyTypeRoomUpdate     = "ROOM_UPDATE"
	LobbyTypeRoomJoinSucces = "ROOM_JOIN_SUCCESS"
	LobbyTypeUserOnline     = "USER_ONLINE"
	LobbyTypeUserLeave      = "USER_LEAVE"
)

func (s *server) LobbyHandleFunc(c *Client, msg *WsMessageRequest) {
	switch msg.Event {
	case LobbyEventInfoList:
		s.mu.RLock()
		onlineUsers := make([]Player, 0, len(s.clients))
		for _, c := range s.clients {
			onlineUsers = append(onlineUsers, *c.player)
		}
		rooms := make([]ChessRoomResponse, 0, len(s.rooms))
		for _, room := range s.rooms {
			rooms = append(rooms, toChessRoomResponse(room))
		}
		s.mu.RUnlock()

		c.write <- WsMessageResponse{
			Type: LobbyTypeInfoList,
			Data: LobbyInfo{
				OnlineUsers: onlineUsers,
				Rooms:       rooms,
			},
		}

	case LobbyEventRoomCreate:
		var payload struct {
			Name          string `json:"name"`
			AllowSpectate bool   `json:"allow_spectate"`
			MaxSpectators int    `json:"max_spectators"`
		}

		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Printf("HandleMsgFunc: %v", err)
			return
		}

		game := s.createRoom(
			payload.Name,
			payload.AllowSpectate,
			payload.MaxSpectators,
		)

		if err := game.join(c.player, PlayerRoleWhite); err != nil {
			return
		}
		s.changeToChessHandleFunc(c)
		c.write <- WsMessageResponse{
			Type: LobbyTypeRoomJoinSucces,
			Data: PlayerRoleWhite,
		}

		s.broadcastToLobby(WsMessageResponse{
			Type: LobbyTypeRoomCreate,
			Data: toChessRoomResponse(game),
		})
	case LobbyEventRoomJoin:
		var payload struct {
			GameID uint   `json:"room_id"`
			Role   string `json:"role"`
		}

		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Printf("HandleMsgFunc: %v", err)
			return
		}

		game := s.getRoom(payload.GameID)
		if game == nil {
			return
		}

		if err := game.join(c.player, payload.Role); err != nil {
			return
		}
		s.changeToChessHandleFunc(c)
		c.write <- WsMessageResponse{
			Type: LobbyTypeRoomJoinSucces,
			Data: payload.Role,
		}

		notifyNewPlayerJoinToRoom(s, c.player)
		s.broadcastToLobby(WsMessageResponse{
			Type: LobbyTypeRoomUpdate,
			Data: toChessRoomResponse(game),
		})
	}
}
