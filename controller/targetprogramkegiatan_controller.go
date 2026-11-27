package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TargetProgramKegiatanController — CRUD endpoint untuk tb_target_program_kegiatan.
// Create dan Update dipisah per jenis (renstra / renja / rankhir / penetapan).
// Update hanya mengizinkan perubahan target dan satuan.
type TargetProgramKegiatanController interface {
	// Create per jenis
	CreateRenstra(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	CreateRenja(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	CreateRankhir(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	CreatePenetapan(w http.ResponseWriter, r *http.Request, p httprouter.Params)

	// Update per jenis (hanya target & satuan)
	UpdateRenstra(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	UpdateRenja(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	UpdateRankhir(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	UpdatePenetapan(w http.ResponseWriter, r *http.Request, p httprouter.Params)

	// Shared
	Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	FindById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	FindByKodeIndikatorAndKodeOpd(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}
