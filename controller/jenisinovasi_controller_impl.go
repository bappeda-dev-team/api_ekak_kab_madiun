package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/jenisinovasi"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type JenisInovasiControllerImpl struct {
	JenisInovasiService service.JenisInovasiService
}

func NewJenisInovasiControllerImpl(jenisinovasiService service.JenisInovasiService) *JenisInovasiControllerImpl {
	return &JenisInovasiControllerImpl{
		JenisInovasiService: jenisinovasiService,
	}
}

func (controller *JenisInovasiControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	jiRequest := jenisinovasi.JenisInovasiRequest{}
	helper.ReadFromRequestBody(request, &jiRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	jenisResponse, err := controller.JenisInovasiService.Create(request.Context(), jiRequest)
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
		Status: "Success Created Jenis Inovasi",
		Data:   jenisResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JenisInovasiControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	UpdateRequest := jenisinovasi.JenisInovasiUpdateRequest{}
	helper.ReadFromRequestBody(request, &UpdateRequest)

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
	UpdateRequest.ID = id

	ppdResponse, err := controller.JenisInovasiService.Update(request.Context(), UpdateRequest)
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
		Status: "Success Updated Jenis Inovasi",
		Data:   ppdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JenisInovasiControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	err = controller.JenisInovasiService.Delete(request.Context(), id)
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
		Status: "Success Deleted Jenis Inovasi",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JenisInovasiControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	bidangUrusanResponses, err := controller.JenisInovasiService.FindAll(request.Context())
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
