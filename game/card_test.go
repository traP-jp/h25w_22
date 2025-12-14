package game

import "testing"

func TestGetCardByID(t *testing.T) {
	tests := []struct {
		name   string
		id     int
		wantID int
		isNil  bool
	}{
		{"Valid card ID 1", 1, 1, false},
		{"Valid card ID 5", 5, 5, false},
		{"Invalid card ID", 999, 0, true},
		{"Zero card ID", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := GetCardByID(tt.id)
			if tt.isNil {
				if card != nil {
					t.Errorf("GetCardByID(%d) = %v, want nil", tt.id, card)
				}
			} else {
				if card == nil {
					t.Errorf("GetCardByID(%d) = nil, want card", tt.id)
				} else if card.ID != tt.wantID {
					t.Errorf("GetCardByID(%d).ID = %d, want %d", tt.id, card.ID, tt.wantID)
				}
			}
		})
	}
}

func TestCardMasterData(t *testing.T) {
	if len(cardMaster) == 0 {
		t.Error("cardMaster should not be empty")
	}

	for _, card := range cardMaster {
		if card.ID <= 0 {
			t.Errorf("Card %s has invalid ID: %d", card.Name, card.ID)
		}
		if card.Name == "" {
			t.Errorf("Card %d has empty name", card.ID)
		}
		if card.Power < 0 {
			t.Errorf("Card %s has negative power: %d", card.Name, card.Power)
		}
		if card.HitRate < 1 || card.HitRate > 100 {
			t.Errorf("Card %s has invalid hit rate: %d (should be 1-100)", card.Name, card.HitRate)
		}
		if card.Type != "ATTACK" && card.Type != "DEFENSE" {
			t.Errorf("Card %s has invalid type: %s", card.Name, card.Type)
		}
	}
}
