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
	// Attack Cards (妨害カード)
	{ID: 1, Name: "騒ぎ立てる", Power: 30, HitRate: 100, Type: "ATTACK"},
	{ID: 2, Name: "ゴミをポイ捨て", Power: 30, HitRate: 100, Type: "ATTACK"},
	{ID: 3, Name: "香水をばら撒く", Power: 30, HitRate: 60, Type: "ATTACK"},
	{ID: 4, Name: "藁人形をばら撒く", Power: 40, HitRate: 90, Type: "ATTACK"},
	{ID: 5, Name: "SNSに悪口を書く", Power: 40, HitRate: 100, Type: "ATTACK"},
	{ID: 6, Name: "すれ違い様に足をかける", Power: 50, HitRate: 100, Type: "ATTACK"},
	{ID: 7, Name: "雨乞いをする", Power: 50, HitRate: 50, Type: "ATTACK"},
	{ID: 8, Name: "ナンパする", Power: 50, HitRate: 70, Type: "ATTACK"},
	{ID: 9, Name: "テニスボールを投げつける", Power: 60, HitRate: 80, Type: "ATTACK"},
	{ID: 10, Name: "Gを投げつける", Power: 60, HitRate: 70, Type: "ATTACK"},
	{ID: 11, Name: "ひったくりをする", Power: 70, HitRate: 70, Type: "ATTACK"},
	{ID: 12, Name: "物乞いを雇って向かわせる", Power: 80, HitRate: 70, Type: "ATTACK"},
	{ID: 13, Name: "マジで思いつかないため保留", Power: 90, HitRate: 60, Type: "ATTACK"},
	{ID: 14, Name: "保留", Power: 100, HitRate: 50, Type: "ATTACK"},
	{ID: 15, Name: "保留", Power: 120, HitRate: 40, Type: "ATTACK"},
	{ID: 16, Name: "保留", Power: 150, HitRate: 30, Type: "ATTACK"},
	{ID: 17, Name: "保留", Power: 150, HitRate: 20, Type: "ATTACK"},
	{ID: 18, Name: "保留", Power: 200, HitRate: 20, Type: "ATTACK"},
	{ID: 19, Name: "保留", Power: 300, HitRate: 10, Type: "ATTACK"},
	// Defense Cards (妨害阻止カード)
	{ID: 20, Name: "服装を褒める", Power: 10, HitRate: 100, Type: "DEFENSE"},
	{ID: 21, Name: "瞳を見つめる", Power: 10, HitRate: 100, Type: "DEFENSE"},
	{ID: 22, Name: "手を繋ぐ", Power: 30, HitRate: 100, Type: "DEFENSE"},
	{ID: 23, Name: "抱きしめる", Power: 50, HitRate: 100, Type: "DEFENSE"},
	{ID: 24, Name: "保留", Power: 50, HitRate: 100, Type: "DEFENSE"},
	{ID: 25, Name: "保留", Power: 80, HitRate: 100, Type: "DEFENSE"},
	{ID: 26, Name: "保留", Power: 100, HitRate: 100, Type: "DEFENSE"},
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
