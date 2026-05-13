package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/ikd"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

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

func (controller *IkdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikdRequest := ikd.ProgramOpdTerpilihCreateRequest{}
	helper.ReadFromRequestBody(request, &ikdRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	ikkResponse, err := controller.IkdService.Create(request.Context(), ikdRequest)
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
		Status: "Success Select Program Opd",
		Data:   ikkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IkdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikdId := params.ByName("id")
	id, err := strconv.Atoi(ikdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.IkdService.Delete(request.Context(), id)
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
		Status: "Success Deleted Program Opd",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}