package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/hollgett/shortURL.git/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouters(handler *HandlerAPI) *chi.Mux {
	r := chi.NewMux()
	r.Post("/", handler.HandlePlainTextRequest)
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handler.HandleJSONRequest)
	})
	r.Get("/{short}", handler.ShortURLGet)

	return r
}

func simulateMockServer(ctrl *gomock.Controller, createURL, getURL string, createError, getError error) *httptest.Server {
	controller := mocks.NewMockShortenerHandler(ctrl)

	controller.EXPECT().CreateShortURL(gomock.AssignableToTypeOf("")).Return(createURL, createError).AnyTimes()
	controller.EXPECT().GetShortURL(gomock.AssignableToTypeOf("")).Return(getURL, getError).AnyTimes()

	api := NewHandlerAPI(controller)
	rtr := setupRouters(api)
	return httptest.NewServer(rtr)
}

func simulateRequest(methodGet bool, baseURL, bodyReq, pathReq string) (*resty.Response, error) {
	client := resty.NewWithClient(
		&http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	).
		SetBaseURL(baseURL).
		SetHeader("Content-Type", "text/plain")
	if methodGet {
		return client.R().Get(pathReq)
	}
	resp := client.R()
	if len(bodyReq) != 0 {
		resp.SetBody(bodyReq)
	}
	return resp.Post(pathReq)
}

func TestHandlerAPI_ShortURLPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name        string
		server      *httptest.Server
		requestBody string
		expected    expected
	}{
		{"positive test", simulateMockServer(ctrl, "test", "", nil, nil), "https://mail.google.com/", expected{http.StatusCreated, "/test"}},
		{"negative test request body", simulateMockServer(ctrl, "test", "", nil, nil), "", expected{http.StatusBadRequest, "request body empty"}},
		{"negative test CreateShortURL", simulateMockServer(ctrl, "", "", errors.New("test error"), nil), "https://mail.google.com/", expected{http.StatusBadRequest, "test error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			resp, err := simulateRequest(false, tt.server.URL, tt.requestBody, "/")
			require.NoError(t, err, "request error catch")

			assert.Equal(t, tt.expected.code, resp.StatusCode(), "request code not equal")
			assert.Contains(t, resp.String(), tt.expected.body, "request body not equal")
		})
	}
}

func TestHandlerAPI_ShortURLGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type expected struct {
		code     int
		location string
	}
	tests := []struct {
		name     string
		server   *httptest.Server
		expected expected
	}{
		{"positive test", simulateMockServer(ctrl, "", "https://mail.google.com/", nil, nil), expected{http.StatusTemporaryRedirect, "https://mail.google.com/"}},
		{"negative GetShortURL error", simulateMockServer(ctrl, "", "", nil, errors.New("test error")), expected{http.StatusBadRequest, ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			resp, err := simulateRequest(true, tt.server.URL, "", "/test")
			require.NoError(t, err, "request error")

			assert.Equal(t, tt.expected.code, resp.StatusCode(), "response status code not equal")
			assert.Equal(t, tt.expected.location, resp.Header().Get("Location"), "response location not equal")
		})
	}
}
