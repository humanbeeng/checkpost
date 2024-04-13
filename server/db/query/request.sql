-- name: CreateNewRequest :one
insert into
    request (
        user_id, endpoint_id, path, response_id, content, method, source_ip, content_size, response_code, headers, query_params, expires_at
    )
values (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
    )
returning
    *;

-- name: GetEndpointHistory :many
select *
from request
    left join endpoint on request.endpoint_id = endpoint.id
where
    endpoint.endpoint = $1
limit $2
offset
    $3;

-- name: GetRequestById :one
select * from request where id = $1 limit 1;

-- name: GetEndpointRequestCount :one
SELECT
    COUNT(*) AS total_count,
    COUNT(
        CASE
            WHEN response_code = 200 THEN 1
        END
    ) AS success_count,
    COUNT(
        CASE
            WHEN response_code != 200 THEN 1
        END
    ) AS failure_count
FROM request r
    LEFT JOIN endpoint e on r.endpoint_id = e.id
where
    endpoint = $1;