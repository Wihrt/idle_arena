package gladiator

type Health struct {
	Current int `json:"current" bson:"current"`
	Max     int `json:"max" bson:"max"`
}

func NewHealth(cMod int) *Health {
	h := &Health{
		Current: 12 + cMod,
		Max:     12 + cMod,
	}

	return h
}
