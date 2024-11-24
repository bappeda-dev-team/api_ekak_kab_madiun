package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"

	"log"

	"fmt"

	"sort"

	"errors"

	"github.com/google/uuid"
)

type PohonKinerjaAdminServiceImpl struct {
	pohonKinerjaRepository repository.PohonKinerjaRepository
	opdRepository          repository.OpdRepository
	DB                     *sql.DB
}

func NewPohonKinerjaAdminServiceImpl(pohonKinerjaRepository repository.PohonKinerjaRepository, opdRepository repository.OpdRepository, DB *sql.DB) *PohonKinerjaAdminServiceImpl {
	return &PohonKinerjaAdminServiceImpl{
		pohonKinerjaRepository: pohonKinerjaRepository,
		opdRepository:          opdRepository,
		DB:                     DB,
	}
}

func (service *PohonKinerjaAdminServiceImpl) Create(ctx context.Context, request pohonkinerja.PohonKinerjaAdminCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	log.Printf("Memulai proses pembuatan PohonKinerja untuk tahun: %s", request.Tahun)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error memulai transaksi: %v", err)
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Logging persiapan indikator
	log.Printf("Mempersiapkan %d indikator", len(request.Indikator))

	// Persiapkan data indikator dan target
	var indikators []domain.Indikator
	for _, ind := range request.Indikator {
		indikatorId := "IND-POKIN-" + uuid.New().String()

		var targets []domain.Target
		for _, t := range ind.Target {
			targetId := "TRGT-IND-POKIN-" + uuid.New().String()
			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       request.Tahun,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:        indikatorId,
			Indikator: ind.NamaIndikator,
			Tahun:     request.Tahun,
			Target:    targets,
		}
		indikators = append(indikators, indikator)
	}

	pohonKinerja := domain.PohonKinerja{
		Parent:     request.Parent,
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    helper.EmptyStringIfNull(request.KodeOpd),
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Indikator:  indikators,
	}

	log.Printf("Menyimpan PohonKinerja dengan NamaPohon: %s, LevelPohon: %d", request.NamaPohon, request.LevelPohon)
	result, err := service.pohonKinerjaRepository.CreatePokinAdmin(ctx, tx, pohonKinerja)
	if err != nil {
		log.Printf("Error saat menyimpan PohonKinerja: %v", err)
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	log.Printf("Berhasil membuat PohonKinerja dengan ID: %d", result.Id)

	// Konversi indikator domain ke IndikatorResponse
	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range result.Indikator {
		var targetResponses []pohonkinerja.TargetResponse
		for _, t := range ind.Target {
			targetResponse := pohonkinerja.TargetResponse{
				Id:              t.Id,
				IndikatorId:     t.IndikatorId,
				TargetIndikator: t.Target,
				SatuanIndikator: t.Satuan,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := pohonkinerja.IndikatorResponse{
			Id:            ind.Id,
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	response := pohonkinerja.PohonKinerjaAdminResponseData{
		Id:         result.Id,
		Parent:     result.Parent,
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Indikators: indikatorResponses,
	}

	log.Printf("Proses pembuatan PohonKinerja selesai")
	return response, nil
}

func (service *PohonKinerjaAdminServiceImpl) Update(ctx context.Context, request pohonkinerja.PohonKinerjaAdminUpdateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists
	_, err = service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	// Persiapkan data indikator dan target
	var indikators []domain.Indikator
	for _, ind := range request.Indikator {
		indikatorId := ""
		if ind.Id != "" {
			indikatorId = ind.Id
		} else {
			indikatorId = "IND-POKIN-" + uuid.New().String()[:4]
		}

		var targets []domain.Target
		for _, t := range ind.Target {
			targetId := ""
			if t.Id != "" {
				targetId = t.Id
			} else {
				targetId = "TRGT-IND-POKIN-" + uuid.New().String()[:4]
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       request.Tahun,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:        indikatorId,
			Indikator: ind.NamaIndikator,
			Tahun:     request.Tahun,
			Target:    targets,
		}
		indikators = append(indikators, indikator)
	}

	pohonKinerja := domain.PohonKinerja{
		Id:         request.Id,
		Parent:     request.Parent,
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    helper.EmptyStringIfNull(request.KodeOpd),
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Indikator:  indikators,
	}

	result, err := service.pohonKinerjaRepository.UpdatePokinAdmin(ctx, tx, pohonKinerja)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	// Konversi indikator domain ke IndikatorResponse
	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range result.Indikator {
		var targetResponses []pohonkinerja.TargetResponse
		for _, t := range ind.Target {
			targetResponse := pohonkinerja.TargetResponse{
				Id:              t.Id,
				IndikatorId:     t.IndikatorId,
				TargetIndikator: t.Target,
				SatuanIndikator: t.Satuan,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := pohonkinerja.IndikatorResponse{
			Id:            ind.Id,
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	response := pohonkinerja.PohonKinerjaAdminResponseData{
		Id:         result.Id,
		Parent:     result.Parent,
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Indikators: indikatorResponses,
	}

	return response, nil
}

func (service *PohonKinerjaAdminServiceImpl) Delete(ctx context.Context, id int) error {
	// Mulai transaksi
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists sebelum dihapus
	pokin, err := service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("data tidak ditemukan: %v", err)
	}

	// Validasi tambahan: pastikan data yang akan dihapus memiliki level yang sesuai
	// Ini opsional, tergantung kebutuhan bisnis
	if pokin.LevelPohon < 0 || pokin.LevelPohon > 6 {
		return fmt.Errorf("level pohon kinerja tidak valid")
	}

	// Lakukan penghapusan secara hierarki
	err = service.pohonKinerjaRepository.DeletePokinAdmin(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data: %v", err)
	}

	return nil
}

func (service *PohonKinerjaAdminServiceImpl) FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	pokin, err := service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	log.Printf("Pohon Kinerja ditemukan: %+v", pokin)

	// Ambil data OPD jika kode OPD ada
	if pokin.KodeOpd != "" {
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
		if err == nil {
			pokin.NamaOpd = opd.NamaOpd
		}
	}

	// Konversi pokin.Id dari int ke string
	pokinIdStr := fmt.Sprint(pokin.Id)

	// Ambil indikator berdasarkan pokin ID
	indikators, err := service.pohonKinerjaRepository.FindIndikatorByPokinId(ctx, tx, pokinIdStr)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	log.Printf("Indikator ditemukan: %+v", indikators)

	// Konversi indikator domain ke IndikatorResponse
	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range indikators {
		// Ambil target berdasarkan indikator ID
		targets, err := service.pohonKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
		if err != nil {
			return pohonkinerja.PohonKinerjaAdminResponseData{}, err
		}

		var targetResponses []pohonkinerja.TargetResponse
		for _, t := range targets {
			targetResponse := pohonkinerja.TargetResponse{
				Id:              t.Id,
				IndikatorId:     t.IndikatorId,
				TargetIndikator: t.Target,
				SatuanIndikator: t.Satuan,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := pohonkinerja.IndikatorResponse{
			Id:            ind.Id,
			IdPokin:       ind.PokinId,
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	response := pohonkinerja.PohonKinerjaAdminResponseData{
		Id:         pokin.Id,
		Parent:     pokin.Parent,
		NamaPohon:  pokin.NamaPohon,
		NamaOpd:    pokin.NamaOpd,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
		Indikators: indikatorResponses,
	}

	return response, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindAll(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua data pohon kinerja
	pokins, err := service.pohonKinerjaRepository.FindPokinAdminAll(ctx, tx, tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}

	// Buat map untuk menyimpan data berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)

	// Kelompokkan data dan ambil data OPD untuk setiap pohon kinerja
	for i := range pokins {
		level := pokins[i].LevelPohon

		// Inisialisasi map untuk level jika belum ada
		if pohonMap[level] == nil {
			pohonMap[level] = make(map[int][]domain.PohonKinerja)
		}

		if pokins[i].KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokins[i].KodeOpd)
			if err == nil {
				pokins[i].NamaOpd = opd.NamaOpd
			}
		}

		pohonMap[level][pokins[i].Parent] = append(
			pohonMap[level][pokins[i].Parent],
			pokins[i],
		)
	}

	// Bangun response dimulai dari Tematik (level 0)
	var tematiks []pohonkinerja.TematikResponse
	for _, tematik := range pohonMap[0][0] {
		tematikResp := helper.BuildTematikResponse(pohonMap, tematik)
		tematiks = append(tematiks, tematikResp)
	}

	return pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   tahun,
		Tematik: tematiks,
	}, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindSubTematik(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua data pohon kinerja
	pokins, err := service.pohonKinerjaRepository.FindPokinAdminAll(ctx, tx, tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}

	// Buat map untuk menyimpan data berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for i := 1; i <= 2; i++ { // Hanya inisialisasi level 1 dan 2
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	// Filter dan kelompokkan data
	for _, p := range pokins {
		// Ambil data OPD jika ada
		if p.KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, p.KodeOpd)
			if err == nil {
				p.NamaOpd = opd.NamaOpd
			}
		}

		// Hanya masukkan data level 1 dan 2
		if p.LevelPohon >= 1 && p.LevelPohon <= 2 {
			pohonMap[p.LevelPohon][p.Parent] = append(pohonMap[p.LevelPohon][p.Parent], p)
		}
	}

	// Bangun response dimulai dari SubTematik (level 1)
	var tematiks []pohonkinerja.TematikResponse
	for _, subTematiks := range pohonMap[1] {
		// Urutkan subTematiks berdasarkan Id
		sort.Slice(subTematiks, func(i, j int) bool {
			return subTematiks[i].Id < subTematiks[j].Id
		})

		for _, subTematik := range subTematiks {
			var childs []interface{}

			// Tambahkan subsubtematik ke childs
			if subSubTematiks := pohonMap[2][subTematik.Id]; len(subSubTematiks) > 0 {
				// Urutkan subSubTematiks berdasarkan Id
				sort.Slice(subSubTematiks, func(i, j int) bool {
					return subSubTematiks[i].Id < subSubTematiks[j].Id
				})

				for _, subSubTematik := range subSubTematiks {
					subSubTematikResp := helper.BuildSubSubTematikResponse(pohonMap, subSubTematik)
					childs = append(childs, subSubTematikResp)
				}
			}

			tematikResp := pohonkinerja.TematikResponse{
				Id:         subTematik.Id,
				Parent:     &subTematik.Parent,
				Tema:       subTematik.NamaPohon,
				JenisPohon: subTematik.JenisPohon,
				LevelPohon: subTematik.LevelPohon,
				Keterangan: subTematik.Keterangan,
				Indikators: helper.ConvertToIndikatorResponses(subTematik.Indikator),
				Child:      childs,
			}
			tematiks = append(tematiks, tematikResp)
		}
	}

	// Urutkan hasil akhir berdasarkan Id
	sort.Slice(tematiks, func(i, j int) bool {
		return tematiks[i].Id < tematiks[j].Id
	})

	return pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   tahun,
		Tematik: tematiks,
	}, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindPokinAdminByIdHierarki(ctx context.Context, idPokin int) (pohonkinerja.TematikResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.TematikResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data pohon kinerja
	pokin, err := service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, idPokin)
	if err != nil {
		return pohonkinerja.TematikResponse{}, err
	}

	// Validasi level pohon harus 0
	if pokin.LevelPohon != 0 {
		return pohonkinerja.TematikResponse{}, fmt.Errorf("id yang diberikan bukan merupakan level tematik (level 0)")
	}

	// Ambil semua data pohon kinerja
	pokins, err := service.pohonKinerjaRepository.FindPokinAdminByIdHierarki(ctx, tx, idPokin)
	if err != nil {
		return pohonkinerja.TematikResponse{}, err
	}

	// Buat map untuk menyimpan data berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)

	// Kelompokkan data
	for _, p := range pokins {
		level := p.LevelPohon

		// Inisialisasi map untuk level jika belum ada
		if pohonMap[level] == nil {
			pohonMap[level] = make(map[int][]domain.PohonKinerja)
		}

		if p.KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, p.KodeOpd)
			if err == nil {
				p.NamaOpd = opd.NamaOpd
			}
		}
		pohonMap[level][p.Parent] = append(pohonMap[level][p.Parent], p)
	}

	// Bangun response hierarki
	var tematikResponse pohonkinerja.TematikResponse
	if tematik, exists := pohonMap[0][0]; exists && len(tematik) > 0 {
		var childs []interface{}

		// Tambahkan strategic langsung ke childs jika ada
		if strategics := pohonMap[4][tematik[0].Id]; len(strategics) > 0 {
			// Urutkan strategic berdasarkan Id
			sort.Slice(strategics, func(i, j int) bool {
				return strategics[i].Id < strategics[j].Id
			})

			for _, strategic := range strategics {
				strategicResp := helper.BuildStrategicResponse(pohonMap, strategic)
				childs = append(childs, strategicResp)
			}
		}

		// Tambahkan subtematik ke childs
		if subTematiks := pohonMap[1][tematik[0].Id]; len(subTematiks) > 0 {
			// Urutkan subtematik berdasarkan Id
			sort.Slice(subTematiks, func(i, j int) bool {
				return subTematiks[i].Id < subTematiks[j].Id
			})

			for _, subTematik := range subTematiks {
				subTematikResp := helper.BuildSubTematikResponse(pohonMap, subTematik)
				childs = append(childs, subTematikResp)
			}
		}

		tematikResponse = pohonkinerja.TematikResponse{
			Id:         tematik[0].Id,
			Parent:     nil,
			Tema:       tematik[0].NamaPohon,
			JenisPohon: tematik[0].JenisPohon,
			LevelPohon: tematik[0].LevelPohon,
			Keterangan: tematik[0].Keterangan,
			Indikators: helper.ConvertToIndikatorResponses(tematik[0].Indikator),
			Child:      childs,
		}
	}

	return tematikResponse, nil
}

func (service *PohonKinerjaAdminServiceImpl) CreateStrategicAdmin(ctx context.Context, request pohonkinerja.PohonKinerjaAdminStrategicCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	existingPokin, err := service.pohonKinerjaRepository.FindPokinToClone(ctx, tx, request.IdToClone)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	err = service.pohonKinerjaRepository.ValidateParentLevel(ctx, tx, request.Parent, existingPokin.LevelPohon)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	// Validasi JenisPohon
	if request.JenisPohon == "" {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, errors.New("jenis pohon tidak boleh kosong")
	}

	newPokin := domain.PohonKinerja{
		Parent:     request.Parent,
		NamaPohon:  existingPokin.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: existingPokin.LevelPohon,
		KodeOpd:    existingPokin.KodeOpd,
		Keterangan: existingPokin.Keterangan,
		Tahun:      existingPokin.Tahun,
	}

	newPokinId, err := service.pohonKinerjaRepository.InsertClonedPokin(ctx, tx, newPokin)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	indikators, err := service.pohonKinerjaRepository.FindIndikatorToClone(ctx, tx, request.IdToClone)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	var indikatorResponses []pohonkinerja.IndikatorResponse

	for _, indikator := range indikators {
		newIndikatorId := "IND-POKIN-" + uuid.New().String()[:6]

		err = service.pohonKinerjaRepository.InsertClonedIndikator(ctx, tx, newIndikatorId, newPokinId, indikator)
		if err != nil {
			return pohonkinerja.PohonKinerjaAdminResponseData{}, err
		}

		targets, err := service.pohonKinerjaRepository.FindTargetToClone(ctx, tx, indikator.Id)
		if err != nil {
			return pohonkinerja.PohonKinerjaAdminResponseData{}, err
		}

		var targetResponses []pohonkinerja.TargetResponse

		for _, target := range targets {
			newTargetId := "TRGT-IND-POKIN-" + uuid.New().String()[:5]
			err = service.pohonKinerjaRepository.InsertClonedTarget(ctx, tx, newTargetId, newIndikatorId, target)
			if err != nil {
				return pohonkinerja.PohonKinerjaAdminResponseData{}, err
			}

			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              newTargetId,
				IndikatorId:     newIndikatorId,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			})
		}

		indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
			Id:            newIndikatorId,
			IdPokin:       fmt.Sprint(newPokinId),
			NamaIndikator: indikator.Indikator,
			Target:        targetResponses,
		})
	}

	response := pohonkinerja.PohonKinerjaAdminResponseData{
		Id:         int(newPokinId),
		Parent:     existingPokin.Parent,
		NamaPohon:  existingPokin.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: existingPokin.LevelPohon,
		KodeOpd:    existingPokin.KodeOpd,
		Keterangan: existingPokin.Keterangan,
		Tahun:      existingPokin.Tahun,
		Indikators: indikatorResponses,
	}

	return response, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindPokinByTematik(ctx context.Context, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	pokins, err := service.pohonKinerjaRepository.FindPokinByJenisPohon(ctx, tx, "Tematik", 0, tahun, "")
	if err != nil {
		return nil, err
	}

	if len(pokins) == 0 {
		return nil, nil
	}

	var result []pohonkinerja.PohonKinerjaAdminResponseData
	for _, pokin := range pokins {
		result = append(result, pohonkinerja.PohonKinerjaAdminResponseData{
			Id:         pokin.Id,
			NamaPohon:  pokin.NamaPohon,
			JenisPohon: pokin.JenisPohon,
			LevelPohon: pokin.LevelPohon,
			Tahun:      pokin.Tahun,
		})
	}

	return result, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindPokinByStrategic(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD jika kodeOpd tidak kosong
	if kodeOpd != "" {
		_, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
		if err != nil {
			return nil, errors.New("kode opd tidak ditemukan")
		}
	}

	// Ambil data pohon kinerja dengan jenis "Strategic" dan level 4
	pokins, err := service.pohonKinerjaRepository.FindPokinByJenisPohon(ctx, tx, "Strategic", 4, tahun, kodeOpd)
	if err != nil {
		return nil, err
	}

	if len(pokins) == 0 {
		return nil, nil
	}

	var result []pohonkinerja.PohonKinerjaAdminResponseData
	for _, pokin := range pokins {
		// Ambil data OPD jika ada kodeOpd
		var namaOpd string
		if pokin.KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
			if err == nil {
				namaOpd = opd.NamaOpd
			}
		}

		result = append(result, pohonkinerja.PohonKinerjaAdminResponseData{
			Id:         pokin.Id,
			Parent:     pokin.Parent,
			NamaPohon:  pokin.NamaPohon,
			JenisPohon: pokin.JenisPohon,
			LevelPohon: pokin.LevelPohon,
			KodeOpd:    pokin.KodeOpd,
			NamaOpd:    namaOpd,
			Tahun:      pokin.Tahun,
		})
	}

	return result, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindPokinByTactical(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD jika kodeOpd tidak kosong
	if kodeOpd != "" {
		_, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
		if err != nil {
			return nil, errors.New("kode opd tidak ditemukan")
		}
	}

	// Ambil data pohon kinerja dengan jenis "Strategic" dan level 4
	pokins, err := service.pohonKinerjaRepository.FindPokinByJenisPohon(ctx, tx, "Tactical", 5, tahun, kodeOpd)
	if err != nil {
		return nil, err
	}

	if len(pokins) == 0 {
		return nil, nil
	}

	var result []pohonkinerja.PohonKinerjaAdminResponseData
	for _, pokin := range pokins {
		var namaOpd string
		if pokin.KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
			if err == nil {
				namaOpd = opd.NamaOpd
			}
		}

		result = append(result, pohonkinerja.PohonKinerjaAdminResponseData{
			Id:         pokin.Id,
			Parent:     pokin.Parent,
			NamaPohon:  pokin.NamaPohon,
			JenisPohon: pokin.JenisPohon,
			LevelPohon: pokin.LevelPohon,
			KodeOpd:    pokin.KodeOpd,
			NamaOpd:    namaOpd,
			Tahun:      pokin.Tahun,
		})
	}

	return result, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindPokinByOperational(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD jika kodeOpd tidak kosong
	if kodeOpd != "" {
		_, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
		if err != nil {
			return nil, errors.New("kode opd tidak ditemukan")
		}
	}

	// Ambil data pohon kinerja dengan jenis "Strategic" dan level 4
	pokins, err := service.pohonKinerjaRepository.FindPokinByJenisPohon(ctx, tx, "Operational", 6, tahun, kodeOpd)
	if err != nil {
		return nil, err
	}

	if len(pokins) == 0 {
		return nil, nil
	}

	var result []pohonkinerja.PohonKinerjaAdminResponseData
	for _, pokin := range pokins {
		// Ambil data OPD jika ada kodeOpd
		var namaOpd string
		if pokin.KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
			if err == nil {
				namaOpd = opd.NamaOpd
			}
		}

		result = append(result, pohonkinerja.PohonKinerjaAdminResponseData{
			Id:         pokin.Id,
			Parent:     pokin.Parent,
			NamaPohon:  pokin.NamaPohon,
			JenisPohon: pokin.JenisPohon,
			LevelPohon: pokin.LevelPohon,
			KodeOpd:    pokin.KodeOpd,
			NamaOpd:    namaOpd,
			Tahun:      pokin.Tahun,
		})
	}

	return result, nil
}
