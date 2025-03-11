package controller

import (
	"context"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ReviewControllerImpl struct {
	ReviewService service.ReviewService
}

func NewReviewControllerImpl(reviewService service.ReviewService) *ReviewControllerImpl {
	return &ReviewControllerImpl{
		ReviewService: reviewService,
	}
}

func (controller *ReviewControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	reviewCreateRequest := pohonkinerja.ReviewCreateRequest{}
	helper.ReadFromRequestBody(request, &reviewCreateRequest)

	idStr := params.ByName("pokinId")
	id, err := strconv.Atoi(idStr)
	helper.PanicIfError(err)
	reviewCreateRequest.IdPohonKinerja = id

	claims := request.Context().Value(helper.UserInfoKey).(web.JWTClaim)
	ctx := context.WithValue(request.Context(), helper.UserInfoKey, claims)

	reviewResponse, err := controller.ReviewService.Create(ctx, reviewCreateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create review",
		Data:   reviewResponse,
	})
}

func (controller *ReviewControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	paramId := params.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
	}

	reviewUpdateRequest := pohonkinerja.ReviewUpdateRequest{}
	helper.ReadFromRequestBody(request, &reviewUpdateRequest)

	reviewUpdateRequest.Id = id

	reviewResponse, err := controller.ReviewService.Update(request.Context(), reviewUpdateRequest)
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
		Status: "success update review",
		Data:   reviewResponse,
	})
}

func (controller *ReviewControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	paramId := params.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
		return
	}

	controller.ReviewService.Delete(request.Context(), id)
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success delete review",
	})
}

func (controller *ReviewControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonkinerjaId := params.ByName("pokin_id")
	id, err := strconv.Atoi(pohonkinerjaId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
		return
	}

	reviewResponse, err := controller.ReviewService.FindAll(request.Context(), id)
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
		Status: "success get all review",
		Data:   reviewResponse,
	})
}

func (controller *ReviewControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	paramId := params.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
	}

	reviewResponse, err := controller.ReviewService.FindById(request.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success get review by id",
		Data:   reviewResponse,
	})
}

func (controller *ReviewControllerImpl) FindAllReviewByTematik(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")

	reviewResponse, err := controller.ReviewService.FindAllReviewByTematik(request.Context(), tahun)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: fmt.Sprintf("success get review tematik tahun %v", tahun),
		Data:   reviewResponse,
	})
}
