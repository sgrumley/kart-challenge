package v1

import (
	"github.com/go-chi/chi/v5"
)

func (s *OrderService) GetRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/order", s.CreateOrder)
	})
}
