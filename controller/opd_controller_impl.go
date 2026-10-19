package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type OpdControllerImpl struct {
	OpdService service.OpdService
}

func NewOpdControllerImpl(opdService service.OpdService) *OpdControllerImpl {
	return &OpdControllerImpl{
		OpdService: opdService,
	}
}

// Create - Membuat data OPD baru
func (controller *OpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	opdCreateRequest := opdmaster.OpdCreateRequest{}
	helper.ReadFromRequestBody(request, &opdCreateRequest)

	opdResponse, err := controller.OpdService.Create(request.Context(), opdCreateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   500,
			Status: "error",
			Data:   err.Error(),
		})
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success create data opd",
		Data:   opdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// Update - Mengupdate data OPD yang ada
func (controller *OpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	opdUpdateRequest := opdmaster.OpdUpdateRequest{}
	helper.ReadFromRequestBody(request, &opdUpdateRequest)

	opdId := params.ByName("opdId")
	opdUpdateRequest.Id = opdId

	opdResponse, err := controller.OpdService.Update(request.Context(), opdUpdateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, err)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success update data opd",
		Data:   opdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// Delete - Menghapus data OPD
func (controller *OpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	opdId := params.ByName("opdId")

	err := controller.OpdService.Delete(request.Context(), opdId)
	if err != nil {
		helper.WriteToResponseBody(writer, err)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success delete data opd",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindById - Mencari OPD berdasarkan ID
func (controller *OpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	opdId := params.ByName("opdId")

	opdResponse, err := controller.OpdService.FindById(request.Context(), opdId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   500,
			Status: "error",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   200,
		Status: "success",
		Data:   opdResponse,
	})
}

// FindAll - Mengambil semua data OPD
func (controller *OpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	opdResponses, err := controller.OpdService.FindAll(request.Context())
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   500,
			Status: "error",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   200,
		Status: "success",
		Data:   opdResponses,
	})
}
