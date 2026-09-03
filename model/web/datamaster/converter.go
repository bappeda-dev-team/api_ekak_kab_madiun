package datamaster

import "ekak_kabupaten_madiun/model/domain/datamaster"

// Convert RBRequest → MasterRB (entity untuk DB)
func ConvertRBRequestToMaster(rbReq RBRequest, userId int) datamaster.MasterRB {
	master := datamaster.MasterRB{
		JenisRB:       rbReq.JenisRB,
		KegiatanUtama: rbReq.KegiatanUtama,
		Keterangan:    rbReq.Keterangan,
		TahunBaseline: rbReq.TahunBaseline,
		TahunNext:     rbReq.TahunNext,
		LastUpdatedBy: userId,
		Indikator:     []datamaster.IndikatorRB{},
	}

	for _, indReq := range rbReq.Indikator {

		ind := datamaster.IndikatorRB{
			Indikator: indReq.Indikator,
			TargetRB:  []datamaster.TargetRB{},
		}

		for _, tReq := range indReq.Target {

			t := datamaster.TargetRB{
				// baseline defaults
				TahunBaseline:     getInt(tReq.TahunBaseline),
				TargetBaseline:    getInt(tReq.TargetBaseline),
				RealisasiBaseline: getFloat32(tReq.RealisasiBaseline),
				SatuanBaseline:    getString(tReq.SatuanBaseline),

				// next defaults
				TahunNext:  getInt(tReq.TahunNext),
				TargetNext: getInt(tReq.TargetNext),
				SatuanNext: getString(tReq.SatuanNext),
			}

			ind.TargetRB = append(ind.TargetRB, t)
		}

		master.Indikator = append(master.Indikator, ind)
	}

	return master
}

// Helpers untuk pointer-safe
func getInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func getFloat32(v *float32) float32 {
	if v == nil {
		return 0
	}
	return *v
}

func getString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
