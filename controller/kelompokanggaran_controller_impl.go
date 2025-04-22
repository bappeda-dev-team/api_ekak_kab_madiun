package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/kelompokanggarans"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type KelompokAnggaranControllerImpl struct {
	kelompokanggaranService service.KelompokAnggaranService
}

func NewKelompokAnggaranControllerImpl(kelompokanggaranService service.KelompokAnggaranService) *KelompokAnggaranControllerImpl {
	return &KelompokAnggaranControllerImpl{
		kelompokanggaranService: kelompokanggaranService,
	}
}

func (controller *KelompokAnggaranControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	KelompokAnggaranCreateRequest := kelompokanggarans.KelompokAnggaranCreateRequest{}
	helper.ReadFromRequestBody(request, &KelompokAnggaranCreateRequest)

	// Panggil service untuk membuat kegiatan baru
	kelompokResponse, err := controller.kelompokanggaranService.Create(request.Context(), KelompokAnggaranCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   kelompokResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *KelompokAnggaranControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	KelompokAnggaranUpdateRequest := kelompokanggarans.KelompokAnggaranUpdateRequest{}
	helper.ReadFromRequestBody(request, &KelompokAnggaranUpdateRequest)

	// Ambil ID dari parameter URL
	idInt, _ := strconv.Atoi(params.ByName("id"))
	KelompokAnggaranUpdateRequest.Id = idInt

	// Panggil service untuk update kegiatan
	kelompokResponse, err := controller.kelompokanggaranService.Update(request.Context(), KelompokAnggaranUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   kelompokResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *KelompokAnggaranControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	Id := params.ByName("id")

	// Panggil service untuk delete kegiatan
	err := controller.kelompokanggaranService.Delete(request.Context(), Id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   nil,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *KelompokAnggaranControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kegiatanId := params.ByName("id")

	// Panggil service untuk mencari kegiatan berdasarkan ID
	kegiatanResponse, err := controller.kelompokanggaranService.FindById(request.Context(), kegiatanId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   kegiatanResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *KelompokAnggaranControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kelompokResponses, err := controller.kelompokanggaranService.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   kelompokResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
