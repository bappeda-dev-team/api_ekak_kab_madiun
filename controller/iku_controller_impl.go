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

// @Summary     Get IKU Renstra OPD
// @Description  Get IKU Renstra by kode OPD, tahun awal, tahun akhir, and jenis periode.
// @Tags         IKU Renstra
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param         tahun_awal path 	string true "Tahun Awal"
// @Param         tahun_akhir path 	string true "Tahun Akhir"
// @Param         jenis_periode path string true "Jenis Periode"
// @Success      200  {object}  web.WebResponse{data=[]iku.IkuOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /indikator_utama/opd/{kode_opd}/{tahun_awal}/{tahun_akhir}/{jenis_periode} [GET]
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

func (controller *IkuControllerImpl) UpdateIkuOpdActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ikuUpdateActiveRequest := iku.IkuUpdateActiveRequest{}
	helper.ReadFromRequestBody(request, &ikuUpdateActiveRequest)

	id := params.ByName("indikator_id")

	err := controller.IkuService.UpdateIkuOpdActive(request.Context(), id, ikuUpdateActiveRequest)
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
