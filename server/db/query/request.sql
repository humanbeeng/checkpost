-- name: CreateNewRequest :one
INSERT INTO
    request (
        user_id,
        endpoint_id,
        PATH,
        response_id,
        CONTENT,
        METHOD,
        source_ip,
        content_size,
        response_code,
        headers,
        query_params,
        expires_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING
    *;

-- name: GetEndpointHistory :many
SELECT
    *
FROM
    request
    LEFT JOIN endpoint ON request.endpoint_id = endpoint.id
WHERE
    endpoint.endpoint = $1
    AND request.is_deleted = FALSE
LIMIT
    $2
OFFSET
    $3;

-- name: GetRequestById :one
SELECT
    *
FROM
    request
WHERE
    id = $1
    AND is_deleted = FALSE
LIMIT
    1;

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
FROM
    request r
    LEFT JOIN endpoint e ON r.endpoint_id = e.id
WHERE
    endpoint = $1
    AND is_deleted = FALSE;