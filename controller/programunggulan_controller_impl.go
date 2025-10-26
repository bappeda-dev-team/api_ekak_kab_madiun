package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programunggulan"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ProgramUnggulanControllerImpl struct {
	ProgramUnggulanService service.ProgramUnggulanService
}

func NewProgramUnggulanControllerImpl(programUnggulanService service.ProgramUnggulanService) *ProgramUnggulanControllerImpl {
	return &ProgramUnggulanControllerImpl{
		ProgramUnggulanService: programUnggulanService,
	}
}
func (controller *ProgramUnggulanControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanCreateRequest := programunggulan.ProgramUnggulanCreateRequest{}
	helper.ReadFromRequestBody(request, &programUnggulanCreateRequest)

	programUnggulanResponse, err := controller.ProgramUnggulanService.Create(request.Context(), programUnggulanCreateRequest)
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
		Status: "Success Created Program Unggulan",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramUnggulanControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanUpdateRequest := programunggulan.ProgramUnggulanUpdateRequest{}
	helper.ReadFromRequestBody(request, &programUnggulanUpdateRequest)

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
	programUnggulanUpdateRequest.Id = id

	programUnggulanResponse, err := controller.ProgramUnggulanService.Update(request.Context(), programUnggulanUpdateRequest)
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
		Status: "Success Updated Program Unggulan",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramUnggulanControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanId := params.ByName("id")
	id, err := strconv.Atoi(programUnggulanId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.ProgramUnggulanService.Delete(request.Context(), id)
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
		Status: "Success Deleted Program Unggulan",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *ProgramUnggulanControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanId := params.ByName("id")
	id, err := strconv.Atoi(programUnggulanId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	programUnggulanResponse, err := controller.ProgramUnggulanService.FindById(request.Context(), id)
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
		Status: "Success Found Program Unggulan",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *ProgramUnggulanControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanResponse, err := controller.ProgramUnggulanService.FindAll(request.Context(), params.ByName("tahun_awal"), params.ByName("tahun_akhir"))
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
		Status: "Success Found Program Unggulan",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *ProgramUnggulanControllerImpl) FindByKodeProgramUnggulan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanResponse, err := controller.ProgramUnggulanService.FindByKodeProgramUnggulan(request.Context(), params.ByName("kode_program_unggulan"))
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
		Status: "Success Found Program Unggulan",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *ProgramUnggulanControllerImpl) FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanResponse, err := controller.ProgramUnggulanService.FindByTahun(request.Context(), params.ByName("tahun"))
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
		Status: "Success Found Program Unggulan",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramUnggulanControllerImpl) FindUnusedByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	programUnggulanResponse, err := controller.ProgramUnggulanService.FindUnusedByTahun(request.Context(), params.ByName("tahun"))
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
		Status: "Success Found Program Unggulan",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramUnggulanControllerImpl) FindByIdTerkait(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	findByIdTerkaitRequest := programunggulan.FindByIdTerkaitRequest{}
	helper.ReadFromRequestBody(request, &findByIdTerkaitRequest)

	programUnggulanResponse, err := controller.ProgramUnggulanService.FindByIdTerkait(request.Context(), findByIdTerkaitRequest)
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
		Status: "Success Found Program Unggulan Terkait",
		Data:   programUnggulanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
