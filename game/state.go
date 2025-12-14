package game

const (
	// RoleAttack represents the attacking role (Jammer)
	RoleAttack = "ATTACK"
	// RoleDefense represents the defending role (Dater)
	RoleDefense = "DEFENSE"
)

// Phase represents the current phase of the game
type Phase int

const (
	// PhaseMatching represents the matching phase
	PhaseMatching Phase = iota
	// PhaseTurnStart represents the turn start phase
	PhaseTurnStart
	// PhaseAction represents the action/battle phase
	PhaseAction
	// PhaseTurnEnd represents the turn end phase
	PhaseTurnEnd
	// PhaseGameEnd represents the game end phase
	PhaseGameEnd
)

// GameState represents the current state of the game
type GameState struct {
	TurnCount       int               // Current turn (1-4)
	Phase           Phase             // Current phase
	Roles           map[string]string // Maps client ID to role (ATTACK or DEFENSE)
	HP              int               // Current HP of the defender
	ScoreMultiplier float64           // Multiplier for score calculation
	CardsPlayed     int               // Number of cards played by attacker this turn
	PendingAttack   *Card             // Attack card waiting for defense response
	Scores          map[string]int    // Scores for each role (ATTACK/DEFENSE)
	ReadyPlayers    map[string]bool   // Tracks which players are ready
	DateSelected    bool              // Whether defender has selected date location
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	return &GameState{
		TurnCount:    0,
		Phase:        PhaseMatching,
		Roles:        make(map[string]string),
		Scores:       map[string]int{RoleAttack: 0, RoleDefense: 0},
		ReadyPlayers: make(map[string]bool),
		DateSelected: false,
	}
}

// SwapRoles swaps the roles between players
func (gs *GameState) SwapRoles() {
	for clientID, role := range gs.Roles {
		if role == RoleAttack {
			gs.Roles[clientID] = RoleDefense
		} else {
			gs.Roles[clientID] = RoleAttack
		}
	}
}

// ResetTurn resets turn-specific state
func (gs *GameState) ResetTurn() {
	gs.CardsPlayed = 0
	gs.PendingAttack = nil
	gs.HP = 0
	gs.ScoreMultiplier = 0
	gs.ReadyPlayers = make(map[string]bool)
	gs.DateSelected = false
}

// GetRoleForClient returns the role for a given client ID
func (gs *GameState) GetRoleForClient(clientID string) string {
	return gs.Roles[clientID]
}

// GetClientByRole returns the client ID for a given role
func (gs *GameState) GetClientByRole(role string) string {
	for clientID, r := range gs.Roles {
		if r == role {
			return clientID
		}
	}
	return ""
}

const maxPlayers = 2

// AllPlayersReady checks if all players are ready
func (gs *GameState) AllPlayersReady() bool {
	if len(gs.ReadyPlayers) != maxPlayers {
		return false
	}
	for _, ready := range gs.ReadyPlayers {
		if !ready {
			return false
		}
	}
	return true
}
