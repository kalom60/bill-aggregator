package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/models"
)

type Service interface {
	InsertOne(userID, providerID, accountIdentifier, encryptedCredential string) error
	FetchLinkedAccountsByUserID(userID string) ([]models.Account, error)
	DeleteLinkedAccountByID(accountID string) error
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

func (s *service) InsertOne(userID, providerID, accountIdentifier, encryptedCredential string) error {
	query := `INSERT INTO linked_accounts (user_id, provider_id, account_identifier, encrypted_credential, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	var id string
	err := s.db.QueryRow(query, userID, providerID, accountIdentifier, encryptedCredential).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FetchLinkedAccountsByUserID(userID string) ([]models.Account, error) {
	query := `SELECT id, user_id, provider_id, account_identifier, created_at FROM linked_accounts WHERE user_id = $1`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.ProviderID, &acc.AccountIdentifier, &acc.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *service) DeleteLinkedAccountByID(accountID string) error {
	query := `DELETE FROM linked_accounts WHERE id = $1`

	result, err := s.db.Exec(query, accountID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no linked account found for account_id: %s", accountID)
	}

	return nil
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}
