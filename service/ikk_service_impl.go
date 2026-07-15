package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/ikk"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type IkkServiceImpl struct {
	IkkRepository repository.IkkRepository
	DB            *sql.DB
	Validate      *validator.Validate
}

func NewIkkServiceImpl(ikkRepository repository.IkkRepository, db *sql.DB, validate *validator.Validate) *IkkServiceImpl {
	return &IkkServiceImpl{
		IkkRepository: ikkRepository,
		DB:            db,
		Validate:      validate,
	}
}

func (service *IkkServiceImpl) Create(ctx context.Context, request ikk.IkkRequest) (ikk.IkkResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// mapping indikator
	indikators := make([]domain.IndikatorIkk, 0)

	for _, indikatorReq := range request.Indikators {

		// mapping targets
		targets := make([]domain.TargetIkk, 0)

		for _, targetReq := range indikatorReq.Targets {
			targets = append(targets, domain.TargetIkk{
				Target:  targetReq.Target,
				Satuan:  targetReq.Satuan,
				Tahun :  targetReq.Tahun,
			})
		}

		indikators = append(indikators, domain.IndikatorIkk{
			Indikator:  indikatorReq.Indikator,
			Targets:    targets,
		})
	}

	data := domain.Ikk{
		KodeBidangUrusan: request.KodeBidangUrusan,
		KodeOpd:          request.KodeOpd,
		Jenis:            request.Jenis,
		Tahun:            request.Tahun,
		Keterangan:       request.Keterangan,
		Indikators:       indikators,
	}

	result, err := service.IkkRepository.Create(ctx, tx, data)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	// mapping response
	responseIndikators := make([]ikk.IndikatorResponse, 0)

	for _, indikator := range result.Indikators {

		responseTargets := make([]ikk.TargetResponse, 0)

		for _, target := range indikator.Targets {
			responseTargets = append(responseTargets, ikk.TargetResponse{
				ID:      target.ID,
				Target:  target.Target,
				Satuan:  target.Satuan,
				Tahun :  target.Tahun,
			})
		}

		responseIndikators = append(responseIndikators, ikk.IndikatorResponse{
			ID:         indikator.ID,
			Indikator:  indikator.Indikator,
			Targets:    responseTargets,
		})
	}

	return ikk.IkkResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Jenis:              result.Jenis,
		Tahun:              result.Tahun,
		Keterangan:         result.Keterangan,
		Indikators:         responseIndikators,
	}, nil
}

func (service *IkkServiceImpl) Update(
	ctx context.Context,
	request ikk.IkkUpdateRequest,
) (ikk.IkkResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.IkkRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	// ================= MAP INDIKATOR =================
	indikators := make([]domain.IndikatorIkk, 0)

	for _, indReq := range request.Indikators {

		targets := make([]domain.TargetIkk, 0)

		for _, tReq := range indReq.Targets {
			targets = append(targets, domain.TargetIkk{
				Target: tReq.Target,
				Satuan: tReq.Satuan,
				Tahun : tReq.Tahun,
			})
		}

		indikators = append(indikators, domain.IndikatorIkk{
			Indikator: indReq.Indikator,
			Targets:   targets,
		})
	}

	data := domain.Ikk{
		ID:                 request.ID,
		KodeBidangUrusan:   request.KodeBidangUrusan,
		KodeOpd:            request.KodeOpd,
		Jenis:              request.Jenis,
		Tahun:              request.Tahun,
		Keterangan:         request.Keterangan,
		Indikators:         indikators,
	}

	result, err := service.IkkRepository.Update(ctx, tx, data)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	// mapping response
	responseIndikators := make([]ikk.IndikatorResponse, 0)

	for _, ind := range result.Indikators {

		responseTargets := make([]ikk.TargetResponse, 0)

		for _, t := range ind.Targets {
			responseTargets = append(responseTargets, ikk.TargetResponse{
				ID:      t.ID,
				Target:  t.Target,
				Satuan:  t.Satuan,
				Tahun :  t.Tahun,
			})
		}

		responseIndikators = append(responseIndikators, ikk.IndikatorResponse{
			ID:         ind.ID,
			Indikator:  ind.Indikator,
			Targets:    responseTargets,
		})
	}

	return ikk.IkkResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Jenis:              result.Jenis,
		Tahun:              result.Tahun,
		Keterangan:         result.Keterangan,
		Indikators:         responseIndikators,
	}, nil
}

func (service *IkkServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IkkRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.IkkRepository.Delete(ctx, tx, id)
}

func (service *IkkServiceImpl) FindById(ctx context.Context, id int) (ikk.IkkResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.IkkRepository.FindById(ctx, tx, id)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	// mapping indikator
	indikators := make([]ikk.IndikatorResponse, 0)

	for _, ind := range result.Indikators {

		// mapping target
		targets := make([]ikk.TargetResponse, 0)

		for _, t := range ind.Targets {
			targets = append(targets, ikk.TargetResponse{
				ID:      t.ID,
				Target:  t.Target,
				Satuan:  t.Satuan,
				Tahun:  t.Tahun,
			})
		}

		indikators = append(indikators, ikk.IndikatorResponse{
			ID:         ind.ID,
			Indikator:  ind.Indikator,
			Targets:    targets,
		})
	}

	return ikk.IkkResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Jenis:              result.Jenis,
		Tahun:              result.Tahun,
		Keterangan:         result.Keterangan,
		Indikators:         indikators,
	}, nil
}


func (service *IkkServiceImpl) FindByKodeOpd(ctx context.Context, levelPohon int, kodeOpd string) ([]ikk.IkkFullResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []ikk.IkkFullResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	mapping := map[int]string{
		5: "outcome",
		6: "output",
	}

	jenis := mapping[levelPohon]

	if jenis == "" {
		return []ikk.IkkFullResponse{}, nil
	}

	bidangUrusans, err := service.IkkRepository.FindByKodeOpd(ctx, tx, jenis, kodeOpd)
	if err != nil {
		return []ikk.IkkFullResponse{}, err
	}

	var bidangUrusanResponses []ikk.IkkFullResponse
	for _, bidangUrusan := range bidangUrusans {
		bidangUrusanResponses = append(bidangUrusanResponses, ikk.IkkFullResponse{
			ID: bidangUrusan.ID,
			KodeBidangUrusan: bidangUrusan.KodeBidangUrusan,
			NamaBidangUrusan: bidangUrusan.NamaBidangUrusan,
			KodeOpd: bidangUrusan.KodeOpd,
			NamaOpd: bidangUrusan.NamaOpd,
			Jenis: bidangUrusan.Jenis,
			Tahun: bidangUrusan.Tahun,
			Keterangan: bidangUrusan.Keterangan,
		})
	}

	return bidangUrusanResponses, nil
}
func (service *IkkServiceImpl) FindAll(ctx context.Context, kodeOpd string) (ikk.IkkMasterResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkMasterResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil selection bidang urusan
	selections, err := service.IkkRepository.FindSelectionByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return ikk.IkkMasterResponse{}, err
	}

	// Ambil data IKK
	ikks, err := service.IkkRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return ikk.IkkMasterResponse{}, err
	}

	// Mapping selection
	selectionResponses := make([]ikk.BidangUrusanSelectionResponse, 0)

	for _, selection := range selections {
		selectionResponses = append(selectionResponses,
			ikk.BidangUrusanSelectionResponse{
				KodeBidangUrusan: selection.KodeBidangUrusan,
				NamaBidangUrusan: selection.NamaBidangUrusan,
				KodeOpd:          selection.KodeOpd,
				NamaOpd:          selection.NamaOpd,
			},
		)
	}

	// Mapping IKK
	ikkResponses := make([]ikk.IkkFullResponse, 0)

	for _, ikkData := range ikks {

		indikators := make([]ikk.IndikatorResponse, 0)

		for _, ind := range ikkData.Indikators {

			targets := make([]ikk.TargetResponse, 0)

			for _, t := range ind.Targets {
				targets = append(targets, ikk.TargetResponse{
					ID:     t.ID,
					Target: t.Target,
					Satuan: t.Satuan,
					Tahun: t.Tahun,
				})
			}

			indikators = append(indikators, ikk.IndikatorResponse{
				ID:        ind.ID,
				Indikator: ind.Indikator,
				Targets:   targets,
			})
		}

		ikkResponses = append(ikkResponses, ikk.IkkFullResponse{
			ID:               ikkData.ID,
			KodeOpd:          ikkData.KodeOpd,
			NamaOpd:          ikkData.NamaOpd,
			KodeBidangUrusan: ikkData.KodeBidangUrusan,
			NamaBidangUrusan: ikkData.NamaBidangUrusan,
			Jenis:            ikkData.Jenis,
			Tahun:            ikkData.Tahun,
			Keterangan:       ikkData.Keterangan,
			Indikators:       indikators,
		})
	}

	return ikk.IkkMasterResponse{
		BidangUrusanSelections: selectionResponses,
		Ikks:                   ikkResponses,
	}, nil
}
func (service *IkkServiceImpl) FindAllByIdPokin(ctx context.Context, pokinId int) ([]ikk.IkkFullResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return []ikk.IkkFullResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data IKK
	ikks, err := service.IkkRepository.FindAllByIdPokin(ctx, tx, pokinId)
	if err != nil {
		return []ikk.IkkFullResponse{}, err
	}

	// Mapping IKK
	ikkResponses := make([]ikk.IkkFullResponse, 0)

	for _, ikkData := range ikks {

		indikators := make([]ikk.IndikatorResponse, 0)

		for _, ind := range ikkData.Indikators {

			targets := make([]ikk.TargetResponse, 0)

			for _, t := range ind.Targets {
				targets = append(targets, ikk.TargetResponse{
					ID:     t.ID,
					Target: t.Target,
					Satuan: t.Satuan,
					Tahun: t.Tahun,
				})
			}

			indikators = append(indikators, ikk.IndikatorResponse{
				ID:        ind.ID,
				Indikator: ind.Indikator,
				Targets:   targets,
			})
		}

		ikkResponses = append(ikkResponses, ikk.IkkFullResponse{
			ID:               ikkData.ID,
			KodeOpd:          ikkData.KodeOpd,
			NamaOpd:          ikkData.NamaOpd,
			KodeBidangUrusan: ikkData.KodeBidangUrusan,
			NamaBidangUrusan: ikkData.NamaBidangUrusan,
			Jenis:            ikkData.Jenis,
			Tahun:            ikkData.Tahun,
			Keterangan:       ikkData.Keterangan,
			Indikators:       indikators,
		})
	}

	return ikkResponses, nil
}

func (service *IkkServiceImpl) FindAllByLevelPohon(
	ctx context.Context,
	levelPohon int,
	kodeOpd string,
) (ikk.IkkMasterResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkMasterResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// ================= mapping levelPohon -> jenis =================
	mapping := map[int]string{
		5: "outcome",
		6: "output",
	}

	jenis := mapping[levelPohon]
	if jenis == "" {
		return ikk.IkkMasterResponse{}, nil
	}

	// ================= selection =================
	selections, err := service.IkkRepository.FindSelectionByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return ikk.IkkMasterResponse{}, err
	}

	// ================= IKK data (FILTER jenis + kodeOpd) =================
	ikks, err := service.IkkRepository.FindAllByJenisAndKodeOpd(ctx, tx, kodeOpd, jenis)
	if err != nil {
		return ikk.IkkMasterResponse{}, err
	}

	// ================= mapping selection =================
	selectionResponses := make([]ikk.BidangUrusanSelectionResponse, 0)
	for _, s := range selections {
		selectionResponses = append(selectionResponses, ikk.BidangUrusanSelectionResponse{
			KodeBidangUrusan: s.KodeBidangUrusan,
			NamaBidangUrusan: s.NamaBidangUrusan,
			KodeOpd:          s.KodeOpd,
			NamaOpd:          s.NamaOpd,
		})
	}

	// ================= mapping IKK nested =================
	ikkResponses := make([]ikk.IkkFullResponse, 0)

	for _, ikkData := range ikks {

		indikators := make([]ikk.IndikatorResponse, 0)

		for _, ind := range ikkData.Indikators {

			targets := make([]ikk.TargetResponse, 0)

			for _, t := range ind.Targets {
				targets = append(targets, ikk.TargetResponse{
					ID:     t.ID,
					Target: t.Target,
					Satuan: t.Satuan,
					Tahun: t.Tahun,
				})
			}

			indikators = append(indikators, ikk.IndikatorResponse{
				ID:        ind.ID,
				Indikator: ind.Indikator,
				Targets:   targets,
			})
		}

		ikkResponses = append(ikkResponses, ikk.IkkFullResponse{
			ID:               ikkData.ID,
			KodeOpd:          ikkData.KodeOpd,
			NamaOpd:          ikkData.NamaOpd,
			KodeBidangUrusan: ikkData.KodeBidangUrusan,
			NamaBidangUrusan: ikkData.NamaBidangUrusan,
			Jenis:            ikkData.Jenis,
			Tahun:            ikkData.Tahun,
			Keterangan:       ikkData.Keterangan,
			Indikators:       indikators,
		})
	}

	return ikk.IkkMasterResponse{
		BidangUrusanSelections: selectionResponses,
		Ikks:                   ikkResponses,
	}, nil
}

func (service *IkkServiceImpl) PilihIkk(ctx context.Context, request ikk.IkkTerpilihCreateRequest) (ikk.IkkTerpilihResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return ikk.IkkTerpilihResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkTerpilihResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.IkkTerpilih{
		PohonKinerjaId: request.PohonKinerjaId,
		IkkId:          request.IkkId,
	}

	result, err := service.IkkRepository.PilihIkk(ctx, tx, data)
	if err != nil {
		return ikk.IkkTerpilihResponse{}, err
	}

	pokinikk, err := service.IkkRepository.FindTerpilihPokinIkkById(
		ctx,
		tx,
		result.Id,
	)

	if err != nil {
		return ikk.IkkTerpilihResponse{}, err
	}

	return ikk.IkkTerpilihResponse{
		Id:              result.Id,
		PokinId:   		 result.PohonKinerjaId,
		IkkId:           result.IkkId,
		NamaPokin:       pokinikk.NamaPokin,
		JenisIkk:        pokinikk.JenisIkk,
		KeteranganIkk:   pokinikk.KeteranganIkk,
	}, nil
}

func (service *IkkServiceImpl) DeletePilihanIkk(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IkkRepository.FindTerpilihById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.IkkRepository.DeletePilihanIkk(ctx, tx, id)
}
