package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenstraController interface {
	GetByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindIndikatorById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
