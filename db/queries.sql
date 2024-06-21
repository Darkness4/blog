-- name: FindPageViews :one
SELECT * FROM page_views WHERE page_id = $1 LIMIT 1;

-- name: FindPageViewsByPageId :many
SELECT * FROM page_views WHERE page_id = ANY(sqlc.arg(page_ids)::string[]);

-- name: CreatePageViewsIPs :many
INSERT INTO page_views_ips (page_id, ip)
VALUES ($1, $2)
ON CONFLICT
DO NOTHING
RETURNING *;

-- name: CreateOrIncrementPageViews :one
INSERT INTO page_views (page_id, views)
VALUES ($1, 1)
ON CONFLICT (page_id)
DO UPDATE SET views = page_views.views + 1 RETURNING *;

-- name: DeletePageViews :exec
DELETE FROM page_views;

-- name: DeletePageViewsIPs :exec
DELETE FROM page_views_ips;
