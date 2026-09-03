package service

import (
	"bytes"
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/strategicarahkebijakan"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type StrategicArahKebijakanPemdaServiceImpl struct {
	opdRepository             repository.OpdRepository
	csfRepository             repository.CSFRepository
	DB                        *sql.DB
	tujuanPemdaRepository     repository.TujuanPemdaRepository
	sasaranPemdaRepository    repository.SasaranPemdaRepository
}

func NewStrategicArahKebijakanPemdaServiceImpl(opdRepository repository.OpdRepository, csfRepository repository.CSFRepository, DB *sql.DB, tujuanPemdaRepository repository.TujuanPemdaRepository, sasaranPemdaRepository repository.SasaranPemdaRepository) *StrategicArahKebijakanPemdaServiceImpl {
	return &StrategicArahKebijakanPemdaServiceImpl{
		opdRepository:             opdRepository,
		DB:                        DB,
		csfRepository:             csfRepository,
		tujuanPemdaRepository: tujuanPemdaRepository,
		sasaranPemdaRepository: sasaranPemdaRepository,
	}
}

func (service *StrategicArahKebijakanPemdaServiceImpl) buildStrategicResponse(sasaranOpds []domain.StrategicPemdaRow) []strategicarahkebijakan.StrategiArahKebijakanPemdaResponse {

	response := make([]strategicarahkebijakan.StrategiArahKebijakanPemdaResponse, 0)

	// map tujuan -> index response
	tujuanMap := make(map[string]int)

	for _, s := range sasaranOpds {

		// =========================
		// TUJUAN
		// =========================
		tujuanIdx, exists := tujuanMap[s.NamaTujuanPemda]

		if !exists {
			response = append(response,
				strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{
					TujuanPemda:   s.NamaTujuanPemda,
					SasaranPemdas: []strategicarahkebijakan.SasaranPemdaResponse{},
				},
			)

			tujuanIdx = len(response) - 1
			tujuanMap[s.NamaTujuanPemda] = tujuanIdx
		}

		// =========================
		// SASARAN
		// =========================
		sasaranIdx := -1

		for i, sasaran := range response[tujuanIdx].SasaranPemdas {
			if sasaran.SasaranPemda == s.NamaSasaranPemda {
				sasaranIdx = i
				break
			}
		}

		if sasaranIdx == -1 {
			response[tujuanIdx].SasaranPemdas = append(
				response[tujuanIdx].SasaranPemdas,
				strategicarahkebijakan.SasaranPemdaResponse{
					SasaranPemda:   s.NamaSasaranPemda,
					StrategiPemdas: []strategicarahkebijakan.StrategiPemdaResponse{},
				},
			)

			sasaranIdx = len(response[tujuanIdx].SasaranPemdas) - 1
		}

		// =========================
		// STRATEGI
		// =========================
		strategiIdx := -1

		for i, strategi := range response[tujuanIdx].SasaranPemdas[sasaranIdx].StrategiPemdas {
			if strategi.StrategiPemda == s.NamaStrategi {
				strategiIdx = i
				break
			}
		}

		if strategiIdx == -1 {
			response[tujuanIdx].SasaranPemdas[sasaranIdx].StrategiPemdas = append(
				response[tujuanIdx].SasaranPemdas[sasaranIdx].StrategiPemdas,
				strategicarahkebijakan.StrategiPemdaResponse{
					StrategiPemda:       s.NamaStrategi,
					ArahKebijakanPemdas: []strategicarahkebijakan.ArahKebijakanPemdaResponse{},
				},
			)

			strategiIdx = len(response[tujuanIdx].SasaranPemdas[sasaranIdx].StrategiPemdas) - 1
		}

		// =========================
		// ARAH KEBIJAKAN
		// =========================
		if s.NamaArahKebijakan != "" {

			sudahAda := false

			for _, ak := range response[tujuanIdx].
				SasaranPemdas[sasaranIdx].
				StrategiPemdas[strategiIdx].
				ArahKebijakanPemdas {

				if ak.ArahKebijakanPemda == s.NamaArahKebijakan {
					sudahAda = true
					break
				}
			}

			if !sudahAda {
				response[tujuanIdx].
					SasaranPemdas[sasaranIdx].
					StrategiPemdas[strategiIdx].
					ArahKebijakanPemdas = append(
					response[tujuanIdx].
						SasaranPemdas[sasaranIdx].
						StrategiPemdas[strategiIdx].
						ArahKebijakanPemdas,
					strategicarahkebijakan.ArahKebijakanPemdaResponse{
						ArahKebijakanPemda: s.NamaArahKebijakan,
					},
				)
			}
		}
	}

	return response
}

func (service *StrategicArahKebijakanPemdaServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) ([]strategicarahkebijakan.StrategiArahKebijakanPemdaResponse, error) {
	
	tx, err := service.DB.Begin()
	if err != nil {
		return []strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Inisialisasi response dasar
	response := []strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{}
	
	sasaranPemdas, err := service.sasaranPemdaRepository.FindStrategicArahKebijakanPemda(ctx, tx, tahunAwal, tahunAkhir, "RPJMD")
	if err != nil {
		return []strategicarahkebijakan.StrategiArahKebijakanPemdaResponse{}, err
	}

	response = service.buildStrategicResponse(sasaranPemdas)

	return response, nil
}

func (service *StrategicArahKebijakanPemdaServiceImpl) ExportExcel(ctx context.Context, tahunAwal string, tahunAkhir string) (*bytes.Buffer, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranPemdas, err := service.sasaranPemdaRepository.FindStrategicArahKebijakanPemda(ctx, tx, tahunAwal, tahunAkhir, "RPJMD")
	if err != nil {
		return nil, err
	}

	responses := service.buildStrategicResponse(sasaranPemdas)

	f := excelize.NewFile()

	sheet := "Strategic Arah Kebijakan"
	f.SetSheetName("Sheet1", sheet)

	// ==========================
	// Style Header
	// ==========================
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
			Size:  12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#10B981"},
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})

	// ==========================
	// Style Body
	// ==========================
	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			WrapText:   true,
			Horizontal: "left",
		},
	})

	// Style nomor
	noStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// ==========================
	// Judul
	// ==========================
	f.MergeCell(sheet, "A1", "E1")
	f.SetCellValue(sheet, "A1", "STRATEGI DAN ARAH KEBIJAKAN PEMDA")

	
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 16,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	f.SetCellStyle(sheet, "A1", "E1", titleStyle)
	f.SetRowHeight(sheet, 1, 30)

	// ==========================
	// Header
	// ==========================
	f.SetCellValue(sheet, "A3", "No")
	f.SetCellValue(sheet, "B3", "Tujuan Pemda")
	f.SetCellValue(sheet, "C3", "Sasaran Pemda")
	f.SetCellValue(sheet, "D3", "Strategi")
	f.SetCellValue(sheet, "E3", "Arah Kebijakan")

	f.SetCellStyle(sheet, "A3", "E3", headerStyle)
	f.SetRowHeight(sheet, 3, 25)

	// ==========================
	// Header
	// ==========================
	f.SetCellValue(sheet, "A4", "1")
	f.SetCellValue(sheet, "B4", "2")
	f.SetCellValue(sheet, "C4", "3")
	f.SetCellValue(sheet, "D4", "4")
	f.SetCellValue(sheet, "E4", "5")

	f.SetCellStyle(sheet, "A4", "E4", headerStyle)
	f.SetRowHeight(sheet, 4, 25)

	// ==========================
	// Isi Data
	// ==========================
	row := 5
	no := 1

	for _, tujuan := range responses {

		awalTujuan := row

		for _, sasaran := range tujuan.SasaranPemdas {

			awalSasaran := row

			for _, strategi := range sasaran.StrategiPemdas {

				awalStrategi := row

				if len(strategi.ArahKebijakanPemdas) == 0 {

					f.SetCellValue(sheet, fmt.Sprintf("D%d", row), strategi.StrategiPemda)

					f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), noStyle)
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("E%d", row), bodyStyle)

					row++

				} else {

					for _, arah := range strategi.ArahKebijakanPemdas {

						f.SetCellValue(sheet, fmt.Sprintf("E%d", row), arah.ArahKebijakanPemda)

						f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), noStyle)
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("E%d", row), bodyStyle)

						row++
					}
				}

				f.SetCellValue(sheet, fmt.Sprintf("D%d", awalStrategi), strategi.StrategiPemda)

				if row-awalStrategi > 1 {
					f.MergeCell(sheet,
						fmt.Sprintf("D%d", awalStrategi),
						fmt.Sprintf("D%d", row-1))
				}
			}

			f.SetCellValue(sheet, fmt.Sprintf("C%d", awalSasaran), sasaran.SasaranPemda)

			if row-awalSasaran > 1 {
				f.MergeCell(sheet,
					fmt.Sprintf("C%d", awalSasaran),
					fmt.Sprintf("C%d", row-1))
			}
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", awalTujuan), no)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", awalTujuan), tujuan.TujuanPemda)

		if row-awalTujuan > 1 {

			f.MergeCell(sheet,
				fmt.Sprintf("A%d", awalTujuan),
				fmt.Sprintf("A%d", row-1))

			f.MergeCell(sheet,
				fmt.Sprintf("B%d", awalTujuan),
				fmt.Sprintf("B%d", row-1))
		}

		no++
	}

	// ==========================
	// Lebar Kolom
	// ==========================
	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 45)
	f.SetColWidth(sheet, "C", "C", 40)
	f.SetColWidth(sheet, "D", "D", 40)
	f.SetColWidth(sheet, "E", "E", 60)

	// Freeze Header
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      3,
		TopLeftCell: "A4",
		ActivePane:  "bottomLeft",
	})

	buffer := new(bytes.Buffer)
	if err := f.Write(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

