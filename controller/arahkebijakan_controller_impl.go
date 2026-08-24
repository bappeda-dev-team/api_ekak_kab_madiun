package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/arahkebijakan"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ArahKebijakanControllerImpl struct {
	ArahKebijakanService service.ArahKebijakanService
}

func NewArahKebijakanControllerImpl(arahkebijakanService service.ArahKebijakanService) *ArahKebijakanControllerImpl {
	return &ArahKebijakanControllerImpl{
		ArahKebijakanService: arahkebijakanService,
	}
}

func (controller *ArahKebijakanControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ppdRequest := arahkebijakan.ArahKebijakanRequest{}
	helper.ReadFromRequestBody(request, &ppdRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	ppdResponse, err := controller.ArahKebijakanService.Create(request.Context(), ppdRequest)
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
		Status: "Success Created Arah Kebijakan",
		Data:   ppdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ArahKebijakanControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ppdUpdateRequest := arahkebijakan.ArahKebijakanUpdateRequest{}
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

	ppdResponse, err := controller.ArahKebijakanService.Update(request.Context(), ppdUpdateRequest)
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
		Status: "Success Updated Arah Kebijakan",
		Data:   ppdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
