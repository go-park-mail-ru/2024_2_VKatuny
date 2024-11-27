package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/stretchr/testify/require"
)

func TestUniversalMarshalExpectSuccessMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	body := struct{
		Val string `json:"value"`
	} {
		Val: "test",
	}
	err := middleware.UniversalMarshal(w, http.StatusOK, body)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.Nil(t, err)
}

func TestUniversalMarshalExpectErrorMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	body := struct{
		Val chan int
	} {
		Val: make(chan int),
	}
	err := middleware.UniversalMarshal(w, http.StatusOK, body)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.Equal(t, fmt.Errorf(dto.MsgUnableToMarshalJSON), err)
}
