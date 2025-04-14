package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/user"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserControllerImpl(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		userService: userService,
	}
}

func (controller *UserControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userCreateRequest := user.UserCreateRequest{}
	helper.ReadFromRequestBody(request, &userCreateRequest)

	userResponse, err := controller.userService.Create(request.Context(), userCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed create user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create user",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Parse ID dari URL parameter
	userId := params.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed update user",
			Data:   "invalid user id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	userUpdateRequest := user.UserUpdateRequest{}
	helper.ReadFromRequestBody(request, &userUpdateRequest)

	userUpdateRequest.Id = id

	userResponse, err := controller.userService.Update(request.Context(), userUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed update user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success update user",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId := params.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed delete user",
			Data:   "invalid user id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	controller.userService.Delete(request.Context(), id)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success delete user",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := request.URL.Query().Get("kode_opd")
	userResponses, err := controller.userService.FindAll(request.Context(), kodeOpd)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find all user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find all user",
		Data:   userResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId := params.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find by id user",
			Data:   "invalid user id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	userResponse, err := controller.userService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find by id user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find by id user",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	loginRequest := user.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &loginRequest)

	loginResponse, err := controller.userService.Login(request.Context(), loginRequest)
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
		Data: map[string]interface{}{
			"token": loginResponse.Token,
		},
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindByKodeOpdAndRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := request.URL.Query().Get("kode_opd")
	roleName := request.URL.Query().Get("role")

	userResponses, err := controller.userService.FindByKodeOpdAndRole(request.Context(), kodeOpd, roleName)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find by kode opd and role",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find by kode opd and role",
		Data:   userResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
