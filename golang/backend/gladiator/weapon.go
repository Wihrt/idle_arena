package gladiator

import "github.com/wihrt/idle_arena/arena/dice"

type WeaponType string

const (
	MeleeWeapon  WeaponType = "melee"
	RangedWeapon WeaponType = "ranged"
)

type Weapon struct {
	Name   string     `json:"name"`
	Number int        `json:"number"`
	Damage int        `json:"damage"`
	Type   WeaponType `json:"type"`
}

func NewWeapon(name string, number int, damage int, weaponType WeaponType) *Weapon {
	var w = &Weapon{
		Name:   name,
		Number: number,
		Damage: damage,
		Type:   weaponType,
	}

	return w
}

func NewRandomWeapon() *Weapon {
	var weapon *Weapon
	roll := dice.Roll(1, 4, -1)

	switch roll {
	case 1:
		weapon = NewWeapon("Shortsword", 1, 6, MeleeWeapon)
	case 2:
		weapon = NewWeapon("Longsword", 1, 8, MeleeWeapon)
	case 3:
		weapon = NewWeapon("Shortbow", 1, 6, RangedWeapon)
	case 4:
		weapon = NewWeapon("Longbow", 1, 8, RangedWeapon)
	}

	return weapon
}
