package service

import (
	"bytes"
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

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

func (service *ReviewServiceImpl) CetakPDFReviewOpd(ctx context.Context, kodeOpd string, tahun string) ([]byte, error) {

	if kodeOpd == "" {
		return nil, fmt.Errorf("kode opd tidak boleh kosong")
	}

	if tahun == "" {
		return nil, fmt.Errorf("tahun tidak boleh kosong")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	reviews, err := service.ReviewRepository.FindAllReviewOpd(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	// Mengikuti response yang digunakan FE.
	reviewResponses := make([]pohonkinerja.ReviewOpdResponse, 0, len(reviews))

	for _, review := range reviews {

		pegawai, _ := service.pegawaiRepository.FindByNip(ctx, tx, review.CreatedBy)

		reviewResponses = append(
			reviewResponses,
			pohonkinerja.ReviewOpdResponse{
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
			},
		)
	}

	// Membuat PDF.
	// pdf := fpdf.New("L", "mm", "A4", "")

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

	pdf.SetTitle("Laporan Review Pohon Kinerja", false)
	pdf.SetAuthor("Pohon Kinerja", false)

	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)

	pdf.AddPage()

	// =========================
	// JUDUL
	// =========================

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "LAPORAN REVIEW POHON KINERJA", "", 1, "C", false, 0, "")

	pdf.Ln(3)

	// =========================
	// INFORMASI
	// =========================

	pdf.SetFont("Arial", "", 10)

	pdf.CellFormat(20, 6, "Kode OPD", "", 0, "L", false, 0, "")

	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")

	pdf.CellFormat(0, 6, kodeOpd, "", 1, "L", false, 0, "")

	pdf.CellFormat(20, 6, "Tahun", "", 0, "L", false, 0, "")

	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")

	pdf.CellFormat(0, 6, tahun, "", 1, "L", false, 0, "")

	pdf.Ln(5)

	// =========================
	// HEADER TABLE
	// =========================

	// Total lebar sekitar 277 mm.
	// A4 landscape = 297 mm
	// margin kiri + kanan = 20 mm
	//
	// No             = 12
	// Nama Pohon     = 48
	// Review         = 75
	// Keterangan     = 55
	// User Pembuat   = 40
	// Waktu Review   = 47
	//
	// Total          = 277

	// colWidths := []float64{12, 48, 75, 55, 40, 47}
	colWidths := []float64{12, 55, 85, 65, 45, 48}

	headers := []string{"No", "Nama Pohon", "Review", "Keterangan", "User Pembuat", "Waktu Review"}

	pdf.SetFont("Arial", "B", 9)

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}

	pdf.Ln(-1)

	// =========================
	// DATA
	// =========================

	pdf.SetFont("Arial", "", 8)

	if len(reviewResponses) == 0 {

		// pdf.CellFormat(277, 10, "Data Kosong / Belum Ditambahkan", "1", 1, "C", false, 0, "")
		pdf.CellFormat(310, 10, "Data Kosong / Belum Ditambahkan", "1", 1, "C", false, 0, "")

	} else {

		for index, data := range reviewResponses {

			namaPohon := data.NamaPohon
			if strings.TrimSpace(namaPohon) == "" {
				namaPohon = "-"
			}

			jenisPohon := data.JenisPohon
			if strings.TrimSpace(jenisPohon) == "" {
				jenisPohon = "-"
			}

			review := data.Review
			if strings.TrimSpace(review) == "" {
				review = "-"
			}

			keterangan := data.Keterangan
			if strings.TrimSpace(keterangan) == "" {
				keterangan = "-"
			}

			namaPegawai := data.NamaPegawai
			if strings.TrimSpace(namaPegawai) == "" {
				namaPegawai = "-"
			}

			waktuReview := formatWaktuReviewPDF(
				data.CreatedAt,
				data.UpdatedAt,
			)

			// Tinggi baris mengikuti isi Review/Keterangan.
			lineReview := pdf.SplitLines(
				[]byte(review),
				colWidths[2],
			)

			lineKeterangan := pdf.SplitLines(
				[]byte(keterangan),
				colWidths[3],
			)

			maxLines := len(lineReview)

			if len(lineKeterangan) > maxLines {
				maxLines = len(lineKeterangan)
			}

			// Minimal 2 baris karena Nama Pohon terdiri dari
			// nama + jenis pohon.
			if maxLines < 2 {
				maxLines = 2
			}

			rowHeight := float64(maxLines) * 4.5

			if rowHeight < 12 {
				rowHeight = 12
			}

			// =========================
			// NO
			// =========================

			pdf.CellFormat(colWidths[0], rowHeight, fmt.Sprintf("%d", index+1), "1", 0, "C", false, 0, "")

			// =========================
			// NAMA POHON
			// =========================

			x := pdf.GetX()
			y := pdf.GetY()

			pdf.MultiCell(colWidths[1], 4.5, namaPohon, "LR", "L", false)

			pdf.SetXY(x, y+4.5)

			pdf.SetFont("Arial", "", 7)

			pdf.MultiCell(colWidths[1], 4, strings.ToUpper(jenisPohon), "LR", "L", false)

			pdf.SetXY(
				x+colWidths[1],
				y,
			)

			pdf.SetFont("Arial", "", 8)

			// =========================
			// REVIEW
			// =========================

			pdf.MultiCell(colWidths[2], 4.5, review, "1", "L", false)

			pdf.SetXY(x+colWidths[1]+colWidths[2], y)

			// =========================
			// KETERANGAN
			// =========================

			pdf.MultiCell(colWidths[3], 4.5, keterangan, "1", "L", false)

			pdf.SetXY(x+colWidths[1]+colWidths[2]+colWidths[3], y)

			// =========================
			// USER PEMBUAT
			// =========================

			pdf.MultiCell(colWidths[4], 4.5, namaPegawai, "1", "L", false)

			pdf.SetXY(x+colWidths[1]+colWidths[2]+colWidths[3]+colWidths[4], y)

			// =========================
			// WAKTU REVIEW
			// =========================

			pdf.MultiCell(colWidths[5], 4.5, waktuReview, "1", "L", false)

			// Pastikan cursor pindah ke awal baris berikutnya.
			pdf.SetXY(x, y+rowHeight)

			pdf.SetFont("Arial", "", 8)
		}
	}

	// =========================
	// OUTPUT PDF
	// =========================

	var buffer bytes.Buffer

	err = pdf.Output(&buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func formatWaktuReviewPDF(createdAt, updatedAt string) string {

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}

	parseTime := func(value string) time.Time {
		for _, layout := range layouts {
			if result, err := time.Parse(layout, value); err == nil {
				return result
			}
		}

		return time.Time{}
	}

	created := parseTime(createdAt)
	updated := parseTime(updatedAt)

	label := "diedit pada :"
	target := updated

	if createdAt == updatedAt {
		label = "dibuat pada :"
		target = created
	}

	if target.IsZero() {
		return label + "\n" + "-"
	}

	months := []string{
		"Januari",
		"Februari",
		"Maret",
		"April",
		"Mei",
		"Juni",
		"Juli",
		"Agustus",
		"September",
		"Oktober",
		"November",
		"Desember",
	}

	waktu := fmt.Sprintf(
		"%d %s %d %02d:%02d",
		target.Day(),
		months[int(target.Month())-1],
		target.Year(),
		target.Hour(),
		target.Minute(),
	)

	return label + "\n" + waktu
}
