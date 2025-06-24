package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CSFControllerImpl struct {
	CSFService service.CSFService
}

func NewCSFControllerImpl(csfService service.CSFService) CSFController {
	return &CSFControllerImpl{
		CSFService: csfService,
	}
}

func (controller *CSFControllerImpl) FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")

	if tahun == "" {
		helper.WriteToResponseBody(writer, "Tahun harus diisi")
		return
	}

	csfResponses, err := controller.CSFService.FindByTahun(request.Context(), tahun)
	if err != nil {
		helper.WriteToResponseBody(writer, err.Error())
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: http.StatusText(200),
		Data:   csfResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
