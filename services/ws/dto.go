package ws

import (
	"encoding/json"
	"time"
)

type WsMessageRequest struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type WsMessageResponse struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type ChessRoomResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	White         string `json:"white"`
	Black         string `json:"black"`
	JoinAble      bool   `json:"join_able"`
	AllowSpectate bool   `json:"allow_spectate"`
	MaxSpectators int    `json:"max_spectators"`
	CurrObserver  int    `json:"curr_observer"`
	IsPlaying     bool   `json:"is_playing"`
}

type ChessRoomInfo struct {
	White        *Player           `json:"white"`
	Black        *Player           `json:"black"`
	Spectators   []Player          `json:"spectators"`
	ChatMessages []ChatMessage     `json:"chat_messages"`
	IsPlaying    bool              `json:"is_playing"`
	GameState    GameStateResponse `json:"game_state"`
}

func toChessRoomInfo(room *ChessRoom) ChessRoomInfo {
	var gameState GameStateResponse
	if room.isPlaying {
		gameState = toGameStateResponse(room.gameState)
	}

	return ChessRoomInfo{
		White:        room.White,
		Black:        room.Black,
		Spectators:   room.spectators,
		ChatMessages: room.historyChatMsg,
		IsPlaying:    room.isPlaying,
		GameState:    gameState,
	}
}

type LobbyInfo struct {
	OnlineUsers []Player            `json:"online_users"`
	Rooms       []ChessRoomResponse `json:"rooms"`
}

type GameStateResponse struct {
	FEN        string        `json:"fen"`
	Turn       string        `json:"turn"`
	WhiteTime  time.Duration `json:"white_time"`
	BlackTime  time.Duration `json:"black_time"`
	ValidMoves []string      `json:"valid_moves"`
}

func toGameStateResponse(gameState *GameState) GameStateResponse {
	validMoves := gameState.Game.ValidMoves()
	validMovesStr := make([]string, len(validMoves))
	for i, move := range validMoves {
		validMovesStr[i] = move.String()
	}

	return GameStateResponse{
		FEN:        gameState.FEN(),
		Turn:       gameState.Turn.Name(),
		WhiteTime:  gameState.WhiteTime,
		BlackTime:  gameState.BlackTime,
		ValidMoves: validMovesStr,
	}
}

type GameResult struct {
	Result string `json:"result"`
	Method string `json:"method"`
}

func toChessRoomResponse(room *ChessRoom) ChessRoomResponse {
	White := ""
	Black := ""
	joinAble := true

	if room.White != nil {
		White = room.White.Username
	}
	if room.Black != nil {
		Black = room.Black.Username
		joinAble = false
	}

	return ChessRoomResponse{
		ID:            room.ID,
		Name:          room.Name,
		White:         White,
		Black:         Black,
		JoinAble:      joinAble,
		AllowSpectate: room.AllowSpectate,
		MaxSpectators: room.maxSpectators,
		CurrObserver:  len(room.spectators),
		IsPlaying:     room.isPlaying,
	}
}
