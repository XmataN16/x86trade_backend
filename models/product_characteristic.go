package models

type ProductCharacteristic struct {
	ID                 int    `json:"id"`
	ProductID          int    `json:"product_id"`
	CharacteristicName string `json:"characteristic_name"` // Было "name"
	CharacteristicUnit string `json:"characteristic_unit"` // Было "unit"
	Value              string `json:"value"`
}
