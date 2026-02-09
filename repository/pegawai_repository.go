package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type PegawaiRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pegawai domainmaster.Pegawai) (domainmaster.Pegawai, error)
	Update(ctx context.Context, tx *sql.Tx, pegawai domainmaster.Pegawai) domainmaster.Pegawai
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Pegawai, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Pegawai, error)
	FindByNip(ctx context.Context, tx *sql.Tx, nip string) (domainmaster.Pegawai, error)
	FindByNipWithJabatan(ctx context.Context, tx *sql.Tx, nip string) (domainmaster.Pegawai, error)
	FindPegawaiByNipsBatch(ctx context.Context, tx *sql.Tx, nips []string) (map[string]*domainmaster.Pegawai, error)
}
