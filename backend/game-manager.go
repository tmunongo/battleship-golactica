package main

import "sync"

type GameManager struct {
	games map[string]*Game
	mu    sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*Game),
	}
}

func (gm *GameManager) createOrJoinGame() *Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// look for game waiting for player
	for _, game := range gm.games {
		if game.Status == WaitingForPlayer {
			return game
		}
	}

	gameId := generateGameID()

	game := &Game{
		ID:     gameId,
		Status: WaitingForPlayer,
	}
	gm.games[gameId] = game

	return game
}
