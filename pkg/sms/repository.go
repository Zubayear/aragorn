package sms

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CallCpSp(ctx context.Context, cli string, username string, result string, nums string, apiCode int, msgType int, unicode int, msgLen int, id string, content string, requestId string) string
}

type repository struct {
	Pool *pgxpool.Pool
}

// CallCpSp implements Repository.
func (r *repository) CallCpSp(ctx context.Context, cli string, username string, result string, nums string, apiCode int, msgType int, unicode int, msgLen int, id string, content string, requestId string) string {
	ctx2, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	responseChan := make(chan string)

	go func() {
		row := r.Pool.QueryRow(ctx2,
			"CALL adareach.process_numbers_cp_api ($1, $2, $3, $4, $5, $6, $7, $8. $9, $10, $11)", cli, username, result, nums, apiCode, msgType, unicode, msgLen, id, content, requestId)
		var res string
		err := row.Scan(&res)
		if err != nil {
			responseChan <- ""
		}
		responseChan <- res
	}()

	for {
		select {
		case <-ctx2.Done():
			return ""
		case response := <-responseChan:
			return response
		}
	}
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{Pool: pool}
}
