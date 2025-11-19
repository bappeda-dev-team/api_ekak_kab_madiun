package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
	"log"
	"strings"
)

type KegiatanRepositoryImpl struct {
}

func NewKegiatanRepositoryImpl() *KegiatanRepositoryImpl {
	return &KegiatanRepositoryImpl{}
}

func (repository *KegiatanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, kegiatan domainmaster.Kegiatan) (domainmaster.Kegiatan, error) {
	scriptKegiatan := "INSERT INTO tb_master_kegiatan (id, nama_kegiatan, kode_kegiatan) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, scriptKegiatan, kegiatan.Id, kegiatan.NamaKegiatan, kegiatan.KodeKegiatan)
	if err != nil {
		return domainmaster.Kegiatan{}, err
	}

	for _, indikator := range kegiatan.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, kegiatan_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator, indikator.Id, kegiatan.Id, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domainmaster.Kegiatan{}, err
		}

		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, tahun, target, satuan) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Tahun, target.Target, target.Satuan)
			if err != nil {
				return domainmaster.Kegiatan{}, err
			}
		}
	}

	return kegiatan, nil
}

func (repository *KegiatanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, kegiatan domainmaster.Kegiatan) (domainmaster.Kegiatan, error) {
	scriptKegiatan := "UPDATE tb_master_kegiatan SET nama_kegiatan = ?, kode_kegiatan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, scriptKegiatan, kegiatan.NamaKegiatan, kegiatan.KodeKegiatan, kegiatan.Id)
	if err != nil {
		return domainmaster.Kegiatan{}, err
	}

	for _, indikator := range kegiatan.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, kegiatan_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator, indikator.Id, kegiatan.Id, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domainmaster.Kegiatan{}, err
		}

		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, tahun, target, satuan) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Tahun, target.Target, target.Satuan)
			if err != nil {
				return domainmaster.Kegiatan{}, err
			}
		}
	}

	return kegiatan, nil
}

func (repository *KegiatanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	// Delete target terlebih dahulu
	scriptTarget := "DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE kegiatan_id = ?)"
	_, err := tx.ExecContext(ctx, scriptTarget, id)
	if err != nil {
		return err
	}

	// Delete indikator
	scriptIndikator := "DELETE FROM tb_indikator WHERE kegiatan_id = ?"
	_, err = tx.ExecContext(ctx, scriptIndikator, id)
	if err != nil {
		return err
	}

	// Delete kegiatan
	scriptKegiatan := "DELETE FROM tb_master_kegiatan WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptKegiatan, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *KegiatanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Kegiatan, error) {
	scriptKegiatan := "SELECT id, nama_kegiatan, kode_kegiatan FROM tb_master_kegiatan WHERE id = ?"
	row := tx.QueryRowContext(ctx, scriptKegiatan, id)
	var kegiatan domainmaster.Kegiatan
	err := row.Scan(&kegiatan.Id, &kegiatan.NamaKegiatan, &kegiatan.KodeKegiatan)
	if err != nil {
		return domainmaster.Kegiatan{}, err
	}
	return kegiatan, nil
}

func (repository *KegiatanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Kegiatan, error) {
	scriptKegiatan := "SELECT id, nama_kegiatan, kode_kegiatan FROM tb_master_kegiatan"
	rows, err := tx.QueryContext(ctx, scriptKegiatan)
	if err != nil {
		return []domainmaster.Kegiatan{}, err
	}
	defer rows.Close()
	var kegiatans []domainmaster.Kegiatan
	for rows.Next() {
		var kegiatan domainmaster.Kegiatan
		rows.Scan(&kegiatan.Id, &kegiatan.NamaKegiatan, &kegiatan.KodeKegiatan)
		kegiatans = append(kegiatans, kegiatan)
	}
	return kegiatans, nil
}

func (repository *KegiatanRepositoryImpl) FindIndikatorByKegiatanId(ctx context.Context, tx *sql.Tx, kegiatanId string) ([]domain.Indikator, error) {
	scriptIndikator := "SELECT id, kegiatan_id, indikator, tahun FROM tb_indikator WHERE kegiatan_id = ?"
	rows, err := tx.QueryContext(ctx, scriptIndikator, kegiatanId)
	if err != nil {
		return []domain.Indikator{}, err
	}
	defer rows.Close()
	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		rows.Scan(&indikator.Id, &indikator.KegiatanId, &indikator.Indikator, &indikator.Tahun)
		indikators = append(indikators, indikator)
	}
	return indikators, nil
}

func (repository *KegiatanRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	scriptTarget := "SELECT id, indikator_id, target, tahun FROM tb_target WHERE indikator_id = ?"
	rows, err := tx.QueryContext(ctx, scriptTarget, indikatorId)
	if err != nil {
		return []domain.Target{}, err
	}
	defer rows.Close()
	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Tahun)
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *KegiatanRepositoryImpl) FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domainmaster.Kegiatan, error) {
	if kodeSubKegiatan == "" {
		return domainmaster.Kegiatan{}, fmt.Errorf("Kode Subkegiatan kosong")
	}
	if len(kodeSubKegiatan) < 17 {
		return domainmaster.Kegiatan{}, fmt.Errorf("kode_subkegiatan '%s' terlalu pendek untuk ekstraksi kode kegiatan", kodeSubKegiatan)
	}
	kodeKegiatan := kodeSubKegiatan[:12]
	log.Printf("KodeKegiatan: %s", kodeKegiatan)
	script := "SELECT kg.nama_kegiatan, kg.kode_kegiatan FROM tb_master_kegiatan kg WHERE kg.kode_kegiatan = ?"
	rows, err := tx.QueryContext(ctx, script, kodeKegiatan)
	if err != nil {
		return domainmaster.Kegiatan{}, err
	}
	defer rows.Close()

	kegiatan := domainmaster.Kegiatan{}

	if rows.Next() {
		if err := rows.Scan(&kegiatan.NamaKegiatan, &kegiatan.KodeKegiatan); err != nil {
			return domainmaster.Kegiatan{}, fmt.Errorf("gagal scan kegiatan: %w", err)
		}
		return kegiatan, nil
	}

	return domainmaster.Kegiatan{}, fmt.Errorf("kegiatan dengan kode %s tidak ditemukan", kodeKegiatan)
}

func (repository *KegiatanRepositoryImpl) FindByKodeSubs(
    ctx context.Context,
    tx *sql.Tx,
    kodeSubs []string,
) (map[string]domainmaster.Kegiatan, error) {

    resultMap := make(map[string]domainmaster.Kegiatan)

    if len(kodeSubs) == 0 {
        return resultMap, nil
    }

    // Step 1. Extract unique kode_kegiatan from kode_sub
    uniqKegiatan := make(map[string]struct{})
    kodeKegiatanList := make([]string, 0)

    // Map kode_kegiatan → list kode_sub_kegiatan yg punya prefix itu
    kodeSubGroup := make(map[string][]string)

    for _, ks := range kodeSubs {
        if len(ks) < 12 {
            // skip invalid kode_sub
            continue
        }

        kodeKegiatan := ks[:12]

        // group kodeSub by kodeKegiatan
        kodeSubGroup[kodeKegiatan] = append(kodeSubGroup[kodeKegiatan], ks)

        if _, exists := uniqKegiatan[kodeKegiatan]; !exists {
            uniqKegiatan[kodeKegiatan] = struct{}{}
            kodeKegiatanList = append(kodeKegiatanList, kodeKegiatan)
        }
    }

    if len(kodeKegiatanList) == 0 {
        return resultMap, nil
    }

    // Step 2. Build placeholders
    placeholders := make([]string, len(kodeKegiatanList))
    args := make([]any, len(kodeKegiatanList))

    for i, kode := range kodeKegiatanList {
        placeholders[i] = "?"
        args[i] = kode
    }

    query := fmt.Sprintf(`
        SELECT
            kg.nama_kegiatan,
            kg.kode_kegiatan
        FROM tb_master_kegiatan kg
        WHERE kg.kode_kegiatan IN (%s)
    `, strings.Join(placeholders, ","))

    rows, err := tx.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Step 3. Scan DB → generate map
    kegiatanByKode := make(map[string]domainmaster.Kegiatan)

    for rows.Next() {
        var k domainmaster.Kegiatan
        if err := rows.Scan(&k.NamaKegiatan, &k.KodeKegiatan); err != nil {
            return nil, err
        }
        kegiatanByKode[k.KodeKegiatan] = k
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Step 4. Map kembali ke kode_subkegiatan asli
    for kodeKegiatan, subs := range kodeSubGroup {
        baseKeg, ok := kegiatanByKode[kodeKegiatan]
        if !ok {
            // Tidak ditemukan di DB → skip
            continue
        }

        for _, kodeSub := range subs {
            resultMap[kodeSub] = domainmaster.Kegiatan{
                NamaKegiatan:    baseKeg.NamaKegiatan,
                KodeKegiatan:    baseKeg.KodeKegiatan,
                KodeSubKegiatan: kodeSub, // ★ aslinya disimpan di struct
            }
        }
    }

    return resultMap, nil
}
