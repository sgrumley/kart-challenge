package v1

import (
	"github.com/go-chi/chi/v5"
)

func (s *ProductService) GetRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		// r.Use(middleware.Authenticate())
		r.Get("/product/{product_id}", s.GetProduct)
		r.Get("/product", s.ListProducts)
	})
}
