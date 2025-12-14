package game

import "math/rand"

const (
	// BaseScoreConstant is the base score added when defender survives a turn
	BaseScoreConstant = 50
	// MaxCardsPerTurn is the maximum number of cards attacker can play per turn
	MaxCardsPerTurn = 4
)

// CalculateHit determines if an attack hits based on the card's hit rate
func CalculateHit(attackCard *Card) bool {
	return (rand.Intn(100) + 1) <= attackCard.HitRate
}

// CalculateDamage returns the damage amount if hit, 0 otherwise
func CalculateDamage(attackCard *Card, hit bool) int {
	if hit {
		return attackCard.Power
	}
	return 0
}

// CalculateScore calculates the score bonus for surviving a turn
func CalculateScore(currentHP int, scoreMultiplier float64) int {
	return BaseScoreConstant + int(float64(currentHP)*scoreMultiplier)
}
