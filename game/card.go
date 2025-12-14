package game

// Card represents a card with its properties
type Card struct {
	ID      int
	Name    string
	Power   int    // Damage amount
	HitRate int    // Percentage (1-100)
	Type    string // "ATTACK" or "DEFENSE"
}

// cardMaster contains all available cards
//
//nolint:mnd // Card stats are game design constants
var cardMaster = []Card{
	{ID: 1, Name: "Weak Attack", Power: 10, HitRate: 90, Type: "ATTACK"},
	{ID: 2, Name: "Medium Attack", Power: 20, HitRate: 70, Type: "ATTACK"},
	{ID: 3, Name: "Strong Attack", Power: 30, HitRate: 50, Type: "ATTACK"},
	{ID: 4, Name: "Critical Strike", Power: 40, HitRate: 30, Type: "ATTACK"},
	{ID: 5, Name: "Weak Defense", Power: 5, HitRate: 80, Type: "DEFENSE"},
	{ID: 6, Name: "Medium Defense", Power: 10, HitRate: 60, Type: "DEFENSE"},
	{ID: 7, Name: "Strong Defense", Power: 15, HitRate: 40, Type: "DEFENSE"},
	{ID: 8, Name: "Ultimate Defense", Power: 20, HitRate: 20, Type: "DEFENSE"},
}

// cardMap provides O(1) lookup for cards by ID
var cardMap = buildCardMap()

func buildCardMap() map[int]*Card {
	m := make(map[int]*Card)
	for i := range cardMaster {
		m[cardMaster[i].ID] = &cardMaster[i]
	}
	return m
}

// GetCardByID returns a card by its ID
func GetCardByID(id int) *Card {
	return cardMap[id]
}
