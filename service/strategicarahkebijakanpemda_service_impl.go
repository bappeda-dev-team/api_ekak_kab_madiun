package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/strategicarahkebijakan"
	"ekak_kabupaten_madiun/repository"
	"log"
	"time"
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

func (service *StrategicArahKebijakanPemdaServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) (strategicarahkebijakan.StrategicArahKebijakanPemdaAllResponse, error) {
	startTime := time.Now()
	serviceName := "StrategicArahKebijakanPemdaService.FindAll"

	tx, err := service.DB.Begin()
	if err != nil {
		return strategicarahkebijakan.StrategicArahKebijakanPemdaAllResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Inisialisasi response dasar
	response := strategicarahkebijakan.StrategicArahKebijakanPemdaAllResponse{
		IsuStrategisPemda: make([]strategicarahkebijakan.IsuStrategiPemdaResponse, 0),
		TujuanPemda:  make([]strategicarahkebijakan.TujuanPemdaResponse, 0),
		StrategiArahKebijakanPemdas: make([]strategicarahkebijakan.StrategiArahKebijakanPemdaResponse, 0),
	}

	csfList, err := service.csfRepository.IsuFindBetweenTahun(ctx, tx, tahunAwal, tahunAkhir)
	if err != nil {
		return strategicarahkebijakan.StrategicArahKebijakanPemdaAllResponse{}, err
	}
	if len(csfList) > 0 {
		Responses := make([]strategicarahkebijakan.IsuStrategiPemdaResponse, 0, len(csfList))
		for _, tujuan := range csfList {
			Responses = append(Responses, strategicarahkebijakan.IsuStrategiPemdaResponse{
				NamaIsu:        tujuan.NamaIsu,
			})
		}
		response.IsuStrategisPemda = Responses
	}
	

	// Ambil data tujuan OPD dengan batch
	tujuanPemdas, err := service.tujuanPemdaRepository.FindAllBetweenTahun(ctx, tx, tahunAwal, tahunAkhir, "RPJMD")
	if err != nil {
		return strategicarahkebijakan.StrategicArahKebijakanPemdaAllResponse{}, err
	}
	if len(tujuanPemdas) > 0 {
		tujuanResponses := make([]strategicarahkebijakan.TujuanPemdaResponse, 0, len(tujuanPemdas))
		for _, tujuan := range tujuanPemdas {
			tujuanResponses = append(tujuanResponses, strategicarahkebijakan.TujuanPemdaResponse{
				Id:        tujuan.Id,
				Tujuan:    tujuan.TujuanPemda,
			})
		}
		response.TujuanPemda = tujuanResponses
	}

	sasaranOpds, err := service.sasaranPemdaRepository.FindStrategicArahKebijakanPemda(ctx, tx, tahunAwal, tahunAkhir, "RPJMD")
	if err != nil {
		return strategicarahkebijakan.StrategicArahKebijakanPemdaAllResponse{}, err
	}
	if len(sasaranOpds) > 0 {
	strategiResponses := make([]strategicarahkebijakan.StrategiArahKebijakanPemdaResponse, 0)

		for _, s := range sasaranOpds {

			strategiResponses = append(strategiResponses, strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{
				TujuanPemda: s.NamaTujuanPemda, // pastikan field ini ada di domain
				SasaranPemdas: []strategicarahkebijakan.SasaranPemdaResponse{
					{
						SasaranPemda: s.NamaSasaranPemda,
						StrategiPemda: s.NamaStrategi, // kalau ada
						ArahKebijakanPemdas: []strategicarahkebijakan.ArahKebijakanPemdaResponse{
							{
								ArahKebijakanPemda: s.NamaArahKebijakan, // kalau ada
							},
						},
					},
				},
			})
		}

		response.StrategiArahKebijakanPemdas = strategiResponses
	}


	log.Printf("[%s] [END] [%s] totalResponseTime=%v, strategicsCount=%d",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(startTime), len(response.TujuanPemda))

	return response, nil
}

