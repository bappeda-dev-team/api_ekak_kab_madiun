package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/taggingpokin"
	"ekak_kabupaten_madiun/repository"
)

type TaggingPokinServiceImpl struct {
	TaggingPokinRepository repository.TaggingPokinRepository
	DB                     *sql.DB
}

func NewTaggingPokinServiceImpl(taggingPokinRepository repository.TaggingPokinRepository, db *sql.DB) *TaggingPokinServiceImpl {
	return &TaggingPokinServiceImpl{
		TaggingPokinRepository: taggingPokinRepository,
		DB:                     db,
	}
}

func (service *TaggingPokinServiceImpl) Create(ctx context.Context, request taggingpokin.TaggingPokinCreateRequest) (taggingpokin.TaggingPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return taggingpokin.TaggingPokinResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tagging := domain.TaggingPokin{
		NamaTagging:       request.NamaTagging,
		KeteranganTagging: request.KeteranganTagging,
	}

	taggingResponse, err := service.TaggingPokinRepository.Create(ctx, tx, tagging)
	if err != nil {
		return taggingpokin.TaggingPokinResponse{}, err
	}
	return taggingpokin.TaggingPokinResponse{
		Id:                taggingResponse.Id,
		NamaTagging:       taggingResponse.NamaTagging,
		KeteranganTagging: taggingResponse.KeteranganTagging,
	}, nil
}

func (service *TaggingPokinServiceImpl) Update(ctx context.Context, request taggingpokin.TaggingPokinUpdateRequest) (taggingpokin.TaggingPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return taggingpokin.TaggingPokinResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tagging := domain.TaggingPokin{
		Id:                request.Id,
		NamaTagging:       request.NamaTagging,
		KeteranganTagging: request.KeteranganTagging,
	}

	taggingResponse, err := service.TaggingPokinRepository.Update(ctx, tx, tagging)
	if err != nil {
		return taggingpokin.TaggingPokinResponse{}, err
	}
	return taggingpokin.TaggingPokinResponse{
		Id:                taggingResponse.Id,
		NamaTagging:       taggingResponse.NamaTagging,
		KeteranganTagging: taggingResponse.KeteranganTagging,
	}, nil
}

func (service *TaggingPokinServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.TaggingPokinRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (service *TaggingPokinServiceImpl) FindById(ctx context.Context, id int) (taggingpokin.TaggingPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return taggingpokin.TaggingPokinResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tagging, err := service.TaggingPokinRepository.FindById(ctx, tx, id)
	if err != nil {
		return taggingpokin.TaggingPokinResponse{}, err
	}
	return taggingpokin.TaggingPokinResponse{
		Id:                tagging.Id,
		NamaTagging:       tagging.NamaTagging,
		KeteranganTagging: tagging.KeteranganTagging,
	}, nil
}

func (service *TaggingPokinServiceImpl) FindAll(ctx context.Context) ([]taggingpokin.TaggingPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []taggingpokin.TaggingPokinResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	taggings, err := service.TaggingPokinRepository.FindAll(ctx, tx)
	if err != nil {
		return []taggingpokin.TaggingPokinResponse{}, err
	}

	var taggingResponses []taggingpokin.TaggingPokinResponse
	for _, tagging := range taggings {
		taggingResponses = append(taggingResponses, taggingpokin.TaggingPokinResponse{
			Id:                tagging.Id,
			NamaTagging:       tagging.NamaTagging,
			KeteranganTagging: tagging.KeteranganTagging,
		})
	}
	return taggingResponses, nil
}
