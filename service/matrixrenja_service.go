package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/model/web/renja"
)

type MatrixRenjaService interface {
	GetRenja(ctx context.Context, kodeOpd, tahun, jenisPagu string) ([]renja.RanwalMatrixResponse, error)
	GetRenjaRankhir(ctx context.Context, kodeOpd string, tahun string) ([]renja.RankhirMatrixResponse, error)
	GetRenjaPenetapan(ctx context.Context, kodeOpd, tahun, jenisPagu string) ([]renja.PenetapanMatrixResponse, error)
	UpsertBatchIndikatorRenja(ctx context.Context, requests []programkegiatan.IndikatorRenjaCreateRequest) ([]programkegiatan.IndikatorUpsertResponse, error)
	UpsertBatchIndikatorRenjaPenetapan(ctx context.Context, requests []programkegiatan.IndikatorRenjaCreateRequest) ([]programkegiatan.IndikatorUpsertResponse, error)
	UpsertAnggaran(ctx context.Context, request programkegiatan.AnggaranRenjaRequest) (programkegiatan.AnggaranRenjaResponse, error)
	// CRUD target per layer (rankhir / penetapan)
	UpsertTargetRenjaLayer(ctx context.Context, jenis string, req programkegiatan.TargetRenjaLayerRequest) (renja.TargetLayerResponse, error)
	DeleteTargetRenjaLayer(ctx context.Context, kodeIndikator, tahun, jenis string) error
	FindTargetRenjaLayer(ctx context.Context, kodeIndikator, tahun, jenis string) (renja.TargetLayerResponse, error)
	// CRUD pagu per jenis (ranwal / rankhir)
	UpsertPaguRenjaByJenis(ctx context.Context, jenis string, req programkegiatan.PaguRenjaLayerRequest) error

	// V2 — sama dengan Get/Rankhir/Penetapan tetapi khusus level Program:
	// indikator dari tb_master_indikator_program_kegiatan, target dari tb_target_program_kegiatan.
	GetRenjaV2(ctx context.Context, kodeOpd, tahun, jenisPagu string) ([]renja.RanwalMatrixResponse, error)
	GetRenjaRankhirV2(ctx context.Context, kodeOpd, tahun string) ([]renja.RankhirMatrixResponse, error)
	GetRenjaPenetapanV2(ctx context.Context, kodeOpd, tahun, jenisPagu string) ([]renja.PenetapanMatrixResponse, error)
}
