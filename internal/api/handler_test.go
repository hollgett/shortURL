package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	gomock "github.com/golang/mock/gomock"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRequest(t *testing.T, ts *httptest.Server, reqBody, method, path string) *resty.Response {
	client := &http.Client{
		Transport: ts.Client().Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	clientResty := resty.NewWithClient(client).
		SetBaseURL(ts.URL).
		SetHeader("Content-Type", "text/plain")
	var resp *resty.Response
	var err error
	if method == http.MethodGet {
		resp, err = clientResty.R().Get(path)
	} else {
		resp, err = clientResty.R().SetBody(reqBody).Post(path)
	}
	require.NoError(t, err)

	return resp
}

func newServer(short app.ShortenerHandler) *httptest.Server {
	rtr := chi.NewRouter()
	ts := httptest.NewServer(rtr)
	cfg := &config.Config{
		Addr:    strings.Split(ts.URL, ":")[2],
		BaseURL: ts.URL,
	}
	api := NewHandlerAPI(short, cfg)

	rtr.Post("/", api.ShortURLPost)
	rtr.Get("/{short}", api.ShortURLGet)

	return ts
}

func simulateShortener(ctrl *gomock.Controller, url string, err error) app.ShortenerHandler {
	controller := mock.NewMockShortenerHandler(ctrl)

	controller.EXPECT().CreateShortURL(gomock.Any()).Return(url, err).AnyTimes()
	controller.EXPECT().GetShortURL(gomock.Any()).Return(url, err).AnyTimes()

	return controller
}

func TestRouters_shortURLPost(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	type want struct {
		expectedCode int
		contentType  string
	}
	tests := []struct {
		name          string
		want          want
		short         app.ShortenerHandler
		request       string
		expectedShort string
	}{
		{
			name: "Positive test",
			want: want{
				expectedCode: http.StatusCreated,
				contentType:  "text/plain",
			},
			short:         simulateShortener(controller, "gg", nil),
			request:       "https://mail.google.com/",
			expectedShort: "/gg",
		},
		{
			name: "negative test without request",
			want: want{
				expectedCode: http.StatusBadRequest,
				contentType:  "text/plain; charset=utf-8",
			},
			short:         simulateShortener(controller, "", errors.New("request body error")),
			request:       "/",
			expectedShort: "request body error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newServer(tt.short)
			defer ts.Close()
			resp := newRequest(t, ts, tt.request, http.MethodPost, `/`)
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"), "Передаваемый тип контента не совпадает с ожидаемым")
			assert.Equal(t, tt.expectedShort, strings.Replace(strings.TrimSpace(resp.String()), ts.URL, "", -1))
		})
	}
}

func TestRouters_shortURLGet(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	type want struct {
		expectedCode     int
		expectedLocation string
	}
	tests := []struct {
		name    string
		short   app.ShortenerHandler
		request string
		want    want
	}{
		{
			name:    "Positive test #1",
			short:   simulateShortener(controller, "https://go.dev/", nil),
			request: "gg",
			want: want{
				expectedCode:     http.StatusTemporaryRedirect,
				expectedLocation: "https://go.dev/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newServer(tt.short)
			defer ts.Close()
			resp := newRequest(t, ts, "", http.MethodGet, "/"+tt.request)

			assert.Equal(t, tt.want.expectedCode, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.expectedLocation, resp.Header().Get("Location"), "ответ header Location не совпадает с ожидаемым")
		})
	}
}
