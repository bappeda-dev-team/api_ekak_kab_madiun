package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pegawai"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PegawaiControllerImpl struct {
	PegawaiService service.PegawaiService
}

func NewPegawaiControllerImpl(pegawaiService service.PegawaiService) *PegawaiControllerImpl {
	return &PegawaiControllerImpl{
		PegawaiService: pegawaiService,
	}
}

func (controller *PegawaiControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiCreateRequest := pegawai.PegawaiCreateRequest{}
	helper.ReadFromRequestBody(request, &pegawaiCreateRequest)

	pegawaiResponse, err := controller.PegawaiService.Create(request.Context(), pegawaiCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pegawaiResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PegawaiControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiId := params.ByName("id")

	pegawaiUpdateRequest := pegawai.PegawaiUpdateRequest{}
	helper.ReadFromRequestBody(request, &pegawaiUpdateRequest)

	pegawaiUpdateRequest.Id = pegawaiId

	pegawaiResponse, err := controller.PegawaiService.Update(request.Context(), pegawaiUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}

		if err.Error() == "INTERNAL_SERVER_ERROR" {
			webResponse.Code = 500
			webResponse.Status = "Internal Server Error"
		} else if err.Error() == "BAD_REQUEST" {
			webResponse.Code = 400
			webResponse.Status = "Bad Request"
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pegawaiResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PegawaiControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiId := params.ByName("id")

	controller.PegawaiService.Delete(request.Context(), pegawaiId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PegawaiControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiId := params.ByName("id")

	pegawaiResponse, err := controller.PegawaiService.FindById(request.Context(), pegawaiId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pegawaiResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PegawaiControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := request.URL.Query().Get("kode_opd")
	pegawaiResponses, err := controller.PegawaiService.FindAll(request.Context(), kodeOpd)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pegawaiResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PegawaiControllerImpl) TambahJabatanPegawai(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tambahJabatanRequest := pegawai.TambahJabatanRequest{}
	helper.ReadFromRequestBody(request, &tambahJabatanRequest)

	pegawaiResponse, err := controller.PegawaiService.TambahJabatan(request.Context(), tambahJabatanRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pegawaiResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
