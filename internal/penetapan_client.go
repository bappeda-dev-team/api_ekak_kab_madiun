package internal

import "context"

type PenetapanClient interface {
	SyncPenetapanPkPegawai(ctx context.Context, pegawaiId string, kodeOpd string, tahun int) error
}
