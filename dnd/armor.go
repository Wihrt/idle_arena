package dnd

type Armor struct {
	Index              string     `json:"index" bson:"index"`
	Name               string     `json:"name" bson:"name"`
	ArmorCategory      string     `json:"armor_category" bson:"armor_category"`
	ArmorClass         ArmorClass `json:"armor_class" bson:"armor_class"`
	StrMinimum         int        `json:"str_minimum" bson:"str_minimum"`
	StealthDisavantage bool       `json:"stealth_disavantage" bson:"stealth_disavantage"`
	Weight             int        `json:"weight" bson:"weight"`
	Cost               Cost       `json:"cost" bson:"cost"`
	URL                string     `json:"url" bson:"url"`
}

type ArmorClass struct {
	Base     int  `json:"base" bson:"base"`
	DexBonus bool `json:"dex_bonus" bson:"dex_bonus"`
	MaxBonus int  `json:"max_bonus" bson:"max_bonus"`
}
