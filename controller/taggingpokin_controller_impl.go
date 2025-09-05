package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/taggingpokin"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TaggingPokinControllerImpl struct {
	TaggingPokinService service.TaggingPokinService
}

func NewTaggingPokinControllerImpl(taggingPokinService service.TaggingPokinService) *TaggingPokinControllerImpl {
	return &TaggingPokinControllerImpl{
		TaggingPokinService: taggingPokinService,
	}
}

func (controller *TaggingPokinControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taggingPokinCreateRequest := taggingpokin.TaggingPokinCreateRequest{}
	helper.ReadFromRequestBody(request, &taggingPokinCreateRequest)

	taggingPokinResponse, err := controller.TaggingPokinService.Create(request.Context(), taggingPokinCreateRequest)
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
		Status: "OK",
		Data:   taggingPokinResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TaggingPokinControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taggingPokinUpdateRequest := taggingpokin.TaggingPokinUpdateRequest{}
	helper.ReadFromRequestBody(request, &taggingPokinUpdateRequest)

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	taggingPokinUpdateRequest.Id = id

	taggingPokinResponse, err := controller.TaggingPokinService.Update(request.Context(), taggingPokinUpdateRequest)
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
		Status: "OK",
		Data:   taggingPokinResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TaggingPokinControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taggingPokinId := params.ByName("id")
	id, err := strconv.Atoi(taggingPokinId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	controller.TaggingPokinService.Delete(request.Context(), id)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TaggingPokinControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taggingPokinId := params.ByName("id")
	id, err := strconv.Atoi(taggingPokinId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	taggingPokinResponse, err := controller.TaggingPokinService.FindById(request.Context(), id)
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
		Status: "OK",
		Data:   taggingPokinResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TaggingPokinControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taggingPokinResponses, err := controller.TaggingPokinService.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   taggingPokinResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
