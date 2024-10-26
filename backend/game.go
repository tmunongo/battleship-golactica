package main

import (
	"math/rand"
	"sync"
)

// Game represents
//
//
type Game struct {
	ID          string     `json:"id"`
	Players     [2]*Player `json:"players"`
	CurrentTurn int        `json:"currentTurn"`
	Status      string     `json:"status"`
	mu          sync.Mutex
}

func (g *Game) placeShips(player *Player) {
	// for now, place randomly
	shipsPlaced := 0
	for shipsPlaced < ShipCount {
		row := rand.Intn(BoardSize)
		col := rand.Intn(BoardSize)

		if !player.Board[row][col].HasShip {
			player.Board[row][col].HasShip = true
			shipsPlaced++
		}
	}
}

// broadcastGameState notifies players of the current game Status
//
// 
func (g *Game) broadcastGameState() {
	for _, player := range g.Players {
		if player != nil && player.Conn != nil {
			gameState := Message{
				Type: "gameState",
				Payload: map[string]interface{}{
					"board":       player.Board,
					"currentTurn": g.CurrentTurn,
					"status":      g.Status,
				},
			}

			player.Conn.WriteJSON(gameState)
		}
	}
}

func (g *Game) handleShot(playerIndex, row, col int) bool {
	targetIndex := (playerIndex + 1) % 2
	targetBoard := &g.Players[targetIndex].Board
	cell := &targetBoard[row][col]

	cell.Hit = true
	hit := cell.HasShip

	g.broadcastGameState()

	return hit
}
