package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sgrumley/kart-challenge/internal/services/order/v1/mapper"
	"github.com/sgrumley/kart-challenge/internal/store"
	"github.com/sgrumley/kart-challenge/pkg/idempotency"
	"github.com/sgrumley/kart-challenge/pkg/logger"
	"github.com/sgrumley/kart-challenge/pkg/models"
	"github.com/sgrumley/kart-challenge/pkg/web"
)

//go:generate moq -out ./mocks_test.go . OrderStorable IdempotencyStore

var (
	_ OrderStorable    = (*store.Store)(nil)
	_ IdempotencyStore = (*idempotency.Store)(nil)
)

type OrderStorable interface {
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	CheckCoupon(ctx context.Context, coupon string) bool
}

type IdempotencyStore interface {
	Exists(key string) bool
	Set(key string)
	Remove(key string)
}

type OrderService struct {
	validate    *validator.Validate
	idemChecker IdempotencyStore
	store       OrderStorable
}

func NewService(store OrderStorable, idemChecker IdempotencyStore) *OrderService {
	return &OrderService{
		store:       store,
		idemChecker: idemChecker,
		validate:    validator.New(),
	}
}

var (
	Err401InvalidRequestBody = &web.Error{
		Status:      http.StatusBadRequest,
		Code:        "invalid_request_body",
		Description: "Invalid input",
	}

	Err409ConflictDuplicateRequest = &web.Error{
		Status:      http.StatusConflict,
		Code:        "request_already_inprogress",
		Description: "The provided Idempotency-Key has already been used",
	}

	Err422Validation = &web.Error{
		Status:      http.StatusUnprocessableEntity,
		Code:        "invalid_order_detail",
		Description: "Validation exception",
	}
)

func (s *OrderService) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		logger.Error(ctx, "missing header", fmt.Errorf("missing Idempotency-Key header"))
		web.RespondJSONError(w, Err409ConflictDuplicateRequest)
		return
	}

	if s.idemChecker.Exists(key) {
		logger.Error(ctx, "request already in progress", fmt.Errorf("duplicate Idempotency-Key"))
		web.RespondJSONError(w, Err409ConflictDuplicateRequest)
		return
	}

	var req mapper.CreateOrderRequest
	if err := web.DecodeBody(r, &req); err != nil {
		logger.Error(ctx, "invalid request body", err)
		web.RespondJSONError(w, Err401InvalidRequestBody)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		logger.Error(ctx, "validation failed", err)
		web.RespondJSONError(w, Err422Validation)
		return
	}

	if req.CouponCode != "" {
		if valid := s.store.CheckCoupon(ctx, req.CouponCode); !valid {
			logger.Error(ctx, "invalid coupon", fmt.Errorf("coupon was not in atleast 2 files"))
			web.RespondJSONError(w, Err422Validation)
			return
		}
	}

	order, err := s.store.CreateOrder(ctx, mapper.CreateOrderFromRequest(req))
	if err != nil {
		logger.Error(ctx, "failed creating order in store", err)
		web.RespondJSONError(w, fmt.Errorf("failed creating order in store: %w", err))
		return
	}

	s.idemChecker.Set(key)

	web.Respond(w, http.StatusCreated, mapper.CreateOrderToResponse(order))
}
