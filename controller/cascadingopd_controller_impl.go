package controller

import (
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type CascadingOpdControllerImpl struct {
	CascadingOpdService service.CascadingOpdService
}

func NewCascadingOpdControllerImpl(cascadingOpdService service.CascadingOpdService) *CascadingOpdControllerImpl {
	return &CascadingOpdControllerImpl{CascadingOpdService: cascadingOpdService}
}

func (controller *CascadingOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	// Jika kodeOpd atau tahun kosong, kembalikan response null
	if kodeOpd == "" || tahun == "" {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service FindAll
	cascadingOpdResponse, err := controller.CascadingOpdService.FindAll(request.Context(), kodeOpd, tahun)
	if err != nil {
		// Jika tidak ada data, kembalikan response sukses dengan data null
		if err == sql.ErrNoRows {
			webResponse := web.WebResponse{
				Code:   200,
				Status: "OK",
				Data:   nil,
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Untuk error lainnya
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get All Cascading OPD",
		Data:   cascadingOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *CascadingOpdControllerImpl) FindByRekinPegawaiAndId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rekinId := params.ByName("rekin_id")

	cascadingOpdResponse, err := controller.CascadingOpdService.FindByRekinPegawaiAndId(request.Context(), rekinId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get By Rekin Pegawai And Id",
		Data:   cascadingOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CascadingOpdControllerImpl) FindByIdPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pokinId, err := strconv.Atoi(params.ByName("pokin_id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   "ID harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	cascadingOpdResponse, err := controller.CascadingOpdService.FindByIdPokin(request.Context(), pokinId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get By Id Pokin",
		Data:   cascadingOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CascadingOpdControllerImpl) FindByNip(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	nip := params.ByName("nip")
	tahun := params.ByName("tahun")

	cascadingOpdResponse, err := controller.CascadingOpdService.FindByNip(request.Context(), nip, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get By Nip",
		Data:   cascadingOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CascadingOpdControllerImpl) FindByMultipleRekinPegawai(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rekinRequest := pohonkinerja.FindByMultipleRekinRequest{}
	helper.ReadFromRequestBody(request, &rekinRequest)

	cascadingOpdResponse, err := controller.CascadingOpdService.FindByMultipleRekinPegawai(request.Context(), rekinRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get By Multiple Rekin Pegawai",
		Data:   cascadingOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
