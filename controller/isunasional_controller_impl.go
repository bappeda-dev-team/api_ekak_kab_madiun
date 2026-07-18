package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/isunasional"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type IsuNasionalControllerImpl struct {
	IsuNasionalService service.IsuNasionalService
}

func NewIsuNasionalControllerImpl(isuNasionalService service.IsuNasionalService) *IsuNasionalControllerImpl {
	return &IsuNasionalControllerImpl{
		IsuNasionalService: isuNasionalService,
	}
}

func (controller *IsuNasionalControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuNasionalRequest := isunasional.IsuNasionalRequest{}
	helper.ReadFromRequestBody(request, &isuNasionalRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	isuNasionalResponse, err := controller.IsuNasionalService.Create(request.Context(), isuNasionalRequest)
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
		Status: "Success Created Isu Nasional",
		Data:   isuNasionalResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuNasionalControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuNasionalUpdateRequest := isunasional.IsuNasionalUpdateRequest{}
	helper.ReadFromRequestBody(request, &isuNasionalUpdateRequest)

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
	isuNasionalUpdateRequest.ID = id

	isuNasionalResponse, err := controller.IsuNasionalService.Update(request.Context(), isuNasionalUpdateRequest)
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
		Status: "Success Updated Isu Nasional",
		Data:   isuNasionalResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuNasionalControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuNasionalId := params.ByName("id")
	id, err := strconv.Atoi(isuNasionalId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.IsuNasionalService.Delete(request.Context(), id)
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
		Status: "Success Deleted Isu Nasional",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}