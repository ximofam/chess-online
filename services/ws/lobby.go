package ws

import "sync"

type Lobby struct {
	players map[*Player]struct{}
	mu      sync.RWMutex
}

func (l *Lobby) register(u *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.players[u] = struct{}{}
}

func (l *Lobby) remove(u *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.players, u)
}

func (l *Lobby) getPlayerIDs() []uint {
	l.mu.RLock()
	defer l.mu.RUnlock()

	PlayerIDs := make([]uint, 0, len(l.players))
	for u := range l.players {
		PlayerIDs = append(PlayerIDs, u.ID)
	}

	return PlayerIDs
}
