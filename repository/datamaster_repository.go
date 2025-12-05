package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/datamaster"
)

type DataMasterRepository interface {
	DataRBByTahun(ctx context.Context, tx *sql.Tx, tahun int) ([]datamaster.MasterRB, error)
	InsertRB(ctx context.Context, tx *sql.Tx, req datamaster.MasterRB, userId int) (int64, error)
	UpdateRB(ctx context.Context, tx *sql.Tx, rb datamaster.MasterRB, rbId int) error
	InsertIndikator(ctx context.Context, tx *sql.Tx, rbId int64, indikator datamaster.IndikatorRB) (string, error)
	InsertTarget(ctx context.Context, tx *sql.Tx, indikatorID string, t datamaster.TargetRB) error
	FindRBById(ctx context.Context, tx *sql.Tx, rbId int) (datamaster.MasterRB, error)
	DeleteAllIndikatorAndTargetByRB(ctx context.Context, tx *sql.Tx, rbId int) error
	DeleteRB(ctx context.Context, tx *sql.Tx, rbId int) error
	PokinByIdRBs(ctx context.Context, tx *sql.Tx, listIdRB []int) ([]datamaster.PokinIdRBTagging, error)
}
