package mapper

import "github.com/sgrumley/kart-challenge/pkg/models"

type GetProductResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Catergory string  `json:"catergory"`
	Price     float32 `json:"price"`
}

func GetProductToResponse(product *models.Product) *Product {
	return &Product{
		ID:       product.ID,
		Name:     product.Name,
		Category: product.Category,
		Price:    product.Price,
	}
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"catergory"`
	Price    float32 `json:"price"`
}

type ListProductsResponse []Product

func ListProductsToResponse(products []models.Product) ListProductsResponse {
	res := make(ListProductsResponse, len(products))
	for i, p := range products {
		res[i] = Product{
			ID:       p.ID,
			Name:     p.Name,
			Category: p.Category,
			Price:    p.Price,
		}
	}
	return res
}
