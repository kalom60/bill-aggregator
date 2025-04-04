package database

import (
	"bill-aggregation-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Get(ctx context.Context, userID string) ([]models.Bill, error)
	Set(ctx context.Context, userID string, bills []models.Bill) error
}

type service struct {
	db  *redis.Client
	ttl time.Duration
}

var (
	address  = os.Getenv("REDIS_ADDRESS")
	password = os.Getenv("REDIS_PASSWORD")
	ctx      = context.Background()
)

func New() Service {
	fullAddress := fmt.Sprintf("%s:%s", address, "6379")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: password,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	s := &service{
		db:  rdb,
		ttl: 10 * time.Minute,
	}

	return s
}

func (s *service) Get(ctx context.Context, userID string) ([]models.Bill, error) {
	val, err := s.db.Get(ctx, userKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	var bills []models.Bill
	if err := json.Unmarshal([]byte(val), &bills); err != nil {
		return nil, err
	}

	return bills, nil
}

func (s *service) Set(ctx context.Context, userID string, bills []models.Bill) error {
	data, err := json.Marshal(bills)
	if err != nil {
		return err
	}

	return s.db.Set(ctx, userKey(userID), data, s.ttl).Err()
}

func userKey(userID string) string {
	return "bills:" + userID
}
