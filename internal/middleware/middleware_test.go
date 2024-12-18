package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/stretchr/testify/require"
)

func TestUniversalMarshalExpectSuccessMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	err := middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK})
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.Nil(t, err)
}