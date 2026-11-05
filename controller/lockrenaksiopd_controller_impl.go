package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/lockrenaksiopd"
	"ekak_kabupaten_madiun/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type LockRenaksiOpdControllerImpl struct {
	LockRenaksiOpdService service.LockRenaksiOpdService
}

func NewLockRenaksiOpdControllerImpl(lockRenaksiOpdService service.LockRenaksiOpdService) *LockRenaksiOpdControllerImpl {
	return &LockRenaksiOpdControllerImpl{LockRenaksiOpdService: lockRenaksiOpdService}
}

// Lock godoc
// @Summary      Lock Renaksi OPD
// @Description  Membuat snapshot lock Renaksi OPD berdasarkan kode OPD, tahun, sasaran, dan rekin.
// @Tags         Lock Renaksi OPD
// @Accept       json
// @Produce      json
// @Param        kode_opd  path      string                                  true  "Kode OPD"
// @Param        tahun     path      string                                  true  "Tahun"
// @Param        request   body      lockrenaksiopd.LockRenaksiOpdRequest    true  "Data lock Renaksi OPD"
// @Success      200       {object}  web.WebResponse{data=lockrenaksiopd.LockRenaksiOpdResponse}
// @Failure      400       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /lock-renaksi-opd/lock/{kode_opd}/{tahun} [post]
func (controller *LockRenaksiOpdControllerImpl) Lock(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	lockRequest := lockrenaksiopd.LockRenaksiOpdRequest{}
	if err := json.NewDecoder(request.Body).Decode(&lockRequest); err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	response, err := controller.LockRenaksiOpdService.Lock(request.Context(), params.ByName("kode_opd"), params.ByName("tahun"), lockRequest)
	if err != nil {
		writeLockRenaksiOpdError(writer, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: response})
}

// Unlock godoc
// @Summary      Unlock Renaksi OPD
// @Description  Menghapus snapshot lock Renaksi OPD berdasarkan kode OPD, tahun, dan ID Renaksi OPD.
// @Tags         Lock Renaksi OPD
// @Produce      json
// @Param        kode_opd  path      string  true  "Kode OPD"
// @Param        tahun     path      string  true  "Tahun"
// @Param        id        path      int     true  "ID Renaksi OPD"
// @Success      200       {object}  web.WebResponse
// @Failure      400       {object}  web.WebResponse
// @Failure      404       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /lock-renaksi-opd/lock/{kode_opd}/{tahun}/{id} [delete]
func (controller *LockRenaksiOpdControllerImpl) Unlock(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	if err := controller.LockRenaksiOpdService.Unlock(request.Context(), params.ByName("kode_opd"), params.ByName("tahun"), id); err != nil {
		writeLockRenaksiOpdError(writer, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: "Lock Renaksi OPD berhasil dibuka"})
}

// FindById godoc
// @Summary      Detail Lock Renaksi OPD
// @Description  Menampilkan detail lock Renaksi OPD berdasarkan kode OPD, tahun, dan ID Renaksi OPD.
// @Tags         Lock Renaksi OPD
// @Produce      json
// @Param        kode_opd  path      string  true  "Kode OPD"
// @Param        tahun     path      string  true  "Tahun"
// @Param        id        path      int     true  "ID Renaksi OPD"
// @Success      200       {object}  web.WebResponse{data=lockrenaksiopd.LockRenaksiOpdResponse}
// @Failure      400       {object}  web.WebResponse
// @Failure      404       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /lock-renaksi-opd/lock/{kode_opd}/{tahun}/{id} [get]
func (controller *LockRenaksiOpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	response, err := controller.LockRenaksiOpdService.FindById(request.Context(), params.ByName("kode_opd"), params.ByName("tahun"), id)
	if err != nil {
		writeLockRenaksiOpdError(writer, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: response})
}

// FindAll godoc
// @Summary      Daftar Lock Renaksi OPD
// @Description  Menampilkan daftar lock Renaksi OPD berdasarkan kode OPD dan tahun.
// @Tags         Lock Renaksi OPD
// @Produce      json
// @Param        kode_opd  path      string  true  "Kode OPD"
// @Param        tahun     path      string  true  "Tahun"
// @Success      200       {object}  web.WebResponse{data=[]lockrenaksiopd.LockRenaksiOpdResponse}
// @Failure      400       {object}  web.WebResponse
// @Failure      500       {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /lock-renaksi-opd/lock/{kode_opd}/{tahun} [get]
func (controller *LockRenaksiOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	response, err := controller.LockRenaksiOpdService.FindAll(request.Context(), params.ByName("kode_opd"), params.ByName("tahun"))
	if err != nil {
		writeLockRenaksiOpdError(writer, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: response})
}

func writeLockRenaksiOpdError(writer http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if errors.Is(err, service.ErrLockRenaksiOpdInvalidParameter) {
		code = http.StatusBadRequest
	}
	if errors.Is(err, service.ErrLockRenaksiOpdNotFound) {
		code = http.StatusNotFound
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   code,
		Status: http.StatusText(code),
		Data:   err.Error(),
	})
}
