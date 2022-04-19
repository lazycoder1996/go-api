package models

type Product struct {
	GUID        string  `json:"guid"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"createdAt"`
}
