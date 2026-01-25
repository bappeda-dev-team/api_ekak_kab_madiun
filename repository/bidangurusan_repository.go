package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type BidangUrusanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, bidangurusan domainmaster.BidangUrusan) domainmaster.BidangUrusan
	Update(ctx context.Context, tx *sql.Tx, bidangurusan domainmaster.BidangUrusan) domainmaster.BidangUrusan
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.BidangUrusan, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.BidangUrusan, error)
	FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.BidangUrusan, error)
	FindByKodeBidangUrusan(ctx context.Context, tx *sql.Tx, kodeBidangUrusan string) (domainmaster.BidangUrusan, error)

	//bidangurusanOPD
	CreateOPD(ctx context.Context, tx *sql.Tx, bidangurusanOpd domainmaster.BidangUrusanOpd) (domainmaster.BidangUrusanOpd, error)
	DeleteOPD(ctx context.Context, tx *sql.Tx, id string) error
	FindBidangUrusanTerpilihByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.BidangUrusanOpd, error)
	IsBidangUrusanMasterExists(ctx context.Context, tx *sql.Tx, kodeBidangUrusan string) (bool, error)
	IsBidangUrusanAlreadySelected(ctx context.Context, tx *sql.Tx, kodeOpd string, kodeBidangUrusan string) (bool, error)
}
