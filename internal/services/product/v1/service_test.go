package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"testing"

	"github.com/sgrumley/kart-challenge/internal/services/product/v1/mapper"
	"github.com/sgrumley/kart-challenge/pkg/logger"
	"github.com/sgrumley/kart-challenge/pkg/models"
	"github.com/sgrumley/kart-challenge/pkg/testhelper"
	"github.com/sgrumley/kart-challenge/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_API_Service_GetProduct(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		productID     string
		storeMock     *ProductStorableMock
		wantAssertion func(t *testing.T, got *http.Response, storeMock *ProductStorableMock)
	}{
		"success/happy_path": {
			productID: "00000000-0000-0000-0000-000000000001",
			storeMock: &ProductStorableMock{
				GetProductFunc: func(ctx context.Context, id string) (models.Product, error) {
					return models.Product{
						ID:       "00000000-0000-0000-0000-000000000001",
						Name:     "eggs",
						Category: "breakfast",
						Price:    8.99,
					}, nil
				},
			},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *ProductStorableMock) {
				requestID := "00000000-0000-0000-0000-000000000001"
				require.Equal(t, http.StatusOK, got.StatusCode)
				require.Len(t, storeMock.GetProductCalls(), 1)
				assert.Equal(t, requestID, storeMock.GetProductCalls()[0].ID)

				want := &mapper.GetProductResponse{
					ID:        requestID,
					Name:      "eggs",
					Catergory: "breakfast",
					Price:     8.99,
				}

				actual := testhelper.PayloadAsType[mapper.GetProductResponse](t, got.Body)
				assert.Equal(t, want, &actual)
			},
		},
		"error/missing_product_id": {
			productID: "",
			storeMock: &ProductStorableMock{},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *ProductStorableMock) {
				require.Equal(t, http.StatusNotFound, got.StatusCode)
				require.Len(t, storeMock.GetProductCalls(), 0)
			},
		},
		"error/invalid_product_id": {
			productID: "invalid-uuid",
			storeMock: &ProductStorableMock{},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *ProductStorableMock) {
				require.Equal(t, http.StatusBadRequest, got.StatusCode)
				require.Len(t, storeMock.GetProductCalls(), 0)

				actual := testhelper.PayloadAsType[web.ErrorResponse](t, got.Body)
				expectedError := testhelper.MapExpectedErrorResponse(Err400InvalidProductID)
				assert.Equal(t, expectedError, actual)
			},
		},
		"error/store_not_found": {
			productID: "00000000-0000-0000-0000-000000000001",
			storeMock: &ProductStorableMock{
				GetProductFunc: func(ctx context.Context, id string) (models.Product, error) {
					return models.Product{}, fmt.Errorf("error")
				},
			},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *ProductStorableMock) {
				require.Equal(t, http.StatusNotFound, got.StatusCode)
				requestID := "00000000-0000-0000-0000-000000000001"
				require.Len(t, storeMock.GetProductCalls(), 1)
				assert.Equal(t, requestID, storeMock.GetProductCalls()[0].ID)

				actual := testhelper.PayloadAsType[web.ErrorResponse](t, got.Body)
				expectedError := testhelper.MapExpectedErrorResponse(Err404ProductNotFound)
				assert.Equal(t, expectedError, actual)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(
				logger.WithLevel(slog.LevelDebug),
				logger.WithFormat(logger.HandlerJSON),
			)

			svc := NewService(tc.storeMock)
			testServer := testhelper.SetupServer(svc, *log)

			url := fmt.Sprintf("%s/api/v1/product/%s", testServer.URL, tc.productID)
			res := testhelper.SendRequest[any](t, "GET", url, nil, nil)
			t.Cleanup(func() {
				if res.Body != nil {
					require.NoError(t, res.Body.Close())
				}
			})
			tc.wantAssertion(t, res, tc.storeMock)
			testServer.Close()
		})
	}
}

func Test_API_Service_ListProducts(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		storeMock     *ProductStorableMock
		wantAssertion func(t *testing.T, got *http.Response, storeMock *ProductStorableMock)
	}{
		"success/happy_path": {
			storeMock: &ProductStorableMock{
				ListProductsFunc: func(ctx context.Context) ([]models.Product, error) {
					return []models.Product{
						{
							ID:       "00000000-0000-0000-0000-000000000001",
							Name:     "eggs",
							Category: "breakfast",
							Price:    8.99,
						},
						{
							ID:       "00000000-0000-0000-0000-000000000002",
							Name:     "bacon",
							Category: "breakfast",
							Price:    6.99,
						},
					}, nil
				},
			},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *ProductStorableMock) {
				require.Equal(t, http.StatusOK, got.StatusCode)
				require.Len(t, storeMock.ListProductsCalls(), 1)

				want := &mapper.ListProductsResponse{
					{
						ID:       "00000000-0000-0000-0000-000000000001",
						Name:     "eggs",
						Category: "breakfast",
						Price:    8.99,
					},
					{
						ID:       "00000000-0000-0000-0000-000000000002",
						Name:     "bacon",
						Category: "breakfast",
						Price:    6.99,
					},
				}

				actual := testhelper.PayloadAsType[mapper.ListProductsResponse](t, got.Body)
				assert.Equal(t, want, &actual)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(
				logger.WithLevel(slog.LevelDebug),
				logger.WithFormat(logger.HandlerJSON),
			)

			svc := NewService(tc.storeMock)
			testServer := testhelper.SetupServer(svc, *log)

			url := fmt.Sprintf("%s/api/v1/product", testServer.URL)
			res := testhelper.SendRequest[any](t, "GET", url, nil, nil)
			t.Cleanup(func() {
				if res.Body != nil {
					require.NoError(t, res.Body.Close())
				}
			})
			tc.wantAssertion(t, res, tc.storeMock)
			testServer.Close()
		})
	}
}
