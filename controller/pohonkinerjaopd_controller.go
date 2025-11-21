package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PohonKinerjaOpdController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindStrategicNoParent(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeletePelaksana(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinByPelaksana(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeletePokinPemdaInOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateParent(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindidPokinWithAllTema(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Clone(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckPokinExistsByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CountPokinPemda(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinAtasan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ControlPokinOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	LeaderboardPokinOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
