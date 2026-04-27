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

type PohonKinerjaOpdControllerImpl struct {
	PohonKinerjaOpdService service.PohonKinerjaOpdService
}

func NewPohonKinerjaOpdControllerImpl(pohonKinerjaOpdService service.PohonKinerjaOpdService) *PohonKinerjaOpdControllerImpl {
	return &PohonKinerjaOpdControllerImpl{
		PohonKinerjaOpdService: pohonKinerjaOpdService,
	}
}

func (controller *PohonKinerjaOpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonKinerjaCreateRequest := pohonkinerja.PohonKinerjaCreateRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaCreateRequest)

	response, err := controller.PohonKinerjaOpdService.Create(request.Context(), pohonKinerjaCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   response,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonKinerjaUpdateRequest := pohonkinerja.PohonKinerjaUpdateRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaUpdateRequest)

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	pohonKinerjaUpdateRequest.Id = id

	// Panggil service Update
	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.Update(request.Context(), pohonKinerjaUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Update Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idStr := params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.PohonKinerjaOpdService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Delete Pohon Kinerja",
		Data:   nil,
	}

	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *PohonKinerjaOpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Dapatkan id dari params
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service FindById
	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindById(request.Context(), id)
	if err != nil {
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
		Status: "Berhasil Mendapatkan Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindAll(request.Context(), kodeOpd, tahun)
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
		Status: "Success Get All Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *PohonKinerjaOpdControllerImpl) FindAllArah(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindAllArah(request.Context(), kodeOpd, tahun)
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
		Status: "Success Get All Strategic Arah Kebijakan",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
// func (controller *PohonKinerjaOpdControllerImpl) FindAllArahPemda(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
// 	kodeOpd := params.ByName("kode_opd")
// 	tahun := params.ByName("tahun")

// 	// Jika kodeOpd atau tahun kosong, kembalikan response null
// 	if kodeOpd == "" || tahun == "" {
// 		webResponse := web.WebResponse{
// 			Code:   200,
// 			Status: "OK",
// 			Data:   nil,
// 		}
// 		helper.WriteToResponseBody(writer, webResponse)
// 		return
// 	}

// 	// Panggil service FindAll
// 	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindAllArahPemda(request.Context(), kodeOpd, tahun)
// 	if err != nil {
// 		// Jika tidak ada data, kembalikan response sukses dengan data null
// 		if err == sql.ErrNoRows {
// 			webResponse := web.WebResponse{
// 				Code:   200,
// 				Status: "OK",
// 				Data:   nil,
// 			}
// 			helper.WriteToResponseBody(writer, webResponse)
// 			return
// 		}

// 		// Untuk error lainnya
// 		webResponse := web.WebResponse{
// 			Code:   404,
// 			Status: "Not Found",
// 			Data:   err.Error(),
// 		}
// 		writer.WriteHeader(http.StatusNotFound)
// 		helper.WriteToResponseBody(writer, webResponse)
// 		return
// 	}

// 	// Kirim response sukses
// 	webResponse := web.WebResponse{
// 		Code:   200,
// 		Status: "Success Get All Strategic Arah Kebijakan",
// 		Data:   pohonKinerjaResponse,
// 	}
// 	helper.WriteToResponseBody(writer, webResponse)
// }

func (controller *PohonKinerjaOpdControllerImpl) FindStrategicNoParent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	if kodeOpd == "" || tahun == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "kode_opd dan tahun harus diisi",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindStrategicNoParent(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get Strategic No Parent",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) DeletePelaksana(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := controller.PohonKinerjaOpdService.DeletePelaksana(request.Context(), params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusInternalServerError)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Delete Pelaksana",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) FindPokinByPelaksana(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pelaksanaId := params.ByName("pegawai_id")
	tahun := params.ByName("tahun")

	// Validasi parameter
	if pelaksanaId == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID Pegawai harus diisi",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	if tahun == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "Parameter tahun harus diisi",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindPokinByPelaksana(request.Context(), pelaksanaId, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Berhasil Mendapatkan Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) DeletePokinPemdaInOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pokinId := params.ByName("id")
	pokinIdInt, err := strconv.Atoi(pokinId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.PohonKinerjaOpdService.DeletePokinPemdaInOpd(request.Context(), pokinIdInt)
	if err != nil {

		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Delete Pohon Kinerja",
		Data:   nil,
	}

	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *PohonKinerjaOpdControllerImpl) UpdateParent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonKinerjaUpdateParentRequest := pohonkinerja.PohonKinerjaUpdateParentRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaUpdateParentRequest)

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	pohonKinerjaUpdateParentRequest.Id = id

	_, err = controller.PohonKinerjaOpdService.UpdateParent(request.Context(), pohonKinerjaUpdateParentRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Update Parent",
		Data:   "berhasil mengupdate parent",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) FindidPokinWithAllTema(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindidPokinWithAllTema(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get Pokin With Ancestors",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Clone Pohon Kinerja opd
// @Description  Melakukan cloning pohon kinerja dari tahun sumber ke tahun tujuan untuk OPD tertentu.
// @Tags         Clone Pohon Kinerja Opd
// @Accept       json
// @Produce      json
// @Param        pohon_kinerja_clone_request  body     pohonkinerja.PohonKinerjaCloneRequest  true  "Data untuk cloning pohon kinerja"
// @Success      200  {object}  web.WebResponse{data=string}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /pohon_kinerja_opd/clone [post]
func (controller *PohonKinerjaOpdControllerImpl) Clone(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cloneRequest := pohonkinerja.PohonKinerjaCloneRequest{}
	helper.ReadFromRequestBody(request, &cloneRequest)

	err := controller.PohonKinerjaOpdService.CloneByKodeOpdAndTahun(request.Context(), cloneRequest)
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
		Data:   "Berhasil cloning pohon kinerja dari tahun " + cloneRequest.TahunSumber + " ke tahun " + cloneRequest.TahunTujuan + " untuk OPD " + cloneRequest.KodeOpd,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) CheckPokinExistsByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	exists, err := controller.PohonKinerjaOpdService.CheckPokinExistsByTahun(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success",
		Data:   exists,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) CountPokinPemda(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	countPokinPemda, err := controller.PohonKinerjaOpdService.CountPokinPemda(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success",
		Data:   countPokinPemda,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) FindPokinAtasan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindPokinAtasan(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) ControlPokinOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	controlPokinOpd, err := controller.PohonKinerjaOpdService.ControlPokinOpd(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get Control Pokin Opd",
		Data:   controlPokinOpd,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) LeaderboardPokinOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")

	response, err := controller.PohonKinerjaOpdService.LeaderboardPokinOpd(request.Context(), tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   response,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Find All Pokin Parent Clone Pokin Opd
// @Description  Mendapatkan semua pohon kinerja parent clone pohon kinerja opd berdasarkan kode OPD dan tahun.
// @Tags         Clone Pohon Kinerja Opd
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Param        level_pohon  path     string  true  "Level Pohon"  example("1")
// @Success      200  {object}  web.WebResponse{data=[]pohonkinerja.PohonKinerjaOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      404  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /pohon_kinerja_opd/pokin_clone_pokin_opd_statistik/{kode_opd}/{tahun}/{level_pohon} [get]
func (controller *PohonKinerjaOpdControllerImpl) FindAllPokinParentClonePokinOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	levelPohon := params.ByName("level_pohon")

	levelPohonInt, err := strconv.Atoi(levelPohon)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "Level Pohon harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindAllPokinParentClonePokinOpd(request.Context(), kodeOpd, tahun, &levelPohonInt)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get All Pokin Parent Clone Pokin Opd",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Update Parent Pohon Kinerja Opd Clone
// @Description  Mengupdate parent pohon kinerja opd clone berdasarkan ID.
// @Tags         Clone Pohon Kinerja Opd
// @Accept       json
// @Produce      json
// @Param        id  path     string  true  "ID Pohon Kinerja"  example("1")
// @Param        pohon_kinerja_update_parent_clone_request  body     pohonkinerja.PohonKinerjaUpdateParentRequest  true  "Data untuk mengupdate parent pohon kinerja"
// @Success      200  {object}  web.WebResponse{data=string}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /pohon_kinerja_opd/update_parent_clone/{id} [put]
func (controller *PohonKinerjaOpdControllerImpl) UpdateParentClone(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonKinerjaUpdateParentCloneRequest := pohonkinerja.PohonKinerjaUpdateParentRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaUpdateParentCloneRequest)

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	pohonKinerjaUpdateParentCloneRequest.Id = id

	updateCloneResponse, err := controller.PohonKinerjaOpdService.UpdateParentClone(request.Context(), pohonKinerjaUpdateParentCloneRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Update Parent",
		Data:   updateCloneResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
