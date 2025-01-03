// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: sq.sql

package repository

import (
	"context"
)

const callProcessNumbersCpApiS = `-- name: CallProcessNumbersCpApiS :one
CALL process_numbers_cp_api ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

type CallProcessNumbersCpApiSParams struct {
	ProcessNumbersCpApi    interface{} `json:"process_numbers_cp_api"`
	ProcessNumbersCpApi_2  interface{} `json:"process_numbers_cp_api_2"`
	ProcessNumbersCpApi_3  interface{} `json:"process_numbers_cp_api_3"`
	ProcessNumbersCpApi_4  interface{} `json:"process_numbers_cp_api_4"`
	ProcessNumbersCpApi_5  interface{} `json:"process_numbers_cp_api_5"`
	ProcessNumbersCpApi_6  interface{} `json:"process_numbers_cp_api_6"`
	ProcessNumbersCpApi_7  interface{} `json:"process_numbers_cp_api_7"`
	ProcessNumbersCpApi_8  interface{} `json:"process_numbers_cp_api_8"`
	ProcessNumbersCpApi_9  interface{} `json:"process_numbers_cp_api_9"`
	ProcessNumbersCpApi_10 interface{} `json:"process_numbers_cp_api_10"`
}

type CallProcessNumbersCpApiSRow struct {
}

func (q *Queries) CallProcessNumbersCpApiS(ctx context.Context, arg CallProcessNumbersCpApiSParams) (CallProcessNumbersCpApiSRow, error) {
	row := q.db.QueryRow(ctx, callProcessNumbersCpApiS,
		arg.ProcessNumbersCpApi,
		arg.ProcessNumbersCpApi_2,
		arg.ProcessNumbersCpApi_3,
		arg.ProcessNumbersCpApi_4,
		arg.ProcessNumbersCpApi_5,
		arg.ProcessNumbersCpApi_6,
		arg.ProcessNumbersCpApi_7,
		arg.ProcessNumbersCpApi_8,
		arg.ProcessNumbersCpApi_9,
		arg.ProcessNumbersCpApi_10,
	)
	var i CallProcessNumbersCpApiSRow
	err := row.Scan()
	return i, err
}
