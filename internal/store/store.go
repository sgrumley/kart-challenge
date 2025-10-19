package store

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/sgrumley/kart-challenge/internal/store/dbgen"
	"github.com/sgrumley/kart-challenge/pkg/logger"
	"github.com/sgrumley/kart-challenge/pkg/models"

	"github.com/google/uuid"
)

type Store struct {
	Queries *dbgen.Queries
	DB      *sqlx.DB
}

func New(client *sqlx.DB) *Store {
	return &Store{
		Queries: dbgen.New(client),
		DB:      client,
	}
}

func (s *Store) GetProduct(ctx context.Context, id string) (models.Product, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return models.Product{}, err
	}

	product, err := s.Queries.GetProductByID(ctx, uid)
	if err != nil {
		return models.Product{}, err
	}

	return ProductFromDB(product), nil
}

func ProductFromDB(product dbgen.Product) models.Product {
	return models.Product{
		ID:       product.ID.String(),
		Name:     product.Name,
		Category: product.Category.String,
		Price:    float32(product.Price),
	}
}

func (s *Store) ListProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.Queries.ListProducts(ctx)
	if err != nil {
		return []models.Product{}, err
	}

	return ProductsFromDB(products), nil
}

func ProductsFromDB(products []dbgen.Product) []models.Product {
	res := make([]models.Product, len(products))
	for i, p := range products {
		res[i] = models.Product{
			ID:       p.ID.String(),
			Name:     p.Name,
			Category: p.Category.String,
			Price:    float32(p.Price),
		}
	}
	return res
}

func (s *Store) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return models.Order{}, err
	}

	qtx := s.Queries.WithTx(tx)

	// create order
	orderID := GenerateUUIDv4()
	logger.Info(ctx, "creating order", slog.String("id", orderID.String()))
	_, err = qtx.CreateOrder(ctx, dbgen.CreateOrderParams{
		ID: orderID,
		CouponCode: sql.NullString{
			String: order.CouponCode,
			Valid:  order.CouponCode != "",
		},
		CreatedAt: int64(TimeStampNow()),
	})
	if err != nil {
		return models.Order{}, nil
	}

	products := make([]models.Product, len(order.Items))

	for _, product := range order.Items {
		pid, err := uuid.Parse(product.ProductID)
		if err != nil {
			return models.Order{}, fmt.Errorf("product id %s was not uuid: %w", product.ProductID, err)
		}

		p, err := s.Queries.GetProductByID(ctx, pid)
		if err != nil {
			return models.Order{}, fmt.Errorf("product id %s not found: %w", product.ProductID, err)
		}

		products = append(products, ProductFromDB(p))

		_, err = qtx.AddProductToOrder(ctx, dbgen.AddProductToOrderParams{
			ID:        GenerateUUIDv4(),
			OrderID:   orderID,
			ProductID: pid,
		})
		if err != nil {
			return models.Order{}, err
		}
	}

	return models.Order{
		ID:       orderID.String(),
		Items:    order.Items,
		Products: products,
	}, tx.Commit()
}

func (s *Store) CheckCoupon(ctx context.Context, coupon string) bool {
	matches := make([]string, 0)
	for i := 1; i < 4; i++ {
		couponID := fmt.Sprintf("%s-%d", coupon, i)
		logger.Info(ctx, "couponID", slog.String("id", couponID))
		match, err := s.Queries.GetCouponByID(ctx, couponID)
		if err == nil {
			matches = append(matches, match)
		}
	}

	logger.Info(ctx, "coupon matches", slog.Any("id", matches))
	return len(matches) > 1
}
