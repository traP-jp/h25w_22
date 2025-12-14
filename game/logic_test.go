package game

import (
	"testing"
)

func TestCalculateDamage(t *testing.T) {
	card := &Card{Power: 20}

	tests := []struct {
		name       string
		hit        bool
		wantDamage int
	}{
		{"Hit", true, 20},
		{"Miss", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			damage := CalculateDamage(card, tt.hit)
			if damage != tt.wantDamage {
				t.Errorf("CalculateDamage(%v, %v) = %d, want %d", card, tt.hit, damage, tt.wantDamage)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name            string
		hp              int
		multiplier      float64
		expectedMinimum int
	}{
		{"Base case", 100, 1.0, BaseScoreConstant + 100},
		{"High multiplier", 50, 2.0, BaseScoreConstant + 100},
		{"Low HP", 10, 1.5, BaseScoreConstant + 15},
		{"Zero HP", 0, 1.0, BaseScoreConstant},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateScore(tt.hp, tt.multiplier)
			if score < tt.expectedMinimum {
				t.Errorf("CalculateScore(%d, %f) = %d, want >= %d", tt.hp, tt.multiplier, score, tt.expectedMinimum)
			}
		})
	}
}

func TestCalculateHit(t *testing.T) {
	// Test with 100% hit rate
	card100 := &Card{HitRate: 100}
	for i := 0; i < 10; i++ {
		if !CalculateHit(card100) {
			t.Error("Card with 100% hit rate should always hit")
		}
	}

	// Test with 0% hit rate (impossible to hit)
	card0 := &Card{HitRate: 0}
	for i := 0; i < 10; i++ {
		if CalculateHit(card0) {
			t.Error("Card with 0% hit rate should never hit")
		}
	}
}

func TestConstants(t *testing.T) {
	if BaseScoreConstant != 50 {
		t.Errorf("BaseScoreConstant = %d, want 50", BaseScoreConstant)
	}
	if MaxCardsPerTurn != 4 {
		t.Errorf("MaxCardsPerTurn = %d, want 4", MaxCardsPerTurn)
	}
}
