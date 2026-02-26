package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/user"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"sort"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository    repository.UserRepository
	RoleRepository    repository.RoleRepository
	PegawaiRepository repository.PegawaiRepository
	OpdRepository     repository.OpdRepository
	DB                *sql.DB
}

func NewUserServiceImpl(userRepository repository.UserRepository, roleRepository repository.RoleRepository, pegawaiRepository repository.PegawaiRepository, opdRepository repository.OpdRepository, db *sql.DB) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository:    userRepository,
		RoleRepository:    roleRepository,
		PegawaiRepository: pegawaiRepository,
		OpdRepository:     opdRepository,
		DB:                db,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request user.UserCreateRequest) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input dasar
	if request.Nip == "" {
		return user.UserResponse{}, errors.New("nip harus diisi")
	}
	if request.Password == "" {
		return user.UserResponse{}, errors.New("password harus diisi")
	}
	if len(request.Role) == 0 {
		return user.UserResponse{}, errors.New("role harus diisi")
	}

	// Validasi NIP dengan data pegawai
	_, err = service.PegawaiRepository.FindByNip(ctx, tx, request.Nip)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.UserResponse{}, errors.New("nip tidak terdaftar di data pegawai")
		}
		return user.UserResponse{}, err
	}

	// Siapkan slice untuk menyimpan roles
	var roles []domain.Roles

	// Validasi dan ambil semua role yang dipilih
	for _, roleRequest := range request.Role {
		role, err := service.RoleRepository.FindById(ctx, tx, roleRequest.RoleId)
		if err != nil {
			if err == sql.ErrNoRows {
				return user.UserResponse{}, errors.New("role tidak ditemukan")
			}
			return user.UserResponse{}, err
		}
		roles = append(roles, role)
	}

	userDomain := domain.Users{
		Nip:      request.Nip,
		Email:    helper.EmptyStringIfNull(request.Email),
		Password: request.Password,
		IsActive: request.IsActive,
		Role:     roles,
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDomain.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.UserResponse{}, err
	}
	userDomain.Password = string(hashedPassword)

	createdUser, err := service.UserRepository.Create(ctx, tx, userDomain)
	if err != nil {
		return user.UserResponse{}, err
	}

	// Konversi role ke response
	var roleResponses []user.RoleResponse
	for _, role := range createdUser.Role {
		roleResponses = append(roleResponses, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	response := user.UserResponse{
		Id:       createdUser.Id,
		Nip:      createdUser.Nip,
		Email:    createdUser.Email,
		IsActive: createdUser.IsActive,
		Role:     roleResponses,
	}

	return response, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request user.UserUpdateRequest) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi user exists
	existingUser, err := service.UserRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return user.UserResponse{}, err
	}
	if existingUser.Id == 0 {
		return user.UserResponse{}, errors.New("user tidak ditemukan")
	}

	// Validasi input dasar
	if request.Nip == "" {
		return user.UserResponse{}, errors.New("nip harus diisi")
	}
	if request.Email == "" {
		return user.UserResponse{}, errors.New("email harus diisi")
	}
	if len(request.Role) == 0 {
		return user.UserResponse{}, errors.New("role harus diisi")
	}

	// Siapkan slice untuk menyimpan roles
	var roles []domain.Roles

	// Validasi dan ambil semua role yang dipilih
	for _, roleRequest := range request.Role {
		role, err := service.RoleRepository.FindById(ctx, tx, roleRequest.RoleId)
		if err != nil {
			if err == sql.ErrNoRows {
				return user.UserResponse{}, errors.New("role tidak ditemukan")
			}
			return user.UserResponse{}, err
		}
		roles = append(roles, role)
	}

	userDomain := domain.Users{
		Id:       existingUser.Id,
		Nip:      request.Nip,
		Email:    request.Email,
		IsActive: request.IsActive,
		Role:     roles,
	}

	// Handle password update
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return user.UserResponse{}, err
		}
		userDomain.Password = string(hashedPassword)
	} else {
		userDomain.Password = existingUser.Password
	}

	updatedUser, err := service.UserRepository.Update(ctx, tx, userDomain)
	if err != nil {
		return user.UserResponse{}, err
	}

	// Konversi role ke response
	var roleResponses []user.RoleResponse
	for _, role := range updatedUser.Role {
		roleResponses = append(roleResponses, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	response := user.UserResponse{
		Id:       updatedUser.Id,
		Nip:      updatedUser.Nip,
		Email:    updatedUser.Email,
		IsActive: updatedUser.IsActive,
		Role:     roleResponses,
	}

	return response, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	existingUser, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}
	if existingUser.Id == 0 {
		return errors.New("user tidak ditemukan")
	}

	err = service.UserRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) FindAll(ctx context.Context, kodeOpd string) ([]user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}

	// Ambil seluruh NIP user secara unik untuk batch query pegawai
	nipSet := make(map[string]struct{})
	var nips []string
	for _, u := range users {
		if u.Nip == "" {
			continue
		}
		if _, exists := nipSet[u.Nip]; !exists {
			nipSet[u.Nip] = struct{}{}
			nips = append(nips, u.Nip)
		}
	}

	pegawaiByNip, err := service.PegawaiRepository.FindPegawaiByNipsBatch(ctx, tx, nips)
	if err != nil {
		return nil, err
	}

	var userResponses []user.UserResponse
	for _, u := range users {
		var roles []user.RoleResponse
		for _, role := range u.Role {
			roles = append(roles, user.RoleResponse{
				Id:   role.Id,
				Role: role.Role,
			})
		}

		var pegawaiId, namaPegawai, namaJabatan, idJabatan string
		if pegawaiDomain, ok := pegawaiByNip[u.Nip]; ok && pegawaiDomain != nil {
			pegawaiId = pegawaiDomain.Id
			namaPegawai = pegawaiDomain.NamaPegawai
			namaJabatan = pegawaiDomain.NamaJabatan
			idJabatan = pegawaiDomain.IdJabatan
		}

		userResponse := user.UserResponse{
			Id:          u.Id,
			PegawaiId:   pegawaiId,
			Nip:         u.Nip,
			Email:       u.Email,
			NamaPegawai: namaPegawai,
			IdJabatan:   idJabatan,
			NamaJabatan: namaJabatan,
			IsActive:    u.IsActive,
			Role:        roles,
		}
		userResponses = append(userResponses, userResponse)
	}

	sort.Slice(userResponses, func(i, j int) bool {
		return userResponses[i].NamaPegawai < userResponses[j].NamaPegawai

	})

	return userResponses, nil
}

func (service *UserServiceImpl) FindById(ctx context.Context, id int) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cari user berdasarkan ID
	userDomain, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return user.UserResponse{}, err
	}

	// Cek apakah user ditemukan
	if userDomain.Id == 0 {
		return user.UserResponse{}, errors.New("user tidak ditemukan")
	}

	// Konversi role domain ke role response
	var roles []user.RoleResponse
	for _, role := range userDomain.Role {
		roles = append(roles, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	pegawaiDomain, _ := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)

	// Convert ke response
	response := user.UserResponse{
		Id:          userDomain.Id,
		Nip:         userDomain.Nip,
		Email:       userDomain.Email,
		NamaPegawai: pegawaiDomain.NamaPegawai,
		IsActive:    userDomain.IsActive,
		Role:        roles,
	}

	return response, nil
}

// func (service *UserServiceImpl) Login(ctx context.Context, request user.UserLoginRequest) (user.UserLoginResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return user.UserLoginResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	if request.Username == "" {
// 		return user.UserLoginResponse{}, errors.New("email atau nip harus diisi")
// 	}
// 	if request.Password == "" {
// 		return user.UserLoginResponse{}, errors.New("password harus diisi")
// 	}

// 	userDomain, err := service.UserRepository.FindByEmailOrNip(ctx, tx, request.Username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return user.UserLoginResponse{}, errors.New("username atau password salah")
// 		}
// 		return user.UserLoginResponse{}, err
// 	}

// 	pegawaiDomain, err := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)
// 	if err != nil {
// 		return user.UserLoginResponse{}, err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(userDomain.Password), []byte(request.Password))
// 	if err != nil {
// 		return user.UserLoginResponse{}, errors.New("username atau password salah")
// 	}

// 	if !userDomain.IsActive {
// 		return user.UserLoginResponse{}, errors.New("akun tidak aktif")
// 	}

// 	roleNames := make([]string, 0, len(userDomain.Role))
// 	for _, role := range userDomain.Role {
// 		roleNames = append(roleNames, role.Role)
// 	}

// 	token := helper.CreateNewJWT(
// 		userDomain.Id,
// 		pegawaiDomain.Id,
// 		userDomain.Email,
// 		userDomain.Nip,
// 		pegawaiDomain.KodeOpd,
// 		roleNames,
// 	)

// 	response := user.UserLoginResponse{
// 		Token: token,
// 	}

// 	return response, nil
// }

func (service *UserServiceImpl) Login(ctx context.Context, request user.UserLoginRequest) (user.UserLoginResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserLoginResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input
	if request.Username == "" {
		return user.UserLoginResponse{}, errors.New("nip harus diisi")
	}
	if request.Password == "" {
		return user.UserLoginResponse{}, errors.New("password harus diisi")
	}

	// Validasi format NIP
	// if !helper.IsValidNIP(request.Username) {
	// 	return user.UserLoginResponse{}, errors.New("format nip tidak valid")
	// }

	// Cari user berdasarkan NIP saja
	userDomain, err := service.UserRepository.FindByEmailOrNip(ctx, tx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.UserLoginResponse{}, errors.New("nip atau password salah")
		}
		return user.UserLoginResponse{}, err
	}

	// Pastikan username yang digunakan adalah NIP
	if userDomain.Nip != request.Username {
		return user.UserLoginResponse{}, errors.New("silakan login menggunakan NIP")
	}

	pegawaiDomain, err := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)
	if err != nil {
		return user.UserLoginResponse{}, err
	}

	opdDomain, err := service.OpdRepository.FindByKodeOpd(ctx, tx, pegawaiDomain.KodeOpd)
	if err != nil {
		return user.UserLoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDomain.Password), []byte(request.Password))
	if err != nil {
		return user.UserLoginResponse{}, errors.New("nip atau password salah")
	}

	if !userDomain.IsActive {
		return user.UserLoginResponse{}, errors.New("akun tidak aktif")
	}

	roleNames := make([]string, 0, len(userDomain.Role))
	for _, role := range userDomain.Role {
		roleNames = append(roleNames, role.Role)
	}

	token := helper.CreateNewJWT(
		userDomain.Id,
		pegawaiDomain.Id,
		userDomain.Email,
		userDomain.Nip,
		pegawaiDomain.KodeOpd,
		opdDomain.NamaOpd,
		pegawaiDomain.NamaPegawai,
		roleNames,
	)

	response := user.UserLoginResponse{
		Token: token,
	}

	return response, nil
}

func (service *UserServiceImpl) FindByKodeOpdAndRole(ctx context.Context, kodeOpd string, roleName string) ([]user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input
	if kodeOpd == "" {
		return nil, errors.New("kode opd harus diisi")
	}
	if roleName == "" {
		return nil, errors.New("role harus diisi")
	}

	users, err := service.UserRepository.FindByKodeOpdAndRole(ctx, tx, kodeOpd, roleName)
	if err != nil {
		return nil, err
	}

	var userResponses []user.UserResponse
	for _, u := range users {
		var roles []user.RoleResponse
		for _, role := range u.Role {
			roles = append(roles, user.RoleResponse{
				Id:   role.Id,
				Role: role.Role,
			})
		}

		userResponse := user.UserResponse{
			Id:          u.Id,
			Nip:         u.Nip,
			IsActive:    u.IsActive,
			PegawaiId:   u.PegawaiId,
			NamaPegawai: u.NamaPegawai,
			Role:        roles,
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

func (service *UserServiceImpl) FindByNip(ctx context.Context, nip string) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	userDomain, err := service.UserRepository.FindByNip(ctx, tx, nip)
	if err != nil {
		return user.UserResponse{}, err
	}

	var roles []user.RoleResponse
	for _, role := range userDomain.Role {
		roles = append(roles, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	pegawaiDomain, err := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)
	if err != nil {
		return user.UserResponse{}, err
	}

	userResponse := user.UserResponse{
		Nip:         userDomain.Nip,
		NamaPegawai: pegawaiDomain.NamaPegawai,
		IsActive:    userDomain.IsActive,
		Role:        roles,
	}

	return userResponse, nil
}

func (service *UserServiceImpl) CekAdminOpd(ctx context.Context) ([]user.CekAdminOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua OPD
	allOpd, err := service.OpdRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Ambil semua user dengan role admin_opd
	adminUsers, err := service.UserRepository.CekAdminOpd(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Buat map untuk grouping user berdasarkan kode_opd
	adminByOpd := make(map[string][]user.AdminOpdUserDetail)
	for _, u := range adminUsers {
		adminDetail := user.AdminOpdUserDetail{
			UserId:      u.Id,
			Nip:         u.Nip,
			NamaPegawai: u.NamaPegawai,
			Email:       u.Email,
			IsActive:    u.IsActive,
		}
		adminByOpd[u.KodeOpd] = append(adminByOpd[u.KodeOpd], adminDetail)
	}

	// Build response: semua OPD dengan admin users (atau array kosong jika tidak ada)
	var response []user.CekAdminOpdResponse
	for _, opd := range allOpd {
		opdResponse := user.CekAdminOpdResponse{
			KodeOpd:    opd.KodeOpd,
			NamaOpd:    opd.NamaOpd,
			AdminUsers: []user.AdminOpdUserDetail{}, // inisialisasi dengan array kosong
		}

		// Jika ada admin di OPD ini, masukkan datanya
		if admins, exists := adminByOpd[opd.KodeOpd]; exists {
			opdResponse.AdminUsers = admins
		}

		response = append(response, opdResponse)
	}

	return response, nil
}
