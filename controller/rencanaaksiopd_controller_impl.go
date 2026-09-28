package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/renaksiopd"
	"ekak_kabupaten_madiun/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type RencanaAksiOpdControllerImpl struct {
	RencanaAksiOpdService service.RencanaAksiOpdService
}

func NewRencanaAksiOpdControllerImpl(rencanaAksiOpdService service.RencanaAksiOpdService) *RencanaAksiOpdControllerImpl {
	return &RencanaAksiOpdControllerImpl{
		RencanaAksiOpdService: rencanaAksiOpdService,
	}
}

// FindBySasaranOpdAndTahun godoc
// @Summary      Daftar Rencana Aksi OPD per Sasaran
// @Description  Menampilkan seluruh rencana aksi pada satu sasaran OPD untuk tahun tertentu, lengkap dengan rencana kinerja, subkegiatan, dan indikatornya.
// @Tags         Renaksi Opd
// @Produce      json
// @Param        sasaran_opd_id  path      int     true  "ID Sasaran OPD"
// @Param        tahun           path      string  true  "Tahun"
// @Success      200             {object}  web.WebResponse{data=[]renaksiopd.RencanaAksiOpdResponse}
// @Failure      400             {object}  web.WebResponse
// @Failure      500             {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /rencana-aksi-opd/{sasaran_opd_id}/{tahun} [get]
func (controller *RencanaAksiOpdControllerImpl) FindBySasaranOpdAndTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	sasaranOpdId := params.ByName("sasaran_opd_id")
	sasaranOpdIdInt, err := strconv.Atoi(sasaranOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	tahun := params.ByName("tahun")
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.FindBySasaranOpdAndTahun(request.Context(), sasaranOpdIdInt, tahun)
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
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// SyncJadwalPelaksanaan godoc
// @Summary      Sinkronisasi Jadwal Pelaksanaan
// @Description  Menghitung ulang bobot triwulan (TW1-TW4) dari rencana aksi lalu menyimpannya pada data renaksi OPD untuk rekin terkait.
// @Tags         Renaksi Opd
// @Produce      json
// @Param        rekin_id  path      string  true  "ID Rencana Kinerja"
// @Success      200       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /rencana-aksi-opd/sync_jadwal/{rekin_id} [post]
func (controller *RencanaAksiOpdControllerImpl) SyncJadwalPelaksanaan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rekinId := params.ByName("rekin_id")
	err := controller.RencanaAksiOpdService.SyncJadwalPelaksanaan(request.Context(), rekinId)
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
		Data:   "Jadwal pelaksanaan berhasil disinkronkan",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// Create godoc
// @Summary      Tambah Rencana Aksi OPD
// @Description  Menautkan satu rencana kinerja ke satu sasaran OPD. Satu rekin hanya boleh terdaftar satu kali dalam satu sasaran OPD, sehingga permintaan yang bentrok akan gagal dengan 409.
// @Tags         Renaksi Opd
// @Accept       json
// @Produce      json
// @Param        request   body      renaksiopd.RencanaAksiOpdCreateRequest  true  "Data rencana aksi OPD"
// @Success      200       {object}  web.WebResponse{data=renaksiopd.RencanaAksiOpdRequestResponse}
// @Failure      400       {object}  web.WebResponse
// @Failure      409       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /rencana-aksi-opd/create [post]
func (controller *RencanaAksiOpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaAksiOpdCreateRequest := renaksiopd.RencanaAksiOpdCreateRequest{}
	err := json.NewDecoder(request.Body).Decode(&rencanaAksiOpdCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.Create(request.Context(), rencanaAksiOpdCreateRequest)
	if err != nil {
		controller.writeMutationError(writer, err)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// Update godoc
// @Summary      Ubah Rencana Aksi OPD
// @Description  Mengganti rencana kinerja dan keterangan pada satu rencana aksi OPD. Sasaran OPD dan tahun tidak dapat diubah. Duplikasi rekin pada sasaran yang sama akan gagal dengan 409.
// @Tags         Renaksi Opd
// @Accept       json
// @Produce      json
// @Param        id        path      int                                       true  "ID Rencana Aksi OPD"
// @Param        request   body      renaksiopd.RencanaAksiOpdUpdateRequest   true  "Data rencana aksi OPD"
// @Success      200       {object}  web.WebResponse{data=renaksiopd.RencanaAksiOpdRequestResponse}
// @Failure      400       {object}  web.WebResponse
// @Failure      409       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /rencana-aksi-opd/update/{id} [put]
func (controller *RencanaAksiOpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaAksiOpdUpdateRequest := renaksiopd.RencanaAksiOpdUpdateRequest{}
	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdUpdateRequest.Id = idInt

	err = json.NewDecoder(request.Body).Decode(&rencanaAksiOpdUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.Update(request.Context(), rencanaAksiOpdUpdateRequest)
	if err != nil {
		controller.writeMutationError(writer, err)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// Delete godoc
// @Summary      Hapus Rencana Aksi OPD
// @Description  Menghapus satu rencana aksi OPD berdasarkan ID.
// @Tags         Renaksi Opd
// @Produce      json
// @Param        id        path      int     true  "ID Rencana Aksi OPD"
// @Success      200       {object}  web.WebResponse
// @Failure      400       {object}  web.WebResponse
// @Failure      409       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /rencana-aksi-opd/delete/{id} [delete]
func (controller *RencanaAksiOpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	err = controller.RencanaAksiOpdService.Delete(request.Context(), idInt)
	if err != nil {
		controller.writeMutationError(writer, err)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Rencana aksi opd berhasil dihapus",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindById godoc
// @Summary      Detail Rencana Aksi OPD
// @Description  Menampilkan detail satu rencana aksi OPD beserta data sasaran OPD dan indikatornya.
// @Tags         Renaksi Opd
// @Produce      json
// @Param        id        path      int                                          true  "ID Rencana Aksi OPD"
// @Success      200       {object}  web.WebResponse{data=renaksiopd.RencanaAksiOpdByIdResponse}
// @Failure      400       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /renaksi-opd/detail/{id} [get]
func (controller *RencanaAksiOpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.FindById(request.Context(), idInt)
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
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindAllSasaranByTahun godoc
// @Summary      Daftar Sasaran OPD untuk Rencana Aksi
// @Description  Menampilkan daftar sasaran OPD milik satu OPD pada tahun tertentu yang dapat dipilih untuk ditambahkan rencana aksi.
// @Tags         Renaksi Opd
// @Produce      json
// @Param        kode_opd   path      string  true  "Kode OPD"
// @Param        tahun      path      string  true  "Tahun"
// @Success      200        {object}  web.WebResponse{data=[]renaksiopd.SasaranOpdDetailResponse}
// @Failure      500        {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/all/{kode_opd}/{tahun} [get]
func (controller *RencanaAksiOpdControllerImpl) FindAllSasaranByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	sasaranList, err := controller.RencanaAksiOpdService.FindAllSasaranByTahun(request.Context(), kodeOpd, tahun)
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
		Data:   sasaranList,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiOpdControllerImpl) writeMutationError(writer http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if errors.Is(err, service.ErrRencanaAksiOpdLocked) || errors.Is(err, service.ErrRekinSudahDigunakan) {
		code = http.StatusConflict
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   code,
		Status: http.StatusText(code),
		Data:   err.Error(),
	})
}
