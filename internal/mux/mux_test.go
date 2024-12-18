package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	//"github.com/go-park-mail-ru/2024_2_VKatuny/internal/mux"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestNotFoundHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", NotFoundHandler)
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}

func TestMethodNotAllowedHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", MethodNotAllowedHandler)
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}
