package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/isuglobal"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type IsuGlobalControllerImpl struct {
	IsuGlobalService service.IsuGlobalService
}

func NewIsuGlobalControllerImpl(isuGlobalService service.IsuGlobalService) *IsuGlobalControllerImpl {
	return &IsuGlobalControllerImpl{
		IsuGlobalService: isuGlobalService,
	}
}

func (controller *IsuGlobalControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuGlobalRequest := isuglobal.IsuGlobalRequest{}
	helper.ReadFromRequestBody(request, &isuGlobalRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	isuGlobalResponse, err := controller.IsuGlobalService.Create(request.Context(), isuGlobalRequest)
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
		Status: "Success Created Isu Global",
		Data:   isuGlobalResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuGlobalControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuGlobalUpdateRequest := isuglobal.IsuGlobalUpdateRequest{}
	helper.ReadFromRequestBody(request, &isuGlobalUpdateRequest)

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
	isuGlobalUpdateRequest.ID = id

	isuGlobalResponse, err := controller.IsuGlobalService.Update(request.Context(), isuGlobalUpdateRequest)
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
		Status: "Success Updated Isu Global",
		Data:   isuGlobalResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuGlobalControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuGlobalId := params.ByName("id")
	id, err := strconv.Atoi(isuGlobalId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.IsuGlobalService.Delete(request.Context(), id)
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
		Status: "Success Deleted Isu Global",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}