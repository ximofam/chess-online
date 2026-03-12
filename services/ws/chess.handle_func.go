package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	ChessEventLeave    = "LEAVE"
	ChessEventChat     = "CHAT"
	ChessEventMove     = "MOVE"
	ChessEventReady    = "READY"
	ChessEventInfoList = "INFO_LIST"
	ChessEventPlay     = "PLAY"
)

const (
	ChessTypeChat        = "CHAT"
	ChessTypePlayerLeave = "PLAYER_LEAVE"
	ChessTypePlayerJoin  = "PLAYER_JOIN"
	ChessTypeRoomDelete  = "ROOM_DELETE"
	ChessTypeBoardUpdate = "BOARD_UPDATE"
	ChessTypeInfoList    = "INFO_LIST"
	ChessTypePlay        = "PLAY"
	ChessTypeGameState   = "GAME_STATE"
	ChessTypeGameover    = "GAMEOVER"
)

func (s *server) ChessHandleFunc(c *Client, msg *WsMessageRequest) {
	player := c.player
	InRoom := player.InRoom
	room := s.getRoom(InRoom.roomID)
	if InRoom.roomID == 0 || room == nil {
		s.changeToLobbyHandleFunc(c)
		return
	}

	switch msg.Event {
	case ChessEventLeave:
		s.changeToLobbyHandleFunc(c)

		if InRoom.role == PlayerRoleWhite {
			roomID := InRoom.roomID
			s.broadcastToRoom(roomID, WsMessageResponse{
				Type: ChessTypeRoomDelete,
			})
			s.deleteRoom(roomID)
			return
		}

		room.remove(player)

		chatMsg := ChatMessage{
			Text:   fmt.Sprintf("%s has left the room", player.Username),
			SendAt: time.Now(),
		}

		room.historyChatMsg = append(room.historyChatMsg, chatMsg)

		s.broadcastToRoom(InRoom.roomID, WsMessageResponse{
			Type: ChessTypeChat,
			Data: chatMsg,
		})

		s.broadcastToRoom(InRoom.roomID, WsMessageResponse{
			Type: ChessTypePlayerLeave,
			Data: map[string]any{
				"id":       player.ID,
				"username": player.Username,
				"role":     InRoom.role,
			},
		})

		s.broadcastToLobby(WsMessageResponse{
			Type: LobbyTypeRoomUpdate,
			Data: toChessRoomResponse(room),
		})

	case ChessEventInfoList:
		res := toChessRoomInfo(room)

		c.write <- WsMessageResponse{
			Type: ChessTypeInfoList,
			Data: res,
		}

	case ChessEventChat:
		var payload struct {
			Text string `json:"text"`
		}

		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}

		chatMsg := ChatMessage{
			From:   c.player.Username,
			Text:   payload.Text,
			SendAt: time.Now(),
		}

		room.historyChatMsg = append(room.historyChatMsg, chatMsg)

		s.broadcastToRoom(InRoom.roomID, WsMessageResponse{
			Type: ChessTypeChat,
			Data: chatMsg,
		})
	case ChessEventPlay:
		if InRoom.role != PlayerRoleWhite {
			return
		}

		if room.isPlaying || room.Black == nil {
			return
		}

		if err := s.playChess(InRoom.roomID); err != nil {
			log.Printf("Play chess err: %v", err)
			return
		}

		s.broadcastToRoom(room.ID, WsMessageResponse{
			Type: ChessTypePlay,
			Data: toGameStateResponse(room.gameState),
		})
		room.isPlaying = true
		s.broadcastToLobby(WsMessageResponse{
			Type: LobbyTypeRoomUpdate,
			Data: toChessRoomResponse(room),
		})
	case ChessEventMove:
		gameState := room.gameState
		if InRoom.role != gameState.Turn.Name() {
			return
		}

		var payload struct {
			UCI string `json:"uci"`
		}

		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}

		if err := gameState.Move(payload.UCI); err != nil {
			return
		}

		s.broadcastToRoom(InRoom.roomID, WsMessageResponse{
			Type: ChessTypeGameState,
			Data: toGameStateResponse(gameState),
		})
	}
}
