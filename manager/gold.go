package manager

type MoneyPouch struct {
	CopperPieces int `json:"copper_pieces" bson:"copper_pieces"`
	SilverPieces int `json:"silver_pieces" bson:"silver_pieces" `
	GoldPieces   int `json:"gold_pieces" bson:"gold_pieces"`
	TotalPieces  int `json:"total_pieces" bson:"total_pieces"`
}

func NewMoneyPouch(difficulty Difficulty) *MoneyPouch {

	var startingCopper = 100

	switch difficulty {
	case DifficultyEasy:
		startingCopper = 25
	case DifficultyNormal:
		startingCopper = 50
	case DifficultyHard:
		startingCopper = 100
	}

	g := &MoneyPouch{
		TotalPieces: startingCopper,
	}
	return g
}

func (m *MoneyPouch) ConvertPieces() {
	// Reset the values
	m.CopperPieces = m.TotalPieces
	m.SilverPieces = 0
	m.GoldPieces = 0
	m.convertCopperPieces()
	m.convertSilverPieces()
}

func (m *MoneyPouch) convertCopperPieces() {
	for {
		if m.CopperPieces < 10 {
			break
		} else {
			m.SilverPieces += 1
			m.CopperPieces -= 10
		}
	}
}

func (m *MoneyPouch) convertSilverPieces() {
	for {
		if m.SilverPieces < 10 {
			break
		} else {
			m.GoldPieces += 1
			m.SilverPieces -= 10
		}
	}
}
