package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/models"
)

type Service interface {
	InsertOne(name, api_url, authentication_typem, api_key string) error
	FetchProviderByID(id string) (*models.Provider, error)
	FetchProviders() ([]models.Provider, error)
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

func (s *service) InsertOne(name, api_url, authentication_type, api_key string) error {
	query := `INSERT INTO providers (name, api_url, authentication_type, api_key, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	var id string
	err := s.db.QueryRow(query, name, api_url, authentication_type, api_key).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

var ErrProviderNotFound = errors.New("provider not found")

func (s *service) FetchProviderByID(id string) (*models.Provider, error) {
	query := `
        SELECT id, name, api_url, authentication_type, api_key, created_at, updated_at
        FROM providers
        WHERE id = $1
    `

	row := s.db.QueryRow(query, id)

	var provider models.Provider
	err := row.Scan(
		&provider.ID,
		&provider.Name,
		&provider.API_URL,
		&provider.Authentication_Type,
		&provider.API_key,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("provider not found: %w", ErrProviderNotFound)
	case err != nil:
		return nil, fmt.Errorf("error fetching provider: %w", err)
	default:
		return &provider, nil
	}
}

func (s *service) FetchProviders() ([]models.Provider, error) {
	query := `SELECT * FROM providers`
	rows, err := s.db.Query(query)
	if err != nil {
		return []models.Provider{}, err
	}
	defer rows.Close()

	var providers []models.Provider
	for rows.Next() {
		var p models.Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.API_URL, &p.Authentication_Type, &p.API_key, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return []models.Provider{}, err
		}
		providers = append(providers, p)
	}

	return providers, nil
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}
