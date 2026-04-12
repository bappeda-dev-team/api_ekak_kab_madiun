package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programprioritaspusat"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ProgramPrioritasPusatControllerImpl struct {
	ProgramPrioritasPusatService service.ProgramPrioritasPusatService
}

func NewProgramPrioritasPusatControllerImpl(programPrioritasPusatService service.ProgramPrioritasPusatService) *ProgramPrioritasPusatControllerImpl {
	return &ProgramPrioritasPusatControllerImpl{
		ProgramPrioritasPusatService: programPrioritasPusatService,
	}
}
func (controller *ProgramPrioritasPusatControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programPrioritasPusatCreateRequest := programprioritaspusat.ProgramPrioritasPusatCreateRequest{}
	helper.ReadFromRequestBody(request, &programPrioritasPusatCreateRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.Create(request.Context(), programPrioritasPusatCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			// TODO: CODE: AMBIL DARI http
			Code: http.StatusInternalServerError,
			// TODO: STATUS: TERJEMAHKAN DARI code
			Status: http.StatusText(http.StatusInternalServerError),
			// TODO: buat nil saja
			Data: err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		// TODO: CODE AMBIL DARI http
		Code:   201,
		Status: "Success Created Program Prioritas Pusat",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramPrioritasPusatControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programPrioritasPusatUpdateRequest := programprioritaspusat.ProgramPrioritasPusatUpdateRequest{}
	helper.ReadFromRequestBody(request, &programPrioritasPusatUpdateRequest)

	idStr := params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	programPrioritasPusatUpdateRequest.Id = id

	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.Update(request.Context(), programPrioritasPusatUpdateRequest)
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
		Status: "Success Updated Program Prioritas Pusat",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramPrioritasPusatControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programPrioritasPusatId := params.ByName("id")
	id, err := strconv.Atoi(programPrioritasPusatId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.ProgramPrioritasPusatService.Delete(request.Context(), id)
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
		Status: "Success Deleted Program Prioritas Pusat",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *ProgramPrioritasPusatControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programPrioritasPusatId := params.ByName("id")
	id, err := strconv.Atoi(programPrioritasPusatId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.FindById(request.Context(), id)
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
		Status: "Success Found Program Prioritas Pusat",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *ProgramPrioritasPusatControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// TODO: guard jika request invalid
	// return 400 Invalid
	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.FindAll(request.Context(), params.ByName("tahun_awal"), params.ByName("tahun_akhir"))
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
		Status: "Success Found Program Prioritas Pusat",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *ProgramPrioritasPusatControllerImpl) FindByKodeProgramPrioritasPusat(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// TODO: guard jika request invalid
	// return 400 Invalid
	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.FindByKodeProgramPrioritasPusat(request.Context(), params.ByName("kode_program_prioritas_pusat"))
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
		Status: "Success Found Program Prioritas Pusat",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *ProgramPrioritasPusatControllerImpl) FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// TODO: guard jika request invalid
	// return 400 Invalid
	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.FindByTahun(request.Context(), params.ByName("tahun"))
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
		Status: "Success Found Program Prioritas Pusat",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramPrioritasPusatControllerImpl) FindUnusedByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.FindUnusedByTahun(request.Context(), params.ByName("tahun"))
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
		Status: "Success Found Program Prioritas Pusat",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramPrioritasPusatControllerImpl) FindByIdTerkait(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	findByIdTerkaitRequest := programprioritaspusat.FindByIdTerkaitRequest{}
	helper.ReadFromRequestBody(request, &findByIdTerkaitRequest)

	programPrioritasPusatResponse, err := controller.ProgramPrioritasPusatService.FindByIdTerkait(request.Context(), findByIdTerkaitRequest)
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
		Status: "Success Found Program Prioritas Pusat Terkait",
		Data:   programPrioritasPusatResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
