package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/renaksiopd"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
)

type RencanaAksiOpdServiceImpl struct {
	RencanaAksiOpdRepository repository.RencanaAksiOpdRepository
	DB                       *sql.DB
	validator                *validator.Validate
}

func NewRencanaAksiOpdServiceImpl(rencanaAksiOpdRepository repository.RencanaAksiOpdRepository, db *sql.DB, validator *validator.Validate) *RencanaAksiOpdServiceImpl {
	return &RencanaAksiOpdServiceImpl{
		RencanaAksiOpdRepository: rencanaAksiOpdRepository,
		DB:                       db,
		validator:                validator,
	}
}

func (service *RencanaAksiOpdServiceImpl) FindBySasaranOpdAndTahun(ctx context.Context, sasaranOpdId int, tahun string) ([]renaksiopd.RencanaAksiOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	rencanaAksi, err := service.RencanaAksiOpdRepository.FindBySasaranOpdAndTahun(ctx, tx, sasaranOpdId, tahun)
	if err != nil {
		return nil, err
	}

	return toRencanaAksiOpdResponses(rencanaAksi), nil
}

func (service *RencanaAksiOpdServiceImpl) SyncJadwalPelaksanaan(ctx context.Context, rekinId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.RencanaAksiOpdRepository.SyncJadwalPelaksanaan(ctx, tx, rekinId)
}

func (service *RencanaAksiOpdServiceImpl) Create(ctx context.Context, request renaksiopd.RencanaAksiOpdCreateRequest) (renaksiopd.RencanaAksiOpdRequestResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		return renaksiopd.RencanaAksiOpdRequestResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return renaksiopd.RencanaAksiOpdRequestResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	var keterangan *string
	if request.Keterangan != "" {
		keterangan = &request.Keterangan
	}

	rencanaAksiOpdDomain := domain.RencanaAksiOpd{
		Id:           rand.Intn(10000000),
		RekinId:      request.RekinId,
		SasaranOpdId: request.SasaranOpdId,
		TahunRenaksi: request.TahunRenaksi,
		Keterangan:   keterangan,
	}

	rencanaAksiOpd, err := service.RencanaAksiOpdRepository.Create(ctx, tx, rencanaAksiOpdDomain)
	if err != nil {
		return renaksiopd.RencanaAksiOpdRequestResponse{}, err
	}

	return toRencanaAksiOpdRequestResponse(rencanaAksiOpd), nil
}

func (service *RencanaAksiOpdServiceImpl) Update(ctx context.Context, request renaksiopd.RencanaAksiOpdUpdateRequest) (renaksiopd.RencanaAksiOpdRequestResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		return renaksiopd.RencanaAksiOpdRequestResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return renaksiopd.RencanaAksiOpdRequestResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	var keterangan *string
	if request.Keterangan != "" {
		keterangan = &request.Keterangan
	}

	rencanaAksiOpdDomain := domain.RencanaAksiOpd{
		Id:         request.Id,
		RekinId:    request.RekinId,
		Keterangan: keterangan,
	}

	rencanaAksiOpd := service.RencanaAksiOpdRepository.Update(ctx, tx, rencanaAksiOpdDomain)

	return toRencanaAksiOpdRequestResponse(rencanaAksiOpd), nil
}

func (service *RencanaAksiOpdServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.RencanaAksiOpdRepository.Delete(ctx, tx, id)
}

func (service *RencanaAksiOpdServiceImpl) FindById(ctx context.Context, id int) (renaksiopd.RencanaAksiOpdByIdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return renaksiopd.RencanaAksiOpdByIdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	rencanaAksiOpd, err := service.RencanaAksiOpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return renaksiopd.RencanaAksiOpdByIdResponse{}, err
	}

	return toRencanaAksiOpdByIdResponse(rencanaAksiOpd), nil
}

func (service *RencanaAksiOpdServiceImpl) FindAllSasaranByTahun(ctx context.Context, kodeOpd string, tahun string) ([]renaksiopd.SasaranOpdDetailResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranList, err := service.RencanaAksiOpdRepository.FindAllSasaranByTahun(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	responses := make([]renaksiopd.SasaranOpdDetailResponse, 0, len(sasaranList))
	for _, sasaran := range sasaranList {
		responses = append(responses, toSasaranOpdDetailResponse(sasaran))
	}

	return responses, nil
}

func toRencanaAksiOpdResponses(items []domain.RencanaAksiOpd) []renaksiopd.RencanaAksiOpdResponse {
	if len(items) == 0 {
		return []renaksiopd.RencanaAksiOpdResponse{}
	}

	grouped := make(map[string]*renaksiopd.RencanaAksiOpdResponse)
	for _, item := range items {
		key := fmt.Sprintf("%d-%s", item.SasaranOpdId, item.TahunRenaksi)
		resp, exists := grouped[key]
		if !exists {
			grouped[key] = &renaksiopd.RencanaAksiOpdResponse{
				SasaranOpdId:   item.SasaranOpdId,
				NamaSasaranOpd: item.NamaSasaranOpd,
				TahunRenaksi:   item.TahunRenaksi,
				RencanaKinerja: []renaksiopd.RencanaKinerjaResponse{},
			}
			resp = grouped[key]
		}

		for _, rk := range item.RencanaKinerja {
			resp.RencanaKinerja = append(resp.RencanaKinerja, toRencanaKinerjaResponse(rk))
		}
	}

	responses := make([]renaksiopd.RencanaAksiOpdResponse, 0, len(grouped))
	for _, resp := range grouped {
		responses = append(responses, *resp)
	}

	return responses
}

func toRencanaKinerjaResponse(rk domain.RencanaKinerjaOpd) renaksiopd.RencanaKinerjaResponse {
	subKegiatan := make([]renaksiopd.SubKegiatanResponse, 0, len(rk.SubKegiatan))
	for _, sk := range rk.SubKegiatan {
		indikators := make([]renaksiopd.IndikatorResponse, 0, len(sk.Indikator))
		for _, ind := range sk.Indikator {
			indikators = append(indikators, renaksiopd.IndikatorResponse{
				Id:        ind.Id,
				Indikator: ind.Indikator,
				Target:    ind.Target,
				Satuan:    ind.Satuan,
			})
		}

		subKegiatan = append(subKegiatan, renaksiopd.SubKegiatanResponse{
			KodeSubKegiatan: sk.KodeSubKegiatan,
			NamaSubKegiatan: sk.NamaSubKegiatan,
			Indikator:       indikators,
		})
	}

	return renaksiopd.RencanaKinerjaResponse{
		Id:                 rk.Id,
		RekinId:            rk.RekinId,
		NamaRencanaKinerja: rk.NamaRencanaKinerja,
		NipPegawai:         rk.NipPegawai,
		NamaPegawai:        rk.NamaPegawai,
		KodeOpd:            rk.KodeOpd,
		TotalAnggaran:      rk.TotalAnggaran,
		Tw1:                rk.Tw1,
		Tw2:                rk.Tw2,
		Tw3:                rk.Tw3,
		Tw4:                rk.Tw4,
		Keterangan:         rk.Keterangan,
		SubKegiatan:        subKegiatan,
	}
}

func toRencanaAksiOpdRequestResponse(item domain.RencanaAksiOpd) renaksiopd.RencanaAksiOpdRequestResponse {
	return renaksiopd.RencanaAksiOpdRequestResponse{
		SasaranOpdId: item.SasaranOpdId,
		RekinId:      item.RekinId,
		TahunRenaksi: item.TahunRenaksi,
		Tw1:          item.Tw1,
		Tw2:          item.Tw2,
		Tw3:          item.Tw3,
		Tw4:          item.Tw4,
		Keterangan:   item.Keterangan,
	}
}

func toRencanaAksiOpdByIdResponse(item domain.RencanaAksiOpd) renaksiopd.RencanaAksiOpdByIdResponse {
	return renaksiopd.RencanaAksiOpdByIdResponse{
		Id:                 item.Id,
		RekinId:            item.RekinId,
		TahunRenaksi:       item.TahunRenaksi,
		Keterangan:         item.Keterangan,
		NamaRencanaKinerja: item.NamaRencanaKinerja,
		SasaranOpd:         toSasaranOpdDetailResponse(item.SasaranOpd),
	}
}

func toSasaranOpdDetailResponse(item domain.SasaranOpdDetailRenaksi) renaksiopd.SasaranOpdDetailResponse {
	indikators := make([]renaksiopd.IndikatorSasaranOpdResponse, 0, len(item.Indikator))
	for _, ind := range item.Indikator {
		indikators = append(indikators, renaksiopd.IndikatorSasaranOpdResponse{
			Id:               ind.Id,
			Indikator:        ind.Indikator,
			RumusPerhitungan: ind.RumusPerhitungan,
			SumberData:       ind.SumberData,
			Target: renaksiopd.TargetResponse{
				Id:          ind.Target.Id,
				IndikatorId: ind.Target.IndikatorId,
				Tahun:       ind.Target.Tahun,
				Target:      ind.Target.Target,
				Satuan:      ind.Target.Satuan,
			},
		})
	}

	return renaksiopd.SasaranOpdDetailResponse{
		Id:             item.Id,
		NamaSasaranOpd: item.NamaSasaranOpd,
		TahunAwal:      item.TahunAwal,
		TahunAkhir:     item.TahunAkhir,
		JenisPeriode:   item.JenisPeriode,
		Indikator:      indikators,
	}
}
