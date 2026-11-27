package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type ProgramServiceImpl struct {
	programRepository repository.ProgramRepository
	DB                *sql.DB
}

func NewProgramServiceImpl(programRepository repository.ProgramRepository, DB *sql.DB) *ProgramServiceImpl {
	return &ProgramServiceImpl{
		programRepository: programRepository,
		DB:                DB,
	}
}

func (service *ProgramServiceImpl) Create(ctx context.Context, request programkegiatan.ProgramKegiatanCreateRequest) (programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, err
	}

	defer helper.CommitOrRollback(tx)

	uuidPrgm := fmt.Sprintf("PRGM-%s", request.KodeProgram)

	program := domainmaster.ProgramKegiatan{
		Id:          uuidPrgm,
		KodeProgram: request.KodeProgram,
		NamaProgram: request.NamaProgram,
		Tahun:       request.Tahun,
		IsActive:    request.IsActive,
	}

	var indikators []domain.Indikator
	for _, indikatorRequest := range request.Indikator {
		uuidIndikator := fmt.Sprintf("IND-PRG-%s", uuid.New().String()[:5])
		indikator := domain.Indikator{
			Id:        uuidIndikator,
			ProgramId: program.Id,
			Indikator: indikatorRequest.Indikator,
			Tahun:     indikatorRequest.Tahun,
		}

		var targets []domain.Target
		for _, targetRequest := range indikatorRequest.Target {
			uuidTarget := fmt.Sprintf("TRGT-PRG-%s", uuid.New().String()[:5])
			target := domain.Target{
				Id:          uuidTarget,
				IndikatorId: indikator.Id,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.Satuan,
				Tahun:       targetRequest.Tahun,
			}
			targets = append(targets, target)
		}
		indikator.Target = targets
		indikators = append(indikators, indikator)
	}
	program.Indikator = indikators

	result, err := service.programRepository.Create(ctx, tx, program)
	if err != nil {
		tx.Rollback()
		return programkegiatan.ProgramKegiatanResponse{}, err
	}
	var indikatorResponses []programkegiatan.IndikatorResponse
	for _, indikator := range result.Indikator {
		var targetResponses []programkegiatan.TargetResponse
		for _, target := range indikator.Target {
			targetResponse := programkegiatan.TargetResponse{
				Id:          target.Id,
				IndikatorId: target.IndikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := programkegiatan.IndikatorResponse{
			Id:        indikator.Id,
			ProgramId: indikator.ProgramId,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	// if err = tx.Commit(); err != nil {
	// 	tx.Rollback()
	// 	return programkegiatan.ProgramKegiatanResponse{}, err
	// }

	return programkegiatan.ProgramKegiatanResponse{
		Id:          result.Id,
		KodeProgram: result.KodeProgram,
		NamaProgram: result.NamaProgram,
		Tahun:       result.Tahun,
		IsActive:    result.IsActive,
		Indikator:   indikatorResponses,
	}, nil
}

func (service *ProgramServiceImpl) Update(ctx context.Context, request programkegiatan.ProgramKegiatanUpdateRequest) (programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	fmt.Printf("\n=== MULAI PROSES UPDATE ===\n")
	fmt.Printf("Request ID Program: %s\n", request.Id)

	program := domainmaster.ProgramKegiatan{
		Id:          request.Id,
		KodeProgram: request.KodeProgram,
		NamaProgram: request.NamaProgram,
		Tahun:       request.Tahun,
		IsActive:    request.IsActive,
	}

	fmt.Printf("\n=== MEMPROSES INDIKATOR ===\n")
	var indikators []domain.Indikator
	for i, indikator := range request.Indikator {
		fmt.Printf("\nIndikator ke-%d:\n", i+1)
		fmt.Printf("ID Indikator dari request: %s\n", indikator.Id)

		indikatorId := indikator.Id
		if indikatorId == "" {
			indikatorId = fmt.Sprintf("IND-KGT-%s", uuid.New().String()[:5])
			fmt.Printf("Generated new Indikator ID: %s\n", indikatorId)
		} else {
			fmt.Printf("Menggunakan ID Indikator existing: %s\n", indikatorId)
		}

		fmt.Printf("Memproses target untuk indikator: %s\n", indikatorId)
		var targets []domain.Target
		for j, target := range indikator.Target {
			fmt.Printf("\nTarget ke-%d:\n", j+1)
			fmt.Printf("ID Target dari request: %s\n", target.Id)

			targetId := target.Id
			if targetId == "" {
				targetId = fmt.Sprintf("TRGT-KGT-%s", uuid.New().String()[:5])
				fmt.Printf("Generated new Target ID: %s\n", targetId)
			} else {
				fmt.Printf("Menggunakan ID Target existing: %s\n", targetId)
			}

			targetDomain := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			fmt.Printf("Target Domain: %+v\n", targetDomain)
			targets = append(targets, targetDomain)
		}

		indikatorDomain := domain.Indikator{
			Id:        indikatorId,
			ProgramId: request.Id,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targets,
		}
		fmt.Printf("Indikator Domain: %+v\n", indikatorDomain)
		indikators = append(indikators, indikatorDomain)
	}
	program.Indikator = indikators

	fmt.Printf("\n=== MEMANGGIL REPOSITORY UPDATE ===\n")
	fmt.Printf("Program yang akan diupdate: %+v\n", program)

	result, err := service.programRepository.Update(ctx, tx, program)
	if err != nil {
		fmt.Printf("ERROR: Gagal melakukan update di repository: %v\n", err)
		return programkegiatan.ProgramKegiatanResponse{}, err
	}
	fmt.Printf("Hasil update dari repository: %+v\n", result)

	fmt.Printf("\n=== MEMBANGUN RESPONSE ===\n")
	var indikatorResponses []programkegiatan.IndikatorResponse
	for i, indikator := range result.Indikator {
		fmt.Printf("\nMemproses response untuk indikator ke-%d:\n", i+1)
		var targetResponses []programkegiatan.TargetResponse
		for j, target := range indikator.Target {
			fmt.Printf("Memproses response untuk target ke-%d:\n", j+1)
			targetResponse := programkegiatan.TargetResponse{
				Id:          target.Id,
				IndikatorId: target.IndikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			fmt.Printf("Target Response: %+v\n", targetResponse)
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := programkegiatan.IndikatorResponse{
			Id:        indikator.Id,
			ProgramId: indikator.ProgramId,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targetResponses,
		}
		fmt.Printf("Indikator Response: %+v\n", indikatorResponse)
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	finalResponse := programkegiatan.ProgramKegiatanResponse{
		Id:          result.Id,
		KodeProgram: result.KodeProgram,
		NamaProgram: result.NamaProgram,

		Tahun:     result.Tahun,
		IsActive:  result.IsActive,
		Indikator: indikatorResponses,
	}
	fmt.Printf("\n=== RESPONSE FINAL ===\n")
	fmt.Printf("%+v\n", finalResponse)
	fmt.Printf("\n=== SELESAI PROSES UPDATE ===\n")

	return finalResponse, nil
}

func (service *ProgramServiceImpl) Delete(ctx context.Context, programId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.programRepository.Delete(ctx, tx, programId)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProgramServiceImpl) FindById(ctx context.Context, programId string) (programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Mengambil data program
	result, err := service.programRepository.FindById(ctx, tx, programId)
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal mengambil data program: %v", err)
	}

	// Mengambil semua indikator untuk program ini
	indikators, err := service.programRepository.FindIndikatorByProgramId(ctx, tx, result.Id)
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal mengambil data indikator: %v", err)
	}

	var indikatorResponses []programkegiatan.IndikatorResponse
	for _, indikator := range indikators {
		// Mengambil semua target untuk setiap indikator
		targets, err := service.programRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal mengambil data target: %v", err)
		}

		var targetResponses []programkegiatan.TargetResponse
		for _, target := range targets {
			targetResponse := programkegiatan.TargetResponse{
				Id:          target.Id,
				IndikatorId: target.IndikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := programkegiatan.IndikatorResponse{
			Id:        indikator.Id,
			ProgramId: indikator.ProgramId,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	return programkegiatan.ProgramKegiatanResponse{
		Id:          result.Id,
		KodeProgram: result.KodeProgram,
		NamaProgram: result.NamaProgram,
		Tahun:       result.Tahun,
		IsActive:    result.IsActive,
		Indikator:   indikatorResponses,
	}, nil
}

func (service *ProgramServiceImpl) FindAll(ctx context.Context) ([]programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Mengambil semua program
	results, err := service.programRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data program: %v", err)
	}

	var programResponses []programkegiatan.ProgramKegiatanResponse

	for _, program := range results {

		// Mengambil semua indikator untuk program ini
		indikators, err := service.programRepository.FindIndikatorByProgramId(ctx, tx, program.Id)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data indikator untuk program %s: %v", program.Id, err)
		}

		var indikatorResponses []programkegiatan.IndikatorResponse
		for _, indikator := range indikators {
			// Mengambil semua target untuk setiap indikator
			targets, err := service.programRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil {
				return nil, fmt.Errorf("gagal mengambil data target untuk indikator %s: %v", indikator.Id, err)
			}

			// Membuat response untuk semua target
			var targetResponses []programkegiatan.TargetResponse
			for _, target := range targets {
				targetResponse := programkegiatan.TargetResponse{
					Id:          target.Id,
					IndikatorId: target.IndikatorId,
					Target:      target.Target,
					Satuan:      target.Satuan,
					Tahun:       target.Tahun,
				}
				targetResponses = append(targetResponses, targetResponse)
			}

			// Membuat response untuk indikator dengan semua targetnya
			indikatorResponse := programkegiatan.IndikatorResponse{
				Id:        indikator.Id,
				ProgramId: indikator.ProgramId,
				Indikator: indikator.Indikator,
				Tahun:     indikator.Tahun,
				Target:    targetResponses,
			}
			indikatorResponses = append(indikatorResponses, indikatorResponse)
		}

		// Membuat response untuk program dengan indikator dan targetnya
		programResponse := programkegiatan.ProgramKegiatanResponse{
			Id:          program.Id,
			KodeProgram: program.KodeProgram,
			NamaProgram: program.NamaProgram,
			Tahun:       program.Tahun,
			IsActive:    program.IsActive,
			Indikator:   indikatorResponses,
		}
		programResponses = append(programResponses, programResponse)
	}

	return programResponses, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// UpsertIndikator — create/update/delete indikator program dalam 1 request
// ─────────────────────────────────────────────────────────────────────────────

func (service *ProgramServiceImpl) UpsertIndikator(
	ctx context.Context,
	programId string,
	req programkegiatan.UpsertIndikatorProgramRequest,
) (programkegiatan.UpsertIndikatorProgramResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.UpsertIndikatorProgramResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Pastikan program ada
	if _, err := service.programRepository.FindById(ctx, tx, programId); err != nil {
		return programkegiatan.UpsertIndikatorProgramResponse{}, fmt.Errorf("program '%s' tidak ditemukan", programId)
	}

	// 2. Ambil semua existing indikator IDs untuk program ini
	existingIds, err := service.programRepository.FindIndikatorIdsByProgramId(ctx, tx, programId)
	if err != nil {
		return programkegiatan.UpsertIndikatorProgramResponse{}, err
	}

	// 3. Proses tiap indikator dari request
	requestIdSet := make(map[string]bool)
	for _, item := range req.Indikator {
		// Generate ID baru jika tidak ada
		if item.Id == "" {
			item.Id = fmt.Sprintf("IND-PRG-%s", uuid.New().String()[:8])
		}
		requestIdSet[item.Id] = true

		// Upsert indikator
		if err := service.programRepository.UpsertIndikator(ctx, tx, domain.Indikator{
			Id:        item.Id,
			ProgramId: programId,
			Indikator: item.Indikator,
			Tahun:     item.Tahun,
		}); err != nil {
			return programkegiatan.UpsertIndikatorProgramResponse{}, err
		}

		// Replace target untuk indikator ini
		targets := make([]domain.Target, 0, len(item.Target))
		for _, t := range item.Target {
			if t.Id == "" {
				t.Id = fmt.Sprintf("TRGT-PRG-%s", uuid.New().String()[:8])
			}
			targets = append(targets, domain.Target{
				Id:          t.Id,
				IndikatorId: item.Id,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       t.Tahun,
			})
		}
		if err := service.programRepository.ReplaceTargets(ctx, tx, item.Id, targets); err != nil {
			return programkegiatan.UpsertIndikatorProgramResponse{}, err
		}
	}

	// 4. Hapus indikator existing yang TIDAK ada di request
	for _, existId := range existingIds {
		if !requestIdSet[existId] {
			if err := service.programRepository.DeleteIndikatorWithTargets(ctx, tx, existId); err != nil {
				return programkegiatan.UpsertIndikatorProgramResponse{}, err
			}
		}
	}

	// 5. Ambil data terbaru untuk response
	indikators, err := service.programRepository.FindIndikatorWithTargetsByProgramId(ctx, tx, programId)
	if err != nil {
		return programkegiatan.UpsertIndikatorProgramResponse{}, err
	}

	// 6. Bangun response
	indResp := make([]programkegiatan.IndikatorProgramUpsertResponse, 0, len(indikators))
	for _, ind := range indikators {
		targetResp := make([]programkegiatan.TargetResponse, 0, len(ind.Target))
		for _, t := range ind.Target {
			targetResp = append(targetResp, programkegiatan.TargetResponse{
				Id:          t.Id,
				IndikatorId: t.IndikatorId,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       t.Tahun,
			})
		}
		indResp = append(indResp, programkegiatan.IndikatorProgramUpsertResponse{
			Id:        ind.Id,
			ProgramId: ind.ProgramId,
			Indikator: ind.Indikator,
			Tahun:     ind.Tahun,
			Target:    targetResp,
		})
	}

	return programkegiatan.UpsertIndikatorProgramResponse{
		ProgramId: programId,
		Indikator: indResp,
	}, nil
}

func (service *ProgramServiceImpl) FindAllWithIndikator(ctx context.Context) ([]programkegiatan.ProgramKegiatanWithIndikatorV2Response, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	programs, err := service.programRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data program: %v", err)
	}

	kodes := make([]string, 0, len(programs))
	for _, p := range programs {
		if p.KodeProgram != "" {
			kodes = append(kodes, p.KodeProgram)
		}
	}

	indikatorByKode, err := service.programRepository.FindMasterIndikatorProgramKegiatanByKodes(ctx, tx, kodes)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil indikator master: %v", err)
	}

	responses := make([]programkegiatan.ProgramKegiatanWithIndikatorV2Response, 0, len(programs))
	for _, p := range programs {
		responses = append(responses, programkegiatan.ProgramKegiatanWithIndikatorV2Response{
			Id:          p.Id,
			KodeProgram: p.KodeProgram,
			NamaProgram: p.NamaProgram,
			Tahun:       p.Tahun,
			IsActive:    p.IsActive,
			Indikator:   programkegiatan.MapMasterIndikatorResponses(indikatorByKode[p.KodeProgram]),
		})
	}
	return responses, nil
}
