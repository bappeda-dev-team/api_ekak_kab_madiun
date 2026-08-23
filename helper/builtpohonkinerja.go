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
		Id:           tematik.Id,
		Parent:       nil,
		Tema:         tematik.NamaPohon,
		JenisPohon:   tematik.JenisPohon,
		LevelPohon:   tematik.LevelPohon,
		Keterangan:   tematik.Keterangan,
		CountReview:  tematik.CountReview,
		IsActive:     tematik.IsActive,
		Indikators:   ConvertToIndikatorResponses(tematik.Indikator),
		TaggingPokin: ConvertToTaggingResponses(tematik.TaggingPokin), // Tambahkan tagging
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
		Id:           subTematik.Id,
		Parent:       subTematik.Parent,
		Tema:         subTematik.NamaPohon,
		JenisPohon:   subTematik.JenisPohon,
		LevelPohon:   subTematik.LevelPohon,
		Keterangan:   subTematik.Keterangan,
		CountReview:  subTematik.CountReview,
		IsActive:     subTematik.IsActive,
		Indikators:   ConvertToIndikatorResponses(subTematik.Indikator),
		TaggingPokin: ConvertToTaggingResponses(subTematik.TaggingPokin),
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
		Id:           subSubTematik.Id,
		Parent:       subSubTematik.Parent,
		Tema:         subSubTematik.NamaPohon,
		JenisPohon:   subSubTematik.JenisPohon,
		LevelPohon:   subSubTematik.LevelPohon,
		Keterangan:   subSubTematik.Keterangan,
		CountReview:  subSubTematik.CountReview,
		IsActive:     subSubTematik.IsActive,
		Indikators:   ConvertToIndikatorResponses(subSubTematik.Indikator),
		TaggingPokin: ConvertToTaggingResponses(subSubTematik.TaggingPokin),
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
		Id:           superSubTematik.Id,
		Parent:       superSubTematik.Parent,
		Tema:         superSubTematik.NamaPohon,
		JenisPohon:   superSubTematik.JenisPohon,
		LevelPohon:   superSubTematik.LevelPohon,
		Keterangan:   superSubTematik.Keterangan,
		CountReview:  superSubTematik.CountReview,
		IsActive:     superSubTematik.IsActive,
		Indikators:   ConvertToIndikatorResponses(superSubTematik.Indikator),
		TaggingPokin: ConvertToTaggingResponses(superSubTematik.TaggingPokin),
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
	// Modifikasi bagian tagging
	var taggingResponses []pohonkinerja.TaggingResponse
	for _, tagging := range strategic.TaggingPokin {
		var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
				Id:                  keterangan.Id,
				IdTagging:           keterangan.IdTagging,
				KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
				RencanaImplementasi: keterangan.RencanaImplementasi,
				Tahun:               keterangan.Tahun,
			})
		}

		taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
			Id:                       tagging.Id,
			IdPokin:                  tagging.IdPokin,
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganResponses,
			CloneFrom:                tagging.CloneFrom,
		})
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
		Pelaksana:    ConvertToPelaksanaResponses(strategic.Pelaksana),
		TaggingPokin: taggingResponses,
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
		Id:           tactical.Id,
		Parent:       tactical.Parent,
		Strategi:     tactical.NamaPohon,
		JenisPohon:   tactical.JenisPohon,
		LevelPohon:   tactical.LevelPohon,
		Keterangan:   &tactical.Keterangan,
		Status:       tactical.Status,
		Indikators:   uniqueIndikators,
		CountReview:  tactical.CountReview,
		IsActive:     tactical.IsActive,
		Pelaksana:    ConvertToPelaksanaResponses(tactical.Pelaksana),
		TaggingPokin: ConvertToTaggingResponses(tactical.TaggingPokin),
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
		Id:           operational.Id,
		Parent:       operational.Parent,
		Strategi:     operational.NamaPohon,
		JenisPohon:   operational.JenisPohon,
		LevelPohon:   operational.LevelPohon,
		Keterangan:   &operational.Keterangan,
		Status:       operational.Status,
		Indikators:   uniqueIndikators,
		CountReview:  operational.CountReview,
		IsActive:     operational.IsActive,
		Pelaksana:    ConvertToPelaksanaResponses(operational.Pelaksana),
		TaggingPokin: ConvertToTaggingResponses(operational.TaggingPokin),
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
		Id:           operationalN.Id,
		Parent:       operationalN.Parent,
		Strategi:     operationalN.NamaPohon,
		JenisPohon:   operationalN.JenisPohon,
		LevelPohon:   operationalN.LevelPohon,
		Keterangan:   &operationalN.Keterangan,
		Status:       operationalN.Status,
		Indikators:   uniqueIndikators,
		CountReview:  operationalN.CountReview,
		IsActive:     operationalN.IsActive,
		Pelaksana:    ConvertToPelaksanaResponses(operationalN.Pelaksana),
		TaggingPokin: ConvertToTaggingResponses(operationalN.TaggingPokin),
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

// ─────────────────────────────────────────────────────────────────────────────
// OPD View helpers – hierarki baru yang menggabungkan pohon pemda (level 0-4)
// dengan pohon OPD (level 4-6) melalui field clone_from pada strategic pemda.
// ─────────────────────────────────────────────────────────────────────────────

// OpdStrategicLookupKey membentuk key lookup strategic OPD dari id strategic pemda + kode_opd.
func OpdStrategicLookupKey(strategicPemdaId int, kodeOpd string) string {
	return fmt.Sprintf("%d|%s", strategicPemdaId, kodeOpd)
}

// resolveOpdStrategic mencari strategic OPD dari map lookup berdasarkan strategic pemda.
func resolveOpdStrategic(sp domain.PohonKinerja, cloneToOpdStrategic map[string]domain.PohonKinerja) (domain.PohonKinerja, bool) {
	candidates := []string{
		OpdStrategicLookupKey(sp.Id, sp.KodeOpd),
		OpdStrategicLookupKey(sp.Id, ""),
	}
	if sp.CloneFrom != 0 {
		candidates = append(candidates,
			OpdStrategicLookupKey(sp.CloneFrom, sp.KodeOpd),
			OpdStrategicLookupKey(sp.CloneFrom, ""),
		)
	}
	for _, key := range candidates {
		if opd, ok := cloneToOpdStrategic[key]; ok {
			return opd, true
		}
	}
	return domain.PohonKinerja{}, false
}

// buildOpdGroupsForParent mengelompokkan strategic pemda (level 4) yang merupakan
// anak langsung dari parentId berdasarkan kode_opd. Di dalam setiap OPD group,
// strategic di-resolve ke OPD asalnya lalu dikelompokkan lagi berdasarkan
// BidangUrusan → TujuanOpd (dengan indikator+target+sasaran_opd) → Strategic.
func buildOpdGroupsForParent(
	pohonMapPemda map[int]map[int][]domain.PohonKinerja,
	pohonMapOpd map[int]map[int][]domain.PohonKinerja,
	cloneToOpdStrategic map[string]domain.PohonKinerja,
	opdNamaMap map[string]string,
	tujuanOpdMap map[int]domain.TujuanOpd,
	parentId int,
) []interface{} {
	strategicPemda, ok := pohonMapPemda[4][parentId]
	if !ok || len(strategicPemda) == 0 {
		return nil
	}

	// Pertahankan urutan kode_opd sesuai kemunculan pertama
	order := []string{}
	groupMap := map[string][]domain.PohonKinerja{}
	for _, sp := range strategicPemda {
		if _, exists := groupMap[sp.KodeOpd]; !exists {
			order = append(order, sp.KodeOpd)
		}
		groupMap[sp.KodeOpd] = append(groupMap[sp.KodeOpd], sp)
	}

	var childs []interface{}
	for _, kodeOpd := range order {
		strategics := groupMap[kodeOpd]

		// Resolve setiap strategic pemda → OPD strategic node
		type resolvedStrategic struct {
			opdNode domain.PohonKinerja
		}
		var resolved []resolvedStrategic
		for _, sp := range strategics {
			opdNode, found := resolveOpdStrategic(sp, cloneToOpdStrategic)
			if !found {
				continue
			}
			resolved = append(resolved, resolvedStrategic{opdNode: opdNode})
		}
		if len(resolved) == 0 {
			continue
		}

		// ── Kelompokkan per BidangUrusan → per TujuanOpd ──
		// Menggunakan kode_bidang_urusan sebagai key pengelompokan.
		// Satu OPD strategic bisa punya multiple SasaranInfo → pilih tujuan_opd pertama per strategic
		// agar tidak tampil di lebih dari satu grup TujuanOpd.
		type tujuanKey struct {
			kodeBidangUrusan string
			idTujuanOpd      int
		}
		bidangOrder := []string{}
		bidangSeen := map[string]bool{}
		// bidangKode → tujuanId order
		tujuanOrder := map[string][]int{}
		tujuanSeen := map[tujuanKey]bool{}
		// tujuanKey → strategic nodes
		tujuanStrategics := map[tujuanKey][]domain.PohonKinerja{}
		// tujuanKey → sasaran names (dedup)
		tujuanSasaran := map[tujuanKey][]string{}
		tujuanSasaranSeen := map[tujuanKey]map[string]bool{}
		// bidangKode → nama
		bidangNama := map[string]string{}

		for _, rs := range resolved {
			node := rs.opdNode
			// Ambil info sasaran dari node (sudah di-populate di service)
			infos := node.SasaranInfo

			if len(infos) == 0 {
				// Tidak ada sasaran OPD terhubung: masukkan ke grup kosong
				tk := tujuanKey{kodeBidangUrusan: "", idTujuanOpd: 0}
				if !bidangSeen[""] {
					bidangSeen[""] = true
					bidangOrder = append(bidangOrder, "")
				}
				if !tujuanSeen[tk] {
					tujuanSeen[tk] = true
					tujuanOrder[""] = append(tujuanOrder[""], 0)
					tujuanSasaranSeen[tk] = map[string]bool{}
				}
				tujuanStrategics[tk] = append(tujuanStrategics[tk], node)
				continue
			}

			// Gunakan tujuan_opd pertama sebagai kunci pengelompokan utama strategic ini
			primary := infos[0]
			tk := tujuanKey{kodeBidangUrusan: primary.KodeBidangUrusan, idTujuanOpd: primary.IdTujuanOpd}

			if !bidangSeen[primary.KodeBidangUrusan] {
				bidangSeen[primary.KodeBidangUrusan] = true
				bidangOrder = append(bidangOrder, primary.KodeBidangUrusan)
				bidangNama[primary.KodeBidangUrusan] = primary.NamaBidangUrusan
			}
			if !tujuanSeen[tk] {
				tujuanSeen[tk] = true
				tujuanOrder[primary.KodeBidangUrusan] = append(tujuanOrder[primary.KodeBidangUrusan], primary.IdTujuanOpd)
				tujuanSasaranSeen[tk] = map[string]bool{}
			}
			tujuanStrategics[tk] = append(tujuanStrategics[tk], node)

			// Kumpulkan semua nama sasaran yang terhubung dengan tujuan_opd ini
			for _, info := range infos {
				if info.IdTujuanOpd != primary.IdTujuanOpd {
					continue
				}
				if !tujuanSasaranSeen[tk][info.NamaSasaranOpd] {
					tujuanSasaranSeen[tk][info.NamaSasaranOpd] = true
					tujuanSasaran[tk] = append(tujuanSasaran[tk], info.NamaSasaranOpd)
				}
			}
		}

		// ── Bangun hierarki BidangUrusan → TujuanOpd → Strategic ──
		var bidangChilds []interface{}
		for _, kodeBidang := range bidangOrder {
			var tujuanChilds []interface{}
			for _, idTujuan := range tujuanOrder[kodeBidang] {
				tk := tujuanKey{kodeBidangUrusan: kodeBidang, idTujuanOpd: idTujuan}

				// Build strategic responses untuk tujuan ini
				var strategicChilds []interface{}
				for _, opdNode := range tujuanStrategics[tk] {
					strategicResp := BuildStrategicResponse(pohonMapOpd, opdNode)
					strategicChilds = append(strategicChilds, strategicResp)
				}

				// Indikator tujuan_opd (dengan target + tahun)
				var indikatorResps []pohonkinerja.IndikatorResponse
				if tj, ok := tujuanOpdMap[idTujuan]; ok {
					for _, ind := range tj.Indikator {
						var targets []pohonkinerja.TargetResponse
						for _, t := range ind.Target {
							targets = append(targets, pohonkinerja.TargetResponse{
								Id:              t.Id,
								IndikatorId:     t.IndikatorId,
								TargetIndikator: t.Target,
								SatuanIndikator: t.Satuan,
								TahunSasaran:    t.Tahun,
							})
						}
						indikatorResps = append(indikatorResps, pohonkinerja.IndikatorResponse{
							Id:            ind.Id,
							NamaIndikator: ind.Indikator,
							Target:        targets,
						})
					}
				}

				// Nama tujuan_opd
				namaTujuan := ""
				if tj, ok := tujuanOpdMap[idTujuan]; ok {
					namaTujuan = tj.Tujuan
				} else if len(tujuanStrategics[tk]) > 0 {
					// fallback: ambil dari SasaranInfo
					for _, info := range tujuanStrategics[tk][0].SasaranInfo {
						if info.IdTujuanOpd == idTujuan {
							namaTujuan = info.NamaTujuanOpd
							break
						}
					}
				}

				tujuanChilds = append(tujuanChilds, pohonkinerja.TujuanOpdStrategicGroupResponse{
					Id:            idTujuan,
					NamaTujuanOpd: namaTujuan,
					Indikators:    indikatorResps,
					SasaranOpd:    tujuanSasaran[tk],
					Childs:        strategicChilds,
				})
			}

			bidangChilds = append(bidangChilds, pohonkinerja.BidangUrusanGroupResponse{
				NamaBidangUrusan: bidangNama[kodeBidang],
				Childs:           tujuanChilds,
			})
		}

		displayKode := kodeOpd
		if displayKode == "" && len(resolved) > 0 {
			displayKode = resolved[0].opdNode.KodeOpd
		}

		childs = append(childs, pohonkinerja.OpdGroupResponse{
			KodeOpd: displayKode,
			NamaOpd: opdNamaMap[displayKode],
			Childs:  bidangChilds,
		})
	}
	return childs
}

// buildSubTematikOpdViewChilds membangun childs untuk sebuah node subtematik (atau
// node tematik) dalam OPD View: sub-level (SubSubTematik, dst) ditampilkan rekursif,
// sedangkan strategic pemda di-resolve ke OPD Group → BidangUrusan → TujuanOpd.
func buildSubTematikOpdViewChilds(
	pohonMapPemda map[int]map[int][]domain.PohonKinerja,
	pohonMapOpd map[int]map[int][]domain.PohonKinerja,
	cloneToOpdStrategic map[string]domain.PohonKinerja,
	opdNamaMap map[string]string,
	tujuanOpdMap map[int]domain.TujuanOpd,
	nodeId int,
	currentLevel int, // level dari node saat ini (1=subtematik, 2=subsubtematik, 3=supersubtematik)
) []interface{} {
	var childs []interface{}

	// OPD Groups dari strategic pemda yang langsung di bawah node ini
	opdGroups := buildOpdGroupsForParent(pohonMapPemda, pohonMapOpd, cloneToOpdStrategic, opdNamaMap, tujuanOpdMap, nodeId)
	childs = append(childs, opdGroups...)

	// Sub-level berikutnya (mis. subsubtematik di bawah subtematik)
	nextLevel := currentLevel + 1
	if nextLevel <= 3 {
		if subNodes, ok := pohonMapPemda[nextLevel][nodeId]; ok {
			sort.Slice(subNodes, func(i, j int) bool { return subNodes[i].Id < subNodes[j].Id })
			for _, sub := range subNodes {
				subChilds := buildSubTematikOpdViewChilds(pohonMapPemda, pohonMapOpd, cloneToOpdStrategic, opdNamaMap, tujuanOpdMap, sub.Id, nextLevel)

				switch nextLevel {
				case 2:
					resp := pohonkinerja.SubSubTematikResponse{
						Id:           sub.Id,
						Parent:       sub.Parent,
						Tema:         sub.NamaPohon,
						JenisPohon:   sub.JenisPohon,
						LevelPohon:   sub.LevelPohon,
						Keterangan:   sub.Keterangan,
						CountReview:  sub.CountReview,
						IsActive:     sub.IsActive,
						Indikators:   ConvertToIndikatorResponses(sub.Indikator),
						TaggingPokin: ConvertToTaggingResponses(sub.TaggingPokin),
						Child:        subChilds,
					}
					childs = append(childs, resp)
				case 3:
					resp := pohonkinerja.SuperSubTematikResponse{
						Id:           sub.Id,
						Parent:       sub.Parent,
						Tema:         sub.NamaPohon,
						JenisPohon:   sub.JenisPohon,
						LevelPohon:   sub.LevelPohon,
						Keterangan:   sub.Keterangan,
						CountReview:  sub.CountReview,
						IsActive:     sub.IsActive,
						Indikators:   ConvertToIndikatorResponses(sub.Indikator),
						TaggingPokin: ConvertToTaggingResponses(sub.TaggingPokin),
						Childs:       subChilds,
					}
					childs = append(childs, resp)
				}
			}
		}
	}

	return childs
}

// BuildTematikOpdViewResponse membangun TematikResponse untuk endpoint OPD View.
// Perbedaan dari BuildTematikResponse biasa:
//   - Strategic pemda (level 4) di-resolve ke strategic OPD via clone_from
//   - Di dalam OpdGroup, strategic dikelompokkan: BidangUrusan → TujuanOpd → Strategic
//   - tujuanOpdMap berisi data tujuan OPD lengkap (indikator + target renstra)
//   - Tactical & operational diambil dari pohon OPD (pohonMapOpd), bukan pohon pemda
func BuildTematikOpdViewResponse(
	pohonMapPemda map[int]map[int][]domain.PohonKinerja,
	pohonMapOpd map[int]map[int][]domain.PohonKinerja,
	cloneToOpdStrategic map[string]domain.PohonKinerja,
	opdNamaMap map[string]string,
	tujuanOpdMap map[int]domain.TujuanOpd,
	tematik domain.PohonKinerja,
) pohonkinerja.TematikResponse {
	var childs []interface{}

	// OPD Groups dari strategic pemda langsung di bawah tematik
	opdGroups := buildOpdGroupsForParent(pohonMapPemda, pohonMapOpd, cloneToOpdStrategic, opdNamaMap, tujuanOpdMap, tematik.Id)
	childs = append(childs, opdGroups...)

	// Subtematik (level 1) di bawah tematik
	if subTematiks, ok := pohonMapPemda[1][tematik.Id]; ok {
		sort.Slice(subTematiks, func(i, j int) bool { return subTematiks[i].Id < subTematiks[j].Id })
		for _, st := range subTematiks {
			subChilds := buildSubTematikOpdViewChilds(pohonMapPemda, pohonMapOpd, cloneToOpdStrategic, opdNamaMap, tujuanOpdMap, st.Id, 1)

			subResp := pohonkinerja.SubtematikResponse{
				Id:           st.Id,
				Parent:       st.Parent,
				Tema:         st.NamaPohon,
				JenisPohon:   st.JenisPohon,
				LevelPohon:   st.LevelPohon,
				Keterangan:   st.Keterangan,
				CountReview:  st.CountReview,
				IsActive:     st.IsActive,
				Indikators:   ConvertToIndikatorResponses(st.Indikator),
				TaggingPokin: ConvertToTaggingResponses(st.TaggingPokin),
				Child:        subChilds,
			}
			childs = append(childs, subResp)
		}
	}

	var uniqueIndikators []pohonkinerja.IndikatorResponse
	seen := make(map[string]bool)
	for _, ind := range tematik.Indikator {
		if !seen[ind.Id] {
			seen[ind.Id] = true
			uniqueIndikators = append(uniqueIndikators, ConvertToIndikatorResponse(ind))
		}
	}

	return pohonkinerja.TematikResponse{
		Id:           tematik.Id,
		Parent:       nil,
		Tema:         tematik.NamaPohon,
		JenisPohon:   tematik.JenisPohon,
		LevelPohon:   tematik.LevelPohon,
		Keterangan:   tematik.Keterangan,
		CountReview:  tematik.CountReview,
		IsActive:     tematik.IsActive,
		Indikators:   uniqueIndikators,
		TaggingPokin: ConvertToTaggingResponses(tematik.TaggingPokin),
		Child:        childs,
	}
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
