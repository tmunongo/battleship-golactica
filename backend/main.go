package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	BoardSize        = 10
	ShipCount        = 5
	WaitingForPlayer = "waiting"
	InProgress       = "in_progress"
	Finished         = "finished"
)

type Player struct {
	ID    string                     `json:"id"`
	Board [BoardSize][BoardSize]Cell `json:"board"`
	Conn  *websocket.Conn            `json:"-"`
	Reay  bool                       `json:"ready"`
}

type Cell struct {
	HasShip bool `json:"hasShip"`
	Hit     bool `json:"hit"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: check origin
		return true
	},
}

func main() {
	gameManager := NewGameManager()

	// handle websocket connections
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("Websocket upgrade error: %v", err)
		}

		game := gameManager.createOrJoinGame()
		player := &Player{
			ID:   generatePlayerID(),
			Conn: conn,
		}

		// add player to game
		if game.Players[0] == nil {
			game.Players[0] = player
		} else {
			game.Players[1] = player
			game.Status = InProgress

			game.placeShips(game.Players[0])
			game.placeShips(game.Players[1])

			game.broadcastGameState()
		}

		go handlePlayerConnection(game, player)
	})

	log.Println("Starting Golactica backend server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateGameID() string {
	return "game_" + generateRandomString(8)
}

func generatePlayerID() string {
	return "player_" + generateRandomString(8)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func handlePlayerConnection(game *Game, player *Player) {
	for {
		var msg Message
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("ws connection error: %v", err)
			return
		}

		switch msg.Type {
		case "fire":
			payload := msg.Payload.(map[string]interface{})
			row := int(payload["row"].(float64))
			col := int(payload["col"].(float64))

			playerIndex := 0
			if game.Players[1] == player {
				playerIndex = 1
			}

			if playerIndex == game.CurrentTurn {
				hit := game.handleShot(playerIndex, row, col)
				game.CurrentTurn = (game.CurrentTurn + 1) % 2

				player.Conn.WriteJSON(Message{
					Type: "shotResult",
					Payload: map[string]interface{}{
						"hit": hit,
						"row": row,
						"col": col},
				})
			}
		}
	}
}
