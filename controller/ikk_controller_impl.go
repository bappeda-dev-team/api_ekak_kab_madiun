package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/ikk"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type IkkControllerImpl struct {
	IkkService service.IkkService
}

func NewIkkControllerImpl(ikkService service.IkkService) *IkkControllerImpl {
	return &IkkControllerImpl{
		IkkService: ikkService,
	}
}

func (controller *IkkControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikkRequest := ikk.IkkRequest{}
	helper.ReadFromRequestBody(request, &ikkRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	ikkResponse, err := controller.IkkService.Create(request.Context(), ikkRequest)
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
		Status: "Success Created IKK",
		Data:   ikkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IkkControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikkUpdateRequest := ikk.IkkUpdateRequest{}
	helper.ReadFromRequestBody(request, &ikkUpdateRequest)

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
	ikkUpdateRequest.ID = id

	ikkResponse, err := controller.IkkService.Update(request.Context(), ikkUpdateRequest)
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
		Status: "Success Updated IKK",
		Data:   ikkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IkkControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikkId := params.ByName("id")
	id, err := strconv.Atoi(ikkId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.IkkService.Delete(request.Context(), id)
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
		Status: "Success Deleted IKK",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *IkkControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikkId := params.ByName("id")
	id, err := strconv.Atoi(ikkId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	ikkResponse, err := controller.IkkService.FindById(request.Context(), id)
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
		Status: "Success Found IKK",
		Data:   ikkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *IkkControllerImpl) FindByKodeOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	levelStr := params.ByName("level_pohon")

	levelPohon, err := strconv.Atoi(levelStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "level_pohon harus berupa angka", // 👈 lebih jelas dari err.Error()
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	kodeOpd := params.ByName("kode_opd")

	bidangUrusanResponses, err := controller.IkkService.FindByKodeOpd(request.Context(), levelPohon, kodeOpd)
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