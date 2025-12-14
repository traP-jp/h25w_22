package handler

import (
	"log"
	"strconv"
	"sync"

	"github.com/traP-jp/h25w_22/game"
)

// Room represents a game room with two players
type Room struct {
	ID         string
	Clients    map[string]*Client
	Game       *game.GameState
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

// NewRoom creates a new room
func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Clients:    make(map[string]*Client),
		Game:       game.NewGameState(),
		Broadcast:  make(chan []byte, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the room's main loop
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.mu.Lock()
			r.Clients[client.ID] = client

			// Assign roles when both players join
			if len(r.Clients) == 2 {
				clientIDs := make([]string, 0, 2)
				for id := range r.Clients {
					clientIDs = append(clientIDs, id)
				}
				r.Game.Roles[clientIDs[0]] = game.RoleAttack
				r.Game.Roles[clientIDs[1]] = game.RoleDefense

				// Send MATCHED event
				msg, _ := MarshalEvent("MATCHED", MatchedPayload{})
				r.broadcastToAll(msg)
			}
			r.mu.Unlock()

		case client := <-r.Unregister:
			r.mu.Lock()
			if _, ok := r.Clients[client.ID]; ok {
				delete(r.Clients, client.ID)
				close(client.Send)
			}
			r.mu.Unlock()

		case message := <-r.Broadcast:
			r.mu.Lock()
			for _, client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client.ID)
				}
			}
			r.mu.Unlock()
		}
	}
}

// HandleMessage processes commands from clients
func (r *Room) HandleMessage(client *Client, cmd string, args []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch cmd {
	case "SELECT_DATE":
		r.handleSelectDate(client, args)
	case "READY":
		r.handleReady(client)
	case "PLAY_CARD":
		r.handlePlayCard(client, args)
	default:
		log.Printf("Unknown command: %s", cmd)
	}
}

func (r *Room) handleSelectDate(client *Client, args []string) {
	if len(args) < 2 {
		return
	}

	role := r.Game.GetRoleForClient(client.ID)
	if role != game.RoleDefense {
		return
	}

	hp, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	multiplier, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return
	}

	r.Game.HP = hp
	r.Game.ScoreMultiplier = multiplier
	r.Game.DateSelected = true

	// Check if we can proceed to ACTION phase
	r.checkTurnStartComplete()
}

func (r *Room) handleReady(client *Client) {
	r.Game.ReadyPlayers[client.ID] = true

	switch r.Game.Phase {
	case "MATCHING":
		if r.Game.AllPlayersReady() {
			r.startGame()
		}
	case "TURN_START":
		r.checkTurnStartComplete()
	case "TURN_END":
		if r.Game.AllPlayersReady() {
			r.advanceToNextTurn()
		}
	}
}

func (r *Room) handlePlayCard(client *Client, args []string) {
	if len(args) < 1 {
		return
	}

	cardID, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	card := game.GetCardByID(cardID)
	if card == nil {
		return
	}

	role := r.Game.GetRoleForClient(client.ID)

	if role == game.RoleAttack {
		// Attacker plays a card
		r.Game.PendingAttack = card
		r.Game.CardsPlayed++

		// Notify defender
		msg, _ := MarshalEvent("OPPONENT_PLAYED", OpponentPlayedPayload{CardID: cardID})
		defenderID := r.Game.GetClientByRole(game.RoleDefense)
		if defender, ok := r.Clients[defenderID]; ok {
			defender.Send <- msg
		}

	} else if role == game.RoleDefense {
		// Defender responds
		if r.Game.PendingAttack == nil {
			return
		}

		// Calculate hit and damage
		hit := game.CalculateHit(r.Game.PendingAttack)
		damage := game.CalculateDamage(r.Game.PendingAttack, hit)
		r.Game.HP -= damage

		// Send result to both players
		msg, _ := MarshalEvent("ACTION_RESULT", ActionResultPayload{
			Hit:       hit,
			Damage:    damage,
			CurrentHP: r.Game.HP,
		})
		r.broadcastToAll(msg)

		// Clear pending attack
		r.Game.PendingAttack = nil

		// Check end conditions
		if r.Game.HP <= 0 {
			r.endTurn("HP_ZERO", 0)
		} else if r.Game.CardsPlayed >= game.MaxCardsPerTurn {
			// Calculate score
			turnScore := game.CalculateScore(r.Game.HP, r.Game.ScoreMultiplier)
			r.Game.Scores[game.RoleDefense] += turnScore
			r.endTurn("LIMIT_REACHED", turnScore)
		}
	}
}

func (r *Room) startGame() {
	r.Game.Phase = "TURN_START"
	r.Game.TurnCount = 1
	r.Game.ResetTurn()

	// Send GAME_START event
	msg, _ := MarshalEvent("GAME_START", GameStartPayload{})
	r.broadcastToAll(msg)

	// Start first turn
	r.startTurn()
}

func (r *Room) startTurn() {
	r.Game.Phase = "TURN_START"
	r.Game.ResetTurn()

	// Send TURN_START to each player with their role
	for clientID, client := range r.Clients {
		role := r.Game.GetRoleForClient(clientID)
		msg, _ := MarshalEvent("TURN_START", TurnStartPayload{
			TurnCount: r.Game.TurnCount,
			Role:      role,
		})
		client.Send <- msg
	}
}

func (r *Room) checkTurnStartComplete() {
	if r.Game.Phase != "TURN_START" {
		return
	}

	// Need both players ready and date selected
	if r.Game.AllPlayersReady() && r.Game.DateSelected {
		r.startAction()
	}
}

func (r *Room) startAction() {
	r.Game.Phase = "ACTION"
	r.Game.ReadyPlayers = make(map[string]bool)

	msg, _ := MarshalEvent("ACTION_START", ActionStartPayload{MaxHP: r.Game.HP})
	r.broadcastToAll(msg)
}

func (r *Room) endTurn(reason string, turnScore int) {
	r.Game.Phase = "TURN_END"

	msg, _ := MarshalEvent("TURN_END", TurnEndPayload{
		Reason:    reason,
		TurnScore: turnScore,
	})
	r.broadcastToAll(msg)
}

func (r *Room) advanceToNextTurn() {
	if r.Game.TurnCount >= 4 {
		r.endGame()
		return
	}

	r.Game.TurnCount++
	r.Game.SwapRoles()
	r.startTurn()
}

func (r *Room) endGame() {
	r.Game.Phase = "GAME_END"

	msg, _ := MarshalEvent("GAME_RESULT", GameResultPayload{
		FinalScores: r.Game.Scores,
	})
	r.broadcastToAll(msg)

	// Close all connections after a short delay
	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		for _, client := range r.Clients {
			close(client.Send)
		}
	}()
}

// IsFull checks if the room has reached maximum capacity
func (r *Room) IsFull() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.Clients) >= 2
}

func (r *Room) broadcastToAll(message []byte) {
	for _, client := range r.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(r.Clients, client.ID)
		}
	}
}
