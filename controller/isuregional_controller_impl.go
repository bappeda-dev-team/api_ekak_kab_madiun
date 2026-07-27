package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/isuregional"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type IsuRegionalControllerImpl struct {
	IsuRegionalService service.IsuRegionalService
}

func NewIsuRegionalControllerImpl(isuRegionalService service.IsuRegionalService) *IsuRegionalControllerImpl {
	return &IsuRegionalControllerImpl{
		IsuRegionalService: isuRegionalService,
	}
}

func (controller *IsuRegionalControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuRegionalRequest := isuregional.IsuRegionalRequest{}
	helper.ReadFromRequestBody(request, &isuRegionalRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	isuRegionalResponse, err := controller.IsuRegionalService.Create(request.Context(), isuRegionalRequest)
	if err != nil {
		webResponse := web.WebResponse{
			// TODO: CODE: AMBIL DARI http
			Code: http.StatusInternalServerError,
			// TODO: STATUS: TERJEMAHKAN DARI code
			Status: http.StatusText(http.StatusInternalServerError),
			// TODO: buat nil saja
			Data: err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		// TODO: CODE AMBIL DARI http
		Code:   201,
		Status: "Success Created Isu Regional",
		Data:   isuRegionalResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuRegionalControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuRegionalUpdateRequest := isuregional.IsuRegionalUpdateRequest{}
	helper.ReadFromRequestBody(request, &isuRegionalUpdateRequest)

	idStr := params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	isuRegionalUpdateRequest.ID = id

	isuRegionalResponse, err := controller.IsuRegionalService.Update(request.Context(), isuRegionalUpdateRequest)
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
		Status: "Success Updated Isu Regional",
		Data:   isuRegionalResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuRegionalControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuRegionalId := params.ByName("id")
	id, err := strconv.Atoi(isuRegionalId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.IsuRegionalService.Delete(request.Context(), id)
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
		Status: "Success Deleted Isu Regional",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuRegionalControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	kodeOpd := params.ByName("kode_opd")

	bidangUrusanResponses, err := controller.IsuRegionalService.FindAll(request.Context(), kodeOpd)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bidangUrusanResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}