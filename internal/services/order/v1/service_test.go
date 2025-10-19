package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"testing"

	"github.com/sgrumley/kart-challenge/internal/services/order/v1/mapper"
	"github.com/sgrumley/kart-challenge/pkg/logger"
	"github.com/sgrumley/kart-challenge/pkg/models"
	"github.com/sgrumley/kart-challenge/pkg/testhelper"
	"github.com/sgrumley/kart-challenge/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type OrderRequestOption func(order *mapper.CreateOrderRequest)

func NewDefaultOrderRequest(opts ...OrderRequestOption) mapper.CreateOrderRequest {
	defaultOrder := mapper.CreateOrderRequest{
		Items: []mapper.Item{
			{
				ProductID: "00000000-0000-0000-0000-000000000001",
				Quantity:  1,
			},
			{
				ProductID: "00000000-0000-0000-0000-000000000002",
				Quantity:  2,
			},
		},
		CouponCode: "FIFTYOFF",
	}

	for _, opt := range opts {
		opt(&defaultOrder)
	}

	return defaultOrder
}

func Test_API_Service_CreateOrder(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		req           func() *mapper.CreateOrderRequest
		headers       map[string]string
		idemMock      *IdempotencyStoreMock
		storeMock     *OrderStorableMock
		wantAssertion func(t *testing.T, got *http.Response, storeMock *OrderStorableMock)
	}{
		"success/happy_path": {
			req: func() *mapper.CreateOrderRequest {
				def := NewDefaultOrderRequest()
				return &def
			},
			headers: map[string]string{
				"Idempotency-Key": "key",
				"api_key":         "a-secret-key",
			},
			idemMock: &IdempotencyStoreMock{
				ExistsFunc: func(key string) bool {
					return false
				},
				SetFunc: func(key string) {},
			},
			storeMock: &OrderStorableMock{
				CreateOrderFunc: func(ctx context.Context, order models.Order) (models.Order, error) {
					return models.Order{
						ID: "12300000-0000-0000-0000-000000000000",
						Items: []models.Item{
							{
								ProductID: "00000000-0000-0000-0000-000000000001",
								Quantity:  1,
							},
							{
								ProductID: "00000000-0000-0000-0000-000000000002",
								Quantity:  2,
							},
						},
						Products: []models.Product{
							{
								ID:       "00000000-0000-0000-0000-000000000001",
								Name:     "Eggs",
								Category: "Breakfast",
								Price:    8.99,
							},
							{
								ID:       "00000000-0000-0000-0000-000000000002",
								Name:     "Bacon",
								Category: "Breakfast",
								Price:    7.99,
							},
						},
					}, nil
				},
			},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *OrderStorableMock) {
				require.Equal(t, http.StatusCreated, got.StatusCode)
				require.Len(t, storeMock.CreateOrderCalls(), 1)
				defaultRequest := NewDefaultOrderRequest()
				expectedItems := []models.Item{
					{
						ProductID: "00000000-0000-0000-0000-000000000001",
						Quantity:  1,
					},
					{
						ProductID: "00000000-0000-0000-0000-000000000002",
						Quantity:  2,
					},
				}

				expectedStoreCalledWith := models.Order{
					CouponCode: defaultRequest.CouponCode,
					Items:      expectedItems,
				}
				assert.Equal(t, expectedStoreCalledWith, storeMock.CreateOrderCalls()[0].Order)

				want := &mapper.CreateOrderResponse{
					ID:    "12300000-0000-0000-0000-000000000000",
					Items: defaultRequest.Items,
					Products: []mapper.Product{
						{
							ID:       "00000000-0000-0000-0000-000000000001",
							Name:     "Eggs",
							Category: "Breakfast",
							Price:    8.99,
						},
						{
							ID:       "00000000-0000-0000-0000-000000000002",
							Name:     "Bacon",
							Category: "Breakfast",
							Price:    7.99,
						},
					},
				}

				actual := testhelper.PayloadAsType[mapper.CreateOrderResponse](t, got.Body)
				assert.Equal(t, want, &actual)
			},
		},
		"error/missing_idempotency_key": {
			req: func() *mapper.CreateOrderRequest {
				def := NewDefaultOrderRequest()
				return &def
			},
			headers: map[string]string{
				"Idempotency-Key": "",
			},
			idemMock:  &IdempotencyStoreMock{},
			storeMock: &OrderStorableMock{},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *OrderStorableMock) {
				require.Equal(t, http.StatusBadRequest, got.StatusCode)
				require.Len(t, storeMock.CreateOrderCalls(), 0)
				actual := testhelper.PayloadAsType[web.ErrorResponse](t, got.Body)
				expectedError := testhelper.MapExpectedErrorResponse(web.Err500Default)
				assert.Equal(t, expectedError, actual)
			},
		},
		"error/repeated_request": {
			req: func() *mapper.CreateOrderRequest {
				def := NewDefaultOrderRequest()
				return &def
			},
			headers: map[string]string{
				"Idempotency-Key": "used",
			},
			idemMock: &IdempotencyStoreMock{
				ExistsFunc: func(key string) bool {
					return true
				},
			},
			storeMock: &OrderStorableMock{},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *OrderStorableMock) {
				require.Equal(t, http.StatusInternalServerError, got.StatusCode)
				require.Len(t, storeMock.CreateOrderCalls(), 0)
				actual := testhelper.PayloadAsType[web.ErrorResponse](t, got.Body)
				expectedError := testhelper.MapExpectedErrorResponse(web.Err500Default)
				assert.Equal(t, expectedError, actual)
			},
		},
		"error/validation_error": {
			req: func() *mapper.CreateOrderRequest {
				def := NewDefaultOrderRequest()
				def.Items[0].Quantity = 0
				def.CouponCode = "inv"
				return &def
			},
			headers: map[string]string{
				"Idempotency-Key": "key",
				"api_key":         "a-secret-key",
			},
			idemMock: &IdempotencyStoreMock{
				ExistsFunc: func(key string) bool {
					return false
				},
			},
			storeMock: &OrderStorableMock{},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *OrderStorableMock) {
				require.Equal(t, http.StatusUnprocessableEntity, got.StatusCode)
				require.Len(t, storeMock.CreateOrderCalls(), 0)
				actual := testhelper.PayloadAsType[web.ErrorResponse](t, got.Body)
				expectedError := testhelper.MapExpectedErrorResponse(Err422Validation)
				assert.Equal(t, expectedError, actual)
			},
		},
		"error/store_failed": {
			req: func() *mapper.CreateOrderRequest {
				def := NewDefaultOrderRequest()
				return &def
			},
			headers: map[string]string{
				"Idempotency-Key": "key",
				"api_key":         "a-secret-key",
			},
			idemMock: &IdempotencyStoreMock{
				ExistsFunc: func(key string) bool {
					return false
				},
			},
			storeMock: &OrderStorableMock{
				CreateOrderFunc: func(ctx context.Context, order models.Order) (models.Order, error) {
					return models.Order{}, fmt.Errorf("error")
				},
			},
			wantAssertion: func(t *testing.T, got *http.Response, storeMock *OrderStorableMock) {
				require.Equal(t, http.StatusInternalServerError, got.StatusCode)
				require.Len(t, storeMock.CreateOrderCalls(), 1)
				actual := testhelper.PayloadAsType[web.ErrorResponse](t, got.Body)
				expectedError := testhelper.MapExpectedErrorResponse(web.Err500Default)
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

			svc := NewService(tc.storeMock, tc.idemMock)
			testServer := testhelper.SetupServer(svc, *log)

			url := fmt.Sprintf("%s/api/v1/order", testServer.URL)
			res := testhelper.SendRequest[mapper.CreateOrderRequest](t, "POST", url, tc.req(), tc.headers)
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
