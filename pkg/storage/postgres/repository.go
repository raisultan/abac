package postgres

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/raisultan/abac/pkg/list"
	"github.com/raisultan/abac/pkg/login"
	"github.com/raisultan/abac/pkg/register"
	"github.com/raisultan/abac/pkg/retrieve"
	"github.com/raisultan/abac/pkg/update"
	"golang.org/x/crypto/bcrypt"
)

var dbUrl = os.Getenv("POSTGRES_URL")

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	var err error
	s := Storage{}

	s.db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Storage) GetAllUsers(limit, offset int) ([]list.UserRetrieveResponse, error) {
	rows, err := s.db.Query(
		"SELECT id, email, firstName, lastName FROM users LIMIT $1 OFFSET $2",
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []list.UserRetrieveResponse{}

	for rows.Next() {
		var u list.UserRetrieveResponse
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *Storage) CreateUser(ru register.UserRegisterRequest) (register.UserRegisterResponse, error) {
	hashingCost := 8
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ru.Password), hashingCost)
	if err != nil {
		return register.UserRegisterResponse{}, err
	}

	hashedPasswordStr := string(hashedPassword)
	err = s.db.QueryRow(
		"INSERT INTO users(email, password, firstName, lastName) VALUES($1, $2, $3, $4) RETURNING id",
		ru.Email,
		hashedPasswordStr,
		ru.FirstName,
		ru.LastName,
	).Scan(&ru.ID)

	if err != nil {
		return register.UserRegisterResponse{}, err
	}

	resp := register.UserRegisterResponse{
		ID:        ru.ID,
		Email:     ru.Email,
		FirstName: ru.FirstName,
		LastName:  ru.LastName,
	}
	return resp, nil
}

func (s *Storage) CheckUserExists(ru register.UserRegisterRequest) (bool, error) {
	err := s.db.QueryRow(
		"SELECT email FROM users WHERE email=$1",
		ru.Email,
	).Scan(&ru.Email)

	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (s *Storage) GetUserByID(uID int) (retrieve.UserRetrieveResponse, error) {
	u := retrieve.UserRetrieveResponse{}

	err := s.db.QueryRow(
		"SELECT id, email, firstName, lastName, isAdmin, isApproved FROM users WHERE id=$1",
		uID,
	).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.IsAdmin, &u.IsApproved)

	if err != nil {
		return retrieve.UserRetrieveResponse{}, err
	}

	return u, nil
}

func (s *Storage) GetUserByEmail(r login.UserLoginRequest) (login.UserLoginRequest, error) {
	u := login.UserLoginRequest{}
	err := s.db.QueryRow(
		"SELECT email, password FROM users WHERE email=$1",
		r.Email,
	).Scan(&u.Email, &u.Password)

	if err != nil {
		return login.UserLoginRequest{}, err
	}

	return u, nil
}

func (s *Storage) UpdateUser(r update.UserUpdateRequest) (update.UserRetrieveResponse, error) {
	_, err := s.db.Exec(
		"UPDATE users SET firstName=$1, lastName=$2 WHERE id=$3",
		r.FirstName,
		r.LastName,
		r.ID,
	)

	if err != nil {
		return update.UserRetrieveResponse{}, err
	}

	u := update.UserRetrieveResponse{}
	err = s.db.QueryRow(
		"SELECT email, firstName, lastName, isAdmin, isApproved FROM users WHERE id=$1",
		r.ID,
	).Scan(&u.Email, &u.FirstName, &u.LastName, &u.IsAdmin, &u.IsApproved)

	if err != nil {
		return update.UserRetrieveResponse{}, err
	}

	return u, nil
}

func (s *Storage) DeleteUser(uID int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id=$1", uID)
	return err
}
