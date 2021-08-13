package dnd

type Cost struct {
	Quantity int    `json:"quantity" bson:"quantity"`
	Unit     string `json:"unit" bson:"unit"`
}
