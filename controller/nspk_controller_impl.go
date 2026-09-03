package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/masternspk"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type NspkControllerImpl struct {
	NspkService service.NspkService
}

func NewNspkControllerImpl(nspkService service.NspkService) *NspkControllerImpl {
	return &NspkControllerImpl{
		NspkService: nspkService,
	}
}

func (controller *NspkControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ppdRequest := masternspk.NspkRequest{}
	helper.ReadFromRequestBody(request, &ppdRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	ppdResponse, err := controller.NspkService.Create(request.Context(), ppdRequest)
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
		Status: "Success Created NSPK",
		Data:   ppdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *NspkControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ppdUpdateRequest := masternspk.NspkUpdateRequest{}
	helper.ReadFromRequestBody(request, &ppdUpdateRequest)

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
	ppdUpdateRequest.ID = id

	ppdResponse, err := controller.NspkService.Update(request.Context(), ppdUpdateRequest)
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
		Status: "Success Updated NSPK",
		Data:   ppdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *NspkControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ppdId := params.ByName("id")
	id, err := strconv.Atoi(ppdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.NspkService.Delete(request.Context(), id)
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
		Status: "Success Deleted NSPK",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *NspkControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	kodeOpd := params.ByName("kode_opd")

	bidangUrusanResponses, err := controller.NspkService.FindAll(request.Context(), kodeOpd)
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
