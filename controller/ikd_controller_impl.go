package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IkdControllerImpl struct {
	IkdService service.IkdService
}

func NewIkdControllerImpl(ikdService service.IkdService) *IkdControllerImpl {
	return &IkdControllerImpl{
		IkdService: ikdService,
	}
}

func (controller *IkdControllerImpl) FindAll(
	writer http.ResponseWriter,
	request *http.Request,
	params httprouter.Params,
) {

	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")

	responses, err := controller.IkdService.FindAll(
		request.Context(),
		kodeOpd,
		tahun,
		jenisPeriode,
	)

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
		Data:   responses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}