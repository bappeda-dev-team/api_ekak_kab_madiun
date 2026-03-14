package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TujuanOpdController interface {
	CreateTujuanOpdRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTujuanOpdRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// CreateTujuanOpdRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// UpdateTujuanOpdRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// CreateTujuanOpdRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// UpdateTujuanOpdRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanOpdOnlyName(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanOpdByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanOpdRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanOpdRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanOpdRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateTujuanRenjaRanwalIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTujuanRenjaRanwalIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteTujuanRenjaIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateTujuanRenjaRankhirIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTujuanRenjaRankhirIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
