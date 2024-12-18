package mux

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/mailru/easyjson"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusNotFound)
	response := &dto.JSONResponse{
		HTTPStatus: http.StatusNotFound,
		Error:      commonerrors.ErrFrontServiceNotFound.Error(),
	}
	JSONResponse, err := easyjson.Marshal(response)
	if err != nil {
		return
	}
	w.Write(JSONResponse)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusMethodNotAllowed)
	response := &dto.JSONResponse{
		HTTPStatus: http.StatusMethodNotAllowed,
		Error:      commonerrors.ErrFrontMethodNotAllowed.Error(),
	}
	JSONResponse, err := easyjson.Marshal(response)
	if err != nil {
		return
	}
	w.Write(JSONResponse)
}
