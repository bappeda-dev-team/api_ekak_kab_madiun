package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/kegiatan"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type KegiatanControllerImpl struct {
	KegiatanService service.KegiatanService
}

func NewKegiatanControllerImpl(kegiatanService service.KegiatanService) *KegiatanControllerImpl {
	return &KegiatanControllerImpl{
		KegiatanService: kegiatanService,
	}
}

func (controller *KegiatanControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Decode request JSON ke struct kegiatan
	kegiatanCreateRequest := kegiatan.KegiatanCreateRequest{}
	helper.ReadFromRequestBody(request, &kegiatanCreateRequest)

	// Panggil service untuk membuat kegiatan baru
	kegiatanResponse, err := controller.KegiatanService.Create(request.Context(), kegiatanCreateRequest)
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
func (controller *KegiatanControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Decode request JSON ke struct kegiatan
	kegiatanUpdateRequest := kegiatan.KegiatanUpdateRequest{}
	helper.ReadFromRequestBody(request, &kegiatanUpdateRequest)

	// Ambil ID dari parameter URL
	kegiatanUpdateRequest.Id = params.ByName("id")

	// Panggil service untuk update kegiatan
	kegiatanResponse, err := controller.KegiatanService.Update(request.Context(), kegiatanUpdateRequest)
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

func (controller *KegiatanControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil ID dari parameter URL
	kegiatanId := params.ByName("id")

	// Panggil service untuk delete kegiatan
	err := controller.KegiatanService.Delete(request.Context(), kegiatanId)
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

func (controller *KegiatanControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil ID dari parameter URL
	kegiatanId := params.ByName("id")

	// Panggil service untuk mencari kegiatan berdasarkan ID
	kegiatanResponse, err := controller.KegiatanService.FindById(request.Context(), kegiatanId)
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

func (controller *KegiatanControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Panggil service untuk mendapatkan semua kegiatan
	kegiatanResponses, err := controller.KegiatanService.FindAll(request.Context())
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
		Data:   kegiatanResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// FindAllWithIndikator godoc
// @Summary      Daftar Kegiatan dengan Indikator Master (V2)
// @Description  Mengambil semua kegiatan dari master beserta indikator dari tb_master_indikator_program_kegiatan (match kode = kode_kegiatan). Tanpa target; data indikator di-inject via database.
// @Tags         Kegiatan
// @Accept       json
// @Produce      json
// @Success      200  {object}  web.WebResponse{data=[]kegiatan.KegiatanWithIndikatorV2Response}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /kegiatan/findall/v2 [get]
func (controller *KegiatanControllerImpl) FindAllWithIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kegiatanResponses, err := controller.KegiatanService.FindAllWithIndikator(request.Context())
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
		Status: "OK",
		Data:   kegiatanResponses,
	})
}
