// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getUser = `-- name: GetUser :one
SELECT username,
       password,
       mt_port,
       status,
       acl_list ->> 'c11_campaign_api' AS c11_campaign_api, acl_list ->> 'c15_campaign_cpapi' AS c15_campaign_cpapi, acl_list ->> 'c12_campaign_sapi' AS c12_campaign_sapi, acl_list ->> 'c18_campaign_v2_api' AS campaign_api_v2, mid_expiry_time
FROM tbl_user
WHERE username=$1
`

type GetUserRow struct {
	Username         string        `json:"username"`
	Password         string        `json:"password"`
	MtPort           string        `json:"mt_port"`
	Status           int16         `json:"status"`
	C11CampaignApi   interface{}   `json:"c11_campaign_api"`
	C15CampaignCpapi interface{}   `json:"c15_campaign_cpapi"`
	C12CampaignSapi  interface{}   `json:"c12_campaign_sapi"`
	CampaignApiV2    interface{}   `json:"campaign_api_v2"`
	MidExpiryTime    pgtype.Float8 `json:"mid_expiry_time"`
}

func (q *Queries) GetUser(ctx context.Context, username string) (GetUserRow, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i GetUserRow
	err := row.Scan(
		&i.Username,
		&i.Password,
		&i.MtPort,
		&i.Status,
		&i.C11CampaignApi,
		&i.C15CampaignCpapi,
		&i.C12CampaignSapi,
		&i.CampaignApiV2,
		&i.MidExpiryTime,
	)
	return i, err
}