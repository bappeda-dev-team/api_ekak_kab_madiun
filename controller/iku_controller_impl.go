package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/iku"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IkuControllerImpl struct {
	IkuService service.IkuService
}

func NewIkuControllerImpl(ikuService service.IkuService) *IkuControllerImpl {
	return &IkuControllerImpl{
		IkuService: ikuService,
	}
}

func (controller *IkuControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")

	if tahunAwal == "" {
		// Handle error jika tahun tidak ada
		helper.WriteToResponseBody(writer, "Tahun harus diisi")
		return
	}

	ikuResponses, err := controller.IkuService.FindAll(request.Context(), tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		helper.WriteToResponseBody(writer, err.Error())
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   ikuResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *IkuControllerImpl) FindAllIkuOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")

	ikuOpdResponses, err := controller.IkuService.FindAllIkuOpd(request.Context(), kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)

	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   ikuOpdResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)

}
func (controller *IkuControllerImpl) UpdateIkuActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikuUpdateActiveRequest := iku.IkuUpdateActiveRequest{}
	helper.ReadFromRequestBody(request, &ikuUpdateActiveRequest)

	id := params.ByName("indikator_id")

	err := controller.IkuService.UpdateIkuActive(request.Context(), id, ikuUpdateActiveRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Berhasil mengupdate status IKU",
	})
}
