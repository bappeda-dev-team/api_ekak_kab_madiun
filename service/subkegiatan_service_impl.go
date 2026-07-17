package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubKegiatanServiceImpl struct {
	subKegiatanRepository   repository.SubKegiatanRepository
	opdRepository           repository.OpdRepository
	rencanaKinerjaRepoitory repository.RencanaKinerjaRepository
	DB                      *sql.DB
	validator               *validator.Validate
}

func NewSubKegiatanServiceImpl(subKegiatanRepository repository.SubKegiatanRepository, opdRepository repository.OpdRepository, rencanaKinerjaRepoitory repository.RencanaKinerjaRepository, DB *sql.DB, validator *validator.Validate) *SubKegiatanServiceImpl {
	return &SubKegiatanServiceImpl{
		subKegiatanRepository:   subKegiatanRepository,
		opdRepository:           opdRepository,
		rencanaKinerjaRepoitory: rencanaKinerjaRepoitory,
		DB:                      DB,
		validator:               validator,
	}
}

func (service *SubKegiatanServiceImpl) Create(ctx context.Context, request subkegiatan.SubKegiatanCreateRequest) (subkegiatan.SubKegiatanResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		log.Println("Validasi gagal:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	uuId := fmt.Sprintf("SUB-KEG-%s", request.KodeSubkegiatan)

	var indikators []domain.Indikator

	for _, indikatorReq := range request.Indikator {
		indikatorId := indikatorReq.Id
		if indikatorId == "" {
			randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			indikatorId = fmt.Sprintf("IND-SUB-%s", randomDigits)
		}

		var targets []domain.Target

		for _, targetReq := range indikatorReq.Target {
			targetId := targetReq.Id
			if targetId == "" {
				randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
				targetId = fmt.Sprintf("TRG-SUB-%s", randomDigits)
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      targetReq.TargetIndikator,
				Satuan:      targetReq.SatuanIndikator,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:            indikatorId,
			SubKegiatanId: uuId,
			Indikator:     indikatorReq.NamaIndikator,
			Target:        targets,
		}
		indikators = append(indikators, indikator)
	}

	subKegiatan := domain.SubKegiatan{
		Id:              uuId,
		KodeSubKegiatan: request.KodeSubkegiatan,
		NamaSubKegiatan: request.NamaSubKegiatan,
		Indikator:       indikators,
	}

	result, err := service.subKegiatanRepository.Create(ctx, tx, subKegiatan)
	if err != nil {
		log.Println("Gagal membuat data sub kegiatan:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	return helper.ToSubKegiatanResponse(result), nil
}

func (service *SubKegiatanServiceImpl) Update(ctx context.Context, request subkegiatan.SubKegiatanUpdateRequest) (subkegiatan.SubKegiatanResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		log.Println("Validasi gagal:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	var indikators []domain.Indikator

	for _, indikatorReq := range request.Indikator {
		indikatorId := indikatorReq.Id
		if indikatorId == "" {
			randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			indikatorId = fmt.Sprintf("IND-%s", randomDigits)
		}

		var targets []domain.Target

		for _, targetReq := range indikatorReq.Target {
			targetId := targetReq.Id
			if targetId == "" {
				randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
				targetId = fmt.Sprintf("TRG-%s", randomDigits)
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      targetReq.TargetIndikator,
				Satuan:      targetReq.SatuanIndikator,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:               indikatorId,
			SubKegiatanId:    request.Id,
			RencanaKinerjaId: indikatorReq.RencanaKinerjaId,
			Indikator:        indikatorReq.NamaIndikator,
			Target:           targets,
		}
		indikators = append(indikators, indikator)
	}

	domainSubKegiatan := domain.SubKegiatan{
		Id:              request.Id,
		KodeSubKegiatan: request.KodeSubkegiatan,
		NamaSubKegiatan: request.NamaSubKegiatan,
		Indikator:       indikators,
	}

	result, err := service.subKegiatanRepository.Update(ctx, tx, domainSubKegiatan)
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("gagal mengupdate sub kegiatan: %v", err)
	}

	response := helper.ToSubKegiatanResponse(result)
	return response, nil
}

func (service *SubKegiatanServiceImpl) FindById(ctx context.Context, subKegiatanId string) (subkegiatan.SubKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data SubKegiatan
	subKegiatan, err := service.subKegiatanRepository.FindById(ctx, tx, subKegiatanId)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("sub kegiatan dengan id %s tidak ditemukan", subKegiatanId)
		}
		log.Println("Gagal mencari data sub kegiatan:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	// Ambil data Indikator
	indikators, err := service.subKegiatanRepository.FindIndikatorBySubKegiatanId(ctx, tx, subKegiatanId)
	if err != nil {
		// Jika tidak ada indikator, gunakan array kosong
		if err == sql.ErrNoRows {
			subKegiatan.Indikator = []domain.Indikator{}
			return helper.ToSubKegiatanResponse(subKegiatan), nil
		}
		log.Printf("Gagal mengambil indikator untuk subkegiatan %s: %v", subKegiatanId, err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	// Untuk setiap Indikator, ambil Target
	for i, indikator := range indikators {
		targets, err := service.subKegiatanRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			// Jika tidak ada target, gunakan array kosong
			if err == sql.ErrNoRows {
				indikators[i].Target = []domain.Target{}
				continue
			}
			log.Printf("Gagal mengambil target untuk indikator %s: %v", indikator.Id, err)
			return subkegiatan.SubKegiatanResponse{}, err
		}
		indikators[i].Target = targets
	}

	// Gabungkan data
	subKegiatan.Indikator = indikators

	return helper.ToSubKegiatanResponse(subKegiatan), nil
}

func normalizeSubKegiatanFindAllFilter(filter subkegiatan.SubKegiatanFindAllFilter) subkegiatan.SubKegiatanFindAllFilter {
	filter.KodeSubKegiatan = strings.TrimSpace(filter.KodeSubKegiatan)
	filter.NamaSubKegiatan = strings.TrimSpace(filter.NamaSubKegiatan)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return filter
}

func attachIndikatorTargets(
	subKegiatans []domain.SubKegiatan,
	indikators []domain.Indikator,
	targets []domain.Target,
) {
	targetByIndikator := make(map[string][]domain.Target)
	for _, target := range targets {
		targetByIndikator[target.IndikatorId] = append(targetByIndikator[target.IndikatorId], target)
	}
	indikatorBySubKegiatan := make(map[string][]domain.Indikator)
	for _, indikator := range indikators {
		indikator.Target = targetByIndikator[indikator.Id]
		if indikator.Target == nil {
			indikator.Target = []domain.Target{}
		}
		indikatorBySubKegiatan[indikator.SubKegiatanId] = append(indikatorBySubKegiatan[indikator.SubKegiatanId], indikator)
	}
	for i := range subKegiatans {
		inds := indikatorBySubKegiatan[subKegiatans[i].Id]
		if inds == nil {
			inds = []domain.Indikator{}
		}
		subKegiatans[i].Indikator = inds
	}
}

func (service *SubKegiatanServiceImpl) FindAll(
	ctx context.Context, filter subkegiatan.SubKegiatanFindAllFilter,
) (subkegiatan.SubKegiatanPaginatedResponse, error) {
	filter = normalizeSubKegiatanFindAllFilter(filter)
	offset := (filter.Page - 1) * filter.Limit

	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanPaginatedResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	total, err := service.subKegiatanRepository.CountAll(ctx, tx, filter.KodeSubKegiatan, filter.NamaSubKegiatan)
	if err != nil {
		log.Println("Gagal menghitung data sub kegiatan:", err)
		return subkegiatan.SubKegiatanPaginatedResponse{}, err
	}

	if total == 0 {
		return subkegiatan.SubKegiatanPaginatedResponse{
			Items:      []subkegiatan.SubKegiatanResponse{},
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      0,
			TotalPages: 0,
		}, nil
	}

	subKegiatans, err := service.subKegiatanRepository.FindAll(
		ctx, tx, filter.KodeSubKegiatan, filter.NamaSubKegiatan, filter.Limit, offset,
	)
	if err != nil {
		log.Println("Gagal mencari data sub kegiatan:", err)
		return subkegiatan.SubKegiatanPaginatedResponse{}, err
	}

	subKegiatanIds := make([]string, 0, len(subKegiatans))
	for _, sp := range subKegiatans {
		subKegiatanIds = append(subKegiatanIds, sp.Id)
	}

	indikators, err := service.subKegiatanRepository.FindIndikatorsBySubKegiatanIds(ctx, tx, subKegiatanIds)
	if err != nil {
		log.Println("Gagal batch load indikator sub kegiatan:", err)
		return subkegiatan.SubKegiatanPaginatedResponse{}, err
	}

	indikatorIds := make([]string, 0, len(indikators))
	for _, ind := range indikators {
		indikatorIds = append(indikatorIds, ind.Id)
	}

	targets, err := service.subKegiatanRepository.FindTargetsByIndikatorIds(ctx, tx, indikatorIds)
	if err != nil {
		log.Println("Gagal batch load target sub kegiatan:", err)
		return subkegiatan.SubKegiatanPaginatedResponse{}, err
	}

	attachIndikatorTargets(subKegiatans, indikators, targets)

	totalPages := total / filter.Limit
	if total%filter.Limit != 0 {
		totalPages++
	}

	hasPrevious := filter.Page > 1 && totalPages > 0
	hasNext := filter.Page < totalPages
	previousPage, nextPage := 0, 0
	if hasPrevious {
		previousPage = filter.Page - 1
	}
	if hasNext {
		nextPage = filter.Page + 1
	}

	return subkegiatan.SubKegiatanPaginatedResponse{
		Items:        helper.ToSubKegiatanResponses(subKegiatans),
		Page:         filter.Page,
		Limit:        filter.Limit,
		Total:        total,
		TotalPages:   totalPages,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
		NextPage:     nextPage,
		PreviousPage: previousPage,
	}, nil
}

func (service *SubKegiatanServiceImpl) Delete(ctx context.Context, subKegiatanId string) error {
	// Validasi ID
	if subKegiatanId == "" {
		return errors.New("subkegiatan id tidak boleh kosong")
	}

	// Mulai transaksi
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Proses delete
	err = service.subKegiatanRepository.Delete(ctx, tx, subKegiatanId)
	if err != nil {
		return fmt.Errorf("gagal menghapus sub kegiatan: %v", err)
	}

	return nil
}

func (service *SubKegiatanServiceImpl) FindSubKegiatanKAK(ctx context.Context, kodeOpd string, kode string, tahun string) (subkegiatan.SubKegiatanKAKResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanKAKResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository dengan parameter kode_opd, kode, dan tahun
	data, err := service.subKegiatanRepository.FindSubKegiatanKAK(ctx, tx, kodeOpd, kode, tahun)
	if err != nil {
		log.Println("Gagal mengambil data subkegiatan KAK:", err)
		return subkegiatan.SubKegiatanKAKResponse{}, err
	}

	// Transform ke response
	response := subkegiatan.SubKegiatanKAKResponse{
		KodeOpd: data.KodeOpd,
		NamaOpd: data.NamaOpd,
		Program: subkegiatan.ProgramKAKResponse{
			Kode: data.KodeProgram,
			Nama: data.NamaProgram,
			IndikatorKinerjaProgram: subkegiatan.IndikatorKinerjaKAKResponse{
				Nama:   data.IndikatorProgram,
				Target: data.TargetProgram,
				Satuan: data.SatuanProgram,
			},
		},
		Kegiatan: subkegiatan.KegiatanKAKResponse{
			Kode: data.KodeKegiatan,
			Nama: data.NamaKegiatan,
			IndikatorKinerjaKegiatan: subkegiatan.IndikatorKinerjaKAKResponse{
				Nama:   data.IndikatorKegiatan,
				Target: data.TargetKegiatan,
				Satuan: data.SatuanKegiatan,
			},
		},
		SubKegiatan: subkegiatan.SubKegiatanDetailKAKResponse{
			Subkegiatan: data.KodeSubKegiatan,
			Nama:        data.NamaSubKegiatan,
			IndikatorKinerjaSubKegiatan: subkegiatan.IndikatorKinerjaKAKResponse{
				Nama:   data.IndikatorSubKegiatan,
				Target: data.TargetSubKegiatan,
				Satuan: data.SatuanSubKegiatan,
			},
		},
		PaguAnggaran: strconv.FormatInt(data.PaguAnggaran, 10),
	}

	return response, nil
}
