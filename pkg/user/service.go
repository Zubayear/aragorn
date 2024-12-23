package user

import (
	"context"
	"github.com/Zubayear/aragorn/internal/repository"
	"github.com/Zubayear/aragorn/pkg/models"
)

type Service interface {
	GetUser(ctx context.Context, username string) (repository.GetUserRow, error)
	GetUserFromCache(ctx context.Context, key string) (*models.User, error)
	SaveToCache(ctx context.Context, key string, value []byte) bool
	KeyExists(ctx context.Context, key string) bool
}

type service struct {
	Repo Repository
}

// KeyExists implements Service.
func (s *service) KeyExists(ctx context.Context, key string) bool {

	entry, err := s.Repo.FindByKey(ctx, key)
	if err != nil {
		return false
	}
	return entry >= 1
}

func (s *service) SaveToCache(ctx context.Context, key string, value []byte) bool {
	return s.Repo.SaveToCache(ctx, key, value)
}

func (s *service) GetUserFromCache(ctx context.Context, key string) (*models.User, error) {
	return s.Repo.FetchUserFromCache(ctx, key)
}

func (s *service) GetUser(ctx context.Context, username string) (repository.GetUserRow, error) {
	return s.Repo.FetchUser(ctx, username)
}

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}
