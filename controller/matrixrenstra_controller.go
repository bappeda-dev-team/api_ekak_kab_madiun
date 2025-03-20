package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenstraController interface {
	GetByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
