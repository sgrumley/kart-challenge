package models

type Product struct {
	ID       string
	Name     string
	Category string
	Price    float32
}

type Order struct {
	ID         string
	CouponCode string
	Items      []Item
	Products   []Product
}

type Item struct {
	ProductID string
	Quantity  int
}
