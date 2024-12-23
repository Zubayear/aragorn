package user

import (
	"context"
	"errors"
	"github.com/Zubayear/aragorn/internal/repository"
	"github.com/Zubayear/aragorn/pkg/models"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

type Repository interface {
	FetchUser(ctx context.Context, username string) (repository.GetUserRow, error)
	FetchUserFromCache(ctx context.Context, key string) (*models.User, error)
	SaveToCache(ctx context.Context, key string, value []byte) bool
	CliExists(ctx context.Context, key, cli string) (bool, error)
	FindByKey(ctx context.Context, key string) (int64, error)
}

type userRepository struct {
	pool        *pgxpool.Pool
	redisClient *redis.Client
	queries     *repository.Queries
}

func NewUserRepository(pool *pgxpool.Pool, redisClient *redis.Client) *userRepository {
	return &userRepository{pool: pool, redisClient: redisClient, queries: repository.New(pool)}
}

func (u *userRepository) FetchUser(ctx context.Context, username string) (repository.GetUserRow, error) {
	return u.queries.GetUser(ctx, username)
}

func (u *userRepository) FetchUserFromCache(ctx context.Context, key string) (*models.User, error) {
	entry, err := u.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, errors.New("user not found")
	}

	var userFromRedis models.User
	err = json.Unmarshal([]byte(entry), &userFromRedis)

	if err != nil {
		return nil, err
	}
	return &userFromRedis, err
}

func (u *userRepository) SaveToCache(ctx context.Context, key string, value []byte) bool {
	err := u.redisClient.Set(ctx, key, value, 12*time.Hour).Err()
	if err != nil {
		return false
	}
	return true
}

func (u *userRepository) CliExists(ctx context.Context, key, cli string) (bool, error) {
	return u.redisClient.SIsMember(ctx, key, cli).Result()
}

func (u *userRepository) FindByKey(ctx context.Context, key string) (int64, error) {
	return u.redisClient.Exists(ctx, key).Result()
}
