package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/rincianbelanja"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RincianBelanjaControllerImpl struct {
	rincianBelanjaService service.RincianBelanjaService
}

func NewRincianBelanjaControllerImpl(rincianBelanjaService service.RincianBelanjaService) *RincianBelanjaControllerImpl {
	return &RincianBelanjaControllerImpl{
		rincianBelanjaService: rincianBelanjaService,
	}
}

func (controller *RincianBelanjaControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rincianBelanjaCreateRequest := rincianbelanja.RincianBelanjaCreateRequest{}
	helper.ReadFromRequestBody(request, &rincianBelanjaCreateRequest)

	rincianBelanjaResponse, err := controller.rincianBelanjaService.Create(request.Context(), rincianBelanjaCreateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create rincian belanja",
		Data:   rincianBelanjaResponse,
	})
}

func (controller *RincianBelanjaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rincianBelanjaUpdateRequest := rincianbelanja.RincianBelanjaUpdateRequest{}
	helper.ReadFromRequestBody(request, &rincianBelanjaUpdateRequest)

	rincianBelanjaUpdateRequest.RenaksiId = params.ByName("renaksiId")
	rincianBelanjaResponse, err := controller.rincianBelanjaService.Update(request.Context(), rincianBelanjaUpdateRequest)
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
		Status: "Success Update Rincian Belanja",
		Data:   rincianBelanjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RincianBelanjaControllerImpl) FindRincianBelanjaAsn(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiId := params.ByName("pegawai_id")
	tahun := params.ByName("tahun")

	response := controller.rincianBelanjaService.FindRincianBelanjaAsn(request.Context(), pegawaiId, tahun)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success",
		Data:   response,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
