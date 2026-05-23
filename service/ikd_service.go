package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/web/ikd"
)

type IkdService interface {
	FindAll(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string) ([]ikd.IkdResponse, error)
	Create(ctx context.Context, request []ikd.ProgramOpdTerpilihCreateRequest) ([]ikd.ProgramOpdTerpilihResponse, error)
	Delete(ctx context.Context, id int) error
	LockProgramOpdTerpilih(ctx context.Context, id int) error
	UnlockProgramOpdTerpilih(ctx context.Context, id int) error
	ensureNotLocked(ctx context.Context, tx *sql.Tx, id int) error
}