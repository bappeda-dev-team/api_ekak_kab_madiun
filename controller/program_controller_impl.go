package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProgramControllerImpl struct {
	ProgramService service.ProgramService
}

func NewProgramControllerImpl(programService service.ProgramService) *ProgramControllerImpl {
	return &ProgramControllerImpl{
		ProgramService: programService,
	}
}

func (controller *ProgramControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programCreateRequest := programkegiatan.ProgramKegiatanCreateRequest{}
	helper.ReadFromRequestBody(request, &programCreateRequest)

	programResponse, err := controller.ProgramService.Create(request.Context(), programCreateRequest)
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
		Code:   201,
		Status: "Success",
		Data:   programResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUpdateRequest := programkegiatan.ProgramKegiatanUpdateRequest{}
	helper.ReadFromRequestBody(request, &programUpdateRequest)

	programUpdateRequest.Id = params.ByName("programId")
	programResponse, err := controller.ProgramService.Update(request.Context(), programUpdateRequest)
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
		Status: "Success",
		Data:   programResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programId := params.ByName("id")
	err := controller.ProgramService.Delete(request.Context(), programId)
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
		Status: "Success",
		Data:   "Program berhasil dihapus",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programId := params.ByName("id")
	programResponse, err := controller.ProgramService.FindById(request.Context(), programId)
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
		Data:   programResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programResponses, err := controller.ProgramService.FindAll(request.Context())
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
		Status: "Success",
		Data:   programResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// UpsertIndikator godoc
// @Summary      Upsert Indikator Program
// @Description  Menambah, mengubah, atau menghapus indikator (beserta target) untuk satu program dalam satu request.
// @Description  Aturan:
// @Description  - Indikator dengan `id` diisi → update indikator tersebut
// @Description  - Indikator dengan `id` kosong → insert sebagai indikator baru
// @Description  - Indikator existing yang TIDAK ada di array → dihapus beserta semua target-nya
// @Tags         Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        program_id  path  string  true  "ID Program (contoh: PRGM-X.XX.XX)"
// @Param        request     body  programkegiatan.UpsertIndikatorProgramRequest  true  "Array indikator"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.UpsertIndikatorProgramResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      404  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /program_kegiatan/indikator/:program_id [put]
func (controller *ProgramControllerImpl) UpsertIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programId := params.ByName("program_id")

	var req programkegiatan.UpsertIndikatorProgramRequest
	helper.ReadFromRequestBody(request, &req)

	resp, err := controller.ProgramService.UpsertIndikator(request.Context(), programId, req)
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
		Data:   resp,
	})
}

// FindAllWithIndikator godoc
// @Summary      Daftar Program dengan Indikator Master (V2)
// @Description  Mengambil semua program dari master beserta indikator dari tb_master_indikator_program_kegiatan (match kode = kode_program). Tanpa target; data indikator di-inject via database.
// @Tags         Program Kegiatan
// @Accept       json
// @Produce      json
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.ProgramKegiatanWithIndikatorV2Response}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /program_kegiatan/findall/v2 [get]
func (controller *ProgramControllerImpl) FindAllWithIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programResponses, err := controller.ProgramService.FindAllWithIndikator(request.Context())
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Internal Server Error",
			Data:   err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   programResponses,
	})
}
