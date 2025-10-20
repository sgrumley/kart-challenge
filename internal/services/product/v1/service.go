package v1

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sgrumley/kart-challenge/internal/services/product/v1/mapper"
	"github.com/sgrumley/kart-challenge/internal/store"
	"github.com/sgrumley/kart-challenge/pkg/logger"
	"github.com/sgrumley/kart-challenge/pkg/models"
	"github.com/sgrumley/kart-challenge/pkg/web"
)

//go:generate moq -out ./mocks_test.go . ProductStorable

var _ ProductStorable = (*store.Store)(nil)

type ProductStorable interface {
	GetProduct(ctx context.Context, id string) (models.Product, error)
	ListProducts(ctx context.Context) ([]models.Product, error)
}

func NewService(store ProductStorable) *ProductService {
	return &ProductService{
		store: store,
	}
}

type ProductService struct {
	store ProductStorable
}

var (
	Err400InvalidProductID = &web.Error{
		Status:      http.StatusBadRequest,
		Code:        "invalid_product_id",
		Description: "Invalid ID supplied",
	}

	Err404ProductNotFound = &web.Error{
		Status:      http.StatusNotFound,
		Code:        "product_not_found",
		Description: "Product not found",
	}
)

func (s *ProductService) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := chi.URLParam(r, "product_id")
	if _, err := uuid.Parse(productID); err != nil {
		logger.Error(ctx, "invalid product id is not uuid", Err400InvalidProductID)
		web.RespondJSONError(w, Err400InvalidProductID)
		return
	}

	product, err := s.store.GetProduct(ctx, productID)
	if err != nil {
		logger.Error(ctx, "could not find product with id: "+productID, err)
		web.RespondJSONError(w, Err404ProductNotFound)
		return
	}

	web.Respond(w, http.StatusOK, mapper.GetProductToResponse(&product))
}

func (s *ProductService) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, _ := s.store.ListProducts(ctx)

	web.Respond(w, http.StatusOK, mapper.ListProductsToResponse(products))
}
