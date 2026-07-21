package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/ppd"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PpdControllerImpl struct {
	PpdService service.PpdService
}

func NewPpdControllerImpl(ppdService service.PpdService) *PpdControllerImpl {
	return &PpdControllerImpl{
		PpdService: ppdService,
	}
}

func (controller *PpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ppdRequest := ppd.PpdRequest{}
	helper.ReadFromRequestBody(request, &ppdRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	ppdResponse, err := controller.PpdService.Create(request.Context(), ppdRequest)
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
		Status: "Success Created Potensi Perangkat Daerah",
		Data:   ppdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ppdUpdateRequest := ppd.PpdUpdateRequest{}
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

	ppdResponse, err := controller.PpdService.Update(request.Context(), ppdUpdateRequest)
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
		Status: "Success Updated Potensi Perangkat Daerah",
		Data:   ppdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	err = controller.PpdService.Delete(request.Context(), id)
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
		Status: "Success Deleted Potensi Perangkat Daerah",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	kodeOpd := params.ByName("kode_opd")

	bidangUrusanResponses, err := controller.PpdService.FindAll(request.Context(), kodeOpd)
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