package arena

type ArmorType string
type WeaponType string

const (
	NoArmor      ArmorType  = "none"
	LightArmor   ArmorType  = "light"
	MediumArmor  ArmorType  = "medium"
	HeavyArmor   ArmorType  = "heavy"
	MeleeWeapon  WeaponType = "melee"
	RangedWeapon WeaponType = "ranged"
)

type Armor struct {
	Name        string    `json:"name"`
	Value       int       `json:"value"`
	Type        ArmorType `json:"type"`
	MaxDexBonus int       `json:"max_dex_bonus"`
}

type Caracteristic struct {
	Name     string `json:"name"`
	Value    int    `json:"value"`
	Modifier int    `json:"modifier"`
}

type Weapon struct {
	Name   string     `json:"name"`
	Number int        `json:"number"`
	Damage int        `json:"damage"`
	Type   WeaponType `json:"type"`
}

type Gladiator struct {
	Armor                 *Armor         `json:"armor"`
	ArmorClass            int            `json:"armor_class"`
	Constitution          *Caracteristic `json:"constitution"`
	CurrentHealth         int            `json:"current_health"`
	Dexterity             *Caracteristic `json:"dexterity"`
	Experience            int            `json:"experience"`
	ExperienceToNextLevel int            `json:"experience_to_next_level"`
	Level                 int            `json:"level"`
	MaxHealth             int            `json:"max_health"`
	Name                  string         `json:"name"`
	Strength              *Caracteristic `json:"strength"`
	Weapon                *Weapon        `json:"weapon"`
}
