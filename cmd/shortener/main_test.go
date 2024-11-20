package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_shortURLmiddleware(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
	}{
		{name: "Test Method Delete", method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed},
		{name: "Test Method Put", method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed},
		{name: "Test Method Post without data", method: http.MethodPost, expectedCode: http.StatusBadRequest},
		{name: "Test Method Get without data", method: http.MethodGet, expectedCode: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, `/`, nil)
			w := httptest.NewRecorder()
			h := shortURLmiddleware(&Routers{})
			h(w, r)
			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestRouters_shortURLPost(t *testing.T) {
	type want struct {
		expectedCode int
		contentType  string
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name: "Positive test #1",
			want: want{
				expectedCode: http.StatusCreated,
				contentType:  "text/plain",
			},
			request: "https://mail.google.com/",
		},
		{
			name: "Positive test #2",
			want: want{
				expectedCode: http.StatusCreated,
				contentType:  "text/plain",
			},
			request: "/",
		},
		{
			name: "Negative test without data #1",
			want: want{
				expectedCode: http.StatusBadRequest,
				contentType:  "text/plain",
			},
			request: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtr := newRouters()
			r := httptest.NewRequest(http.MethodPost, `/`, bytes.NewBuffer([]byte(tt.request)))
			w := httptest.NewRecorder()
			r.Header.Set("Content-Type", tt.want.contentType)
			rtr.shortURLPost(w, r)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.expectedCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
			if result.StatusCode != http.StatusBadRequest {
				assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"), "Передаваемый тип контента не совпадает с ожидаемым")

				res, err := io.ReadAll(result.Body)
				require.NoError(t, err, "Тело запроса имеет ошибку")
				link := strings.Replace(string(res), "http://localhost:8080/", "", -1)
				assert.Contains(t, rtr.routes, link, "Ответ не совпадает данными")
				assert.Equal(t, rtr.routes[link], tt.request, "Передаваемый ответ не совпадает с ожидаемым")
			}
		})
	}
}

func TestRouters_shortURLGet(t *testing.T) {
	type want struct {
		expectedCode     int
		expectedLocation string
	}
	tests := []struct {
		name    string
		rtr     *Routers
		request string
		want    want
	}{
		{
			name: "Positive test #1",
			rtr: &Routers{
				map[string]string{
					"gg": "https://go.dev/",
				},
			},
			request: `/gg`,
			want: want{
				expectedCode:     http.StatusTemporaryRedirect,
				expectedLocation: "https://go.dev/",
			},
		},
		{
			name: "Positive test with a lot of routers #2",
			rtr: &Routers{
				map[string]string{
					"gg":    "https://go.dev1/",
					"ggg":   "https://go.dev2/",
					"gggg":  "https://go.de3v/",
					"gggd":  "https://go.dev4/",
					"ggh":   "https://go.dev5/",
					"gggj":  "https://go.dev6/",
					"ggdgd": "https://go.dev7/",
					"gggh":  "https://go.dev8/",
					"ggjgj": "https://go.dev9/",
				},
			},
			request: `/ggdgd`,
			want: want{
				expectedCode:     http.StatusTemporaryRedirect,
				expectedLocation: "https://go.dev7/",
			},
		},
		{
			name: "Negative test without data routers #1",
			rtr: &Routers{
				map[string]string{},
			},
			request: `/gg`,
			want: want{
				expectedCode:     http.StatusBadRequest,
				expectedLocation: "",
			},
		},
		{
			name: "Negative test bad request #2",
			rtr: &Routers{
				map[string]string{
					"gg": "https://go.dev/",
				},
			},
			request: `/ggg`,
			want: want{
				expectedCode:     http.StatusBadRequest,
				expectedLocation: "",
			},
		},
		{
			name: "Negative test bad request #3",
			rtr: &Routers{
				map[string]string{
					"gg": "https://go.dev/",
				},
			},
			request: `/`,
			want: want{
				expectedCode:     http.StatusBadRequest,
				expectedLocation: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			tt.rtr.shortURLGet(w, r)
			assert.Equal(t, tt.want.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.expectedLocation, w.Header().Get("Location"), "ответ header Location не совпадает с ожидаемым")
		})
	}
}
