package sms

import "context"

type Service interface {
	ProcessNumbersCp(ctx context.Context, cli string, username string, result string, nums string, apiCode int, msgType int, unicode int, msgLen int, id string, content string, requestId string) string
}

type service struct {
	Repo Repository
}

// ProcessNumbersCp implements Service.
func (s *service) ProcessNumbersCp(ctx context.Context, cli string, username string, result string, nums string, apiCode int, msgType int, unicode int, msgLen int, id string, content string, requestId string) string {
	return s.Repo.CallCpSp(ctx, cli, username, result, nums, apiCode, msgType, unicode, msgLen, id, content, requestId)
}

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}
