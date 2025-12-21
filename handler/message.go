package handler

import "encoding/json"

// Event represents a server-to-client event
type Event struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload,omitempty"`
}

// MatchedPayload is sent when both players have joined
type MatchedPayload struct{}

// GameStartPayload is sent when the game starts
type GameStartPayload struct{}

// TurnStartPayload contains turn information
type TurnStartPayload struct {
	TurnCount int    `json:"turnCount"`
	Role      string `json:"role"`
}

// ActionStartPayload is sent when the battle phase begins
type ActionStartPayload struct {
	MaxHP int `json:"maxHp"`
}

// OpponentPlayedPayload notifies about opponent's card play
type OpponentPlayedPayload struct {
	CardID int `json:"cardId"`
}

// ActionResultPayload contains the result of a card exchange
type ActionResultPayload struct {
	Hit       bool `json:"hit"`
	Damage    int  `json:"damage"`
	CurrentHP int  `json:"currentHp"`
}

// TurnEndPayload is sent when a turn ends
type TurnEndPayload struct {
	Reason    string `json:"reason"` // "HP_ZERO" or "LIMIT_REACHED"
	TurnScore int    `json:"turnScore"`
}

// GameResultPayload contains final game results
type GameResultPayload struct {
	FinalScores map[string]int `json:"finalScores"`
}

// OpponentDisconnectedPayload notifies that the opponent has disconnected
type OpponentDisconnectedPayload struct {
	Reason string `json:"reason"` // Descriptive message (e.g., "Player disconnected")
}

// MarshalEvent converts an event to JSON bytes
func MarshalEvent(eventName string, payload interface{}) ([]byte, error) {
	event := Event{
		Event:   eventName,
		Payload: payload,
	}
	return json.Marshal(event)
}
