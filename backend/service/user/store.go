package user

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/utsavll0/ecom/types"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM public.users WHERE email = $1", email)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) CreateUser(u types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstname, lastname, email, password) VALUES ($1,$2,$3,$4)", u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM public.users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}
