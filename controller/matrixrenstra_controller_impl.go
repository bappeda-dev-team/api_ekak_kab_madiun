package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenstraControllerImpl struct {
	MatrixRenstraService service.MatrixRenstraService
}

func NewMatrixRenstraControllerImpl(matrixRenstraService service.MatrixRenstraService) *MatrixRenstraControllerImpl {
	return &MatrixRenstraControllerImpl{MatrixRenstraService: matrixRenstraService}
}

// @Summary      Matrix Renstra
// @Description  Mendapatkan data matrix renstra berdasarkan kode OPD dan tahun.
// @Tags         Matrix Renstra
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun_awal  query     string  true  "Tahun Awal"      example("2025")
// @Param        tahun_akhir  query     string  true  "Tahun Akhir"      example("2026")
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.UrusanDetailResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renstra/opd/{kode_opd} [get]
func (controller *MatrixRenstraControllerImpl) GetByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahunAwal := request.URL.Query().Get("tahun_awal")
	tahunAkhir := request.URL.Query().Get("tahun_akhir")

	matrixRenstraResponses, err := controller.MatrixRenstraService.GetByKodeSubKegiatan(request.Context(), kodeOpd, tahunAwal, tahunAkhir)
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
		Status: "OK",
		Data:   matrixRenstraResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Delete Indikator Renstra
// @Description  Menghapus data indikator renstra yang sudah ada berdasarkan Kode Indikator.
// @Tags         Matrix Renstra
// @Accept       json
// @Produce      json
// @Param        kode_indikator  path     string  true  "Kode Indikator"   example("RENS-1.01.1.01.0.00.01.0000-2025-01")
// @Success      200  {object}  web.WebResponse{data=string}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renstra/indikator/delete/{kode_indikator} [delete]
func (controller *MatrixRenstraControllerImpl) DeleteIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kode_indikator")

	err := controller.MatrixRenstraService.DeleteIndikator(request.Context(), kodeIndikator)
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
		Status: "success delete indikator",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Upsert Anggaran Renstra
// @Description  Upsert anggaran renstra.
// @Tags         Matrix Renstra
// @Accept       json
// @Produce      json
// @Param        request  body  programkegiatan.AnggaranRenstraRequest  true  "Anggaran Renstra"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.AnggaranRenstraResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renstra/anggaran/upsert [post]
func (controller *MatrixRenstraControllerImpl) UpsertAnggaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	AnggaranRenstraRequest := programkegiatan.AnggaranRenstraRequest{}
	helper.ReadFromRequestBody(request, &AnggaranRenstraRequest)

	AnggaranRenstraResponse, err := controller.MatrixRenstraService.UpsertAnggaran(request.Context(), AnggaranRenstraRequest)
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
		Status: "success upsert anggaran renstra",
		Data:   AnggaranRenstraResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

// @Summary      Upsert Batch Indikator Renstra
// @Description  kode_indikator kosong = create baru, isi = update
// @Tags         Matrix Renstra
// @Accept       json
// @Produce      json
// @Param        request  body  []programkegiatan.IndikatorRenstraCreateRequest  true  "Array of indikator"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.IndikatorResponse}
// @Security     BearerAuth
// @Router       /matrix_renstra/indikator/upsert [post]
func (controller *MatrixRenstraControllerImpl) UpsertBatchIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var requests []programkegiatan.IndikatorRenstraCreateRequest
	helper.ReadFromRequestBody(request, &requests)
	resp, err := controller.MatrixRenstraService.UpsertBatchIndikator(request.Context(), requests)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success upsert indikator renstra", Data: resp,
	})
}
