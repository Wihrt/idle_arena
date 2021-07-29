package gladiator

import "github.com/wihrt/idle_arena/dice"

type Caracteristic struct {
	Name     string `json:"name"`
	Value    int    `json:"value"`
	Modifier int    `json:"modifier"`
}

func NewCaracteristic(name string, number int, dices int, keep_highest int) *Caracteristic {
	var (
		c = &Caracteristic{
			Name:  name,
			Value: 0,
		}
	)
	c.Add(dice.Roll(number, dices, keep_highest))
	return c
}

func (c *Caracteristic) Add(number int) {
	c.Value += number
	c.Modifier = calculateModifier(c.Value)
}

func calculateModifier(value int) int {
	var modifier = float64((value - 10) / 2)
	return int(modifier)
}
