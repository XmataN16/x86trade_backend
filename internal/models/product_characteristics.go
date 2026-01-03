package models

type ProductCharacteristic struct {
	ID                 int    `json:"id"`
	ProductID          int    `json:"product_id"`
	ProductName        string `json:"product_name,omitempty"`
	CharacteristicName string `json:"characteristic_name"`
	CharacteristicUnit string `json:"characteristic_unit"`
	Value              string `json:"value"`
}

type ProductCharacteristicInput struct {
	ID                   int    `json:"id,omitempty"`
	ProductID            int    `json:"product_id"`
	CharacteristicTypeID int    `json:"characteristic_type_id"`
	Value                string `json:"value"`
}

type CharacteristicType struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Unit       string `json:"unit,omitempty"`
	CategoryID int    `json:"category_id,omitempty"`
}
