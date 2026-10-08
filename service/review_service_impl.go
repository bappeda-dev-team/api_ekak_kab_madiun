package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"math/rand"
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
