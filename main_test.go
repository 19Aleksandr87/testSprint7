package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeStatusBadRequest(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe?", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tyla&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, r := range request {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", r.request, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, r.status, response.Code)
		assert.Equal(t, r.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		count           string
		status, lengtch int // ожидаемое количество кафе в ответе
	}{
		{"/cafe?city=moscow&count=0", http.StatusOK, 0},
		{"/cafe?city=moscow&count=1", http.StatusOK, 1},
		{"/cafe?city=moscow&count=2", http.StatusOK, 2},
		{"/cafe?city=moscow&count=100", http.StatusOK, min(100, len(cafeList["moscow"]))},
	}

	for _, r := range request {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", r.count, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, r.status, response.Code)
		city := req.URL.Query().Get("city")

		i := strings.Split(response.Body.String(), ",")
		j := len(i)
		if i[0] == "" {
			j = 0
		}

		if r.lengtch > len(cafeList[city]) {
			assert.Equal(t, r.lengtch, j)
			continue
		}
		assert.Equal(t, j, r.lengtch)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search            string
		status, wantCount int
	}{
		{"/cafe?city=moscow&search=фасоль", http.StatusOK, 0},
		{"/cafe?city=moscow&search=кофе", http.StatusOK, 2},
		{"/cafe?city=moscow&search=вилка", http.StatusOK, 1},
	}

	for _, r := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", r.search, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, r.status, response.Code)

		strings.Split(response.Body.String(), ",")
		s := strings.Split(response.Body.String(), ",")
		j := len(s)
		if s[0] == "" {
			j = 0
		}
		assert.Equal(t, r.wantCount, j)

		if j == 0 {
			continue
		}

		for _, i := range s {
			search := req.URL.Query().Get("search")
			str := strings.ToLower(strings.ToUpper(i))
			assert.Equal(t, true, strings.Contains(str, search))
		}

	}

}
