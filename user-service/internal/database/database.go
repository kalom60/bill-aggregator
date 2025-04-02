package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kalom60/bill-aggregator/user-service/internal/models"
)

type Service interface {
	InsertOne(email, password, firstname, lastName string) error
	GetUserByEmail(email string) (*models.User, error)
	Close() error
}

type service struct {
	db *sql.DB
}

var dbInstance *service

func New() Service {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSLMODE"),
		os.Getenv("DB_SCHEMA"),
	)

	if dbInstance != nil {
		return dbInstance
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}

	err = RunMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	return dbInstance
}

func (s *service) InsertOne(email, password, firstName, lastName string) error {
	query := `INSERT INTO users (email, password, first_name, last_name, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	var id string
	err := s.db.QueryRow(query, email, password, firstName, lastName).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, first_name, last_name, created_at FROM users WHERE email = $1`
	row := s.db.QueryRow(query, email)

	var user models.User

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}
