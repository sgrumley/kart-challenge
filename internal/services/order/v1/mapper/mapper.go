package mapper

import "github.com/sgrumley/kart-challenge/pkg/models"

type Item struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type CreateOrderRequest struct {
	Items      []Item `json:"items" validate:"required"`
	CouponCode string `json:"coupon_code" validate:"omitempty,min=8,max=10"`
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"catergory"`
	Price    float32 `json:"price"`
}

type CreateOrderResponse struct {
	ID       string    `json:"id"`
	Items    []Item    `json:"items"`
	Products []Product `json:"products"`
}

func ItemsFromRequest(items []Item) []models.Item {
	modelItems := make([]models.Item, len(items))
	for i := range items {
		modelItems[i] = models.Item{
			ProductID: items[i].ProductID,
			Quantity:  items[i].Quantity,
		}
	}
	return modelItems
}

func CreateOrderFromRequest(req CreateOrderRequest) models.Order {
	return models.Order{
		CouponCode: req.CouponCode,
		Items:      ItemsFromRequest(req.Items),
	}
}

func ItemsToResponse(items []models.Item) []Item {
	responseItems := make([]Item, len(items))
	for i, it := range items {
		responseItems[i] = Item{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
		}
	}

	return responseItems
}

func ProductsToResponse(products []models.Product) []Product {
	responseProducts := make([]Product, len(products))
	for i, p := range products {
		responseProducts[i] = Product{
			ID:       p.ID,
			Name:     p.Name,
			Category: p.Category,
			Price:    p.Price,
		}
	}

	return responseProducts
}

func CreateOrderToResponse(res models.Order) CreateOrderResponse {
	return CreateOrderResponse{
		ID:       res.ID,
		Items:    ItemsToResponse(res.Items),
		Products: ProductsToResponse(res.Products),
	}
}
