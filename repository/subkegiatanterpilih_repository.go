package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanTerpilihRepository interface {
	Update(ctx context.Context, tx *sql.Tx, subKegiatanTerpilih domain.SubKegiatanTerpilih) (domain.SubKegiatanTerpilih, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) error
	FindByIdAndKodeSubKegiatan(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) (domain.SubKegiatanTerpilih, error)
	CreateRekin(ctx context.Context, tx *sql.Tx, idSubKegiatan string, rekinId string, kodeSubKegiatan string) error
	DeleteSubKegiatanTerpilih(ctx context.Context, tx *sql.Tx, idSubKegiatan string) error
	FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.SubKegiatanTerpilih, error)
	//subkegiatan opd
	CreateOPD(ctx context.Context, tx *sql.Tx, subkegiatanOpd domain.SubKegiatanOpd) (domain.SubKegiatanOpd, error)
	UpdateOPD(ctx context.Context, tx *sql.Tx, subkegiatanOpd domain.SubKegiatanOpd) (domain.SubKegiatanOpd, error)
	FindallOpd(ctx context.Context, tx *sql.Tx, kodeOpd, tahun *string) ([]domain.SubKegiatanOpd, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.SubKegiatanOpd, error)
	DeleteSubOpd(ctx context.Context, tx *sql.Tx, id int) error
	FindAllSubkegiatanByBidangUrusanOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.SubKegiatan, error)
	FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domain.SubKegiatan, error)
	CheckExists(ctx context.Context, tx *sql.Tx, kodeSubkegiatan, kodeOpd, tahun string) (bool, error)
}
