package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ReviewController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllReviewByTematik(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllReviewOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	CreateTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
