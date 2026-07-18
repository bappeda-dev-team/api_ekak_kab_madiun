package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/isuklhs"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type IsuKlhsControllerImpl struct {
	IsuKlhsService service.IsuKlhsService
}

func NewIsuKlhsControllerImpl(isuKlhsService service.IsuKlhsService) *IsuKlhsControllerImpl {
	return &IsuKlhsControllerImpl{
		IsuKlhsService: isuKlhsService,
	}
}

func (controller *IsuKlhsControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuKlhsRequest := isuklhs.IsuKlhsRequest{}
	helper.ReadFromRequestBody(request, &isuKlhsRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	isuKlhsResponse, err := controller.IsuKlhsService.Create(request.Context(), isuKlhsRequest)
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
		Status: "Success Created Isu KLHS",
		Data:   isuKlhsResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuKlhsControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuKlhsUpdateRequest := isuklhs.IsuKlhsUpdateRequest{}
	helper.ReadFromRequestBody(request, &isuKlhsUpdateRequest)

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
	isuKlhsUpdateRequest.ID = id

	isuKlhsResponse, err := controller.IsuKlhsService.Update(request.Context(), isuKlhsUpdateRequest)
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
		Status: "Success Updated Isu KLHS",
		Data:   isuKlhsResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IsuKlhsControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	isuKlhsId := params.ByName("id")
	id, err := strconv.Atoi(isuKlhsId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.IsuKlhsService.Delete(request.Context(), id)
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
		Status: "Success Deleted Isu KLHS",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}