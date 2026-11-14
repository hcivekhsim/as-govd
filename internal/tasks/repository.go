package tasks

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

type dbTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (u *UserRepo) getByLogin(ctx context.Context, name string) (bool, error) {
	var exists int
	err := u.pool.QueryRow(ctx, "SELECT 1 FROM users WHERE login=$1", name).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}

func (u *UserRepo) getUserByLogin(ctx context.Context, name string) (*Login, error) {
	var user Login
	err := u.pool.QueryRow(ctx, `SELECT id,login,password FROM users WHERE login=$1`, name).Scan(
		&user.UserID,
		&user.Login,
		&user.Password,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, dbError(err)
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) createUser(ctx context.Context, user Register, password string) (int64, error) {
	var id int64
	err := u.pool.QueryRow(ctx, `INSERT INTO users (user_type, name_organization, login, password, email, phone,last_name, first_name, patronymic, is_abroad,locality, street, house, apartment, note_address, postal_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`,
		user.UserType, user.NameOrganization, user.Login, password, user.Email, user.Phone,
		user.LastName, user.FirstName, user.Patronymic, user.IsAbroad,
		user.Locality, user.Street, user.House, user.Apartment, user.NoteAddress, user.PostalCode,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *UserRepo) createUserAuth(ctx context.Context, user LoginAuth) error {
	_, err := u.pool.Exec(ctx, `INSERT INTO users_auth (user_id, session_token_hash, csrf_token_hash, expires_at) VAlUES ($1, $2, $3, $4)`,
		user.UserID,
		user.SessionTokenHash,
		user.CSRFTokenHash,
		user.ExpiresAt,
	)
	if err != nil {
		return err
	}
	return err
}

func (u *UserRepo) deleteUserAuthBySessionTokenHash(ctx context.Context, sessionTokenHash string) error {
	_, err := u.pool.Exec(ctx,
		`DELETE FROM users_auth WHERE session_token_hash = $1`,
		sessionTokenHash,
	)
	if err != nil {
		return dbError(err)
	}
	return nil
}

func (u *UserRepo) getUserAuthBySessionTokenHash(ctx context.Context, sessionTokenHash string) (*LoginAuth, error) {
	var userAuth LoginAuth
	err := u.pool.QueryRow(ctx, `
		SELECT ua.user_id, u.login, u.role, u.preferred_language, ua.session_token_hash, ua.csrf_token_hash, ua.expires_at
		FROM users_auth ua
		JOIN users u ON u.id = ua.user_id
		WHERE ua.session_token_hash = $1
		ORDER BY ua.expires_at DESC
		LIMIT 1
	`, sessionTokenHash).Scan(
		&userAuth.UserID,
		&userAuth.Login,
		&userAuth.Role,
		&userAuth.PreferredLanguage,
		&userAuth.SessionTokenHash,
		&userAuth.CSRFTokenHash,
		&userAuth.ExpiresAt,
	)
	if err != nil {
		return nil, dbError(err)
	}

	return &userAuth, nil
}

// getPreferredLanguage возвращает сохранённый в профиле язык интерфейса.
func (u *UserRepo) getPreferredLanguage(ctx context.Context, userID int64) (Lang, error) {
	var lang string
	err := u.pool.QueryRow(ctx, `SELECT preferred_language FROM users WHERE id=$1`, userID).Scan(&lang)
	if err != nil {
		return LangRu, dbError(err)
	}
	result := Lang(lang)
	if !result.valid() {
		return LangRu, nil
	}
	return result, nil
}

// updatePreferredLanguage сохраняет выбранный язык интерфейса в профиле.
func (u *UserRepo) updatePreferredLanguage(ctx context.Context, userID int64, lang Lang) error {
	_, err := u.pool.Exec(ctx, `UPDATE users SET preferred_language=$2 WHERE id=$1`, userID, string(lang))
	if err != nil {
		return dbError(err)
	}
	return nil
}

func (u *UserRepo) getByAll(ctx context.Context) ([]Register, error) {
	var users []Register

	rows, err := u.pool.Query(ctx, "SELECT login FROM users") //...
	if err != nil {
		return nil, dbError(err)
	}

	for rows.Next() {
		var user Register
		if err := rows.Scan(&user.Login); err != nil {
			return nil, dbError(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, dbError(err)
	}

	return users, nil
}
