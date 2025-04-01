package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/models"
)

type Service interface {
	InsertOne(name, api_url, authentication_typem, api_key string) error
	FetchProviders() (*[]models.Provider, error)
	Close() error
}

type service struct {
	db *sql.DB
}

var dbInstance *service

func New() Service {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSLMODE"),
	)

	RunMigrations(dsn)

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
	return dbInstance
}

func (s *service) InsertOne(name, api_url, authentication_type, api_key string) error {
	query := `INSERT INTO providers (name, api_url, authentication_type, api_key, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	var id string
	err := s.db.QueryRow(query, name, api_url, authentication_type, api_key).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FetchProviders() (*[]models.Provider, error) {
	query := `SELECT * FROM providers`
	rows, err := s.db.Query(query)
	if err != nil {
		return &[]models.Provider{}, err
	}
	defer rows.Close()

	var providers []models.Provider
	for rows.Next() {
		var p models.Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.API_URL, &p.Authentication_Type, &p.CreatedAt); err != nil {
			return &[]models.Provider{}, err
		}
		providers = append(providers, p)
	}

	return &providers, nil
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}
