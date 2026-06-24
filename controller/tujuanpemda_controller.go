package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TujuanPemdaController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdatePeriode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllWithPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinWithPeriode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllWithPokinRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	FindTujuanPemdaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanPemdaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanPemdaPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertTargetPemdaLayer(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanPemdaRankhirDual(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindTujuanPemdaPenetapanDual(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateTargetRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateTargetPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTargetRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTargetPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
