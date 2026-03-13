package ws

import (
	"errors"
	"sync"
	"time"

	"github.com/notnil/chess"
)

const (
	PlayerRoleWhite     = "White"
	PlayerRoleBlack     = "Black"
	PlayerRoleSpectator = "Spectator"
)

type InRoom struct {
	roomID uint
	role   string
}

type Player struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	InRoom   InRoom `json:"-"`
}

type ChatMessage struct {
	From   string    `json:"from"`
	Text   string    `json:"text"`
	SendAt time.Time `json:"send_at"`
}

type GameState struct {
	Game *chess.Game

	WhiteTime time.Duration
	BlackTime time.Duration

	Turn chess.Color

	LastMoveTime time.Time

	timer *time.Timer

	mu sync.Mutex

	GameOver chan chess.Outcome

	once sync.Once
}

func NewGameState(initial time.Duration) *GameState {

	g := &GameState{
		Game:         chess.NewGame(),
		WhiteTime:    initial,
		BlackTime:    initial,
		Turn:         chess.White,
		LastMoveTime: time.Now(),
		GameOver:     make(chan chess.Outcome, 1),
	}

	g.startTimer()

	return g
}

func (g *GameState) Close() {
	g.once.Do(func() {
		close(g.GameOver)
		g.timer.Stop()
	})
}

func (g *GameState) startTimer() {

	remaining := g.remainingTime()

	g.timer = time.AfterFunc(remaining, func() {

		g.mu.Lock()
		defer g.mu.Unlock()

		if g.Turn == chess.White {
			g.GameOver <- chess.BlackWon
		} else {
			g.GameOver <- chess.WhiteWon
		}
	})
}

func (g *GameState) remainingTime() time.Duration {

	elapsed := time.Since(g.LastMoveTime)

	if g.Turn == chess.White {
		return g.WhiteTime - elapsed
	}

	return g.BlackTime - elapsed
}

func (g *GameState) Move(uci string) error {

	g.mu.Lock()
	defer g.mu.Unlock()

	move, err := chess.UCINotation{}.Decode(g.Game.Position(), uci)
	if err != nil {
		return err
	}

	elapsed := time.Since(g.LastMoveTime)

	if g.Turn == chess.White {
		g.WhiteTime -= elapsed
	} else {
		g.BlackTime -= elapsed
	}

	if err := g.Game.Move(move); err != nil {
		return err
	}

	if g.Game.Outcome() != chess.NoOutcome {
		g.GameOver <- g.Game.Outcome()
		return nil
	}

	if g.timer != nil {
		g.timer.Stop()
	}

	g.Turn = g.Game.Position().Turn()
	g.LastMoveTime = time.Now()

	g.startTimer()

	return nil
}

func (g *GameState) FEN() string {
	return g.Game.FEN()
}

type ChessRoom struct {
	ID             uint
	Name           string
	White          *Player
	Black          *Player
	joinAble       bool
	AllowSpectate  bool
	spectators     []Player
	gameState      *GameState
	maxSpectators  int
	historyChatMsg []ChatMessage
	isPlaying      bool
	mu             sync.RWMutex
	ClientManager
}

func NewChessRoom(id uint, name string, allowSpectate bool, maxSpectators int) *ChessRoom {
	game := &ChessRoom{
		ID:             id,
		Name:           name,
		joinAble:       true,
		spectators:     make([]Player, 0, maxSpectators),
		maxSpectators:  maxSpectators,
		AllowSpectate:  allowSpectate,
		historyChatMsg: make([]ChatMessage, 0),
		ClientManager: ClientManager{
			clients: make(map[*Client]struct{}),
		},
	}

	return game
}

func (g *ChessRoom) join(c *Client, role string) error {
	g.mu.Lock()

	u := c.player

	u.InRoom = InRoom{
		roomID: g.ID,
		role:   role,
	}

	success := false

	switch role {
	case PlayerRoleWhite:
		if g.White == nil {
			g.White = u
			success = true
		}
	case PlayerRoleBlack:
		if g.Black == nil {
			g.Black = u
			success = true
		}
	case PlayerRoleSpectator:
		if g.AllowSpectate && (g.maxSpectators == 0 || len(g.spectators) < g.maxSpectators) {
			g.spectators = append(g.spectators, *u)
			success = true
		}
	}

	if !success {
		u.InRoom = InRoom{}
		g.mu.Unlock()
		return errors.New("failed to join room")
	}
	g.mu.Unlock()

	g.register(c)

	return nil
}

func (g *ChessRoom) remove(c *Client) {
	g.mu.Lock()

	u := c.player
	u.InRoom = InRoom{}

	if g.White == u {
		g.White = nil
	} else if g.Black == u {
		g.Black = nil
	} else if g.AllowSpectate {
		for i, player := range g.spectators {
			if u.ID == player.ID {
				g.spectators = append(g.spectators[:i], g.spectators[i+1:]...)
				break
			}
		}
	}
	g.mu.Unlock()

	g.ClientManager.remove(c)
}
