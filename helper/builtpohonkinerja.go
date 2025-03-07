package helper

import (
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"fmt"
	"sort"
)

func BuildTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, tematik domain.PohonKinerja) pohonkinerja.TematikResponse {
	tematikResp := pohonkinerja.TematikResponse{
		Id:          tematik.Id,
		Parent:      nil,
		Tema:        tematik.NamaPohon,
		JenisPohon:  tematik.JenisPohon,
		LevelPohon:  tematik.LevelPohon,
		Keterangan:  tematik.Keterangan,
		CountReview: tematik.CountReview,
		IsActive:    tematik.IsActive,
		Indikators:  ConvertToIndikatorResponses(tematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 0
	if strategics := pohonMap[4][tematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	// Tambahkan subtematik (level 1)
	if subTematiks := pohonMap[1][tematik.Id]; len(subTematiks) > 0 {
		for _, subTematik := range subTematiks {
			subTematikResp := BuildSubTematikResponse(pohonMap, subTematik)
			childs = append(childs, subTematikResp)
		}
	}

	tematikResp.Child = childs
	return tematikResp
}

func BuildSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subTematik domain.PohonKinerja) pohonkinerja.SubtematikResponse {
	subTematikResp := pohonkinerja.SubtematikResponse{
		Id:          subTematik.Id,
		Parent:      subTematik.Parent,
		Tema:        subTematik.NamaPohon,
		JenisPohon:  subTematik.JenisPohon,
		LevelPohon:  subTematik.LevelPohon,
		Keterangan:  subTematik.Keterangan,
		CountReview: subTematik.CountReview,
		IsActive:    subTematik.IsActive,
		Indikators:  ConvertToIndikatorResponses(subTematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 1
	if strategics := pohonMap[4][subTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	// Tambahkan subsubtematik (level 2)
	if subSubTematiks := pohonMap[2][subTematik.Id]; len(subSubTematiks) > 0 {
		for _, subSubTematik := range subSubTematiks {
			subSubTematikResp := BuildSubSubTematikResponse(pohonMap, subSubTematik)
			childs = append(childs, subSubTematikResp)
		}
	}

	subTematikResp.Child = childs
	return subTematikResp
}

func BuildSubSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subSubTematik domain.PohonKinerja) pohonkinerja.SubSubTematikResponse {
	subSubTematikResp := pohonkinerja.SubSubTematikResponse{
		Id:          subSubTematik.Id,
		Parent:      subSubTematik.Parent,
		Tema:        subSubTematik.NamaPohon,
		JenisPohon:  subSubTematik.JenisPohon,
		LevelPohon:  subSubTematik.LevelPohon,
		Keterangan:  subSubTematik.Keterangan,
		CountReview: subSubTematik.CountReview,
		IsActive:    subSubTematik.IsActive,
		Indikators:  ConvertToIndikatorResponses(subSubTematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 2
	if strategics := pohonMap[4][subSubTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	// Tambahkan supersubtematik (level 3)
	if superSubTematiks := pohonMap[3][subSubTematik.Id]; len(superSubTematiks) > 0 {
		for _, superSubTematik := range superSubTematiks {
			superSubTematikResp := BuildSuperSubTematikResponse(pohonMap, superSubTematik)
			childs = append(childs, superSubTematikResp)
		}
	}

	subSubTematikResp.Child = childs
	return subSubTematikResp
}

func BuildSuperSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, superSubTematik domain.PohonKinerja) pohonkinerja.SuperSubTematikResponse {
	superSubTematikResp := pohonkinerja.SuperSubTematikResponse{
		Id:          superSubTematik.Id,
		Parent:      superSubTematik.Parent,
		Tema:        superSubTematik.NamaPohon,
		JenisPohon:  superSubTematik.JenisPohon,
		LevelPohon:  superSubTematik.LevelPohon,
		Keterangan:  superSubTematik.Keterangan,
		CountReview: superSubTematik.CountReview,
		IsActive:    superSubTematik.IsActive,
		Indikators:  ConvertToIndikatorResponses(superSubTematik.Indikator),
	}

	var childs []interface{}

	// Tambahkan strategic (level 4) yang memiliki parent level 3
	if strategics := pohonMap[4][superSubTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	superSubTematikResp.Childs = childs
	return superSubTematikResp
}

func BuildStrategicResponse(pohonMap map[int]map[int][]domain.PohonKinerja, strategic domain.PohonKinerja) pohonkinerja.StrategicResponse {
	// Tambahkan map untuk melacak indikator yang sudah diproses
	processedIndikators := make(map[string]bool)
	var uniqueIndikators []pohonkinerja.IndikatorResponse

	// Proses indikator dengan pengecekan duplikasi
	for _, ind := range strategic.Indikator {
		if !processedIndikators[ind.Id] {
			processedIndikators[ind.Id] = true

			// Buat map untuk melacak target yang unik
			processedTargets := make(map[string]bool)
			var uniqueTargets []pohonkinerja.TargetResponse

			// Proses target dengan pengecekan duplikasi
			for _, target := range ind.Target {
				if !processedTargets[target.Id] {
					processedTargets[target.Id] = true
					targetResp := pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					}
					uniqueTargets = append(uniqueTargets, targetResp)
				}
			}

			indResp := pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				IdPokin:       fmt.Sprint(strategic.Id),
				NamaIndikator: ind.Indikator,
				Target:        uniqueTargets,
			}
			uniqueIndikators = append(uniqueIndikators, indResp)
		}
	}

	strategicResp := pohonkinerja.StrategicResponse{
		Id:          strategic.Id,
		Parent:      strategic.Parent,
		Strategi:    strategic.NamaPohon,
		JenisPohon:  strategic.JenisPohon,
		LevelPohon:  strategic.LevelPohon,
		Keterangan:  strategic.Keterangan,
		Status:      strategic.Status,
		Indikators:  uniqueIndikators,
		CountReview: strategic.CountReview,
		IsActive:    strategic.IsActive,
		KodeOpd: &opdmaster.OpdResponseForAll{
			KodeOpd: strategic.KodeOpd,
			NamaOpd: strategic.NamaOpd,
		},
		Pelaksana: ConvertToPelaksanaResponses(strategic.Pelaksana),
	}

	var childs []interface{}

	// Tambahkan tactical (level 5) ke childs
	if tacticals := pohonMap[5][strategic.Id]; len(tacticals) > 0 {
		// Urutkan tactical berdasarkan Id
		sort.Slice(tacticals, func(i, j int) bool {
			return tacticals[i].Id < tacticals[j].Id
		})

		for _, tactical := range tacticals {
			tacticalResp := BuildTacticalResponse(pohonMap, tactical)
			childs = append(childs, tacticalResp)
		}
	}

	strategicResp.Childs = childs
	return strategicResp
}

func BuildTacticalResponse(pohonMap map[int]map[int][]domain.PohonKinerja, tactical domain.PohonKinerja) pohonkinerja.TacticalResponse {
	// Proses indikator dengan pengecekan duplikasi
	processedIndikators := make(map[string]bool)
	var uniqueIndikators []pohonkinerja.IndikatorResponse

	for _, ind := range tactical.Indikator {
		if !processedIndikators[ind.Id] {
			processedIndikators[ind.Id] = true

			// Buat map untuk melacak target yang unik
			processedTargets := make(map[string]bool)
			var uniqueTargets []pohonkinerja.TargetResponse

			for _, target := range ind.Target {
				if !processedTargets[target.Id] {
					processedTargets[target.Id] = true
					targetResp := pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					}
					uniqueTargets = append(uniqueTargets, targetResp)
				}
			}

			indResp := pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				IdPokin:       fmt.Sprint(tactical.Id),
				NamaIndikator: ind.Indikator,
				Target:        uniqueTargets,
			}
			uniqueIndikators = append(uniqueIndikators, indResp)
		}
	}

	tacticalResp := pohonkinerja.TacticalResponse{
		Id:          tactical.Id,
		Parent:      tactical.Parent,
		Strategi:    tactical.NamaPohon,
		JenisPohon:  tactical.JenisPohon,
		LevelPohon:  tactical.LevelPohon,
		Keterangan:  &tactical.Keterangan,
		Status:      tactical.Status,
		Indikators:  uniqueIndikators,
		CountReview: tactical.CountReview,
		IsActive:    tactical.IsActive,
		Pelaksana:   ConvertToPelaksanaResponses(tactical.Pelaksana),
	}

	// Tambahkan data OPD jika ada
	if tactical.KodeOpd != "" {
		tacticalResp.KodeOpd = &opdmaster.OpdResponseForAll{
			KodeOpd: tactical.KodeOpd,
			NamaOpd: tactical.NamaOpd,
		}
	}

	var childs []interface{}

	// Tambahkan operational ke childs
	if operationals := pohonMap[6][tactical.Id]; len(operationals) > 0 {
		for _, operational := range operationals {
			operationalResp := BuildOperationalResponse(pohonMap, operational)
			childs = append(childs, operationalResp)
		}
	}

	tacticalResp.Childs = childs
	return tacticalResp
}

func BuildOperationalResponse(pohonMap map[int]map[int][]domain.PohonKinerja, operational domain.PohonKinerja) pohonkinerja.OperationalResponse {
	// Proses indikator dengan pengecekan duplikasi
	processedIndikators := make(map[string]bool)
	var uniqueIndikators []pohonkinerja.IndikatorResponse

	for _, ind := range operational.Indikator {
		if !processedIndikators[ind.Id] {
			processedIndikators[ind.Id] = true

			// Buat map untuk melacak target yang unik
			processedTargets := make(map[string]bool)
			var uniqueTargets []pohonkinerja.TargetResponse

			for _, target := range ind.Target {
				if !processedTargets[target.Id] {
					processedTargets[target.Id] = true
					targetResp := pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					}
					uniqueTargets = append(uniqueTargets, targetResp)
				}
			}

			indResp := pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				IdPokin:       fmt.Sprint(operational.Id),
				NamaIndikator: ind.Indikator,
				Target:        uniqueTargets,
			}
			uniqueIndikators = append(uniqueIndikators, indResp)
		}
	}

	operationalResp := pohonkinerja.OperationalResponse{
		Id:          operational.Id,
		Parent:      operational.Parent,
		Strategi:    operational.NamaPohon,
		JenisPohon:  operational.JenisPohon,
		LevelPohon:  operational.LevelPohon,
		Keterangan:  &operational.Keterangan,
		Status:      operational.Status,
		Indikators:  uniqueIndikators,
		CountReview: operational.CountReview,
		IsActive:    operational.IsActive,
		Pelaksana:   ConvertToPelaksanaResponses(operational.Pelaksana),
	}

	// Tambahkan data OPD jika ada
	if operational.KodeOpd != "" {
		operationalResp.KodeOpd = &opdmaster.OpdResponseForAll{
			KodeOpd: operational.KodeOpd,
			NamaOpd: operational.NamaOpd,
		}
	}

	var childs []interface{}

	// Cek level berikutnya (operational-n)
	nextLevel := operational.LevelPohon + 1
	if operationalNs := pohonMap[nextLevel][operational.Id]; len(operationalNs) > 0 {
		// Urutkan berdasarkan Id
		sort.Slice(operationalNs, func(i, j int) bool {
			return operationalNs[i].Id < operationalNs[j].Id
		})

		for _, opN := range operationalNs {
			operationalNResp := BuildOperationalNResponse(pohonMap, opN)
			childs = append(childs, operationalNResp)
		}
	}

	operationalResp.Childs = childs
	return operationalResp
}

func BuildOperationalNResponse(pohonMap map[int]map[int][]domain.PohonKinerja, operationalN domain.PohonKinerja) pohonkinerja.OperationalNResponse {
	// Proses indikator dengan pengecekan duplikasi
	processedIndikators := make(map[string]bool)
	var uniqueIndikators []pohonkinerja.IndikatorResponse

	for _, ind := range operationalN.Indikator {
		if !processedIndikators[ind.Id] {
			processedIndikators[ind.Id] = true

			// Buat map untuk melacak target yang unik
			processedTargets := make(map[string]bool)
			var uniqueTargets []pohonkinerja.TargetResponse

			for _, target := range ind.Target {
				if !processedTargets[target.Id] {
					processedTargets[target.Id] = true
					targetResp := pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					}
					uniqueTargets = append(uniqueTargets, targetResp)
				}
			}

			indResp := pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				IdPokin:       fmt.Sprint(operationalN.Id),
				NamaIndikator: ind.Indikator,
				Target:        uniqueTargets,
			}
			uniqueIndikators = append(uniqueIndikators, indResp)
		}
	}

	operationalNResp := pohonkinerja.OperationalNResponse{
		Id:          operationalN.Id,
		Parent:      operationalN.Parent,
		Strategi:    operationalN.NamaPohon,
		JenisPohon:  operationalN.JenisPohon,
		LevelPohon:  operationalN.LevelPohon,
		Keterangan:  &operationalN.Keterangan,
		Status:      operationalN.Status,
		Indikators:  uniqueIndikators,
		CountReview: operationalN.CountReview,
		IsActive:    operationalN.IsActive,
		Pelaksana:   ConvertToPelaksanaResponses(operationalN.Pelaksana),
	}

	// Tambahkan data OPD jika ada
	if operationalN.KodeOpd != "" {
		operationalNResp.KodeOpd = &opdmaster.OpdResponseForAll{
			KodeOpd: operationalN.KodeOpd,
			NamaOpd: operationalN.NamaOpd,
		}
	}

	// Cek level berikutnya secara rekursif
	nextLevel := operationalN.LevelPohon + 1
	if nextOperationalNs := pohonMap[nextLevel][operationalN.Id]; len(nextOperationalNs) > 0 {
		// Urutkan berdasarkan Id
		sort.Slice(nextOperationalNs, func(i, j int) bool {
			return nextOperationalNs[i].Id < nextOperationalNs[j].Id
		})

		var childs []pohonkinerja.OperationalNResponse
		for _, nextOpN := range nextOperationalNs {
			childResp := BuildOperationalNResponse(pohonMap, nextOpN)
			childs = append(childs, childResp)
		}
		operationalNResp.Childs = childs
	}

	return operationalNResp
}

func BuildSubTematikResponseLimited(pohonMap map[int]map[int][]domain.PohonKinerja, subTematik domain.PohonKinerja) pohonkinerja.SubtematikResponse {
	subTematikResp := pohonkinerja.SubtematikResponse{
		Id:          subTematik.Id,
		Parent:      subTematik.Parent,
		Tema:        subTematik.NamaPohon,
		JenisPohon:  subTematik.JenisPohon,
		LevelPohon:  subTematik.LevelPohon,
		Keterangan:  subTematik.Keterangan,
		CountReview: subTematik.CountReview,
		IsActive:    subTematik.IsActive,
		Indikators:  ConvertToIndikatorResponses(subTematik.Indikator),
	}

	var childs []interface{}

	// Hanya tambahkan strategic (level 4) yang memiliki parent level 1
	if strategics := pohonMap[4][subTematik.Id]; len(strategics) > 0 {
		for _, strategic := range strategics {
			strategicResp := BuildStrategicResponse(pohonMap, strategic)
			childs = append(childs, strategicResp)
		}
	}

	subTematikResp.Child = childs
	return subTematikResp
}

func ConvertToPelaksanaResponses(pelaksanas []domain.PelaksanaPokin) []pohonkinerja.PelaksanaOpdResponse {
	var responses []pohonkinerja.PelaksanaOpdResponse
	for _, p := range pelaksanas {
		responses = append(responses, pohonkinerja.PelaksanaOpdResponse{
			Id:          p.Id,
			PegawaiId:   p.PegawaiId,
			NamaPegawai: p.NamaPegawai,
		})
	}
	return responses
}
