package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/tujuanopd"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TujuanOpdControllerImpl struct {
	TujuanOpdService service.TujuanOpdService
}

func NewTujuanOpdControllerImpl(tujuanOpdService service.TujuanOpdService) *TujuanOpdControllerImpl {
	return &TujuanOpdControllerImpl{
		TujuanOpdService: tujuanOpdService,
	}
}

// CreateTujuanOpd godoc
// @Summary      Tambah Tujuan Opd Renstra
// @Description  Memasukkan data tujuan opd renstra baru ke dalam sistem.
// @Tags         Tujuan Opd Renstra
// @Accept       json
// @Produce      json
// @Param        request  body      tujuanopd.TujuanOpdCreateRequest  true  "Payload Create Tujuan OPD"
// @Success      201      {object}  web.WebResponse{data=tujuanopd.TujuanOpdResponse}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renstra/create [post]
func (controller *TujuanOpdControllerImpl) CreateTujuanOpdRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanOpdCreateRequest := tujuanopd.TujuanOpdCreateRequest{}
	helper.ReadFromRequestBody(request, &tujuanOpdCreateRequest)

	for i := range tujuanOpdCreateRequest.Indikator {
		tujuanOpdCreateRequest.Indikator[i].Jenis = "renstra"
	}

	// Panggil service Create
	tujuanOpdResponse, err := controller.TujuanOpdService.Create(request.Context(), tujuanOpdCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create tujuan opd",
		Data:   tujuanOpdResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// UpdateTujuanOpd godoc
// @Summary      Update Tujuan Opd Renstra
// @Description  Memperbarui data tujuan opd renstra yang sudah ada berdasarkan ID.
// @Tags         Tujuan Opd Renstra
// @Accept       json
// @Produce      json
// @Param        tujuanOpdId       path      int                              true  "ID Tujuan OPD" example(1)
// @Param        request  body      tujuanopd.TujuanOpdUpdateRequest  true  "Payload Update Tujuan OPD"
// @Success      200      {object}  web.WebResponse{data=tujuanopd.TujuanOpdResponse}
// @Failure      400      {object}  web.WebResponse
// @Failure      404      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renstra/update/{tujuanOpdId} [put]
func (controller *TujuanOpdControllerImpl) UpdateTujuanOpdRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil ID dari params
	tujuanOpdId := params.ByName("tujuanOpdId")
	tujuanOpdIdInt, err := strconv.Atoi(tujuanOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
	}

	// Baca request body
	tujuanOpdUpdateRequest := tujuanopd.TujuanOpdUpdateRequest{}
	helper.ReadFromRequestBody(request, &tujuanOpdUpdateRequest)

	for i := range tujuanOpdUpdateRequest.Indikator {
		tujuanOpdUpdateRequest.Indikator[i].Jenis = "renstra"
	}

	// Set ID dari params ke request
	tujuanOpdUpdateRequest.Id = tujuanOpdIdInt

	// Panggil service Update
	tujuanOpdResponse, err := controller.TujuanOpdService.Update(request.Context(), tujuanOpdUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tujuanOpdResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanOpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanOpdId := params.ByName("tujuanOpdId")
	tujuanOpdIdInt, err := strconv.Atoi(tujuanOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.TujuanOpdService.Delete(request.Context(), tujuanOpdIdInt)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Data berhasil dihapus",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanOpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanOpdId := params.ByName("tujuanOpdId")
	tujuanOpdIdInt, err := strconv.Atoi(tujuanOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	tujuanOpdResponse, err := controller.TujuanOpdService.FindById(request.Context(), tujuanOpdIdInt)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find by id tujuan opd",
		Data:   tujuanOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")

	tujuanOpdResponses, err := controller.TujuanOpdService.FindAll(request.Context(), kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find all tujuan opd",
		Data:   tujuanOpdResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanOpdControllerImpl) FindTujuanOpdOnlyName(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")

	tujuanOpdResponses, err := controller.TujuanOpdService.FindTujuanOpdOnlyName(request.Context(), kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find tujuan opd",
		Data:   tujuanOpdResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanOpdControllerImpl) FindTujuanOpdByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")

	tujuanOpdResponses, err := controller.TujuanOpdService.FindTujuanOpdByTahun(request.Context(), kodeOpd, tahun, jenisPeriode)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find tujuan opd by tahun",
		Data:   tujuanOpdResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// GetRenjaRanwal godoc
// @Summary      Tujuan Opd Renstra
// @Description  Mendapatkan data tujuan opd renstra berdasarkan kode OPD dan tahun.
// @Tags         Tujuan Opd Renstra
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun_awal     path     string  true  "Tahun Awal"      example("2025")
// @Param        tahun_akhir     path     string  true  "Tahun Akhir"      example("2030")
// @Success      200  {object}  web.WebResponse{data=[]tujuanopd.TujuanOpdwithBidangUrusanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renstra/{kode_opd}/{tahun_awal}/{tahun_akhir} [get]
func (controller *TujuanOpdControllerImpl) FindTujuanOpdRenstra(
	writer http.ResponseWriter,
	request *http.Request,
	params httprouter.Params,
) {
	kodeOpd := params.ByName("kode_opd")
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	// jenisPeriode := params.ByName("jenis_periode") // dari URL — nilai di tb_tujuan_opd
	tujuanOpdResponses, err := controller.TujuanOpdService.FindTujuanRenstra(
		request.Context(), kodeOpd, tahunAwal, tahunAkhir, "RPJMD",
	)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   tujuanOpdResponses,
	})
}

// GetRenjaRanwal godoc
// @Summary      Tujuan Opd Renja Ranwal
// @Description  Mendapatkan data tujuan opd ranwal berdasarkan kode OPD dan tahun.
// @Tags         Tujuan Opd Renja Ranwal
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]tujuanopd.TujuanOpdwithBidangUrusanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/ranwal/{kode_opd}/{tahun} [get]
func (controller *TujuanOpdControllerImpl) FindTujuanOpdRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	tujuanOpdResponses, err := controller.TujuanOpdService.FindTujuanRanwal(request.Context(), kodeOpd, tahun, "RPJMD")
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   tujuanOpdResponses,
	})
}

// GetRenjaRankhir godoc
// @Summary      Tujuan Opd Renja Rankhir
// @Description  Mendapatkan data tujuan opd rankhir berdasarkan kode OPD dan tahun.
// @Tags         Tujuan Opd Renja Rankhir
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]tujuanopd.TujuanOpdwithBidangUrusanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/rankhir/{kode_opd}/{tahun} [get]
func (controller *TujuanOpdControllerImpl) FindTujuanOpdRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisPeriode := "RPJMD"                                                   // dari URL — nilai di tb_tujuan_opd
	tujuanOpdResponses, err := controller.TujuanOpdService.FindTujuanRankhir( // method BARU
		request.Context(), kodeOpd, tahun, jenisPeriode,
	)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   tujuanOpdResponses,
	})
}

// Create Indikator TujuanOpd godoc
// @Summary      Tambah Indikator Tujuan Opd Renja Ranwal
// @Description  Memasukkan data indikator tujuan opd renja ranwal baru ke dalam sistem.
// @Tags         Tujuan Opd Renja Ranwal
// @Accept       json
// @Produce      json
// @Param        tujuanOpdId       path      int                              true  "ID Tujuan OPD" example(1)
// @Param        request  body      []tujuanopd.IndikatorCreateRequest  true  "Payload Create Indikator Tujuan OPD"
// @Success      201      {object}  web.WebResponse{data=tujuanopd.IndikatorResponse}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renja/ranwal/indikator/create/{tujuanOpdId} [post]
func (controller *TujuanOpdControllerImpl) CreateTujuanRenjaRanwalIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanOpdId := params.ByName("tujuanOpdId")
	tujuanOpdIdInt, err := strconv.Atoi(tujuanOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	indikatorCreateRequests := []tujuanopd.IndikatorCreateRequest{}
	helper.ReadFromRequestBody(request, &indikatorCreateRequests)
	indikatorResponses, err := controller.TujuanOpdService.CreateTujuanRenjaIndikator(request.Context(), tujuanOpdIdInt, "ranwal", indikatorCreateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success create tujuan renja ranwal indikator",
		Data:   indikatorResponses,
	})
}

// Update Indikator TujuanOpd godoc
// @Summary      Update Indikator Tujuan Opd Renja Ranwal
// @Description  Memperbarui data indikator tujuan opd renja ranwal yang sudah ada berdasarkan ID.
// @Tags         Tujuan Opd Renja Ranwal
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Param        request  body      tujuanopd.IndikatorUpdateRequest  true  "Payload Update Indikator Tujuan OPD"
// @Success      201      {object}  web.WebResponse{data=tujuanopd.IndikatorResponse}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renja/ranwal/indikator/update/{kodeIndikator} [put]
func (controller *TujuanOpdControllerImpl) UpdateTujuanRenjaRanwalIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kodeIndikator")

	indikatorUpdateRequests := tujuanopd.IndikatorUpdateRequest{}
	helper.ReadFromRequestBody(request, &indikatorUpdateRequests)
	indikatorResponses, err := controller.TujuanOpdService.UpdateTujuanRenjaIndikator(request.Context(), kodeIndikator, "ranwal", indikatorUpdateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update tujuan renja ranwal indikator",
		Data:   indikatorResponses,
	})
}

// Create Indikator TujuanOpd Rankhir godoc
// @Summary      Tambah Indikator Tujuan Opd Renja Rankhir
// @Description  Memasukkan data indikator tujuan opd renja rankhir baru ke dalam sistem.
// @Tags         Tujuan Opd Renja Rankhir
// @Accept       json
// @Produce      json
// @Param        tujuanOpdId       path      int                              true  "ID Tujuan OPD" example(1)
// @Param        request  body      tujuanopd.IndikatorCreateRequest  true  "Payload Create Indikator Tujuan OPD"
// @Success      201      {object}  web.WebResponse{data=tujuanopd.IndikatorResponse}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renja/rankhir/indikator/create/{tujuanOpdId} [post]
func (controller *TujuanOpdControllerImpl) CreateTujuanRenjaRankhirIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanOpdId := params.ByName("tujuanOpdId")
	tujuanOpdIdInt, err := strconv.Atoi(tujuanOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	indikatorCreateRequests := []tujuanopd.IndikatorCreateRequest{}
	helper.ReadFromRequestBody(request, &indikatorCreateRequests)
	indikatorResponses, err := controller.TujuanOpdService.CreateTujuanRenjaIndikator(request.Context(), tujuanOpdIdInt, "rankhir", indikatorCreateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success create tujuan renja rankhir indikator",
		Data:   indikatorResponses,
	})
}

// Update Indikator TujuanOpd godoc
// @Summary      Update Indikator Tujuan Opd Renja Rankhir
// @Description  Memperbarui data indikator tujuan opd renja rankhir yang sudah ada berdasarkan ID.
// @Tags         Tujuan Opd Renja Rankhir
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Param        request  body      tujuanopd.IndikatorUpdateRequest  true  "Payload Update Indikator Tujuan OPD"
// @Success      201      {object}  web.WebResponse{data=tujuanopd.IndikatorResponse}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renja/rankhir/indikator/update/{kodeIndikator} [put]
func (controller *TujuanOpdControllerImpl) UpdateTujuanRenjaRankhirIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kodeIndikator")

	indikatorUpdateRequests := tujuanopd.IndikatorUpdateRequest{}
	helper.ReadFromRequestBody(request, &indikatorUpdateRequests)
	indikatorResponses, err := controller.TujuanOpdService.UpdateTujuanRenjaIndikator(request.Context(), kodeIndikator, "rankhir", indikatorUpdateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update tujuan renja ranwal indikator",
		Data:   indikatorResponses,
	})
}

// Delete Indikator TujuanOpd godoc
// @Summary      Delete Indikator Tujuan Opd Renja
// @Description  Menghapus data indikator tujuan opd renja yang sudah ada berdasarkan Kode Indikator.
// @Tags         Tujuan Opd Renja Rankhir
// @Tags         Tujuan Opd Renja Ranwal
// @Tags         Tujuan Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Success      200      {object}  web.WebResponse{data=string}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renja/indikator/delete/{kodeIndikator} [delete]
func (controller *TujuanOpdControllerImpl) DeleteTujuanRenjaIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	kodeIndikator := params.ByName("kodeIndikator")
	err := controller.TujuanOpdService.DeleteTujuanRenjaIndikator(request.Context(), kodeIndikator)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success delete tujuan renja indikator",
		Data:   nil,
	})

}

// GetRenjaRanwal godoc
// @Summary      Tujuan Opd Renja Penetapan
// @Description  Mendapatkan data tujuan opd penetapan berdasarkan kode OPD dan tahun.
// @Tags         Tujuan Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]tujuanopd.TujuanOpdwithBidangUrusanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/penetapan/{kode_opd}/{tahun} [get]
func (controller *TujuanOpdControllerImpl) FindTujuanOpdPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	tujuanOpdResponses, err := controller.TujuanOpdService.FindTujuanPenetapan(request.Context(), kodeOpd, tahun, "RPJMD")
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   tujuanOpdResponses,
	})
}

// Create Indikator TujuanOpd godoc
// @Summary      Tambah Indikator Tujuan Opd Renja Penetapan
// @Description  Memasukkan data indikator tujuan opd renja penetapan baru ke dalam sistem.
// @Tags         Tujuan Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        tujuanOpdId       path      int                              true  "ID Tujuan OPD" example(1)
// @Param        request  body      []tujuanopd.IndikatorCreateRequest  true  "Payload Create Indikator Tujuan OPD"
// @Success      201      {object}  web.WebResponse{data=tujuanopd.IndikatorResponse}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renja/penetapan/indikator/create/{tujuanOpdId} [post]
func (controller *TujuanOpdControllerImpl) CreateTujuanRenjaPenetapanIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanOpdId := params.ByName("tujuanOpdId")
	tujuanOpdIdInt, err := strconv.Atoi(tujuanOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	indikatorCreateRequests := []tujuanopd.IndikatorCreateRequest{}
	helper.ReadFromRequestBody(request, &indikatorCreateRequests)
	indikatorResponses, err := controller.TujuanOpdService.CreateTujuanRenjaIndikator(request.Context(), tujuanOpdIdInt, "penetapan", indikatorCreateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success create tujuan renja penetapan indikator",
		Data:   indikatorResponses,
	})
}

// Update Indikator TujuanOpd godoc
// @Summary      Update Indikator Tujuan Opd Renja Penetapan
// @Description  Memperbarui data indikator tujuan opd renja penetapan yang sudah ada berdasarkan ID.
// @Tags         Tujuan Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Param        request  body      tujuanopd.IndikatorUpdateRequest  true  "Payload Update Indikator Tujuan OPD"
// @Success      201      {object}  web.WebResponse{data=tujuanopd.IndikatorResponse}
// @Failure      400      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_opd/renja/penetapan/indikator/update/{kodeIndikator} [put]
func (controller *TujuanOpdControllerImpl) UpdateTujuanRenjaPenetapanIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kodeIndikator")

	indikatorUpdateRequests := tujuanopd.IndikatorUpdateRequest{}
	helper.ReadFromRequestBody(request, &indikatorUpdateRequests)
	indikatorResponses, err := controller.TujuanOpdService.UpdateTujuanRenjaIndikator(request.Context(), kodeIndikator, "penetapan", indikatorUpdateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update tujuan renja penetapan indikator",
		Data:   indikatorResponses,
	})
}
