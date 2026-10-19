package service

import (
	"bytes"
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/laporanreviewtujuanopd"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-pdf/fpdf"
)

type ReviewServiceImpl struct {
	ReviewRepository       repository.ReviewRepository
	DB                     *sql.DB
	PohonKinerjaRepository repository.PohonKinerjaRepository
	pegawaiRepository      repository.PegawaiRepository
}

func NewReviewServiceImpl(reviewRepository repository.ReviewRepository, db *sql.DB, pohonkinerjaRepository repository.PohonKinerjaRepository, pegawaiRepository repository.PegawaiRepository) *ReviewServiceImpl {
	return &ReviewServiceImpl{
		ReviewRepository:       reviewRepository,
		DB:                     db,
		PohonKinerjaRepository: pohonkinerjaRepository,
		pegawaiRepository:      pegawaiRepository,
	}
}

func (service *ReviewServiceImpl) Create(ctx context.Context, request pohonkinerja.ReviewCreateRequest) (pohonkinerja.ReviewResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}
	defer tx.Rollback()

	// Mendapatkan claims dari context
	claims, ok := ctx.Value(helper.UserInfoKey).(web.JWTClaim)
	if !ok {
		return pohonkinerja.ReviewResponse{}, errors.New("unauthorized: invalid user info in context")
	}
	if claims.Nip == "" {
		return pohonkinerja.ReviewResponse{}, errors.New("unauthorized: NIP tidak ditemukan")
	}

	err = service.PohonKinerjaRepository.ValidatePokinId(ctx, tx, request.IdPohonKinerja)
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}

	// Generate random ID
	randomId := rand.Intn(1000000)

	_, err = service.ReviewRepository.FindById(ctx, tx, randomId)
	for err == nil {
		randomId = rand.Intn(1000000)
		_, err = service.ReviewRepository.FindById(ctx, tx, randomId)
	}

	review := domain.Review{
		Id:             randomId,
		IdPohonKinerja: request.IdPohonKinerja,
		Review:         request.Review,
		Keterangan:     request.Keterangan,
		CreatedBy:      claims.Nip,
		Jenis_pokin:    request.JenisPokin,
		Tahun:          request.Tahun,
	}

	// Simpan ke database
	result, err := service.ReviewRepository.Create(ctx, tx, review)
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}

	// Konversi hasil ke response
	response := pohonkinerja.ReviewResponse{
		Id:             result.Id,
		IdPohonKinerja: result.IdPohonKinerja,
		Review:         result.Review,
		Keterangan:     result.Keterangan,
		CreatedBy:      result.CreatedBy,
		JenisPokin:     result.Jenis_pokin,
	}

	return response, nil
}

func (service *ReviewServiceImpl) Update(ctx context.Context, request pohonkinerja.ReviewUpdateRequest) (pohonkinerja.ReviewResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}
	defer tx.Rollback()

	// Cek apakah review ada
	_, err = service.ReviewRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.ReviewResponse{}, errors.New("review tidak ditemukan")
	}

	review := domain.Review{
		Id:         request.Id,
		Review:     request.Review,
		Keterangan: request.Keterangan,
	}

	result, err := service.ReviewRepository.Update(ctx, tx, review)
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}

	response := pohonkinerja.ReviewResponse{
		Id:             result.Id,
		IdPohonKinerja: result.IdPohonKinerja,
		Review:         result.Review,
		Keterangan:     result.Keterangan,
		CreatedBy:      result.CreatedBy,
	}

	return response, nil
}

func (service *ReviewServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Cek apakah review ada
	_, err = service.ReviewRepository.FindById(ctx, tx, id)
	if err != nil {
		return errors.New("review tidak ditemukan")
	}

	err = service.ReviewRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (service *ReviewServiceImpl) FindAll(ctx context.Context, idPohonKinerja int) ([]pohonkinerja.ReviewResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	reviews, err := service.ReviewRepository.FindByPohonKinerja(ctx, tx, idPohonKinerja)
	if err != nil {
		return nil, err
	}

	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, reviews[0].CreatedBy)
	if err != nil {
		return nil, err
	}

	var reviewResponses []pohonkinerja.ReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
			Id:             review.Id,
			IdPohonKinerja: review.IdPohonKinerja,
			Review:         review.Review,
			Keterangan:     review.Keterangan,
			// CreatedBy:      review.CreatedBy,
			JenisPokin:  review.Jenis_pokin,
			NamaPegawai: pegawai.NamaPegawai,
		})
	}

	return reviewResponses, nil
}

func (service *ReviewServiceImpl) FindById(ctx context.Context, id int) (pohonkinerja.ReviewResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}
	defer tx.Rollback()

	review, err := service.ReviewRepository.FindById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		return pohonkinerja.ReviewResponse{}, err
	}

	response := pohonkinerja.ReviewResponse{
		Id:             review.Id,
		IdPohonKinerja: review.IdPohonKinerja,
		Review:         review.Review,
		Keterangan:     review.Keterangan,
		CreatedBy:      review.CreatedBy,
		JenisPokin:     review.Jenis_pokin,
	}

	return response, nil
}

func (service *ReviewServiceImpl) FindAllReviewByTematik(ctx context.Context, tahun string) ([]pohonkinerja.ReviewTematikResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tahun
	if tahun == "" {
		return nil, errors.New("tahun harus diisi")
	}

	reviews, err := service.ReviewRepository.FindAllReviewByTematik(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var response []pohonkinerja.ReviewTematikResponse
	for _, tematik := range reviews {
		var reviewDetails []pohonkinerja.ReviewDetailResponse
		for _, review := range tematik.Review {
			pegawai, _ := service.pegawaiRepository.FindByNip(ctx, tx, review.CreatedBy)

			reviewDetails = append(reviewDetails, pohonkinerja.ReviewDetailResponse{
				IdPohon:     review.IdPohon,
				Parent:      review.Parent,
				NamaPohon:   review.NamaPohon,
				LevelPohon:  review.LevelPohon,
				JenisPohon:  review.JenisPohon,
				Review:      review.Review,
				Keterangan:  review.Keterangan,
				NamaPegawai: pegawai.NamaPegawai,
				CreatedAt:   review.CreatedAt,
				UpdatedAt:   review.UpdatedAt,
			})
		}

		response = append(response, pohonkinerja.ReviewTematikResponse{
			IdTematik:  tematik.IdTematik,
			NamaPohon:  tematik.NamaPohon,
			LevelPohon: tematik.LevelPohon,
			Review:     reviewDetails,
		})
	}

	return response, nil
}

func (service *ReviewServiceImpl) FindAllReviewOpd(ctx context.Context, kodeOpd, tahun string) ([]pohonkinerja.ReviewOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	if kodeOpd == "" {
		return []pohonkinerja.ReviewOpdResponse{}, nil
	}
	if tahun == "" {
		return []pohonkinerja.ReviewOpdResponse{}, nil
	}

	reviews, err := service.ReviewRepository.FindAllReviewOpd(ctx, tx, kodeOpd, tahun)
	if err != nil {
		if err == sql.ErrNoRows {
			return []pohonkinerja.ReviewOpdResponse{}, nil
		}
		return nil, err
	}

	var reviewResponses []pohonkinerja.ReviewOpdResponse
	for _, review := range reviews {
		pegawai, _ := service.pegawaiRepository.FindByNip(ctx, tx, review.CreatedBy)

		reviewResponses = append(reviewResponses, pohonkinerja.ReviewOpdResponse{
			IdPohon:     review.IdPohon,
			Parent:      review.Parent,
			NamaPohon:   review.NamaPohon,
			LevelPohon:  review.LevelPohon,
			JenisPohon:  review.JenisPohon,
			Review:      review.Review,
			Keterangan:  review.Keterangan,
			NamaPegawai: pegawai.NamaPegawai, // Akan kosong jika pegawai tidak ditemukan
			CreatedAt:   review.CreatedAt,
			UpdatedAt:   review.UpdatedAt,
		})
	}

	if len(reviewResponses) == 0 {
		return []pohonkinerja.ReviewOpdResponse{}, nil
	}

	return reviewResponses, nil
}

func toReviewTujuanOpdResponse(review domain.Review, namaPegawai string) pohonkinerja.ReviewTujuanOpdResponse {
	return pohonkinerja.ReviewTujuanOpdResponse{
		Id:             review.Id,
		IdTujuanOpd:    review.IdTujuanOpd,
		IdPohonKinerja: review.IdPohonKinerja,
		Review:         review.Review,
		Keterangan:     review.Keterangan,
		CreatedBy:      review.CreatedBy,
		NamaPegawai:    namaPegawai,
	}
}

func (service *ReviewServiceImpl) CreateTujuanOpd(ctx context.Context, request pohonkinerja.ReviewTujuanOpdCreateRequest) (pohonkinerja.ReviewTujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, err
	}
	defer tx.Rollback()

	claims, ok := ctx.Value(helper.UserInfoKey).(web.JWTClaim)
	if !ok {
		return pohonkinerja.ReviewTujuanOpdResponse{}, errors.New("unauthorized: invalid user info in context")
	}
	if claims.Nip == "" {
		return pohonkinerja.ReviewTujuanOpdResponse{}, errors.New("unauthorized: NIP tidak ditemukan")
	}
	if request.IdTujuanOpd <= 0 {
		return pohonkinerja.ReviewTujuanOpdResponse{}, errors.New("id_tujuan_opd harus diisi")
	}

	randomId := rand.Intn(1000000)
	_, err = service.ReviewRepository.FindById(ctx, tx, randomId)
	for err == nil {
		randomId = rand.Intn(1000000)
		_, err = service.ReviewRepository.FindById(ctx, tx, randomId)
	}

	review := domain.Review{
		Id:             randomId,
		IdPohonKinerja: 0,
		IdTujuanOpd:    request.IdTujuanOpd,
		Review:         request.Review,
		Keterangan:     request.Keterangan,
		Catatan:        request.Catatan,
		CreatedBy:      claims.Nip,
		Tahun:          request.Tahun,
	}

	result, err := service.ReviewRepository.CreateTujuanOpd(ctx, tx, review)
	if err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, err
	}
	return toReviewTujuanOpdResponse(result, ""), nil
}

func (service *ReviewServiceImpl) UpdateTujuanOpd(ctx context.Context, request pohonkinerja.ReviewTujuanOpdUpdateRequest) (pohonkinerja.ReviewTujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, err
	}
	defer tx.Rollback()

	existing, err := service.ReviewRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, errors.New("review tidak ditemukan")
	}

	result, err := service.ReviewRepository.Update(ctx, tx, domain.Review{
		Id:         request.Id,
		Review:     request.Review,
		Keterangan: request.Keterangan,
		Catatan:    request.Catatan,
	})
	if err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, err
	}
	result.IdPohonKinerja = existing.IdPohonKinerja
	result.IdTujuanOpd = existing.IdTujuanOpd
	result.CreatedBy = existing.CreatedBy
	return toReviewTujuanOpdResponse(result, ""), nil
}

func (service *ReviewServiceImpl) DeleteTujuanOpd(ctx context.Context, id int) error {
	return service.Delete(ctx, id)
}

func (service *ReviewServiceImpl) FindAllTujuanOpd(ctx context.Context, idTujuanOpd int) ([]pohonkinerja.ReviewTujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	reviews, err := service.ReviewRepository.FindByTujuanOpd(ctx, tx, idTujuanOpd)
	if err != nil {
		return nil, err
	}

	responses := make([]pohonkinerja.ReviewTujuanOpdResponse, 0, len(reviews))
	for _, review := range reviews {
		namaPegawai := ""
		if review.CreatedBy != "" {
			if pegawai, errPegawai := service.pegawaiRepository.FindByNip(ctx, tx, review.CreatedBy); errPegawai == nil {
				namaPegawai = pegawai.NamaPegawai
			}
		}
		responses = append(responses, toReviewTujuanOpdResponse(review, namaPegawai))
	}
	return responses, nil
}

func (service *ReviewServiceImpl) FindByIdTujuanOpd(ctx context.Context, id int) (pohonkinerja.ReviewTujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	review, err := service.ReviewRepository.FindById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.ReviewTujuanOpdResponse{}, errors.New("review tidak ditemukan")
	}
	namaPegawai := ""
	if review.CreatedBy != "" {
		if pegawai, errPegawai := service.pegawaiRepository.FindByNip(ctx, tx, review.CreatedBy); errPegawai == nil {
			namaPegawai = pegawai.NamaPegawai
		}
	}
	return toReviewTujuanOpdResponse(review, namaPegawai), nil
}


func (service *ReviewServiceImpl) FindAllLaporanTujuanOpd(ctx context.Context, kodeOpd string) ([]laporanreviewtujuanopd.LaporanReviewTujuanOpd, error) {

	tx, err := service.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	results, err := service.ReviewRepository.FindLaporanReviewTujuanOpd(
		ctx, tx, kodeOpd,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]laporanreviewtujuanopd.LaporanReviewTujuanOpd,
		0,
	)

	for _, result := range results {
		reviewTujuan := laporanreviewtujuanopd.LaporanReviewTujuan{}

		for _, review := range result.Review {
			detail := &laporanreviewtujuanopd.LaporanReviewTujuanDetail{
				Id:             &review.Id,
				IdTujuanOpd:    &review.IdTujuanOpd,
				IdPohonKinerja: &review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				Catatan:        review.Catatan,
				CreatedBy:      &review.CreatedBy,
				NamaPegawai:    &review.NamaPegawai,
				Tahun:          &review.Tahun,
			}

			switch strings.ToLower(strings.TrimSpace(review.Review)) {
			case "validasi":
				reviewTujuan.Validasi = detail
			case "kesesuaian":
				reviewTujuan.Kesesuaian = detail
			case "kelengkapan":
				reviewTujuan.Kelengkapan = detail
			case "catatan":
				reviewTujuan.Catatan = detail
			}
		}

		
		// Jika validasi Ideal, hasil kesesuaian dan kelengkapan otomatis.
		if reviewTujuan.Validasi != nil &&
			strings.EqualFold(
				strings.TrimSpace(reviewTujuan.Validasi.Keterangan),
				"Ideal",
			) {

			reviewTujuan.Kesesuaian = &laporanreviewtujuanopd.LaporanReviewTujuanDetail{
				Review:     "KESESUAIAN",
				Keterangan: "Sesuai",
			}

			reviewTujuan.Kelengkapan = &laporanreviewtujuanopd.LaporanReviewTujuanDetail{
				Review:     "KELENGKAPAN",
				Keterangan: "Lengkap",
			}

			reviewTujuan.Catatan = nil
		} else {
			// Gabungkan catatan kesesuaian dan kelengkapan.
			var catatanList []string

			if reviewTujuan.Kesesuaian != nil &&
				strings.TrimSpace(reviewTujuan.Kesesuaian.Catatan) != "" {
				catatanList = append(
					catatanList,
					reviewTujuan.Kesesuaian.Catatan,
				)
			}

			if reviewTujuan.Kelengkapan != nil &&
				strings.TrimSpace(reviewTujuan.Kelengkapan.Catatan) != "" {
				catatanList = append(
					catatanList,
					reviewTujuan.Kelengkapan.Catatan,
				)
			}

			if len(catatanList) > 0 {
				reviewTujuan.Catatan = &laporanreviewtujuanopd.LaporanReviewTujuanDetail{
					Review:  "CATATAN",
					Catatan: strings.Join(catatanList, "; "),
				}
			} else {
				reviewTujuan.Catatan = nil
			}
		}

		responses = append(responses, laporanreviewtujuanopd.LaporanReviewTujuanOpd{
			Id:               result.Id,
			KodeOpd:          result.KodeOpd,
			NamaOpd:          result.NamaOpd,
			KodeBidangUrusan: result.KodeBidangUrusan,
			NamaBidangUrusan: result.NamaBidangUrusan,
			Tujuan:           result.Tujuan,
			Review:           reviewTujuan,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	committed = true

	return responses, nil
}


func (service *ReviewServiceImpl) CetakPDFReviewTujuanOpd(ctx context.Context, kodeOpd string) ([]byte, error) {
	if strings.TrimSpace(kodeOpd) == "" {
		return nil, fmt.Errorf("kode OPD tidak boleh kosong")
	}

	data, err := service.FindAllLaporanTujuanOpd(ctx, kodeOpd)
	if err != nil {
		return nil, err
	}

	// =========================
	// MEMBUAT PDF F4 LANDSCAPE
	// =========================
	pdf := fpdf.NewCustom(&fpdf.InitType{
		OrientationStr: "L",
		UnitStr:        "mm",
		SizeStr:        "",
		Size: fpdf.SizeType{
			Wd: 210,
			Ht: 330,
		},
		FontDirStr: "",
	})

	pdf.SetTitle("Laporan Hasil Analisa Tujuan OPD", false)
	pdf.SetAuthor("Review Tujuan OPD", false)
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)

	pdf.AddPage()

	// =========================
	// JUDUL LAPORAN
	// =========================
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(
		330,
		8,
		"Laporan Hasil Analisa Tujuan OPD",
		"",
		1,
		"C",
		false,
		0,
		"",
	)

	pdf.Ln(7)

	// =========================
	// LEBAR KOLOM
	// =========================
	colWidths := []float64{
		10, // No
		45, // Nama OPD
		45, // Bidang Urusan
		30, // Lengkap
		30, // Tidak Lengkap
		45, // Sesuai
		45, // Belum Sesuai
		60, // Catatan
	}

	tableWidth := 0.0
	for _, width := range colWidths {
		tableWidth += width
	}

	// Lebar halaman F4 landscape 330 mm.
	// Margin kiri dan kanan masing-masing 10 mm.
	// Lebar tabel disesuaikan agar tidak melewati area cetak.
	if tableWidth > 310 {
		return nil, fmt.Errorf("total lebar tabel melebihi area cetak")
	}

	startX := (330 - tableWidth) / 2

	// =========================
	// FUNGSI HEADER TABEL
	// =========================
	renderHeader := func() {
		pdf.SetFont("Arial", "B", 9)

		startY := pdf.GetY()
		x := startX
		groupHeight := 7.0
		subHeight := 7.0
		numberHeight := 6.0
		totalHeaderHeight := groupHeight + subHeight + numberHeight

		// Kolom No.
		pdf.Rect(x, startY, colWidths[0], totalHeaderHeight, "D")
		pdf.SetXY(x, startY)
		pdf.CellFormat(
			colWidths[0],
			totalHeaderHeight,
			"No.",
			"",
			0,
			"C",
			false,
			0,
			"",
		)
		x += colWidths[0]

		// Kolom Nama OPD
		pdf.Rect(x, startY, colWidths[1], totalHeaderHeight, "D")
		pdf.SetXY(x, startY)
		pdf.CellFormat(
			colWidths[1],
			totalHeaderHeight,
			"Nama OPD",
			"",
			0,
			"C",
			false,
			0,
			"",
		)
		x += colWidths[1]

		// Kolom Bidang Urusan
		pdf.Rect(x, startY, colWidths[2], totalHeaderHeight, "D")
		pdf.SetXY(x, startY)
		pdf.CellFormat(
			colWidths[2],
			totalHeaderHeight,
			"Bidang Urusan",
			"",
			0,
			"C",
			false,
			0,
			"",
		)
		x += colWidths[2]

		// Kelompok Kelengkapan Tujuan OPD
		kelengkapanWidth := colWidths[3] + colWidths[4]
		pdf.Rect(x, startY, kelengkapanWidth, groupHeight, "D")
		pdf.SetXY(x, startY)
		pdf.CellFormat(
			kelengkapanWidth,
			groupHeight,
			"KELENGKAPAN TUJUAN OPD",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Subkolom Lengkap
		pdf.Rect(x, startY+groupHeight, colWidths[3], subHeight, "D")
		pdf.SetXY(x, startY+groupHeight)
		pdf.CellFormat(
			colWidths[3],
			subHeight,
			"Lengkap",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Subkolom Tidak Lengkap
		pdf.Rect(
			x+colWidths[3],
			startY+groupHeight,
			colWidths[4],
			subHeight,
			"D",
		)
		pdf.SetXY(x+colWidths[3], startY+groupHeight)
		pdf.CellFormat(
			colWidths[4],
			subHeight,
			"Tidak Lengkap",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Nomor referensi 4 dan 5
		pdf.SetFont("Arial", "I", 9)
		pdf.Rect(
			x,
			startY+groupHeight+subHeight,
			colWidths[3],
			numberHeight,
			"D",
		)
		pdf.SetXY(x, startY+groupHeight+subHeight)
		pdf.CellFormat(
			colWidths[3],
			numberHeight,
			"4",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		pdf.Rect(
			x+colWidths[3],
			startY+groupHeight+subHeight,
			colWidths[4],
			numberHeight,
			"D",
		)
		pdf.SetXY(x+colWidths[3], startY+groupHeight+subHeight)
		pdf.CellFormat(
			colWidths[4],
			numberHeight,
			"5",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		x += kelengkapanWidth

		// Kelompok Kesesuaian
		kesesuaianWidth := colWidths[5] + colWidths[6]
		pdf.SetFont("Arial", "B", 9)
		pdf.Rect(x, startY, kesesuaianWidth, groupHeight, "D")
		pdf.SetXY(x, startY)
		pdf.CellFormat(
			kesesuaianWidth,
			groupHeight,
			"KESESUAIAN DENGAN INMENDAGRI NO 2 TAHUN 2025",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Subkolom Sesuai
		pdf.Rect(x, startY+groupHeight, colWidths[5], subHeight, "D")
		pdf.SetXY(x, startY+groupHeight)
		pdf.CellFormat(
			colWidths[5],
			subHeight,
			"Sesuai",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Subkolom Belum Sesuai
		pdf.Rect(
			x+colWidths[5],
			startY+groupHeight,
			colWidths[6],
			subHeight,
			"D",
		)
		pdf.SetXY(x+colWidths[5], startY+groupHeight)
		pdf.CellFormat(
			colWidths[6],
			subHeight,
			"Belum Sesuai",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Nomor referensi 6 dan 7
		pdf.SetFont("Arial", "I", 9)
		pdf.Rect(
			x,
			startY+groupHeight+subHeight,
			colWidths[5],
			numberHeight,
			"D",
		)
		pdf.SetXY(x, startY+groupHeight+subHeight)
		pdf.CellFormat(
			colWidths[5],
			numberHeight,
			"6",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		pdf.Rect(
			x+colWidths[5],
			startY+groupHeight+subHeight,
			colWidths[6],
			numberHeight,
			"D",
		)
		pdf.SetXY(x+colWidths[5], startY+groupHeight+subHeight)
		pdf.CellFormat(
			colWidths[6],
			numberHeight,
			"7",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		x += kesesuaianWidth

				// Kolom Catatan
		pdf.SetFont("Arial", "B", 9)

		// Judul CATATAN
		pdf.Rect(
			x,
			startY,
			colWidths[7],
			groupHeight+subHeight,
			"D",
		)
		pdf.SetXY(x, startY)
		pdf.CellFormat(
			colWidths[7],
			groupHeight+subHeight,
			"CATATAN",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Nomor kolom 8
		pdf.SetFont("Arial", "I", 9)
		numberY := startY + groupHeight + subHeight

		pdf.Rect(
			x,
			numberY,
			colWidths[7],
			numberHeight,
			"D",
		)
		pdf.SetXY(x, numberY)
		pdf.CellFormat(
			colWidths[7],
			numberHeight,
			"8",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Nomor referensi 1, 2, dan 3
		pdf.SetFont("Arial", "I", 9)

		numberY = startY + groupHeight + subHeight

		// Nomor kolom 1 - No.
		pdf.Rect(
			startX,
			numberY,
			colWidths[0],
			numberHeight,
			"D",
		)
		pdf.SetXY(startX, numberY)
		pdf.CellFormat(
			colWidths[0],
			numberHeight,
			"1",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Nomor kolom 2 - Nama OPD
		xNomor2 := startX + colWidths[0]
		pdf.Rect(
			xNomor2,
			numberY,
			colWidths[1],
			numberHeight,
			"D",
		)
		pdf.SetXY(xNomor2, numberY)
		pdf.CellFormat(
			colWidths[1],
			numberHeight,
			"2",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Nomor kolom 3 - Bidang Urusan
		xNomor3 := startX + colWidths[0] + colWidths[1]
		pdf.Rect(
			xNomor3,
			numberY,
			colWidths[2],
			numberHeight,
			"D",
		)
		pdf.SetXY(xNomor3, numberY)
		pdf.CellFormat(
			colWidths[2],
			numberHeight,
			"3",
			"",
			0,
			"C",
			false,
			0,
			"",
		)

		// Posisi awal data
		pdf.SetXY(startX, startY+totalHeaderHeight)
	}

	renderHeader()

	// =========================
	// FUNGSI MENGAMBIL STATUS
	// =========================
	getKeterangan := func(
		detail *laporanreviewtujuanopd.LaporanReviewTujuanDetail,
	) string {
		if detail == nil {
			return ""
		}
		return strings.TrimSpace(detail.Keterangan)
	}

	getCatatan := func(
		detail *laporanreviewtujuanopd.LaporanReviewTujuanDetail,
	) string {
		if detail == nil {
			return ""
		}
		return strings.TrimSpace(detail.Catatan)
	}

	// =========================
	// ISI TABEL
	// =========================
	pdf.SetFont("Arial", "", 8)

	for i, item := range data {
		// Status kelengkapan
		kelengkapan := getKeterangan(item.Review.Kelengkapan)
		lengkap := ""
		tidakLengkap := ""

		if strings.EqualFold(kelengkapan, "Lengkap") {
			lengkap = "Lengkap"
		} else if strings.EqualFold(kelengkapan, "Tidak Lengkap") {
			tidakLengkap = "Tidak Lengkap"
		}

		// Status kesesuaian
		kesesuaian := getKeterangan(item.Review.Kesesuaian)
		sesuai := ""
		belumSesuai := ""

		if strings.EqualFold(kesesuaian, "Sesuai") {
			sesuai = "Sesuai"
		} else if strings.EqualFold(kesesuaian, "Tidak Sesuai") {
			belumSesuai = "TIdak Sesuai"
		}

		catatan := getCatatan(item.Review.Catatan)

		row := []string{
			fmt.Sprintf("%d", i+1),
			item.NamaOpd,
			item.NamaBidangUrusan,
			lengkap,
			tidakLengkap,
			sesuai,
			belumSesuai,
			catatan,
		}

		// Hitung tinggi baris sesuai teks terpanjang
		lineHeight := 5.0
		maxLines := 1

		for col, value := range row {
			lines := pdf.SplitLines(
				[]byte(value),
				colWidths[col]-2,
			)
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
		}

		rowHeight := float64(maxLines)*lineHeight + 2

		// Jika tidak cukup ruang, pindah halaman
		if pdf.GetY()+rowHeight > 200 {
			pdf.AddPage()
			renderHeader()
			pdf.SetFont("Arial", "", 8)
		}

		y := pdf.GetY()
		x := startX

		for col, value := range row {
			width := colWidths[col]

			// Gambar border sel
			pdf.Rect(x, y, width, rowHeight, "D")

			// Isi teks sel
			pdf.SetXY(x+1, y+1)

			align := "L"
			if col == 0 || (col >= 3 && col <= 6) {
				align = "C"
			}

			pdf.MultiCell(
				width-2,
				lineHeight,
				value,
				"",
				align,
				false,
			)

			x += width
		}

		pdf.SetXY(startX, y+rowHeight)
	}

	// =========================
	// OUTPUT PDF
	// =========================
	var buffer bytes.Buffer

	err = pdf.Output(&buffer)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat PDF: %w", err)
	}

	return buffer.Bytes(), nil
}

