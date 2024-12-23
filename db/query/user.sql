-- name: GetUser :one
SELECT username,
       password,
       mt_port,
       status,
       acl_list ->> 'c11_campaign_api' AS c11_campaign_api, acl_list ->> 'c15_campaign_cpapi' AS c15_campaign_cpapi, acl_list ->> 'c12_campaign_sapi' AS c12_campaign_sapi, acl_list ->> 'c18_campaign_v2_api' AS campaign_api_v2, mid_expiry_time
FROM tbl_user
WHERE username=$1;