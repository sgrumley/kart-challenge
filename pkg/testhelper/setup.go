package testhelper

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sgrumley/kart-challenge/pkg/middleware"
	"github.com/sgrumley/kart-challenge/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Service interface {
	GetRoutes(r chi.Router)
}

func SetupServer(service Service, log slog.Logger) *httptest.Server {
	testRouter := chi.NewRouter()
	testRouter.Use(middleware.AddLogger(log))
	v1Router := chi.NewRouter()

	service.GetRoutes(v1Router)
	testRouter.Mount("/api/v1", v1Router)

	return httptest.NewServer(testRouter)
}

func SendRequest[T any](t *testing.T, method string, path string, body *T, headers map[string]string) *http.Response {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	ctx := context.Background()
	var jsonReq []byte
	var err error

	if body == nil && method != "GET" {
		jsonReq = []byte("{invalid request body}")
	} else {
		jsonReq, err = json.Marshal(&body)
		require.NoError(t, err)
	}

	req, err := http.NewRequestWithContext(ctx, method, path, strings.NewReader(string(jsonReq)))
	require.NoError(t, err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	require.NoError(t, err)

	return res
}

func PayloadAsString(t *testing.T, body io.Reader) string {
	b, err := io.ReadAll(body)
	require.NoError(t, err)
	return string(b)
}

func PayloadAsType[T any](t *testing.T, body io.Reader) T {
	str := PayloadAsString(t, body)
	var actualResponse T
	err := json.Unmarshal([]byte(str), &actualResponse)
	require.NoError(t, err)

	return actualResponse
}

func MapExpectedErrorResponse(err *web.Error) web.ErrorResponse {
	return web.ErrorResponse{
		Error: &web.ErrorPayload{
			Code:    err.Code,
			Message: err.Description,
		},
	}
}

func AssertErrorHTTP(t *testing.T, expectedError web.ErrorResponse, responseBody string) {
	var actualErr web.ErrorResponse
	err := json.Unmarshal([]byte(responseBody), &actualErr)
	require.NoError(t, err)

	assert.Equal(t, expectedError, actualErr)
}
