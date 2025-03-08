package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"log"
	"time"
)

type UserRepositoryImpl struct {
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, users domain.Users) (domain.Users, error) {
	script := "INSERT INTO tb_users(nip, email, password, is_active) VALUES (?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, users.Nip, users.Email, users.Password, users.IsActive)
	if err != nil {
		return users, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return users, err
	}
	users.Id = int(id)

	scriptRole := "INSERT INTO tb_user_role(user_id, role_id) VALUES (?, ?)"
	for _, role := range users.Role {
		_, err = tx.ExecContext(ctx, scriptRole, users.Id, role.Id)
		if err != nil {
			return users, err
		}
	}

	return users, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, users domain.Users) (domain.Users, error) {
	script := "UPDATE tb_users SET nip = ?, email = ?, password = ?, is_active = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, users.Nip, users.Email, users.Password, users.IsActive, users.Id)
	if err != nil {
		return users, err
	}

	scriptDeleteRoles := "DELETE FROM tb_user_role WHERE user_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteRoles, users.Id)
	if err != nil {
		return users, err
	}

	scriptRole := "INSERT INTO tb_user_role(user_id, role_id) VALUES (?, ?)"
	for _, role := range users.Role {
		_, err = tx.ExecContext(ctx, scriptRole, users.Id, role.Id)
		if err != nil {
			return users, err
		}
	}

	return repository.FindById(ctx, tx, users.Id)
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.Users, error) {
	script := `
        SELECT DISTINCT u.id, u.nip, u.email, u.is_active, ur.role_id, r.role 
        FROM tb_users u
        LEFT JOIN tb_user_role ur ON u.id = ur.user_id
        LEFT JOIN tb_role r ON ur.role_id = r.id
        INNER JOIN tb_pegawai p ON u.nip = p.nip
        WHERE 1=1
    `
	var params []interface{}

	if kodeOpd != "" {
		script += " AND p.kode_opd = ?"
		params = append(params, kodeOpd)
	}

	script += " ORDER BY u.id, ur.role_id"

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return []domain.Users{}, err
	}
	defer rows.Close()

	var users []domain.Users
	userMap := make(map[int]*domain.Users)

	for rows.Next() {
		var userId int
		var nip, email string
		var isActive bool
		var roleId sql.NullInt64
		var roleName sql.NullString

		err := rows.Scan(
			&userId,
			&nip,
			&email,
			&isActive,
			&roleId,
			&roleName,
		)
		if err != nil {
			return []domain.Users{}, err
		}

		user, exists := userMap[userId]
		if !exists {
			user = &domain.Users{
				Id:       userId,
				Nip:      nip,
				Email:    email,
				IsActive: isActive,
				Role:     []domain.Roles{},
			}
			userMap[userId] = user
		}

		if roleId.Valid && roleName.Valid {
			user.Role = append(user.Role, domain.Roles{
				Id:   int(roleId.Int64),
				Role: roleName.String,
			})
		}
	}

	for _, user := range userMap {
		users = append(users, *user)
	}

	return users, nil
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Users, error) {
	script := `
		SELECT u.id, u.nip, u.email, u.password, u.is_active, ur.role_id, r.role 
		FROM tb_users u
			LEFT JOIN tb_user_role ur ON u.id = ur.user_id
			LEFT JOIN tb_role r ON ur.role_id = r.id
		WHERE u.id = ?
		ORDER BY ur.role_id
	`
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.Users{}, err
	}
	defer rows.Close()

	var user domain.Users
	first := true

	for rows.Next() {
		var roleId sql.NullInt64
		var roleName sql.NullString

		if first {
			err := rows.Scan(
				&user.Id,
				&user.Nip,
				&user.Email,
				&user.Password,
				&user.IsActive,
				&roleId,
				&roleName,
			)
			if err != nil {
				return domain.Users{}, err
			}
			first = false
		} else {
			var userId int
			var nip, email, password string
			var isActive bool
			err := rows.Scan(
				&userId,
				&nip,
				&email,
				&password,
				&isActive,
				&roleId,
				&roleName,
			)
			if err != nil {
				return domain.Users{}, err
			}
		}

		if roleId.Valid && roleName.Valid {
			user.Role = append(user.Role, domain.Roles{
				Id:   int(roleId.Int64),
				Role: roleName.String,
			})
		}
	}

	return user, nil
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	scriptRole := "DELETE FROM tb_user_role WHERE user_id = ?"
	_, err := tx.ExecContext(ctx, scriptRole, id)
	if err != nil {
		return err
	}

	scriptUser := "DELETE FROM tb_users WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptUser, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepositoryImpl) FindByEmailOrNip(ctx context.Context, tx *sql.Tx, username string) (domain.Users, error) {
	startTime := time.Now()
	log.Printf("Start finding user by email/NIP: %s", username)

	script := `
		SELECT u.id, u.nip, u.email, u.password, u.is_active, ur.role_id, r.role 
		FROM tb_users u
		LEFT JOIN tb_user_role ur ON u.id = ur.user_id
		LEFT JOIN tb_role r ON ur.role_id = r.id
		WHERE u.email = ? OR u.nip = ?
		ORDER BY ur.role_id
	`
	rows, err := tx.QueryContext(ctx, script, username, username)
	if err != nil {
		log.Printf("Error querying user by email/NIP: %v", err)
		return domain.Users{}, err
	}
	defer rows.Close()

	var user domain.Users
	first := true

	for rows.Next() {
		var roleId sql.NullInt64
		var roleName sql.NullString

		if first {
			err := rows.Scan(
				&user.Id,
				&user.Nip,
				&user.Email,
				&user.Password,
				&user.IsActive,
				&roleId,
				&roleName,
			)
			if err != nil {
				log.Printf("Error scanning first user row: %v", err)
				return domain.Users{}, err
			}
			first = false
		} else {
			var userId int
			var nip, email, password string
			var isActive bool
			err := rows.Scan(
				&userId,
				&nip,
				&email,
				&password,
				&isActive,
				&roleId,
				&roleName,
			)
			if err != nil {
				return domain.Users{}, err
			}
		}

		if roleId.Valid && roleName.Valid {
			user.Role = append(user.Role, domain.Roles{
				Id:   int(roleId.Int64),
				Role: roleName.String,
			})
		}
	}

	log.Printf("Successfully found user by email/NIP: %s, execution time: %v", username, time.Since(startTime))
	return user, nil
}
