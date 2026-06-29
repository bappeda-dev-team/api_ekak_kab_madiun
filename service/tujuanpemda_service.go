package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/tujuanpemda"
)

type TujuanPemdaService interface {
	Create(ctx context.Context, request tujuanpemda.TujuanPemdaCreateRequest) (tujuanpemda.TujuanPemdaResponse, error)
	Update(ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest) (tujuanpemda.TujuanPemdaResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, tujuanPemdaId int) (tujuanpemda.TujuanPemdaResponse, error)
	FindAll(ctx context.Context, tahun string, jenisPeriode string) ([]tujuanpemda.TujuanPemdaResponse, error)
	UpdatePeriode(ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest) (tujuanpemda.TujuanPemdaResponse, error)
	FindAllWithPokin(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]tujuanpemda.TujuanPemdaWithPokinResponse, error)
	FindAllWithPokinRenstra(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]tujuanpemda.TujuanPemdaWithPokinResponse, error)
	FindPokinWithPeriode(ctx context.Context, pokinId int, jenisPeriode string) (tujuanpemda.PokinWithPeriodeResponse, error)

	FindTujuanPemdaRanwal(ctx context.Context, tahun, jenisPeriode string) ([]tujuanpemda.TujuanPemdaResponse, error)
	FindTujuanPemdaRankhir(ctx context.Context, tahun, jenisPeriode string) ([]tujuanpemda.TujuanPemdaResponse, error)
	FindTujuanPemdaPenetapan(ctx context.Context, tahun, jenisPeriode string) ([]tujuanpemda.TujuanPemdaResponse, error)
	UpsertTargetPemdaLayer(ctx context.Context, jenis string, request tujuanpemda.LayerTargetBatchRequest) ([]tujuanpemda.TargetResponse, error)
	CreateTargetPemdaLayer(ctx context.Context, jenis string, request tujuanpemda.LayerTargetBatchRequest) ([]tujuanpemda.TargetResponse, error)
	UpdateTargetPemdaLayer(ctx context.Context, jenis string, request tujuanpemda.LayerTargetUpdateBatchRequest) ([]tujuanpemda.TargetResponse, error)

	// Opsi B — tampilkan 2 jenis target sekaligus (tanpa fallback)
	FindTujuanPemdaRankhirDual(ctx context.Context, tahun, jenisPeriode string) ([]tujuanpemda.TujuanPemdaResponse, error)
	FindTujuanPemdaPenetapanDual(ctx context.Context, tahun, jenisPeriode string) ([]tujuanpemda.TujuanPemdaResponse, error)

	//lock pemda
	LockTujuanPemda(ctx context.Context, tahun string) error
	UnlockTujuanPemda(ctx context.Context, tahun string) error
	IsTujuanPemdaLocked(ctx context.Context, tahun string) (bool, error)
	FindAllLockTujuanPemda(ctx context.Context) ([]tujuanpemda.LockDataPemdaResponse, error)
}
