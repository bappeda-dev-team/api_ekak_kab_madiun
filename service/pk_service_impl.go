package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pegawai"
	"ekak_kabupaten_madiun/model/web/pkopd"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type PkServiceImpl struct {
	pkOpdRepository              repository.PkRepository
	pegawaiService               PegawaiService
	rekinService                 RencanaKinerjaService
	opdService                   OpdService
	strukturOrganisasiRepository repository.StrukturOrganisasiRepository
	Validate                     *validator.Validate
	DB                           *sql.DB
}

func NewPkServiceImpl(
	pkOpdRepository repository.PkRepository,
	pegawaiService PegawaiService,
	rekinService RencanaKinerjaService,
	opdService OpdService,
	strukturOrganisasiRepository repository.StrukturOrganisasiRepository,
	validate *validator.Validate,
	DB *sql.DB,
) *PkServiceImpl {
	return &PkServiceImpl{
		pkOpdRepository:              pkOpdRepository,
		pegawaiService:               pegawaiService,
		rekinService:                 rekinService,
		opdService:                   opdService,
		strukturOrganisasiRepository: strukturOrganisasiRepository,
		Validate:                     validate,
		DB:                           DB,
	}
}

func (service *PkServiceImpl) FindByKodeOpdTahun(ctx context.Context, kodeOpd string, tahun int) (pkopd.PkOpdResponse, error) {
	log.Printf("[INFO] PK OPD FIND BY KODE OPD TAHUN")
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return pkopd.PkOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// cek opd dulu
	opd, err := service.opdService.FindByKodeOpd(ctx, kodeOpd)
	if err != nil {
		log.Printf("[ERROR] Find OPD by kodeOpd: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("OPD TIDAK DITEMUKAN")
	}
	// base info nama opd dan kepala opd
	namaOpd := opd.NamaOpd
	kepalaOpd := opd.NamaKepalaOpd
	nipKepalaOpd := opd.NIPKepalaOpd
	// end check opd

	// all pegawai in opd
	pegawais, err := service.pegawaiService.FindAll(ctx, kodeOpd)
	if err != nil {
		log.Printf("[ERROR] Find Pegawai kodeOpd: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("terjadi kesalahan sistem")
	}
	// rekin in opd by tahun
	// filter params
	filterParams := domain.FilterParams{
		"kode_opd": kodeOpd,
		"tahun":    strconv.Itoa(tahun),
	}

	log.Printf("FILTER PARAMS: %v \n", filterParams)
	rekins, err := service.rekinService.FindByFilter(ctx, filterParams)
	if err != nil {
		log.Printf("[ERROR] Find Rekin by kodeOpd and tahun: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("terjadi kesalahan sistem")
	}

	// data struktur untuk penyusunan
	// lookup pegawai by nip untuk susun nama atasan
	pegawaiByNip := make(map[string]pegawai.PegawaiResponse)

	// pegawaiId = nip
	rekinByPegawaiId := make(map[string][]rencanakinerja.RencanaKinerjaResponse)
	// agar jika tidak ada rekin, bisa empty
	for _, peg := range pegawais {
		pegawaiByNip[peg.Nip] = peg
		rekinByPegawaiId[peg.Nip] = []rencanakinerja.RencanaKinerjaResponse{}
	}
	rekinById := make(map[string]rencanakinerja.RencanaKinerjaResponse)

	for _, rekin := range rekins {
		rekinById[rekin.Id] = rekin

		pegawaiId := rekin.PegawaiId // nip

		// skip kalau pegawai tidak ada (defensive)
		if _, exists := rekinByPegawaiId[pegawaiId]; !exists {
			continue
		}

		rekinByPegawaiId[pegawaiId] = append(
			rekinByPegawaiId[pegawaiId],
			rencanakinerja.RencanaKinerjaResponse{
				Id:                 rekin.Id,
				IdPohon:            rekin.IdPohon,
				IdParentPohon:      rekin.IdParentPohon,
				LevelPohon:         rekin.LevelPohon,
				NamaRencanaKinerja: rekin.NamaRencanaKinerja,
				NamaPegawai:        rekin.NamaPegawai,
				PegawaiId:          rekin.PegawaiId,
				Indikator:          rekin.Indikator,
			},
		)
	}

	// pk yang sudah tersimpan di opd dan tahun
	// grouping by level
	pkOpds, err := service.pkOpdRepository.FindByKodeOpdTahun(ctx, tx, kodeOpd, tahun)
	if err != nil {
		log.Printf("[ERROR] Find PK OPD by kodeOpd and tahun: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("terjadi kesalahan sistem")
	}

	// atasan by pegawai
	// grouped by id pegawai (nip)
	atasans, err := service.strukturOrganisasiRepository.AtasanBawahanByKodeOpdTahun(ctx, tx, kodeOpd, tahun)
	if err != nil {
		log.Printf("[ERROR] Find Struktur Organisasi: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("terjadi kesalahan sistem")
	}

	// processing DTO
	// merge pegawais dan rekins
	pkByRekinPemilik := make(map[string]domain.PkOpd)

	for _, pkList := range pkOpds {
		for _, pk := range pkList {
			if pk.IdRekinPemilikPk != "" {
				pkByRekinPemilik[pk.IdRekinPemilikPk] = pk
			}
		}
	}
	// DTO PK yang sudah pilih rekin atasan
	// level -> nip -> pegawai node
	pkByLevel := make(map[int]map[string]*pkopd.PkPegawai)

	for _, rekin := range rekins {

		level := rekin.LevelPohon
		nip := rekin.PegawaiId
		nama := rekin.NamaPegawai
		nipAtasan := atasans[nip]
		namaAtasan := ""
		jabatanAtasan := ""

		if nipAtasan != "" {
			if peg, ok := pegawaiByNip[nipAtasan]; ok {
				namaAtasan = peg.NamaPegawai
				jabatanAtasan = peg.NamaJabatan
			}
		}

		// init level jika belum ada
		if _, ok := pkByLevel[level]; !ok {
			pkByLevel[level] = make(map[string]*pkopd.PkPegawai)
		}

		// data pegawai
		jabatanPegawai := pegawaiByNip[nip].NamaJabatan

		// init pegawai jika belum ada
		// input atasan sekalian kalau ada
		if _, ok := pkByLevel[level][nip]; !ok {
			pkByLevel[level][nip] = &pkopd.PkPegawai{
				JenisItem:      translateJenisItem(level),
				NipAtasan:      nipAtasan,
				NamaAtasan:     namaAtasan,
				JabatanAtasan:  jabatanAtasan,
				Nip:            nip,
				Nama:           nama,
				JabatanPegawai: jabatanPegawai,
				Pks:            []pkopd.PkAsn{},
				Subkegiatan:    []pkopd.SubkegiatanPk{},
			}
		}
		indikatorPk := []pkopd.IndikatorPk{}
		for _, ind := range rekin.Indikator {
			targetIndPk := []pkopd.TargetIndPk{}

			for _, tar := range ind.Target {
				targetIndPk = append(targetIndPk, pkopd.TargetIndPk{
					Target: tar.TargetIndikator,
					Satuan: tar.SatuanIndikator,
				})
			}
			indikatorPk = append(indikatorPk, pkopd.IndikatorPk{
				Indikator: ind.NamaIndikator,
				Targets:   targetIndPk,
			})
		}

		// default PK (BELUM ADA)
		pkAsn := pkopd.PkAsn{
			Id:               "",
			IdPohon:          rekin.IdPohon,
			IdParentPohon:    rekin.IdParentPohon,
			KodeOpd:          rekin.KodeOpd.KodeOpd,
			NamaOpd:          rekin.KodeOpd.NamaOpd,
			LevelPk:          rekin.LevelPohon,
			IdRekinPemilikPk: rekin.Id,
			RekinPemilikPk:   rekin.NamaRencanaKinerja,
			NipPemilikPk:     rekin.PegawaiId,
			NamaPemilikPk:    rekin.NamaPegawai,
			Tahun:            tahun,
			Indikators:       indikatorPk,
			PaguAnggaran:     0,
		}

		// enrich dari PK jika ada
		if pk, ok := pkByRekinPemilik[rekin.Id]; ok {
			pkAsn.Id = pk.Id
			pkAsn.NipAtasan = pk.NipAtasan
			pkAsn.NamaAtasan = pk.NamaAtasan
			pkAsn.IdRekinAtasan = pk.IdRekinAtasan
			pkAsn.Keterangan = pk.Keterangan

			if rAtasan, ok := rekinById[pk.IdRekinAtasan]; ok {
				pkAsn.RekinAtasan = rAtasan.NamaRencanaKinerja
			}
		}

		// append PK ke pegawai
		pkByLevel[level][nip].Pks = append(
			pkByLevel[level][nip].Pks,
			pkAsn,
		)
	}

	// sort rekin
	for _, pegawaiMap := range pkByLevel {
		for _, peg := range pegawaiMap {
			sort.Slice(peg.Pks, func(i, j int) bool {
				return peg.Pks[i].IdRekinPemilikPk <
					peg.Pks[j].IdRekinPemilikPk
			})
		}
	}

	pkItems := make([]pkopd.PkOpdByLevel, 0, len(pkByLevel))

	// 1. ambil semua level
	levels := make([]int, 0, len(pkByLevel))
	for level := range pkByLevel {
		// skip empty level
		// artinya pokin dari rekin tersebut tidak ditemukan
		if level == 0 {
			continue
		}
		levels = append(levels, level)
	}

	// 2. sort ascending
	sort.Ints(levels)

	// 3. build response sesuai urutan level
	for _, level := range levels {
		pegawaiMap := pkByLevel[level]

		pegawais := make([]pkopd.PkPegawai, 0, len(pegawaiMap))
		for _, peg := range pegawaiMap {
			pegawais = append(pegawais, *peg)
		}

		sort.Slice(pegawais, func(i, j int) bool {
			return pegawais[i].Nip < pegawais[j].Nip
		})

		pkItems = append(pkItems, pkopd.PkOpdByLevel{
			LevelPk:  level,
			Pegawais: pegawais,
		})
	}

	result := pkopd.PkOpdResponse{
		KodeOpd:      kodeOpd,
		NamaOpd:      namaOpd,
		KepalaOpd:    kepalaOpd,
		NipKepalaOpd: nipKepalaOpd,
		Tahun:        tahun,
		PkItem:       pkItems,
	}

	return result, nil
}

func (service *PkServiceImpl) HubungkanRekin(
	ctx context.Context,
	request pkopd.PkOpdRequest,
) (pkopd.PkOpdResponse, error) {

	// 1. validasi
	if err := service.Validate.Struct(request); err != nil {
		log.Printf("Invalid hubungkan rekin request: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("validasi gagal")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("gagal memulai transaksi")
	}
	// transaksi ini selalu commit di akhir
	// defer helper.CommitOrRollback(tx)

	kodeOpd := request.KodeOpd
	tahun := request.Tahun
	tahunStr := strconv.Itoa(tahun)

	// 2. ambil OPD
	opd, err := service.opdService.FindByKodeOpd(ctx, kodeOpd)
	if err != nil {
		log.Printf("[ERROR] Find OPD: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("OPD tidak ditemukan")
	}

	// 3. ambil rekin atasan
	rekinAtasan, err := service.rekinService.FindById(
		ctx,
		request.IdRekinAtasan,
		kodeOpd,
		tahunStr,
	)
	if err != nil {
		log.Printf("[ERROR] Find Rekin Atasan: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("rekin atasan tidak ditemukan")
	}

	// 4. ambil rekin pemilik PK
	rekinPemilik, err := service.rekinService.FindById(
		ctx,
		request.IdRekinPemilikPk,
		kodeOpd,
		tahunStr,
	)
	if err != nil {
		log.Printf("[ERROR] Find Rekin Pemilik: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("rekin pemilik tidak ditemukan")
	}

	// 5. bentuk domain PK OPD
	pk := domain.PkOpd{
		KodeOpd: kodeOpd,
		NamaOpd: opd.NamaOpd,
		LevelPk: request.LevelPk,
		Tahun:   tahun,

		NipAtasan:     rekinAtasan.PegawaiId,
		NamaAtasan:    rekinAtasan.NamaPegawai,
		IdRekinAtasan: rekinAtasan.Id,
		RekinAtasan:   rekinAtasan.NamaRencanaKinerja,

		NipPemilikPk:     rekinPemilik.PegawaiId,
		NamaPemilikPk:    rekinPemilik.NamaPegawai,
		IdRekinPemilikPk: rekinPemilik.Id,
		RekinPemilikPk:   rekinPemilik.NamaRencanaKinerja,
	}

	// 6. simpan / update relasi
	if err := service.pkOpdRepository.HubungkanRekin(ctx, tx, pk); err != nil {
		log.Printf("[ERROR] HubungkanRekin repo: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("gagal menghubungkan rekin")
	}

	// 6.1 Commit dulu
	if err := tx.Commit(); err != nil {
		return pkopd.PkOpdResponse{}, err
	}

	// 7. return FULL RESPONSE TERBARU
	return service.FindByKodeOpdTahun(ctx, kodeOpd, tahun)
}

func (service *PkServiceImpl) HubungkanAtasan(
	ctx context.Context,
	request pkopd.HubungkanAtasanRequest,
) (pkopd.PkOpdResponse, error) {

	// 1. validasi
	if err := service.Validate.Struct(request); err != nil {
		log.Printf("Invalid hubungkan atasan request: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("validasi gagal")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("gagal memulai transaksi")
	}
	// transaksi ini selalu commit di akhir
	// defer helper.CommitOrRollback(tx)

	kodeOpd := request.KodeOpd
	tahun := request.Tahun

	// 5. bentuk domain PK OPD
	strukturOrganisasi := domain.StrukturOrganisasi{
		KodeOpd:    kodeOpd,
		Tahun:      tahun,
		NipAtasan:  request.NipAtasan,
		NipBawahan: request.NipBawahan,
	}

	// 6. simpan / update relasi
	if err := service.strukturOrganisasiRepository.Create(ctx, tx, strukturOrganisasi); err != nil {
		log.Printf("[ERROR] HubungkanRekin repo: %v", err)
		return pkopd.PkOpdResponse{}, fmt.Errorf("gagal menghubungkan rekin")
	}

	// 6.1 Commit dulu
	if err := tx.Commit(); err != nil {
		return pkopd.PkOpdResponse{}, err
	}

	// 7. return FULL RESPONSE TERBARU
	return service.FindByKodeOpdTahun(ctx, kodeOpd, tahun)
}

func translateJenisItem(level int) string {
	switch level {
	case 4:
		return "Strategic"
	case 5:
		return "Tactical"
	case 6:
		return "Operational"
	case 7:
		return "Operational N"
	default:
		return "-"
	}
}
