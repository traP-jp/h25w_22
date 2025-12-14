package game

import "testing"

func TestNewGameState(t *testing.T) {
	gs := NewGameState()

	if gs.TurnCount != 0 {
		t.Errorf("TurnCount = %d, want 0", gs.TurnCount)
	}
	if gs.Phase != "MATCHING" {
		t.Errorf("Phase = %s, want MATCHING", gs.Phase)
	}
	if gs.Roles == nil {
		t.Error("Roles should not be nil")
	}
	if gs.Scores == nil {
		t.Error("Scores should not be nil")
	}
	if gs.ReadyPlayers == nil {
		t.Error("ReadyPlayers should not be nil")
	}
}

func TestSwapRoles(t *testing.T) {
	gs := NewGameState()
	gs.Roles["client1"] = RoleAttack
	gs.Roles["client2"] = RoleDefense

	gs.SwapRoles()

	if gs.Roles["client1"] != RoleDefense {
		t.Errorf("client1 role = %s, want DEFENSE", gs.Roles["client1"])
	}
	if gs.Roles["client2"] != RoleAttack {
		t.Errorf("client2 role = %s, want ATTACK", gs.Roles["client2"])
	}
}

func TestResetTurn(t *testing.T) {
	gs := NewGameState()
	gs.CardsPlayed = 3
	gs.HP = 100
	gs.ScoreMultiplier = 1.5
	gs.DateSelected = true
	gs.ReadyPlayers["client1"] = true

	gs.ResetTurn()

	if gs.CardsPlayed != 0 {
		t.Errorf("CardsPlayed = %d, want 0", gs.CardsPlayed)
	}
	if gs.HP != 0 {
		t.Errorf("HP = %d, want 0", gs.HP)
	}
	if gs.ScoreMultiplier != 0 {
		t.Errorf("ScoreMultiplier = %f, want 0", gs.ScoreMultiplier)
	}
	if gs.DateSelected {
		t.Error("DateSelected should be false")
	}
	if len(gs.ReadyPlayers) != 0 {
		t.Errorf("ReadyPlayers length = %d, want 0", len(gs.ReadyPlayers))
	}
}

func TestGetRoleForClient(t *testing.T) {
	gs := NewGameState()
	gs.Roles["client1"] = RoleAttack
	gs.Roles["client2"] = RoleDefense

	tests := []struct {
		clientID string
		want     string
	}{
		{"client1", RoleAttack},
		{"client2", RoleDefense},
		{"client3", ""},
	}

	for _, tt := range tests {
		t.Run(tt.clientID, func(t *testing.T) {
			got := gs.GetRoleForClient(tt.clientID)
			if got != tt.want {
				t.Errorf("GetRoleForClient(%s) = %s, want %s", tt.clientID, got, tt.want)
			}
		})
	}
}

func TestGetClientByRole(t *testing.T) {
	gs := NewGameState()
	gs.Roles["client1"] = RoleAttack
	gs.Roles["client2"] = RoleDefense

	attacker := gs.GetClientByRole(RoleAttack)
	if attacker != "client1" {
		t.Errorf("GetClientByRole(ATTACK) = %s, want client1", attacker)
	}

	defender := gs.GetClientByRole(RoleDefense)
	if defender != "client2" {
		t.Errorf("GetClientByRole(DEFENSE) = %s, want client2", defender)
	}
}

func TestAllPlayersReady(t *testing.T) {
	gs := NewGameState()

	// No players ready
	if gs.AllPlayersReady() {
		t.Error("AllPlayersReady() = true, want false when no players are ready")
	}

	// One player ready
	gs.ReadyPlayers["client1"] = true
	if gs.AllPlayersReady() {
		t.Error("AllPlayersReady() = true, want false when only one player is ready")
	}

	// Two players ready
	gs.ReadyPlayers["client2"] = true
	if !gs.AllPlayersReady() {
		t.Error("AllPlayersReady() = false, want true when both players are ready")
	}

	// One player not ready
	gs.ReadyPlayers["client2"] = false
	if gs.AllPlayersReady() {
		t.Error("AllPlayersReady() = true, want false when one player is not ready")
	}
}
