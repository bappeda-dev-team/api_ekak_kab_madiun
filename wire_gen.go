// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"ekak_kabupaten_madiun/app"
	"ekak_kabupaten_madiun/controller"
	"ekak_kabupaten_madiun/dataseeder"
	"ekak_kabupaten_madiun/middleware"
	"ekak_kabupaten_madiun/repository"
	"ekak_kabupaten_madiun/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"net/http"
)

// Injectors from injector.go:

func InitializeServer() *http.Server {
	rencanaKinerjaRepositoryImpl := repository.NewRencanaKinerjaRepositoryImpl()
	db := app.GetConnection()
	v := _wireValue
	validate := validator.New(v...)
	opdRepositoryImpl := repository.NewOpdRepositoryImpl()
	rencanaAksiRepositoryImpl := repository.NewRencanaAksiRepositoryImpl()
	usulanMusrebangRepositoryImpl := repository.NewUsulanMusrebangRepositoryImpl()
	usulanMandatoriRepositoryImpl := repository.NewUsulanMandatoriRepositoryImpl()
	usulanPokokPikiranRepositoryImpl := repository.NewUsulanPokokPikiranRepositoryImpl()
	usulanInisiatifRepositoryImpl := repository.NewUsulanInisiatifRepositoryImpl()
	subKegiatanRepositoryImpl := repository.NewSubKegiatanRepositoryImpl()
	dasarHukumRepositoryImpl := repository.NewDasarHukumRepositoryImpl()
	gambaranUmumRepositoryImpl := repository.NewGambaranUmumRepositoryImpl()
	inovasiRepositoryImpl := repository.NewInovasiRepositoryImpl()
	pelaksanaanRencanaAksiRepositoryImpl := repository.NewPelaksanaanRencanaAksiRepositoryImpl()
	pegawaiRepositoryImpl := repository.NewPegawaiRepositoryImpl()
	pohonKinerjaRepositoryImpl := repository.NewPohonKinerjaRepositoryImpl()
	rencanaKinerjaServiceImpl := service.NewRencanaKinerjaServiceImpl(rencanaKinerjaRepositoryImpl, db, validate, opdRepositoryImpl, rencanaAksiRepositoryImpl, usulanMusrebangRepositoryImpl, usulanMandatoriRepositoryImpl, usulanPokokPikiranRepositoryImpl, usulanInisiatifRepositoryImpl, subKegiatanRepositoryImpl, dasarHukumRepositoryImpl, gambaranUmumRepositoryImpl, inovasiRepositoryImpl, pelaksanaanRencanaAksiRepositoryImpl, pegawaiRepositoryImpl, pohonKinerjaRepositoryImpl)
	rencanaKinerjaControllerImpl := controller.NewRencanaKinerjaControllerImpl(rencanaKinerjaServiceImpl)
	rencanaAksiServiceImpl := service.NewRencanaAksiServiceImpl(rencanaAksiRepositoryImpl, db, validate, pelaksanaanRencanaAksiRepositoryImpl)
	rencanaAksiControllerImpl := controller.NewRencanaAksiControllerImpl(rencanaAksiServiceImpl)
	pelaksanaanRencanaAksiServiceImpl := service.NewPelaksanaanRencanaAksiServiceImpl(pelaksanaanRencanaAksiRepositoryImpl, rencanaAksiRepositoryImpl, db)
	pelaksanaanRencanaAksiControllerImpl := controller.NewPelaksanaanRencanaAksiControllerImpl(pelaksanaanRencanaAksiServiceImpl)
	usulanMusrebangServiceImpl := service.NewUsulanMusrebangServiceImpl(usulanMusrebangRepositoryImpl, db)
	usulanMusrebangControllerImpl := controller.NewUsulanMusrebangControllerImpl(usulanMusrebangServiceImpl)
	usulanMandatoriServiceImpl := service.NewUsulanMandatoriServiceImpl(usulanMandatoriRepositoryImpl, db)
	usulanMandatoriControllerImpl := controller.NewUsulanMandatoriControllerImpl(usulanMandatoriServiceImpl)
	usulanPokokPikiranServiceImpl := service.NewUsulanPokokPikiranServiceImpl(usulanPokokPikiranRepositoryImpl, db)
	usulanPokokPikiranControllerImpl := controller.NewUsulanPokokPikiranControllerImpl(usulanPokokPikiranServiceImpl)
	usulanInisiatifServiceImpl := service.NewUsulanInisiatifServiceImpl(usulanInisiatifRepositoryImpl, db)
	usulanInisiatifControllerImpl := controller.NewUsulanInisiatifControllerImpl(usulanInisiatifServiceImpl)
	usulanTerpilihRepositoryImpl := repository.NewUsulanTerpilihRepositoryImpl()
	usulanTerpilihServiceImpl := service.NewUsulanTerpilihServiceImpl(usulanTerpilihRepositoryImpl, db, validate)
	usulanTerpilihControllerImpl := controller.NewUsulanTerpilihControllerImpl(usulanTerpilihServiceImpl)
	gambaranUmumServiceImpl := service.NewGambaranUmumServiceImpl(gambaranUmumRepositoryImpl, db)
	gambaranUmumControllerImpl := controller.NewGambaranUmumControllerImpl(gambaranUmumServiceImpl)
	dasarHukumServiceImpl := service.NewDasarHukumServiceImpl(dasarHukumRepositoryImpl, db)
	dasarHukumControllerImpl := controller.NewDasarHukumControllerImpl(dasarHukumServiceImpl)
	inovasiServiceImpl := service.NewInovasiServiceImpl(inovasiRepositoryImpl, db)
	inovasiControllerImpl := controller.NewInovasiControllerImpl(inovasiServiceImpl)
	subKegiatanServiceImpl := service.NewSubKegiatanServiceImpl(subKegiatanRepositoryImpl, opdRepositoryImpl, db, validate)
	subKegiatanControllerImpl := controller.NewSubKegiatanControllerImpl(subKegiatanServiceImpl)
	subKegiatanTerpilihRepositoryImpl := repository.NewSubKegiatanTerpilihRepositoryImpl()
	subKegiatanTerpilihServiceImpl := service.NewSubKegiatanTerpilihServiceImpl(rencanaKinerjaRepositoryImpl, subKegiatanRepositoryImpl, subKegiatanTerpilihRepositoryImpl, db)
	subKegiatanTerpilihControllerImpl := controller.NewSubKegiatanTerpilihControllerImpl(subKegiatanTerpilihServiceImpl)
	pohonKinerjaOpdServiceImpl := service.NewPohonKinerjaOpdServiceImpl(pohonKinerjaRepositoryImpl, opdRepositoryImpl, pegawaiRepositoryImpl, db)
	pohonKinerjaOpdControllerImpl := controller.NewPohonKinerjaOpdControllerImpl(pohonKinerjaOpdServiceImpl)
	pegawaiServiceImpl := service.NewPegawaiServiceImpl(pegawaiRepositoryImpl, opdRepositoryImpl, db)
	pegawaiControllerImpl := controller.NewPegawaiControllerImpl(pegawaiServiceImpl)
	lembagaRepositoryImpl := repository.NewLembagaRepositoryImpl()
	lembagaServiceImpl := service.NewLembagaServiceImpl(lembagaRepositoryImpl, db, validate)
	lembagaControllerImpl := controller.NewLembagaControllerImpl(lembagaServiceImpl)
	jabatanRepositoryImpl := repository.NewJabatanRepositoryImpl()
	jabatanServiceImpl := service.NewJabatanServiceImpl(jabatanRepositoryImpl, opdRepositoryImpl, db)
	jabatanControllerImpl := controller.NewJabatanControllerImpl(jabatanServiceImpl)
	pohonKinerjaAdminServiceImpl := service.NewPohonKinerjaAdminServiceImpl(pohonKinerjaRepositoryImpl, opdRepositoryImpl, db, pegawaiRepositoryImpl)
	pohonKinerjaAdminControllerImpl := controller.NewPohonKinerjaAdminControllerImpl(pohonKinerjaAdminServiceImpl)
	opdServiceImpl := service.NewOpdServiceImpl(opdRepositoryImpl, lembagaRepositoryImpl, db, validate)
	opdControllerImpl := controller.NewOpdControllerImpl(opdServiceImpl)
	programRepositoryImpl := repository.NewProgramRepositoryImpl()
	programServiceImpl := service.NewProgramServiceImpl(programRepositoryImpl, opdRepositoryImpl, db)
	programControllerImpl := controller.NewProgramControllerImpl(programServiceImpl)
	urusanRepositoryImpl := repository.NewUrusanRepositoryImpl()
	urusanServiceImpl := service.NewUrusanServiceImpl(urusanRepositoryImpl, db)
	urusanControllerImpl := controller.NewUrusanControllerImpl(urusanServiceImpl)
	bidangUrusanRepositoryImpl := repository.NewBidangUrusanRepositoryImpl()
	bidangUrusanServiceImpl := service.NewBidangUrusanServiceImpl(bidangUrusanRepositoryImpl, db)
	bidangUrusanControllerImpl := controller.NewBidangUrusanControllerImpl(bidangUrusanServiceImpl)
	kegiatanRepositoryImpl := repository.NewKegiatanRepositoryImpl()
	kegiatanServiceImpl := service.NewKegiatanServiceImpl(kegiatanRepositoryImpl, opdRepositoryImpl, db)
	kegiatanControllerImpl := controller.NewKegiatanControllerImpl(kegiatanServiceImpl)
	userRepositoryImpl := repository.NewUserRepositoryImpl()
	roleRepositoryImpl := repository.NewRoleRepositoryImpl()
	userServiceImpl := service.NewUserServiceImpl(userRepositoryImpl, roleRepositoryImpl, pegawaiRepositoryImpl, db)
	userControllerImpl := controller.NewUserControllerImpl(userServiceImpl)
	roleServiceImpl := service.NewRoleServiceImpl(roleRepositoryImpl, db)
	roleControllerImpl := controller.NewRoleControllerImpl(roleServiceImpl)
	tujuanOpdRepositoryImpl := repository.NewTujuanOpdRepositoryImpl()
	tujuanOpdServiceImpl := service.NewTujuanOpdServiceImpl(tujuanOpdRepositoryImpl, opdRepositoryImpl, db)
	tujuanOpdControllerImpl := controller.NewTujuanOpdControllerImpl(tujuanOpdServiceImpl)
	router := app.NewRouter(rencanaKinerjaControllerImpl, rencanaAksiControllerImpl, pelaksanaanRencanaAksiControllerImpl, usulanMusrebangControllerImpl, usulanMandatoriControllerImpl, usulanPokokPikiranControllerImpl, usulanInisiatifControllerImpl, usulanTerpilihControllerImpl, gambaranUmumControllerImpl, dasarHukumControllerImpl, inovasiControllerImpl, subKegiatanControllerImpl, subKegiatanTerpilihControllerImpl, pohonKinerjaOpdControllerImpl, pegawaiControllerImpl, lembagaControllerImpl, jabatanControllerImpl, pohonKinerjaAdminControllerImpl, opdControllerImpl, programControllerImpl, urusanControllerImpl, bidangUrusanControllerImpl, kegiatanControllerImpl, userControllerImpl, roleControllerImpl, tujuanOpdControllerImpl)
	authMiddleware := middleware.NewAuthMiddleware(router)
	server := NewServer(authMiddleware)
	return server
}

var (
	_wireValue = []validator.Option{}
)

func InitializeSeeder() dataseeder.Seeder {
	db := app.GetConnection()
	roleRepositoryImpl := repository.NewRoleRepositoryImpl()
	roleSeederImpl := dataseeder.NewRoleSeederImpl(roleRepositoryImpl)
	userRepositoryImpl := repository.NewUserRepositoryImpl()
	userSeederImpl := dataseeder.NewUserSeederImpl(userRepositoryImpl, roleRepositoryImpl)
	pegawaiRepositoryImpl := repository.NewPegawaiRepositoryImpl()
	pegawaiSeederImpl := dataseeder.NewPegawaiSeederImpl(db, pegawaiRepositoryImpl)
	seederImpl := dataseeder.NewSeederImpl(db, roleSeederImpl, userSeederImpl, pegawaiSeederImpl)
	return seederImpl
}

// injector.go:

var rencanaKinerjaSet = wire.NewSet(repository.NewRencanaKinerjaRepositoryImpl, wire.Bind(new(repository.RencanaKinerjaRepository), new(*repository.RencanaKinerjaRepositoryImpl)), service.NewRencanaKinerjaServiceImpl, wire.Bind(new(service.RencanaKinerjaService), new(*service.RencanaKinerjaServiceImpl)), controller.NewRencanaKinerjaControllerImpl, wire.Bind(new(controller.RencanaKinerjaController), new(*controller.RencanaKinerjaControllerImpl)))

var rencanaAksiSet = wire.NewSet(repository.NewRencanaAksiRepositoryImpl, wire.Bind(new(repository.RencanaAksiRepository), new(*repository.RencanaAksiRepositoryImpl)), service.NewRencanaAksiServiceImpl, wire.Bind(new(service.RencanaAksiService), new(*service.RencanaAksiServiceImpl)), controller.NewRencanaAksiControllerImpl, wire.Bind(new(controller.RencanaAksiController), new(*controller.RencanaAksiControllerImpl)))

var pelaksanaanRencanaAksiSet = wire.NewSet(repository.NewPelaksanaanRencanaAksiRepositoryImpl, wire.Bind(new(repository.PelaksanaanRencanaAksiRepository), new(*repository.PelaksanaanRencanaAksiRepositoryImpl)), service.NewPelaksanaanRencanaAksiServiceImpl, wire.Bind(new(service.PelaksanaanRencanaAksiService), new(*service.PelaksanaanRencanaAksiServiceImpl)), controller.NewPelaksanaanRencanaAksiControllerImpl, wire.Bind(new(controller.PelaksanaanRencanaAksiController), new(*controller.PelaksanaanRencanaAksiControllerImpl)))

var usulanMusrebangSet = wire.NewSet(repository.NewUsulanMusrebangRepositoryImpl, wire.Bind(new(repository.UsulanMusrebangRepository), new(*repository.UsulanMusrebangRepositoryImpl)), service.NewUsulanMusrebangServiceImpl, wire.Bind(new(service.UsulanMusrebangService), new(*service.UsulanMusrebangServiceImpl)), controller.NewUsulanMusrebangControllerImpl, wire.Bind(new(controller.UsulanMusrebangController), new(*controller.UsulanMusrebangControllerImpl)))

var usulanMandatoriSet = wire.NewSet(repository.NewUsulanMandatoriRepositoryImpl, wire.Bind(new(repository.UsulanMandatoriRepository), new(*repository.UsulanMandatoriRepositoryImpl)), service.NewUsulanMandatoriServiceImpl, wire.Bind(new(service.UsulanMandatoriService), new(*service.UsulanMandatoriServiceImpl)), controller.NewUsulanMandatoriControllerImpl, wire.Bind(new(controller.UsulanMandatoriController), new(*controller.UsulanMandatoriControllerImpl)))

var usulanPokokPikiranSet = wire.NewSet(repository.NewUsulanPokokPikiranRepositoryImpl, wire.Bind(new(repository.UsulanPokokPikiranRepository), new(*repository.UsulanPokokPikiranRepositoryImpl)), service.NewUsulanPokokPikiranServiceImpl, wire.Bind(new(service.UsulanPokokPikiranService), new(*service.UsulanPokokPikiranServiceImpl)), controller.NewUsulanPokokPikiranControllerImpl, wire.Bind(new(controller.UsulanPokokPikiranController), new(*controller.UsulanPokokPikiranControllerImpl)))

var usulanInisiatifSet = wire.NewSet(repository.NewUsulanInisiatifRepositoryImpl, wire.Bind(new(repository.UsulanInisiatifRepository), new(*repository.UsulanInisiatifRepositoryImpl)), service.NewUsulanInisiatifServiceImpl, wire.Bind(new(service.UsulanInisiatifService), new(*service.UsulanInisiatifServiceImpl)), controller.NewUsulanInisiatifControllerImpl, wire.Bind(new(controller.UsulanInisiatifController), new(*controller.UsulanInisiatifControllerImpl)))

var usulanTerpilihSet = wire.NewSet(repository.NewUsulanTerpilihRepositoryImpl, wire.Bind(new(repository.UsulanTerpilihRepository), new(*repository.UsulanTerpilihRepositoryImpl)), service.NewUsulanTerpilihServiceImpl, wire.Bind(new(service.UsulanTerpilihService), new(*service.UsulanTerpilihServiceImpl)), controller.NewUsulanTerpilihControllerImpl, wire.Bind(new(controller.UsulanTerpilihController), new(*controller.UsulanTerpilihControllerImpl)))

var gambaranUmumSet = wire.NewSet(repository.NewGambaranUmumRepositoryImpl, wire.Bind(new(repository.GambaranUmumRepository), new(*repository.GambaranUmumRepositoryImpl)), service.NewGambaranUmumServiceImpl, wire.Bind(new(service.GambaranUmumService), new(*service.GambaranUmumServiceImpl)), controller.NewGambaranUmumControllerImpl, wire.Bind(new(controller.GambaranUmumController), new(*controller.GambaranUmumControllerImpl)))

var dasarHukumSet = wire.NewSet(repository.NewDasarHukumRepositoryImpl, wire.Bind(new(repository.DasarHukumRepository), new(*repository.DasarHukumRepositoryImpl)), service.NewDasarHukumServiceImpl, wire.Bind(new(service.DasarHukumService), new(*service.DasarHukumServiceImpl)), controller.NewDasarHukumControllerImpl, wire.Bind(new(controller.DasarHukumController), new(*controller.DasarHukumControllerImpl)))

var inovasiSet = wire.NewSet(repository.NewInovasiRepositoryImpl, wire.Bind(new(repository.InovasiRepository), new(*repository.InovasiRepositoryImpl)), service.NewInovasiServiceImpl, wire.Bind(new(service.InovasiService), new(*service.InovasiServiceImpl)), controller.NewInovasiControllerImpl, wire.Bind(new(controller.InovasiController), new(*controller.InovasiControllerImpl)))

var subKegiatanSet = wire.NewSet(repository.NewSubKegiatanRepositoryImpl, wire.Bind(new(repository.SubKegiatanRepository), new(*repository.SubKegiatanRepositoryImpl)), service.NewSubKegiatanServiceImpl, wire.Bind(new(service.SubKegiatanService), new(*service.SubKegiatanServiceImpl)), controller.NewSubKegiatanControllerImpl, wire.Bind(new(controller.SubKegiatanController), new(*controller.SubKegiatanControllerImpl)))

var subKegiatanTerpilihSet = wire.NewSet(repository.NewSubKegiatanTerpilihRepositoryImpl, wire.Bind(new(repository.SubKegiatanTerpilihRepository), new(*repository.SubKegiatanTerpilihRepositoryImpl)), service.NewSubKegiatanTerpilihServiceImpl, wire.Bind(new(service.SubKegiatanTerpilihService), new(*service.SubKegiatanTerpilihServiceImpl)), controller.NewSubKegiatanTerpilihControllerImpl, wire.Bind(new(controller.SubKegiatanTerpilihController), new(*controller.SubKegiatanTerpilihControllerImpl)))

var pohonKinerjaOpdSet = wire.NewSet(repository.NewPohonKinerjaRepositoryImpl, wire.Bind(new(repository.PohonKinerjaRepository), new(*repository.PohonKinerjaRepositoryImpl)), service.NewPohonKinerjaOpdServiceImpl, wire.Bind(new(service.PohonKinerjaOpdService), new(*service.PohonKinerjaOpdServiceImpl)), controller.NewPohonKinerjaOpdControllerImpl, wire.Bind(new(controller.PohonKinerjaOpdController), new(*controller.PohonKinerjaOpdControllerImpl)))

var pegawaiSet = wire.NewSet(repository.NewPegawaiRepositoryImpl, wire.Bind(new(repository.PegawaiRepository), new(*repository.PegawaiRepositoryImpl)), service.NewPegawaiServiceImpl, wire.Bind(new(service.PegawaiService), new(*service.PegawaiServiceImpl)), controller.NewPegawaiControllerImpl, wire.Bind(new(controller.PegawaiController), new(*controller.PegawaiControllerImpl)))

var lembagaSet = wire.NewSet(repository.NewLembagaRepositoryImpl, wire.Bind(new(repository.LembagaRepository), new(*repository.LembagaRepositoryImpl)), service.NewLembagaServiceImpl, wire.Bind(new(service.LembagaService), new(*service.LembagaServiceImpl)), controller.NewLembagaControllerImpl, wire.Bind(new(controller.LembagaController), new(*controller.LembagaControllerImpl)))

var jabatanSet = wire.NewSet(repository.NewJabatanRepositoryImpl, wire.Bind(new(repository.JabatanRepository), new(*repository.JabatanRepositoryImpl)), service.NewJabatanServiceImpl, wire.Bind(new(service.JabatanService), new(*service.JabatanServiceImpl)), controller.NewJabatanControllerImpl, wire.Bind(new(controller.JabatanController), new(*controller.JabatanControllerImpl)))

var pohonKinerjaAdminSet = wire.NewSet(service.NewPohonKinerjaAdminServiceImpl, wire.Bind(new(service.PohonKinerjaAdminService), new(*service.PohonKinerjaAdminServiceImpl)), controller.NewPohonKinerjaAdminControllerImpl, wire.Bind(new(controller.PohonKinerjaAdminController), new(*controller.PohonKinerjaAdminControllerImpl)))

var opdSet = wire.NewSet(repository.NewOpdRepositoryImpl, wire.Bind(new(repository.OpdRepository), new(*repository.OpdRepositoryImpl)), service.NewOpdServiceImpl, wire.Bind(new(service.OpdService), new(*service.OpdServiceImpl)), controller.NewOpdControllerImpl, wire.Bind(new(controller.OpdController), new(*controller.OpdControllerImpl)))

var programSet = wire.NewSet(repository.NewProgramRepositoryImpl, wire.Bind(new(repository.ProgramRepository), new(*repository.ProgramRepositoryImpl)), service.NewProgramServiceImpl, wire.Bind(new(service.ProgramService), new(*service.ProgramServiceImpl)), controller.NewProgramControllerImpl, wire.Bind(new(controller.ProgramController), new(*controller.ProgramControllerImpl)))

var urusanSet = wire.NewSet(repository.NewUrusanRepositoryImpl, wire.Bind(new(repository.UrusanRepository), new(*repository.UrusanRepositoryImpl)), service.NewUrusanServiceImpl, wire.Bind(new(service.UrusanService), new(*service.UrusanServiceImpl)), controller.NewUrusanControllerImpl, wire.Bind(new(controller.UrusanController), new(*controller.UrusanControllerImpl)))

var bidangUrusanSet = wire.NewSet(repository.NewBidangUrusanRepositoryImpl, wire.Bind(new(repository.BidangUrusanRepository), new(*repository.BidangUrusanRepositoryImpl)), service.NewBidangUrusanServiceImpl, wire.Bind(new(service.BidangUrusanService), new(*service.BidangUrusanServiceImpl)), controller.NewBidangUrusanControllerImpl, wire.Bind(new(controller.BidangUrusanController), new(*controller.BidangUrusanControllerImpl)))

var kegiatanSet = wire.NewSet(repository.NewKegiatanRepositoryImpl, wire.Bind(new(repository.KegiatanRepository), new(*repository.KegiatanRepositoryImpl)), service.NewKegiatanServiceImpl, wire.Bind(new(service.KegiatanService), new(*service.KegiatanServiceImpl)), controller.NewKegiatanControllerImpl, wire.Bind(new(controller.KegiatanController), new(*controller.KegiatanControllerImpl)))

var roleSet = wire.NewSet(repository.NewRoleRepositoryImpl, wire.Bind(new(repository.RoleRepository), new(*repository.RoleRepositoryImpl)), service.NewRoleServiceImpl, wire.Bind(new(service.RoleService), new(*service.RoleServiceImpl)), controller.NewRoleControllerImpl, wire.Bind(new(controller.RoleController), new(*controller.RoleControllerImpl)))

var userSet = wire.NewSet(repository.NewUserRepositoryImpl, wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)), service.NewUserServiceImpl, wire.Bind(new(service.UserService), new(*service.UserServiceImpl)), controller.NewUserControllerImpl, wire.Bind(new(controller.UserController), new(*controller.UserControllerImpl)))

var seederProviderSet = wire.NewSet(dataseeder.NewSeederImpl, wire.Bind(new(dataseeder.Seeder), new(*dataseeder.SeederImpl)), dataseeder.NewRoleSeederImpl, wire.Bind(new(dataseeder.RoleSeeder), new(*dataseeder.RoleSeederImpl)), dataseeder.NewUserSeederImpl, wire.Bind(new(dataseeder.UserSeeder), new(*dataseeder.UserSeederImpl)), dataseeder.NewPegawaiSeederImpl, wire.Bind(new(dataseeder.PegawaiSeeder), new(*dataseeder.PegawaiSeederImpl)))

var tujuanOpdSet = wire.NewSet(repository.NewTujuanOpdRepositoryImpl, wire.Bind(new(repository.TujuanOpdRepository), new(*repository.TujuanOpdRepositoryImpl)), service.NewTujuanOpdServiceImpl, wire.Bind(new(service.TujuanOpdService), new(*service.TujuanOpdServiceImpl)), controller.NewTujuanOpdControllerImpl, wire.Bind(new(controller.TujuanOpdController), new(*controller.TujuanOpdControllerImpl)))
