package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/strategicarahkebijakan"
	"ekak_kabupaten_madiun/repository"
)

type StrategicArahKebijakanPemdaServiceImpl struct {
	csfRepository             repository.CSFRepository
	DB                        *sql.DB
	tujuanPemdaRepository     repository.TujuanPemdaRepository
	sasaranPemdaRepository    repository.SasaranPemdaRepository
}

func NewStrategicArahKebijakanPemdaServiceImpl(csfRepository repository.CSFRepository, DB *sql.DB, tujuanPemdaRepository repository.TujuanPemdaRepository, sasaranPemdaRepository repository.SasaranPemdaRepository) *StrategicArahKebijakanPemdaServiceImpl {
	return &StrategicArahKebijakanPemdaServiceImpl{
		DB:                        DB,
		csfRepository:             csfRepository,
		tujuanPemdaRepository: tujuanPemdaRepository,
		sasaranPemdaRepository: sasaranPemdaRepository,
	}
}

func (service *StrategicArahKebijakanPemdaServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) ([]strategicarahkebijakan.StrategiArahKebijakanPemdaResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return []strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Inisialisasi response dasar
	response := []strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{}

	sasaranOpds, err := service.sasaranPemdaRepository.FindStrategicArahKebijakanPemda(ctx, tx, tahunAwal, tahunAkhir, "RPJMD")
	if err != nil {
		return []strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{}, err
	}
	if len(sasaranOpds) > 0 {
		strategiResponses := make([]strategicarahkebijakan.StrategiArahKebijakanPemdaResponse, 0)

		for _, s := range sasaranOpds {

			// arah kebijakan (bisa null)
			var arahKebijakan []strategicarahkebijakan.ArahKebijakanPemdaResponse
			if s.NamaArahKebijakan != "" {
				arahKebijakan = []strategicarahkebijakan.ArahKebijakanPemdaResponse{
					{
						ArahKebijakanPemda: s.NamaArahKebijakan,
					},
				}
			}

			// sasaran (bisa null)
			var sasaran []strategicarahkebijakan.SasaranPemdaResponse
			if s.NamaSasaranPemda != "" {
				sasaran = []strategicarahkebijakan.SasaranPemdaResponse{
					{
						SasaranPemda:        s.NamaSasaranPemda,
						StrategiPemda:       s.NamaStrategi,
						ArahKebijakanPemdas: arahKebijakan,
					},
				}
			}

			strategiResponses = append(strategiResponses, strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{
				TujuanPemda:   s.NamaTujuanPemda,
				SasaranPemdas: sasaran,
			})
		}

		response = strategiResponses
	}

	return response, nil
}

