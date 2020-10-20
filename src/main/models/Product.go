package models

// Product is the representation of a product
//
// swagger:model
type Product struct {
	// The product's id
	ID int `json:"id" pg:"id,pk"`

	//The product's name
	//
	// required: true
	Name string `json:"name" validate:"required" pg:"name,unique,notnull"`

	//The product's price
	//
	// required: true
	//
	// min: 0
	Price float64 `json:"price" validate:"required,gt=0" pg:"price"`
}
